package bazi

import "github.com/ml8s/mingligo/pkg/ganzhi"

// RelationPosition is the engine-owned strength hint for a paired relation:
// adjacent pairs are tight combinations, others are candidates.
type RelationPosition uint8

const (
	positionSeparated RelationPosition = iota
	positionAdjacent
	positionRemote
)

// DecodeRelationPosition decodes a RelationPosition symbol into its display text.
func DecodeRelationPosition(p RelationPosition) string {
	switch p {
	case positionAdjacent:
		return "adjacent"
	case positionRemote:
		return "remote"
	default:
		return "separated"
	}
}

func EncodeRelationPosition(name string) (RelationPosition, bool) {
	switch name {
	case "adjacent":
		return positionAdjacent, true
	case "remote":
		return positionRemote, true
	case "separated":
		return positionSeparated, true
	default:
		return 0, false
	}
}

// SanQiKind identifies the 三奇贵人 pattern kind (a bazi-domain symbol).
type SanQiKind uint8

const (
	sanQiNone      SanQiKind = iota
	sanQiTianShang           // 天上三奇（甲戊庚）
	sanQiDiXia               // 地下三奇（乙丙丁）
	sanQiRenZhong            // 人中三奇（壬癸辛）
)

// Code returns the short type code (天上/地下/人中) of a SanQiKind.
func (k SanQiKind) Code() string {
	switch k {
	case sanQiTianShang:
		return "天上"
	case sanQiDiXia:
		return "地下"
	case sanQiRenZhong:
		return "人中"
	default:
		return ""
	}
}

// DecodeSanQiKind decodes a SanQiKind symbol into its full display name.
func DecodeSanQiKind(k SanQiKind) string {
	switch k {
	case sanQiTianShang:
		return "天上三奇（甲戊庚）"
	case sanQiDiXia:
		return "地下三奇（乙丙丁）"
	case sanQiRenZhong:
		return "人中三奇（壬癸辛）"
	default:
		return ""
	}
}

// ── Rule-table domain symbols ──────────────────────────────────────────────
// These are closed domain symbols used as keys/values in strength & yong-ji
// rule tables. Their text form is the stable rule-table identifier; text↔symbol
// conversion happens only via Encode (inbound) / Decode (outbound).

// StrengthLevel is the qualitative day-master strength verdict.
type StrengthLevel uint8

const (
	strengthUnknown StrengthLevel = iota
	strengthStrong                // 身强
	strengthBalance               // 中和
	strengthWeak                  // 身弱
)

var strengthLevelNames = [...]string{"", "身强", "中和", "身弱"}

func DecodeStrengthLevel(s StrengthLevel) string {
	if int(s) < len(strengthLevelNames) {
		return strengthLevelNames[s]
	}
	return ""
}

func EncodeStrengthLevel(name string) (StrengthLevel, bool) {
	for i, n := range strengthLevelNames {
		if n != "" && n == name {
			return StrengthLevel(i), true
		}
	}
	return strengthUnknown, false
}

// RootKind is the day-master root classification (also used as a rule-table key).
type RootKind uint8

const (
	rootUnknown   RootKind = iota
	rootAny                // any (wildcard in rules)
	rootNone               // none
	rootMid                // mid (wildcard for month_mid|zhi_mid in rules)
	rootMonthMain          // month_main
	rootMonthMid           // month_mid
	rootZhiMain            // zhi_main
	rootZhiMid             // zhi_mid
	rootMinor              // minor (root_detail)
)

var rootKindNames = [...]string{"", "any", "none", "mid", "month_main", "month_mid", "zhi_main", "zhi_mid", "minor"}

func DecodeRootKind(r RootKind) string {
	if int(r) < len(rootKindNames) {
		return rootKindNames[r]
	}
	return ""
}

func EncodeRootKind(name string) (RootKind, bool) {
	for i, n := range rootKindNames {
		if n != "" && n == name {
			return RootKind(i), true
		}
	}
	return rootUnknown, false
}

