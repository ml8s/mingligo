package bazi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

var updateBaziSnapshot = os.Getenv("UPDATE_BAZI_SNAPSHOT") == "1"

// TestBaziChartFullSnapshot guards the complete ComputeChart output as a byte
// snapshot. The older golden files only assert selected anchors; this catches
// field-level and serialization drift across all layers (神煞/十神/合冲刑害/
// 用神证据/纳音/三元/大运等).
func TestBaziChartFullSnapshot(t *testing.T) {
	for _, g := range goldenCharts {
		t.Run(g.Name, func(t *testing.T) {
			chart := computeChart(g)
			got, err := json.MarshalIndent(chart, "", "  ")
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			path := filepath.Join("testdata", "full_"+g.Name+".json")
			if updateBaziSnapshot {
				if err := os.WriteFile(path, append(got, '\n'), 0644); err != nil {
					t.Fatalf("write snapshot: %v", err)
				}
				return
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read snapshot %s: %v — run with UPDATE_BAZI_SNAPSHOT=1", path, err)
			}
			if string(got)+"\n" != string(want) {
				t.Errorf("chart output differs from snapshot %s", path)
			}
		})
	}
}
