package bazi

import "github.com/ml8s/mingligo/pkg/ganzhi"

// This file defines the domain-owned, domain-aggregated output view of a Bazi
// chart. Every top-level key is a命理域; computation-only trace fields stay in
// 方法口径 and never appear among chart facts. jsonschema tags are inert
// metadata read only when a consumer reflects the type for its schema.

// PillarsView groups the four pillars (干支/藏干/十神投影/神煞/纳音/旬空/长生落点).
type PillarsView struct {
	Nian PillarView `json:"nian" jsonschema:"年柱。"`
	Yue  PillarView `json:"yue" jsonschema:"月柱。"`
	Ri   PillarView `json:"ri" jsonschema:"日柱。"`
	Shi  PillarView `json:"shi" jsonschema:"时柱。"`
}

// FuYiView is the扶抑 candidate. 旺衰 is published once by WangShuaiView.
type FuYiView struct {
	WuxingCount map[string]int  `json:"wuxing_count" jsonschema:"五行个数。"`
	Yong        CandidateWuxing `json:"yong" jsonschema:"扶抑用神候选。"`
	Xi          CandidateWuxing `json:"xi" jsonschema:"扶抑喜神候选。"`
	Ji          CandidateWuxing `json:"ji" jsonschema:"扶抑忌神候选。"`
	Strength    string          `json:"qiang_ruo" jsonschema:"身强身弱。"`
	Pattern     string          `json:"ge_ju,omitempty" jsonschema:"特殊格局。"`
	Basis       FuYiBasisView   `json:"basis" jsonschema:"扶抑推导线索。"`
}

// TiaoHouAvailabilityView is the readable调候 availability. Pointer fields
// allow the all-null secondary form (meaning "checked, found none").
type TiaoHouAvailabilityView struct {
	Stem        *string              `json:"gan" jsonschema:"天干。"`
	Transparent *bool                `json:"transparent" jsonschema:"是否透干。"`
	Hidden      *bool                `json:"hidden" jsonschema:"是否藏干。"`
	Occurrences []StemOccurrenceView `json:"occurrences" jsonschema:"出现位置。"`
	GuanXiLinks []RelationLinkView   `json:"guan_xi_links" jsonschema:"关系证据引用。"`
}

// StemOccurrenceView is a readable stem occurrence.
type StemOccurrenceView struct {
	Pillar *string `json:"zhu" jsonschema:"柱位。"`
	Branch *string `json:"zhi" jsonschema:"地支。"`
	Stem   *string `json:"gan" jsonschema:"天干。"`
	Source *string `json:"lai_yuan" jsonschema:"来源。"`
}

// TiaoHouView is the调候 candidate with an explicit all-null secondary form.
type TiaoHouView struct {
	PrimaryWuxing      string                  `json:"zhu_yong" jsonschema:"主调候五行。"`
	SecondaryWuxing    string                  `json:"fu_yong" jsonschema:"辅调候五行。"`
	Season             string                  `json:"season" jsonschema:"季节状态。"`
	Primary            TiaoHouAvailabilityView `json:"zhu" jsonschema:"主调候。"`
	Secondary          TiaoHouAvailabilityView `json:"fu" jsonschema:"辅调候（无候选时为全 null）。"`
	SecondaryCondition string                  `json:"fu_tiao_jian,omitempty" jsonschema:"辅调候条件。"`
}

// YongShenQuFaView groups the three取法 sources for用神 (扶抑/调候/格局).
type YongShenQuFaView struct {
	FuYi          FuYiView       `json:"fu_yi" jsonschema:"扶抑候选（旺衰见 wang_shuai）。"`
	TiaoHou       TiaoHouView    `json:"tiao_hou" jsonschema:"调候候选。"`
	GeJu          GeJuResultView `json:"ge_ju" jsonschema:"月令结构证据。"`
	GuanShaQuQing GuanSha        `json:"guan_sha_qu_qing" jsonschema:"官杀取清（格局结论）。"`
	CaiKu         CaiKu          `json:"cai_ku" jsonschema:"财库（格局结论）。"`
}

