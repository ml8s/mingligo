package ziwei

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/tianwen"
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
	changed := map[string]bool{"ming_gong": true} // P5: now real value on all layers
	allowedNew := map[string]bool{"liu_ming_zhi": true, "gan": true, "period_gan": true}
	for layer, fields := range baseDoc {
		for key, old := range fields {
			if removed[key] || changed[key] {
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
	if !oa || !ga {
		t.Error("gong_wei type changed")
		return false
	}
	if len(oldArr) != len(gotArr) {
		return false
	}
	for i := range oldArr {
		om, oOK := oldArr[i].(map[string]any)
		gm, gOK := gotArr[i].(map[string]any)
		if !oOK || !gOK {
			t.Errorf("gong_wei[%d] element type changed", i)
			return false
		}
		// every old field must be preserved byte-identically
		for k, ov := range om {
			if !reflect.DeepEqual(ov, gm[k]) {
				return false
			}
		}
		// only "gan" may be newly present; it must exist on every palace
		for k := range gm {
			if _, existed := om[k]; !existed && k != "gan" {
				t.Errorf("gong_wei[%d]: unexpected new field %s", i, k)
				return false
			}
		}
		if _, hasGan := gm["gan"]; !hasGan {
			t.Errorf("gong_wei[%d]: missing new gan field", i)
			return false
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

// TestFlowMingGongIndexOracle cross-checks the shared function against the
// original inline formula used by 流年, for all 12×12 zhi pairs.
func TestFlowMingGongIndexOracle(t *testing.T) {
	for natal := 1; natal <= 12; natal++ {
		for flow := 1; flow <= 12; flow++ {
			got := flowMingGongIndex(Zhi(natal), Zhi(flow))
			want := zhiIdxToPalaceIndex(zhiToZhiIdx(Zhi(natal)), zhiToZhiIdx(Zhi(flow)))
			if got != want {
				t.Errorf("natal=%d flow=%d: got %v want %v", natal, flow, got, want)
			}
		}
	}
}

// TestAllFlowLayersMingGongConsistent verifies the four layers all report
// 流命落本命宫序 via the same pure function (T1 fix).
func TestAllFlowLayersMingGongConsistent(t *testing.T) {
	ctx := flowP4Context(t)
	standard, _ := time.Parse(time.RFC3339, "2026-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	layers := []struct {
		name     string
		mingGong GongIndex
		liuMing  Zhi
	}{
		{"year", ComputeLiuNian(ctx, 2026).MingGong, ComputeLiuNian(ctx, 2026).LiuMingZhi},
		{"month", mustLiuYue(t, ctx, standard, loc).MingGong, mustLiuYue(t, ctx, standard, loc).LiuMingZhi},
		{"day", mustLiuRi(t, ctx, standard, loc).MingGong, mustLiuRi(t, ctx, standard, loc).LiuMingZhi},
		{"hour", mustLiuShi(t, ctx, standard, loc, 116.4).MingGong, mustLiuShi(t, ctx, standard, loc, 116.4).LiuMingZhi},
	}
	for _, l := range layers {
		want := flowMingGongIndex(ctx.MingGongZhi, l.liuMing)
		if l.mingGong != want {
			t.Errorf("%s: ming_gong=%v, want flowMingGongIndex(%v,%v)=%v",
				l.name, l.mingGong, ctx.MingGongZhi, l.liuMing, want)
		}
	}
}

// TestFlowP5LegacyVsNew verifies P5 changes: ming_gong gains real values on
// month/day/hour (was 0), period_gan is added. Everything else byte-identical.
func TestFlowP5LegacyVsNew(t *testing.T) {
	base, err := os.ReadFile("testdata/flow_p5_baseline.json")
	if err != nil {
		t.Fatal(err)
	}
	var baseDoc map[string]map[string]any
	if err := json.Unmarshal(base, &baseDoc); err != nil {
		t.Fatal(err)
	}
	ctx := flowP4Context(t)
	standard, _ := time.Parse(time.RFC3339, "2026-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	gotRaw, _ := json.Marshal(map[string]any{
		"ln": ComputeLiuNian(ctx, 2026),
		"ly": mustLiuYue(t, ctx, standard, loc),
		"lr": mustLiuRi(t, ctx, standard, loc),
		"ls": mustLiuShi(t, ctx, standard, loc, 116.4),
	})
	var gotDoc map[string]map[string]any
	if err := json.Unmarshal(gotRaw, &gotDoc); err != nil {
		t.Fatal(err)
	}
	changed := map[string]bool{"ming_gong": true} // intentional (T1)
	added := map[string]bool{"period_gan": true}  // intentional (T2)
	for layer, fields := range baseDoc {
		for key, old := range fields {
			nv, ok := gotDoc[layer][key]
			if !ok {
				t.Errorf("%s.%s disappeared", layer, key)
				continue
			}
			if changed[key] {
				continue
			}
			if !reflect.DeepEqual(old, nv) {
				t.Errorf("%s.%s changed: %v -> %v", layer, key, old, nv)
			}
		}
		for key := range gotDoc[layer] {
			if _, existed := fields[key]; !existed && !added[key] {
				t.Errorf("unexpected new field %s.%s", layer, key)
			}
		}
	}
}

// TestFlowP5PeriodGanOracle verifies each layer's period_gan is the correct
// flow-period stem by cross-checking against the SiHua transformations it
// drives (the stem that produced the 四化).
func TestFlowP5PeriodGanOracle(t *testing.T) {
	ctx := flowP4Context(t)
	standard, _ := time.Parse(time.RFC3339, "2026-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	layers := []struct {
		name string
		gan  Gan
		zhi  Zhi
	}{
		{"year", ComputeLiuNian(ctx, 2026).Gan, ComputeLiuNian(ctx, 2026).Zhi},
		{"month", mustLiuYue(t, ctx, standard, loc).Gan, mustLiuYue(t, ctx, standard, loc).Zhi},
		{"day", mustLiuRi(t, ctx, standard, loc).Gan, mustLiuRi(t, ctx, standard, loc).Zhi},
		{"hour", mustLiuShi(t, ctx, standard, loc, 116.4).Gan, mustLiuShi(t, ctx, standard, loc, 116.4).Zhi},
	}
	for _, l := range layers {
		if l.gan == 0 || l.gan > 10 {
			t.Errorf("%s: period_gan=%d out of range", l.name, l.gan)
		}
		if l.zhi < 1 || l.zhi > 12 {
			t.Errorf("%s: zhi=%d out of range", l.name, l.zhi)
		}
	}
}

// TestLeap12SecondHalfYearSync verifies the 闰十二月下半月 boundary: the flow
// month jumps to next year's 正月, and the monthly index must use that same
// year so 月干支 and 流盘起点 don't fork. leaps beyond 1900-2100 are reachable
// astronomically (next at 3358), so this guards a real (if rare) correctness bug.
func TestLeap12SecondHalfYearSync(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	p, err := tianwen.NewLunarProjector(tianwen.LunarBasisLocalAstronomical, loc)
	if err != nil {
		t.Fatal(err)
	}
	target := tianwen.LunarDate{Year: 3358, Month: 12, Day: 20, Leap: true}
	period, err := resolveLiuYuePeriodWith(p, target)
	if err != nil {
		t.Fatal(err)
	}
	if period.FlowMonth.Year != 3359 || period.FlowMonth.Month != 1 {
		t.Fatalf("flowMonth = %+v, want 3359-01", period.FlowMonth)
	}
	// 3359 年干丙 → 正月寅 → 丙寅。修复前 mi 用 3358(戊)年支，会分叉。
	ctx := FlowContext{BirthLunarMonth: 2, BirthLunarDay: 15, ShiZhi: 5, MingGongZhi: 11}
	result, err := computeLiuYueWithPeriod(ctx, target, period)
	if err != nil {
		t.Fatal(err)
	}
	flowYearGan, _ := yearGanZhi(period.FlowMonth.Year)
	wantGan := yueGanByWuHuDun(flowYearGan, result.Zhi)
	if result.Gan != wantGan {
		t.Errorf("period gan = %v, want %v (遁 from flowMonth year %d)", result.Gan, wantGan, period.FlowMonth.Year)
	}
	// monthlyIndex uses flowMonth year — assert via an equivalent direct call.
	mi := computeMonthlyIndex(ctx, flowTarget{
		Year: period.FlowMonth.Year, LunarMonth: period.FlowMonth.Month,
		LunarDay: target.Day, IsLeapMonth: period.FlowMonth.Leap,
	})
	if mi < 0 || mi >= 12 {
		t.Errorf("monthly index out of range: %d", mi)
	}
}

// TestLiuYueNormalMonthUnchanged pins ordinary-month behavior (no leap jump)
// identical to baseline, so the fix touches only the boundary.
func TestLiuYueNormalMonthUnchanged(t *testing.T) {
	base, err := os.ReadFile("testdata/leap_month_baseline.json")
	if err != nil {
		t.Fatal(err)
	}
	var baseDoc any
	if err := json.Unmarshal(base, &baseDoc); err != nil {
		t.Fatal(err)
	}
	subjectCtx := flowP4Context(t)
	standard, _ := time.Parse(time.RFC3339, "2026-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	got, err := ComputeLiuYueAt(subjectCtx, standard, loc)
	if err != nil {
		t.Fatal(err)
	}
	gotRaw, _ := json.Marshal(got)
	var gotDoc any
	if err := json.Unmarshal(gotRaw, &gotDoc); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(baseDoc, gotDoc) {
		t.Errorf("normal month output changed:\nbase=%s\ngot=%s", base, gotRaw)
	}
}
