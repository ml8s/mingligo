package mcpcontract

import "github.com/ml8s/mingligo/pkg/bazi"

// bazi.bond wire contract — bazi_chart now uses the domain-owned bazi.ChartView.

// ── bazi.bond input ──

type BaziBondInput struct {
	A BaziBondSideInput `json:"a" jsonschema:"甲方：已归一化的出生主体。"`
	B BaziBondSideInput `json:"b" jsonschema:"乙方：已归一化的出生主体。"`
}

type BaziBondSideInput struct {
	Subject NatalSubject `json:"subject" jsonschema:"已归一化的出生主体。"`
}

// ── bazi.bond output ──

// ── bazi.liupan input ──

type BaziLiupanInput struct {
	Scope        BaziLiupanScopeName `json:"scope" jsonschema:"流盘层级。"`
	Subject      NatalSubject        `json:"subject" jsonschema:"已归一化的出生主体。"`
	StandardTime StandardTime        `json:"standard_time,omitempty" jsonschema:"流盘选择起点；liu_nian/liu_yue/liu_ri/liu_shi 必填，local standard time，RFC3339 fixed offset。"`
	Count        *SeriesCount        `json:"count,omitempty" jsonschema:"连续查询数量。"`
	Longitude    *Longitude          `json:"longitude,omitempty" jsonschema:"地点经度；liu_ri/liu_shi 必填（-180..180），真太阳时校正用于日柱/时柱。"`
}

// ── bazi.liupan liu_nian output ──

type BaziLiuNianResolved struct {
	LiuNianStandardTime StandardTime `json:"standard_time" jsonschema:"流年选择器时区标准时间（无 DST）。"`
	LiuNianStartDate    DateString   `json:"start_date" jsonschema:"按当地标准日取的流年选择日期。"`
	LichunYear          int          `json:"lichun_year" jsonschema:"该日期所在的立春年。"`
	LichunStart         DateString   `json:"lichun_start" jsonschema:"立春年起始公历日（含）。"`
	LichunEnd           DateString   `json:"lichun_end" jsonschema:"下一立春日（不含）。"`
	BoundaryDay         bool         `json:"boundary_day" jsonschema:"输入日期是否为立春边界日（归入新年）。"`
}

type BaziLiuNianOutput struct {
	Year              int                                 `json:"year" jsonschema:"年。"`
	ShiShen           TenGodName                          `json:"shi_shen" jsonschema:"十神。"`
	NianGan           *GanName                            `json:"gan,omitempty" jsonschema:"流年天干"`
	NianZhi           *ZhiName                            `json:"zhi,omitempty" jsonschema:"流年地支"`
	NianGanWuxing     *string                             `json:"gan_wuxing,omitempty" jsonschema:"流年天干五行。"`
	LiuNianNaYin      *string                             `json:"na_yin,omitempty" jsonschema:"流年干支纳音。"`
	ShengRi           *int                                `json:"sheng_ri,omitempty" jsonschema:"流年是否生扶日主：0/1。"`
	KeRi              *int                                `json:"ke_ri,omitempty" jsonschema:"流年是否克制日主：0/1。"`
	ShenSha           *[]BaziSeriesShenshaItem            `json:"shen_sha,omitempty" jsonschema:"流年神煞名称与类别（吉/凶/中性）；取象见 chart 层 shen_sha description。年支+日支双查；值年按流年地支查命局四柱"`
	FuYinFanYin       *[]BaziLiuNianFuyinFanyinItem       `json:"fu_yin_fan_yin,omitempty" jsonschema:"伏吟反吟"`
	DaYunInteractions *[]BaziLiuNianDayunInteractionsItem `json:"da_yun_interactions,omitempty" jsonschema:"查询年份所在大运与流年交互（未起运或大运过完时为空数组）"`
	NatalInteractions *[]BaziLiuNianNatalInteractionsItem `json:"natal_interactions,omitempty" jsonschema:"流年与命局交互"`
	NianJie           BaziLiuNianResolved                 `json:"nian_jie" jsonschema:"流年立春年界。"`
}

type BaziSeriesShenshaItem struct {
	Name        ShenshaName         `json:"name" jsonschema:"名称。"`
	Category    ShenShaCategoryName `json:"category" jsonschema:"类别。"`
	Description *string             `json:"description,omitempty" jsonschema:"神煞取象说明。"`
}

type BaziLiuNianFuyinFanyinItem struct {
	Explanation *string `json:"relation_explanation" jsonschema:"关系解释。"`
	NatalIndex  *int    `json:"natal_index" jsonschema:"命局四柱索引：0 年、1 月、2 日、3 时。"`
	Kind        *string `json:"relation_kind" jsonschema:"类型。"`
}

type BaziLiuNianDayunInteractionsItem struct {
	GanRels     []BaziLiuNianDayunInteractionsItemGanRelsItem `json:"gan_rels" jsonschema:"天干关系列表。"`
	PillarLabel *string                                       `json:"zhu_label" jsonschema:"命局柱位及干支标签。"`
	ZhiRels     []BaziLiuNianDayunInteractionsItemZhiRelsItem `json:"zhi_rels" jsonschema:"地支关系列表。"`
}

