package bazi

import "github.com/ml8s/mingligo/pkg/ganzhi"

// This file defines the domain-owned, domain-aggregated output view of a Bazi
// chart. It groups related facts into命理子域 objects and is the single output
// model that every transport adapter transmits unchanged. jsonschema tags are
// inert metadata read only when a consumer reflects the type for its schema.

// PillarsView groups the four pillars (干支/藏干/十神投影/神煞/纳音/旬空/长生落点).
type PillarsView struct {
	Nian Pillar `json:"nian" jsonschema:"年柱。"`
	Yue  Pillar `json:"yue" jsonschema:"月柱。"`
	Ri   Pillar `json:"ri" jsonschema:"日柱。"`
	Shi  Pillar `json:"shi" jsonschema:"时柱。"`
}

// FuYiView is the扶抑 candidate. 旺衰 deliberately omitted: the season map is
// published once in StatesView.WuxingStates (domain source field kept for the
// computation that derives WuxingStates).
type FuYiView struct {
	WuxingCount map[string]int  `json:"wuxing_count" jsonschema:"五行个数。"`
	Yong        CandidateWuxing `json:"yong" jsonschema:"扶抑用神候选。"`
	Xi          CandidateWuxing `json:"xi" jsonschema:"扶抑喜神候选。"`
	Ji          CandidateWuxing `json:"ji" jsonschema:"扶抑忌神候选。"`
	Strength    string          `json:"qiang_ruo" jsonschema:"身强身弱。"`
	Pattern     string          `json:"ge_ju,omitempty" jsonschema:"特殊格局。"`
	ModelID     string          `json:"fu_yi_model_id" jsonschema:"扶抑模型 ID。"`
	Basis       FuYiBasis       `json:"basis" jsonschema:"扶抑推导线索。"`
}

// TiaoHouAvailabilityView is the readable调候 availability. Pointer fields
// allow the all-null secondary form (meaning "checked, found none").
type TiaoHouAvailabilityView struct {
	Stem        *string              `json:"gan" jsonschema:"天干。"`
	Transparent *bool                `json:"transparent" jsonschema:"是否透干。"`
	Hidden      *bool                `json:"hidden" jsonschema:"是否藏干。"`
	Occurrences []StemOccurrenceView `json:"occurrences" jsonschema:"出现位置。"`
	GuanXiLinks []RelationLink       `json:"guan_xi_links" jsonschema:"关系证据引用。"`
}

// StemOccurrenceView is a readable stem occurrence.
type StemOccurrenceView struct {
	Pillar *string `json:"zhu" jsonschema:"柱位。"`
	Branch *string `json:"zhi" jsonschema:"地支。"`
	Stem   *string `json:"gan" jsonschema:"天干。"`
	Source *string `json:"source" jsonschema:"来源。"`
}

// TiaoHouView is the调候 candidate with an explicit all-null secondary form.
type TiaoHouView struct {
	PrimaryWuxing      string                  `json:"primary_wuxing" jsonschema:"主调候五行。"`
	SecondaryWuxing    string                  `json:"secondary_wuxing" jsonschema:"辅调候五行。"`
	Season             string                  `json:"season" jsonschema:"季节状态。"`
	ModelID            string                  `json:"tiao_hou_model_id" jsonschema:"调候模型 ID。"`
	Primary            TiaoHouAvailabilityView `json:"primary" jsonschema:"主调候。"`
	Secondary          TiaoHouAvailabilityView `json:"secondary" jsonschema:"辅调候（无候选时为全 null）。"`
	SecondaryCondition string                  `json:"secondary_condition,omitempty" jsonschema:"辅调候条件。"`
}

// YongShenQuFaView groups the three取法 sources for用神 (扶抑/调候/格局).
type YongShenQuFaView struct {
	FuYi    FuYiView    `json:"fu_yi" jsonschema:"扶抑候选（旺衰见 states）。"`
	TiaoHou TiaoHouView `json:"tiao_hou" jsonschema:"调候候选。"`
	GeJu    GeJuResult  `json:"ge_ju" jsonschema:"月令结构证据。"`
}

// GuanXiView groups canonical干支关系 facts (合会冲刑/禄根/拱夹/纳音).
type GuanXiView struct {
	GuanXi   []NatalRelationFact `json:"guan_xi" jsonschema:"合会冲刑等命理关系。"`
	LuRoots  []LuRoot            `json:"lu_roots" jsonschema:"十神通根禄位。"`
	GongJia  []GongJia           `json:"gong_jia,omitempty" jsonschema:"三合三会拱局。"`
	NayinRel []NayinRelEntry     `json:"na_yin_rel" jsonschema:"纳音生克关系。"`
}

// WangShuaiView groups十神/五行旺衰 facts. 旺衰 published only here.
type WangShuaiView struct {
	ShiShenStates []TenGodState  `json:"shi_shen_states" jsonschema:"十神显隐通根得令旺弱。"`
	WuxingStates  []ElementState `json:"wuxing_states" jsonschema:"五行旺衰（唯一出处）。"`
}

// SanYuanView groups 胎元/命宫/身宫.
type SanYuanView struct {
	SanYuan SanYuan `json:"san_yuan" jsonschema:"胎元命宫身宫。"`
}

