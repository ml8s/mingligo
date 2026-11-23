package qimen

import "github.com/ml8s/mingligo/pkg/ganzhi"

// This file defines display mappings for machine-code values at the View
// boundary only. Internal codes, rule tables, comparisons, map keys and MCP
// input enums are unchanged — these functions are pure presentation lookups.

var methodSourceDisplay = map[string]string{
	"solar_term":         "节气",
	"quarter_rule":       "刻家规则",
	"month_cycle":        "月建",
	"year_cycle":         "岁循环",
	"solstice_day_cycle": "二至日循环",
}

var dayBoundaryDisplay = map[string]string{
	"late_zi_rolls_day_pillar": "晚子时日柱进位",
	"zi_zheng":                 "子正换日",
}

var geometryDisplay = map[string]string{
	"luoshu_ring_rotation":            "洛书飞泊",
	"palace_number_shift":             "宫数位移",
	"mingfa_star_door_forward_flight": "鸣法星门顺飞",
	"solstice_day_star_flight":        "二至日星飞",
}

var starModeDisplay = map[string]string{
	"tian_qin_lodges_kun":            "天禽寄坤",
	"star_gan_shift_together":        "星干同移",
	"nine_stars_forward_with_center": "九星顺飞寄中",
	"taiyi_nine_stars_day_flight":    "太乙九星日飞",
}

var doorModeDisplay = map[string]string{
	"center_projects_kun":            "中宫寄坤",
	"eight_doors_real_center_step":   "八门实中宫步进",
	"nine_doors_with_center_forward": "九门寄中顺飞",
	"three_day_eight_door_rotation":  "三日一换八门",
}

var spiritModeDisplay = map[string]string{
	"eight_spirits_yang_gouchen_zhuque_yin_baihu_xuanwu": "八神（阳勾陈朱雀·阴白虎玄武）",
	"nine_spirits_baihu_taichang_xuanwu":                 "九神（白虎太常玄武）",
	"nine_spirits_gouchen_taichang_zhuque":               "九神（勾陈太常朱雀）",
	"day_gan_twelve_spirits":                             "日干十二神",
}

var spiritFlightDisplay = map[string]string{
	"yang_forward_yin_reverse":          "阳顺阴逆",
	"solstice_yang_forward_yin_reverse": "二至阳顺阴逆",
}

var jinhanDunSourceDisplay = map[string]string{
	"winter_solstice_yang_summer_solstice_yin": "冬至后阳遁·夏至后阴遁",
}

var wuxingRelationDisplay = map[string]string{
	"same": "比和", "gong_generates_xing": "宫生星", "xing_generates_gong": "星生宫",
	"xing_controls_gong": "星克宫", "gong_controls_xing": "宫克星",
	"generates": "相生", "controls": "相克", "generated_by": "受生", "controlled_by": "受克",
	"ri_generates_shi": "日生时", "shi_generates_ri": "时生日",
	"ri_controls_shi": "日克时", "shi_controls_ri": "时克日",
}

var layerDisplay = map[string]string{"heaven": "天盘", "earth": "地盘"}

var roleDisplay = map[string]string{
	"pillar": "柱", "lead": "值符", "duty": "值使", "yong_shen": "用神",
}

var pillarLevelDisplay = map[string]string{
	"year": "年", "month": "月", "day": "日", "hour": "时",
}

var yingQiKindDisplay = map[string]string{
	"ma_xing": "马星", "kong_wang": "空亡", "duty_star": "值符星", "duty_door": "值使门",
}

var yingQiMatchDisplay = map[string]string{"same": "同支", "opposite": "对支"}

var patternKindDisplay = map[string]string{
	"san_qi_ru_mu": "三奇入墓", "liu_yi_ji_xing": "六仪击刑", "you_hun_ru_mu": "游魂入墓",
}

var domainDisplay = map[string]string{"inner": "内盘", "outer": "外盘", "center": "中宫"}
var rangeDisplay = map[string]string{"low": "下", "high": "上"}

func displayFrom(m map[string]string, code string) string {
	if v, ok := m[code]; ok {
		return v
	}
	return code
}

// View-boundary types for domain structs that carry machine codes. Each keeps
// the same json keys; only code-valued fields are mapped to 命理 display terms.
type GongProfileView struct {
	Gong      GongIndex `json:"gong"`
	Name      string    `json:"name"`
	Direction string    `json:"direction"`
	Person    string    `json:"person,omitempty"`
	Range     string    `json:"range"`
	Domain    string    `json:"domain"`
}

type GongRelationView struct {
	From     GongIndex `json:"from"`
	To       GongIndex `json:"to"`
	Relation string    `json:"relation"`
}

type XingGongWuXingView struct {
	Star             StarIndex `json:"xing"`
	Gong             GongIndex `json:"gong"`
	Relation         string    `json:"relation"`
	RelationName     string    `json:"relation_name"`
	TraditionalLabel string    `json:"traditional_label"`
}

type GanGongFactView struct {
	Palace GongIndex `json:"gong"`
	Layer  string    `json:"layer"`
}

