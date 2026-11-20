package bazi

import (
	"strings"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// RelationDetection holds complete relation detection evidence for a Bazi
// chart: 合冲刑害破暗合 are detector-owned facts before canonical projection.
type RelationDetection struct {
	GanHe        []GanHeDetection
	GanChong     []GanChongDetection
	ZhiLiuHe     []ZhiPairDetection
	SanHe        []TripleRelationDetection
	SanHui       []TripleRelationDetection
	LiuChong     []ZhiPairDetection
	LiuHai       []ZhiPairDetection
	LiuXing      []ZhiPairDetection
	LiuPo        []ZhiPairDetection
	AnHe         []ZhiPairDetection
	SanHePartial []TripleRelationDetection
}

// GanHeDetection describes a 天干五合 pair. Adjacency is an engine-owned fact;
// non-adjacent pairs remain candidates and must not be promoted to tight combinations.
type GanHeDetection struct {
	GanA      ganzhi.Gan
	GanB      ganzhi.Gan
	PillarA   int
	PillarB   int
	Position  RelationPosition
	Contested bool
	HeElement ganzhi.Wuxing
}

// GanChongDetection describes a paired gan relationship (天干相冲). Position is an
// engine-owned strength hint: tight pairs are adjacent, others are candidates.
type GanChongDetection struct {
	GanA     ganzhi.Gan
	GanB     ganzhi.Gan
	PillarA  int
	PillarB  int
	Position RelationPosition
}

// ZhiPairDetection describes a paired zhi relationship (六合/六冲/六害/相刑/六破/暗合)
// between two pillars.
type ZhiPairDetection struct {
	ZhiA    ganzhi.Zhi
	ZhiB    ganzhi.Zhi
	PillarA int
	PillarB int
	Element ganzhi.Wuxing
}

// TripleRelationDetection describes a complete 三合局/三会方 or a 半合局 candidate.
type TripleRelationDetection struct {
	Kind     string
	Name     string
	Element  ganzhi.Wuxing
	Branches []ganzhi.Zhi
	Pillars  []string
}

// DetectRelations detects all published relation candidates from four pillars.
func DetectRelations(bz ganzhi.Bazi) RelationDetection {
	return RelationDetection{
		GanHe:        detectGanHe(bz),
		GanChong:     detectGanChong(bz),
		ZhiLiuHe:     detectZhiPairs(bz, ganzhi.IsZhiHe, true),
		SanHe:        detectTriple(bz, ganzhi.TripleHeList, relSanHe, "局"),
		SanHui:       detectTriple(bz, ganzhi.TripleHuiList, relSanHui, "方"),
		LiuChong:     detectZhiPairs(bz, ganzhi.IsLiuChong, false),
		LiuHai:       detectZhiPairs(bz, ganzhi.IsHai, false),
		LiuXing:      detectZhiPairs(bz, ganzhi.IsXing, false),
		LiuPo:        detectZhiPairs(bz, ganzhi.IsPo, false),
		AnHe:         detectZhiPairs(bz, ganzhi.IsAnHe, false),
		SanHePartial: detectPartialTriple(bz),
	}
}

const (
	zhuNian = 0
	zhuYue  = 1
	zhuRi   = 2
	zhuShi  = 3
)

func detectGanHe(bz ganzhi.Bazi) []GanHeDetection {
	pairs := []GanHeDetection{}
	zhus := bz.Slice()
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			a, b := zhus[i].Gan, zhus[j].Gan
			if !ganzhi.IsGanHe(a, b) {
				continue
			}
			position := positionSeparated
			switch j - i {
			case 1:
				position = positionAdjacent
			case 3:
				position = positionRemote
			}
			pairs = append(pairs, GanHeDetection{
				GanA: a, GanB: b,
				PillarA: i, PillarB: j, Position: position,
				HeElement: ganHeResult(a, b),
			})
		}
	}
	for i := range pairs {
		participations := 0
		for _, other := range pairs {
			if other.PillarA == pairs[i].PillarA || other.PillarA == pairs[i].PillarB ||
				other.PillarB == pairs[i].PillarA || other.PillarB == pairs[i].PillarB {
				participations++
			}
		}
		pairs[i].Contested = participations > 1
	}
	return pairs
}

func detectGanChong(bz ganzhi.Bazi) []GanChongDetection {
	zhus := bz.Slice()
	pairs := []GanChongDetection{}
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			a, b := zhus[i].Gan, zhus[j].Gan
			if !ganzhi.IsGanChong(a, b) {
				continue
			}
			position := positionSeparated
			switch j - i {
			case 1:
				position = positionAdjacent
			case 3:
				position = positionRemote
			}
			pairs = append(pairs, GanChongDetection{
				GanA: a, GanB: b,
				PillarA: i, PillarB: j, Position: position,
			})
		}
	}
	return pairs
}

