package tianwen

import (
	"fmt"
	"math"
)

// JieQi is the canonical solar-term symbol, indexed from 冬至(0) through
// 大雪(23) — the same order as AllSolarTerms and qimen's canonical index.
// It is defined once here in tianwen (the time domain) and shared globally.
type JieQi int

const (
	JQDongZhi     JieQi = iota // 冬至 0  (270°)
	JQXiaoHan                  // 小寒 1  (285°)
	JQDaHan                    // 大寒 2  (300°)
	JQLiChun                   // 立春 3  (315°)
	JQYuShui                   // 雨水 4  (330°)
	JQJingZhe                  // 惊蛰 5  (345°)
	JQChunFen                  // 春分 6  (0°)
	JQQingMing                 // 清明 7  (15°)
	JQGuYu                     // 谷雨 8  (30°)
	JQLiXia                    // 立夏 9  (45°)
	JQXiaoMan                  // 小满 10 (60°)
	JQMangZhong                // 芒种 11 (75°)
	JQXiaZhi                   // 夏至 12 (90°)
	JQXiaoShu                  // 小暑 13 (105°)
	JQDaShu                    // 大暑 14 (120°)
	JQLiQiu                    // 立秋 15 (135°)
	JQChuShu                   // 处暑 16 (150°)
	JQBaiLu                    // 白露 17 (165°)
	JQQiuFen                   // 秋分 18 (180°)
	JQHanLu                    // 寒露 19 (195°)
	JQShuangJiang              // 霜降 20 (210°)
	JQLiDong                   // 立冬 21 (225°)
	JQXiaoXue                  // 小雪 22 (240°)
	JQDaXue                    // 大雪 23 (255°)
)

var jieQiNamesFromDongZhi = [24]string{
	"冬至", "小寒", "大寒", "立春", "雨水", "惊蛰", "春分", "清明",
	"谷雨", "立夏", "小满", "芒种", "夏至", "小暑", "大暑", "立秋",
	"处暑", "白露", "秋分", "寒露", "霜降", "立冬", "小雪", "大雪",
}

// DecodeJieQi decodes a JieQi symbol into its display text.
func (j JieQi) String() string { return DecodeJieQi(j) }

func DecodeJieQi(j JieQi) string {
	if int(j) >= 0 && int(j) < len(jieQiNamesFromDongZhi) {
		return jieQiNamesFromDongZhi[j]
	}
	return ""
}

func (j JieQi) Longitude() float64 {
	return JieQiLongitudes[(int(j)+21)%24]
}

// JieQiOfLongitude maps a solar longitude (multiple of 15°) to its canonical
// JieQi (冬至=0). Returns ok=false for invalid longitudes.
func JieQiOfLongitude(longitude float64) (JieQi, bool) {
	if longitude < 0 || longitude >= 360 || math.Mod(longitude, 15) != 0 {
		return 0, false
	}
	return JieQi(int(math.Mod(longitude-270+360, 360) / 15)), true
}

// EncodeJieQi maps a Chinese solar-term name back to its canonical symbol.
func EncodeJieQi(name string) (JieQi, bool) {
	for i, n := range jieQiNamesFromDongZhi {
		if n == name {
			return JieQi(i), true
		}
	}
	return 0, false
}

// LunarBasis selects the date frame used to project an instant into the
// Chinese lunisolar calendar.
type LunarBasis string

const (
	// LunarBasisOfficial uses the Beijing civil date and the UTC+8 calendar.
	LunarBasisOfficial LunarBasis = "official"
	// LunarBasisLocalLabel uses the caller's local date label with the UTC+8
	// calendar table. This is the ziwei-compatible school convention.
	LunarBasisLocalLabel LunarBasis = "local_label"
	// LunarBasisLocalAstronomical uses the caller's local date and a calendar
	// whose new moons and 中气 are assigned in the caller's zone.
	LunarBasisLocalAstronomical LunarBasis = "local_astronomical"
)

// EncodeLunarBasis validates a basis name. The empty string is the default
// local_astronomical basis, which is the globally aligned projection.
func EncodeLunarBasis(basis string) (LunarBasis, error) {
	switch LunarBasis(basis) {
	case "", LunarBasisLocalAstronomical:
		return LunarBasisLocalAstronomical, nil
	case LunarBasisOfficial, LunarBasisLocalLabel:
		return LunarBasis(basis), nil
	default:
		return "", fmt.Errorf(
			"lunar basis must be official, local_label, or local_astronomical, got %q", basis,
		)
	}
}
