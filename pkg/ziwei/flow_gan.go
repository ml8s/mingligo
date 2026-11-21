package ziwei

import "github.com/ml8s/mingligo/pkg/ganzhi"

// flowPalaceGan derives the heavenly stem of a flow-chart palace by 五虎遁.
// flowGan is the flow period's stem; zhiIdx is the palace's earth-zhi index in
// 安星序 (寅=0 ... 丑=11). Pure function — the single authority for flow palace
// stems across year/month/day/hour layers.
func flowPalaceGan(flowGan ganzhi.Gan, zhiIdx int) ganzhi.Gan {
	yinGan := yueGanByWuHuDun(flowGan, 3) // 寅月(支=3)干
	offset := (zhiIdx%12 + 12) % 12       // 寅=0 relative to 寅
	return Gan(((int(yinGan)-1+offset)%10+10)%10 + 1)
}

// flowMingZhi returns the earth-zhi of the flow chart's 命宫 from the flow
// start index (安星序). Pure; used to cross-check IsMing in the same chart.
func flowMingZhi(liuPanIdx int) Zhi {
	return zhiIdxToZhi(anXingIdxToZhiIdx(liuPanIdx))
}

// flowMingGongIndex returns the natal palace index where the flow chart's
// 命宫 falls. Pure; shared by all four flow layers so "流命落本命何宫" is
// computed by one authority rather than duplicated per layer.
func flowMingGongIndex(natalMingZhi, liuMingZhi Zhi) GongIndex {
	return zhiIdxToPalaceIndex(zhiToZhiIdx(natalMingZhi), zhiToZhiIdx(liuMingZhi))
}
