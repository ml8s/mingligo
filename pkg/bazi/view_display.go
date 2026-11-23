package bazi

import "github.com/ml8s/mingligo/pkg/ganzhi"

// Display mappings at the View boundary only. Internal codes, rule tables,
// comparisons, map keys and MCP input enums are unchanged — these are pure
// presentation lookups for machine-code strings that reach the wire.

var tenGodSourceDisplay = map[string]string{
	"gan": "透干", "main_qi": "本气", "mid_qi": "中气", "minor_qi": "余气",
}

var rootKindDisplay = map[string]string{
	"any": "任一", "none": "无根", "mid": "中气",
	"month_main": "月令本气", "month_mid": "月令中气",
	"zhi_main": "地支本气", "zhi_mid": "地支中气", "minor": "余气",
}

var relationPositionDisplay = map[string]string{
	"adjacent": "相邻", "remote": "隔位", "separated": "远离",
}

var strengthDisplay = map[string]string{
	"strong": "强", "weak": "弱", "neutral": "中和",
}

var reasonDisplay = map[string]string{
	"timely": "得令", "untimely": "失令", "transparent_rooted": "透干通根",
	"not_transparent": "不透", "not_rooted": "无根",
}

var targetDisplay = map[string]string{
	"day_master": "日主", "yin_bi": "印比", "day_branch": "日支",
	"day_master_root": "日主根", "primary": "主", "secondary": "辅",
	"pattern": "格局", "controller": "制化", "generator": "生扶",
}

var patternSourceDisplay = map[string]string{
	"main_qi": "本气", "mid_qi": "中气", "minor_qi": "余气",
	"main_qi_not_transparent": "本气不透", "month_blade": "月刃", "month_lu": "月禄",
}

var policyDisplay = map[string]string{
	"union_of_year_and_day_references": "年日双参照并取",
}

var xiaoYunLiuPaiDisplay = map[string]string{
	"san_ming_tong_hui": "三命通会", "xing_ping_hui_hai": "星平会海",
}

var dayunRootSourceDisplay = map[string]string{
	"natal_hidden_stem": "命局藏干", "sitting_branch_main_qi": "坐支本气",
}

func displayFromBazi(m map[string]string, code string) string {
	if v, ok := m[code]; ok {
		return v
	}
	return code
}

func displayStringList(m map[string]string, codes []string) []string {
	if codes == nil {
		return nil
	}
	out := make([]string, 0, len(codes))
	for _, code := range codes {
		out = append(out, displayFromBazi(m, code))
	}
	return out
}

// View-boundary types for domain structs carrying machine codes. Same json
// keys; only code-valued fields are mapped to 命理 display terms.
type TenGodProjectionView struct {
	ShiShen ganzhi.ShiShen `json:"shi_shen"`
	Name    string         `json:"name"`
	Source  string         `json:"source"`
	Gan     ganzhi.Gan     `json:"gan"`
}

type PillarView struct {
	ganzhi.Zhu
	DayMasterTrend string                 `json:"ri_zhu_trend"`
	Xun            string                 `json:"xun"`
	XunKong        string                 `json:"xun_kong"`
	NaYin          string                 `json:"na_yin"`
	CangGan        HiddenStemProjection   `json:"cang_gan"`
	ShiShens       []TenGodProjectionView `json:"shi_shens"`
	SelfSitting    string                 `json:"self_sitting"`
	ShenSha        []ShenShaProjection    `json:"shen_sha"`
	IsVoid         bool                   `json:"is_void"`
	IsSelfHe       bool                   `json:"is_self_he"`
	IsKuiGang      bool                   `json:"is_kui_gang"`
	SelfHeName     string                 `json:"self_he_name"`
}

func pillarView(p Pillar) PillarView {
	v := PillarView{
		Zhu: p.Zhu, DayMasterTrend: p.DayMasterTrend, Xun: p.Xun, XunKong: p.XunKong,
		NaYin: p.NaYin, CangGan: p.CangGan, SelfSitting: p.SelfSitting,
		ShenSha: p.ShenSha, IsVoid: p.IsVoid, IsSelfHe: p.IsSelfHe,
		IsKuiGang: p.IsKuiGang, SelfHeName: p.SelfHeName,
	}
	if p.ShiShens != nil {
		v.ShiShens = make([]TenGodProjectionView, 0, len(p.ShiShens))
		for _, s := range p.ShiShens {
			v.ShiShens = append(v.ShiShens, TenGodProjectionView{
				ShiShen: s.ShiShen, Name: s.Name,
				Source: displayFromBazi(tenGodSourceDisplay, string(s.Source)), Gan: s.Gan,
			})
		}
	}
	return v
}

type TenGodStateView struct {
	ShiShen     string `json:"shi_shen"`
	Wuxing      string `json:"wuxing"`
	Transparent bool   `json:"transparent"`
	Hidden      bool   `json:"hidden"`
	Rooted      bool   `json:"rooted"`
	Timely      bool   `json:"timely"`
	Count       int    `json:"count"`
	Strength    string `json:"strength"`
}