func ganHeResult(a, b ganzhi.Gan) ganzhi.Wuxing {
	for _, p := range ganzhi.GanHes {
		if (a == p.A && b == p.B) || (a == p.B && b == p.A) {
			return p.Result
		}
	}
	return 0
}

// detectZhiPairs checks all pillar zhi pairs using check. When withElem is true,
// also computes the 六合 element produced by the pair.
func detectZhiPairs(bz ganzhi.Bazi, check func(ganzhi.Zhi, ganzhi.Zhi) bool, withElem bool) []ZhiPairDetection {
	zhus := bz.Slice()
	rels := []ZhiPairDetection{}
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			za, zb := zhus[i].Zhi, zhus[j].Zhi
			if check(za, zb) {
				var elem ganzhi.Wuxing
				if withElem {
					elem = zhiHeElement(za, zb)
				}
				rels = append(rels, ZhiPairDetection{
					ZhiA:    za,
					ZhiB:    zb,
					PillarA: i,
					PillarB: j,
					Element: elem,
				})
			}
		}
	}
	return rels
}

func detectTriple(bz ganzhi.Bazi, list []ganzhi.SanHeHui, typ, suffix string) []TripleRelationDetection {
	bs := zhiSet(bz)
	results := []TripleRelationDetection{}
	for _, tr := range list {
		if countZhi(bs, tr.Zhi...) == len(tr.Zhi) {
			results = append(results, TripleRelationDetection{
				Kind:     typ,
				Name:     tripleName(tr.Zhi, tr.Element, suffix),
				Element:  tr.Element,
				Branches: tr.Zhi,
				Pillars:  pillarsWithAnyZhi(bz, tr.Zhi...),
			})
		}
	}
	return results
}

// detectPartialTriple finds 半合 (partial 三合: any 2 of 3 branches present).
// Only reports when the complete 三合 is NOT already satisfied.
func detectPartialTriple(bz ganzhi.Bazi) []TripleRelationDetection {
	bs := zhiSet(bz)
	results := []TripleRelationDetection{}
	for _, tr := range ganzhi.TripleHeList {
		matched := countZhi(bs, tr.Zhi...)
		if matched == len(tr.Zhi) {
			continue // complete 三合, skip
		}
		if matched < 2 {
			continue
		}
		// 半合必须包含旺支（三合表中第二支）。缺旺支的生墓二支是拱合候选，
		// 不标为半合，避免把 巳丑 误判成 巳酉丑 半合。
		if !bs[tr.Zhi[1]] {
			continue
		}
		// Find which 2 branches are present
		var present []ganzhi.Zhi
		for _, z := range tr.Zhi {
			if bs[z] {
				present = append(present, z)
			}
		}
		results = append(results, TripleRelationDetection{
			Kind:     "半合",
			Name:     tripleName(present, tr.Element, "半合局"),
			Element:  tr.Element,
			Branches: present,
			Pillars:  pillarsWithAnyZhi(bz, present...),
		})
	}
	return results
}

func pillarsWithAnyZhi(bz ganzhi.Bazi, targets ...ganzhi.Zhi) []string {
	pillars := make([]string, 0, 4)
	for index, pillar := range bz.Slice() {
		for _, target := range targets {
			if pillar.Zhi == target {
				pillars = append(pillars, zhuLabels[index])
				break
			}
		}
	}
	return pillars
}

func zhiHeElement(a, b ganzhi.Zhi) ganzhi.Wuxing {
	for _, p := range ganzhi.ZhiHes {
		if (a == p.A && b == p.B) || (a == p.B && b == p.A) {
			return p.Result
		}
	}
	return 0
}

func zhiSet(bz ganzhi.Bazi) [13]bool {
	zhus := bz.Slice()
	var bs [13]bool
	for _, p := range zhus {
		if b := int(p.Zhi); b >= 1 && b <= 12 {
			bs[b] = true
		}
	}
	return bs
}

func countZhi(bs [13]bool, targets ...ganzhi.Zhi) int {
	c := 0
	for _, t := range targets {
		if t >= 1 && t <= 12 && bs[int(t)] {
			c++
		}
	}
	return c
}

func tripleName(zhi []ganzhi.Zhi, element ganzhi.Wuxing, suffix string) string {
	parts := make([]string, len(zhi))
	for i, b := range zhi {
		parts[i] = ganzhi.ZhiName(b)
	}
	return strings.Join(parts, "") + ganzhi.DecodeWuxing(element) + suffix
}

func containsPair(list []ganzhi.Zhi, a, b ganzhi.Zhi) bool {
	hasA, hasB := false, false
	for _, v := range list {
		if v == a {
			hasA = true
		}
		if v == b {
			hasB = true
		}
	}
	return hasA && hasB
}
