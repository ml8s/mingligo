package bazi

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// NayinRelEntry describes a nayin relation between two pillars.
type NayinRelEntry struct {
	A        string             `json:"a"`
	B        string             `json:"b"`
	Relation GanZhiRelationKind `json:"relation"`
}

// ChangShengStage describes one of the 12 life stages.
type ChangShengStage struct {
	Name string     `json:"name"`
	Zhi  ganzhi.Zhi `json:"index" jsonschema:"地支符号（json key 沿用 index 以兼容消费端）。"`
}

func withChartExtraFields(chart Chart, core Core) Chart {
	bz := core.ToBazi()

	stages := [12]ChangShengStage{}
	zhi := ganzhi.ChangShengTable[bz.Ri.Gan]
	for i := 0; i < 12; i++ {
		stages[i] = ChangShengStage{
			Name: ganzhi.StageNamesZH[i],
			Zhi:  zhi[i],
		}
	}

	nayins := core.NaYin
	var nayinRels []NayinRelEntry
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			ae := ganzhi.NayinWuxing(nayins[i])
			be := ganzhi.NayinWuxing(nayins[j])
			rel := relSame
			if ae != 0 && be != 0 && ae != be {
				if ganzhi.Sheng(ae, be) || ganzhi.Sheng(be, ae) {
					rel = relSheng
				} else {
					rel = relKe
				}
			}
			nayinRels = append(nayinRels, NayinRelEntry{
				A: zhuLabels[i], B: zhuLabels[j], Relation: rel,
			})
		}
	}

	chart.SanYuan = computeSanYuan(bz.Yue, bz.Nian.Gan, bz.Shi.Zhi)
	chart.TaiXi = computeTaiXi(bz.Ri)
	chart.GongJia = computeGongJia(bz)
	chart.NayinRel = nayinRels
	chart.ChangSheng = stages
	chart.SanQiName = DecodeSanQiKind(sanQiType(bz))
	return chart
}
