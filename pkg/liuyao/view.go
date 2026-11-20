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
	Moving      bool   `json:"moving"`
	YuePo       bool   `json:"yue_po"`
	XunKong     bool   `json:"xun_kong"`
	MuKu        bool   `json:"mu_ku"`
	MuKuBranch  string `json:"mu_ku_zhi,omitempty"`
	MuKuElement string `json:"mu_ku_wuxing,omitempty"`
	DongSheng   bool   `json:"dong_sheng"`
	DongKe      bool   `json:"dong_ke"`
}

type LineView struct {
	Position        int               `json:"position"`
	YaoValue        int               `json:"yao_value"`
	GanZhi          string            `json:"gan_zhi"`
	Wuxing          string            `json:"wuxing"`
	LiuQin          string            `json:"liu_qin"`
	ShiYing         string            `json:"shi_ying,omitempty"`
	LiuShou         string            `json:"liu_shou"`
	RiChenRelations map[string][]string `json:"ri_chen_relations,omitempty"`
	WangShuai       string            `json:"wang_shuai,omitempty"`
	ChangShengYue   string            `json:"chang_sheng_yue,omitempty"`
	Flags           LineFlags         `json:"flags"`
}

type PatternView struct {
	ID          string `json:"id"`
	PatternKind string `json:"ge_ju_kind"`
	SubType     string `json:"sub_type"`
	Position    int    `json:"position"`
	IsTrue      bool   `json:"is_true"`
}

type GuaCiView struct {
	Name  string   `json:"name"`
	Text  string   `json:"text"`
	YaoCi []string `json:"yao_ci"`
}

type FuShenView struct {
	Position int    `json:"position"`
	LiuQin   string `json:"liu_qin"`
	Zhi      string `json:"zhi"`
}

type YongShenView struct {
	Name             string        `json:"name"`
	Position         int           `json:"position"`
	IsHidden         bool          `json:"is_hidden,omitempty"`
	WangShuai        string        `json:"wang_shuai"`
	YuePo            bool          `json:"yue_po,omitempty"`
	XunKong          bool          `json:"xun_kong,omitempty"`
	MuKu             bool          `json:"mu_ku,omitempty"`
	LiuShou          string        `json:"liu_shou,omitempty"`
	ChangSheng       string        `json:"chang_sheng,omitempty"`
	FuShen           *FuShenView   `json:"fu_shen,omitempty"`
	FuShenCandidates []FuShenView  `json:"fu_shen_candidates,omitempty"`
}

type AnalysisView struct {
	TimingCandidates      []TimingCandidate       `json:"timing_candidates,omitempty"`
	DongYaoRelations      []DongYaoRelation       `json:"dong_yao_relations,omitempty"`
	HiddenLines           []HiddenLine            `json:"hidden_lines,omitempty"`
	BranchRelationFacts   []BranchRelationFact    `json:"zhi_relation_facts,omitempty"`
	TombSchool            TombSchool              `json:"tomb_school"`
	DayClashFacts         []DayClashFact          `json:"day_clash_facts,omitempty"`
	MovingTransformations []MovingTransformation  `json:"moving_transformations,omitempty"`
	SanHeCandidates       []SanHeCandidate        `json:"san_he_candidates,omitempty"`
	ForceChain            *ForceChain             `json:"yong_shen_force_chain,omitempty"`
	YongShenCandidates    []YongShenCandidate     `json:"yong_shen_candidates,omitempty"`
	Conflicts             []Conflict              `json:"conflicts,omitempty"`
}

type ChartView struct {
	GuaName      string         `json:"gua_name"`
	BianGua      *string        `json:"bian_gua"`
	Gong         string         `json:"gong"`
	GongWuxing   string         `json:"gong_wuxing"`
	Lines        []LineView     `json:"lines"`
	BianLines    []LineView     `json:"bian_lines"`
	RiChenGan    string         `json:"ri_chen_gan"`
	RiChenZhi    string         `json:"ri_chen_zhi"`
	YueJianGan   string         `json:"yue_jian_gan"`
	YueJianZhi   string         `json:"yue_jian_zhi"`
	XunKong      []string       `json:"xun_kong"`
	Patterns     []PatternView  `json:"ge_ju"`
	YongShen     YongShenView   `json:"yong_shen"`
	GuaCi        *GuaCiView     `json:"gua_ci,omitempty"`
	Analysis     AnalysisView   `json:"analysis"`
}

