package liuyao

import (
	"encoding/json"
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
	MuKuBranch  string `json:"mu_ku_branch,omitempty"`
	MuKuElement string `json:"mu_ku_element,omitempty"`
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
	PatternKind string `json:"pattern_kind"`
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

type TimingView struct {
	ID            string   `json:"id"`
	Mechanism     string   `json:"mechanism"`
	Position      int      `json:"position,omitempty"`
	TriggerBranch string   `json:"trigger_branch,omitempty"`
	Window        string   `json:"window,omitempty"`
	Basis         []string `json:"basis,omitempty"`
}

type DongRelationView struct {
	Position  int      `json:"position"`
	Relations []string `json:"relations"`
}

type HiddenLineView struct {
	Position       int    `json:"position"`
	LiuQin         string `json:"liu_qin"`
	Gan            string `json:"gan"`
	Zhi            string `json:"zhi"`
	Wuxing         string `json:"wuxing"`
	FlyingPosition int    `json:"flying_position"`
	FlyingLiuQin   string `json:"flying_liu_qin"`
	FlyingGan      string `json:"flying_gan"`
	FlyingZhi      string `json:"flying_zhi"`
	FlyingWuxing   string `json:"flying_wuxing"`
	FeiFuRelation  string `json:"fei_fu_relation"`
	XunKong        bool   `json:"xun_kong,omitempty"`
	YuePo          bool   `json:"yue_po,omitempty"`
	MuKu           bool   `json:"mu_ku,omitempty"`
	ExitCondition  string `json:"exit_condition,omitempty"`
}

type BranchRelationView struct {
	LeftPosition  int      `json:"left_position"`
	RightPosition int      `json:"right_position"`
	Relations     []string `json:"relations"`
	Targets       []string `json:"targets,omitempty"`
	Basis         []string `json:"basis,omitempty"`
}

type TombView struct {
	EarthBranch string `json:"earth_branch"`
	School      string `json:"school"`
}

type DayClashView struct {
	Position  int      `json:"position"`
	Kind      string   `json:"kind"`
	Basis     []string `json:"basis,omitempty"`
	Condition string   `json:"condition,omitempty"`
}

type MovingTransformationView struct {
	Position       int      `json:"position"`
	FromBranch     string   `json:"from_branch"`
	ToBranch       string   `json:"to_branch"`
	ReturnRelation string   `json:"return_relation,omitempty"`
	AdvanceRetreat string   `json:"advance_retreat,omitempty"`
	ChangedVoid    bool     `json:"changed_void,omitempty"`
	ChangedBreak   bool     `json:"changed_break,omitempty"`
	Basis          []string `json:"basis,omitempty"`
}

type SanHeView struct {
	Positions       []int    `json:"positions"`
	Branches        []string `json:"branches"`
	Element         string   `json:"wuxing"`
	Complete        bool     `json:"complete"`
	MovingCount     int      `json:"moving_count"`
	Qualification   string   `json:"qualification"`
	LineComplete    bool     `json:"line_complete"`
	DayPresent      bool     `json:"day_present"`
	MonthPresent    bool     `json:"month_present"`
	Missing         []string `json:"missing,omitempty"`
	StaticCount     int      `json:"static_count"`
	VoidCount       int      `json:"void_count"`
	BreakCount      int      `json:"break_count"`
	Activated       bool     `json:"activated"`
	Targets         []string `json:"targets,omitempty"`
	ConclusionScope string   `json:"conclusion_scope"`
}

type ForceEntryView struct {
	Role      string   `json:"role"`
	Element   string   `json:"wuxing"`
	Positions []int    `json:"positions,omitempty"`
	Branches  []string `json:"branches,omitempty"`
	Relations []string `json:"relations,omitempty"`
	States    []string `json:"states,omitempty"`
	Note      string   `json:"note,omitempty"`
}

type ForceChainView struct {
	YongShen    string           `json:"yong_shen"`
	YongElement string           `json:"yong_element"`
	Position    int              `json:"position"`
	IsHidden    bool             `json:"is_hidden"`
	Entries     []ForceEntryView `json:"entries"`
}

type CandidateView struct {
	Position  int      `json:"position"`
	Branch    string   `json:"branch"`
	IsHidden  bool     `json:"is_hidden"`
	Selected  bool     `json:"selected"`
	WangShuai string   `json:"wang_shuai,omitempty"`
	Reason    string   `json:"reason"`
	Basis     []string `json:"basis,omitempty"`
}

type ConflictView struct {
	ID     string   `json:"id"`
	Reason string   `json:"reason"`
	Facts  []string `json:"facts,omitempty"`
}

type AnalysisView struct {
	TimingCandidates      []TimingView              `json:"timing_candidates,omitempty"`
	DongYaoRelations      []DongRelationView        `json:"dong_yao_relations,omitempty"`
	HiddenLines           []HiddenLineView          `json:"hidden_lines,omitempty"`
	BranchRelationFacts   []BranchRelationView      `json:"branch_relation_facts,omitempty"`
	TombSchool            TombView                  `json:"tomb_school"`
	DayClashFacts         []DayClashView            `json:"day_clash_facts,omitempty"`
	MovingTransformations []MovingTransformationView `json:"moving_transformations,omitempty"`
	SanHeCandidates       []SanHeView               `json:"san_he_candidates,omitempty"`
	ForceChain            *ForceChainView           `json:"yong_shen_force_chain,omitempty"`
	YongShenCandidates    []CandidateView           `json:"yong_shen_candidates,omitempty"`
	Conflicts             []ConflictView            `json:"conflicts,omitempty"`
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
	Patterns     []PatternView  `json:"patterns"`
	YongShen     YongShenView   `json:"yong_shen"`
	GuaCi        *GuaCiView     `json:"gua_ci,omitempty"`
	Analysis     AnalysisView   `json:"analysis"`
}

// encodedName returns the readable name of an enum via its JSON encoding.
func encodedName(value any) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	var name string
	if err := json.Unmarshal(raw, &name); err != nil {
		return ""
	}
	return name
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
		ShiYing: line.ShiYing, LiuShou: line.LiuShou.String(),
		WangShuai: wangShuai.String(), ChangShengYue: line.ChangShengYue,
		Flags: flags,
	}
	if len(rel.Relations) > 0 {
		v.RiChenRelations = map[string][]string{"relations": append([]string(nil), rel.Relations...)}
	}
	return v
}