type ElementStateView struct {
	Wuxing             string   `json:"wuxing"`
	Season             string   `json:"season"`
	SeasonStrength     string   `json:"season_strength"`
	Strength           string   `json:"strength"`
	Transparent        bool     `json:"transparent"`
	Rooted             bool     `json:"rooted"`
	Controls           string   `json:"controls"`
	ControlledBy       string   `json:"controlled_by"`
	ControllerStrength string   `json:"controller_strength"`
	Reasons            []string `json:"reasons"`
}

func tenGodStateViews(src []TenGodState) []TenGodStateView {
	if src == nil {
		return nil
	}
	out := make([]TenGodStateView, 0, len(src))
	for _, s := range src {
		out = append(out, TenGodStateView{
			ShiShen: s.ShiShen, Wuxing: s.Wuxing, Transparent: s.Transparent,
			Hidden: s.Hidden, Rooted: s.Rooted, Timely: s.Timely, Count: s.Count,
			Strength: displayFromBazi(strengthDisplay, s.Strength),
		})
	}
	return out
}

func elementStateViews(src []ElementState) []ElementStateView {
	if src == nil {
		return nil
	}
	out := make([]ElementStateView, 0, len(src))
	for _, s := range src {
		out = append(out, ElementStateView{
			Wuxing: s.Wuxing, Season: s.Season,
			SeasonStrength: displayFromBazi(strengthDisplay, s.SeasonStrength),
			Strength:       displayFromBazi(strengthDisplay, s.Strength),
			Transparent:    s.Transparent, Rooted: s.Rooted, Controls: s.Controls,
			ControlledBy:       s.ControlledBy,
			ControllerStrength: displayFromBazi(strengthDisplay, s.ControllerStrength),
			Reasons:            displayStringList(reasonDisplay, s.Reasons),
		})
	}
	return out
}

type NatalRelationFactView struct {
	ID        string                `json:"id"`
	Kind      NatalRelationKind     `json:"relation_kind"`
	Group     string                `json:"group"`
	Members   []NatalRelationMember `json:"members"`
	Wuxing    string                `json:"wuxing,omitempty"`
	Position  string                `json:"position,omitempty"`
	Contested bool                  `json:"contested,omitempty"`
}

func natalRelationFactViews(src []NatalRelationFact) []NatalRelationFactView {
	if src == nil {
		return nil
	}
	out := make([]NatalRelationFactView, 0, len(src))
	for _, r := range src {
		out = append(out, NatalRelationFactView{
			ID: r.ID, Kind: r.Kind, Group: r.Group, Members: r.Members,
			Wuxing: r.Wuxing, Position: displayFromBazi(relationPositionDisplay, r.Position),
			Contested: r.Contested,
		})
	}
	return out
}

type RelationLinkView struct {
	ID      string   `json:"id"`
	Targets []string `json:"targets"`
}

type FuYiBasisView struct {
	RiGanGen      string               `json:"ri_gan_gen"`
	YinBiShuLiang int                  `json:"yin_bi_shu_liang"`
	RiGanTongGen  []StemOccurrenceView `json:"ri_gan_tong_gen"`
	GuanXiLinks   []RelationLinkView   `json:"guan_xi_links"`
}

func stemOccurrenceViews(src []StemOccurrence) []StemOccurrenceView {
	if src == nil {
		return nil
	}
	out := make([]StemOccurrenceView, 0, len(src))
	for _, s := range src {
		zhu := DecodePillarIndex(s.Pillar)
		zhi := ganzhi.ZhiName(s.Branch)
		gan := ganzhi.GanName(s.Stem)
		src := displayFromBazi(tenGodSourceDisplay, string(s.Source))
		out = append(out, StemOccurrenceView{
			Pillar: &zhu, Branch: &zhi, Stem: &gan, Source: &src,
		})
	}
	return out
}

func relationLinkViews(src []RelationLink) []RelationLinkView {
	if src == nil {
		return nil
	}
	out := make([]RelationLinkView, 0, len(src))
	for _, r := range src {
		out = append(out, RelationLinkView{ID: r.ID, Targets: displayStringList(targetDisplay, r.Targets)})
	}
	return out
}

func fuYiBasisView(b FuYiBasis) FuYiBasisView {
	return FuYiBasisView{
		RiGanGen:      displayFromBazi(rootKindDisplay, b.RiGanGen),
		YinBiShuLiang: b.YinBiShuLiang,
		RiGanTongGen:  stemOccurrenceViews(b.RiGanTongGen),
		GuanXiLinks:   relationLinkViews(b.GuanXiLinks),
	}
}

type ElementStructureView struct {
	Wuxing             ganzhi.Wuxing        `json:"wuxing"`
	TransparentPillars []PillarIndex        `json:"transparent_pillars"`
	Roots              []StemOccurrenceView `json:"roots"`
	Season             ganzhi.WangShuai     `json:"season"`
	Timely             bool                 `json:"timely"`
	Strength           string               `json:"strength"`
}