func (c Chart) lineView(line Line, rel DayRelation, wangShuai ganzhi.WangShuai) LineView {
	flags := LineFlags{
		Moving: line.DongSelf, YuePo: line.YuePo, XunKong: line.XunKong,
		MuKu: line.MuKu, MuKuBranch: line.MuKuBranch, MuKuElement: line.MuKuElement,
		DongSheng: line.DongSheng, DongKe: line.DongKe,
	}
	v := LineView{
		Position: line.Position, YaoValue: int(line.Type),
		GanZhi: ganzhi.GanName(line.Gan) + ganzhi.ZhiName(line.Zhi),
		Wuxing: line.Wuxing.String(), LiuQin: line.LiuQin.String(),
		ShiYing: line.ShiYing.String(), LiuShou: line.LiuShou.String(),
		WangShuai: wangShuai.String(), ChangShengYue: line.ChangShengYue,
		Flags: flags,
	}
	if len(rel.Relations) > 0 {
		names := make([]string, 0, len(rel.Relations))
		for _, item := range rel.Relations {
			names = append(names, item.String())
		}
		v.RiChenRelations = map[string][]string{"relations": names}
	}
	return v
}

func fuShenView(v *FuShen) *FuShenView {
	if v == nil {
		return nil
	}
	return &FuShenView{Position: v.Position, LiuQin: v.LiuQin.String(), Zhi: ganzhi.ZhiName(v.Zhi)}
}

func yongShenView(v YongShenResult) YongShenView {
	view := YongShenView{
		Name: v.Name, Position: v.Position, IsHidden: v.IsHidden,
		WangShuai: v.WangShuai.String(), YuePo: v.YuePo, XunKong: v.XunKong,
		MuKu: v.MuKu, LiuShou: v.LiuShou.String(), ChangSheng: v.ChangSheng,
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
		Gong: c.Palace, GongWuxing: c.PalaceWuxing.String(),
		Lines: make([]LineView, 0, 6), BianLines: make([]LineView, 0, 6),
		RiChenGan: ganzhi.GanName(c.RiGan), RiChenZhi: ganzhi.ZhiName(c.RiZhi),
		YueJianGan: ganzhi.GanName(c.YueGan), YueJianZhi: ganzhi.ZhiName(c.YueZhi),
		XunKong: []string{ganzhi.ZhiName(c.XunKong[0]), ganzhi.ZhiName(c.XunKong[1])},
		Patterns: make([]PatternView, 0, len(c.Patterns)),
		YongShen: yongShenView(c.YongShen),
	}
	if int(c.BianGua) != 0 {
		name := c.BianGua.Name()
		view.BianGua = &name
	}
	for i, line := range c.Lines {
		view.Lines = append(view.Lines, c.lineView(line, c.DayRelations[i], c.WangShuai[i]))
	}
	for _, line := range c.BianLines {
		if line.Position != 0 {
			view.BianLines = append(view.BianLines, c.lineView(line, DayRelation{}, ganzhi.WangShuai(-1)))
		}
	}
	for i, p := range c.Patterns {
		view.Patterns = append(view.Patterns, PatternView{
			ID: fmt.Sprintf("pattern-%d", i+1), PatternKind: string(p.Type),
			SubType: p.SubType, Position: p.Position, IsTrue: p.IsTrue,
		})
	}
	if c.GuaCi.Name != "" {
		view.GuaCi = &GuaCiView{
			Name: c.GuaCi.Name, Text: c.GuaCi.GuaCi,
			YaoCi: append([]string(nil), c.GuaCi.YaoCi[:]...),
		}
	}
	view.Analysis = c.Analysis()
	return view
}

// Analysis returns the domain-aggregated analysis view: the deterministic
// judgment facts grouped into one section.
func (c Chart) Analysis() AnalysisView {
	return AnalysisView{
		TimingCandidates:      c.TimingCandidates,
		DongYaoRelations:      c.DongYaoRelations,
		HiddenLines:           c.HiddenLines,
		BranchRelationFacts:   c.BranchRelationFacts,
		TombSchool:            c.TombSchool,
		DayClashFacts:         c.DayClashFacts,
		MovingTransformations: c.MovingTransformations,
		SanHeCandidates:       c.SanHeCandidates,
		ForceChain:            c.ForceChain,
		YongShenCandidates:    c.YongShenCandidates,
		Conflicts:             c.Conflicts,
	}
}
