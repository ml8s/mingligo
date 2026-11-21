package fengshui

// YuanLong is one of the three元龙 in the 24-mountain system (三元龙).
// Its value IS the stable display text.
type YuanLong string

const (
	YuanLongTian YuanLong = "天元龙"
	YuanLongDi   YuanLong = "地元龙"
	YuanLongRen  YuanLong = "人元龙"
)

var yuanLongSet = map[YuanLong]bool{
	YuanLongTian: true, YuanLongDi: true, YuanLongRen: true,
}

// EncodeYuanLong converts a display name to its YuanLong symbol (boundary helper).
func EncodeYuanLong(s string) (YuanLong, bool) {
	v := YuanLong(s)
	return v, yuanLongSet[v]
}
