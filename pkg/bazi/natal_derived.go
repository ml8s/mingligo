package bazi

import "github.com/ml8s/mingligo/pkg/ganzhi"

// SpousePalaceState and DayBranchType are deterministic natal-chart facts.
// They are computed as symbols (PalaceState / BranchImageKind) in the domain
// and emitted as their stable text values at the output boundary.

// spousePalaceState classifies the day branch by the first matching relation
// in priority order: 冲 > 合 > 刑 > 害 > 静 (same order the analysis layer uses).
func spousePalaceState(chart Chart) PalaceState {
	priority := []struct {
		kind  NatalRelationKind
		state PalaceState
	}{
		{relationKindLiuChong, palaceStateChong},
		{relationKindZhiLiuHe, palaceStateHe},
		{relationKindLiuXing, palaceStateXing},
		{relationKindLiuHai, palaceStateHai},
	}
	for _, entry := range priority {
		for _, relation := range chart.GuanXi {
			if relation.Kind != entry.kind {
				continue
			}
			for _, member := range relation.Members {
				if member.Pillar == pillarRi {
					return entry.state
				}
			}
		}
	}
	return palaceStateQuiet
}

// dayBranchType classifies the day branch by shen-sha hits and tomb-branch
// membership. Tomb branches are the four 墓库 zhi (辰戌丑未).
func dayBranchType(chart Chart) BranchImageKind {
	dayPillar := chart.Ri
	for _, sha := range dayPillar.ShenSha {
		switch sha.Name {
		case "桃花":
			return branchImagePeach
		case "驿马":
			return branchImageTravel
		}
	}
	zhi := dayPillar.Zhi
	if zhi == ganzhi.ZhiChen || zhi == ganzhi.ZhiXu ||
		zhi == ganzhi.ZhiChou || zhi == ganzhi.ZhiWei {
		return branchImageTomb
	}
	return branchImageNone
}
