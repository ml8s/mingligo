package bazi

import (
	"encoding/json"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// Domain-owned output views of the bazi.liupan series. Each flow layer is
// aggregated by命理子域 (流期干支/日主作用/神煞/伏吟反吟/干支关系/历法界),
// not a flat bag of fields. jsonschema tags are inert metadata.

type LiuNianResolvedView struct {
	LiuNianStandardTime string `json:"standard_time" jsonschema:"流年选择器时区标准时间（无 DST）。"`
	LiuNianStartDate    string `json:"start_date" jsonschema:"按当地标准日取的流年选择日期。"`
	LichunYear          int    `json:"lichun_year" jsonschema:"该日期所在的立春年。"`
	LichunStart         string `json:"lichun_start" jsonschema:"立春年起始公历日（含）。"`
	LichunEnd           string `json:"lichun_end" jsonschema:"下一立春日（不含）。"`
	BoundaryDay         bool   `json:"boundary_day" jsonschema:"输入日期是否为立春边界日（归入新年）。"`
}

// LiuNianShenShaName is the closed 神煞 name set used by liupan series output.
type LiuNianShenShaName string

type SeriesShenshaItemView struct {
	Name        LiuNianShenShaName `json:"name" jsonschema:"名称。"`
	Category    string             `json:"category" jsonschema:"类别。"`
	Description *string            `json:"description,omitempty" jsonschema:"神煞取象说明。"`
}

type LiuNianFuyinFanyinItemView struct {
	Explanation *string `json:"relation_explanation" jsonschema:"关系解释。"`
	NatalIndex  *int    `json:"natal_index" jsonschema:"命局四柱索引：0 年、1 月、2 日、3 时。"`
	Kind        *string `json:"relation_kind" jsonschema:"类型。"`
}

type LiuNianDayunInteractionsItemView struct {
	GanRels  []LiuNianDayunInteractionsItemGanRelsItemView `json:"gan_rels" jsonschema:"天干关系列表。"`
	ZhuLabel *string                                       `json:"zhu_label" jsonschema:"参照柱（命局/大运）标签。"`
	ZhiRels  []LiuNianDayunInteractionsItemZhiRelsItemView `json:"zhi_rels" jsonschema:"地支关系列表。"`
}

type LiuNianDayunInteractionsItemGanRelsItemView struct {
	GanA     *ganzhi.Gan `json:"gan_a" jsonschema:"甲干。"`
	GanB     *ganzhi.Gan `json:"gan_b" jsonschema:"乙干。"`
	Relation *string     `json:"relation" jsonschema:"关系。"`
	Kind     *string     `json:"relation_kind" jsonschema:"类型。"`
}

type LiuNianDayunInteractionsItemZhiRelsItemView struct {
	Explanation *string     `json:"relation_explanation" jsonschema:"关系解释。"`
	Kind        *string     `json:"relation_kind" jsonschema:"类型。"`
	ZhiA        *ganzhi.Zhi `json:"zhi_a" jsonschema:"甲支。"`
	ZhiB        *ganzhi.Zhi `json:"zhi_b" jsonschema:"乙支。"`
}

type LiuNianNatalInteractionsItemView struct {
	GanRels  []SeriesGanRelationItemView `json:"gan_rels" jsonschema:"天干关系列表。"`
	ZhuLabel *string                     `json:"zhu_label" jsonschema:"参照柱（命局/大运）标签。"`
	ZhiRels  []SeriesZhiRelationItemView `json:"zhi_rels" jsonschema:"地支关系列表。"`
}

type SeriesGanRelationItemView struct {
	GanA     *string     `json:"gan_a" jsonschema:"甲干。"`
	GanB     *ganzhi.Gan `json:"gan_b" jsonschema:"乙干。"`
	Relation *string     `json:"relation" jsonschema:"关系。"`
	Kind     *string     `json:"relation_kind" jsonschema:"类型。"`
}

type SeriesZhiRelationItemView struct {
	Explanation *string     `json:"relation_explanation" jsonschema:"关系解释。"`
	Kind        *string     `json:"relation_kind" jsonschema:"类型。"`
	ZhiA        *string     `json:"zhi_a" jsonschema:"甲支。"`
	ZhiB        *ganzhi.Zhi `json:"zhi_b" jsonschema:"乙支。"`
}

type LiuYueJieQiView struct {
	LiuYueStandardTime string `json:"standard_time" jsonschema:"流月选择器时区标准时间（无 DST）。"`
	LiuYueStartDate    string `json:"start_date" jsonschema:"按当地标准日取的流月选择日期。"`
	JieqiStart         string `json:"jieqi_start" jsonschema:"节气月起日（节，含）。"`
	JieqiEnd           string `json:"jieqi_end" jsonschema:"下一节（不含）。"`
	BoundaryDay        bool   `json:"boundary_day" jsonschema:"输入日期是否为节边界日（归入新月）。"`
}

type LiuRiJieView struct {
	LiuRiStandardTime string `json:"standard_time" jsonschema:"流日选择器时区标准时间（无 DST）。"`
	LiuRiStandardDate string `json:"standard_date" jsonschema:"流日选择器的时区标准日（无 DST）。"`
}

// zhuInteractionView is the two-layer (gan+zhi) interaction with one reference
// pillar; shared by flow-day/flow-hour upper-layer facts.
type zhuInteractionView struct {
	ZhuLabel string                      `json:"zhu_label" jsonschema:"参照柱标签。"`
	GanRels  []SeriesGanRelationItemView `json:"gan_rels" jsonschema:"天干关系。"`
	ZhiRels  []SeriesZhiRelationItemView `json:"zhi_rels" jsonschema:"地支关系。"`
}

type LiuShiChenView struct {
	LiuShiStandardTime string     `json:"standard_time" jsonschema:"流时选择器时区标准时间（无 DST）。"`
	LiuShiStandardDate string     `json:"standard_date" jsonschema:"流时选择器的时区标准日（无 DST）。"`
	Shichen            ganzhi.Zhi `json:"shichen" jsonschema:"真太阳时时辰。"`
}

type XiaoYunItemView struct {
	Method string                    `json:"method" jsonschema:"流派标识"`
	Basis  string                    `json:"basis" jsonschema:"起法依据"`
	Zhus   []XiaoYunItemZhusItemView `json:"zhus" jsonschema:"该学派小运柱列表。"`
}

type XiaoYunItemZhusItemView struct {
	Age     int             `json:"age" jsonschema:"虚岁。"`
	Gan     ganzhi.Gan      `json:"gan" jsonschema:"天干。"`
	Zhi     ganzhi.Zhi      `json:"zhi" jsonschema:"地支。"`
	Name    string          `json:"name" jsonschema:"名称。"`
	ShiShen *ganzhi.ShiShen `json:"shi_shen,omitempty" jsonschema:"十神。"`
}

// LiuQiView groups 流期干支 facts (干/支/五行/纳音/十神).
type LiuQiView struct {
	Gan       *ganzhi.Gan    `json:"gan,omitempty" jsonschema:"流期天干。"`
	Zhi       *ganzhi.Zhi    `json:"zhi,omitempty" jsonschema:"流期地支。"`
	GanWuxing *string        `json:"gan_wuxing,omitempty" jsonschema:"流期天干五行。"`
	NaYin     *string        `json:"na_yin,omitempty" jsonschema:"流期干支纳音。"`
	ShiShen   ganzhi.ShiShen `json:"shi_shen" jsonschema:"流期十神。"`
}

// RiZhuZuoYongView groups 流期对日主的生扶克制 (命理事实).
type RiZhuZuoYongView struct {
	Sheng int `json:"sheng" jsonschema:"流期是否生扶日主：0/1。"`
	Ke    int `json:"ke" jsonschema:"流期是否克制/耗泄日主：0/1。"`
}

// GuanXiView groups 流期与命局/大运/流年的干支关系.
type LiuPanGuanXiView struct {
	GanRels             *[]SeriesGanRelationItemView `json:"gan_rels,omitempty" jsonschema:"流期与命局天干关系。"`
	ZhiRels             *[]SeriesZhiRelationItemView `json:"zhi_rels,omitempty" jsonschema:"流期与命局地支关系。"`
	DaYunInteractions   *[]zhuInteractionView        `json:"da_yun_interactions,omitempty" jsonschema:"与大运的干支交互。"`
	LiuNianInteractions *[]zhuInteractionView        `json:"liu_nian_interactions,omitempty" jsonschema:"与流年的干支交互。"`
	NatalInteractions   *[]zhuInteractionView        `json:"natal_interactions,omitempty" jsonschema:"与命局的干支交互。"`
	LiuRiInteraction    *[]zhuInteractionView        `json:"liu_ri_interaction,omitempty" jsonschema:"与流日的干支交互。"`
}

// LiFaJieView groups 历法边界 (年/年界/节气界/日界/时辰定位).
type LiFaJieView struct {
	Year      int                  `json:"year,omitempty" jsonschema:"流期公历年。"`
	Month     int                  `json:"month,omitempty" jsonschema:"流期公历月。"`
	Date      string               `json:"date,omitempty" jsonschema:"流期公历日。"`
	HourRange string               `json:"hour_range,omitempty" jsonschema:"流时时辰区间。"`
	NianJie   *LiuNianResolvedView `json:"nian_jie,omitempty" jsonschema:"立春年界。"`
	JieQi     *LiuYueJieQiView     `json:"jie_qi,omitempty" jsonschema:"节气月界。"`
	RiJie     *LiuRiJieView        `json:"ri_jie,omitempty" jsonschema:"流日日界。"`
	ShiChen   *LiuShiChenView      `json:"shi_chen,omitempty" jsonschema:"真太阳时辰。"`
}

// LiuNianOutputView aggregates 流年 by命理子域.
type LiuNianOutputView struct {
	LiuQi        LiuQiView                     `json:"liu_qi" jsonschema:"流年干支。"`
	RiZhuZuoYong RiZhuZuoYongView              `json:"ri_zhu_zuo_yong" jsonschema:"生扶克制日主。"`
	ShenSha      *[]SeriesShenshaItemView      `json:"shen_sha,omitempty" jsonschema:"流年神煞。"`
	FuYinFanYin  *[]LiuNianFuyinFanyinItemView `json:"fu_yin_fan_yin,omitempty" jsonschema:"伏吟反吟。"`
	GuanXi       LiuPanGuanXiView              `json:"guan_xi" jsonschema:"干支关系。"`
	LiFaJie      LiFaJieView                   `json:"li_fa_jie" jsonschema:"立春年界。"`
}

// LiuYueOutputView aggregates 流月 by命理子域.
type LiuYueOutputView struct {
	LiuQi        LiuQiView                `json:"liu_qi" jsonschema:"流月干支。"`
	RiZhuZuoYong RiZhuZuoYongView         `json:"ri_zhu_zuo_yong" jsonschema:"生扶克制日主。"`
	ShenSha      *[]SeriesShenshaItemView `json:"shen_sha,omitempty" jsonschema:"流月神煞。"`
	GuanXi       LiuPanGuanXiView         `json:"guan_xi" jsonschema:"干支关系。"`
	LiFaJie      LiFaJieView              `json:"li_fa_jie" jsonschema:"节气月界。"`
}

// LiuRiOutputView aggregates 流日 by命理子域.
type LiuRiOutputView struct {
	LiuQi        LiuQiView                `json:"liu_qi" jsonschema:"流日干支。"`
	RiZhuZuoYong RiZhuZuoYongView         `json:"ri_zhu_zuo_yong" jsonschema:"生扶克制日主。"`
	ShenSha      *[]SeriesShenshaItemView `json:"shen_sha,omitempty" jsonschema:"流日神煞。"`
	GuanXi       LiuPanGuanXiView         `json:"guan_xi" jsonschema:"干支关系。"`
	LiFaJie      LiFaJieView              `json:"li_fa_jie" jsonschema:"流日日界。"`
}

// LiuShiOutputView aggregates 流时 by命理子域.
type LiuShiOutputView struct {
	LiuQi        LiuQiView        `json:"liu_qi" jsonschema:"流时干支。"`
	RiZhuZuoYong RiZhuZuoYongView `json:"ri_zhu_zuo_yong" jsonschema:"生扶克制日主。"`
	GuanXi       LiuPanGuanXiView `json:"guan_xi" jsonschema:"干支关系。"`
	LiFaJie      LiFaJieView      `json:"li_fa_jie" jsonschema:"真太阳时辰。"`
}

func encodeShiShenName(name string) ganzhi.ShiShen {
	for v := ganzhi.ShiShen(1); v <= 10; v++ {
		if ganzhi.DecodeShiShen(v) == name {
			return v
		}
	}
	return 0
}

func liuQiView(gan ganzhi.Gan, zhi ganzhi.Zhi, wuxing, naYin string, shiShen ganzhi.ShiShen) LiuQiView {
	g, z := gan, zhi
	w, n := wuxing, naYin
	return LiuQiView{
		Gan: &g, Zhi: &z,
		GanWuxing: &w, NaYin: &n,
		ShiShen: shiShen,
	}
}

func liupanView[T any](src any, dst *T) T {
	raw, err := json.Marshal(src)
	if err != nil {
		return *dst
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return *dst
	}
	return *dst
}

func zhiInteractions(v []zhuInteraction) *[]zhuInteractionView {
	if v == nil {
		return nil
	}
	out := make([]zhuInteractionView, 0, len(v))
	for _, item := range v {
		out = append(out, liupanView(item, &zhuInteractionView{}))
	}
	return &out
}

func oneZhiInteraction(v zhuInteraction) *[]zhuInteractionView {
	if v.ZhuLabel == "" && len(v.GanRels) == 0 && len(v.ZhiRels) == 0 {
		return nil
	}
	r := liupanView(v, &zhuInteractionView{})
	return &[]zhuInteractionView{r}
}

func shenShaView(v []ShenShaProjection) *[]SeriesShenshaItemView {
	if v == nil {
		return nil
	}
	out := make([]SeriesShenshaItemView, 0, len(v))
	for _, item := range v {
		out = append(out, liupanView(item, &SeriesShenshaItemView{}))
	}
	return &out
}

func fuYinFanYinView(v []FuYinFanYin) *[]LiuNianFuyinFanyinItemView {
	if v == nil {
		return nil
	}
	out := make([]LiuNianFuyinFanyinItemView, 0, len(v))
	for _, item := range v {
		out = append(out, liupanView(item, &LiuNianFuyinFanyinItemView{}))
	}
	return &out
}

func nianJieView(v LiuNianJie) *LiuNianResolvedView {
	r := liupanView(v, &LiuNianResolvedView{})
	return &r
}
func yueJieView(v LiuYueJieQi) *LiuYueJieQiView {
	r := liupanView(v, &LiuYueJieQiView{})
	return &r
}
func riJieView(v LiuRiJie) *LiuRiJieView {
	r := liupanView(v, &LiuRiJieView{})
	return &r
}
func shiChenView(v LiuShiChen) *LiuShiChenView {
	r := liupanView(v, &LiuShiChenView{})
	return &r
}

func zuoYongView(sheng, ke int) RiZhuZuoYongView {
	return RiZhuZuoYongView{Sheng: sheng, Ke: ke}
}

// View aggregates 流年 by命理子域.
func (v LiuNian) View() LiuNianOutputView {
	return LiuNianOutputView{
		LiuQi:        liuQiView(v.NianGan, v.NianZhi, v.NianGanWuxing, v.LiuNianNaYin, encodeShiShenName(v.ShiShen)),
		RiZhuZuoYong: zuoYongView(v.ShengRi, v.KeRi),
		ShenSha:      shenShaView(v.ShenSha),
		FuYinFanYin:  fuYinFanYinView(v.FuYinFanYin),
		GuanXi:       LiuPanGuanXiView{NatalInteractions: zhiInteractions(v.NatalInteractions), DaYunInteractions: zhiInteractions(v.DaYunInteractions)},
		LiFaJie:      LiFaJieView{Year: v.Year, NianJie: nianJieView(v.NianJie)},
	}
}

// View aggregates 流月 by命理子域.
func (v LiuYue) View() LiuYueOutputView {
	var ganRels *[]SeriesGanRelationItemView
	if v.GanRels != nil {
		out := make([]SeriesGanRelationItemView, 0, len(v.GanRels))
		for _, item := range v.GanRels {
			out = append(out, liupanView(item, &SeriesGanRelationItemView{}))
		}
		ganRels = &out
	}
	var zhiRels *[]SeriesZhiRelationItemView
	if v.ZhiRels != nil {
		out := make([]SeriesZhiRelationItemView, 0, len(v.ZhiRels))
		for _, item := range v.ZhiRels {
			out = append(out, liupanView(item, &SeriesZhiRelationItemView{}))
		}
		zhiRels = &out
	}
	return LiuYueOutputView{
		LiuQi:        liuQiView(v.YueGan, v.YueZhi, v.YueGanWuxing, "", encodeShiShenName(v.ShiShen)),
		RiZhuZuoYong: zuoYongView(v.ShengRi, v.KeRi),
		ShenSha:      shenShaView(v.ShenSha),
		GuanXi: LiuPanGuanXiView{
			GanRels: ganRels, ZhiRels: zhiRels,
			DaYunInteractions:   oneZhiInteraction(v.DaYunInteraction),
			LiuNianInteractions: oneZhiInteraction(v.LiuNianInteraction),
		},
		LiFaJie: LiFaJieView{Year: v.Year, Month: v.Month, JieQi: yueJieView(v.JieQi)},
	}
}

// View aggregates 流日 by命理子域.
func (v LiuRi) View() LiuRiOutputView {
	var ganRels *[]SeriesGanRelationItemView
	if v.GanRels != nil {
		out := make([]SeriesGanRelationItemView, 0, len(v.GanRels))
		for _, item := range v.GanRels {
			out = append(out, liupanView(item, &SeriesGanRelationItemView{}))
		}
		ganRels = &out
	}
	var zhiRels *[]SeriesZhiRelationItemView
	if v.ZhiRels != nil {
		out := make([]SeriesZhiRelationItemView, 0, len(v.ZhiRels))
		for _, item := range v.ZhiRels {
			out = append(out, liupanView(item, &SeriesZhiRelationItemView{}))
		}
		zhiRels = &out
	}
	return LiuRiOutputView{
		LiuQi:        liuQiView(v.RiGan, v.RiZhi, "", v.DayNaYin, encodeShiShenName(v.ShiShen)),
		RiZhuZuoYong: zuoYongView(v.ShengRi, v.KeRi),
		ShenSha:      shenShaView(v.ShenSha),
		GuanXi: LiuPanGuanXiView{
			GanRels: ganRels, ZhiRels: zhiRels,
			DaYunInteractions:   oneZhiInteraction(v.DaYunInteraction),
			LiuNianInteractions: oneZhiInteraction(v.LiuNianInteraction),
		},
		LiFaJie: LiFaJieView{Date: v.Date, RiJie: riJieView(v.RiJie)},
	}
}

// View aggregates 流时 by命理子域.
func (v LiuShi) View() LiuShiOutputView {
	var ganRels *[]SeriesGanRelationItemView
	if v.GanRels != nil {
		out := make([]SeriesGanRelationItemView, 0, len(v.GanRels))
		for _, item := range v.GanRels {
			out = append(out, liupanView(item, &SeriesGanRelationItemView{}))
		}
		ganRels = &out
	}
	var zhiRels *[]SeriesZhiRelationItemView
	if v.ZhiRels != nil {
		out := make([]SeriesZhiRelationItemView, 0, len(v.ZhiRels))
		for _, item := range v.ZhiRels {
			out = append(out, liupanView(item, &SeriesZhiRelationItemView{}))
		}
		zhiRels = &out
	}
	return LiuShiOutputView{
		LiuQi:        liuQiView(v.ShiGan, v.ShiZhi, "", "", encodeShiShenName(v.ShiShen)),
		RiZhuZuoYong: zuoYongView(v.ShengRi, v.KeRi),
		GuanXi: LiuPanGuanXiView{
			GanRels: ganRels, ZhiRels: zhiRels,
			DaYunInteractions:   oneZhiInteraction(v.DaYunInteraction),
			LiuNianInteractions: oneZhiInteraction(v.LiuNianInteraction),
			LiuRiInteraction:    oneZhiInteraction(v.LiuRiInteraction),
		},
		LiFaJie: LiFaJieView{HourRange: v.Time, ShiChen: shiChenView(v.ShiChen)},
	}
}

func (v XiaoYunSet) View() XiaoYunItemView { return liupanView(v, &XiaoYunItemView{}) }
