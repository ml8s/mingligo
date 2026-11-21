package ziwei

// This file defines the domain-owned, domain-aggregated output view of a Ziwei
// chart. It groups related facts into术数子域 objects and is the single output
// model that every transport adapter transmits unchanged. jsonschema tags are
// inert metadata read only when a consumer reflects the type for its schema.

// CalendarView groups历法 facts.
type CalendarView struct {
	NianGan         Gan  `json:"nian_gan" jsonschema:"年干。"`
	NianZhi         Zhi  `json:"nian_zhi,omitempty" jsonschema:"年支。"`
	ShiZhi          Zhi  `json:"shi_zhi" jsonschema:"时支。"`
	BirthYear       int  `json:"birth_year" jsonschema:"出生公历年。"`
	LunarMonth      int  `json:"lunar_month,omitempty" jsonschema:"排盘农历月。"`
	LunarDay        int  `json:"lunar_day,omitempty" jsonschema:"排盘农历日。"`
	BirthLunarMonth int  `json:"birth_lunar_month,omitempty" jsonschema:"出生原始农历月。"`
	BirthIsLeap     bool `json:"birth_is_leap,omitempty" jsonschema:"出生农历月是否闰月。"`
}

// SchoolView groups排盘口径 facts.
type SchoolView struct {
	KouJing   ChartSchool `json:"kou_jing" jsonschema:"排盘口径。"`
	JuShu     JuShu       `json:"ju_shu" jsonschema:"局数机器码。"`
	JuShuName string      `json:"ju_shu_name" jsonschema:"局数名。"`
}

// MetaView groups元辰/格局 facts (宫位引用，星曜由源宫承载).
type MetaView struct {
	KongGong       []KongGongInfo `json:"kong_gong,omitempty" jsonschema:"空宫借星。"`
	Patterns       []Pattern      `json:"ge_ju,omitempty" jsonschema:"格局候选。"`
	SanFangSiZheng SanFangSiZheng `json:"san_fang_si_zheng" jsonschema:"命宫三方四正宫位关系。"`
	DaXian         []DaXianStep   `json:"da_xian" jsonschema:"十年大限序列。"`
	MingZhu        string         `json:"ming_zhu,omitempty" jsonschema:"命主。"`
	ShenZhu        string         `json:"shen_zhu,omitempty" jsonschema:"身主。"`
}

// ChartView is the domain-aggregated output of a Ziwei chart.
type ChartView struct {
	Calendar CalendarView `json:"calendar" jsonschema:"历法事实。"`
	School   SchoolView   `json:"school" jsonschema:"排盘口径与局数。"`
	GongWei  [12]Gong     `json:"shi_er_gong" jsonschema:"十二宫；星曜由各宫唯一承载。"`
	Meta     MetaView     `json:"meta" jsonschema:"格局与宫位关系事实。"`
}

// View produces the domain-aggregated output from the computed chart. Filling
// order is free; each fact is placed into its owning术数子域.
func (c Chart) View() ChartView {
	return ChartView{
		Calendar: CalendarView{
			NianGan: c.NianGan, NianZhi: c.NianZhi, ShiZhi: c.ShiZhi,
			BirthYear: c.BirthYear, LunarMonth: c.LunarMonth, LunarDay: c.LunarDay,
			BirthLunarMonth: c.BirthLunarMonth, BirthIsLeap: c.BirthIsLeap,
		},
		School:  SchoolView{KouJing: c.KouJing, JuShu: c.JuShu, JuShuName: c.JuShuName},
		GongWei: c.GongWei,
		Meta: MetaView{
			KongGong: c.KongGong, Patterns: c.Patterns,
			SanFangSiZheng: c.SanFangSiZheng, DaXian: c.DaXian,
			MingZhu: c.MingZhu, ShenZhu: c.ShenZhu,
		},
	}
}
