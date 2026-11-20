package liuyao

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

func validateYaoValues(yaos [6]int) bool {
	for _, value := range yaos {
		if value < 6 || value > 9 {
			return false
		}
	}
	return true
}

func newBaseChart(yaos [6]int, bz ganzhi.Bazi) Chart {
	chart := computeGuaPan(shakeCoinsFixed(yaos), bz.Ri)
	chart.TombSchool = TombSchool{EarthBranch: "辰", School: "engine_default_chen"}
	chart.YueZhi = bz.Yue.Zhi
	chart.YueGan = bz.Yue.Gan
	chart.XunKong = ganzhi.XunKong(bz.Ri.Gan, bz.Ri.Zhi)
	return chart
}

func withYaoStrength(chart Chart) Chart {
	for i := 0; i < 6; i++ {
		chart.WangShuai[i] = ganzhi.WangShuaiOf(
			ganzhi.ZhiWuxing(chart.Lines[i].Zhi),
			chart.YueZhi,
		)
		chart.DayRelations[i] = dayInteraction(chart.Lines[i].Zhi, chart.RiZhi)
	}
	return withLineDerived(chart)
}

func withYongshen(chart Chart, yongShen YongShen) Chart {
	position := chart.findYongShen(yongShen)
	chart.YongShen = YongShenResult{Name: DecodeYongShen(yongShen), Position: position}

	if position == 0 {
		chart.YongShen.FuShen = chart.findFuShen(yongShen)
		chart.YongShen.FuShenCandidates = chart.findFuShenAll(yongShen)
		chart = withHiddenYongshen(chart, position)
		return chart
	}

	line := chart.Lines[position-1]
	chart.YongShen.WangShuai = chart.WangShuai[position-1]
	chart.YongShen.YuePo = line.YuePo
	chart.YongShen.XunKong = line.XunKong
	chart.YongShen.MuKu = line.MuKu
	chart.YongShen.LiuShou = line.LiuShou
	chart.YongShen.ChangSheng = line.ChangShengYue
	return chart
}

func withHiddenYongshen(chart Chart, position int) Chart {
	if position != 0 {
		return chart
	}
	fuShen := chart.YongShen.FuShen
	if fuShen == nil {
		return chart
	}

	fuZhi := fuShen.Zhi
	fuElement := ganzhi.ZhiWuxing(fuZhi)
	chart.YongShen.IsHidden = true
	chart.YongShen.WangShuai = ganzhi.WangShuaiOf(fuElement, chart.YueZhi)
	chart.YongShen.YuePo = ganzhi.IsLiuChong(fuZhi, chart.YueZhi)
	chart.YongShen.XunKong = fuZhi == chart.XunKong[0] || fuZhi == chart.XunKong[1]
	chart.YongShen.MuKu = len(fuTombSources(&chart, fuElement)) > 0
	chart.YongShen.ChangSheng = lifeStageOf(fuElement, chart.YueZhi)
	return chart
}

func withAdvancedFacts(chart Chart, yongShen YongShen) Chart {
	chart.DongYaoRelations = computeDongYaoRelations(&chart, yongShen)
	chart.Patterns = ComputePatterns(&chart, yongShen)
	chart.TimingCandidates = computeTimingCandidates(&chart, yongShen)
	chart.DayClashFacts = computeDayClashFacts(&chart)
	chart.MovingTransformations = computeMovingTransformations(&chart)
	chart.SanHeCandidates = computeSanHeCandidates(&chart)
	chart.ForceChain = computeForceChain(&chart, yongShen)
	chart.YongShenCandidates = computeYongShenCandidates(&chart, yongShen)
	chart.HiddenLines = computeHiddenLines(&chart)
	chart.BranchRelationFacts = computeBranchRelationFacts(&chart)
	chart.Conflicts = computeConflicts(&chart)
	return chart
}

func withGuaCi(chart Chart) Chart {
	if guaCi, err := GetGuaCi(int(chart.BenGua)); err == nil {
		chart.GuaCi = guaCi
	}
	return chart
}