// RelationKind is the five-element relation of a yong-ji rule ("克我者" etc).
type RelationKind uint8

const (
	relationUnknown   RelationKind = iota
	relationKeWo                   // 克我者
	relationShengWo                // 生我者
	relationTongWo                 // 同我者
	relationShengKeWo              // 生克我者
)

var relationKindNames = [...]string{"", "克我者", "生我者", "同我者", "生克我者"}

func DecodeRelationKind(r RelationKind) string {
	if int(r) < len(relationKindNames) {
		return relationKindNames[r]
	}
	return ""
}

func EncodeRelationKind(name string) (RelationKind, bool) {
	for i, n := range relationKindNames {
		if n != "" && n == name {
			return RelationKind(i), true
		}
	}
	return relationUnknown, false
}

// YongXiJiKind selects a five-element relative to the day master in 从格 rules.
type YongXiJiKind uint8

const (
	yxjUnknown  YongXiJiKind = iota
	yxjDmElem                // dm_elem
	yxjKeDm                  // ke_dm
	yxjXiFromDm              // xi_from_dm
	yxjGuansha               // guansha
	yxjCai                   // cai
	yxjYin                   // yin
	yxjBiJie                 // bijie
	yxjShiShang              // shishang
)

var yongXiJiKindNames = [...]string{"", "dm_elem", "ke_dm", "xi_from_dm", "guansha", "cai", "yin", "bijie", "shishang"}

func DecodeYongXiJiKind(y YongXiJiKind) string {
	if int(y) < len(yongXiJiKindNames) {
		return yongXiJiKindNames[y]
	}
	return ""
}

func EncodeYongXiJiKind(name string) (YongXiJiKind, bool) {
	for i, n := range yongXiJiKindNames {
		if n != "" && n == name {
			return YongXiJiKind(i), true
		}
	}
	return yxjUnknown, false
}

// PillarIndex identifies one of the four pillars (年/月/日/时).
type PillarIndex uint8

const (
	pillarNian PillarIndex = iota
	pillarYue
	pillarRi
	pillarShi
)

var pillarCodes = [...]string{"nian", "yue", "ri", "shi"}

func DecodePillarIndex(p PillarIndex) string { return pillarCodes[p] }
func (p PillarIndex) MarshalJSON() ([]byte, error) {
	return []byte(`"` + DecodePillarIndex(p) + `"`), nil
}
func EncodePillarIndex(code string) (PillarIndex, bool) {
	for i, c := range pillarCodes {
		if c == code {
			return PillarIndex(i), true
		}
	}
	return 0, false
}

// RelationKind is a relation fact kind (六合/六冲/三合/...).
type NatalRelationKind uint8

const (
	natalRelationUnknown          NatalRelationKind = iota
	relationKindZhiLiuHe                            // zhi_liu_he
	relationKindLiuChong                            // liu_chong
	relationKindLiuHai                              // liu_hai
	relationKindLiuPo                               // liu_po
	relationKindAnHe                                // an_he
	relationKindLiuXing                             // liu_xing
	relationKindGanHe                               // gan_he
	relationKindGanHeCandidate                      // gan_he_candidate
	relationKindGanChong                            // gan_chong
	relationKindGanChongCandidate                   // gan_chong_candidate
	relationKindSanHe                               // san_he
	relationKindSanHePartial                        // san_he_partial
	relationKindSanHui                              // san_hui
)

var relationKindCodes = [...]string{
	"", "zhi_liu_he", "liu_chong", "liu_hai", "liu_po", "an_he", "liu_xing",
	"gan_he", "gan_he_candidate", "gan_chong", "gan_chong_candidate",
	"san_he", "san_he_partial", "san_hui",
}

func DecodeNatalRelationKind(r NatalRelationKind) string { return relationKindCodes[r] }
func (r NatalRelationKind) MarshalJSON() ([]byte, error) {
	return []byte(`"` + DecodeNatalRelationKind(r) + `"`), nil
}
func EncodeNatalRelationKind(code string) (NatalRelationKind, bool) {
	for i, c := range relationKindCodes {
		if c != "" && c == code {
			return NatalRelationKind(i), true
		}
	}
	return natalRelationUnknown, false
}

