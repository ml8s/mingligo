package ziwei

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// Type aliases from ganzhi.
type (
	Gan    = ganzhi.Gan
	Zhi    = ganzhi.Zhi
	Wuxing = ganzhi.Wuxing
)

const (
	Male   = ganzhi.Male
	Female = ganzhi.Female
)

// juShuFromWuxing converts a five-element to its bureau number.
func juShuFromWuxing(w Wuxing) JuShu {
	switch w {
	case ganzhi.WxShui:
		return JuWater
	case ganzhi.WxMu:
		return JuWood
	case ganzhi.WxJin:
		return JuMetal
	case ganzhi.WxTu:
		return JuEarth
	case ganzhi.WxHuo:
		return JuFire
	}
	return 0
}

// gong holds all computed data for one gong.
type Gong struct {
	Index        GongIndex  `json:"-"`
	Name         string     `json:"ming"`
	Gan          Gan        `json:"gan"`
	Zhi          Zhi        `json:"zhi"`
	IsBodyPalace bool       `json:"shen_gong,omitempty"`
	IsYuanGong   bool       `json:"lai_yin_gong,omitempty"`
	Stars        []StarInfo `json:"xing_yao"`
	Ages         []int      `json:"xiao_xian,omitempty"`
	ChangSheng   string     `json:"chang_sheng,omitempty"`
	BoShi        string     `json:"bo_shi,omitempty"`
	JiangQian    string     `json:"jiang_qian,omitempty"`
	SuiQian      string     `json:"sui_qian,omitempty"`
	ZaYao        []string   `json:"za_yao,omitempty"`
}

// StarInfo is one star entry in a gong.
type StarInfo struct {
	Star       StarIndex `json:"-"`
	Name       string    `json:"ming"`
	IsMajor    bool      `json:"-"`
	Lei        []string  `json:"lei,omitempty"`
	SiHua      SiHua     `json:"si_hua,omitempty"`   // 禄/权/科/忌 or empty
	Brightness string    `json:"liang_du,omitempty"` // "庙"/"旺"/"利"/"平"/"陷"
}

// KongGongInfo 描述一个无主星宫位及其借星（对宫主星）。
type KongGongInfo struct {
	KongGong GongIndex `json:"kong_gong"`
	JieGong  GongIndex `json:"jie_gong"`
}

// Chart holds the complete ziwei chart.
type Chart struct {
	KouJing  ChartSchool `json:"kou_jing"`
	GongWei  [12]Gong    `json:"shi_er_gong"`
	MingGong GongIndex   `json:"-"`
	ShenGong GongIndex   `json:"-"`
	JuShu    JuShu       `json:"ju_shu"`
	// JuShuName is the human-readable bureau name; JuShu.code is machine-facing.
	JuShuName string      `json:"ju_shu_name,omitempty"`
	ZiweiPos  GongIndex   `json:"-"`
	SiHua     siHuaResult `json:"-"`

	// 空宫借星：只记录无主星宫和借星来源宫；星曜仍由来源宫唯一承载。
	KongGong        []KongGongInfo `json:"kong_gong,omitempty"`
	NianGan         Gan            `json:"nian_gan"`
	NianZhi         Zhi            `json:"nian_zhi,omitempty"`
	ShiZhi          Zhi            `json:"shi_zhi"`
	BirthYear       int            `json:"birth_year"`
	LunarMonth      int            `json:"lunar_month,omitempty"`
	LunarDay        int            `json:"lunar_day,omitempty"`
	BirthLunarMonth int            `json:"birth_lunar_month,omitempty"`
	BirthIsLeap     bool           `json:"birth_is_leap,omitempty"`
	Gender          ganzhi.Gender  `json:"-"`
	MingZhu         string         `json:"ming_zhu,omitempty"`
	ShenZhu         string         `json:"shen_zhu,omitempty"`
	Patterns        []Pattern      `json:"ge_ju,omitempty"`
	SanFangSiZheng  SanFangSiZheng `json:"san_fang_si_zheng"`
	DaXian          []DaXianStep   `json:"da_xian"`
}

// SanFangSiZheng is the palace relation, not a copied star view. Consumers
// read star facts from the referenced palaces in ShiErGong.
type SanFangSiZheng struct {
	MingGong []GongIndex `json:"ming_gong"`
}

// ChartSchool makes the leap-month convention explicit. The engine currently
// follows the common iztro-compatible rule rather than silently claiming the
// Quanshu year-start rule.
type ChartSchool struct {
	LeapMonth string `json:"leap_month"`
	Source    string `json:"source"`
}

// siHuaResult maps star → transformation.
type siHuaResult map[StarIndex]SiHua

// DaXianStep records one 10-year da-xian segment.
type DaXianStep struct {
	QiSui     int       `json:"qi_sui"`
	ZhiSui    int       `json:"zhi_sui"`
	StartYear int       `json:"start_year"` // 本限起始公历年（birthYear + QiSui − 1）
	EndYear   int       `json:"end_year"`   // 本限结束公历年（StartYear + 9）
	MingGong  GongIndex `json:"ming_gong"`
}

// LiuNianJie locates the flow lunar year: the observer's standard clock and
// civil selection date, plus the lunar-year boundary. Values are readable.
type LiuNianJie struct {
	StandardTime   string `json:"standard_time"`
	StartDate      string `json:"start_date"`
	LunarYear      int    `json:"lunar_year"`
	LunarYearStart string `json:"lunar_year_start"`
	LunarYearEnd   string `json:"lunar_year_end"`
	BoundaryDay    bool   `json:"boundary_day"`
}

// LiuNian is the annual fate analysis.
type LiuNian struct {
	MingGong     GongIndex            `json:"ming_gong"`
	MingGongName string               `json:"ming_gong_name"` // 流年命宫名
	Zhi          Zhi                  `json:"zhi"`            // 流年地支
	SiHua        []FlowTransformation `json:"si_hua"`
	GongWei      [12]FlowPalace       `json:"gong_wei"` // 流年盘（地支坐标 12 宫）
	NianJie      LiuNianJie           `json:"nian_jie"`
}