type BaziLiuNianDayunInteractionsItemGanRelsItem struct {
	GanA     *GanNameNullable `json:"gan_a" jsonschema:"甲干。"`
	GanB     *GanNameNullable `json:"gan_b" jsonschema:"乙干。"`
	Relation *string          `json:"relation" jsonschema:"关系。"`
	Kind     *string          `json:"relation_kind" jsonschema:"类型。"`
}

type BaziLiuNianDayunInteractionsItemZhiRelsItem struct {
	Explanation *string          `json:"relation_explanation" jsonschema:"关系解释。"`
	Kind        *string          `json:"relation_kind" jsonschema:"类型。"`
	ZhiA        *ZhiNameNullable `json:"zhi_a" jsonschema:"甲支。"`
	ZhiB        *ZhiNameNullable `json:"zhi_b" jsonschema:"乙支。"`
}

type BaziLiuNianNatalInteractionsItem struct {
	GanRels     []BaziSeriesGanRelationItem `json:"gan_rels" jsonschema:"天干关系列表。"`
	PillarLabel *string                     `json:"zhu_label" jsonschema:"命局柱位及干支标签。"`
	ZhiRels     []BaziSeriesZhiRelationItem `json:"zhi_rels" jsonschema:"地支关系列表。"`
}

type BaziSeriesGanRelationItem struct {
	GanA     *string          `json:"gan_a" jsonschema:"甲干。"`
	GanB     *GanNameNullable `json:"gan_b" jsonschema:"乙干。"`
	Relation *string          `json:"relation" jsonschema:"关系。"`
	Kind     *string          `json:"relation_kind" jsonschema:"类型。"`
}

type BaziSeriesZhiRelationItem struct {
	Explanation *string          `json:"relation_explanation" jsonschema:"关系解释。"`
	Kind        *string          `json:"relation_kind" jsonschema:"类型。"`
	ZhiA        *string          `json:"zhi_a" jsonschema:"甲支。"`
	ZhiB        *ZhiNameNullable `json:"zhi_b" jsonschema:"乙支。"`
}

// ── bazi.liupan liu_yue output ──

type BaziLiuYueJieQi struct {
	LiuYueStandardTime StandardTime `json:"standard_time" jsonschema:"流月选择器时区标准时间（无 DST）。"`
	LiuYueStartDate    DateString   `json:"start_date" jsonschema:"按当地标准日取的流月选择日期。"`
	JieqiStart         DateString   `json:"jieqi_start" jsonschema:"节气月起日（节，含）。"`
	JieqiEnd           DateString   `json:"jieqi_end" jsonschema:"下一节（不含）。"`
	BoundaryDay        bool         `json:"boundary_day" jsonschema:"输入日期是否为节边界日（归入新月）。"`
}

type BaziLiuYueOutput struct {
	Year         int                          `json:"year" jsonschema:"流月公历年。"`
	Month        int                          `json:"month" jsonschema:"流月公历月。"`
	ShiShen      TenGodName                   `json:"shi_shen" jsonschema:"十神。"`
	YueGan       *GanName                     `json:"gan,omitempty" jsonschema:"流月天干"`
	YueZhi       *ZhiName                     `json:"zhi,omitempty" jsonschema:"流月地支"`
	YueGanWuxing *string                      `json:"gan_wuxing,omitempty" jsonschema:"流月天干五行。"`
	ShengRi      *int                         `json:"sheng_ri,omitempty" jsonschema:"流月是否生扶日主：0/1。"`
	KeRi         *int                         `json:"ke_ri,omitempty" jsonschema:"流月是否克制日主：0/1。"`
	ShenSha      *[]BaziSeriesShenshaItem     `json:"shen_sha,omitempty" jsonschema:"流月神煞名称与类别；取象见 chart 层 shen_sha description"`
	GanRels      *[]BaziSeriesGanRelationItem `json:"gan_rels,omitempty" jsonschema:"天干关系"`
	ZhiRels      *[]BaziSeriesZhiRelationItem `json:"zhi_rels,omitempty" jsonschema:"地支关系"`
	JieQi        BaziLiuYueJieQi              `json:"jie_qi" jsonschema:"流月节气月界。"`
}

type BaziLiuRiJie struct {
	LiuRiStandardTime StandardTime `json:"standard_time" jsonschema:"流日选择器时区标准时间（无 DST）。"`
	LiuRiStandardDate DateString   `json:"standard_date" jsonschema:"流日选择器的时区标准日（无 DST）。"`
}

// ── bazi.liupan liu_ri output ──

