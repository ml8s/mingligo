package ziwei

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// TestFlowPalaceGanWuHuDunOracle pins the full 10×12 table against the
// classical 五虎遁 rule (寅宫干 by year stem, then sequential by zhi).
func TestFlowPalaceGanWuHuDunOracle(t *testing.T) {
	// 寅宫干 base: 甲己→丙, 乙庚→戊, 丙辛→庚, 丁壬→壬, 戊癸→甲
	base := map[ganzhi.Gan]ganzhi.Gan{
		ganzhi.GanJia: ganzhi.GanBing, ganzhi.GanJi: ganzhi.GanBing,
		ganzhi.GanYi: ganzhi.GanWu, ganzhi.GanGeng: ganzhi.GanWu,
		ganzhi.GanBing: ganzhi.GanGeng, ganzhi.GanXin: ganzhi.GanGeng,
		ganzhi.GanDing: ganzhi.GanRen, ganzhi.GanRen: ganzhi.GanRen,
		ganzhi.GanWu: ganzhi.GanJia, ganzhi.GanGui: ganzhi.GanJia,
	}
	for flow := ganzhi.GanJia; flow <= ganzhi.GanGui; flow++ {
		for zidx := 0; zidx < 12; zidx++ {
			got := flowPalaceGan(flow, zidx)
			want := Gan(((int(base[flow])-1+zidx)%10+10)%10 + 1)
			if got != want {
				t.Errorf("flow=%v zidx=%d: gan=%v want %v", flow, zidx, got, want)
			}
		}
	}
}

// TestFlowMingZhiInvariant verifies the derived 流命地支 matches IsMing palace.
func TestFlowMingZhiInvariant(t *testing.T) {
	for idx := 0; idx < 12; idx++ {
		mingZhi := flowMingZhi(idx)
		palaces := buildFlowPalaces(idx, ganzhi.GanJia, map[int][]FlowStar{})
		found := 0
		for _, p := range palaces {
			if p.IsMing {
				found++
				if p.Zhi != mingZhi {
					t.Errorf("idx=%d: IsMing palace zhi=%v, flowMingZhi=%v", idx, p.Zhi, mingZhi)
				}
			}
		}
		if found != 1 {
			t.Errorf("idx=%d: %d palaces marked IsMing, want 1", idx, found)
		}
	}
}

// TestFlowP4LegacyVsNew is the differential gate for P4: every pre-existing
// field must remain except ming_gong_name (intentionally removed — constant
// "命宫", zero information). New fields: palace gan and liu_ming_zhi.
func TestFlowP4LegacyVsNew(t *testing.T) {
	base, err := os.ReadFile("testdata/flow_p4_baseline.json")
	if err != nil {
		t.Fatal(err)
	}
	var baseDoc map[string]map[string]any
	if err := json.Unmarshal(base, &baseDoc); err != nil {
		t.Fatal(err)
	}
	standard, _ := time.Parse(time.RFC3339, "2026-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	ctx := flowP4Context(t)
	got := map[string]any{
		"ln": ComputeLiuNian(ctx, 2026),
		"ly": mustLiuYue(t, ctx, standard, loc),
		"lr": mustLiuRi(t, ctx, standard, loc),
		"ls": mustLiuShi(t, ctx, standard, loc, 116.4),
	}
	gotRaw, _ := json.Marshal(got)
	var gotDoc map[string]map[string]any
	if err := json.Unmarshal(gotRaw, &gotDoc); err != nil {
		t.Fatal(err)
	}
	removed := map[string]bool{"ming_gong_name": true}
	allowedNew := map[string]bool{"liu_ming_zhi": true, "gan": true}
	for layer, fields := range baseDoc {
		for key, old := range fields {
			if removed[key] {
				continue
			}
			nv, ok := gotDoc[layer][key]
			if !ok {
				t.Errorf("%s.%s disappeared", layer, key)
				continue
			}
			// gong_wei items gained "gan"; compare with that key excluded.
			if key == "gong_wei" {
				if !gongWeiSubsetEqual(t, old, nv) {
					t.Errorf("%s.gong_wei changed beyond new gan field", layer)
				}
				continue
			}
			if !reflect.DeepEqual(old, nv) {
				t.Errorf("%s.%s changed: %v -> %v", layer, key, old, nv)
			}
		}
		for key := range gotDoc[layer] {
			if _, existed := fields[key]; !existed && !allowedNew[key] {
				t.Errorf("unexpected new field %s.%s", layer, key)
			}
		}
	}
}

func gongWeiSubsetEqual(t *testing.T, old, got any) bool {
	t.Helper()
	oldArr, oa := old.([]any)
	gotArr, ga := got.([]any)
	if !oa || !ga || len(oldArr) != len(gotArr) {
		return false
	}
	for i := range oldArr {
		om, gm := oldArr[i].(map[string]any), gotArr[i].(map[string]any)
		for k, ov := range om {
			if k == "gan" {
				continue
			}
			if !reflect.DeepEqual(ov, gm[k]) {
				return false
			}
		}
	}
	return true
}

func flowP4Context(t *testing.T) FlowContext {
	t.Helper()
	standard, _ := time.Parse(time.RFC3339, "1984-02-15T00:00:00Z")
	chart, err := ComputeChartAt(standard, 116.4, ganzhi.Male)
	if err != nil {
		t.Fatal(err)
	}
	return chart.FlowContext()
}

func mustLiuYue(t *testing.T, ctx FlowContext, s time.Time, l *time.Location) LiuYue {
	v, err := ComputeLiuYueAt(ctx, s, l)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func mustLiuRi(t *testing.T, ctx FlowContext, s time.Time, l *time.Location) LiuRi {
	v, err := ComputeLiuRiAt(ctx, s, l)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
func mustLiuShi(t *testing.T, ctx FlowContext, s time.Time, l *time.Location, lon float64) LiuShi {
	v, err := ComputeLiuShiAt(ctx, s, l, lon)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
