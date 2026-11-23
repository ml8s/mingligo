package liuyao

import "github.com/ml8s/mingligo/pkg/ganzhi"

// This file defines the domain-owned, domain-aggregated output view of a 六爻
// chart. Top-level keys are命理域: 本卦/六爻/日辰月建/用神/格局/伏神/应期/关系/
// 信号/方法口径. Every chart fact has exactly one owning domain.

// GanZhiView is a stem-branch pair used by 日辰/月建.
type GanZhiView struct {
	Gan string `json:"gan" jsonschema:"天干。"`
	Zhi string `json:"zhi" jsonschema:"地支。"`
}

// BenGuaView groups the本卦 identity facts.
type BenGuaView struct {
	GuaMing    string  `json:"gua_ming" jsonschema:"本卦卦名。"`
	BianGua    *string `json:"bian_gua,omitempty" jsonschema:"变卦。"`
	Gong       string  `json:"gong" jsonschema:"宫位。"`
	GongWuxing string  `json:"gong_wuxing" jsonschema:"宫位五行。"`
}

// RiChenYueJianView groups 日辰/月建/旬空 — the time-axis facts of the cast.
type RiChenYueJianView struct {
	RiChen  GanZhiView `json:"ri_chen" jsonschema:"日辰。"`
	YueJian GanZhiView `json:"yue_jian" jsonschema:"月建。"`
	XunKong []string   `json:"xun_kong" jsonschema:"旬空地支（两字）。"`
}

// LineXiangView is the状态 of one line: 动/破/空/墓/生克.
type LineXiangView struct {
	Dong        bool     `json:"dong" jsonschema:"是否发动。"`
	YuePo       bool     `json:"yue_po" jsonschema:"是否月破。"`
	XunKong     bool     `json:"xun_kong" jsonschema:"是否旬空。"`
	MuKu        bool     `json:"mu_ku" jsonschema:"是否墓库。"`
	MuKuZhi     string   `json:"mu_ku_zhi,omitempty" jsonschema:"墓库地支。"`
	MuKuWuxing  string   `json:"mu_ku_wuxing,omitempty" jsonschema:"墓库五行。"`
	MuKuLaiYuan []string `json:"mu_ku_lai_yuan,omitempty" jsonschema:"入墓来源：日墓/动墓/化墓。"`
	DongSheng   bool     `json:"dong_sheng" jsonschema:"动爻生。"`
	DongKe      bool     `json:"dong_ke" jsonschema:"动爻克。"`
}

// LineView is one whole line: 干支/六亲/旺衰/状态 together.
type LineView struct {
	YaoWei        int           `json:"yao_wei" jsonschema:"爻位。"`
	YaoZhi        int           `json:"yao_zhi" jsonschema:"爻值。"`
	GanZhi        string        `json:"gan_zhi" jsonschema:"干支。"`
	Wuxing        string        `json:"wuxing" jsonschema:"五行。"`
	LiuQin        string        `json:"liu_qin" jsonschema:"六亲。"`
	ShiYing       string        `json:"shi_ying,omitempty" jsonschema:"世应。"`
	LiuShou       string        `json:"liu_shou" jsonschema:"六神。"`
	RiChenGuanXi  []string      `json:"ri_chen_guan_xi,omitempty" jsonschema:"日辰关系。"`
	WangShuai     string        `json:"wang_shuai,omitempty" jsonschema:"旺衰。"`
	ChangShengYue string        `json:"chang_sheng_yue,omitempty" jsonschema:"月令长生。"`
	Xiang         LineXiangView `json:"xiang" jsonschema:"爻状态。"`
}

// GeJuView groups 格局与三合结构 (格局域, no "_ceng" engineering layer suffix).
type GeJuView struct {
	Patterns        []PatternView        `json:"ge_ju_ming_xi,omitempty" jsonschema:"六爻格局。"`
	SanHeCandidates []SanHeCandidateView `json:"san_he_candidates,omitempty" jsonschema:"三合结构候选。"`
}

type PatternView struct {
	GeJuKind  string `json:"ge_ju_kind" jsonschema:"格局类型。"`
	ZiLei     string `json:"zi_lei" jsonschema:"格局子类。"`
	YaoWei    int    `json:"yao_wei" jsonschema:"爻位。"`
	ShiZhenGe bool   `json:"shi_zhen_ge" jsonschema:"是否真格局。"`
}

type FuShenItemView struct {
	YaoWei int    `json:"yao_wei" jsonschema:"爻位。"`
	LiuQin string `json:"liu_qin" jsonschema:"六亲。"`
	Zhi    string `json:"zhi" jsonschema:"地支。"`
}

