package liuyao

import (
	"strings"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// DayRelation describes 日建与爻的关系.
type DayRelation struct {
	Relations []DayRelationType `json:"relations"` // 冲/合/生/扶/克/平；冲合可与其他关系并存
}

func dayInteraction(lineZhi ganzhi.Zhi, riZhi ganzhi.Zhi) DayRelation {
	le := ganzhi.ZhiWuxing(lineZhi)
	de := ganzhi.ZhiWuxing(riZhi)

	rel := DayRelation{Relations: []DayRelationType{}}

	if ganzhi.IsLiuChong(lineZhi, riZhi) {
		rel.Relations = append(rel.Relations, DayChong)
	}
	if ganzhi.IsZhiHe(lineZhi, riZhi) {
		rel.Relations = append(rel.Relations, DayHe)
	}
	if ganzhi.Sheng(de, le) {
		rel.Relations = append(rel.Relations, DaySheng)
	}
	if le == de {
		rel.Relations = append(rel.Relations, DayFu)
	}
	if ganzhi.Ke(de, le) {
		rel.Relations = append(rel.Relations, DayKe)
	}
	if len(rel.Relations) == 0 {
		rel.Relations = append(rel.Relations, DayPing)
	}
	return rel
}

func (p *Chart) dayRelationHas(position int, types ...DayRelationType) bool {
	if position < 1 || position > 6 {
		return false
	}
	for _, relation := range p.DayRelations[position-1].Relations {
		for _, typ := range types {
			if relation == typ {
				return true
			}
		}
	}
	return false
}

func (r DayRelation) String() string {
	parts := make([]string, 0, len(r.Relations))
	for _, item := range r.Relations {
		parts = append(parts, DecodeDayRelationType(item))
	}
	return strings.Join(parts, "兼")
}

func chongZhi(z ganzhi.Zhi) ganzhi.Zhi {
	return ganzhi.ChongZhi(z)
}
