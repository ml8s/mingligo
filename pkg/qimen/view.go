package qimen

// This file defines the domain-owned, domain-aggregated output view of a Qimen
// chart. It groups related facts into术数子域 objects and is the single output
// model that every transport adapter transmits unchanged. jsonschema tags are
// inert metadata read only when a consumer reflects the type for its schema.

// PaiPanFangFaView is the 排盘方法域. It mirrors the method facts with命理
// pinyin keys and never leaks implementation field names into the wire.
type PaiPanFangFaView struct {
	FanChou       string        `json:"fan_chou" jsonschema:"排盘范围。"`
	LiuPai        string        `json:"liu_pai" jsonschema:"流派。"`
	DingJuFangFa  string        `json:"ding_ju_fang_fa,omitempty" jsonschema:"定局方法。"`
	JieQiGuiZe    string        `json:"jie_qi_gui_ze,omitempty" jsonschema:"寄宫/分季规则。"`
	JiChuDingJu   string        `json:"ji_chu_ding_ju,omitempty" jsonschema:"基础定局。"`
	FangFaLaiYuan string        `json:"fang_fa_lai_yuan" jsonschema:"方法来源。"`
	RiJie         string        `json:"ri_jie" jsonschema:"日界。"`
	ShenPanMoShi  string        `json:"shen_pan_mo_shi" jsonschema:"神盘模式。"`
	JiHeMoShi     string        `json:"ji_he_mo_shi" jsonschema:"几何排布模式。"`
	XingPanMoShi  string        `json:"xing_pan_mo_shi" jsonschema:"星盘模式。"`
	MenPanMoShi   string        `json:"men_pan_mo_shi" jsonschema:"门盘模式。"`
	ShenFeiMoShi  string        `json:"shen_fei_mo_shi" jsonschema:"神飞模式。"`
	TianQin       *TianQinRule  `json:"tian_qin,omitempty" jsonschema:"天禽规则。"`
	JieQi         string        `json:"jie_qi" jsonschema:"节气。"`
	YongJuJieQi   string        `json:"yong_ju_jie_qi" jsonschema:"用局节气。"`
	Yuan          string        `json:"yuan" jsonschema:"三元。"`
	ZhiRun        string        `json:"zhi_run,omitempty" jsonschema:"置闰状态。"`
	LingZhu       PillarInfo    `json:"ling_zhu" jsonschema:"值符柱。"`
	XunShouZhu    PillarInfo    `json:"xun_shou_zhu" jsonschema:"旬首柱。"`
	XunShouGong   GongIndex     `json:"xun_shou_gong" jsonschema:"旬首落宫。"`
	RiFuTou       PillarInfo    `json:"ri_fu_tou" jsonschema:"日符头。"`
	Jie           *ChartQuarter `json:"jie,omitempty" jsonschema:"分季。"`
}

func paiPanFangFaView(m ChartMethod) PaiPanFangFaView {
	return PaiPanFangFaView{
		FanChou: m.Scope, LiuPai: m.School, DingJuFangFa: m.DingjuMethod,
		JieQiGuiZe: m.QuarterRule, JiChuDingJu: m.BaseDingjuMethod,
		FangFaLaiYuan: displayFrom(methodSourceDisplay, m.MethodSource),
		RiJie:         displayFrom(dayBoundaryDisplay, m.DayBoundary),
		ShenPanMoShi:  displayFrom(spiritModeDisplay, m.SpiritMode),
		JiHeMoShi:     displayFrom(geometryDisplay, m.Geometry),
		XingPanMoShi:  displayFrom(starModeDisplay, m.StarMode),
		MenPanMoShi:   displayFrom(doorModeDisplay, m.DoorMode),
		ShenFeiMoShi:  displayFrom(spiritFlightDisplay, m.SpiritFlight), TianQin: m.TianQin,
		JieQi: m.JieQi, YongJuJieQi: m.YongJuJieQi, Yuan: m.Yuan,
		ZhiRun: m.ZhiRunState, LingZhu: m.LeadPillar,
		XunShouZhu: m.LeadXunShou, XunShouGong: m.LeadXunShouGong,
		RiFuTou: m.DayFuTou, Jie: m.Quarter,
	}
}

// KeYingView groups 克应 facts (十干/八门/九星克应 + 星宫五行).
type KeYingView struct {
	GanKeYing      []GanInteraction     `json:"gan_ke_ying,omitempty" jsonschema:"十干克应。"`
	MenKeYing      []MenInteraction     `json:"men_ke_ying,omitempty" jsonschema:"八门克应。"`
	XingKeYing     []XingInteraction    `json:"xing_ke_ying,omitempty" jsonschema:"九星克应。"`
	XingGongGuanXi []XingGongWuXingView `json:"xing_gong_guan_xi,omitempty" jsonschema:"星宫五行关系。"`
}

// JiuGongGuYouView groups 九宫固有象义 (不随本盘变化): 宫位属性/宫际关系/门星吉凶.
type JiuGongGuYouView struct {
	GongProfiles  [9]GongProfileView `json:"gong_xiang_yi" jsonschema:"九宫固有象义（方位/人物/范围/内外盘）。"`
	GongRelations []GongRelationView `json:"gong_guan_xi" jsonschema:"九宫五行关系矩阵。"`
	MenJiXiong    []MenJiXiong       `json:"men_ji_xiong" jsonschema:"八门固有吉凶。"`
	XingJiXiong   []XingJiXiong      `json:"xing_ji_xiong" jsonschema:"九星固有吉凶。"`
}

