package pkg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFactCoreDoesNotImportAdapters keeps MingliGo independently consumable as
// a Go library. MCP transport and product-facing adapters stay under examples.
func TestFactCoreDoesNotImportAdapters(t *testing.T) {
	root, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		source, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, line := range strings.Split(string(source), "\n") {
			examplesPrefix := "github.com/ml8s/mingligo/" + "examples"
			internalPrefix := "github.com/ml8s/mingligo/" + "internal"
			if strings.Contains(line, examplesPrefix) ||
				strings.Contains(line, internalPrefix) {
				t.Fatalf("%s imports an adapter package: %s", path, strings.TrimSpace(line))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