func fuShenView(v *FuShen) *FuShenView {
	if v == nil {
		return nil
	}
	return &FuShenView{Position: v.Position, LiuQin: v.LiuQin.String(), Zhi: v.Zhi}
}

func yongShenView(v YongShenResult) YongShenView {
	view := YongShenView{
		Name: v.Name, Position: v.Position, IsHidden: v.IsHidden,
		WangShuai: v.WangShuai, YuePo: v.YuePo, XunKong: v.XunKong,
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
		GuaName: encodedName(c.BenGua),
		Gong: c.Palace, GongWuxing: c.PalaceWuxing.String(),
		Lines: make([]LineView, 0, 6), BianLines: make([]LineView, 0, 6),
		RiChenGan: ganzhi.GanName(c.RiGan), RiChenZhi: ganzhi.ZhiName(c.RiZhi),
		YueJianGan: ganzhi.GanName(c.YueGan), YueJianZhi: ganzhi.ZhiName(c.YueZhi),
		XunKong: []string{ganzhi.ZhiName(c.XunKong[0]), ganzhi.ZhiName(c.XunKong[1])},
		Patterns: make([]PatternView, 0, len(c.Patterns)),
		YongShen: yongShenView(c.YongShen),
	}
	if int(c.BianGua) != 0 {
		name := encodedName(c.BianGua)
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
			Name: encodedName(c.GuaCi.Name), Text: c.GuaCi.GuaCi,
			YaoCi: append([]string(nil), c.GuaCi.YaoCi[:]...),
		}
	}
	view.Analysis = c.Analysis()
	return view
}

