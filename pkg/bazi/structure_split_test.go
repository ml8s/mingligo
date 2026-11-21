package bazi

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// TestStructureSplitFieldRelocation verifies the StructureView split moved
// fields to their owning subdomains without changing any value: conclusions
// -> ConclusionsView, 口径/birth_year -> CalendarView, structural facts stay.
func TestStructureSplitFieldRelocation(t *testing.T) {
	rawBase, err := os.ReadFile("testdata/structure_split_baseline.json")
	if err != nil {
		t.Fatal(err)
	}
	var base map[string]any
	if err := json.Unmarshal(rawBase, &base); err != nil {
		t.Fatal(err)
	}
	baseD := base
	if d, ok := base["data"].(map[string]any); ok {
		baseD = d
	}

	standard, _ := time.Parse(time.RFC3339, "1990-05-20T03:49:00Z")
	c := ComputeCoreAt(standard, 116.4, ganzhi.Male)
	v := ComputeChartFromCore(c).View()
	gotRaw, _ := json.Marshal(v)
	var got map[string]any
	if err := json.Unmarshal(gotRaw, &got); err != nil {
		t.Fatal(err)
	}

	// 迁移字段值不变
	checks := []struct{ from, key, to string }{
		{"structure", "guan_sha", "conclusions"}, {"structure", "cai_ku", "conclusions"},
		{"structure", "spouse_palace_state", "conclusions"}, {"structure", "day_branch_type", "conclusions"},
		{"structure", "shen_sha_school", "calendar"}, {"", "birth_year", "calendar"},
	}
	moved := map[string]bool{"guan_sha": true, "cai_ku": true, "spouse_palace_state": true, "day_branch_type": true, "shen_sha_school": true}
	for _, ck := range checks {
		var ov, nv any
		if ck.from == "" {
			ov = baseD[ck.key]
		} else {
			ov = baseD[ck.from].(map[string]any)[ck.key]
		}
		nv = got[ck.to].(map[string]any)[ck.key]
		if !reflect.DeepEqual(ov, nv) {
			t.Errorf("%s.%s -> %s.%s: %v -> %v", ck.from, ck.key, ck.to, ck.key, ov, nv)
		}
	}
	// structure 保留字段
	for key, ov := range baseD["structure"].(map[string]any) {
		if moved[key] {
			continue
		}
		if !reflect.DeepEqual(ov, got["structure"].(map[string]any)[key]) {
			t.Errorf("structure.%s changed", key)
		}
	}
}
