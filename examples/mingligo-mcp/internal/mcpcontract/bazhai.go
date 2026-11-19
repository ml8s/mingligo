package mcpcontract

import (
	"context"
	"fmt"

	"github.com/ml8s/mingligo/pkg/bazhai"
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// BazhaiSubjectInput is the minimal calibrated identity used by Bazhai. It does
// not carry longitude because Bazhai does not use local apparent solar time.
type BazhaiSubjectInput struct {
	StandardTime StandardTime `json:"standard_time" jsonschema:"出生年代表标准时间（无 DST），RFC3339 带 offset；用于立春年界。"`
	Gender       GenderName   `json:"gender" jsonschema:"性别。"`
}

// BazhaiChartInput consumes the Bazhai subject.
type BazhaiChartInput struct {
	Subject BazhaiSubjectInput `json:"subject" jsonschema:"八宅命卦主体。"`
}

type BazhaiGuaOutput struct {
	Index   GuaNumber    `json:"index" jsonschema:"洛书宫数。"`
	Name    EightGuaName `json:"name" jsonschema:"卦名。"`
	Wuxing  ElementName  `json:"wuxing" jsonschema:"卦五行。"`
	YinYang YinYangName  `json:"yin_yang" jsonschema:"卦阴阳。"`
}

type BazhaiMingGuaOutput struct {
	Gua          BazhaiGuaOutput  `json:"gua" jsonschema:"命卦。"`
	Group        LifeGroupName    `json:"group" jsonschema:"东四命或西四命。"`
	YearBoundary YearBoundaryName `json:"year_boundary" jsonschema:"年界：立春。"`
}

type BazhaiDirectionsOutput struct {
	ShengQi []DirectionName `json:"sheng_qi" jsonschema:"生气方。"`
	TianYi  []DirectionName `json:"tian_yi" jsonschema:"天医方。"`
	YanNian []DirectionName `json:"yan_nian" jsonschema:"延年方。"`
	FuWei   []DirectionName `json:"fu_wei" jsonschema:"伏位方。"`
	HuoHai  []DirectionName `json:"huo_hai" jsonschema:"祸害方。"`
	WuGui   []DirectionName `json:"wu_gui" jsonschema:"五鬼方。"`
	LiuSha  []DirectionName `json:"liu_sha" jsonschema:"六煞方。"`
	JueMing []DirectionName `json:"jue_ming" jsonschema:"绝命方。"`
}

type BazhaiAnnualStarOutput struct {
	GongNum  GuaNumber        `json:"gong_num" jsonschema:"洛书宫数。"`
	Xing     FlyingStarNumber `json:"xing" jsonschema:"紫白星数。"`
	XingName FlyingStarName   `json:"xing_name" jsonschema:"紫白星名。"`
	Wuxing   ElementName      `json:"wuxing" jsonschema:"星五行。"`
	RuZhong  bool             `json:"ru_zhong" jsonschema:"是否入中。"`
}

type BazhaiYearStarsOutput struct {
	Year         BazhaiBirthYear           `json:"year" jsonschema:"立春年份。"`
	RuZhong      FlyingStarName            `json:"ru_zhong" jsonschema:"入中紫白星。"`
	YearBoundary YearBoundaryName          `json:"year_boundary" jsonschema:"年界：立春。"`
	Palaces      [9]BazhaiAnnualStarOutput `json:"gong_wei" jsonschema:"九宫紫白星。"`
}

type BazhaiChartOutput struct {
	MingGua    BazhaiMingGuaOutput    `json:"ming_gua" jsonschema:"命卦。"`
	BaZhaiDirs BazhaiDirectionsOutput `json:"ba_zhai_dirs" jsonschema:"四吉四凶方向。"`
	YearStars  BazhaiYearStarsOutput  `json:"birth_year_xing" jsonschema:"出生年紫白飞星。"`
}

type BazhaiLayoutInput struct {
	MingGua   EightGuaName `json:"ming_gua" jsonschema:"命卦。"`
	DoorGua   EightGuaName `json:"door_gua" jsonschema:"门卦。"`
	MasterGua EightGuaName `json:"master_gua" jsonschema:"主卧卦。"`
	StoveGua  EightGuaName `json:"stove_gua" jsonschema:"灶卦。"`
}

type BazhaiDoorStoveOutput struct {
	Direction DirectionName    `json:"direction" jsonschema:"方位。"`
	GuaName   EightGuaName     `json:"gua_name" jsonschema:"卦名。"`
	Wuxing    ElementName      `json:"wuxing" jsonschema:"卦五行。"`
	YouXing   YouxingName      `json:"youxing" jsonschema:"游年九星。"`
	Rating    BazhaiRatingName `json:"rating" jsonschema:"吉凶等级。"`
	Group     TrigramGroupName `json:"group" jsonschema:"东四卦或西四卦。"`
	Match     BazhaiMatchName  `json:"match" jsonschema:"与命卦同组吉，异组凶。"`
}

type BazhaiLayoutOutput struct {
	Group   LifeGroupName         `json:"group" jsonschema:"东四命或西四命。"`
	MingGua EightGuaName          `json:"ming_gua" jsonschema:"命卦。"`
	Door    BazhaiDoorStoveOutput `json:"door" jsonschema:"门位判断。"`
	Master  BazhaiDoorStoveOutput `json:"master" jsonschema:"主卧判断。"`
	Stove   BazhaiDoorStoveOutput `json:"stove" jsonschema:"灶位判断。"`
}

func BazhaiChart(_ context.Context, input BazhaiChartInput) (bazhai.Chart, error) {
	standard, _, err := parseStandardTime(input.Subject.StandardTime)
	if err != nil {
		return bazhai.Chart{}, fmt.Errorf("bazhai.chart: %w", err)
	}
	gender, err := ganzhi.ParseGender(string(input.Subject.Gender))
	if err != nil {
		return bazhai.Chart{}, fmt.Errorf("bazhai.chart: %w", err)
	}
	chart, err := bazhai.ComputeChartAtStandard(gender, standard)
	if err != nil {
		return bazhai.Chart{}, fmt.Errorf("bazhai.chart: %w", err)
	}
	return chart, nil
}

func BazhaiLayout(_ context.Context, input BazhaiLayoutInput) (bazhai.LayoutResult, error) {
	layout, err := bazhai.ComputeLayout(
		string(input.MingGua), string(input.DoorGua),
		string(input.MasterGua), string(input.StoveGua),
	)
	if err != nil {
		return bazhai.LayoutResult{}, err
	}
	return layout, nil
}

