package bazi

import (
	"fmt"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// GanRelation describes a single gan-to-gan relationship.
type GanRelation struct {
	GanA     ganzhi.Gan         `json:"gan_a"`
	GanB     ganzhi.Gan         `json:"gan_b"`
	Kind     GanZhiRelationKind `json:"relation_kind"`
	Relation string             `json:"relation"`
}

// ZhiRelation describes a single zhi-to-zhi relationship.
type ZhiRelation struct {
	ZhiA   ganzhi.Zhi         `json:"zhi_a"`
	ZhiB   ganzhi.Zhi         `json:"zhi_b"`
	Kind   GanZhiRelationKind `json:"relation_kind"`
	Detail string             `json:"relation_explanation"`
}

// Relation type constants for gan and zhi interactions.
const (
	relSame         GanZhiRelationKind = "相同"
	relSheng        GanZhiRelationKind = "相生"
	relKe           GanZhiRelationKind = "相克"
	relGanHe        GanZhiRelationKind = "天干五合"
	relNone         GanZhiRelationKind = "无"
	relLiuHe        GanZhiRelationKind = "六合"
	relSanHe        GanZhiRelationKind = "三合"
	relSanHePartial GanZhiRelationKind = "半合"
	relSanHeGong    GanZhiRelationKind = "三合拱"
	relSanHui       GanZhiRelationKind = "三会"
	relSanHuiGong   GanZhiRelationKind = "三会拱"
	relLiuChong     GanZhiRelationKind = "六冲"
	relXing         GanZhiRelationKind = "相刑"
	relLiuHai       GanZhiRelationKind = "六害"
	relAnHe         GanZhiRelationKind = "暗合"
	relPo           GanZhiRelationKind = "破"
)

// zhuInteraction holds gan and zhi relations for one pillar against the bazi chart.
type zhuInteraction struct {
	ZhuLabel string        `json:"zhu_label"`
	GanRels  []GanRelation `json:"gan_rels"`
	ZhiRels  []ZhiRelation `json:"zhi_rels"`
}

// analyzeGanRelation checks the relationship between two gan.
func analyzeGanRelation(a, b ganzhi.Gan) GanRelation {
	r := GanRelation{GanA: a, GanB: b}
	if a == b {
		r.Kind = relSame
		r.Relation = fmt.Sprintf("%s%s同气", ganzhi.GanName(a), ganzhi.GanName(b))
		return r
	}

	for _, p := range ganzhi.GanHes {
		if (a == p.A && b == p.B) || (a == p.B && b == p.A) {
			r.Kind = relGanHe
			r.Relation = fmt.Sprintf("%s%s合化%s", ganzhi.GanName(a), ganzhi.GanName(b), ganzhi.DecodeWuxing(p.Result))
			return r
		}
	}

	aElem, bElem := ganzhi.GanWuxing(a), ganzhi.GanWuxing(b)
	if aElem == bElem {
		r.Kind = relSame
		r.Relation = fmt.Sprintf("%s%s同气", ganzhi.GanName(a), ganzhi.GanName(b))
	} else if ganzhi.Sheng(aElem, bElem) {
		r.Kind = relSheng
		r.Relation = fmt.Sprintf("%s生%s", ganzhi.GanName(a), ganzhi.GanName(b))
	} else if ganzhi.Sheng(bElem, aElem) {
		r.Kind = relSheng
		r.Relation = fmt.Sprintf("%s生%s", ganzhi.GanName(b), ganzhi.GanName(a))
	} else if ganzhi.Ke(aElem, bElem) {
		r.Kind = relKe
		r.Relation = fmt.Sprintf("%s克%s", ganzhi.GanName(a), ganzhi.GanName(b))
	} else if ganzhi.Ke(bElem, aElem) {
		r.Kind = relKe
		r.Relation = fmt.Sprintf("%s克%s", ganzhi.GanName(b), ganzhi.GanName(a))
	} else {
		r.Kind = relNone
		r.Relation = "无特殊关系"
	}
	return r
}

// analyzeZhiRelation checks relationship types between exactly two zhi.
// A pair alone can never be a complete 三合 / 三会: two members with the
// imperial branch are 半合; the first-and-last pair is a 拱 candidate.  Full
// groups must be detected from the whole chart by DetectRelations.
func analyzeZhiRelation(a, b ganzhi.Zhi) ZhiRelation {
	r := ZhiRelation{ZhiA: a, ZhiB: b}
	if a == b {
		// 同支：辰午酉亥 同支相见为自刑，其余为同气。
		if ganzhi.IsXing(a, b) {
			r.Kind = relXing
			r.Detail = fmt.Sprintf("%s%s自刑", ganzhi.ZhiName(a), ganzhi.ZhiName(b))
			return r
		}
		r.Kind = relSame
		r.Detail = fmt.Sprintf("%s%s同气", ganzhi.ZhiName(a), ganzhi.ZhiName(b))
		return r
	}

	for _, p := range ganzhi.ZhiHes {
		if (a == p.A && b == p.B) || (a == p.B && b == p.A) {
			r.Kind = relLiuHe
			r.Detail = fmt.Sprintf("%s%s合化%s", ganzhi.ZhiName(a), ganzhi.ZhiName(b), ganzhi.DecodeWuxing(p.Result))
			return r
		}
	}

	if rel, ok := pairTripleHeRelation(a, b); ok {
		return rel
	}

	for _, p := range ganzhi.ChongPairs {
		if (a == p.A && b == p.B) || (a == p.B && b == p.A) {
			r.Kind = relLiuChong
			r.Detail = fmt.Sprintf("%s%s相冲", ganzhi.ZhiName(a), ganzhi.ZhiName(b))
			return r
		}
	}

	// 相刑：与 ganzhi.IsXing 同口径——自刑（辰午酉亥）仅同支相见为刑，
	// 组内不同支（如午亥、辰酉）不构成刑。
	// 优先级：六合/半合 / 拱 / 六冲 > 相刑 > 六害（寅巳相刑显著于六害；寅申冲/丑未冲显著于刑）。
	if ganzhi.IsXing(a, b) {
		r.Kind = relXing
		xType := "刑"
		for _, x := range ganzhi.XingGroups {
			if containsPair(x.Zhi, a, b) {
				xType = ganzhi.DecodeXingKind(x.Type)
				break
			}
		}
		r.Detail = fmt.Sprintf("%s%s%s", ganzhi.ZhiName(a), ganzhi.ZhiName(b), xType)
		return r
	}

	for _, p := range ganzhi.HaiPairs {
		if (a == p.A && b == p.B) || (a == p.B && b == p.A) {
			r.Kind = relLiuHai
			r.Detail = fmt.Sprintf("%s%s相害", ganzhi.ZhiName(a), ganzhi.ZhiName(b))
			return r
		}
	}

	if rel, ok := pairTripleHuiRelation(a, b); ok {
		return rel
	}

	if ganzhi.IsAnHe(a, b) {
		r.Kind = relAnHe
		r.Detail = fmt.Sprintf("%s%s暗合", ganzhi.ZhiName(a), ganzhi.ZhiName(b))
		return r
	}

	if ganzhi.IsPo(a, b) {
		r.Kind = relPo
		r.Detail = fmt.Sprintf("%s%s相破", ganzhi.ZhiName(a), ganzhi.ZhiName(b))
		return r
	}

	r.Kind = relNone
	r.Detail = "无特殊关系"
	return r
}

// pairTripleHeRelation classifies two members of a trine. The trine's second
// branch is the imperial / prosperous branch; a half combination must include
// it. The first and last branches form the classical arch candidate only.
func pairTripleHeRelation(a, b ganzhi.Zhi) (ZhiRelation, bool) {
	r := ZhiRelation{ZhiA: a, ZhiB: b}
	for _, tr := range ganzhi.TripleHeList {
		if !containsPair(tr.Zhi, a, b) {
			continue
		}
		imperial := tr.Zhi[1]
		if a == imperial || b == imperial {
			var missing ganzhi.Zhi
			for _, z := range tr.Zhi {
				if z != a && z != b {
					missing = z
					break
				}
			}
			r.Kind = relSanHePartial
			r.Detail = fmt.Sprintf("%s%s半合%s局，缺%s", ganzhi.ZhiName(a), ganzhi.ZhiName(b), ganzhi.DecodeWuxing(tr.Element), ganzhi.ZhiName(missing))
			return r, true
		}
		if (a == tr.Zhi[0] && b == tr.Zhi[2]) || (a == tr.Zhi[2] && b == tr.Zhi[0]) {
			r.Kind = relSanHeGong
			r.Detail = fmt.Sprintf("%s%s拱%s，%s局候选", ganzhi.ZhiName(a), ganzhi.ZhiName(b), ganzhi.ZhiName(imperial), ganzhi.DecodeWuxing(tr.Element))
			return r, true
		}
	}
	return ZhiRelation{}, false
}

// pairTripleHuiRelation classifies the first-and-last seasonal pair as an arch
// candidate. Other two-branch subsets are not a 三会 and return no relation.
func pairTripleHuiRelation(a, b ganzhi.Zhi) (ZhiRelation, bool) {
	r := ZhiRelation{ZhiA: a, ZhiB: b}
	for _, group := range ganzhi.TripleHuiList {
		if (a == group.Zhi[0] && b == group.Zhi[2]) || (a == group.Zhi[2] && b == group.Zhi[0]) {
			r.Kind = relSanHuiGong
			r.Detail = fmt.Sprintf("%s%s拱%s，%s方候选", ganzhi.ZhiName(a), ganzhi.ZhiName(b), ganzhi.ZhiName(group.Zhi[1]), ganzhi.DecodeWuxing(group.Element))
			return r, true
		}
	}
	return ZhiRelation{}, false
}

// analyzeZhuWithBazi analyzes one pillar against all 4 bazi chart pillars.
func analyzeZhuWithBazi(zhu ganzhi.Zhu, bz ganzhi.Bazi) ([]GanRelation, []ZhiRelation) {
	ganRels := make([]GanRelation, 4)
	zhiRels := make([]ZhiRelation, 4)
	for i, np := range bz.Slice() {
		ganRels[i] = analyzeGanRelation(zhu.Gan, np.Gan)
		zhiRels[i] = analyzeZhiRelation(zhu.Zhi, np.Zhi)
	}
	return ganRels, zhiRels
}
