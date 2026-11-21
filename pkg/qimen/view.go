package qimen

// This file defines the domain-owned, domain-aggregated output view of a Qimen
// chart. It groups related facts into术数子域 objects and is the single output
// model that every transport adapter transmits unchanged. jsonschema tags are
// inert metadata read only when a consumer reflects the type for its schema.

// InteractionsView groups 克应 facts (干/门/星 interaction + 星宫关系).
type InteractionsView struct {
	GanInteraction  []GanInteraction  `json:"gan_interaction,omitempty" jsonschema:"十干克应。"`
	MenInteraction  []MenInteraction  `json:"men_interaction,omitempty" jsonschema:"八门克应。"`
	XingInteraction []XingInteraction `json:"xing_interaction,omitempty" jsonschema:"九星克应。"`
	XingGongWuXing  []XingGongWuXing  `json:"xing_gong_wuxing,omitempty" jsonschema:"星宫五行关系。"`
}

// PalaceFactsView groups per-palace属性与静态关系事实.
type PalaceFactsView struct {
	GongProfiles  [9]GongProfile  `json:"gong_profiles" jsonschema:"九宫固有事实（方位/人物/范围/内外盘）。"`
	GongWangShuai []GongWangShuai `json:"gong_wang_shuai" jsonschema:"九宫旺相休囚死。"`
	GongRelations []GongRelation  `json:"gong_relations" jsonschema:"九宫五行关系矩阵。"`
	MenJiXiong    []MenJiXiong    `json:"men_ji_xiong" jsonschema:"八门固有吉凶。"`
	XingJiXiong   []XingJiXiong   `json:"xing_ji_xiong" jsonschema:"九星固有吉凶。"`
	MenPo         []GongIndex     `json:"men_po,omitempty" jsonschema:"门迫宫位。"`
	MenZhi        []GongIndex     `json:"men_zhi,omitempty" jsonschema:"门制宫位。"`
}

// StemsView groups 干落宫与日时关系 facts.
type StemsView struct {
	DayStemPalace  StemPalace    `json:"ri_gan_gong" jsonschema:"日干落宫。"`
	HourStemPalace StemPalace    `json:"shi_gan_gong" jsonschema:"时干落宫。"`
	RiShiRelation  RiShiRelation `json:"ri_shi_relation" jsonschema:"日时宫关系。"`
}

// FuShiView groups 值符(星)/值使(门)落宫 facts — the pivot of the 神盘,
// distinct from stem palaces.
type FuShiView struct {
	DutyStarPalace GongIndex `json:"zhi_fu_xing_gong" jsonschema:"值符星落宫。"`
	DutyDoorPalace GongIndex `json:"zhi_shi_men_gong" jsonschema:"值使门落宫。"`
}

// MarkerView groups 神煞/格局 facts.
type MarkerView struct {
	KongWangAffected []AffectedSymbol `json:"kong_wang_affected" jsonschema:"空亡影响符号（仅日干/时干/主柱三焦点扫描；用神落宫状态见 yong_shen.symbols[].kong_wang_zhi）。"`
	MaXingAffected   []AffectedSymbol `json:"ma_xing_affected" jsonschema:"马星影响符号（仅日干/时干/主柱三焦点扫描）。"`
	Patterns         []Pattern        `json:"ge_ju" jsonschema:"格局。"`
}

// ChartView is the domain-aggregated output of a Qimen standard chart.
type ChartView struct {
	Method       ChartMethod        `json:"method" jsonschema:"排盘方法元数据。"`
	Pan          Board              `json:"pan" jsonschema:"盘面。"`
	Interactions InteractionsView   `json:"interactions" jsonschema:"克应事实。"`
	PalaceFacts  PalaceFactsView    `json:"palace_facts" jsonschema:"宫位属性与静态关系。"`
	Stems        StemsView          `json:"stems" jsonschema:"干落宫事实。"`
	FuShi        FuShiView          `json:"fu_shi" jsonschema:"值符值使落宫。"`
	Markers      MarkerView         `json:"markers" jsonschema:"神煞格局事实。"`
	YingQi       YingQi             `json:"ying_qi" jsonschema:"应期候选。"`
	YongShen     *YongShenResult    `json:"yong_shen,omitempty" jsonschema:"用神落宫事实。"`
	Specialized  SpecializedContext `json:"special_facts" jsonschema:"专占领域原子事实。"`
}

// View produces the domain-aggregated output from the computed chart. Filling
// order is free; each fact is placed into its owning术数子域.
func (c Chart) View() ChartView {
	return ChartView{
		Method: c.Method,
		Pan:    c.Board,
		Interactions: InteractionsView{
			GanInteraction: c.GanInteractions, MenInteraction: c.MenInteractions,
			XingInteraction: c.XingInteractions, XingGongWuXing: c.XingGongWuXing,
		},
		PalaceFacts: PalaceFactsView{
			GongProfiles: c.GongProfiles, GongWangShuai: c.GongWangShuai,
			GongRelations: c.GongRelations, MenJiXiong: c.MenJiXiong,
			XingJiXiong: c.XingJiXiong, MenPo: c.MenPo, MenZhi: c.MenZhi,
		},
		Stems: StemsView{
			DayStemPalace: c.DayStemPalace, HourStemPalace: c.HourStemPalace,
			RiShiRelation: c.RiShiRelation,
		},
		FuShi: FuShiView{DutyStarPalace: c.DutyStarPalace, DutyDoorPalace: c.DutyDoorPalace},
		Markers: MarkerView{
			KongWangAffected: c.KongWangAffected, MaXingAffected: c.MaXingAffected,
			Patterns: c.Patterns,
		},
		YingQi: c.YingQi, YongShen: c.YongShen, Specialized: c.Specialized,
	}
}
