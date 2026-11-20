package bazi

import _ "embed"

import (
	"encoding/json"
	"fmt"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// ── Data-driven strength & congge rule lookup ────────────────

// strengthRule is the symbol-typed strength classification rule.
type strengthRule struct {
	Rule      int
	Root      RootKind
	Season    []ganzhi.WangShuai
	Strength  StrengthLevel
	Upgrade   *StrengthLevel
	NeedYinBi int
	Note      string
}

// congGeRule is the symbol-typed 从格 detection rule.
type congGeRule struct {
	Rule       int
	Type       CongGeKind
	Conditions congGeConditions
	Yong       YongXiJiKind
	Xi         YongXiJiKind
	Ji         YongXiJiKind
	Source     string
}

type congGeConditions struct {
	Root       RootKind
	Season     []ganzhi.WangShuai
	Gan        ganConditions
	RootDetail []RootKind
}

type ganConditions struct {
	YinBiMin      *int `json:"yin_bi_min,omitempty"`
	GuanShaTouGan *int `json:"guan_sha_tou_gan,omitempty"`
	GuanShaMin    *int `json:"guan_sha_min,omitempty"`
	YinTouGan     *int `json:"yin_tou_gan,omitempty"` // 印透干数量 (excluding 比劫)
	CaiMin        *int `json:"cai_min,omitempty"`
	BiJieTouGan   *int `json:"bi_jie_tou_gan,omitempty"`
	ShiShangMin   *int `json:"shi_shang_min,omitempty"`
}

// ── Embed and load ──

//go:embed data/strength_rules.json
var strengthRulesJSON []byte

//go:embed data/congge_rules.json
var congGeRulesJSON []byte

var (
	strengthRules []strengthRule
	congGeRules   []congGeRule
)

func loadStrengthRules() error {
	var raw []struct {
		Rule      int      `json:"rule"`
		Root      string   `json:"root"`
		Season    []string `json:"season"`
		Strength  string   `json:"strength"`
		Upgrade   *string  `json:"upgrade,omitempty"`
		NeedYinBi int      `json:"need_yinbi,omitempty"`
		Note      string   `json:"note"`
	}
	if err := json.Unmarshal(strengthRulesJSON, &raw); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	for _, r := range raw {
		rule := strengthRule{Rule: r.Rule, NeedYinBi: r.NeedYinBi, Note: r.Note}
		var ok bool
		if rule.Root, ok = EncodeRootKind(r.Root); !ok {
			return fmt.Errorf("strength rule %d: invalid root %q", r.Rule, r.Root)
		}
		if rule.Strength, ok = EncodeStrengthLevel(r.Strength); !ok {
			return fmt.Errorf("strength rule %d: invalid strength %q", r.Rule, r.Strength)
		}
		for _, name := range r.Season {
			ws, err := ganzhi.EncodeWangShuai(name)
			if err != nil {
				return fmt.Errorf("strength rule %d: invalid season %q", r.Rule, name)
			}
			rule.Season = append(rule.Season, ws)
		}
		if r.Upgrade != nil {
			up, ok := EncodeStrengthLevel(*r.Upgrade)
			if !ok {
				return fmt.Errorf("strength rule %d: invalid upgrade %q", r.Rule, *r.Upgrade)
			}
			rule.Upgrade = &up
		}
		strengthRules = append(strengthRules, rule)
	}
	if len(strengthRules) == 0 {
		return fmt.Errorf("empty strength rules")
	}
	return nil
}

func loadCongGeRules() error {
	var raw []struct {
		Rule       int    `json:"rule"`
		Type       string `json:"type"`
		Conditions struct {
			Root       string        `json:"root"`
			Season     []string      `json:"season"`
			Gan        ganConditions `json:"gan"`
			RootDetail []string      `json:"root_detail,omitempty"`
		} `json:"conditions"`
		Yong   string `json:"yong"`
		Xi     string `json:"xi"`
		Ji     string `json:"ji"`
		Source string `json:"source"`
	}
	if err := json.Unmarshal(congGeRulesJSON, &raw); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	for _, r := range raw {
		congKind, ok := EncodeCongGeKind(r.Type)
		if !ok {
			return fmt.Errorf("invalid congge type %q", r.Type)
		}
		rule := congGeRule{Rule: r.Rule, Type: congKind, Source: r.Source}
		if rule.Conditions.Root, ok = EncodeRootKind(r.Conditions.Root); !ok {
			return fmt.Errorf("congge rule %d: invalid root %q", r.Rule, r.Conditions.Root)
		}
		if rule.Yong, ok = EncodeYongXiJiKind(r.Yong); !ok {
			return fmt.Errorf("congge rule %d: invalid yong %q", r.Rule, r.Yong)
		}
		if rule.Xi, ok = EncodeYongXiJiKind(r.Xi); !ok {
			return fmt.Errorf("congge rule %d: invalid xi %q", r.Rule, r.Xi)
		}
		if rule.Ji, ok = EncodeYongXiJiKind(r.Ji); !ok {
			return fmt.Errorf("congge rule %d: invalid ji %q", r.Rule, r.Ji)
		}
		for _, name := range r.Conditions.Season {
			ws, err := ganzhi.EncodeWangShuai(name)
			if err != nil {
				return fmt.Errorf("congge rule %d: invalid season %q", r.Rule, name)
			}
			rule.Conditions.Season = append(rule.Conditions.Season, ws)
		}
		for _, name := range r.Conditions.RootDetail {
			rk, ok := EncodeRootKind(name)
			if !ok {
				return fmt.Errorf("congge rule %d: invalid root_detail %q", r.Rule, name)
			}
			rule.Conditions.RootDetail = append(rule.Conditions.RootDetail, rk)
		}
		rule.Conditions.Gan = r.Conditions.Gan
		congGeRules = append(congGeRules, rule)
	}
	if len(congGeRules) == 0 {
		return fmt.Errorf("empty congge rules")
	}
	return nil
}

// ── Root type classification ──

// classifyRoot determines the root type for the day master.
func classifyRoot(riYuan ganzhi.Gan, yueZhi ganzhi.Zhi, cangGan [4]HiddenStemProjection) RootKind {
	dmElem := ganzhi.GanWuxing(riYuan)

	monthCG := cangGan[1]
	if ganzhi.GanWuxing(monthCG.Main) == dmElem {
		return rootMonthMain
	}
	if (monthCG.Mid != nil && ganzhi.GanWuxing(*monthCG.Mid) == dmElem) ||
		(monthCG.Minor != nil && ganzhi.GanWuxing(*monthCG.Minor) == dmElem) {
		return rootMonthMid
	}

	for i, hs := range cangGan {
		if i == 1 {
			continue
		}
		if ganzhi.GanWuxing(hs.Main) == dmElem {
			return rootZhiMain
		}
		if (hs.Mid != nil && ganzhi.GanWuxing(*hs.Mid) == dmElem) ||
			(hs.Minor != nil && ganzhi.GanWuxing(*hs.Minor) == dmElem) {
			return rootZhiMid
		}
	}
	return rootNone
}

// classifySeason returns the season label (旺相休囚死) for the day master.
func classifySeason(riYuan ganzhi.Gan, yueZhi ganzhi.Zhi) ganzhi.WangShuai {
	dmElem := ganzhi.GanWuxing(riYuan)
	monthElem := ganzhi.ZhiWuxing(yueZhi)
	genElem := elementThatGenerates(dmElem)
	ctrlElem := elementThatControls(dmElem)

	switch {
	case monthElem == dmElem:
		return ganzhi.WSWang
	case monthElem == genElem:
		return ganzhi.WSXiang
	case monthElem == ctrlElem:
		return ganzhi.WSSi
	case elementThatControls(monthElem) == dmElem:
		return ganzhi.WSQiu
	default:
		return ganzhi.WSXiu
	}
}

// DayMasterStrength is the qualitative strength verdict used by 扶抑 and 从格.
type DayMasterStrength struct {
	RootType   RootKind
	Season     ganzhi.WangShuai
	Label      StrengthLevel
	YinBiCount int
}

// AnalyzeDayMasterStrength derives the data-driven day-master strength facts.
func AnalyzeDayMasterStrength(
	bz ganzhi.Bazi,
	hiddenStems [4]HiddenStemProjection,
) DayMasterStrength {
	strength := DayMasterStrength{
		RootType:   classifyRoot(bz.Ri.Gan, bz.Yue.Zhi, hiddenStems),
		Season:     classifySeason(bz.Ri.Gan, bz.Yue.Zhi),
		YinBiCount: countYinBi(bz),
	}
	strength.Label = lookupStrength(strength.RootType, strength.Season, strength.YinBiCount)
	return strength
}

// WuxingCounts are transparent and hidden-branch element counts.
type WuxingCounts map[ganzhi.Wuxing]int

// StringCount projects the typed counts into the published wire shape.
func (counts WuxingCounts) StringCount() map[string]int {
	stringCount := make(map[string]int, len(counts))
	for element, count := range counts {
		stringCount[ganzhi.DecodeWuxing(element)] = count
	}
	return stringCount
}

func WuxingCountsAt(
	bz ganzhi.Bazi,
	hiddenStems [4]HiddenStemProjection,
) WuxingCounts {
	counts := make(WuxingCounts, 5)
	for _, pillar := range bz.Slice() {
		counts[ganzhi.GanWuxing(pillar.Gan)]++
	}
	for _, stems := range hiddenStems {
		counts[ganzhi.GanWuxing(stems.Main)]++
		if stems.Mid != nil {
			counts[ganzhi.GanWuxing(*stems.Mid)]++
		}
		if stems.Minor != nil {
			counts[ganzhi.GanWuxing(*stems.Minor)]++
		}
	}
	return counts
}

// countYinBi returns印比 count (比劫同五行+印生我) from the 4 gan.
func countYinBi(bz ganzhi.Bazi) int {
	// Count 印(生我) and 比劫(同我) from年/月/时干 only (excluding日主自身).
	dmElem := ganzhi.GanWuxing(bz.Ri.Gan)
	genElem := elementThatGenerates(dmElem)
	count := 0
	for _, p := range []struct{ gan ganzhi.Gan }{
		{bz.Nian.Gan}, {bz.Yue.Gan}, {bz.Shi.Gan},
	} {
		e := ganzhi.GanWuxing(p.gan)
		if e == dmElem || e == genElem {
			count++
		}
	}
	return count
}

// ── Lookup methods ──

// lookupStrength returns the day master's strength by matching
// root type, season, and 印比 against the strength rules JSON.
func lookupStrength(rootType RootKind, season ganzhi.WangShuai, yinBiCount int) StrengthLevel {
	for _, rule := range strengthRules {
		if rule.Root != rootType {
			continue
		}
		seasonMatch := false
		for _, s := range rule.Season {
			if s == season {
				seasonMatch = true
				break
			}
		}
		if !seasonMatch {
			continue
		}
		if rule.Upgrade != nil && rule.NeedYinBi > 0 && yinBiCount >= rule.NeedYinBi {
			return *rule.Upgrade
		}
		return rule.Strength
	}
	return strengthWeak
}

// lookupCongGe checks the day master against 从格 rules in priority order.
func lookupCongGe(bz ganzhi.Bazi) (CongGeKind, CandidateWuxing, CandidateWuxing, CandidateWuxing) {
	riYuan := bz.Ri.Gan
	dmElem := ganzhi.GanWuxing(riYuan)
	rootType := classifyRoot(riYuan, bz.Yue.Zhi, HiddenStemFacts(bz))
	season := classifySeason(riYuan, bz.Yue.Zhi)
	gan := countGan(bz)

	for _, rule := range congGeRules {
		if !matchRoot(rule.Conditions.Root, rootType, rule.Conditions.RootDetail) {
			continue
		}
		if !matchSeasons(rule.Conditions.Season, season) {
			continue
		}
		if !matchGanConditions(rule.Conditions.Gan, gan) {
			continue
		}
		yong, xi, ji := resolveCongGeYongXiJi(rule, dmElem)
		return rule.Type, yong, xi, ji
	}
	return congGeUnknown, candidateNone, candidateNone, candidateNone
}

// ── Helpers ──

type ganCounts struct {
	biBi     int // 比劫 (same as dmElem)
	yin      int // 印 (generates dmElem)
	guanSha  int // 官杀 (controls dmElem)
	cai      int // 财 (controlled by dmElem)
	shiShang int // 食伤 (drains dmElem)
}

func (sc ganCounts) yinBi() int { return sc.biBi + sc.yin }

func countGan(bz ganzhi.Bazi) ganCounts {
	dmElem := ganzhi.GanWuxing(bz.Ri.Gan)
	var sc ganCounts
	for _, p := range []struct{ gan ganzhi.Gan }{
		{bz.Nian.Gan}, {bz.Yue.Gan}, {bz.Shi.Gan},
	} {
		e := ganzhi.GanWuxing(p.gan)
		switch {
		case e == dmElem:
			sc.biBi++
		case e == elementThatGenerates(dmElem):
			sc.yin++
		case e == elementThatControls(dmElem):
			sc.guanSha++
		case e == elementControlledBy(dmElem):
			sc.cai++
		case e == elementThatDrains(dmElem):
			sc.shiShang++
		}
	}
	return sc
}

func matchRoot(ruleRoot, actualRoot RootKind, rootDetail []RootKind) bool {
	switch ruleRoot {
	case rootAny:
		return true
	case rootNone:
		return actualRoot == rootNone
	case rootMid:
		return actualRoot == rootMonthMid || actualRoot == rootZhiMid
	default:
		return ruleRoot == actualRoot
	}
}

func matchSeasons(allowed []ganzhi.WangShuai, actual ganzhi.WangShuai) bool {
	for _, s := range allowed {
		if s == actual {
			return true
		}
	}
	return false
}

func matchGanConditions(cond ganConditions, sc ganCounts) bool {
	if cond.YinBiMin != nil && sc.yinBi() < *cond.YinBiMin {
		return false
	}
	if cond.GuanShaTouGan != nil && sc.guanSha > *cond.GuanShaTouGan {
		return false
	}
	if cond.GuanShaMin != nil && sc.guanSha < *cond.GuanShaMin {
		return false
	}
	if cond.YinTouGan != nil && sc.yin > *cond.YinTouGan {
		return false
	}
	if cond.CaiMin != nil && sc.cai < *cond.CaiMin {
		return false
	}
	if cond.BiJieTouGan != nil && sc.biBi > *cond.BiJieTouGan {
		return false
	}
	if cond.ShiShangMin != nil && sc.shiShang < *cond.ShiShangMin {
		return false
	}
	return true
}

func resolveCongGeYongXiJi(rule congGeRule, dmElem ganzhi.Wuxing) (yong, xi, ji CandidateWuxing) {
	resolve := func(symbol YongXiJiKind) CandidateWuxing {
		switch symbol {
		case yxjDmElem:
			return candidateFromWuxing(dmElem)
		case yxjKeDm:
			return candidateFromWuxing(elementThatControls(dmElem))
		case yxjXiFromDm:
			return candidateFromWuxing(elementThatDrains(dmElem))
		case yxjGuansha:
			return candidateFromWuxing(elementThatControls(dmElem))
		case yxjCai:
			return candidateFromWuxing(elementControlledBy(dmElem))
		case yxjYin:
			return candidateFromWuxing(elementThatGenerates(dmElem))
		case yxjBiJie:
			return candidateFromWuxing(dmElem)
		case yxjShiShang:
			return candidateFromWuxing(elementThatDrains(dmElem))
		default:
			return candidateNone
		}
	}
	return resolve(rule.Yong), resolve(rule.Xi), resolve(rule.Ji)
}
