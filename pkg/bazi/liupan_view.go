package bazi

import (
	"encoding/json"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// Domain-owned output views of the bazi.liupan series. jsonschema tags are
// inert metadata read only when a consumer reflects the type for its schema.

type LiuNianResolvedView struct {
	LiuNianStandardTime string `json:"standard_time" jsonschema:"流年选择器时区标准时间（无 DST）。"`
	LiuNianStartDate    string `json:"start_date" jsonschema:"按当地标准日取的流年选择日期。"`
	LichunYear          int    `json:"lichun_year" jsonschema:"该日期所在的立春年。"`
	LichunStart         string `json:"lichun_start" jsonschema:"立春年起始公历日（含）。"`
	LichunEnd           string `json:"lichun_end" jsonschema:"下一立春日（不含）。"`
	BoundaryDay         bool   `json:"boundary_day" jsonschema:"输入日期是否为立春边界日（归入新年）。"`
}

type LiuNianOutputView struct {
	Year              int                                 `json:"year" jsonschema:"年。"`
	ShiShen           ganzhi.ShiShen                      `json:"shi_shen" jsonschema:"十神。"`
	NianGan           *ganzhi.Gan                         `json:"gan,omitempty" jsonschema:"流年天干"`
	NianZhi           *ganzhi.Zhi                         `json:"zhi,omitempty" jsonschema:"流年地支"`
	NianGanWuxing     *string                             `json:"gan_wuxing,omitempty" jsonschema:"流年天干五行。"`
	LiuNianNaYin      *string                             `json:"na_yin,omitempty" jsonschema:"流年干支纳音。"`
	ShengRi           *int                                `json:"sheng_ri,omitempty" jsonschema:"流年是否生扶日主：0/1。"`
	KeRi              *int                                `json:"ke_ri,omitempty" jsonschema:"流年是否克制日主：0/1。"`
	ShenSha           *[]SeriesShenshaItemView            `json:"shen_sha,omitempty" jsonschema:"流年神煞名称与类别（吉/凶/中性）；取象见 chart 层 shen_sha description。年支+日支双查；值年按流年地支查命局四柱"`
	FuYinFanYin       *[]LiuNianFuyinFanyinItemView       `json:"fu_yin_fan_yin,omitempty" jsonschema:"伏吟反吟"`
	DaYunInteractions *[]LiuNianDayunInteractionsItemView `json:"da_yun_interactions,omitempty" jsonschema:"查询年份所在大运与流年交互（未起运或大运过完时为空数组）"`
	NatalInteractions *[]LiuNianNatalInteractionsItemView `json:"natal_interactions,omitempty" jsonschema:"流年与命局交互"`
	NianJie           LiuNianResolvedView                 `json:"nian_jie" jsonschema:"流年立春年界。"`
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
	GanRels     []LiuNianDayunInteractionsItemGanRelsItemView `json:"gan_rels" jsonschema:"天干关系列表。"`
	PillarLabel *string                                       `json:"zhu_label" jsonschema:"命局柱位及干支标签。"`
	ZhiRels     []LiuNianDayunInteractionsItemZhiRelsItemView `json:"zhi_rels" jsonschema:"地支关系列表。"`
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
	GanRels     []SeriesGanRelationItemView `json:"gan_rels" jsonschema:"天干关系列表。"`
	PillarLabel *string                     `json:"zhu_label" jsonschema:"命局柱位及干支标签。"`
	ZhiRels     []SeriesZhiRelationItemView `json:"zhi_rels" jsonschema:"地支关系列表。"`
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

type LiuYueOutputView struct {
	Year         int                          `json:"year" jsonschema:"流月公历年。"`
	Month        int                          `json:"month" jsonschema:"流月公历月。"`
	ShiShen      ganzhi.ShiShen               `json:"shi_shen" jsonschema:"十神。"`
	YueGan       *ganzhi.Gan                  `json:"gan,omitempty" jsonschema:"流月天干"`
	YueZhi       *ganzhi.Zhi                  `json:"zhi,omitempty" jsonschema:"流月地支"`
	YueGanWuxing *string                      `json:"gan_wuxing,omitempty" jsonschema:"流月天干五行。"`
	ShengRi      *int                         `json:"sheng_ri,omitempty" jsonschema:"流月是否生扶日主：0/1。"`
	KeRi         *int                         `json:"ke_ri,omitempty" jsonschema:"流月是否克制日主：0/1。"`
	ShenSha      *[]SeriesShenshaItemView     `json:"shen_sha,omitempty" jsonschema:"流月神煞名称与类别；取象见 chart 层 shen_sha description"`
	GanRels      *[]SeriesGanRelationItemView `json:"gan_rels,omitempty" jsonschema:"天干关系"`
	ZhiRels      *[]SeriesZhiRelationItemView `json:"zhi_rels,omitempty" jsonschema:"地支关系"`
	JieQi        LiuYueJieQiView              `json:"jie_qi" jsonschema:"流月节气月界。"`
}

type LiuRiJieView struct {
	LiuRiStandardTime string `json:"standard_time" jsonschema:"流日选择器时区标准时间（无 DST）。"`
	LiuRiStandardDate string `json:"standard_date" jsonschema:"流日选择器的时区标准日（无 DST）。"`
}

type LiuRiOutputView struct {
	Date        string                       `json:"date" jsonschema:"流日日期 YYYY-MM-DD。"`
	ShiShen     ganzhi.ShiShen               `json:"shi_shen" jsonschema:"十神。"`
	RiGan       *ganzhi.Gan                  `json:"gan,omitempty" jsonschema:"流日天干"`
	RiZhi       *ganzhi.Zhi                  `json:"zhi,omitempty" jsonschema:"流日地支"`
	DayNayin    *string                      `json:"na_yin,omitempty" jsonschema:"流日干支纳音。"`
	Shensha     *[]SeriesShenshaItemView     `json:"shen_sha,omitempty" jsonschema:"流日神煞名称与类别；取象见 chart 层 shen_sha description"`
	GanRels     *[]SeriesGanRelationItemView `json:"gan_rels,omitempty" jsonschema:"天干关系"`
	ZhiRels     *[]SeriesZhiRelationItemView `json:"zhi_rels,omitempty" jsonschema:"地支关系"`
	DaYunRels   *[]SeriesZhiRelationItemView `json:"da_yun_rels,omitempty" jsonschema:"与大运关系"`
	LiuNianRels *[]LiuRiLiunianRelsItemView  `json:"liu_nian_rels,omitempty" jsonschema:"与流年关系"`
	RiJie       LiuRiJieView                 `json:"ri_jie" jsonschema:"流日日界。"`
}

type LiuRiLiunianRelsItemView struct {
	Explanation *string `json:"relation_explanation" jsonschema:"关系解释。"`
	Kind        *string `json:"relation_kind" jsonschema:"类型。"`
	ZhiA        *string `json:"zhi_a" jsonschema:"甲支。"`
	ZhiB        *string `json:"zhi_b" jsonschema:"乙支。"`
}

type LiuShiChenView struct {
	LiuShiStandardTime string     `json:"standard_time" jsonschema:"流时选择器时区标准时间（无 DST）。"`
	LiuShiStandardDate string     `json:"standard_date" jsonschema:"流时选择器的时区标准日（无 DST）。"`
	Shichen            ganzhi.Zhi `json:"shichen" jsonschema:"真太阳时时辰。"`
}

type LiuShiOutputView struct {
	LiuShiTime string                       `json:"hour_range" jsonschema:"流时时辰区间。"`
	ShiShen    ganzhi.ShiShen               `json:"shi_shen" jsonschema:"十神。"`
	ShiGan     *ganzhi.Gan                  `json:"gan,omitempty" jsonschema:"流时天干"`
	ShiZhi     *ganzhi.Zhi                  `json:"zhi,omitempty" jsonschema:"流时地支"`
	GanRels    *[]SeriesGanRelationItemView `json:"gan_rels,omitempty" jsonschema:"天干关系"`
	ZhiRels    *[]SeriesZhiRelationItemView `json:"zhi_rels,omitempty" jsonschema:"地支关系"`
	ShiChen    LiuShiChenView               `json:"shi_chen" jsonschema:"流时时辰定位。"`
}

type XiaoYunOutputView struct {
	XiaoYun []XiaoYunItemView `json:"xiao_yun" jsonschema:"各学派小运列表；长度 = count。"`
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

// viewRoundtrip maps a pure output series value into its tagged View contract.
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

func (v LiuNian) View() LiuNianOutputView  { return liupanView(v, &LiuNianOutputView{}) }
func (v LiuYue) View() LiuYueOutputView    { return liupanView(v, &LiuYueOutputView{}) }
func (v LiuRi) View() LiuRiOutputView      { return liupanView(v, &LiuRiOutputView{}) }
func (v LiuShi) View() LiuShiOutputView    { return liupanView(v, &LiuShiOutputView{}) }
func (v XiaoYunSet) View() XiaoYunItemView { return liupanView(v, &XiaoYunItemView{}) }
