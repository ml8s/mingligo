package bazi

import "github.com/ml8s/mingligo/pkg/ganzhi"

// computeFuYi computes the FuYi (扶抑) yongshen analysis from pillars.
// Uses data-driven lookup tables for strength, congGe, and yong/xi/ji.
func FuYi(
	bz ganzhi.Bazi,
	counts WuxingCounts,
	seasonStrength SeasonStrength,
	strength DayMasterStrength,
	relations []NatalRelationFact,
) FuYiResult {
	riYuan := bz.Ri.Gan
	basis := buildFuYiBasis(bz, strength.RootType, strength.YinBiCount, relations)

	// Try congGe rules first (qualitative rule chain, not percentages).
	pat, congYong, congXi, congJi := lookupCongGe(bz)
	if pat != "" {
		return FuYiResult{
			WuxingCount: counts.StringCount(),
			WangShuai:   seasonStrength,
			Yong:        candidateString(congYong),
			Xi:          candidateString(congXi),
			Ji:          candidateString(congJi),
			Strength:    DecodeStrengthLevel(strength.Label),
			Pattern:     pat,
			ModelID:     "support_control_with_day_master_strength",
			Basis:       basis,
		}
	}

	// Normal扶抑 (qualitative, no element count comparison).
	yong, xi, ji := computeNormalYongJi(riYuan, strength.Label)

	return FuYiResult{
		WuxingCount: counts.StringCount(),
		WangShuai:   seasonStrength,
		Yong:        yong,
		Xi:          xi,
		Ji:          ji,
		Strength:    DecodeStrengthLevel(strength.Label),
		ModelID:     "support_control_with_day_master_strength",
		Basis:       basis,
	}
}

// computeNormalYongJi determines yong/xi/ji by strength (qualitative, no scoring).
func computeNormalYongJi(riYuan ganzhi.Gan, strength StrengthLevel) (yongShen, xiShen, jiShen string) {
	dmElem := ganzhi.GanWuxing(riYuan)
	rule, ok := yongJiRules[strength]
	if !ok {
		return NoCandidateWuxing, NoCandidateWuxing, NoCandidateWuxing
	}
	// 五行关系 → 具体五行：克我者官杀、生我者印、同我者比劫、生克我者财生官杀。
	resolve := func(rel RelationKind) ganzhi.Wuxing {
		switch rel {
		case relationKeWo:
			return elementThatControls(dmElem)
		case relationShengWo:
			return elementThatGenerates(dmElem)
		case relationTongWo:
			return dmElem
		case relationShengKeWo:
			return elementThatGenerates(elementThatControls(dmElem))
		}
		return 0
	}
	return candidateString(ganzhi.DecodeWuxing(resolve(rule.Yong))),
		candidateString(ganzhi.DecodeWuxing(resolve(rule.Xi))),
		candidateString(ganzhi.DecodeWuxing(resolve(rule.Ji)))
}
