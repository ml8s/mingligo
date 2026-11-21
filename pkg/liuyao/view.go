package liuyao

import (
	"fmt"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// This file defines the domain-owned, domain-aggregated output view of a 六爻
// chart. The view groups related facts (a line's 干支/六亲/旺衰/flags together,
// analysis facts into one section) and uses readable values. It is the single
// output model that every transport adapter (MCP, REST) transmits unchanged.

type LineFlags struct {
	Moving      bool     `json:"moving" jsonschema:"是否发动。"`
	YuePo       bool     `json:"yue_po" jsonschema:"是否月破。"`
	XunKong     bool     `json:"xun_kong" jsonschema:"是否旬空。"`
	MuKu        bool     `json:"mu_ku" jsonschema:"是否墓库。"`
	MuKuBranch  string   `json:"mu_ku_zhi,omitempty" jsonschema:"墓库地支。"`
	MuKuElement string   `json:"mu_ku_wuxing,omitempty" jsonschema:"墓库五行。"`
	MuKuTypes   []string `json:"mu_ku_types,omitempty" jsonschema:"入墓来源：日墓/动墓/化墓。"`
	DongSheng   bool     `json:"dong_sheng" jsonschema:"动爻生。"`
	DongKe      bool     `json:"dong_ke" jsonschema:"动爻克。"`
}

type LineView struct {
	Position        int                 `json:"position" jsonschema:"爻位。"`
	YaoValue        int                 `json:"yao_value" jsonschema:"爻值。"`
	GanZhi          string              `json:"gan_zhi" jsonschema:"干支。"`
	Wuxing          string              `json:"wuxing" jsonschema:"五行。"`
	LiuQin          string              `json:"liu_qin" jsonschema:"六亲。"`
	ShiYing         string              `json:"shi_ying,omitempty" jsonschema:"世应。"`
	LiuShou         string              `json:"liu_shou" jsonschema:"六神。"`
	RiChenRelations map[string][]string `json:"ri_chen_relations,omitempty" jsonschema:"日辰关系。"`
	WangShuai       string              `json:"wang_shuai,omitempty" jsonschema:"旺衰。"`
	ChangShengYue   string              `json:"chang_sheng_yue,omitempty" jsonschema:"月令长生。"`
	Flags           LineFlags           `json:"flags" jsonschema:"爻状态标志。"`
}

type PatternView struct {
	ID          string `json:"id" jsonschema:"标识。"`
	PatternKind string `json:"ge_ju_kind" jsonschema:"格局类型。"`
	SubType     string `json:"sub_type" jsonschema:"格局子类。"`
	Position    int    `json:"position" jsonschema:"爻位。"`
	IsTrue      bool   `json:"is_true" jsonschema:"是否真格局。"`
}

type GuaCiView struct {
	Name  string   `json:"name" jsonschema:"名称。"`
	Text  string   `json:"text" jsonschema:"卦辞文本。"`
	YaoCi []string `json:"yao_ci" jsonschema:"爻辞。"`
}

type FuShenView struct {
	Position int    `json:"position" jsonschema:"爻位。"`
	LiuQin   string `json:"liu_qin" jsonschema:"六亲。"`
	Zhi      string `json:"zhi" jsonschema:"地支。"`
}

type YongShenView struct {
	Name             string       `json:"name" jsonschema:"名称。"`
	Position         int          `json:"position" jsonschema:"爻位。"`
	IsHidden         bool         `json:"is_hidden,omitempty" jsonschema:"是否伏藏。"`
	WangShuai        string       `json:"wang_shuai" jsonschema:"旺衰。"`
	YuePo            bool         `json:"yue_po,omitempty" jsonschema:"是否月破。"`
	XunKong          bool         `json:"xun_kong,omitempty" jsonschema:"是否旬空。"`
	MuKu             bool         `json:"mu_ku,omitempty" jsonschema:"是否墓库。"`
	LiuShou          string       `json:"liu_shou,omitempty" jsonschema:"六神。"`
	ChangSheng       string       `json:"chang_sheng,omitempty" jsonschema:"长生。"`
	FuShen           *FuShenView  `json:"fu_shen,omitempty" jsonschema:"本宫全量伏神表（含飞神/飞伏关系/出入条件）；所选用神的伏神见 yong_shen.fu_shen。"`
	FuShenCandidates []FuShenView `json:"fu_shen_candidates,omitempty" jsonschema:"伏神候选。"`
}

// ShiQiView groups 应期候选 (candidates, not a verdict).
type ShiQiView struct {
	Candidates []TimingCandidate `json:"candidates,omitempty" jsonschema:"应期候选（非定论）。"`
}

// YongShenCengView groups 用神层 facts (力量链与取舍候选).
type YongShenCengView struct {
	ForceChain *ForceChain         `json:"force_chain,omitempty" jsonschema:"用神力量链。"`
	Candidates []YongShenCandidate `json:"candidates,omitempty" jsonschema:"用神候选（含确定性取舍）。"`
}

// GeJuCengView groups 格局结构层 (六爻格局与三合结构).
type GeJuCengView struct {
	Patterns        []PatternView    `json:"patterns,omitempty" jsonschema:"六爻格局。"`
	SanHeCandidates []SanHeCandidate `json:"san_he_candidates,omitempty" jsonschema:"三合结构候选。"`
}

// DongYaoGuanXiView groups 动爻关系.
type DongYaoGuanXiView struct {
	Relations       []DongYaoRelation      `json:"relations,omitempty" jsonschema:"动爻关系。"`
	Transformations []MovingTransformation `json:"transformations,omitempty" jsonschema:"动变。"`
}

// FuShenCengView groups 伏神层 facts.
type FuShenCengView struct {
	HiddenLines []HiddenLine `json:"hidden_lines,omitempty" jsonschema:"本宫全量伏神表。"`
}

// ZhiGuanXiView groups 地支关系与日冲.
type ZhiGuanXiView struct {
	Relations  []BranchRelationFact `json:"relations,omitempty" jsonschema:"地支关系。"`
	DayClashes []DayClashFact       `json:"day_clashes,omitempty" jsonschema:"日冲事实。"`
}

// MaoDunXinHao groups 矛盾信号 (待解释层裁决).
type MaoDunXinHao struct {
	Signals []Conflict `json:"signals,omitempty" jsonschema:"矛盾信号（非吉凶裁决）。"`
}

type ChartView struct {
	GuaName       string            `json:"gua_name" jsonschema:"卦名。"`
	BianGua       *string           `json:"bian_gua" jsonschema:"变卦。"`
	Gong          string            `json:"gong" jsonschema:"宫位。"`
	GongWuxing    string            `json:"gong_wuxing" jsonschema:"宫位五行。"`
	LiuYao        []LineView        `json:"liu_yao" jsonschema:"本卦六爻。"`
	BianYao       []LineView        `json:"bian_yao" jsonschema:"变卦六爻（仅动变事实）。"`
	RiChenGan     string            `json:"ri_chen_gan" jsonschema:"日辰天干。"`
	RiChenZhi     string            `json:"ri_chen_zhi" jsonschema:"日辰地支。"`
	YueJianGan    string            `json:"yue_jian_gan" jsonschema:"月建天干。"`
	YueJianZhi    string            `json:"yue_jian_zhi" jsonschema:"月建地支。"`
	XunKong       []string          `json:"xun_kong" jsonschema:"旬空地支（两字）。"`
	GeJuCeng      GeJuCengView      `json:"ge_ju_ceng" jsonschema:"格局结构层。"`
	YongShen      YongShenView      `json:"yong_shen" jsonschema:"用神。"`
	YongShenCeng  YongShenCengView  `json:"yong_shen_ceng" jsonschema:"用神层。"`
	ShiQi         ShiQiView         `json:"shi_qi" jsonschema:"应期。"`
	DongYaoGuanXi DongYaoGuanXiView `json:"dong_yao_guan_xi" jsonschema:"动爻关系。"`
	FuShenCeng    FuShenCengView    `json:"fu_shen_ceng" jsonschema:"伏神层。"`
	ZhiGuanXi     ZhiGuanXiView     `json:"zhi_guan_xi" jsonschema:"地支关系。"`
	MaoDunXinHao  MaoDunXinHao      `json:"mao_dun_xin_hao" jsonschema:"矛盾信号。"`
	MuKuLiuPai    TombSchool        `json:"mu_ku_liu_pai" jsonschema:"墓库流派（方法元数据）。"`
	GuaCi         *GuaCiView        `json:"gua_ci,omitempty" jsonschema:"卦辞。"`
}

func (c Chart) lineView(line Line, rel DayRelation, wangShuai ganzhi.WangShuai) LineView {
	flags := LineFlags{
		Moving: line.DongSelf, YuePo: line.YuePo, XunKong: line.XunKong,
		MuKu: line.MuKu, MuKuBranch: line.MuKuBranch, MuKuElement: line.MuKuElement,
		DongSheng: line.DongSheng, DongKe: line.DongKe,
	}
	for _, source := range line.MuKuTypes {
		flags.MuKuTypes = append(flags.MuKuTypes, DecodeMuKuSource(source))
	}
	v := LineView{
		Position: line.Position, YaoValue: int(line.Type),
		GanZhi: ganzhi.GanName(line.Gan) + ganzhi.ZhiName(line.Zhi),
		Wuxing: ganzhi.DecodeWuxing(line.Wuxing), LiuQin: DecodeLiuQin(line.LiuQin),
		ShiYing: DecodeShiYing(line.ShiYing), LiuShou: DecodeLiuShou(line.LiuShou),
		WangShuai: ganzhi.DecodeWangShuai(wangShuai), ChangShengYue: string(line.ChangShengYue),
		Flags: flags,
	}
	if len(rel.Relations) > 0 {
		names := make([]string, 0, len(rel.Relations))
		for _, item := range rel.Relations {
			names = append(names, DecodeDayRelationType(item))
		}
		v.RiChenRelations = map[string][]string{"relations": names}
	}
	return v
}

func fuShenView(v *FuShen) *FuShenView {
	if v == nil {
		return nil
	}
	return &FuShenView{Position: v.Position, LiuQin: DecodeLiuQin(v.LiuQin), Zhi: ganzhi.ZhiName(v.Zhi)}
}

func yongShenView(v YongShenResult) YongShenView {
	view := YongShenView{
		Name: v.Name, Position: v.Position, IsHidden: v.IsHidden,
		WangShuai: ganzhi.DecodeWangShuai(v.WangShuai), YuePo: v.YuePo, XunKong: v.XunKong,
		MuKu: v.MuKu, LiuShou: DecodeLiuShou(v.LiuShou), ChangSheng: string(v.ChangSheng),
		FuShen: fuShenView(v.FuShen),
	}
	for _, cand := range v.FuShenCandidates {
		view.FuShenCandidates = append(view.FuShenCandidates, *fuShenView(&cand))
	}
	return view
}

// View returns the domain-aggregated, readable chart view. It groups the
// parallel line arrays into whole-line objects and encodes enum names; this is
// domain knowledge, not presentation.
func (c Chart) View() ChartView {
	view := ChartView{
		GuaName: c.BenGua.Name(),
		Gong:    c.Palace, GongWuxing: ganzhi.DecodeWuxing(c.PalaceWuxing),
		LiuYao: make([]LineView, 0, 6), BianYao: make([]LineView, 0, 6),
		RiChenGan: ganzhi.GanName(c.RiGan), RiChenZhi: ganzhi.ZhiName(c.RiZhi),
		YueJianGan: ganzhi.GanName(c.YueGan), YueJianZhi: ganzhi.ZhiName(c.YueZhi),
		XunKong:       []string{ganzhi.ZhiName(c.XunKong[0]), ganzhi.ZhiName(c.XunKong[1])},
		GeJuCeng:      GeJuCengView{Patterns: make([]PatternView, 0, len(c.Patterns)), SanHeCandidates: c.SanHeCandidates},
		YongShen:      yongShenView(c.YongShen),
		YongShenCeng:  YongShenCengView{ForceChain: c.ForceChain, Candidates: c.YongShenCandidates},
		ShiQi:         ShiQiView{Candidates: c.TimingCandidates},
		DongYaoGuanXi: DongYaoGuanXiView{Relations: c.DongYaoRelations, Transformations: c.MovingTransformations},
		FuShenCeng:    FuShenCengView{HiddenLines: c.HiddenLines},
		ZhiGuanXi:     ZhiGuanXiView{Relations: c.BranchRelationFacts, DayClashes: c.DayClashFacts},
		MaoDunXinHao:  MaoDunXinHao{Signals: c.Conflicts},
		MuKuLiuPai:    c.TombSchool,
	}
	if int(c.BianGua) != 0 {
		name := c.BianGua.Name()
		view.BianGua = &name
	}
	for i, line := range c.Lines {
		view.LiuYao = append(view.LiuYao, c.lineView(line, c.DayRelations[i], c.WangShuai[i]))
	}
	for _, line := range c.BianLines {
		if line.Position != 0 {
			view.BianYao = append(view.BianYao, c.lineView(line, DayRelation{}, ganzhi.WangShuai(-1)))
		}
	}
	for i, p := range c.Patterns {
		view.GeJuCeng.Patterns = append(view.GeJuCeng.Patterns, PatternView{
			ID: fmt.Sprintf("pattern-%d", i+1), PatternKind: string(p.Type),
			SubType: string(p.SubType), Position: p.Position, IsTrue: p.IsTrue,
		})
	}
	if c.GuaCi.Name != "" {
		view.GuaCi = &GuaCiView{
			Name: c.GuaCi.Name, Text: c.GuaCi.GuaCi,
			YaoCi: append([]string(nil), c.GuaCi.YaoCi[:]...),
		}
	}
	return view
}
