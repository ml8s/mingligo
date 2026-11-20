package qimen

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

// TianPanSymbol binds a heaven-plate star to the gan carried from its home palace.
type TianPanSymbol struct {
	Gan  ganzhi.Gan `json:"gan"`
	Star StarIndex  `json:"xing"`
}

type PalaceIdentity struct {
	Name      string `json:"name"`
	Luoshu    int    `json:"luoshu"`
	RingIndex *int   `json:"ring_index,omitempty"`
}

type BranchPalace struct {
	Branch ganzhi.Zhi `json:"zhi"`
	Gong   GongIndex  `json:"gong"`
}

// Gong holds all layers of information for one 宫。
type Gong struct {
	AnGan        ganzhi.Gan      `json:"an_gan"`
	HiddenPillar string          `json:"an_gan_zhi,omitempty"`
	DiPanGan     ganzhi.Gan      `json:"di_pan_gan"`
	Gong         PalaceIdentity  `json:"gong"`
	Door         DoorIndex       `json:"men,omitempty"`
	DoorSet      bool            `json:"men_present"`
	Shen         string          `json:"shen,omitempty"`
	SpiritSet    bool            `json:"shen_present"`
	TianPan      []TianPanSymbol `json:"tian_pan,omitempty"`
	TianPanGan   *ganzhi.Gan     `json:"tian_pan_gan,omitempty"`
	Spirit       SpiritIndex     `json:"-"`
}

// board is the complete 奇门遁甲排盘。
type Board struct {
	GongWei        [9]Gong         `json:"gong_wei"`
	Jushu          int             `json:"jushu"`
	KongWang       [2]BranchPalace `json:"kong_wang"`
	MaXing         BranchPalace    `json:"ma_xing"`
	NianGan        ganzhi.Gan      `json:"nian_gan"`
	NianZhi        ganzhi.Zhi      `json:"nian_zhi"`
	RiGan          ganzhi.Gan      `json:"ri_gan"`
	RiZhi          ganzhi.Zhi      `json:"ri_zhi"`
	HourGan        ganzhi.Gan      `json:"shi_gan"`
	HourZhi        ganzhi.Zhi      `json:"shi_zhi"`
	WuBuYuShi      bool            `json:"wu_bu_yu_shi"`
	YinDun         bool            `json:"yin_dun"`
	YueGan         ganzhi.Gan      `json:"yue_gan"`
	YueZhi         ganzhi.Zhi      `json:"yue_zhi"`
	DutyStar       StarIndex       `json:"zhi_fu_xing"`
	DutyDoor       DoorIndex       `json:"zhi_shi_men"`
	School         School          `json:"-"`
	LeadGan        ganzhi.Gan      `json:"-"`
	LeadZhi        ganzhi.Zhi      `json:"-"`
	DutyStarPalace GongIndex       `json:"-"`
	DutyDoorPalace GongIndex       `json:"-"`
}

// duty holds the value符 star and value使 door.
type duty struct {
	Star   StarIndex
	Door   DoorIndex
	Palace GongIndex
}

// juShu holds the result of dingju determination.
type juShu struct {
	Method      Method
	Number      int
	YinDun      bool
	Yuan        Yuan
	JieQi       tianwen.JieQi
	YongJuJieQi string
	ZhiRunState ZhiRunState
	Quarter     *quarterInfo
}

type TianQinRule struct {
	Star                  string `json:"xing"`
	HomeGong              string `json:"home_gong"`
	LodgingGong           string `json:"lodging_gong"`
	Follows               string `json:"follows"`
	CarriesCenterEarthGan bool   `json:"carries_center_earth_gan"`
	LodgingRule           string `json:"lodging_rule"`
}

type ChartMethod struct {
	Scope            string        `json:"scope"`
	School           string        `json:"school"`
	DingjuMethod     string        `json:"dingju_method,omitempty"`
	QuarterRule      string        `json:"quarter_rule,omitempty"`
	BaseDingjuMethod string        `json:"base_dingju_method,omitempty"`
	MethodSource     string        `json:"method_source"`
	DayBoundary      string        `json:"day_boundary"`
	SpiritMode       string        `json:"shen_mode"`
	Geometry         string        `json:"geometry"`
	StarMode         string        `json:"xing_mode"`
	DoorMode         string        `json:"men_mode"`
	SpiritFlight     string        `json:"shen_flight"`
	TianQin          *TianQinRule  `json:"tian_qin,omitempty"`
	JieQi            string        `json:"jie_qi"`
	YongJuJieQi      string        `json:"yong_ju_jie_qi"`
	Yuan             string        `json:"yuan"`
	ZhiRunState      string        `json:"zhi_run_state,omitempty"`
	LeadPillar       PillarInfo    `json:"lead_pillar"`
	LeadXunShou      PillarInfo    `json:"lead_xun_shou"`
	LeadXunShouGong  GongIndex     `json:"lead_xun_shou_gong"`
	DayFuTou         PillarInfo    `json:"day_fu_tou"`
	Quarter          *ChartQuarter `json:"quarter,omitempty"`
}

type PillarInfo struct {
	Name  string `json:"name"`
	Gan   string `json:"gan"`
	Zhi   string `json:"zhi"`
	LiuYi string `json:"liu_yi,omitempty"`
}

// GanInteraction represents a 十干克应 between earth and heaven gan.
type GanInteraction struct {
	Gong       GongIndex  `json:"gong"`
	DiPanGan   ganzhi.Gan `json:"di_pan_gan"`
	TianPanGan ganzhi.Gan `json:"tian_pan_gan"`
	Name       string     `json:"name"`
	Auspicious bool       `json:"auspicious"`
}

// MenInteraction represents an 八门克应 for a door in a specific gong.
type MenInteraction struct {
	Door DoorIndex `json:"men,omitempty"`
	Gong GongIndex `json:"gong,omitempty"`
	Name string    `json:"name"`
}

// Pattern represents a detected 格局 in the board.
type Pattern struct {
	Name       string      `json:"name"`
	YiJu       string      `json:"basis,omitempty"`
	Auspicious bool        `json:"auspicious"`
	GongWei    []GongIndex `json:"gong_wei,omitempty"`
}

type patternRule struct {
	Name             string
	YiJu             string
	Auspicious       bool
	RequiresDutyDoor bool
	DutyStarPosition DutyStarPosition
	Conditions       []patternCondition
}

type patternCondition struct {
	HeavenGan []ganzhi.Gan
	EarthGan  []ganzhi.Gan
	Doors     []DoorIndex
	Spirits   []SpiritIndex
	Palaces   []GongIndex
}