// ForceChainView keeps only facts not already published by YongShenView.
// 用神名/爻位/伏藏 are the single responsibility of yong_shen itself.
type ForceChainView struct {
	YongWuxing string       `json:"yong_wuxing" jsonschema:"用神五行。"`
	Entries    []ForceEntry `json:"entries" jsonschema:"生扶克泄链条。"`
}

type YongShenCandidateView struct {
	YaoWei    int      `json:"yao_wei" jsonschema:"爻位。"`
	Zhi       string   `json:"zhi" jsonschema:"地支。"`
	IsHidden  bool     `json:"is_hidden" jsonschema:"是否伏藏。"`
	XuanZhong bool     `json:"xuan_zhong" jsonschema:"是否为当前取用。"`
	WangShuai string   `json:"wang_shuai,omitempty" jsonschema:"旺衰。"`
	LiYou     string   `json:"li_you" jsonschema:"取舍理由。"`
	YiJu      []string `json:"yi_ju,omitempty" jsonschema:"取舍依据。"`
}

// YongShenView is the single owner of every用神 fact: selected state, force
// chain and auditable candidates all live under one命理域.
type YongShenView struct {
	GuaMing          string                  `json:"gua_ming" jsonschema:"用神名称。"`
	YaoWei           int                     `json:"yao_wei" jsonschema:"爻位。"`
	IsHidden         bool                    `json:"is_hidden,omitempty" jsonschema:"是否伏藏。"`
	WangShuai        string                  `json:"wang_shuai" jsonschema:"旺衰。"`
	YuePo            bool                    `json:"yue_po,omitempty" jsonschema:"是否月破。"`
	XunKong          bool                    `json:"xun_kong,omitempty" jsonschema:"是否旬空。"`
	MuKu             bool                    `json:"mu_ku,omitempty" jsonschema:"是否墓库。"`
	LiuShou          string                  `json:"liu_shou,omitempty" jsonschema:"六神。"`
	ChangSheng       string                  `json:"chang_sheng,omitempty" jsonschema:"长生。"`
	FuShen           *FuShenItemView         `json:"fu_shen,omitempty" jsonschema:"所选用神的伏神。"`
	FuShenCandidates []FuShenItemView        `json:"fu_shen_candidates,omitempty" jsonschema:"伏神候选。"`
	ForceChain       *ForceChainView         `json:"liang_chain,omitempty" jsonschema:"用神力量链。"`
	Candidates       []YongShenCandidateView `json:"candidates,omitempty" jsonschema:"用神候选。"`
}

// FuShenView groups 伏神事实 (本宫全量伏神表).
type FuShenView struct {
	FuShenBiao []HiddenLine `json:"fu_shen_biao,omitempty" jsonschema:"本宫全量伏神表。"`
}

type ShiQiView struct {
	Candidates []TimingCandidate `json:"candidates,omitempty" jsonschema:"应期候选（非定论）。"`
}

type DongYaoGuanXiView struct {
	Relations       []DongYaoRelation      `json:"guan_xi,omitempty" jsonschema:"动爻关系。"`
	Transformations []MovingTransformation `json:"transformations,omitempty" jsonschema:"动变。"`
}

type ZhiGuanXiView struct {
	Relations  []BranchRelationFact `json:"guan_xi,omitempty" jsonschema:"地支关系。"`
	DayClashes []DayClashFact       `json:"day_clashes,omitempty" jsonschema:"日冲事实。"`
}

// MaoDunXinHaoView groups 矛盾信号 (待解释层裁决).
type MaoDunXinHaoView struct {
	Signals []Conflict `json:"signals,omitempty" jsonschema:"矛盾信号（非吉凶裁决）。"`
}

// FangFaKouJingView isolates method metadata from chart facts.
type FangFaKouJingView struct {
	MuKuLiuPai TombSchool `json:"mu_ku_liu_pai" jsonschema:"墓库流派（方法口径）。"`
}

// GuaCiView groups 卦辞/爻辞.
type GuaCiView struct {
	GuaMing string   `json:"gua_ming" jsonschema:"卦名。"`
	GuaCi   string   `json:"gua_ci" jsonschema:"卦辞文本。"`
	YaoCi   []string `json:"yao_ci" jsonschema:"爻辞。"`
}

