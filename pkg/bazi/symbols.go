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
	sanQiNone SanQiKind = iota
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
