package bazhai

// LifeGroup is the 命卦 group (东四命/西四命). Its value IS the stable text.
type LifeGroup string

const (
	LifeGroupEast LifeGroup = "东四命"
	LifeGroupWest LifeGroup = "西四命"
)

// TrigramGroup is the 卦 group (东四卦/西四卦) of a door/stove palace.
type TrigramGroup string

const (
	TrigramGroupEast TrigramGroup = "东四卦"
	TrigramGroupWest TrigramGroup = "西四卦"
)

// YouXing is one of the eight 游年 stars (八宅游星). Its value IS the stable text.
type YouXing string

const (
	YouXingShengQi YouXing = "生气"
	YouXingTianYi  YouXing = "天医"
	YouXingYanNian YouXing = "延年"
	YouXingFuWei   YouXing = "伏位"
	YouXingHuoHai  YouXing = "祸害"
	YouXingWuGui   YouXing = "五鬼"
	YouXingLiuSha  YouXing = "六煞"
	YouXingJueMing YouXing = "绝命"
)

// Rating is the 吉凶等级 of a 游星 (大吉/吉/小吉/平/凶/大凶). Its value IS the text.
type Rating string

const (
	RatingGreatLucky Rating = "大吉"
	RatingLucky      Rating = "吉"
	RatingSmallLucky Rating = "小吉"
	RatingPlain      Rating = "平"
	RatingBad        Rating = "凶"
	RatingGreatBad   Rating = "大凶"
)

// Match is whether a palace shares the 命卦 group (吉=同组, 凶=异组).
type Match string

const (
	MatchLucky Match = "吉"
	MatchBad   Match = "凶"
)
