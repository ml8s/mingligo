package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/ml8s/mingligo/pkg/tables"

	"github.com/ml8s/mingligo/examples/mingligo-api/internal/apicontract"
)

// route binds one MCP-equivalent method name to a typed pure handler.
// The HTTP layer only decodes the body, calls the domain handler, and writes
// the returned domain View back — no 命理 logic, no key rewriting.
type route struct {
	method  string
	handler func(context.Context, json.RawMessage) (any, error)
}

func typedRoute[In any, Out any](
	name string,
	fn func(context.Context, In) (Out, error),
) route {
	return route{
		method: name,
		handler: func(ctx context.Context, raw json.RawMessage) (any, error) {
			var in In
			if err := json.Unmarshal(raw, &in); err != nil {
				return nil, fmt.Errorf("invalid input: %w", err)
			}
			return fn(ctx, in)
		},
	}
}

func routes() []route {
	return []route{
		typedRoute("time_now", apicontract.TimeNow),
		typedRoute("subject_candidates", apicontract.TianwenSubjectCandidates),
		typedRoute("city_info", apicontract.CityInfo),
		typedRoute("huangli_days", apicontract.HuangliDays),
		typedRoute("bazi_chart", apicontract.BaziChart),
		typedRoute("bazi_bond", apicontract.BaziBond),
		typedRoute("bazi_liupan", apicontract.BaziLiupan),
		typedRoute("bazhai_chart", apicontract.BazhaiChart),
		typedRoute("bazhai_layout", apicontract.BazhaiLayout),
		typedRoute("xuankong_chart", apicontract.XuankongChartCompute),
		typedRoute("xuankong_liunian", apicontract.XuankongLiunian),
		typedRoute("liuyao_qigua", apicontract.LiuyaoQigua),
		typedRoute("liuyao_chart", apicontract.LiuyaoChart),
		typedRoute("ziwei_chart", apicontract.ZiweiChart),
		typedRoute("ziwei_liupan", apicontract.ZiweiLiupan),
		typedRoute("ziwei_bond", apicontract.ZiweiBond),
		typedRoute("qimen_chart", apicontract.QimenChart),
	}
}

var tablesOnce sync.Once
var tablesErr error

func ensureTables() error {
	tablesOnce.Do(func() { tablesErr = tables.Load() })
	return tablesErr
}

// outputValidators resolves every published response schema once. Runtime
// validation turns a handler/type drift into an HTTP 500 instead of silently
// publishing a contract-breaking View.
func newOutputValidators() (map[string]*jsonschema.Resolved, error) {
	methods, err := apicontract.MethodSchemas()
	if err != nil {
		return nil, err
	}
	validators := make(map[string]*jsonschema.Resolved, len(methods))
	for _, method := range methods {
		var schema jsonschema.Schema
		if err := json.Unmarshal(method.Output, &schema); err != nil {
			return nil, fmt.Errorf("decode %s output schema: %w", method.Name, err)
		}
		resolved, err := schema.Resolve(nil)
		if err != nil {
			return nil, fmt.Errorf("resolve %s output schema: %w", method.Name, err)
		}
		validators[method.Name] = resolved
	}
	return validators, nil
}

func validateOutput(validators map[string]*jsonschema.Resolved, method string, body []byte) error {
	resolved, ok := validators[method]
	if !ok {
		return fmt.Errorf("output schema missing for %s", method)
	}
	var instance any
	if err := json.Unmarshal(body, &instance); err != nil {
		return fmt.Errorf("decode %s output: %w", method, err)
	}
	if err := resolved.Validate(instance); err != nil {
		return fmt.Errorf("validate %s output: %w", method, err)
	}
	return nil
}

const defaultBasePath = "/mingligo/api"

