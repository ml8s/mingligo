package liuyao

import (

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

func feiFuRelation(flying, hidden ganzhi.Wuxing) FeiFuRelation {
	switch {
	case flying == hidden:
		return FeiFuBiHe
	case ganzhi.Sheng(flying, hidden):
		return FeiFuShengFu
	case ganzhi.Ke(flying, hidden):
		return FeiFuKeFu
	case ganzhi.Sheng(hidden, flying):
		return FeiFuFuSheng
	case ganzhi.Ke(hidden, flying):
		return FeiFuFuKe
	default:
		return FeiFuNone
	}
}

func hiddenExitCondition(relation FeiFuRelation) string {
	switch relation {
	case FeiFuKeFu:
		return "待飞神受冲或合化，伏神旺相方可出伏"
	case FeiFuShengFu:
		return "飞神生扶，待值日或引拔而出"
	case FeiFuFuKe:
		return "伏神有力克飞，待冲开飞神"
	case FeiFuFuSheng:
		return "伏气外泄，待生扶值日"
	default:
		return "待值日或冲引而出"
	}
}

func computeHiddenLines(p *Chart) []HiddenLine {
	palaceBase := guaTable[p.BenGua].PalaceIdx * 8
	if palaceBase < 0 || palaceBase >= len(guaTable) {
		return nil
	}
	baseLines := zhuangGua(guaIndex(palaceBase), p.RiGan, false, p.PalaceWuxing)
	visible := map[LiuQin]bool{}
	for _, line := range p.Lines {
		visible[line.LiuQin] = true
	}

	var facts []HiddenLine
	for i := range baseLines {
		base := baseLines[i]
		if visible[base.LiuQin] {
			continue
		}
		flying := p.Lines[i]
		relation := feiFuRelation(flying.Wuxing, base.Wuxing)
		tombTypes := fuTombSources(p, base.Wuxing)
		facts = append(facts, HiddenLine{
			Position:       i + 1,
			LiuQin:         base.LiuQin.String(),
			Gan:            ganzhi.GanName(base.Gan),
			Zhi:            ganzhi.ZhiName(base.Zhi),
			Wuxing:         base.Wuxing.String(),
			FlyingPosition: flying.Position,
			FlyingLiuQin:   flying.LiuQin.String(),
			FlyingGan:      ganzhi.GanName(flying.Gan),
			FlyingZhi:      ganzhi.ZhiName(flying.Zhi),
			FlyingWuxing:   flying.Wuxing.String(),
			FeiFuRelation:  relation.String(),
			XunKong:        base.Zhi == p.XunKong[0] || base.Zhi == p.XunKong[1],
			YuePo:          ganzhi.IsLiuChong(base.Zhi, p.YueZhi),
			MuKu:           len(tombTypes) > 0,
			ExitCondition:  hiddenExitCondition(relation),
		})
	}
	return facts
}

func fuTombSources(p *Chart, element ganzhi.Wuxing) []MuKuSource {
	tomb := tombOf(element)
	if tomb == 0 {
		return nil
	}
	var types []MuKuSource
	if p.RiZhi == tomb {
		types = append(types, MuKuDay)
	}
	for _, position := range p.DongYao {
		if position >= 1 && position <= 6 && p.Lines[position-1].Zhi == tomb {
			types = append(types, MuKuMoving)
			break
		}
	}
	return types
}