// JiuGongBenPanView groups 九宫本盘状态 (换局即变): 宫旺衰/门迫/门制.
type JiuGongBenPanView struct {
	GongWangShuai []GongWangShuai `json:"gong_wang_shuai" jsonschema:"九宫旺相休囚死（按当前月令）。"`
	MenPo         []GongIndex     `json:"men_po,omitempty" jsonschema:"门迫宫位。"`
	MenZhi        []GongIndex     `json:"men_zhi,omitempty" jsonschema:"门制宫位。"`
}

// GanLuoGongView groups 干落宫与日时宫关系 facts.
type GanLuoGongView struct {
	DayStemPalace  StemPalaceView    `json:"ri_gan_gong" jsonschema:"日干落宫。"`
	HourStemPalace StemPalaceView    `json:"shi_gan_gong" jsonschema:"时干落宫。"`
	RiShiRelation  RiShiRelationView `json:"ri_shi_relation" jsonschema:"日时宫关系。"`
}

// FuShiView groups 值符(星)/值使(门)落宫 facts.
type FuShiView struct {
	DutyStarPalace GongIndex `json:"zhi_fu_xing_gong" jsonschema:"值符星落宫。"`
	DutyDoorPalace GongIndex `json:"zhi_shi_men_gong" jsonschema:"值使门落宫。"`
}

// ShenShaView groups 神煞 facts (空亡/马星影响扫描).
type ShenShaView struct {
	KongWangAffected []AffectedSymbolView `json:"kong_wang_ying_xiang" jsonschema:"空亡影响符号。"`
	MaXingAffected   []AffectedSymbolView `json:"ma_xing_ying_xiang" jsonschema:"马星影响符号。"`
}

// GeJuView groups 格局 facts.
type GeJuView struct {
	Patterns []PatternView `json:"ge_ju_ming_xi" jsonschema:"格局。"`
}

// ChartView is the domain-aggregated output of a Qimen chart, grouped by命理域.
type ChartView struct {
	PaiPanFangFa  PaiPanFangFaView       `json:"pai_pan_fang_fa" jsonschema:"排盘方法元数据。"`
	Pan           Board                  `json:"jiu_gong_pan" jsonschema:"盘面。"`
	KeYing        KeYingView             `json:"ke_ying" jsonschema:"克应事实。"`
	JiuGongGuYou  JiuGongGuYouView       `json:"jiu_gong_gu_you" jsonschema:"九宫固有象义。"`
	JiuGongBenPan JiuGongBenPanView      `json:"jiu_gong_ben_pan" jsonschema:"九宫本盘状态。"`
	GanLuoGong    GanLuoGongView         `json:"gan_luo_gong" jsonschema:"干落宫事实。"`
	FuShi         FuShiView              `json:"fu_shi" jsonschema:"值符值使落宫。"`
	ShenSha       ShenShaView            `json:"shen_sha" jsonschema:"神煞。"`
	GeJu          GeJuView               `json:"ge_ju" jsonschema:"格局。"`
	YingQi        YingQiView             `json:"ying_qi" jsonschema:"应期候选。"`
	YongShen      *YongShenResultView    `json:"yong_shen,omitempty" jsonschema:"用神落宫事实。"`
	ZhuanZhan     SpecializedContextView `json:"zhuan_zhan" jsonschema:"专占领域原子事实。"`
}

// View produces the domain-aggregated output from the computed chart. Filling
// order is free; each fact is placed into its owning术数子域.
func (c Chart) View() ChartView {
	return ChartView{
		PaiPanFangFa: paiPanFangFaView(c.Method),
		Pan:          c.Board,
		KeYing: KeYingView{
			GanKeYing: c.GanInteractions, MenKeYing: c.MenInteractions,
			XingKeYing: c.XingInteractions, XingGongGuanXi: xingGongWuXingViews(c.XingGongWuXing),
		},
		JiuGongGuYou: JiuGongGuYouView{
			GongProfiles: gongProfileViews(c.GongProfiles), GongRelations: gongRelationViews(c.GongRelations),
			MenJiXiong: c.MenJiXiong, XingJiXiong: c.XingJiXiong,
		},
		JiuGongBenPan: JiuGongBenPanView{
			GongWangShuai: c.GongWangShuai, MenPo: c.MenPo, MenZhi: c.MenZhi,
		},
		GanLuoGong: GanLuoGongView{
			DayStemPalace: stemPalaceView(c.DayStemPalace), HourStemPalace: stemPalaceView(c.HourStemPalace),
			RiShiRelation: riShiRelationView(c.RiShiRelation),
		},
		FuShi:   FuShiView{DutyStarPalace: c.DutyStarPalace, DutyDoorPalace: c.DutyDoorPalace},
		ShenSha: ShenShaView{KongWangAffected: affectedSymbolViews(c.KongWangAffected), MaXingAffected: affectedSymbolViews(c.MaXingAffected)},
		GeJu:    GeJuView{Patterns: patternViews(c.Patterns)},
		YingQi:  yingQiView(c.YingQi), YongShen: yongShenResultView(c.YongShen), ZhuanZhan: specializedContextView(c.Specialized),
	}
}