type problem struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail,omitempty"`
	Instance  string `json:"instance,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func normalizeBasePath(value string) string {
	if value == "" {
		return defaultBasePath
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	return strings.TrimSuffix(value, "/")
}

func newServer(openapi []byte) (http.Handler, string, error) {
	basePath := normalizeBasePath(os.Getenv("BASE_PATH"))
	if err := ensureTables(); err != nil {
		return nil, basePath, fmt.Errorf("load domain tables: %w", err)
	}
	validators, err := newOutputValidators()
	if err != nil {
		return nil, basePath, fmt.Errorf("load output schemas: %w", err)
	}
	mux := http.NewServeMux()
	knownMethods := make(map[string]bool)
	for _, r := range routes() {
		rt := r // capture
		path := basePath + "/" + rt.method
		knownMethods[path] = true
		mux.HandleFunc("POST "+path, func(w http.ResponseWriter, req *http.Request) {
			if !hasJSONContentType(req.Header.Get("Content-Type")) {
				writeProblem(w, basePath+"/"+rt.method, http.StatusUnsupportedMediaType,
					"urn:mingligo:problem:unsupported-media-type",
					"Content-Type must be application/json", requestID(req))
				return
			}
			req.Body = http.MaxBytesReader(w, req.Body, 2<<20)
			var raw json.RawMessage
			dec := json.NewDecoder(req.Body)
			if err := dec.Decode(&raw); err != nil {
				writeProblem(w, req.URL.Path, http.StatusBadRequest,
					"urn:mingligo:problem:invalid-json",
					"request body is not valid JSON", requestID(req))
				return
			}
			out, err := rt.handler(req.Context(), raw)
			if err != nil {
				writeProblem(w, req.URL.Path, http.StatusUnprocessableEntity,
					"urn:mingligo:problem:invalid-input",
					err.Error(), requestID(req))
				return
			}
			body, err := json.Marshal(out)
			if err != nil {
				writeInternalProblem(w, req, fmt.Errorf("marshal %s output: %w", rt.method, err))
				return
			}
			if err := validateOutput(validators, rt.method, body); err != nil {
				writeInternalProblem(w, req, fmt.Errorf("validate %s output: %w", rt.method, err))
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(body)
		})
	}
	mux.HandleFunc("GET "+basePath+"/openapi.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(openapi)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	// Exact routes above win over this fallback. Here we turn protocol misses
	// into the public Problem Details contract instead of ServeMux's plain text.
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch {
		case req.URL.Path == basePath || strings.HasPrefix(req.URL.Path, basePath+"/"):
			if knownMethods[req.URL.Path] {
				writeProblem(w, req.URL.Path, http.StatusMethodNotAllowed,
					"urn:mingligo:problem:method-not-allowed",
					"domain methods require POST", requestID(req))
				return
			}
			writeProblem(w, req.URL.Path, http.StatusNotFound,
				"urn:mingligo:problem:method-not-found",
				"unknown domain method", requestID(req))
		default:
			writeProblem(w, req.URL.Path, http.StatusNotFound,
				"urn:mingligo:problem:not-found", "not found", requestID(req))
		}
	})
	return recoverInternal(basePath, mux), basePath, nil
}

func hasJSONContentType(value string) bool {
	mediaType, _, err := mime.ParseMediaType(value)
	return err == nil && mediaType == "application/json"
}

func requestID(req *http.Request) string {
	id := strings.TrimSpace(req.Header.Get("X-Request-Id"))
	if id != "" && len(id) <= 128 {
		return id
	}
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return ""
	}
	return hex.EncodeToString(b[:])
}

func writeProblem(w http.ResponseWriter, instance string, status int, kind, detail, id string) {
	problem := problem{
		Type: kind, Title: problemTitle(kind), Status: status,
		Detail: detail, Instance: instance, RequestID: id,
	}
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.Header().Set("X-Request-Id", id)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem)
}

func writeInternalProblem(w http.ResponseWriter, req *http.Request, err error) {
	log.Printf("request_id=%s path=%s internal_error=%v", requestID(req), req.URL.Path, err)
	writeProblem(w, req.URL.Path, http.StatusInternalServerError,
		"urn:mingligo:problem:internal-error", "internal error", requestID(req))
}

func problemTitle(kind string) string {
	switch kind {
	case "urn:mingligo:problem:invalid-json":
		return "Invalid JSON"
	case "urn:mingligo:problem:unsupported-media-type":
		return "Unsupported media type"
	case "urn:mingligo:problem:method-not-found":
		return "Method not found"
	case "urn:mingligo:problem:method-not-allowed":
		return "Method not allowed"
	case "urn:mingligo:problem:invalid-input":
		return "Invalid input"
	case "urn:mingligo:problem:internal-error":
		return "Internal error"
	default:
		return "Request failed"
	}
}

// recoverInternal ensures an adapter panic becomes a sanitized Problem Details
// response rather than a dropped connection.
func recoverInternal(basePath string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := requestID(req)
		req.Header.Set("X-Request-Id", id)
		w.Header().Set("X-Request-Id", id)
		defer func() {
			if recovered := recover(); recovered != nil {
				writeInternalProblem(w, req, fmt.Errorf("panic: %v", recovered))
			}
		}()
		next.ServeHTTP(w, req)
	})
}

func main() {
	healthcheck := len(os.Args) == 2 && os.Args[1] == "-healthcheck"
	if healthcheck {
		addr := os.Getenv("ADDR")
		if addr == "" {
			addr = ":18082"
		}
		url := "http://127.0.0.1:" + strings.TrimPrefix(addr, ":") + "/healthz"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("health endpoint returned %s", resp.Status)
		}
		return
	}
	openapi, err := loadOpenAPI()
	if err != nil {
		log.Fatalf("read openapi.json (run gen_openapi first): %v", err)
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":18082"
	}
	server, basePath, err := newServer(bytes.Clone(openapi))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("mingligo-api listening on %s%s (%d routes)", addr, basePath, len(routes()))
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatal(err)
	}
}

func loadOpenAPI() ([]byte, error) {
	if path := os.Getenv("MINGLIGO_OPENAPI_PATH"); path != "" {
		if data, err := os.ReadFile(path); err == nil {
			return data, nil
		}
	}
	data, err := os.ReadFile("openapi.json")
	if err == nil {
		return data, nil
	}
	// deployed layout: openapi.json sits next to the binary
	if exe, exeErr := os.Executable(); exeErr == nil {
		if data, err := os.ReadFile(filepath.Join(filepath.Dir(exe), "openapi.json")); err == nil {
			return data, nil
		}
	}
	// source-tree layout when run with `go run` from the repo root
	return os.ReadFile(filepath.Join("examples", "mingligo-api", "openapi.json"))
}
