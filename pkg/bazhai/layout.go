package bazhai

import (
	"fmt"
	"github.com/ml8s/mingligo/pkg/ganzhi"

	"github.com/ml8s/mingligo/pkg/luoshu"
)

// ── 门主灶判断 ──

type LayoutResult struct {
	Group   LifeGroup     `json:"group"` // 东四命/西四命
	MingGua string        `json:"ming_gua"`
	Door    DoorStoveInfo `json:"men"`
	Master  DoorStoveInfo `json:"zhu"`
	Stove   DoorStoveInfo `json:"zao"`
}

// DoorStoveInfo is one door/master/stove position assessment.
type DoorStoveInfo struct {
	Direction string       `json:"direction"`
	GuaName   string       `json:"gua_name"`
	Wuxing    string       `json:"wuxing"`
	YouXing   YouXing      `json:"youxing"`
	Rating    Rating       `json:"rating"`
	Group     TrigramGroup `json:"group"` // 东四卦/西四卦
	Match     Match        `json:"match"` // 吉/凶(与命卦同组不同组)
}

// ComputeLayout analyzes 门主灶 in八宅风水. Trigram inputs are domain symbols;
// display text is only produced at the boundary.
func ComputeLayout(mingGua, doorGua, masterGua, stoveGua luoshu.Trigram) (LayoutResult, error) {
	inputs := []struct {
		name    string
		trigram luoshu.Trigram
	}{
		{"ming_gua", mingGua}, {"door_gua", doorGua},
		{"master_gua", masterGua}, {"stove_gua", stoveGua},
	}
	nums := make([]int, len(inputs))
	for i, in := range inputs {
		nums[i] = luoshu.TrigramNumber(in.trigram)
		if nums[i] == 0 {
			return LayoutResult{}, fmt.Errorf("invalid %s %q", in.name, in.trigram)
		}
	}
	mg := nums[0]
	mgGroup := LifeGroupEast
	if westGroup[mg] {
		mgGroup = LifeGroupWest
	}
	return LayoutResult{
		Group:   mgGroup,
		MingGua: string(mingGua),
		Door:    evalPosition(nums[1], mg),
		Master:  evalPosition(nums[2], mg),
		Stove:   evalPosition(nums[3], mg),
	}, nil
}

func evalPosition(guaNum, mingGua int) DoorStoveInfo {
	group := TrigramGroupEast
	if westGroup[guaNum] {
		group = TrigramGroupWest
	}
	match := MatchBad
	if westGroup[guaNum] == westGroup[mingGua] {
		match = MatchLucky
	}
	youxing, rating := youxingForGua(mingGua, guaNum)
	return DoorStoveInfo{
		Direction: luoshu.PalaceTable[guaNum].Direction,
		GuaName:   luoshu.PalaceTable[guaNum].Name,
		Wuxing:    ganzhi.DecodeWuxing(luoshu.PalaceTable[guaNum].Element),
		YouXing:   youxing,
		Rating:    rating,
		Group:     group,
		Match:     match,
	}
}

func youxingForGua(mingGua, guaNum int) (YouXing, Rating) {
	p, ok := eightMansionPatterns[mingGua]
	if !ok {
		return "", ""
	}
	var youxing YouXing
	switch guaNum {
	case p.shengQi:
		youxing = YouXingShengQi
	case p.tianYi:
		youxing = YouXingTianYi
	case p.yanNian:
		youxing = YouXingYanNian
	case p.fuWei:
		youxing = YouXingFuWei
	case p.huoHai:
		youxing = YouXingHuoHai
	case p.wuGui:
		youxing = YouXingWuGui
	case p.liuSha:
		youxing = YouXingLiuSha
	case p.jueMing:
		youxing = YouXingJueMing
	default:
		return "", ""
	}
	return youxing, youxingRatings[youxing]
}
