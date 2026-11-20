package bazi

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// --- per-pillar data ---

type HiddenStemProjection struct {
	Main  ganzhi.Gan  `json:"main"`
	Mid   *ganzhi.Gan `json:"mid"`
	Minor *ganzhi.Gan `json:"minor"`
}

// Pillar is the single published representation of one Bazi pillar. It carries
// the pillar itself and all pillar-owned projections in one domain value.
type Pillar struct {
	ganzhi.Zhu
	DayMasterTrend string               `json:"ri_zhu_trend"`
	Xun            string               `json:"xun"`
	XunKong        string               `json:"xun_kong"`
	NaYin          string               `json:"na_yin"`
	CangGan        HiddenStemProjection `json:"cang_gan"`
	ShiShens       []TenGodProjection   `json:"shi_shens"`
	SelfSitting    string               `json:"self_sitting"`
	ShenSha        []ShenShaProjection  `json:"shen_sha"`
	IsVoid         bool                 `json:"is_void"`
	IsSelfHe       bool                 `json:"is_self_he"`
	IsKuiGang      bool                 `json:"is_kui_gang"`
	SelfHeName     string               `json:"self_he_name"`
}

type TenGodProjection struct {
	ShiShen ganzhi.ShiShen `json:"shi_shen"`
	Name    string         `json:"name"`
	Source  TenGodSource   `json:"source"`
	Gan     ganzhi.Gan     `json:"gan"`
}

// Ten god source constants.
const (
	sourceGan    TenGodSource = "gan"
	sourceMainQi TenGodSource = "main_qi"
	sourceMidQi  TenGodSource = "mid_qi"
	sourceMinQi  TenGodSource = "minor_qi"
)

// --- chart ---

// Chart holds a complete bazi chart: four pillars, dayun, and gender.
// Chart is the domain-owned, complete Bazi chart. Its fields are the canonical
// bazi.chart projection; derived methods receive the narrower Core.
type Chart struct {
	Nian   Pillar        `json:"nian"`
	Yue    Pillar        `json:"yue"`
	Ri     Pillar        `json:"ri"`
	Shi    Pillar        `json:"shi"`
	DaYun  *DaYun        `json:"da_yun"`
	Gender ganzhi.Gender `json:"-"`

	// 出生公历年份。derived methods 用它定位大运；bazi.chart 同时输出该字段。
	BirthYear int `json:"birth_year"`

	// 用神候选与结构证据独立输出：FuYi 是扶抑候选，
	// TiaoHou 是调候候选，GeJu 是月令结构证据；最终用神由判断层裁决。
	FuYi    FuYiResult    `json:"fu_yi"`
	TiaoHou TiaoHouResult `json:"tiao_hou"`
	GeJu    GeJuResult    `json:"ge_ju"`
	// 神煞双参照流派：年参照与日参照并取，不在解释层隐含二选一。
	ShenShaSchool ShenShaSchool `json:"shen_sha_school"`

	// 确定性命盘事实：透干十神得十干禄位。
	LuRoots []LuRoot `json:"lu_roots"`
	// 合会冲刑 canonical facts；旧并列数组是 detection evidence，待消费方迁移。
	GuanXi []NatalRelationFact `json:"guan_xi"`
	// 十神显隐、通根、得令与旺弱状态；Python 因子层只读。
	TenGodStates []TenGodState `json:"shi_shen_states"`
	// 五行季节旺弱、组合旺弱与生克方向；Python 因子层只读。
	ElementStates []ElementState `json:"wuxing_states"`
	// 官杀与财库是具体命理结论，不是通用事实袋。
	GuanSha GuanSha `json:"guan_sha"`
	CaiKu   CaiKu   `json:"cai_ku"`
	// 从 guan_xi/shen_sha 派生的确定性命理事实（符号空间内计算）。
	SpousePalaceState PalaceState     `json:"spouse_palace_state"`
	DayBranchType     BranchImageKind `json:"day_branch_type"`

	// 补充信息（原 bazi.chart_extra）
	SanYuan    SanYuan             `json:"san_yuan"`
	TaiXi      ganzhi.Zhu          `json:"tai_xi"`
	GongJia    []GongJia           `json:"gong_jia,omitempty"`
	NayinRel   []NayinRelEntry     `json:"na_yin_rel"`
	ChangSheng [12]ChangShengStage `json:"chang_sheng"`
	SanQiName  string              `json:"san_qi_name,omitempty"`
	// 生肖（年支对应动物，如酉→鸡）。
	Zodiac string `json:"sheng_xiao"`

	// 日柱所在旬与日柱旬空。is_void 语义固定为该柱地支值日柱旬空；
	// 各柱 xun/xun_kong 是该柱自身所在旬，两者是不同命理事实。
	DayXun     string `json:"ri_xun"`
	DayXunKong string `json:"ri_xun_kong"`

	// 子时换日规则说明（引擎口径说明；供上层核验换日规则，非计算输入）。
	ZiShiRule string `json:"zi_shi_rule"`
}

// ShenShaSchool records the union policy used by dual-reference shensha.
type ShenShaSchool struct {
	DualReference []string `json:"dual_reference"`
	Policy        string   `json:"policy"`
}

func (c Chart) ToBazi() ganzhi.Bazi {
	return ganzhi.Bazi{
		Nian: c.Nian.Zhu,
		Yue:  c.Yue.Zhu,
		Ri:   c.Ri.Zhu,
		Shi:  c.Shi.Zhu,
	}
}
func (c Chart) NaYinArray() [4]string {
	return [4]string{c.Nian.NaYin, c.Yue.NaYin, c.Ri.NaYin, c.Shi.NaYin}
}

var zhuLabels = [4]string{"nian", "yue", "ri", "shi"}