// GeJuKind is one of the 格局 (正官格/七杀格/...). Its value IS the stable
// display text, so it formats and serializes as text without Decode.
type GeJuKind string

const (
	geJuUnknown   GeJuKind = ""
	geJuZhengGuan GeJuKind = "正官格"
	geJuQiSha     GeJuKind = "七杀格"
	geJuZhengCai  GeJuKind = "正财格"
	geJuPianCai   GeJuKind = "偏财格"
	geJuZhengYin  GeJuKind = "正印格"
	geJuPianYin   GeJuKind = "偏印格"
	geJuShiShen   GeJuKind = "食神格"
	geJuShangGuan GeJuKind = "伤官格"
	geJuZa        GeJuKind = "杂格"
	geJuYueRen    GeJuKind = "月刃格"
	geJuJianLu    GeJuKind = "建禄格"
)

func DecodeGeJuKind(g GeJuKind) string { return string(g) }
func EncodeGeJuKind(name string) (GeJuKind, bool) {
	k := GeJuKind(name)
	if k == geJuUnknown {
		return geJuUnknown, false
	}
	return k, true
}

// CongGeKind is one of the 从格 variants. Its value IS the stable text.
type CongGeKind string

const (
	congGeUnknown CongGeKind = ""
	congGeWang    CongGeKind = "从旺格"
	congGeSha     CongGeKind = "从杀格"
	congGeCai     CongGeKind = "从财格"
	congGeEr      CongGeKind = "从儿格"
	congGeJiaWang CongGeKind = "假从旺格"
	congGeJiaSha  CongGeKind = "假从杀格"
	congGeJiaCai  CongGeKind = "假从财格"
	congGeJiaEr   CongGeKind = "假从儿格"
)

func DecodeCongGeKind(c CongGeKind) string { return string(c) }
func EncodeCongGeKind(name string) (CongGeKind, bool) {
	k := CongGeKind(name)
	if k == congGeUnknown {
		return congGeUnknown, false
	}
	return k, true
}

// CandidateWuxing is a 扶抑候选五行 symbol. Its value IS the output text;
// the closed set is 木/火/土/金/水 plus the explicit absence marker 无候选.
type CandidateWuxing string

const (
	NoCandidateWuxing CandidateWuxing = "无候选"
	candidateNone                     = NoCandidateWuxing
	candidateWood     CandidateWuxing = "木"
	candidateFire     CandidateWuxing = "火"
	candidateEarth    CandidateWuxing = "土"
	candidateMetal    CandidateWuxing = "金"
	candidateWater    CandidateWuxing = "水"
)

func DecodeCandidateWuxing(c CandidateWuxing) string { return string(c) }

func candidateFromWuxing(w ganzhi.Wuxing) CandidateWuxing {
	switch w {
	case ganzhi.WxMu:
		return candidateWood
	case ganzhi.WxHuo:
		return candidateFire
	case ganzhi.WxTu:
		return candidateEarth
	case ganzhi.WxJin:
		return candidateMetal
	case ganzhi.WxShui:
		return candidateWater
	default:
		return candidateNone
	}
}

// TenGodSource is where a ten-god occurrence comes from: 透干/本气/中气/余气.
// The output-only patternSource (main_qi_not_transparent) stays a plain string.
type TenGodSource string

// TripleRelationKind is a 三合/三会/半合 relation kind.
type TripleRelationKind string

const (
	tripleSanHe     TripleRelationKind = "三合"
	tripleSanHui    TripleRelationKind = "三会"
	tripleSanHePart TripleRelationKind = "半合"
)

func (k TripleRelationKind) suffix() string {
	switch k {
	case tripleSanHe:
		return "局"
	case tripleSanHui:
		return "方"
	default:
		return "半合局"
	}
}