type RiShiRelationView struct {
	Zhu      string `json:"zhu"`
	Ke       string `json:"ke"`
	Relation string `json:"relation"`
	Name     string `json:"name"`
}

type AffectedSymbolView struct {
	Symbol string     `json:"symbol"`
	Role   string     `json:"role"`
	Branch ganzhi.Zhi `json:"zhi"`
	Gong   GongIndex  `json:"gong"`
}

type PatternView struct {
	Kind        string      `json:"kind,omitempty"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	YiJu        string      `json:"basis,omitempty"`
	Auspicious  bool        `json:"auspicious"`
	GongWei     []GongIndex `json:"gong_wei,omitempty"`
}

func gongProfileViews(src [9]GongProfile) [9]GongProfileView {
	var out [9]GongProfileView
	for i, p := range src {
		out[i] = GongProfileView{
			Gong: p.Gong, Name: p.Name, Direction: p.Direction, Person: p.Person,
			Range: displayFrom(rangeDisplay, p.Range), Domain: displayFrom(domainDisplay, p.Domain),
		}
	}
	return out
}

func gongRelationViews(src []GongRelation) []GongRelationView {
	if src == nil {
		return nil
	}
	out := make([]GongRelationView, 0, len(src))
	for _, r := range src {
		out = append(out, GongRelationView{
			From: r.From, To: r.To, Relation: displayFrom(wuxingRelationDisplay, string(r.Relation)),
		})
	}
	return out
}

func xingGongWuXingViews(src []XingGongWuXing) []XingGongWuXingView {
	if src == nil {
		return nil
	}
	out := make([]XingGongWuXingView, 0, len(src))
	for _, x := range src {
		out = append(out, XingGongWuXingView{
			Star: x.Star, Gong: x.Gong,
			Relation:     displayFrom(wuxingRelationDisplay, x.Relation),
			RelationName: x.RelationName, TraditionalLabel: x.TraditionalLabel,
		})
	}
	return out
}

func ganGongFactViews(src []GanGongFact) []GanGongFactView {
	if src == nil {
		return nil
	}
	out := make([]GanGongFactView, 0, len(src))
	for _, f := range src {
		out = append(out, GanGongFactView{Palace: f.Palace, Layer: displayFrom(layerDisplay, f.Layer)})
	}
	return out
}

func riShiRelationView(src RiShiRelation) RiShiRelationView {
	palaceDisplay := map[string]string{
		"day_stem_palace": "日干宫", "hour_stem_palace": "时干宫",
		"year_stem_palace": "年干宫", "month_stem_palace": "月干宫",
	}
	return RiShiRelationView{
		Zhu: displayFrom(palaceDisplay, src.Zhu), Ke: displayFrom(palaceDisplay, src.Ke),
		Relation: displayFrom(wuxingRelationDisplay, src.Relation), Name: src.Name,
	}
}

func affectedSymbolViews(src []AffectedSymbol) []AffectedSymbolView {
	if src == nil {
		return nil
	}
	out := make([]AffectedSymbolView, 0, len(src))
	for _, s := range src {
		out = append(out, AffectedSymbolView{
			Symbol: s.Symbol, Role: displayFrom(roleDisplay, s.Role), Branch: s.Branch, Gong: s.Gong,
		})
	}
	return out
}

func patternViews(src []Pattern) []PatternView {
	if src == nil {
		return nil
	}
	out := make([]PatternView, 0, len(src))
	for _, p := range src {
		out = append(out, PatternView{
			Kind: displayFrom(patternKindDisplay, p.Kind), Name: p.Name,
			Description: p.Description, YiJu: p.YiJu, Auspicious: p.Auspicious, GongWei: p.GongWei,
		})
	}
	return out
}

type StemPalaceView struct {
	Palace    GongIndex         `json:"gong"`
	Facts     []GanGongFactView `json:"facts"`
	WangShuai GongWangShuai     `json:"wang_shuai"`
}

type YingQiCandidateView struct {
	Type      string              `json:"ying_qi_kind"`
	Mechanism string              `json:"mechanism"`
	Branch    string              `json:"zhi"`
	Gong      GongIndex           `json:"gong"`
	RelatedTo []RelatedSymbolView `json:"related_to"`
	Dates     []YingQiDateView    `json:"dates"`
}

type RelatedSymbolView struct {
	Symbol string `json:"symbol"`
	Role   string `json:"role"`
}

type YingQiDateView struct {
	Date   string `json:"trigger_date"`
	Branch string `json:"zhi"`
	Match  string `json:"match"`
}

type YingQiView struct {
	Candidates []YingQiCandidateView `json:"candidates"`
}

type GengGeFactView struct {
	Level    string     `json:"level"`
	Gong     GongIndex  `json:"gong"`
	Door     string     `json:"men,omitempty"`
	EarthGan ganzhi.Gan `json:"earth_gan"`
}

type TianwangFactView struct {
	Kind           string     `json:"kind"`
	Gan            ganzhi.Gan `json:"gan"`
	Gong           GongIndex  `json:"gong"`
	HourGong       GongIndex  `json:"hour_gong"`
	RelationToHour string     `json:"relation_to_hour"`
}

type SpecializedContextView struct {
	GengGe      []GengGeFactView   `json:"geng_ge"`
	LostContext []LostContextFact  `json:"lost_facts"`
	Thief       []ThiefFact        `json:"thief_facts"`
	Tianwang    []TianwangFactView `json:"tianwang_facts"`
}

type YongShenResultView struct {
	NianGanPalace    *GongIndex        `json:"nian_gan_gong,omitempty"`
	NianGanGongFacts []GanGongFactView `json:"nian_gan_gong_facts,omitempty"`
	Symbols          []SymbolResult    `json:"symbols"`
}

func stemPalaceView(src StemPalace) StemPalaceView {
	return StemPalaceView{Palace: src.Palace, Facts: ganGongFactViews(src.Facts), WangShuai: src.WangShuai}
}

func yingQiView(src YingQi) YingQiView {
	if src.Candidates == nil {
		return YingQiView{}
	}
	out := make([]YingQiCandidateView, 0, len(src.Candidates))
	for _, c := range src.Candidates {
		rel := make([]RelatedSymbolView, 0, len(c.RelatedTo))
		for _, r := range c.RelatedTo {
			rel = append(rel, RelatedSymbolView{Symbol: r.Symbol, Role: displayFrom(roleDisplay, r.Role)})
		}
		dates := make([]YingQiDateView, 0, len(c.Dates))
		for _, d := range c.Dates {
			dates = append(dates, YingQiDateView{
				Date: d.Date, Branch: d.Branch, Match: displayFrom(yingQiMatchDisplay, string(d.Match)),
			})
		}
		out = append(out, YingQiCandidateView{
			Type: displayFrom(yingQiKindDisplay, string(c.Type)), Mechanism: c.Mechanism,
			Branch: c.Branch, Gong: c.Gong, RelatedTo: rel, Dates: dates,
		})
	}
	return YingQiView{Candidates: out}
}

func specializedContextView(src SpecializedContext) SpecializedContextView {
	out := SpecializedContextView{LostContext: src.LostContext, Thief: src.Thief}
	if src.GengGe != nil {
		out.GengGe = make([]GengGeFactView, 0, len(src.GengGe))
		for _, g := range src.GengGe {
			out.GengGe = append(out.GengGe, GengGeFactView{
				Level: displayFrom(pillarLevelDisplay, string(g.Level)),
				Gong:  g.Gong, Door: g.Door, EarthGan: g.EarthGan,
			})
		}
	}
	if src.Tianwang != nil {
		out.Tianwang = make([]TianwangFactView, 0, len(src.Tianwang))
		for _, t := range src.Tianwang {
			out.Tianwang = append(out.Tianwang, TianwangFactView{
				Kind: t.Kind, Gan: t.Gan, Gong: t.Gong, HourGong: t.HourGong,
				RelationToHour: displayFrom(wuxingRelationDisplay, string(t.RelationToHour)),
			})
		}
	}
	return out
}

func yongShenResultView(src *YongShenResult) *YongShenResultView {
	if src == nil {
		return nil
	}
	return &YongShenResultView{
		NianGanPalace: src.NianGanPalace, NianGanGongFacts: ganGongFactViews(src.NianGanGongFacts),
		Symbols: src.Symbols,
	}
}

// JinhanChartMethodView maps the output-only method codes of 金函玉镜 to 命理
// display terms. scope/school are symmetric input codes and stay unchanged.
type JinhanChartMethodView struct {
	Scope        Scope  `json:"scope"`
	School       School `json:"school"`
	DingjuMethod string `json:"dingju_method"`
	MethodSource string `json:"method_source"`
	DayBoundary  string `json:"day_boundary"`
	Geometry     string `json:"geometry"`
	StarMode     string `json:"xing_mode"`
	DoorMode     string `json:"men_mode"`
	SpiritMode   string `json:"shen_mode"`
	SpiritFlight string `json:"shen_flight"`
	DunSource    string `json:"dun_source"`
	JieQi        string `json:"jie_qi"`
	Basis        string `json:"basis"`
}

func jinhanChartMethodView(m JinhanChartMethod) JinhanChartMethodView {
	return JinhanChartMethodView{
		Scope: m.Scope, School: m.School,
		DingjuMethod: string(m.DingjuMethod),
		MethodSource: displayFrom(methodSourceDisplay, m.MethodSource),
		DayBoundary:  displayFrom(dayBoundaryDisplay, m.DayBoundary),
		Geometry:     displayFrom(geometryDisplay, m.Geometry),
		StarMode:     displayFrom(starModeDisplay, m.StarMode),
		DoorMode:     displayFrom(doorModeDisplay, m.DoorMode),
		SpiritMode:   displayFrom(spiritModeDisplay, m.SpiritMode),
		SpiritFlight: displayFrom(spiritFlightDisplay, m.SpiritFlight),
		DunSource:    displayFrom(jinhanDunSourceDisplay, m.DunSource),
		JieQi:        m.JieQi, Basis: m.MethodYiJu,
	}
}
