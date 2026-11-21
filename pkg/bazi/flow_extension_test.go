package bazi

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

func flowP1Core(t *testing.T) Core {
	t.Helper()
	standard, err := time.Parse(time.RFC3339, "1984-02-15T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	return ComputeCoreAt(standard, 116.4, ganzhi.Male)
}

func flowP1TestResults(t *testing.T) (*LiuRi, *LiuShi) {
	t.Helper()
	core := flowP1Core(t)
	standard, _ := time.Parse(time.RFC3339, "2024-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	ri, err := ComputeLiuRiAt(core, standard, loc, 116.4)
	if err != nil {
		t.Fatal(err)
	}
	shi, err := ComputeLiuShiAt(core, standard, loc, 116.4)
	if err != nil {
		t.Fatal(err)
	}
	return ri, shi
}

// TestFlowP1LegacyFieldsUnchanged is the differential-invariance gate: every
// pre-existing field in the baseline must be byte-identical; only the two
// newly added fields may appear.
func TestFlowP1LegacyFieldsUnchanged(t *testing.T) {
	rawBase, err := os.ReadFile("testdata/flow_p1_baseline.json")
	if err != nil {
		t.Fatal(err)
	}
	var base map[string]map[string]any
	if err := json.Unmarshal(rawBase, &base); err != nil {
		t.Fatal(err)
	}
	ri, shi := flowP1TestResults(t)
	got := map[string]any{"ri": ri, "shi": shi}
	gotRaw, _ := json.Marshal(got)
	var gotDoc map[string]map[string]any
	if err := json.Unmarshal(gotRaw, &gotDoc); err != nil {
		t.Fatal(err)
	}
	allowedNew := map[string]bool{
		"sheng_ri": true, "ke_ri": true,
		"da_yun_interaction": true, "liu_nian_interaction": true,
	}
	removed := map[string]bool{"da_yun_rels": true, "liu_nian_rels": true}
	for layer, baseFields := range base {
		for key, oldVal := range baseFields {
			if removed[key] {
				continue // P2 intentionally replaced zhi-only arrays with two-layer objects
			}
			newVal, ok := gotDoc[layer][key]
			if !ok {
				t.Errorf("%s.%s disappeared", layer, key)
				continue
			}
			if !reflect.DeepEqual(oldVal, newVal) {
				t.Errorf("%s.%s changed: %v -> %v", layer, key, oldVal, newVal)
			}
		}
		for key := range gotDoc[layer] {
			if _, existed := baseFields[key]; !existed && !allowedNew[key] {
				t.Errorf("%s: unexpected new field %s", layer, key)
			}
		}
	}
}

// TestFlowP1NewFieldsOracle pins the new values against the pure function,
// so the wiring (not just the function) is verified.
func TestFlowP1NewFieldsOracle(t *testing.T) {
	ri, shi := flowP1TestResults(t)
	dm := flowP1Core(t).ToBazi().Ri.Gan
	for name, tc := range map[string]struct {
		flow ganzhi.Gan
		got  int
	}{
		"ri":  {ri.RiGan, ri.ShengRi},
		"shi": {shi.ShiGan, shi.ShengRi},
	} {
		wantSheng, _, _ := flowDayMasterFacts(tc.flow, dm)
		if tc.got != wantSheng {
			t.Errorf("%s sheng_ri = %d, want %d", name, tc.got, wantSheng)
		}
	}
	if _, wantKe, _ := flowDayMasterFacts(ri.RiGan, dm); ri.KeRi != wantKe {
		t.Errorf("ri ke_ri = %d, want %d", ri.KeRi, wantKe)
	}
	if _, wantKe, _ := flowDayMasterFacts(shi.ShiGan, dm); shi.KeRi != wantKe {
		t.Errorf("shi ke_ri = %d, want %d", shi.KeRi, wantKe)
	}
}

// TestFlowP2LegacyVsNew is the differential gate for P2: the only allowed
// changes are the two old zhi-only arrays being replaced by two-layer objects.
// Every other field must remain byte-identical to the P2 baseline.
func TestFlowP2LegacyVsNew(t *testing.T) {
	rawBase, err := os.ReadFile("testdata/flow_p2_baseline.json")
	if err != nil {
		t.Fatal(err)
	}
	var base map[string]any
	if err := json.Unmarshal(rawBase, &base); err != nil {
		t.Fatal(err)
	}
	ri, _ := flowP1TestResults(t)
	gotRaw, _ := json.Marshal(ri)
	var got map[string]any
	if err := json.Unmarshal(gotRaw, &got); err != nil {
		t.Fatal(err)
	}
	removed := map[string]bool{"da_yun_rels": true, "liu_nian_rels": true}
	added := map[string]bool{"da_yun_interaction": true, "liu_nian_interaction": true}
	for key, oldVal := range base {
		if removed[key] {
			if _, still := got[key]; still {
				t.Errorf("%s should be removed in P2", key)
			}
			continue
		}
		newVal, ok := got[key]
		if !ok {
			t.Errorf("%s disappeared", key)
			continue
		}
		if !reflect.DeepEqual(oldVal, newVal) {
			t.Errorf("%s changed: %v -> %v", key, oldVal, newVal)
		}
	}
	for key := range got {
		if _, existed := base[key]; !existed && !added[key] {
			t.Errorf("unexpected new field %s", key)
		}
	}
}

// TestFlowP2InteractionOracle verifies the new two-layer objects contain BOTH
// the gan-layer and zhi-layer results from the pure function.
func TestFlowP2InteractionOracle(t *testing.T) {
	ri, _ := flowP1TestResults(t)
	if len(ri.DaYunInteraction.GanRels) != 1 || len(ri.DaYunInteraction.ZhiRels) != 1 {
		t.Errorf("da_yun_interaction missing a layer: %+v", ri.DaYunInteraction)
	}
	if len(ri.LiuNianInteraction.GanRels) != 1 || len(ri.LiuNianInteraction.ZhiRels) != 1 {
		t.Errorf("liu_nian_interaction missing a layer: %+v", ri.LiuNianInteraction)
	}
	if ri.DaYunInteraction.ZhuLabel == "" || ri.LiuNianInteraction.ZhuLabel == "" {
		t.Error("interaction zhu_label empty")
	}
}
