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

type BaziBondOutput struct {
	DayMaster    BaziBondDayMaster    `json:"ri_zhu" jsonschema:"双方日主及日主生克关系。"`
	SpousePalace BaziBondSpousePalace `json:"pei_ou_gong" jsonschema:"双方配偶宫（日支）及其关系。"`
	ZhuCross     BaziBondZhuCross     `json:"zhu_cross" jsonschema:"双方四柱逐柱干支交叉。"`
	ShiShenCross BaziBondShiShenCross `json:"shi_shen_cross" jsonschema:"双方日主互视的十神柱位。"`
	WuxingCross  BaziBondWuxingCross  `json:"wuxing_cross" jsonschema:"双方实盘五行数量与扶抑对照。"`
	NaYinCross   BaziBondNaYinCross   `json:"na_yin_cross" jsonschema:"双方四柱纳音交叉。"`
	SpouseStar   BaziBondSpouseStar   `json:"pei_ou_xing" jsonschema:"双方配偶星候选与落柱。"`
	ShenshaCross BaziBondShenshaCross `json:"shen_sha_cross" jsonschema:"双方互见神煞事实。"`
}

type BaziBondDayMaster struct {
	A        BaziBondSideIdentity `json:"a" jsonschema:"甲方日主。"`
	B        BaziBondSideIdentity `json:"b" jsonschema:"乙方日主。"`
	Relation BaziBondRelation     `json:"relation" jsonschema:"日主关系。"`
}

type BaziBondSideIdentity struct {
	Gan    GanName     `json:"gan" jsonschema:"天干。"`
	Wuxing ElementName `json:"wuxing" jsonschema:"五行。"`
	Gender GenderName  `json:"gender" jsonschema:"性别。"`
}

type BaziBondRelation struct {
	Kind     string `json:"kind" jsonschema:"关系类型。"`
	Relation string `json:"relation" jsonschema:"关系描述。"`
}

type BaziBondSpousePalace struct {
	A        BaziBondZhiSide           `json:"a" jsonschema:"甲方日支。"`
	B        BaziBondZhiSide           `json:"b" jsonschema:"乙方日支。"`
	Relation BaziBondExplainedRelation `json:"relation" jsonschema:"日支关系。"`
}

type BaziBondZhiSide struct {
	Zhi ZhiName `json:"zhi" jsonschema:"日支。"`
}

type BaziBondExplainedRelation struct {
	Kind        string `json:"kind" jsonschema:"关系类型。"`
	Explanation string `json:"explanation" jsonschema:"关系解释。"`
}

type BaziBondZhuCross struct {
	Pairs []BaziBondZhuCrossPairsItem `json:"pairs" jsonschema:"四柱配对；按年、月、日、时对齐。"`
}

type BaziBondZhuCrossPairsItem struct {
	A           BaziBondZhuSide     `json:"a" jsonschema:"甲方柱。"`
	B           BaziBondZhuSide     `json:"b" jsonschema:"乙方柱。"`
	GanRelation BaziBondRelation    `json:"gan_relation" jsonschema:"天干关系。"`
	ZhiRelation BaziBondZhiRelation `json:"zhi_relation" jsonschema:"地支关系。"`
}

type BaziBondZhuSide struct {
	Zhu PillarName `json:"zhu" jsonschema:"柱位。"`
	Gan GanName    `json:"gan" jsonschema:"天干。"`
	Zhi ZhiName    `json:"zhi" jsonschema:"地支。"`
}

type BaziBondZhiRelation struct {
	Kind        string `json:"kind" jsonschema:"关系类型。"`
	Explanation string `json:"explanation" jsonschema:"关系解释。"`
}

type BaziBondShiShenCross struct {
	AToB BaziBondTenGodPillars `json:"a_to_b" jsonschema:"甲方日主视乙方四柱的十神。"`
	BToA BaziBondTenGodPillars `json:"b_to_a" jsonschema:"乙方日主视甲方四柱的十神。"`
}

type BaziBondTenGodPillars struct {
	NianGan TenGodName `json:"nian_gan" jsonschema:"年干。"`
	YueGan  TenGodName `json:"yue_gan" jsonschema:"月干。"`
	RiGan   TenGodName `json:"ri_gan" jsonschema:"日干。"`
	ShiGan  TenGodName `json:"shi_gan" jsonschema:"时干。"`
}

type BaziBondWuxingCross struct {
	A           BaziBondElementCounts          `json:"a" jsonschema:"甲方计数或候选。"`
	B           BaziBondElementCounts          `json:"b" jsonschema:"乙方计数或候选。"`
	Combined    BaziBondElementCounts          `json:"combined" jsonschema:"双方合并计数。"`
	FuYiInOther BaziBondWuxingCrossFuYiInOther `json:"fu_yi_in_other" jsonschema:"一方扶抑候选在另一方实盘中的命中数量。"`
}