// ChartView is the domain-aggregated, readable chart view. Every direct child
// is one命理域 and every fact appears exactly once.
type ChartView struct {
	BenGua        BenGuaView        `json:"ben_gua" jsonschema:"本卦。"`
	LiuYao        []LineView        `json:"liu_yao" jsonschema:"本卦六爻。"`
	BianYao       []LineView        `json:"bian_yao" jsonschema:"变卦六爻（仅动变事实）。"`
	RiChenYueJian RiChenYueJianView `json:"ri_chen_yue_jian" jsonschema:"日辰月建与旬空。"`
	YongShen      YongShenView      `json:"yong_shen" jsonschema:"用神（唯一事实源）。"`
	GeJu          GeJuView          `json:"ge_ju" jsonschema:"格局结构。"`
	FuShen        FuShenView        `json:"fu_shen" jsonschema:"伏神。"`
	ShiQi         ShiQiView         `json:"shi_qi" jsonschema:"应期。"`
	DongYaoGuanXi DongYaoGuanXiView `json:"dong_yao_guan_xi" jsonschema:"动爻关系。"`
	ZhiGuanXi     ZhiGuanXiView     `json:"zhi_guan_xi" jsonschema:"地支关系。"`
	MaoDunXinHao  MaoDunXinHaoView  `json:"mao_dun_xin_hao" jsonschema:"矛盾信号。"`
	GuaCi         *GuaCiView        `json:"gua_ci,omitempty" jsonschema:"卦辞。"`
	FangFaKouJing FangFaKouJingView `json:"fang_fa_kou_jing" jsonschema:"方法口径。"`
}

func (c Chart) lineView(line Line, rel DayRelation, wangShuai ganzhi.WangShuai) LineView {
	xiang := LineXiangView{
		Dong: line.DongSelf, YuePo: line.YuePo, XunKong: line.XunKong,
		MuKu: line.MuKu, MuKuZhi: line.MuKuBranch, MuKuWuxing: line.MuKuElement,
		DongSheng: line.DongSheng, DongKe: line.DongKe,
	}
	for _, source := range line.MuKuTypes {
		xiang.MuKuLaiYuan = append(xiang.MuKuLaiYuan, DecodeMuKuSource(source))
	}
	v := LineView{
		YaoWei: line.Position, YaoZhi: int(line.Type),
		GanZhi: ganzhi.GanName(line.Gan) + ganzhi.ZhiName(line.Zhi),
		Wuxing: ganzhi.DecodeWuxing(line.Wuxing), LiuQin: DecodeLiuQin(line.LiuQin),
		ShiYing: DecodeShiYing(line.ShiYing), LiuShou: DecodeLiuShou(line.LiuShou),
		WangShuai: ganzhi.DecodeWangShuai(wangShuai), ChangShengYue: string(line.ChangShengYue),
		Xiang: xiang,
	}
	for _, item := range rel.Relations {
		v.RiChenGuanXi = append(v.RiChenGuanXi, DecodeDayRelationType(item))
	}
	return v
}

func fuShenItemView(v *FuShen) *FuShenItemView {
	if v == nil {
		return nil
	}
	return &FuShenItemView{YaoWei: v.Position, LiuQin: DecodeLiuQin(v.LiuQin), Zhi: ganzhi.ZhiName(v.Zhi)}
}

func forceChainView(v *ForceChain) *ForceChainView {
	if v == nil {
		return nil
	}
	return &ForceChainView{YongWuxing: v.YongElement, Entries: v.Entries}
}

func yongShenView(v YongShenResult, c Chart) YongShenView {
	view := YongShenView{
		GuaMing: v.Name, YaoWei: v.Position, IsHidden: v.IsHidden,
		WangShuai: ganzhi.DecodeWangShuai(v.WangShuai), YuePo: v.YuePo, XunKong: v.XunKong,
		MuKu: v.MuKu, LiuShou: DecodeLiuShou(v.LiuShou), ChangSheng: string(v.ChangSheng),
		FuShen:     fuShenItemView(v.FuShen),
		ForceChain: forceChainView(c.ForceChain),
	}
	for _, cand := range v.FuShenCandidates {
		if item := fuShenItemView(&cand); item != nil {
			view.FuShenCandidates = append(view.FuShenCandidates, *item)
		}
	}
	for _, cand := range c.YongShenCandidates {
		view.Candidates = append(view.Candidates, YongShenCandidateView{
			YaoWei: cand.Position, Zhi: cand.Branch, IsHidden: cand.IsHidden,
			XuanZhong: cand.Selected, WangShuai: cand.WangShuai,
			LiYou: cand.Reason, YiJu: cand.Basis,
		})
	}
	return view
}

