package bazi

import (
	"testing"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// TestFlowDayMasterFactsExhaustive verifies the new pure function against its
// two authoritative sources for every stem pair (10×10). This is the first
// correctness gate: the composed function must equal the primitives it wraps.
func TestFlowDayMasterFactsExhaustive(t *testing.T) {
	for flow := ganzhi.GanJia; flow <= ganzhi.GanGui; flow++ {
		for dm := ganzhi.GanJia; dm <= ganzhi.GanGui; dm++ {
			sheng, ke, shiShen := flowDayMasterFacts(flow, dm)

			wantGen, wantRest := countGenRest(ganzhi.GanWuxing(flow), ganzhi.GanWuxing(dm))
			if sheng != wantGen || ke != wantRest {
				t.Errorf("flow=%v dm=%v: got gen/rest=%d/%d want %d/%d",
					flow, dm, sheng, ke, wantGen, wantRest)
			}
			if want := ganzhi.DecodeShiShen(ganzhi.ShiShenFromGan(dm, flow)); shiShen != want {
				t.Errorf("flow=%v dm=%v: shi_shen=%q want %q", flow, dm, shiShen, want)
			}
		}
	}
}

// TestFlowDayMasterFactsKnownCases pins classical values as an oracle,
// independent of countGenRest/ShiShenFromGan, so a shared bug cannot hide.
func TestFlowDayMasterFactsKnownCases(t *testing.T) {
	cases := []struct {
		flow, dm  ganzhi.Gan
		sheng, ke int
		shiShen   string
	}{
		{ganzhi.GanJia, ganzhi.GanJia, 1, 0, "比肩"},  // 甲见甲：同我
		{ganzhi.GanYi, ganzhi.GanJia, 1, 0, "劫财"},   // 乙见甲：同五行异阴阳
		{ganzhi.GanBing, ganzhi.GanJia, 0, 1, "食神"}, // 甲(木)生丙(火)：日主生流年=泄，不生扶日主
		{ganzhi.GanGeng, ganzhi.GanJia, 0, 1, "七杀"}, // 庚(金)克甲(木)：克我者
	}
	for _, tc := range cases {
		sheng, ke, shiShen := flowDayMasterFacts(tc.flow, tc.dm)
		if sheng != tc.sheng || ke != tc.ke || shiShen != tc.shiShen {
			t.Errorf("flow=%v dm=%v: got %d/%d/%q want %d/%d/%q",
				tc.flow, tc.dm, sheng, ke, shiShen, tc.sheng, tc.ke, tc.shiShen)
		}
	}
}

// TestAnalyzeFlowAgainstRefEquivalence verifies the composed function equals
// the primitive relation analyzers for all 60×60 pillar pairs.
func TestAnalyzeFlowAgainstRefEquivalence(t *testing.T) {
	for f := 1; f <= 60; f++ {
		for r := 1; r <= 60; r++ {
			flow := ganzhi.SixtyToZhu(f)
			ref := ganzhi.SixtyToZhu(r)
			got := analyzeFlowAgainstRef(flow, ref)
			if len(got.GanRels) != 1 || got.GanRels[0] != analyzeGanRelation(flow.Gan, ref.Gan) {
				t.Fatalf("flow=%d ref=%d gan mismatch", f, r)
			}
			if len(got.ZhiRels) != 1 || got.ZhiRels[0] != analyzeZhiRelation(flow.Zhi, ref.Zhi) {
				t.Fatalf("flow=%d ref=%d zhi mismatch", f, r)
			}
		}
	}
}
