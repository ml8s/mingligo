package qimen

// This file defines the domain-owned, domain-aggregated output view of a Qimen
// chart. It groups related facts into术数子域 objects and is the single output
// model that every transport adapter transmits unchanged. jsonschema tags are
// inert metadata read only when a consumer reflects the type for its schema.

// KeYingView groups 克应 facts (十干/八门/九星克应 + 星宫五行).
type KeYingView struct {
	GanInteraction  []GanInteraction  `json:"gan_interaction,omitempty" jsonschema:"十干克应。"`
	MenInteraction  []MenInteraction  `json:"men_interaction,omitempty" jsonschema:"八门克应。"`
	XingInteraction []XingInteraction `json:"xing_interaction,omitempty" jsonschema:"九星克应。"`
	XingGongWuXing  []XingGongWuXing  `json:"xing_gong_wuxing,omitempty" jsonschema:"星宫五行关系。"`
}

// JiuGongGuYouView groups 九宫固有象义 (不随本盘变化): 宫位属性/宫际关系/门星吉凶.
type JiuGongGuYouView struct {
	GongProfiles  [9]GongProfile `json:"gong_profiles" jsonschema:"九宫固有象义（方位/人物/范围/内外盘）。"`
	GongRelations []GongRelation `json:"gong_relations" jsonschema:"九宫五行关系矩阵。"`
	MenJiXiong    []MenJiXiong   `json:"men_ji_xiong" jsonschema:"八门固有吉凶。"`
	XingJiXiong   []XingJiXiong  `json:"xing_ji_xiong" jsonschema:"九星固有吉凶。"`
}

// JiuGongBenPanView groups 九宫本盘状态 (换局即变): 宫旺衰/门迫/门制.
type JiuGongBenPanView struct {
	GongWangShuai []GongWangShuai `json:"gong_wang_shuai" jsonschema:"九宫旺相休囚死（按当前月令）。"`
	MenPo         []GongIndex     `json:"men_po,omitempty" jsonschema:"门迫宫位。"`
	MenZhi        []GongIndex     `json:"men_zhi,omitempty" jsonschema:"门制宫位。"`
}

// GanLuoGongView groups 干落宫与日时宫关系 facts.
type GanLuoGongView struct {
	DayStemPalace  StemPalace    `json:"ri_gan_gong" jsonschema:"日干落宫。"`
	HourStemPalace StemPalace    `json:"shi_gan_gong" jsonschema:"时干落宫。"`
	RiShiRelation  RiShiRelation `json:"ri_shi_relation" jsonschema:"日时宫关系。"`
}

// FuShiView groups 值符(星)/值使(门)落宫 facts.
type FuShiView struct {
	DutyStarPalace GongIndex `json:"zhi_fu_xing_gong" jsonschema:"值符星落宫。"`
	DutyDoorPalace GongIndex `json:"zhi_shi_men_gong" jsonschema:"值使门落宫。"`
}

// ShenShaView groups 神煞 facts (空亡/马星影响扫描).
type ShenShaView struct {
	KongWangAffected []AffectedSymbol `json:"kong_wang_affected" jsonschema:"空亡影响符号。"`
	MaXingAffected   []AffectedSymbol `json:"ma_xing_affected" jsonschema:"马星影响符号。"`
}

// GeJuView groups 格局 facts.
type GeJuView struct {
	Patterns []Pattern `json:"patterns" jsonschema:"格局。"`
}

// ChartView is the domain-aggregated output of a Qimen chart, grouped by命理域.
type ChartView struct {
	PaiPanFangFa  ChartMethod        `json:"pai_pan_fang_fa" jsonschema:"排盘方法元数据。"`
	Pan           Board              `json:"pan" jsonschema:"盘面。"`
	KeYing        KeYingView         `json:"ke_ying" jsonschema:"克应事实。"`
	JiuGongGuYou  JiuGongGuYouView   `json:"jiu_gong_gu_you" jsonschema:"九宫固有象义。"`
	JiuGongBenPan JiuGongBenPanView  `json:"jiu_gong_ben_pan" jsonschema:"九宫本盘状态。"`
	GanLuoGong    GanLuoGongView     `json:"gan_luo_gong" jsonschema:"干落宫事实。"`
	FuShi         FuShiView          `json:"fu_shi" jsonschema:"值符值使落宫。"`
	ShenSha       ShenShaView        `json:"shen_sha" jsonschema:"神煞。"`
	GeJu          GeJuView           `json:"ge_ju" jsonschema:"格局。"`
	YingQi        YingQi             `json:"ying_qi" jsonschema:"应期候选。"`
	YongShen      *YongShenResult    `json:"yong_shen,omitempty" jsonschema:"用神落宫事实。"`
	ZhuanZhan     SpecializedContext `json:"zhuan_zhan" jsonschema:"专占领域原子事实。"`
}

// View produces the domain-aggregated output from the computed chart. Filling
// order is free; each fact is placed into its owning术数子域.
func (c Chart) View() ChartView {
	return ChartView{
		PaiPanFangFa: c.Method,
		Pan:          c.Board,
		KeYing: KeYingView{
			GanInteraction: c.GanInteractions, MenInteraction: c.MenInteractions,
			XingInteraction: c.XingInteractions, XingGongWuXing: c.XingGongWuXing,
		},
		JiuGongGuYou: JiuGongGuYouView{
			GongProfiles: c.GongProfiles, GongRelations: c.GongRelations,
			MenJiXiong: c.MenJiXiong, XingJiXiong: c.XingJiXiong,
		},
		JiuGongBenPan: JiuGongBenPanView{
			GongWangShuai: c.GongWangShuai, MenPo: c.MenPo, MenZhi: c.MenZhi,
		},
		GanLuoGong: GanLuoGongView{
			DayStemPalace: c.DayStemPalace, HourStemPalace: c.HourStemPalace,
			RiShiRelation: c.RiShiRelation,
		},
		FuShi:   FuShiView{DutyStarPalace: c.DutyStarPalace, DutyDoorPalace: c.DutyDoorPalace},
		ShenSha: ShenShaView{KongWangAffected: c.KongWangAffected, MaXingAffected: c.MaXingAffected},
		GeJu:    GeJuView{Patterns: c.Patterns},
		YingQi:  c.YingQi, YongShen: c.YongShen, ZhuanZhan: c.Specialized,
	}
}