type BaziLiuRiOutput struct {
	Date        string                       `json:"date" jsonschema:"流日日期 YYYY-MM-DD。"`
	ShiShen     TenGodName                   `json:"shi_shen" jsonschema:"十神。"`
	RiGan       *GanName                     `json:"gan,omitempty" jsonschema:"流日天干"`
	RiZhi       *ZhiName                     `json:"zhi,omitempty" jsonschema:"流日地支"`
	DayNayin    *string                      `json:"na_yin,omitempty" jsonschema:"流日干支纳音。"`
	Shensha     *[]BaziSeriesShenshaItem     `json:"shen_sha,omitempty" jsonschema:"流日神煞名称与类别；取象见 chart 层 shen_sha description"`
	GanRels     *[]BaziSeriesGanRelationItem `json:"gan_rels,omitempty" jsonschema:"天干关系"`
	ZhiRels     *[]BaziSeriesZhiRelationItem `json:"zhi_rels,omitempty" jsonschema:"地支关系"`
	DaYunRels   *[]BaziSeriesZhiRelationItem `json:"da_yun_rels,omitempty" jsonschema:"与大运关系"`
	LiuNianRels *[]BaziLiuRiLiunianRelsItem  `json:"liu_nian_rels,omitempty" jsonschema:"与流年关系"`
	RiJie       BaziLiuRiJie                 `json:"ri_jie" jsonschema:"流日日界。"`
}

type BaziLiuRiLiunianRelsItem struct {
	Explanation *string `json:"relation_explanation" jsonschema:"关系解释。"`
	Kind        *string `json:"relation_kind" jsonschema:"类型。"`
	ZhiA        *string `json:"zhi_a" jsonschema:"甲支。"`
	ZhiB        *string `json:"zhi_b" jsonschema:"乙支。"`
}

// ── bazi.liupan liu_shi output ──

type BaziLiuShiChen struct {
	LiuShiStandardTime StandardTime `json:"standard_time" jsonschema:"流时选择器时区标准时间（无 DST）。"`
	LiuShiStandardDate DateString   `json:"standard_date" jsonschema:"流时选择器的时区标准日（无 DST）。"`
	Shichen            ZhiName      `json:"shichen" jsonschema:"真太阳时时辰。"`
}

type BaziLiuShiOutput struct {
	LiuShiTime string                       `json:"hour_range" jsonschema:"流时时辰区间。"`
	ShiShen    TenGodName                   `json:"shi_shen" jsonschema:"十神。"`
	ShiGan     *GanName                     `json:"gan,omitempty" jsonschema:"流时天干"`
	ShiZhi     *ZhiName                     `json:"zhi,omitempty" jsonschema:"流时地支"`
	GanRels    *[]BaziSeriesGanRelationItem `json:"gan_rels,omitempty" jsonschema:"天干关系"`
	ZhiRels    *[]BaziSeriesZhiRelationItem `json:"zhi_rels,omitempty" jsonschema:"地支关系"`
	ShiChen    BaziLiuShiChen               `json:"shi_chen" jsonschema:"流时时辰定位。"`
}

// ── bazi.liupan xiao_yun output ──

type BaziXiaoYunOutput struct {
	XiaoYun []BaziXiaoYunItem `json:"xiao_yun" jsonschema:"各学派小运列表；长度 = count。"`
}

type BaziXiaoYunItem struct {
	Method string                    `json:"method" jsonschema:"流派标识"`
	Basis  string                    `json:"basis" jsonschema:"起法依据"`
	Zhus   []BaziXiaoYunItemZhusItem `json:"zhus" jsonschema:"该学派小运柱列表。"`
}

type BaziXiaoYunItemZhusItem struct {
	Age     int         `json:"age" jsonschema:"虚岁。"`
	Gan     GanName     `json:"gan" jsonschema:"天干。"`
	Zhi     ZhiName     `json:"zhi" jsonschema:"地支。"`
	Name    string      `json:"name" jsonschema:"名称。"`
	ShiShen *TenGodName `json:"shi_shen,omitempty" jsonschema:"十神。"`
}

// ── bazi.liupan union output ──

// BaziLiupanOutput carries exactly one selected flow layer. The explicit
// output schema keeps the current single-root wire face while allowing one
// typed handler to serve the five scope forms.
type BaziLiupanOutput struct {
	LiuNian []bazi.LiuNian    `json:"liu_nian,omitempty" jsonschema:"流年数组，长度 = count。"`
	LiuYue  []bazi.LiuYue     `json:"liu_yue,omitempty" jsonschema:"流月数组，长度 = count。"`
	LiuRi   []bazi.LiuRi      `json:"liu_ri,omitempty" jsonschema:"流日数组，长度 = count。"`
	LiuShi  []bazi.LiuShi     `json:"liu_shi,omitempty" jsonschema:"流时数组，长度 = count。"`
	XiaoYun []bazi.XiaoYunSet `json:"xiao_yun,omitempty" jsonschema:"各学派小运列表；长度 = count。"`
}

// ── named natal series roots ──

type BaziLiuNianSeries struct {
	LiuNian []BaziLiuNianOutput `json:"liu_nian" jsonschema:"流年数组，长度 = count。"`
}
