package bazi

import "github.com/ml8s/mingligo/pkg/ganzhi"

// flowDayMasterFacts is the single pure entry point for a flow pillar's stem
// relative to the day master. All four flow layers (year/month/day/hour) use
// it so the 生扶/克制 and 十神 dimensions never drift apart.
//
// It composes countGenRest (五行生克) and ShiShenFromGan (十神) — both pure —
// and adds no layer-specific state.
func flowDayMasterFacts(flowGan, dayMaster ganzhi.Gan) (sheng, ke int, shiShen string) {
	flowElem := ganzhi.GanWuxing(flowGan)
	dmElem := ganzhi.GanWuxing(dayMaster)
	gen, rest := countGenRest(flowElem, dmElem)
	return gen, rest, ganzhi.DecodeShiShen(ganzhi.ShiShenFromGan(dayMaster, flowGan))
}

// analyzeFlowAgainstRef computes the two-layer (gan + zhi) interaction of a
// flow pillar against a single reference pillar. Pure; ref==nil yields zero
// value. This is the single entry point used by every flow layer when it
// relates itself to 大运/流年/流月/流日.
func analyzeFlowAgainstRef(flow, ref ganzhi.Zhu) zhuInteraction {
	label := ganzhi.GanName(ref.Gan) + ganzhi.ZhiName(ref.Zhi)
	ganRels := []GanRelation{analyzeGanRelation(flow.Gan, ref.Gan)}
	zhiRels := []ZhiRelation{analyzeZhiRelation(flow.Zhi, ref.Zhi)}
	return zhuInteraction{ZhuLabel: label, GanRels: ganRels, ZhiRels: zhiRels}
}