// View returns the domain-aggregated, readable chart view.
func (c Chart) View() ChartView {
	view := ChartView{
		BenGua: BenGuaView{
			GuaMing: c.BenGua.Name(), Gong: c.Palace, GongWuxing: ganzhi.DecodeWuxing(c.PalaceWuxing),
		},
		LiuYao: make([]LineView, 0, 6), BianYao: make([]LineView, 0, 6),
		RiChenYueJian: RiChenYueJianView{
			RiChen:  GanZhiView{Gan: ganzhi.GanName(c.RiGan), Zhi: ganzhi.ZhiName(c.RiZhi)},
			YueJian: GanZhiView{Gan: ganzhi.GanName(c.YueGan), Zhi: ganzhi.ZhiName(c.YueZhi)},
			XunKong: []string{ganzhi.ZhiName(c.XunKong[0]), ganzhi.ZhiName(c.XunKong[1])},
		},
		YongShen: yongShenView(c.YongShen, c),
		GeJu: GeJuView{
			Patterns: make([]PatternView, 0, len(c.Patterns)), SanHeCandidates: sanHeCandidateViews(c.SanHeCandidates),
		},
		FuShen:        FuShenView{FuShenBiao: c.HiddenLines},
		ShiQi:         ShiQiView{Candidates: c.TimingCandidates},
		DongYaoGuanXi: DongYaoGuanXiView{Relations: c.DongYaoRelations, Transformations: c.MovingTransformations},
		ZhiGuanXi:     ZhiGuanXiView{Relations: c.BranchRelationFacts, DayClashes: c.DayClashFacts},
		MaoDunXinHao:  MaoDunXinHaoView{Signals: c.Conflicts},
		FangFaKouJing: FangFaKouJingView{MuKuLiuPai: c.TombSchool},
	}
	if int(c.BianGua) != 0 {
		name := c.BianGua.Name()
		view.BenGua.BianGua = &name
	}
	for i, line := range c.Lines {
		view.LiuYao = append(view.LiuYao, c.lineView(line, c.DayRelations[i], c.WangShuai[i]))
	}
	for _, line := range c.BianLines {
		if line.Position != 0 {
			view.BianYao = append(view.BianYao, c.lineView(line, DayRelation{}, ganzhi.WangShuai(-1)))
		}
	}
	for _, p := range c.Patterns {
		view.GeJu.Patterns = append(view.GeJu.Patterns, PatternView{
			GeJuKind: string(p.Type), ZiLei: string(p.SubType),
			YaoWei: p.Position, ShiZhenGe: p.IsTrue,
		})
	}
	if c.GuaCi.Name != "" {
		view.GuaCi = &GuaCiView{
			GuaMing: c.GuaCi.Name, GuaCi: c.GuaCi.GuaCi,
			YaoCi: append([]string(nil), c.GuaCi.YaoCi[:]...),
		}
	}
	return view
}

// SanHeCandidateView maps the conclusion-scope code at the View boundary.
type SanHeCandidateView struct {
	Positions       []int        `json:"positions"`
	Branches        []string     `json:"zhi"`
	Element         string       `json:"wuxing"`
	Complete        bool         `json:"complete"`
	MovingCount     int          `json:"moving_count"`
	Qualification   string       `json:"qualification"`
	LineComplete    bool         `json:"line_complete"`
	DayPresent      bool         `json:"day_present"`
	MonthPresent    bool         `json:"month_present"`
	Missing         []string     `json:"missing,omitempty"`
	StaticCount     int          `json:"static_count"`
	VoidCount       int          `json:"void_count"`
	BreakCount      int          `json:"break_count"`
	Activated       bool         `json:"activated"`
	Targets         []TargetRole `json:"targets,omitempty"`
	ConclusionScope string       `json:"conclusion_scope"`
}

func sanHeCandidateViews(src []SanHeCandidate) []SanHeCandidateView {
	if src == nil {
		return nil
	}
	out := make([]SanHeCandidateView, 0, len(src))
	for _, c := range src {
		scope := c.ConclusionScope
		if scope == "structure_candidate_only_not_outcome" {
			scope = "仅结构候选（非吉凶结论）"
		}
		out = append(out, SanHeCandidateView{
			Positions: c.Positions, Branches: c.Branches, Element: c.Element,
			Complete: c.Complete, MovingCount: c.MovingCount, Qualification: c.Qualification,
			LineComplete: c.LineComplete, DayPresent: c.DayPresent, MonthPresent: c.MonthPresent,
			Missing: c.Missing, StaticCount: c.StaticCount, VoidCount: c.VoidCount,
			BreakCount: c.BreakCount, Activated: c.Activated, Targets: c.Targets,
			ConclusionScope: scope,
		})
	}
	return out
}
