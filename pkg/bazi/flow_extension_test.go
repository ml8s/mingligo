package bazi

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/tianwen"
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
		"da_yun_interaction": true, "liu_nian_interaction": true, "liu_ri_interaction": true,
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

// TestFlowP3LegacyVsNew verifies flow-month/flow-hour upper interactions were
// added without changing any pre-existing field.
func TestFlowP3LegacyVsNew(t *testing.T) {
	rawBase, err := os.ReadFile("testdata/flow_p3_baseline.json")
	if err != nil {
		t.Fatal(err)
	}
	var base map[string]map[string]any
	if err := json.Unmarshal(rawBase, &base); err != nil {
		t.Fatal(err)
	}
	core := flowP1Core(t)
	standard, _ := time.Parse(time.RFC3339, "2024-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	yue, err := ComputeLiuYueAt(core, standard, loc)
	if err != nil {
		t.Fatal(err)
	}
	shi, err := ComputeLiuShiAt(core, standard, loc, 116.4)
	if err != nil {
		t.Fatal(err)
	}
	gotRaw, _ := json.Marshal(map[string]any{"yue": yue, "shi": shi})
	var got map[string]map[string]any
	if err := json.Unmarshal(gotRaw, &got); err != nil {
		t.Fatal(err)
	}
	allowed := map[string]bool{
		"da_yun_interaction": true, "liu_nian_interaction": true, "liu_ri_interaction": true,
	}
	for layer, fields := range base {
		for key, old := range fields {
			nv, ok := got[layer][key]
			if !ok {
				t.Errorf("%s.%s disappeared", layer, key)
				continue
			}
			if !reflect.DeepEqual(old, nv) {
				t.Errorf("%s.%s changed: %v -> %v", layer, key, old, nv)
			}
		}
		for key := range got[layer] {
			if _, existed := fields[key]; !existed && !allowed[key] {
				t.Errorf("unexpected new field %s.%s", layer, key)
			}
		}
	}
}

// TestFlowP3InteractionOracle verifies all upper-layer interactions carry both
// gan and zhi layers and never panic when a reference is absent.
func TestFlowP3InteractionOracle(t *testing.T) {
	core := flowP1Core(t)
	standard, _ := time.Parse(time.RFC3339, "2024-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	yue, _ := ComputeLiuYueAt(core, standard, loc)
	shi, _ := ComputeLiuShiAt(core, standard, loc, 116.4)
	for name, inter := range map[string]zhuInteraction{
		"yue.da_yun": yue.DaYunInteraction, "yue.liu_nian": yue.LiuNianInteraction,
		"shi.da_yun": shi.DaYunInteraction, "shi.liu_nian": shi.LiuNianInteraction,
		"shi.liu_ri": shi.LiuRiInteraction,
	} {
		if len(inter.GanRels) != 1 || len(inter.ZhiRels) != 1 || inter.ZhuLabel == "" {
			t.Errorf("%s incomplete: %+v", name, inter)
		}
	}
}

// TestFlowRefsBoundaryBases verifies the corrected time bases: each layer's
// upper-layer references must share the same civil/true-solar authority as
// its own pillar, including late-zi cross-day and 立春 boundary.
func TestFlowRefsBoundaryBases(t *testing.T) {
	core := flowP1Core(t)
	loc, _ := time.LoadLocation("Asia/Shanghai")

	t.Run("liushi late-zi uses next day day-pillar", func(t *testing.T) {
		// Directly construct the resolved time with LateZi=true to test liuShiRefs
		// without depending on longitude/timezone resolution details.
		st := LiuShiTime{
			LastDate: tianwen.NewDate(2024, 6, 15),
			LateZi:   true,
			ShiZhi:   ganzhi.ZhiZi,
		}
		refs, err := liuShiRefs(core, st)
		if err != nil {
			t.Fatal(err)
		}
		wantDay := tianwen.RiZhu(st.LastDate.AddDays(1).Time())
		if refs.liuri == nil || refs.liuri.Gan != wantDay.Gan || refs.liuri.Zhi != wantDay.Zhi {
			t.Errorf("liushi liuri = %+v, want %+v (next day)", refs.liuri, wantDay)
		}
		// 大运/流年 also keyed to the true-solar day (not civil next day)
		wantNian, _ := ResolveLiuNianTime(st.LastDate)
		if refs.liunian == nil || refs.liunian.Gan != wantNian.NianGan || refs.liunian.Zhi != wantNian.NianZhi {
			t.Errorf("liushi liunian ref = %+v, want %+v", refs.liunian, wantNian)
		}
	})

	t.Run("liuyue refs nian matches pillar nian", func(t *testing.T) {
		standard, _ := time.Parse(time.RFC3339, "2024-02-04T02:00:00Z") // 立春 boundary 10:00 CST
		nianInstant, _ := ResolveLiuNianTimeAt(standard, loc)
		nianDate, _ := ResolveLiuNianTime(tianwen.StandardDateAt(standard, loc))
		// Document the two authorities may differ at the boundary; the fix makes
		// liu-yue refs use the instant-based one (same as its month pillar).
		_ = nianDate
		day := tianwen.RiZhu(tianwen.StandardDateAt(standard, loc).Time())
		refs, err := resolveFlowRefs(core, nianInstant, ganzhi.Zhu{Gan: day.Gan, Zhi: day.Zhi})
		if err != nil {
			t.Fatal(err)
		}
		if refs.liunian == nil {
			t.Fatal("liunian ref nil")
		}
		wantZhu := nianInstant
		if refs.liunian.Gan != wantZhu.NianGan || refs.liunian.Zhi != wantZhu.NianZhi {
			t.Errorf("liuyue liunian ref = %+v, want nian-year pillar %+v", refs.liunian, wantZhu)
		}
	})
}