// GuanXiView groups canonical干支关系 facts. The nested key names are
// themselves命理域 terms, so no parent/child name collision remains.
type GuanXiView struct {
	HeHuiChongXing []NatalRelationFactView `json:"he_hui_chong_xing" jsonschema:"合会冲刑等命理关系。"`
	LuRoots        []LuRoot                `json:"lu_gen" jsonschema:"十神通根禄位。"`
	GongJia        []GongJia               `json:"gong_jia,omitempty" jsonschema:"三合三会拱局。"`
	NayinRel       []NayinRelEntry         `json:"na_yin_guan_xi" jsonschema:"纳音生克关系。"`
	PeiOuGong      PalaceState             `json:"pei_ou_gong" jsonschema:"夫妻宫状态（由日支关系派生）。"`
}

// WangShuaiView is the single source of旺衰 facts. The parent key already
// carries the domain, so children are named by their subject only.
type WangShuaiView struct {
	ShiShen []TenGodStateView  `json:"shi_shen" jsonschema:"十神显隐通根得令旺弱。"`
	Wuxing  []ElementStateView `json:"wuxing" jsonschema:"五行旺衰（唯一出处）。"`
}

// SanYuanView groups 胎元/命宫/身宫 without a redundant wrapper level.
type SanYuanView struct {
	TaiYuan  ganzhi.Zhu `json:"tai_yuan" jsonschema:"胎元。"`
	MingGong ganzhi.Zhu `json:"ming_gong" jsonschema:"命宫。"`
	ShenGong ganzhi.Zhu `json:"shen_gong" jsonschema:"身宫。"`
}

// ChangShengStageView is one stage of日主十二长生.
type ChangShengStageView struct {
	Ming string     `json:"ming" jsonschema:"长生阶段名。"`
	Zhi  ganzhi.Zhi `json:"zhi" jsonschema:"该阶段地支。"`
}

// ShenShaView groups 胎息/三奇 (盘面神煞事实).
type ShenShaView struct {
	TaiXi      ganzhi.Zhu      `json:"tai_xi" jsonschema:"胎息。"`
	SanQiName  string          `json:"san_qi_name,omitempty" jsonschema:"三奇名称。"`
	RiZhiXiang BranchImageKind `json:"ri_zhi_xiang" jsonschema:"日支取象（由日支神煞/墓库派生）。"`
}

// LiFaView groups历法事实 only.
type LiFaView struct {
	Zodiac    string `json:"sheng_xiao" jsonschema:"生肖。"`
	RiXun     string `json:"ri_xun" jsonschema:"日柱所在旬。"`
	RiXunKong string `json:"ri_xun_kong" jsonschema:"日柱旬空。"`
	BirthYear int    `json:"birth_year" jsonschema:"出生公历年份。"`
}

// FangFaKouJingView groups 方法/模型口径 separately from历法事实. Model IDs
// are engineering provenance, not chart facts, and belong here only.
type FangFaKouJingView struct {
	ZiShiRule     string            `json:"zi_shi_rule" jsonschema:"子时换日规则说明。"`
	ShenShaSchool ShenShaSchoolView `json:"shen_sha_school" jsonschema:"神煞双参照口径。"`
	FuYiMoXing    string            `json:"fu_yi_mo_xing" jsonschema:"扶抑模型标识。"`
	TiaoHouMoXing string            `json:"tiao_hou_mo_xing" jsonschema:"调候模型标识。"`
}

// ChartView is the domain-aggregated output of a Bazi chart. Every direct
// child is a命理域; method provenance is isolated in fang_fa_kou_jing.
type ChartView struct {
	Pillars       PillarsView           `json:"pillars" jsonschema:"四柱。"`
	DaYun         *DaYunView            `json:"da_yun,omitempty" jsonschema:"大运。"`
	YongShenQuFa  YongShenQuFaView      `json:"yong_shen_qu_fa" jsonschema:"用神三法（扶抑/调候/格局）。"`
	GuanXi        GuanXiView            `json:"guan_xi" jsonschema:"干支关系事实。"`
	WangShuai     WangShuaiView         `json:"wang_shuai" jsonschema:"十神与五行旺衰。"`
	SanYuan       SanYuanView           `json:"san_yuan" jsonschema:"胎元命宫身宫。"`
	ChangSheng    []ChangShengStageView `json:"chang_sheng" jsonschema:"日主十二长生。"`
	ShenSha       ShenShaView           `json:"shen_sha" jsonschema:"胎息三奇。"`
	LiFa          LiFaView              `json:"li_fa" jsonschema:"历法事实。"`
	FangFaKouJing FangFaKouJingView     `json:"fang_fa_kou_jing" jsonschema:"方法与模型口径。"`
}

