package mcpcontract

import (
	"encoding/json"
	"os"
	"slices"
	"testing"
)

func TestMCPCatalogMatchesRegisteredTools(t *testing.T) {
	raw, err := os.ReadFile("../../../../contracts/mcp-tools.json")
	if err != nil {
		t.Fatal(err)
	}
	var catalog struct {
		Tools []string `json:"tools"`
	}
	if err := json.Unmarshal(raw, &catalog); err != nil {
		t.Fatal(err)
	}
	registered := ToolNames()
	slices.Sort(catalog.Tools)
	if !slices.Equal(catalog.Tools, registered) {
		t.Fatalf("neutral catalog drift: catalog=%v registered=%v", catalog.Tools, registered)
	}
}
