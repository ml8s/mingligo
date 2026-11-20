package bazi

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
