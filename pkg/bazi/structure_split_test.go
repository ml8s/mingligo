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

	// 迁移字段值不变（最终命理域路径）
	checks := []struct{ from, key, to, toKey string }{
		{"structure", "guan_sha", "", "guan_sha"}, {"structure", "cai_ku", "", "cai_ku"},
		{"structure", "spouse_palace_state", "", "pei_ou_gong"}, {"structure", "day_branch_type", "", "ri_zhi_xiang"},
		{"structure", "shen_sha_school", "li_fa_kou_jing", "shen_sha_school"}, {"", "birth_year", "li_fa_kou_jing", "birth_year"},
		{"structure", "san_yuan", "san_yuan_view", "san_yuan"}, {"structure", "tai_xi", "shen_sha", "tai_xi"},
		{"structure", "chang_sheng", "chang_sheng_view", "chang_sheng"}, {"structure", "san_qi_name", "shen_sha", "san_qi_name"},
	}
	for _, ck := range checks {
		var ov, nv any
		if ck.from == "" {
			ov = baseD[ck.key]
		} else {
			ov = baseD[ck.from].(map[string]any)[ck.key]
		}
		if ck.to == "" {
			nv = got[ck.toKey]
		} else {
			nv = got[ck.to].(map[string]any)[ck.toKey]
		}
		if !reflect.DeepEqual(ov, nv) {
			t.Errorf("%s.%s -> %s.%s: %v -> %v", ck.from, ck.key, ck.to, ck.toKey, ov, nv)
		}
	}
}
