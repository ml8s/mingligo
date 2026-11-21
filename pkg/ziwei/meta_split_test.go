package ziwei

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// TestMetaSplitRelocation verifies the MetaView split moved every field to its
// owning subdomain without changing any value.
func TestMetaSplitRelocation(t *testing.T) {
	rawBase, err := os.ReadFile("testdata/meta_split_baseline.json")
	if err != nil {
		t.Fatal(err)
	}
	var base map[string]any
	if err := json.Unmarshal(rawBase, &base); err != nil {
		t.Fatal(err)
	}
	standard, _ := time.Parse(time.RFC3339, "1984-02-15T00:00:00Z")
	chart, err := ComputeChartAt(standard, 116.4, ganzhi.Male)
	if err != nil {
		t.Fatal(err)
	}
	gotRaw, _ := json.Marshal(chart.View())
	var got map[string]any
	if err := json.Unmarshal(gotRaw, &got); err != nil {
		t.Fatal(err)
	}
	bm, ok := base["meta"].(map[string]any)
	if !ok {
		t.Fatal("baseline meta not object")
	}
	moves := []struct{ key, to string }{
		{"kong_gong", "palace_relations"}, {"san_fang_si_zheng", "palace_relations"},
		{"ge_ju", "patterns"}, {"da_xian", "cycles"}, {"ming_zhu", "soul_body"}, {"shen_zhu", "soul_body"},
	}
	for _, m := range moves {
		ov, had := bm[m.key]
		target, ok := got[m.to].(map[string]any)
		if !ok {
			if had {
				t.Errorf("%s missing", m.to)
			}
			continue
		}
		nv, has := target[m.key]
		if had != has {
			t.Errorf("meta.%s presence changed", m.key)
			continue
		}
		if had && !reflect.DeepEqual(ov, nv) {
			t.Errorf("meta.%s -> %s.%s: %v -> %v", m.key, m.to, m.key, ov, nv)
		}
	}
	for key, ov := range base {
		if key == "meta" {
			continue
		}
		if !reflect.DeepEqual(ov, got[key]) {
			t.Errorf("top-level %s changed", key)
		}
	}
}
