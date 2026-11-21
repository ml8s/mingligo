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
