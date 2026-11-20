package bazi

import "github.com/ml8s/mingligo/pkg/ganzhi"

// NatalDerivedFacts are deterministic facts derived from the natal chart's
// canonical facts (guan_xi relations and shen_sha hits). Judgment layers
// consume these directly rather than re-deriving them from raw chart fields.
type NatalDerivedFacts struct {
	SpousePalaceState string `json:"spouse_palace_state"`
	DayBranchType     string `json:"day_branch_type"`
}

// computeNatalDerivedFacts derives spouse-palace state and day-branch type
// using domain symbols (not display strings).
func computeNatalDerivedFacts(chart Chart) NatalDerivedFacts {
	return NatalDerivedFacts{
		SpousePalaceState: spousePalaceState(chart),
		DayBranchType:     dayBranchType(chart),
	}
}

// spousePalaceState classifies the day branch by the first matching relation
// in priority order: 冲 > 合 > 刑 > 害 > 静 (same order the analysis layer uses).
func spousePalaceState(chart Chart) string {
	priority := []struct {
		kind NatalRelationKind
		name string
	}{
		{relationKindLiuChong, "冲"},
		{relationKindZhiLiuHe, "合"},
		{relationKindLiuXing, "刑"},
		{relationKindLiuHai, "害"},
	}
	for _, entry := range priority {
		for _, relation := range chart.GuanXi {
			if relation.Kind != entry.kind {
				continue
			}
			for _, member := range relation.Members {
				if member.Pillar == pillarRi {
					return entry.name
				}
			}
		}
	}
	return "静"
}

// dayBranchType classifies the day branch by shen-sha hits and tomb-branch
// membership. Tomb branches are the four 墓库 zhi (辰戌丑未).
func dayBranchType(chart Chart) string {
	dayPillar := chart.Ri
	for _, sha := range dayPillar.ShenSha {
		switch sha.Name {
		case "桃花":
			return "桃花"
		case "驿马":
			return "驿马"
		}
	}
	zhi := dayPillar.Zhi
	if zhi == ganzhi.ZhiChen || zhi == ganzhi.ZhiXu ||
		zhi == ganzhi.ZhiChou || zhi == ganzhi.ZhiWei {
		return "墓库"
	}
	return "无"
}