type GeJuStructureView struct {
	Pattern     ElementStructureView `json:"ge_ju_shen"`
	Controller  ElementStructureView `json:"controller"`
	Generator   ElementStructureView `json:"generator"`
	GuanXiLinks []RelationLinkView   `json:"guan_xi_links"`
}

type GeJuResultView struct {
	Pattern          GeJuKind          `json:"ge_ju"`
	GeShenTou        bool              `json:"ge_shen_tou"`
	Usage            string            `json:"usage"`
	PatternGod       string            `json:"ge_ju_shen,omitempty"`
	PatternGodTenGod string            `json:"ge_ju_shen_shi_shen,omitempty"`
	PatternGodSource string            `json:"ge_ju_shen_source"`
	Structure        GeJuStructureView `json:"structure"`
}

func elementStructureView(e ElementStructure) ElementStructureView {
	roots := make([]StemOccurrenceView, 0, len(e.Roots))
	for _, r := range e.Roots {
		zhu := DecodePillarIndex(r.Pillar)
		zhi := ganzhi.ZhiName(r.Branch)
		gan := ganzhi.GanName(r.Stem)
		src := displayFromBazi(tenGodSourceDisplay, string(r.Source))
		roots = append(roots, StemOccurrenceView{Pillar: &zhu, Branch: &zhi, Stem: &gan, Source: &src})
	}
	return ElementStructureView{
		Wuxing: e.Wuxing, TransparentPillars: e.TransparentPillars, Roots: roots,
		Season: e.Season, Timely: e.Timely, Strength: displayFromBazi(strengthDisplay, e.Strength),
	}
}

func geJuResultView(g GeJuResult) GeJuResultView {
	return GeJuResultView{
		Pattern: g.Pattern, GeShenTou: g.GeShenTou, Usage: g.Usage,
		PatternGod: g.PatternGod, PatternGodTenGod: g.PatternGodTenGod,
		PatternGodSource: displayFromBazi(patternSourceDisplay, g.PatternGodSource),
		Structure: GeJuStructureView{
			Pattern:     elementStructureView(g.Structure.Pattern),
			Controller:  elementStructureView(g.Structure.Controller),
			Generator:   elementStructureView(g.Structure.Generator),
			GuanXiLinks: relationLinkViews(g.Structure.GuanXiLinks),
		},
	}
}

type ShenShaSchoolView struct {
	DualReference []string `json:"dual_reference"`
	Policy        string   `json:"policy"`
}

func shenShaSchoolView(s ShenShaSchool) ShenShaSchoolView {
	return ShenShaSchoolView{
		DualReference: s.DualReference,
		Policy:        displayFromBazi(policyDisplay, s.Policy),
	}
}

type DaYunRootRefView struct {
	Source string `json:"source"`
	Pillar string `json:"zhu,omitempty"`
	Stem   string `json:"gan"`
	Branch string `json:"zhi"`
}

type DaYunStepView struct {
	Gan       ganzhi.Gan         `json:"gan"`
	Zhi       ganzhi.Zhi         `json:"zhi"`
	StartDate string             `json:"start_date"`
	EndDate   string             `json:"end_date"`
	StartYear int                `json:"start_year"`
	EndYear   int                `json:"end_year"`
	Name      string             `json:"name"`
	Element   string             `json:"wuxing"`
	ShiShen   string             `json:"shi_shen"`
	Rooted    *bool              `json:"rooted,omitempty"`
	RootRefs  []DaYunRootRefView `json:"root_refs,omitempty"`
}

type DaYunView struct {
	StartDate       string          `json:"start_date"`
	StartYearAfter  int             `json:"start_year_after"`
	StartMonthAfter int             `json:"start_month_after"`
	StartDayAfter   int             `json:"start_day_after"`
	Direction       string          `json:"direction"`
	Steps           []DaYunStepView `json:"steps"`
}

func daYunView(d *DaYun) *DaYunView {
	if d == nil {
		return nil
	}
	v := &DaYunView{
		StartDate: d.StartDate, StartYearAfter: d.StartYearAfter,
		StartMonthAfter: d.StartMonthAfter, StartDayAfter: d.StartDayAfter,
		Direction: d.Direction, Steps: make([]DaYunStepView, 0, len(d.Steps)),
	}
	for _, s := range d.Steps {
		sv := DaYunStepView{
			Gan: s.Gan, Zhi: s.Zhi, StartDate: s.StartDate, EndDate: s.EndDate,
			StartYear: s.StartYear, EndYear: s.EndYear, Name: s.Name,
			Element: s.Element, ShiShen: s.ShiShen, Rooted: s.Rooted,
		}
		if s.RootRefs != nil {
			sv.RootRefs = make([]DaYunRootRefView, 0, len(s.RootRefs))
			for _, r := range s.RootRefs {
				sv.RootRefs = append(sv.RootRefs, DaYunRootRefView{
					Source: displayFromBazi(dayunRootSourceDisplay, r.Source),
					Pillar: r.Pillar, Stem: r.Stem, Branch: r.Branch,
				})
			}
		}
		v.Steps = append(v.Steps, sv)
	}
	return v
}
