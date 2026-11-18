package mcpcontract

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Definition mirrors the public metadata exposed by tools/list.
type Definition struct {
	Name         string
	Description  string
	InputSchema  json.RawMessage
	OutputSchema json.RawMessage
}

// ContractServer is a read-only test view over one official SDK server session. It
// lets focused contract tests exercise the same validation path as production.
type ContractServer struct {
	session     *mcp.ClientSession
	definitions map[string]Definition
	names       []string
}

var testEngine = sync.OnceValue(newContractServer)

func newContractServer() *ContractServer {
	logger := slog.New(slog.NewTextHandler(&discardWriter{}, nil))
	server := mcp.NewServer(&mcp.Implementation{Name: "test/mingligo-mcp", Version: "test"}, &mcp.ServerOptions{Logger: logger})
	if err := RegisterTools(server); err != nil {
		panic(err)
	}
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	if _, err := server.Connect(context.Background(), serverTransport, nil); err != nil {
		panic(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1"}, &mcp.ClientOptions{Logger: logger})
	session, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		panic(err)
	}

	contractServer := &ContractServer{session: session, definitions: make(map[string]Definition)}
	for tool, err := range session.Tools(context.Background(), nil) {
		if err != nil {
			panic(err)
		}
		inputRaw, err := json.Marshal(tool.InputSchema)
		if err != nil {
			panic(fmt.Sprintf("marshal %s input schema: %v", tool.Name, err))
		}
		outputRaw, err := json.Marshal(tool.OutputSchema)
		if err != nil {
			panic(fmt.Sprintf("marshal %s output schema: %v", tool.Name, err))
		}
		contractServer.names = append(contractServer.names, tool.Name)
		contractServer.definitions[tool.Name] = Definition{
			Name:         tool.Name,
			Description:  tool.Description,
			InputSchema:  inputRaw,
			OutputSchema: outputRaw,
		}
	}
	return contractServer
}

type discardWriter struct{}

func (discardWriter) Write(_ []byte) (int, error) { return 0, nil }

func (r *ContractServer) Names() []string { return append([]string(nil), r.names...) }

func (r *ContractServer) MaterializeSchemas() error {
	for _, name := range r.names {
		if _, err := r.Definition(name); err != nil {
			return err
		}
	}
	return nil
}

func (r *ContractServer) Definition(name string) (Definition, error) {
	definition, ok := r.definitions[name]
	if !ok {
		return Definition{}, fmt.Errorf("unknown tool: %s", name)
	}
	return definition, nil
}

func (r *ContractServer) Call(_ context.Context, name string, arguments json.RawMessage) (json.RawMessage, error) {
	if _, ok := r.definitions[name]; !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	var args any
	if len(arguments) > 0 {
		if err := json.Unmarshal(arguments, &args); err != nil {
			return nil, fmt.Errorf("decode arguments: %w", err)
		}
	}
	result, err := r.session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		return nil, err
	}
	if result.IsError {
		var detail []string
		for _, content := range result.Content {
			if text, ok := content.(*mcp.TextContent); ok {
				detail = append(detail, text.Text)
			}
		}
		return nil, fmt.Errorf("tool error: %s", strings.Join(detail, "; "))
	}
	if result.StructuredContent == nil {
		return nil, fmt.Errorf("tool %s returned no structured content", name)
	}
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (r *ContractServer) ValidateOutput(name string, output json.RawMessage) error {
	definition, err := r.Definition(name)
	if err != nil {
		return err
	}
	var schema jsonschema.Schema
	if err := json.Unmarshal(definition.OutputSchema, &schema); err != nil {
		return err
	}
	resolved, err := schema.Resolve(&jsonschema.ResolveOptions{ValidateDefaults: true})
	if err != nil {
		return err
	}
	var document any
	if err := json.Unmarshal(output, &document); err != nil {
		return err
	}
	if err := resolved.Validate(document); err != nil {
		return fmt.Errorf("output contract violation: %w", err)
	}
	return nil
}

func decode(t *testing.T, raw json.RawMessage, target any) {
	t.Helper()
	if err := json.Unmarshal(raw, target); err != nil {
		t.Fatalf("decode JSON %s: %v", raw, err)
	}
}