type BaziBondElementCounts struct {
	Mu   int `json:"木" jsonschema:"木的数量。"`
	Huo  int `json:"火" jsonschema:"火的数量。"`
	Tu   int `json:"土" jsonschema:"土的数量。"`
	Jin  int `json:"金" jsonschema:"金的数量。"`
	Shui int `json:"水" jsonschema:"水的数量。"`
}

type BaziBondWuxingCrossFuYiInOther struct {
	A BaziBondFuYiInOtherSide `json:"a" jsonschema:"甲方计数或候选。"`
	B BaziBondFuYiInOtherSide `json:"b" jsonschema:"乙方计数或候选。"`
}

type BaziBondFuYiInOtherSide struct {
	ModelID     string              `json:"fu_yi_model_id" jsonschema:"扶抑候选模型 ID。"`
	Yong        CandidateWuxingName `json:"yong" jsonschema:"用神；缺省为 无候选。"`
	Xi          CandidateWuxingName `json:"xi" jsonschema:"喜神；缺省为 无候选。"`
	Ji          CandidateWuxingName `json:"ji" jsonschema:"忌神；缺省为 无候选。"`
	YongInOther int                 `json:"yong_in_other" jsonschema:"对方用神五行在本方实盘中的数量。"`
	XiInOther   int                 `json:"xi_in_other" jsonschema:"对方喜神五行在本方实盘中的数量。"`
	JiInOther   int                 `json:"ji_in_other" jsonschema:"对方忌神五行在本方实盘中的数量。"`
}

type BaziBondNaYinCross struct {
	Pairs        []BaziBondNaYinCrossPairsItem  `json:"pairs" jsonschema:"四柱配对；按年、月、日、时对齐。"`
	WuxingCounts BaziBondNaYinCrossWuxingCounts `json:"wuxing_counts" jsonschema:"双方纳音五行计数。"`
}

type BaziBondNaYinCrossPairsItem struct {
	A        BaziBondNaYinSide `json:"a" jsonschema:"甲方柱。"`
	B        BaziBondNaYinSide `json:"b" jsonschema:"乙方柱。"`
	Relation string            `json:"relation" jsonschema:"纳音关系。"`
}

type BaziBondNaYinSide struct {
	Zhu   PillarName `json:"zhu" jsonschema:"柱位。"`
	NaYin string     `json:"na_yin" jsonschema:"纳音。"`
}

type BaziBondNaYinCrossWuxingCounts struct {
	A map[string]int `json:"a" jsonschema:"甲方纳音五行计数。"`
	B map[string]int `json:"b" jsonschema:"乙方纳音五行计数。"`
}

type BaziBondSpouseStar struct {
	A BaziBondSpouseStarSide `json:"a" jsonschema:"甲方配偶星。"`
	B BaziBondSpouseStarSide `json:"b" jsonschema:"乙方配偶星。"`
}

type BaziBondSpouseStarSide struct {
	Gender    GenderName                  `json:"gender" jsonschema:"性别。"`
	Primary   BaziBondSpouseStarCandidate `json:"primary" jsonschema:"主调候。"`
	Secondary BaziBondSpouseStarCandidate `json:"secondary" jsonschema:"辅调候。"`
}

type BaziBondSpouseStarOccurrence struct {
	Pillar PillarName `json:"zhu" jsonschema:"柱位。"`
	Branch ZhiName    `json:"zhi" jsonschema:"地支。"`
	Stem   GanName    `json:"gan" jsonschema:"天干。"`
	Source string     `json:"source" jsonschema:"来源。"`
	IsVoid bool       `json:"is_void" jsonschema:"是否旬空。"`
}

type BaziBondSpouseStarCandidate struct {
	TenGod      TenGodName                     `json:"shi_shen" jsonschema:"配偶星十神。"`
	Occurrences []BaziBondSpouseStarOccurrence `json:"occurrences" jsonschema:"出现位置。"`
}

type BaziBondShenshaCross struct {
	TianYi   BaziBondMutualFact `json:"tian_yi" jsonschema:"天乙贵人互见。"`
	Lu       BaziBondMutualFact `json:"lu" jsonschema:"禄神互见。"`
	TaoHua   BaziBondMutualFact `json:"tao_hua" jsonschema:"桃花互见。"`
	YiMa     BaziBondMutualFact `json:"yi_ma" jsonschema:"驿马互见。"`
	KongWang BaziBondMutualFact `json:"kong_wang" jsonschema:"空亡。"`
	KuiGang  BaziBondMutualFact `json:"kui_gang" jsonschema:"魁罡互见。"`
	RiDe     BaziBondMutualFact `json:"ri_de" jsonschema:"日德互见。"`
	RiGui    BaziBondMutualFact `json:"ri_gui" jsonschema:"日贵互见。"`
}

type BaziBondMutualFact struct {
	AToB bool `json:"a_to_b" jsonschema:"甲方神煞是否见于乙方。"`
	BToA bool `json:"b_to_a" jsonschema:"乙方神煞是否见于甲方。"`
}

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
