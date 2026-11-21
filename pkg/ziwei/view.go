package ziwei

// This file defines the domain-owned, domain-aggregated output view of a Ziwei
// chart. It groups related facts into术数子域 objects and is the single output
// model that every transport adapter transmits unchanged. jsonschema tags are
// inert metadata read only when a consumer reflects the type for its schema.

// LiFaKouJingView groups历法 facts.
type LiFaKouJingView struct {
	NianGan         Gan  `json:"nian_gan" jsonschema:"年干。"`
	NianZhi         Zhi  `json:"nian_zhi,omitempty" jsonschema:"年支。"`
	ShiZhi          Zhi  `json:"shi_zhi" jsonschema:"时支。"`
	BirthYear       int  `json:"birth_year" jsonschema:"出生公历年。"`
	LunarMonth      int  `json:"lunar_month,omitempty" jsonschema:"排盘农历月。"`
	LunarDay        int  `json:"lunar_day,omitempty" jsonschema:"排盘农历日。"`
	BirthLunarMonth int  `json:"birth_lunar_month,omitempty" jsonschema:"出生原始农历月。"`
	BirthIsLeap     bool `json:"birth_is_leap,omitempty" jsonschema:"出生农历月是否闰月。"`
}

// PaiPanKouJingView groups排盘口径 facts.
type PaiPanKouJingView struct {
	KouJing   ChartSchool `json:"kou_jing" jsonschema:"排盘口径。"`
	JuShu     JuShu       `json:"ju_shu" jsonschema:"局数机器码。"`
	JuShuName string      `json:"ju_shu_name" jsonschema:"局数名。"`
}

// KongGongSanFangView groups宫位级引用事实 (星曜仍由 GongWei 源宫唯一承载).
type KongGongSanFangView struct {
	KongGong       []KongGongInfo `json:"kong_gong,omitempty" jsonschema:"空宫借星。"`
	SanFangSiZheng SanFangSiZheng `json:"san_fang_si_zheng" jsonschema:"命宫三方四正宫位关系。"`
}

// GeJuView groups格局候选 (不裁决成败).
type GeJuView struct {
	Candidates []Pattern `json:"candidates,omitempty" jsonschema:"格局候选。"`
}

// DaXianView groups时间轴序列 (与八字顶层 DaYun 对齐).
type DaXianView struct {
	DaXian []DaXianStep `json:"da_xian" jsonschema:"十年大限序列。"`
}

// MingZhuShenZhuView groups命主身主定星 (命宫支/年支的固定映射).
type MingZhuShenZhuView struct {
	MingZhu string `json:"ming_zhu,omitempty" jsonschema:"命主。"`
	ShenZhu string `json:"shen_zhu,omitempty" jsonschema:"身主。"`
}

// ChartView is the domain-aggregated output of a Ziwei chart.
type ChartView struct {
	LiFaKouJing     LiFaKouJingView     `json:"li_fa_kou_jing" jsonschema:"历法事实。"`
	PaiPanKouJing   PaiPanKouJingView   `json:"pai_pan_kou_jing" jsonschema:"排盘口径与局数。"`
	GongWei         [12]Gong            `json:"shi_er_gong" jsonschema:"十二宫；星曜由各宫唯一承载。"`
	KongGongSanFang KongGongSanFangView `json:"kong_gong_san_fang" jsonschema:"宫位引用事实。"`
	GeJu            GeJuView            `json:"ge_ju" jsonschema:"格局候选。"`
	DaXian          DaXianView          `json:"da_xian" jsonschema:"时间轴序列。"`
	MingZhuShenZhu  MingZhuShenZhuView  `json:"ming_zhu_shen_zhu" jsonschema:"命主身主定星。"`
}

// View produces the domain-aggregated output from the computed chart. Filling
// order is free; each fact is placed into its owning术数子域.
func (c Chart) View() ChartView {
	return ChartView{
		LiFaKouJing: LiFaKouJingView{
			NianGan: c.NianGan, NianZhi: c.NianZhi, ShiZhi: c.ShiZhi,
			BirthYear: c.BirthYear, LunarMonth: c.LunarMonth, LunarDay: c.LunarDay,
			BirthLunarMonth: c.BirthLunarMonth, BirthIsLeap: c.BirthIsLeap,
		},
		PaiPanKouJing: PaiPanKouJingView{KouJing: c.KouJing, JuShu: c.JuShu, JuShuName: c.JuShuName},
		GongWei:       c.GongWei,
		KongGongSanFang: KongGongSanFangView{
			KongGong: c.KongGong, SanFangSiZheng: c.SanFangSiZheng,
		},
		GeJu:           GeJuView{Candidates: c.Patterns},
		DaXian:         DaXianView{DaXian: c.DaXian},
		MingZhuShenZhu: MingZhuShenZhuView{MingZhu: c.MingZhu, ShenZhu: c.ShenZhu},
	}
}