// ChangShengView groups 日主十二长生全表.
type ChangShengView struct {
	ChangSheng [12]ChangShengStage `json:"chang_sheng" jsonschema:"日主十二长生全表。"`
}

// ShenShaView groups 胎息/三奇 (盘面神煞事实).
type ShenShaView struct {
	TaiXi     ganzhi.Zhu `json:"tai_xi" jsonschema:"胎息。"`
	SanQiName string     `json:"san_qi_name,omitempty" jsonschema:"三奇名称。"`
}

// LiFaKouJingView groups历法/排盘口径 facts only.
type LiFaKouJingView struct {
	Zodiac        string        `json:"sheng_xiao" jsonschema:"生肖。"`
	RiXun         string        `json:"ri_xun" jsonschema:"日柱所在旬。"`
	RiXunKong     string        `json:"ri_xun_kong" jsonschema:"日柱旬空。"`
	ZiShiRule     string        `json:"zi_shi_rule" jsonschema:"子时换日规则说明。"`
	ShenShaSchool ShenShaSchool `json:"shen_sha_school" jsonschema:"神煞双参照口径。"`
	BirthYear     int           `json:"birth_year" jsonschema:"出生公历年份。"`
}

// ChartView is the domain-aggregated output of a Bazi chart.
type ChartView struct {
	Pillars      PillarsView      `json:"pillars" jsonschema:"四柱。"`
	DaYun        *DaYun           `json:"da_yun,omitempty" jsonschema:"大运。"`
	YongShenQuFa YongShenQuFaView `json:"yong_shen_qu_fa" jsonschema:"用神三法（扶抑/调候/格局）。"`
	GuanXi       GuanXiView       `json:"guan_xi" jsonschema:"干支关系事实。"`
	WangShuai    WangShuaiView    `json:"wang_shuai" jsonschema:"十神与五行旺衰。"`
	SanYuan      SanYuanView      `json:"san_yuan_view" jsonschema:"胎元命宫身宫。"`
	ChangSheng   ChangShengView   `json:"chang_sheng_view" jsonschema:"日主十二长生。"`
	ShenSha      ShenShaView      `json:"shen_sha" jsonschema:"胎息三奇。"`
	GuanSha      GuanSha          `json:"guan_sha" jsonschema:"官杀取清。"`
	CaiKu        CaiKu            `json:"cai_ku" jsonschema:"财库。"`
	PeiOuGong    PalaceState      `json:"pei_ou_gong" jsonschema:"夫妻宫状态。"`
	RiZhiXiang   BranchImageKind  `json:"ri_zhi_xiang" jsonschema:"日支取象。"`
	LiFaKouJing  LiFaKouJingView  `json:"li_fa_kou_jing" jsonschema:"历法与排盘口径。"`
}

// View produces the domain-aggregated output from the computed chart. Filling
// order is free; each fact is placed into its owning命理子域.
func (c Chart) View() ChartView {
	return ChartView{
		Pillars: PillarsView{Nian: c.Nian, Yue: c.Yue, Ri: c.Ri, Shi: c.Shi},
		DaYun:   c.DaYun,
		YongShenQuFa: YongShenQuFaView{
			FuYi:    fuyiView(c.FuYi),
			TiaoHou: tiaoHouView(c.TiaoHou),
			GeJu:    c.GeJu,
		},
		GuanXi: GuanXiView{
			GuanXi: c.GuanXi, LuRoots: c.LuRoots, GongJia: c.GongJia, NayinRel: c.NayinRel,
		},
		WangShuai:  WangShuaiView{ShiShenStates: c.TenGodStates, WuxingStates: c.ElementStates},
		SanYuan:    SanYuanView{SanYuan: c.SanYuan},
		ChangSheng: ChangShengView{ChangSheng: c.ChangSheng},
		ShenSha:    ShenShaView{TaiXi: c.TaiXi, SanQiName: c.SanQiName},
		GuanSha:    c.GuanSha,
		CaiKu:      c.CaiKu,
		PeiOuGong:  c.SpousePalaceState,
		RiZhiXiang: c.DayBranchType,
		LiFaKouJing: LiFaKouJingView{
			Zodiac: c.Zodiac, RiXun: c.DayXun, RiXunKong: c.DayXunKong, ZiShiRule: c.ZiShiRule,
			ShenShaSchool: c.ShenShaSchool, BirthYear: c.BirthYear,
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
		ModelID:     f.ModelID,
		Basis:       f.Basis,
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
		GuanXiLinks: a.GuanXiLinks,
	}
	v.Occurrences = make([]StemOccurrenceView, 0, len(a.Occurrences))
	for _, occ := range a.Occurrences {
		pillar := DecodePillarIndex(occ.Pillar)
		branch := ganzhi.ZhiName(occ.Branch)
		stem := ganzhi.GanName(occ.Stem)
		source := string(occ.Source)
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
		ModelID:            th.ModelID,
		Primary:            tiaoHouAvailabilityView(th.Primary),
		Secondary:          tiaoHouAvailabilityView(secondary),
		SecondaryCondition: th.SecondaryCondition,
	}
}
