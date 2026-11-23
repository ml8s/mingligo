// Command gen_openapi regenerates examples/mingligo-api/openapi.json from the
// code's input DTOs and domain Views. The document is generated, never
// hand-written, so it cannot drift from the implementation.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ml8s/mingligo/examples/mingligo-api/internal/apicontract"
)

type openAPIDoc struct {
	OpenAPI    string                    `json:"openapi"`
	Info       map[string]string         `json:"info"`
	Servers    []map[string]string       `json:"servers"`
	Paths      map[string]map[string]any `json:"paths"`
	Components map[string]map[string]any `json:"components"`
}

func main() {
	schemas, err := apicontract.MethodSchemas()
	if err != nil {
		panic(err)
	}
	doc := openAPIDoc{
		OpenAPI: "3.1.0",
		Info: map[string]string{
			"title":       "MingliGo Domain Facts API",
			"version":     "1.0.0",
			"description": "Generated from domain Views. The API layer only transports domain facts; no 命理 logic.",
		},
		Servers: []map[string]string{{"url": "https://liki.tokflux.com"}},
		Paths:   make(map[string]map[string]any),
		Components: map[string]map[string]any{
			"schemas": {
				"Problem": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"type":       map[string]any{"type": "string"},
						"title":      map[string]any{"type": "string"},
						"status":     map[string]any{"type": "integer"},
						"detail":     map[string]any{"type": "string"},
						"instance":   map[string]any{"type": "string"},
						"request_id": map[string]any{"type": "string"},
					},
					"required": []string{"type", "title", "status"},
				},
			},
		},
	}
	for _, ms := range schemas {
		var input, output any
		if err := json.Unmarshal(ms.Input, &input); err != nil {
			panic(err)
		}
		if err := json.Unmarshal(ms.Output, &output); err != nil {
			panic(err)
		}
		doc.Paths["/mingligo/api/"+ms.Name] = map[string]any{
			"post": map[string]any{
				"operationId": ms.Name,
				"summary":     "Execute " + ms.Name + " (domain fact computation)",
				"requestBody": map[string]any{
					"required": true,
					"content":  map[string]any{"application/json": map[string]any{"schema": input}},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Domain View",
						"content":     map[string]any{"application/json": map[string]any{"schema": output}},
					},
					"400": problemRef("Invalid JSON"),
					"415": problemRef("Unsupported media type"),
					"422": problemRef("Invalid input"),
					"500": problemRef("Internal error"),
				},
			},
		}
	}
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	out := filepath.Join("examples", "mingligo-api", "openapi.json")
	if err := os.WriteFile(out, append(raw, '\n'), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("generated %s (%d methods)\n", out, len(schemas))
}

func problemRef(title string) map[string]any {
	return map[string]any{
		"description": title,
		"content": map[string]any{
			"application/problem+json": map[string]any{
				"schema": map[string]any{"$ref": "#/components/schemas/Problem"},
			},
		},
	}
}