// Analysis returns the domain-aggregated analysis view: the deterministic
// judgment facts grouped into one section.
func (c Chart) Analysis() AnalysisView {
	out := AnalysisView{
		TombSchool: TombView{EarthBranch: c.TombSchool.EarthBranch, School: c.TombSchool.School},
	}
	for _, v := range c.TimingCandidates {
		out.TimingCandidates = append(out.TimingCandidates, TimingView{
			ID: v.ID, Mechanism: v.Mechanism, Position: v.Position,
			TriggerBranch: v.TriggerBranch, Window: v.Window,
			Basis: append([]string(nil), v.Basis...),
		})
	}
	for _, v := range c.DongYaoRelations {
		rels := make([]string, 0, len(v.Relations))
		for _, r := range v.Relations {
			rels = append(rels, string(r))
		}
		out.DongYaoRelations = append(out.DongYaoRelations, DongRelationView{Position: v.Position, Relations: rels})
	}
	for _, v := range c.HiddenLines {
		out.HiddenLines = append(out.HiddenLines, HiddenLineView{
			Position: v.Position, LiuQin: v.LiuQin,
			Gan: v.Gan, Zhi: v.Zhi, Wuxing: v.Wuxing,
			FlyingPosition: v.FlyingPosition, FlyingLiuQin: v.FlyingLiuQin,
			FlyingGan: v.FlyingGan, FlyingZhi: v.FlyingZhi,
			FlyingWuxing: v.FlyingWuxing, FeiFuRelation: v.FeiFuRelation,
			XunKong: v.XunKong, YuePo: v.YuePo, MuKu: v.MuKu, ExitCondition: v.ExitCondition,
		})
	}
	for _, v := range c.BranchRelationFacts {
		out.BranchRelationFacts = append(out.BranchRelationFacts, BranchRelationView{
			LeftPosition: v.LeftPosition, RightPosition: v.RightPosition,
			Relations: append([]string(nil), v.Relations...),
			Targets: append([]string(nil), v.Targets...), Basis: append([]string(nil), v.Basis...),
		})
	}
	for _, v := range c.DayClashFacts {
		out.DayClashFacts = append(out.DayClashFacts, DayClashView{
			Position: v.Position, Kind: v.Kind,
			Basis: append([]string(nil), v.Basis...), Condition: v.Condition,
		})
	}
	for _, v := range c.MovingTransformations {
		out.MovingTransformations = append(out.MovingTransformations, MovingTransformationView{
			Position: v.Position, FromBranch: v.FromBranch, ToBranch: v.ToBranch,
			ReturnRelation: v.ReturnRelation, AdvanceRetreat: v.AdvanceRetreat,
			ChangedVoid: v.ChangedVoid, ChangedBreak: v.ChangedBreak,
			Basis: append([]string(nil), v.Basis...),
		})
	}
	for _, v := range c.SanHeCandidates {
		positions := make([]int, 0, len(v.Positions))
		for _, p := range v.Positions {
			positions = append(positions, p)
		}
		out.SanHeCandidates = append(out.SanHeCandidates, SanHeView{
			Positions: positions, Branches: append([]string(nil), v.Branches...),
			Element: v.Element, Complete: v.Complete, MovingCount: v.MovingCount,
			Qualification: v.Qualification, LineComplete: v.LineComplete,
			DayPresent: v.DayPresent, MonthPresent: v.MonthPresent,
			Missing: append([]string(nil), v.Missing...), StaticCount: v.StaticCount,
			VoidCount: v.VoidCount, BreakCount: v.BreakCount, Activated: v.Activated,
			Targets: append([]string(nil), v.Targets...), ConclusionScope: v.ConclusionScope,
		})
	}
	if c.ForceChain != nil {
		force := ForceChainView{
			YongShen: c.ForceChain.YongShen, YongElement: c.ForceChain.YongElement,
			Position: c.ForceChain.Position, IsHidden: c.ForceChain.IsHidden,
		}
		for _, e := range c.ForceChain.Entries {
			positions := make([]int, 0, len(e.Positions))
			for _, p := range e.Positions {
				positions = append(positions, p)
			}
			force.Entries = append(force.Entries, ForceEntryView{
				Role: string(e.Role), Element: e.Element, Positions: positions,
				Branches: append([]string(nil), e.Branches...),
				Relations: append([]string(nil), e.Relations...),
				States: append([]string(nil), e.States...), Note: e.Note,
			})
		}
		out.ForceChain = &force
	}
	for _, v := range c.YongShenCandidates {
		out.YongShenCandidates = append(out.YongShenCandidates, CandidateView{
			Position: v.Position, Branch: v.Branch, IsHidden: v.IsHidden,
			Selected: v.Selected, WangShuai: v.WangShuai, Reason: v.Reason,
			Basis: append([]string(nil), v.Basis...),
		})
	}
	for _, v := range c.Conflicts {
		out.Conflicts = append(out.Conflicts, ConflictView{
			ID: v.ID, Reason: v.Reason, Facts: append([]string(nil), v.Facts...),
		})
	}
	return out
}
