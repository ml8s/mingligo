package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ml8s/mingligo/examples/mingligo-api/internal/apicontract"
)

func fixtureTool(name string) string {
	domain, _, ok := strings.Cut(name, "_")
	if !ok {
		return ""
	}
	if domain == "qimen" {
		return "qimen_chart"
	}
	rest := name[len(domain)+1:]
	tokens := strings.Split(rest, "_")
	for i := len(tokens) - 1; i >= 0; i-- {
		switch tokens[i] {
		case "liunian", "liuyue", "liuri", "liushi", "xiaoyun", "bond", "chart":
			switch tokens[i] {
			case "liunian", "liuyue", "liuri", "liushi", "xiaoyun":
				return domain + "_liupan"
			default:
				return domain + "_" + tokens[i]
			}
		}
	}
	return ""
}

func TestOpenAPIMatchesCode(t *testing.T) {
	raw, err := os.ReadFile("openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Paths map[string]map[string]struct {
			RequestBody struct {
				Content map[string]struct {
					Schema json.RawMessage `json:"schema"`
				} `json:"content"`
			} `json:"requestBody"`
			Responses map[string]struct {
				Content map[string]struct {
					Schema json.RawMessage `json:"schema"`
				} `json:"content"`
			} `json:"responses"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	schemas, err := apicontract.MethodSchemas()
	if err != nil {
		t.Fatal(err)
	}
	for _, ms := range schemas {
		p, ok := doc.Paths["/mingligo/api/"+ms.Name]
		if !ok {
			t.Errorf("openapi missing /mingligo/api/%s", ms.Name)
			continue
		}
		in := p["post"].RequestBody.Content["application/json"].Schema
		out := p["post"].Responses["200"].Content["application/json"].Schema
		if !bytes.Equal(compact(in), compact(ms.Input)) {
			t.Errorf("%s input schema drift: openapi vs code", ms.Name)
		}
		if !bytes.Equal(compact(out), compact(ms.Output)) {
			t.Errorf("%s output schema drift: openapi vs code", ms.Name)
		}
	}
}

func compact(b json.RawMessage) []byte {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return b
	}
	out, _ := json.Marshal(v)
	return out
}

func TestAllOutputSchemasResolveForValidation(t *testing.T) {
	validators, err := newOutputValidators()
	if err != nil {
		t.Fatal(err)
	}
	names := apicontract.MethodNames()
	if len(validators) != len(names) {
		t.Fatalf("validator count = %d, want %d", len(validators), len(names))
	}
	for _, name := range names {
		if validators[name] == nil {
			t.Errorf("%s has no resolved output validator", name)
		}
		if err := validateOutput(validators, name, []byte(`{"__unexpected__":true}`)); err == nil {
			t.Errorf("%s accepted an output outside its contract", name)
		}
	}
}

func TestAllRoutesHTTPContract(t *testing.T) {
	openapi, err := os.ReadFile("openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("BASE_PATH", "/mingligo/api")
	server, basePath, err := newServer(openapi)
	if err != nil {
		t.Fatal(err)
	}
	if basePath != "/mingligo/api" {
		t.Fatalf("base path = %q", basePath)
	}
	srv := httptest.NewServer(server)
	defer srv.Close()

	dir := filepath.Join("testdata", "wire_inputs")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	covered := map[string]bool{}
	// Methods without wire fixtures: exercise with inline representative inputs.
	inline := map[string]string{
		"time_now":           `{"timezone":"Asia/Shanghai"}`,
		"subject_candidates": `{"calendar_kind":"solar","recorded_date":"1990-05-20","record_timezone":"Asia/Shanghai","recorded_time":"11:49","longitude":116.4,"subject_gender":"male"}`,
		"city_info":          `{"city":"北京"}`,
		"huangli_days":       `{"standard_time":"2024-06-15T12:00:00+08:00","count":1,"event":"sign"}`,
		"bazhai_chart":       `{"subject":{"standard_time":"1984-06-15T12:00:00+08:00","gender":"male"}}`,
		"bazhai_layout":      `{"ming_gua":"乾","men_gua":"坤","master_gua":"艮","stove_gua":"兑"}`,
		"xuankong_chart":     `{"standard_time":"2026-07-31T12:00:00+08:00","zuo_shan":2,"xiang_shan":14}`,
		"xuankong_liunian":   `{"standard_time":"2026-06-15T12:00:00+08:00","zuo_shan":2,"xiang_shan":14}`,
		"liuyao_qigua":       `{"origin":{"fang_shi":"yaos","yao_zhi":[7,8,7,8,7,8]},"qi_gua_qing_jing":{"shi_jian":"2026-06-15T12:00:00+08:00","jing_du":116.4,"focus":{"kind":"qu_yong","qu_yong":"官鬼"}}}`,
	}
	// liuyao_chart consumes the qigua event produced above.
	for name, in := range inline {
		body, code := callHTTP(t, srv, name, json.RawMessage(in))
		if code != http.StatusOK {
			t.Fatalf("%s inline: status %d body %s", name, code, body)
		}
		covered[name] = true
	}
	// liuyao_chart consumes the qigua event produced by liuyao_qigua (chained).
	qiguaOut, code := callHTTP(t, srv, "liuyao_qigua", json.RawMessage(inline["liuyao_qigua"]))
	if code != http.StatusOK {
		t.Fatalf("qigua chain: status %d", code)
	}
	var qg struct {
		QiGua json.RawMessage `json:"qi_gua"`
	}
	if err := json.Unmarshal(qiguaOut, &qg); err != nil {
		t.Fatal(err)
	}
	chartIn, _ := json.Marshal(map[string]any{"qi_gua": qg.QiGua})
	body, code := callHTTP(t, srv, "liuyao_chart", chartIn)
	if code != http.StatusOK {
		t.Fatalf("liuyao_chart chain: status %d body %s", code, body)
	}
	covered["liuyao_chart"] = true

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		method := fixtureTool(strings.TrimSuffix(e.Name(), ".json"))
		if method == "" {
			t.Fatalf("cannot map fixture %s", e.Name())
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		var f struct {
			Input json.RawMessage `json:"input"`
		}
		if err := json.Unmarshal(raw, &f); err != nil {
			t.Fatal(err)
		}
		body, code := callHTTP(t, srv, method, f.Input)
		if code != http.StatusOK {
			t.Fatalf("%s (%s): status %d body %s", method, e.Name(), code, body)
		}
		if len(body) == 0 || string(body) == "null" {
			t.Fatalf("%s: empty output", method)
		}
		body2, code2 := callHTTP(t, srv, method, f.Input)
		if code2 != http.StatusOK || !bytes.Equal(body, body2) {
			t.Errorf("%s: non-deterministic output", method)
		}
		covered[method] = true
	}
	// every published method must have at least one exercised fixture
	for _, name := range apicontract.MethodNames() {
		if !covered[name] {
			t.Errorf("%s has no exercised fixture", name)
		}
	}
}

func callHTTP(t *testing.T, srv *httptest.Server, method string, input json.RawMessage) ([]byte, int) {
	t.Helper()
	resp, err := http.Post(srv.URL+"/mingligo/api/"+method, "application/json", bytes.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(resp.Body)
	return buf.Bytes(), resp.StatusCode
}