// View produces the domain-aggregated output from the computed chart.
func (c Chart) View() ChartView {
	stages := make([]ChangShengStageView, 0, len(c.ChangSheng))
	for _, stage := range c.ChangSheng {
		stages = append(stages, ChangShengStageView{Ming: stage.Name, Zhi: stage.Zhi})
	}
	return ChartView{
		Pillars: PillarsView{Nian: pillarView(c.Nian), Yue: pillarView(c.Yue), Ri: pillarView(c.Ri), Shi: pillarView(c.Shi)},
		DaYun:   daYunView(c.DaYun),
		YongShenQuFa: YongShenQuFaView{
			FuYi:          fuyiView(c.FuYi),
			TiaoHou:       tiaoHouView(c.TiaoHou),
			GeJu:          geJuResultView(c.GeJu),
			GuanShaQuQing: c.GuanSha,
			CaiKu:         c.CaiKu,
		},
		GuanXi: GuanXiView{
			HeHuiChongXing: natalRelationFactViews(c.GuanXi), LuRoots: c.LuRoots, GongJia: c.GongJia,
			NayinRel: c.NayinRel, PeiOuGong: c.SpousePalaceState,
		},
		WangShuai: WangShuaiView{ShiShen: tenGodStateViews(c.TenGodStates), Wuxing: elementStateViews(c.ElementStates)},
		SanYuan: SanYuanView{
			TaiYuan: c.SanYuan.TaiYuan, MingGong: c.SanYuan.MingGong, ShenGong: c.SanYuan.ShenGong,
		},
		ChangSheng: stages,
		ShenSha:    ShenShaView{TaiXi: c.TaiXi, SanQiName: c.SanQiName, RiZhiXiang: c.DayBranchType},
		LiFa: LiFaView{
			Zodiac: c.Zodiac, RiXun: c.DayXun, RiXunKong: c.DayXunKong, BirthYear: c.BirthYear,
		},
		FangFaKouJing: FangFaKouJingView{
			ZiShiRule: c.ZiShiRule, ShenShaSchool: shenShaSchoolView(c.ShenShaSchool),
			FuYiMoXing: c.FuYi.ModelID, TiaoHouMoXing: c.TiaoHou.ModelID,
		},
	}
}

func fuyiView(f FuYiResult) FuYiView {
	return FuYiView{
		WuxingCount: f.WuxingCount,
		Yong:        f.Yong,
		Xi:          f.Xi,
		Ji:          f.Ji,
		Strength:    f.Strength,
		Pattern:     f.Pattern,
		Basis:       fuYiBasisView(f.Basis),
	}
}

func stringPtr(s string) *string { return &s }

func tiaoHouAvailabilityView(a TiaoHouAvailability) TiaoHouAvailabilityView {
	var stemPtr *string
	if a.Stem != "" {
		stemPtr = stringPtr(a.Stem)
	}
	v := TiaoHouAvailabilityView{
		Stem:        stemPtr,
		Transparent: &a.Transparent,
		Hidden:      &a.Hidden,
		GuanXiLinks: relationLinkViews(a.GuanXiLinks),
	}
	v.Occurrences = make([]StemOccurrenceView, 0, len(a.Occurrences))
	for _, occ := range a.Occurrences {
		pillar := DecodePillarIndex(occ.Pillar)
		branch := ganzhi.ZhiName(occ.Branch)
		stem := ganzhi.GanName(occ.Stem)
		source := displayFromBazi(tenGodSourceDisplay, string(occ.Source))
		v.Occurrences = append(v.Occurrences, StemOccurrenceView{
			Pillar: &pillar, Branch: &branch, Stem: &stem, Source: &source,
		})
	}
	return v
}

func tiaoHouView(th TiaoHouResult) TiaoHouView {
	var secondary TiaoHouAvailability
	if th.Secondary != nil {
		secondary = *th.Secondary
	}
	return TiaoHouView{
		PrimaryWuxing:      th.PrimaryWuxing,
		SecondaryWuxing:    th.SecondaryWuxing,
		Season:             th.Season,
		Primary:            tiaoHouAvailabilityView(th.Primary),
		Secondary:          tiaoHouAvailabilityView(secondary),
		SecondaryCondition: th.SecondaryCondition,
	}
}
