package bazi

import (
	"sort"
	"strings"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// LuRoot records a visible stem whose traditional Lu branch appears in one of
// the four palaces. Hidden stems are not promoted to Lu roots.
type LuRoot struct {
	ShiShen string `json:"shi_shen"`
	Gan     string `json:"gan"`
	Zhi     string `json:"zhi"`
}

// NatalRelationMember locates one stem or branch participant in its palace.
type NatalRelationMember struct {
	Pillar PillarIndex `json:"zhu"`
	Gan    string `json:"gan,omitempty"`
	Zhi    string `json:"zhi,omitempty"`
}

// NatalRelationFact is the canonical natal relation projection. Group is the
// stable match key; members carry palace-resolved evidence. Name-like text is
// presentation and is intentionally not part of the domain fact.
type NatalRelationFact struct {
	ID        string                `json:"id"`
	Kind      NatalRelationKind     `json:"relation_kind"`
	Group     string                `json:"group"`
	Members   []NatalRelationMember `json:"members"`
	Wuxing    string                `json:"wuxing,omitempty"`
	Position  string                `json:"position,omitempty"`
	Contested bool                  `json:"contested,omitempty"`
}

// TenGodState aggregates all visible and hidden occurrences of one ten god.
// Rooted, timely and strength follow the underlying element; transparent and
// hidden describe this specific ten god.
type TenGodState struct {
	ShiShen     string `json:"shi_shen"`
	Wuxing      string `json:"wuxing"`
	Transparent bool   `json:"transparent"`
	Hidden      bool   `json:"hidden"`
	Rooted      bool   `json:"rooted"`
	Timely      bool   `json:"timely"`
	Count       int    `json:"count"`
	Strength    string `json:"strength"`
}

// ElementState records engine-owned seasonal, rooted and control facts for one
// element. Combined strength is strong when timely, or transparent and rooted;
// weak only when untimely, not transparent and rootless.
type ElementState struct {
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

type tenGodAccumulator struct {
	wuxing      ganzhi.Wuxing
	transparent bool
	hidden      bool
	count       int
}

func computeTenGodStates(c Chart) ([]TenGodState, map[ganzhi.Wuxing]ElementState) {
	pillars := [4]Pillar{c.Nian, c.Yue, c.Ri, c.Shi}
	transparent := map[ganzhi.Wuxing]bool{}
	rooted := map[ganzhi.Wuxing]bool{}
	accumulators := map[ganzhi.ShiShen]*tenGodAccumulator{}

	for _, pillar := range pillars {
		transparent[ganzhi.GanWuxing(pillar.Gan)] = true
		hiddenStems := []ganzhi.Gan{pillar.CangGan.Main}
		if pillar.CangGan.Mid != nil {
			hiddenStems = append(hiddenStems, *pillar.CangGan.Mid)
		}
		if pillar.CangGan.Minor != nil {
			hiddenStems = append(hiddenStems, *pillar.CangGan.Minor)
		}
		for _, gan := range hiddenStems {
			if gan != 0 {
				rooted[ganzhi.GanWuxing(gan)] = true
			}
		}
		for _, item := range pillar.ShiShens {
			if item.Source != sourceGan && !hiddenTenGodSource(item.Source) {
				continue
			}
			state := accumulators[item.ShiShen]
			if state == nil {
				state = &tenGodAccumulator{wuxing: ganzhi.GanWuxing(item.Gan)}
				accumulators[item.ShiShen] = state
			}
			if item.Source == sourceGan {
				state.transparent = true
			} else {
				state.hidden = true
			}
			state.count++
		}
	}

	elements := computeElementStates(c, transparent, rooted)
	states := make([]TenGodState, 0, len(accumulators))
	for shiShen := ganzhi.ShiShen(0); shiShen <= ganzhi.ShiShenZhengYin; shiShen++ {
		accumulator := accumulators[shiShen]
		if accumulator == nil {
			continue
		}
		element := elements[accumulator.wuxing]
		// 十神强度按 ten_god_strength_rules 表判定（《子平真诠》旺衰强弱：
		// 得令且透/根为强，失令且不透不根为弱；旺而不强判 neutral）。
		strength := "neutral"
		seasonStrong := element.SeasonStrength == "strong"
		for _, rule := range tenGodStrengthRules {
			if rule.Season != nil && *rule.Season != seasonStrong {
				continue
			}
			if rule.Transparent != nil && *rule.Transparent != accumulator.transparent {
				continue
			}
			if rule.Rooted != nil && *rule.Rooted != element.Rooted {
				continue
			}
			strength = rule.Kind
			break
		}
		states = append(states, TenGodState{
			ShiShen:     ganzhi.DecodeShiShen(shiShen),
			Wuxing:      ganzhi.DecodeWuxing(accumulator.wuxing),
			Transparent: accumulator.transparent,
			Hidden:      accumulator.hidden,
			Rooted:      element.Rooted,
			Timely:      element.SeasonStrength == "strong",
			Count:       accumulator.count,
			Strength:    strength,
		})
	}
	return states, elements
}

func hiddenTenGodSource(source string) bool {
	return source == sourceMainQi || source == sourceMidQi || source == sourceMinQi
}

func computeElementStates(c Chart, transparent, rooted map[ganzhi.Wuxing]bool) map[ganzhi.Wuxing]ElementState {
	all := []ganzhi.Wuxing{ganzhi.WxMu, ganzhi.WxHuo, ganzhi.WxTu, ganzhi.WxJin, ganzhi.WxShui}
	result := make(map[ganzhi.Wuxing]ElementState, len(all))
	for _, element := range all {
		state := ElementState{
			Wuxing:       ganzhi.DecodeWuxing(element),
			Season:       c.FuYi.WangShuai[ganzhi.DecodeWuxing(element)],
			Transparent:  transparent[element],
			Rooted:       rooted[element],
			Controls:     ganzhi.DecodeWuxing(elementControlledBy(element)),
			ControlledBy: ganzhi.DecodeWuxing(elementThatControls(element)),
			Reasons:      []string{},
		}
		switch c.FuYi.WangShuai[ganzhi.DecodeWuxing(element)] {
		case ganzhi.DecodeWangShuai(ganzhi.WSWang), ganzhi.DecodeWangShuai(ganzhi.WSXiang):
			state.SeasonStrength = "strong"
		case ganzhi.DecodeWangShuai(ganzhi.WSXiu), ganzhi.DecodeWangShuai(ganzhi.WSQiu), ganzhi.DecodeWangShuai(ganzhi.WSSi):
			state.SeasonStrength = "weak"
		}
		switch {
		case state.SeasonStrength == "strong":
			state.Strength = "strong"
			state.Reasons = append(state.Reasons, "timely")
		case state.Transparent && state.Rooted:
			state.Strength = "strong"
			state.Reasons = append(state.Reasons, "transparent_rooted")
		case state.SeasonStrength == "weak" && !state.Transparent && !state.Rooted:
			state.Strength = "weak"
			state.Reasons = append(state.Reasons, "untimely", "not_transparent", "not_rooted")
		default:
			state.Strength = "neutral"
		}
		result[element] = state
	}
	for _, element := range all {
		state := result[element]
		controller := elementThatControls(element)
		state.ControllerStrength = result[controller].SeasonStrength
		result[element] = state
	}
	return result
}

func elementStateList(states map[ganzhi.Wuxing]ElementState) []ElementState {
	result := make([]ElementState, 0, len(states))
	for _, element := range []ganzhi.Wuxing{ganzhi.WxMu, ganzhi.WxHuo, ganzhi.WxTu, ganzhi.WxJin, ganzhi.WxShui} {
		result = append(result, states[element])
	}
	return result
}

func computeLuRoots(c Chart) []LuRoot {
	roots := make([]LuRoot, 0, 4)
	pillars := [4]Pillar{c.Nian, c.Yue, c.Ri, c.Shi}
	branches := map[ganzhi.Zhi]bool{}
	for _, pillar := range pillars {
		branches[pillar.Zhi] = true
	}
	for _, pillar := range pillars {
		for _, item := range pillar.ShiShens {
			if item.Source != sourceGan || !branches[ganzhi.LuZhi(item.Gan)] {
				continue
			}
			roots = append(roots, LuRoot{
				ShiShen: ganzhi.DecodeShiShen(item.ShiShen),
				Gan:     ganzhi.GanName(item.Gan),
				Zhi:     ganzhi.ZhiName(ganzhi.LuZhi(item.Gan)),
			})
		}
	}
	return roots
}

// ChartConclusions are the small set of natal conclusions used to fill the
// published GuanSha, CaiKu, and GeJu transparency fields.
type ChartConclusions struct {
	PatternGodTransparent bool
	OfficerKillingCleaned bool
	WealthTombPresent     bool
	WealthStarInTomb      bool
}

func wealthTomb(element ganzhi.Wuxing) ganzhi.Zhi {
	switch element {
	case ganzhi.WxJin:
		return ganzhi.ZhiChou
	case ganzhi.WxMu:
		return ganzhi.ZhiWei
	case ganzhi.WxShui:
		return ganzhi.ZhiChen
	case ganzhi.WxHuo:
		return ganzhi.ZhiXu
	case ganzhi.WxTu:
		return ganzhi.ZhiChen
	default:
		return 0
	}
}

func computeChartConclusions(c Chart) ChartConclusions {
	var conclusions ChartConclusions
	pillars := [4]Pillar{c.Nian, c.Yue, c.Ri, c.Shi}

	dayElement := ganzhi.GanWuxing(c.Ri.Gan)
	wealthElement := ganzhi.Wuxing(0)
	for _, element := range []ganzhi.Wuxing{ganzhi.WxMu, ganzhi.WxHuo, ganzhi.WxTu, ganzhi.WxJin, ganzhi.WxShui} {
		if ganzhi.Ke(dayElement, element) {
			wealthElement = element
			break
		}
	}
	tomb := wealthTomb(wealthElement)
	if tomb != 0 {
		for _, pillar := range pillars {
			if pillar.Zhi != tomb {
				continue
			}
			conclusions.WealthTombPresent = true
			for _, item := range pillar.ShiShens {
				if item.ShiShen == ganzhi.ShiShenZhengCai || item.ShiShen == ganzhi.ShiShenPianCai {
					conclusions.WealthStarInTomb = true
					break
				}
			}
		}
	}

	hasOfficer := false
	hasKilling := false
	killingPillars := make(map[int]bool)
	for index, pillar := range pillars {
		for _, item := range pillar.ShiShens {
			if item.Source != sourceGan {
				continue
			}
			switch item.ShiShen {
			case ganzhi.ShiShenZhengGuan:
				hasOfficer = true
			case ganzhi.ShiShenQiSha:
				hasKilling = true
				killingPillars[index] = true
			}
		}
	}
	for _, relation := range c.GuanXi {
		if relation.Kind != relationKindZhiLiuHe && relation.Kind != relationKindLiuChong {
			continue
		}
		for _, member := range relation.Members {
			index := int(member.Pillar)
			if index >= 0 && killingPillars[index] {
				conclusions.OfficerKillingCleaned = hasOfficer && hasKilling
			}
		}
	}

	conclusions.PatternGodTransparent = patternGodTransparent(c)
	return conclusions
}

func computeNatalRelations(bz ganzhi.Bazi) []NatalRelationFact {
	return buildNatalRelationsFromDetection(bz, DetectRelations(bz))
}

func buildNatalRelationsFromDetection(bz ganzhi.Bazi, detection RelationDetection) []NatalRelationFact {
	facts := make([]NatalRelationFact, 0, 8)
	seen := map[string]bool{}
	add := func(fact NatalRelationFact) {
		if fact.Group == "" || len(fact.Members) == 0 {
			return
		}
		fact.ID = natalRelationFactID(fact)
		if seen[fact.ID] {
			return
		}
		seen[fact.ID] = true
		facts = append(facts, fact)
	}

	for _, relation := range detection.GanHe {
		kind := relationKindGanHeCandidate
		if relation.Position == positionAdjacent {
			kind = relationKindGanHe
		}
		add(NatalRelationFact{
			ID:    DecodeNatalRelationKind(kind) + ":" + ganHeGroupLabel(relation.GanA, relation.GanB),
			Kind:  kind,
			Group: ganHeGroupLabel(relation.GanA, relation.GanB),
			Members: []NatalRelationMember{
				{Pillar: PillarIndex(relation.PillarA), Gan: ganzhi.GanName(relation.GanA)},
				{Pillar: PillarIndex(relation.PillarB), Gan: ganzhi.GanName(relation.GanB)},
			},
			Wuxing:    wuxingOutput(relation.HeElement),
			Position:  DecodeRelationPosition(relation.Position),
			Contested: relation.Contested,
		})
	}
	for _, relation := range detection.GanChong {
		kind := relationKindGanChongCandidate
		if relation.Position == positionAdjacent {
			kind = relationKindGanChong
		}
		group := ganChongGroupLabel(relation.GanA, relation.GanB)
		add(NatalRelationFact{
			ID:    DecodeNatalRelationKind(kind) + ":" + group,
			Kind:  kind,
			Group: group,
			Members: []NatalRelationMember{
				{Pillar: PillarIndex(relation.PillarA), Gan: ganzhi.GanName(relation.GanA)},
				{Pillar: PillarIndex(relation.PillarB), Gan: ganzhi.GanName(relation.GanB)},
			},
			Position: DecodeRelationPosition(relation.Position),
		})
	}
	for _, relation := range detection.ZhiLiuHe {
		group := canonicalZhiPairGroup(relationKindZhiLiuHe, relation.ZhiA, relation.ZhiB)
		add(newZhiPairFact(relationKindZhiLiuHe, group, relation.ZhiA, relation.PillarA, relation.ZhiB, relation.PillarB, relation.Element))
	}
	for _, relation := range detection.SanHe {
		add(newTripleRelationFact(relationKindSanHe, bz, relation))
	}
	for _, relation := range detection.SanHePartial {
		add(newTripleRelationFact(relationKindSanHePartial, bz, relation))
	}
	for _, relation := range detection.SanHui {
		add(newTripleRelationFact(relationKindSanHui, bz, relation))
	}
	for _, relation := range detection.LiuChong {
		group := canonicalZhiPairGroup(relationKindLiuChong, relation.ZhiA, relation.ZhiB)
		add(newZhiPairFact(relationKindLiuChong, group, relation.ZhiA, relation.PillarA, relation.ZhiB, relation.PillarB, relation.Element))
	}
	for _, relation := range detection.LiuHai {
		group := canonicalZhiPairGroup(relationKindLiuHai, relation.ZhiA, relation.ZhiB)
		add(newZhiPairFact(relationKindLiuHai, group, relation.ZhiA, relation.PillarA, relation.ZhiB, relation.PillarB, relation.Element))
	}
	for _, relation := range detection.LiuXing {
		group := canonicalXingGroup(relation.ZhiA, relation.ZhiB)
		add(newZhiPairFact(relationKindLiuXing, group, relation.ZhiA, relation.PillarA, relation.ZhiB, relation.PillarB, relation.Element))
	}
	// LiuXing detection keeps pairwise evidence. A complete three-branch 刑组
	// is an additional canonical fact and must not be lost when the pairwise
	// projection is later retired.
	zhiCounts := map[ganzhi.Zhi]int{}
	for _, pillar := range bz.Slice() {
		zhiCounts[pillar.Zhi]++
	}
	for _, group := range ganzhi.XingGroups {
		if group.Type == ganzhi.XingKindZi || !allBranchesPresent(zhiCounts, group.Zhi) {
			continue
		}
		groupLabel := xingGroupLabel(group.Zhi)
		members := make([]NatalRelationMember, 0, len(group.Zhi))
		for _, branch := range group.Zhi {
			members = append(members, NatalRelationMember{
				Pillar: pillarWithZhi(bz, branch),
				Zhi:    ganzhi.ZhiName(branch),
			})
		}
		add(NatalRelationFact{
			ID: "liu_xing:" + groupLabel, Kind: relationKindLiuXing, Group: groupLabel, Members: members,
		})
	}
	for _, relation := range detection.LiuPo {
		group := canonicalZhiPairGroup(relationKindLiuPo, relation.ZhiA, relation.ZhiB)
		add(newZhiPairFact(relationKindLiuPo, group, relation.ZhiA, relation.PillarA, relation.ZhiB, relation.PillarB, relation.Element))
	}
	for _, relation := range detection.AnHe {
		group := canonicalZhiPairGroup(relationKindAnHe, relation.ZhiA, relation.ZhiB)
		add(newZhiPairFact(relationKindAnHe, group, relation.ZhiA, relation.PillarA, relation.ZhiB, relation.PillarB, relation.Element))
	}
	sort.Slice(facts, func(i, j int) bool {
		if facts[i].Kind != facts[j].Kind {
			return DecodeNatalRelationKind(facts[i].Kind) < DecodeNatalRelationKind(facts[j].Kind)
		}
		if facts[i].Group != facts[j].Group {
			return facts[i].Group < facts[j].Group
		}
		return facts[i].ID < facts[j].ID
	})
	return facts
}

func natalRelationFactID(fact NatalRelationFact) string {
	keys := make([]string, 0, len(fact.Members))
	for _, member := range fact.Members {
		value := member.Gan
		if value == "" {
			value = member.Zhi
		}
		keys = append(keys, DecodePillarIndex(member.Pillar)+":"+value)
	}
	return DecodeNatalRelationKind(fact.Kind) + ":" + fact.Group + "@" + strings.Join(keys, "|")
}

func wuxingOutput(w ganzhi.Wuxing) string {
	if w == 0 {
		return ""
	}
	return ganzhi.DecodeWuxing(w)
}

func newZhiPairFact(kind NatalRelationKind, group string, zhiA ganzhi.Zhi, pillarA int, zhiB ganzhi.Zhi, pillarB int, wuxing ganzhi.Wuxing) NatalRelationFact {
	return NatalRelationFact{
		ID:    DecodeNatalRelationKind(kind) + ":" + group,
		Kind:  kind,
		Group: group,
		Members: []NatalRelationMember{
			{Pillar: PillarIndex(pillarA), Zhi: ganzhi.ZhiName(zhiA)},
			{Pillar: PillarIndex(pillarB), Zhi: ganzhi.ZhiName(zhiB)},
		},
		Wuxing: wuxingOutput(wuxing),
	}
}

func newTripleRelationFact(kind NatalRelationKind, bz ganzhi.Bazi, relation TripleRelationDetection) NatalRelationFact {
	group := branchGroupLabel(relation.Branches)
	members := make([]NatalRelationMember, 0, len(relation.Branches))
	for _, branch := range relation.Branches {
		members = append(members, NatalRelationMember{
			Pillar: pillarWithZhi(bz, branch),
			Zhi:    ganzhi.ZhiName(branch),
		})
	}
	return NatalRelationFact{
		ID:      DecodeNatalRelationKind(kind) + ":" + group,
		Kind:    kind,
		Group:   group,
		Members: members,
		Wuxing:  wuxingOutput(relation.Element),
	}
}

func canonicalZhiPairGroup(kind NatalRelationKind, a, b ganzhi.Zhi) string {
	switch kind {
	case relationKindZhiLiuHe:
		for _, relation := range ganzhi.ZhiHes {
			if (relation.A == a && relation.B == b) || (relation.A == b && relation.B == a) {
				return ganzhi.ZhiName(relation.A) + ganzhi.ZhiName(relation.B)
			}
		}
	case relationKindLiuChong:
		for _, relation := range ganzhi.ChongPairs {
			if (relation.A == a && relation.B == b) || (relation.A == b && relation.B == a) {
				return ganzhi.ZhiName(relation.A) + ganzhi.ZhiName(relation.B)
			}
		}
	case relationKindLiuHai:
		for _, relation := range ganzhi.HaiPairs {
			if (relation.A == a && relation.B == b) || (relation.A == b && relation.B == a) {
				return ganzhi.ZhiName(relation.A) + ganzhi.ZhiName(relation.B)
			}
		}
	case relationKindLiuPo:
		for _, relation := range ganzhi.PoPairs {
			if (relation.A == a && relation.B == b) || (relation.A == b && relation.B == a) {
				return ganzhi.ZhiName(relation.A) + ganzhi.ZhiName(relation.B)
			}
		}
	case relationKindAnHe:
		for _, relation := range ganzhi.AnHePairs {
			if (relation.A == a && relation.B == b) || (relation.A == b && relation.B == a) {
				return ganzhi.ZhiName(relation.A) + ganzhi.ZhiName(relation.B)
			}
		}
	}
	return ""
}

func canonicalXingGroup(a, b ganzhi.Zhi) string {
	if a == b {
		return ganzhi.ZhiName(a) + ganzhi.ZhiName(b)
	}
	for _, group := range ganzhi.XingGroups {
		if group.Type == ganzhi.XingKindZi {
			continue
		}
		present := make([]ganzhi.Zhi, 0, len(group.Zhi))
		for _, branch := range group.Zhi {
			if branch == a || branch == b {
				present = append(present, branch)
			}
		}
		if len(present) == 2 {
			return branchGroupLabel(present)
		}
	}
	return ""
}

func parseZhiNames(names []string) []ganzhi.Zhi {
	values := make([]ganzhi.Zhi, 0, len(names))
	for _, name := range names {
		value, err := ganzhi.EncodeZhi(name)
		if err != nil {
			continue
		}
		values = append(values, value)
	}
	return values
}

func pillarWithZhi(bz ganzhi.Bazi, target ganzhi.Zhi) PillarIndex {
	for index, pillar := range bz.Slice() {
		if pillar.Zhi == target {
			return PillarIndex(index)
		}
	}
	return 0
}

// ganHeGroupLabel canonicalizes an observed adjacent stem pair to the table's
// stable group order. Non-adjacent stems never enter this function.
func ganHeGroupLabel(a, b ganzhi.Gan) string {
	for _, relation := range ganzhi.GanHes {
		if (relation.A == a && relation.B == b) || (relation.A == b && relation.B == a) {
			return ganzhi.GanName(relation.A) + ganzhi.GanName(relation.B)
		}
	}
	return ""
}

func ganChongGroupLabel(a, b ganzhi.Gan) string {
	for _, relation := range ganzhi.GanChongPairs {
		if (relation.A == a && relation.B == b) || (relation.A == b && relation.B == a) {
			return ganzhi.GanName(relation.A) + ganzhi.GanName(relation.B)
		}
	}
	return ""
}

func allBranchesPresent(counts map[ganzhi.Zhi]int, branches []ganzhi.Zhi) bool {
	for _, branch := range branches {
		if counts[branch] == 0 {
			return false
		}
	}
	return true
}

func branchGroupLabel(branches []ganzhi.Zhi) string {
	parts := make([]string, 0, len(branches))
	for _, branch := range branches {
		parts = append(parts, ganzhi.ZhiName(branch))
	}
	return strings.Join(parts, "")
}

func xingGroupLabel(branches []ganzhi.Zhi) string {
	present := map[ganzhi.Zhi]bool{}
	for _, branch := range branches {
		present[branch] = true
	}
	if present[ganzhi.ZhiChou] && present[ganzhi.ZhiXu] && present[ganzhi.ZhiWei] {
		return "丑戌未"
	}
	return branchGroupLabel(branches)
}

func patternGodTransparent(c Chart) bool {
	target, ok := patternTenGods[c.GeJu.Pattern]
	if !ok {
		return false
	}
	for _, pillar := range []Pillar{c.Nian, c.Yue, c.Ri, c.Shi} {
		for _, item := range pillar.ShiShens {
			if item.Source == sourceGan && item.ShiShen == target {
				return true
			}
		}
	}
	return false
}

var patternTenGods = map[GeJuKind]ganzhi.ShiShen{
	geJuZhengGuan: ganzhi.ShiShenZhengGuan,
	geJuQiSha: ganzhi.ShiShenQiSha,
	geJuZhengCai: ganzhi.ShiShenZhengCai,
	geJuPianCai: ganzhi.ShiShenPianCai,
	geJuZhengYin: ganzhi.ShiShenZhengYin,
	geJuPianYin: ganzhi.ShiShenPianYin,
	geJuShiShen: ganzhi.ShiShenShiShen,
	geJuShangGuan: ganzhi.ShiShenShangGuan,
}
