package qimen

import (
	"encoding/json"
	"fmt"
)

// ParsePalaceIndex converts a Chinese gong name to a GongIndex value.
func ParsePalaceIndex(s string) (GongIndex, error) {
	for i := 1; i <= 9; i++ {
		if gongNames[i] == s {
			return GongIndex(i), nil
		}
	}
	return 0, fmt.Errorf("unknown gong: %q", s)
}

// ParseStarIndex converts a Chinese star name to a StarIndex value.
func ParseStarIndex(s string) (StarIndex, error) {
	for i := 1; i <= 9; i++ {
		if starNames[i] == s {
			return StarIndex(i), nil
		}
	}
	return 0, fmt.Errorf("unknown star: %q", s)
}

// ParseDoorIndex converts a Chinese door name to a DoorIndex value.
func ParseDoorIndex(s string) (DoorIndex, error) {
	for i := 1; i <= len(doorNames)-1; i++ {
		if doorNames[i]+"门" == s {
			return DoorIndex(i), nil
		}
	}
	return 0, fmt.Errorf("unknown door: %q", s)
}

// ParseSpiritIndex converts a Chinese spirit name to a SpiritIndex value.
func ParseSpiritIndex(s string) (SpiritIndex, error) {
	if s == "太常" {
		return SpiritTaiChang, nil
	}
	for i := 1; i <= 8; i++ {
		if SpiritIndex(i).YangName() == s || SpiritIndex(i).YinName() == s {
			return SpiritIndex(i), nil
		}
	}
	return 0, fmt.Errorf("unknown spirit: %q", s)
}

// WuxingRelation is a qimen-domain symbol for a five-element relation between
// two members (a star and a gong, or 日 and 时). The relation key is stable and
// has a reverse mapping to its catalog identifier.
type WuxingRelation uint8

const (
	wrSame             WuxingRelation = iota // same
	wrGongGeneratesXing                       // gong_generates_xing
	wrXingGeneratesGong                       // xing_generates_gong
	wrXingControlsGong                        // xing_controls_gong
	wrGongControlsXing                        // gong_controls_xing
	wrRiGeneratesShi                          // ri_generates_shi
	wrShiGeneratesRi                          // shi_generates_ri
	wrRiControlsShi                           // ri_controls_shi
	wrShiControlsRi                           // shi_controls_ri
)

var wuxingRelationNames = [...]string{
	"same", "gong_generates_xing", "xing_generates_gong",
	"xing_controls_gong", "gong_controls_xing",
	"ri_generates_shi", "shi_generates_ri",
	"ri_controls_shi", "shi_controls_ri",
}

func (r WuxingRelation) String() string {
	if int(r) < len(wuxingRelationNames) {
		return wuxingRelationNames[r]
	}
	return ""
}

func ParseWuxingRelation(name string) (WuxingRelation, bool) {
	for i, n := range wuxingRelationNames {
		if n == name {
			return WuxingRelation(i), true
		}
	}
	return 0, false
}

// Yuan is the qimen-domain symbol for the upper/middle/lower 元 of a
// solar-term or fu-tou cycle.
type Yuan uint8

const (
	YuanShang Yuan = iota // 上元
	YuanZhong             // 中元
	YuanXia               // 下元
)

var yuanLabels = [...]string{"上元", "中元", "下元"}

func (y Yuan) String() string {
	if int(y) < len(yuanLabels) {
		return yuanLabels[y]
	}
	return ""
}

func ParseYuan(name string) (Yuan, bool) {
	for i, n := range yuanLabels {
		if n == name {
			return Yuan(i), true
		}
	}
	return 0, false
}

// ZhiRunState is the qimen-domain symbol for the 置闰 state of a dingju.
type ZhiRunState uint8

const (
	zsNone ZhiRunState = iota // 无置闰状态（非置闰法）
	zsExact                   // exact 正授
	zsAhead                   // ahead 超神
	zsBehind                  // behind 接气
	zsRepeat                  // repeat 置闰
)

var zhiRunStateNamesByKey = [...]string{"", "exact", "ahead", "behind", "repeat"}
var zhiRunStateLabels = [...]string{"", "正授", "超神", "接气", "置闰"}

func (s ZhiRunState) String() string {
	if int(s) < len(zhiRunStateLabels) {
		return zhiRunStateLabels[s]
	}
	return ""
}

func (s ZhiRunState) Key() string {
	if int(s) < len(zhiRunStateNamesByKey) {
		return zhiRunStateNamesByKey[s]
	}
	return ""
}

func ParseZhiRunState(key string) (ZhiRunState, bool) {
	for i, n := range zhiRunStateNamesByKey {
		if n == key && key != "" {
			return ZhiRunState(i), true
		}
	}
	return 0, false
}

// DunSource is the qimen-domain symbol for how 遁阴阳 is selected in the
// twelve-minute quarter method.
type DunSource uint8

const (
	dsHourBranch DunSource = iota // hour_branch
	dsSolarTerm                    // solar_term
)

var dunSourceNames = [...]string{"hour_branch", "solar_term"}

func (d DunSource) String() string {
	if int(d) < len(dunSourceNames) {
		return dunSourceNames[d]
	}
	return ""
}

func (d DunSource) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

func ParseDunSource(name string) (DunSource, bool) {
	for i, n := range dunSourceNames {
		if n == name {
			return DunSource(i), true
		}
	}
	return 0, false
}

// HourBoundary is the qimen-domain symbol for the hour-start boundary used by
// the twelve-minute quarter method.
type HourBoundary uint8

const (
	hbLateZi  HourBoundary = iota // late_zi
	hbZiZheng                     // zi_zheng
)

var hourBoundaryNames = [...]string{"late_zi", "zi_zheng"}

func (h HourBoundary) String() string {
	if int(h) < len(hourBoundaryNames) {
		return hourBoundaryNames[h]
	}
	return ""
}

func (h HourBoundary) MarshalJSON() ([]byte, error) { return json.Marshal(h.String()) }

func ParseHourBoundary(name string) (HourBoundary, bool) {
	for i, n := range hourBoundaryNames {
		if n == name {
			return HourBoundary(i), true
		}
	}
	return 0, false
}

// DingjuMethod selects a stable method for determining the Qimen ju (局数).
type DingjuMethod string

const (
	DingjuChaiBu  DingjuMethod = "chaibu"
	DingjuZhiRun  DingjuMethod = "zhirun"
	DingjuMaoShan DingjuMethod = "maoshan"
	DingjuNone    DingjuMethod = "none"
)

// ParseDingjuMethod validates a dingju-method identifier.
func ParseDingjuMethod(value string) (DingjuMethod, error) {
	switch DingjuMethod(value) {
	case DingjuChaiBu, DingjuZhiRun, DingjuMaoShan, DingjuNone:
		return DingjuMethod(value), nil
	default:
		return "", fmt.Errorf(
			"unknown dingju method %q: must be chaibu, zhirun, maoshan or none", value,
		)
	}
}

func (m DingjuMethod) String() string { return string(m) }

// QuarterRule selects a named Kejia time-division rule.
type QuarterRule string

const (
	QuarterTenMinuteSanYuan        QuarterRule = "ten_minute_sanyuan"
	QuarterTwelveMinuteTenDivision QuarterRule = "twelve_minute_ten_division"
)

// ParseQuarterRule validates a Kejia rule identifier.
func ParseQuarterRule(value string) (QuarterRule, error) {
	switch QuarterRule(value) {
	case QuarterTenMinuteSanYuan, QuarterTwelveMinuteTenDivision:
		return QuarterRule(value), nil
	default:
		return "", fmt.Errorf(
			"unknown quarter rule %q: must be ten_minute_sanyuan or twelve_minute_ten_division", value,
		)
	}
}

func (r QuarterRule) String() string { return string(r) }

// Scope selects the chart's time layer and duty pillar.
type Scope string

const (
	ScopeHour    Scope = "hour"
	ScopeDay     Scope = "day"
	ScopeQuarter Scope = "quarter"
	ScopeMonth   Scope = "month"
	ScopeYear    Scope = "year"
)

// ParseScope validates a time-scope identifier.
func ParseScope(value string) (Scope, error) {
	switch Scope(value) {
	case ScopeHour, ScopeDay, ScopeQuarter, ScopeMonth, ScopeYear:
		return Scope(value), nil
	default:
		return "", fmt.Errorf("unknown scope %q: must be hour, day, quarter, month or year", value)
	}
}

func (s Scope) String() string { return string(s) }

// School selects a plate geometry and its spirit convention.
type School string

const (
	SchoolZhuanPan     School = "zhuanpan"
	SchoolLuoShuFeiPan School = "luoshu_feipan"
	SchoolMingFaFeiPan School = "mingfa_feipan"
	SchoolJinhanYuJing School = "jinhan_yujing"
)

// ParseSchool validates a plate-school identifier.
func ParseSchool(value string) (School, error) {
	switch School(value) {
	case SchoolZhuanPan, SchoolLuoShuFeiPan, SchoolMingFaFeiPan, SchoolJinhanYuJing:
		return School(value), nil
	default:
		return "", fmt.Errorf("unknown school %q: must be zhuanpan, luoshu_feipan, mingfa_feipan or jinhan_yujing", value)
	}
}

func (s School) String() string { return string(s) }

// GongIndex is a 洛书九宫 index:
// 1=坎, 2=坤, 3=震, 4=巽, 5=中, 6=乾, 7=兑, 8=艮, 9=离.
type GongIndex int

const (
	GongKan GongIndex = 1 + iota
	GongKun
	GongZhen
	GongXun
	GongZhong
	GongQian
	GongDui
	GongGen
	GongLi
)

var gongNames = [10]string{"", "坎", "坤", "震", "巽", "中", "乾", "兑", "艮", "离"}

func (p GongIndex) String() string {
	if p >= 1 && p <= 9 {
		return gongNames[p]
	}
	return "?"
}

func (p GongIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *GongIndex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("gong must be a string, got %s", string(data))
	}
	for i, name := range gongNames {
		if i > 0 && name == s {
			*p = GongIndex(i)
			return nil
		}
	}
	return fmt.Errorf("unknown gong: %q", s)
}

// StarIndex represents one of the 九星 (nine stars).
type StarIndex int

const (
	StarTianPeng StarIndex = 1 + iota
	StarTianRui
	StarTianChong
	StarTianFu
	StarTianQin
	StarTianXin
	StarTianZhu
	StarTianRen
	StarTianYing
)

var starNames = [10]string{"", "天蓬", "天芮", "天冲", "天辅", "天禽", "天心", "天柱", "天任", "天英"}

func (s StarIndex) String() string {
	if s >= 1 && s <= 9 {
		return starNames[s]
	}
	return "?"
}

func (s StarIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *StarIndex) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("star must be a string (e.g. \"天蓬\"), got %s", string(data))
	}
	for i, n := range starNames {
		if i > 0 && n == name {
			*s = StarIndex(i)
			return nil
		}
	}
	return fmt.Errorf("unknown star: %q", name)
}

// DoorIndex represents the eight doors and the Ming Fa center door.
type DoorIndex int

const (
	DoorXiu DoorIndex = 1 + iota
	DoorSheng
	DoorShang
	DoorDu
	DoorJing
	DoorSi
	DoorJingMen
	DoorKai
	DoorZhong
)

var doorNames = [10]string{"", "休", "生", "伤", "杜", "景", "死", "惊", "开", "中"}

func (d DoorIndex) String() string {
	if d >= 1 && d <= 9 {
		return doorNames[d] + "门"
	}
	return "?"
}

func (d DoorIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *DoorIndex) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("door must be a string (e.g. \"休门\"), got %s", string(data))
	}
	for i, n := range doorNames {
		if i > 0 && n+"门" == name {
			*d = DoorIndex(i)
			return nil
		}
	}
	return fmt.Errorf("unknown door: %q", name)
}

// SpiritIndex represents the rotate-plate 八神 and the fly-plate ninth spirit.
type SpiritIndex int

const (
	SpiritZhiFu SpiritIndex = 1 + iota
	SpiritTengShe
	SpiritTaiYin
	SpiritLiuHe
	SpiritGouChen // 阳遁=勾陈, 阴遁=白虎
	SpiritZhuQue  // 阳遁=朱雀, 阴遁=玄武
	SpiritJiuDi
	SpiritJiuTian
	SpiritTaiChang
)

// YangSpiritNames returns the spirit name for 阳遁.
func (s SpiritIndex) YangName() string {
	if s >= 1 && int(s) < len(spiritYangNames) {
		if spiritYangNames[s] != "" {
			return spiritYangNames[s]
		}
	}
	if s >= 1 && int(s) < len(flySpiritNames) && flySpiritNames[s] != "" {
		return flySpiritNames[s]
	}
	return "?"
}

// YinName returns the spirit name for 阴遁.
func (s SpiritIndex) YinName() string {
	if s >= 1 && int(s) < len(spiritYinNames) {
		if spiritYinNames[s] != "" {
			return spiritYinNames[s]
		}
	}
	if s >= 1 && int(s) < len(flySpiritNames) && flySpiritNames[s] != "" {
		return flySpiritNames[s]
	}
	return "?"
}

// UnmarshalJSON accepts a spirit name (值符/螣蛇/…).
func (s *SpiritIndex) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("spirit must be a string, got %s", string(data))
	}
	for i := 1; i < len(spiritYangNames); i++ {
		if SpiritIndex(i).YangName() == name || SpiritIndex(i).YinName() == name {
			*s = SpiritIndex(i)
			return nil
		}
	}
	return fmt.Errorf("unknown spirit: %q", name)
}

// DutyStarPosition is the qimen-domain symbol for a duty star's required
// position relative to its home palace.
type DutyStarPosition uint8

const (
	dspNone     DutyStarPosition = iota // 未指定
	dspHome                             // home 本位
	dspOpposite                         // opposite 对冲
)

func ParseDutyStarPosition(name string) (DutyStarPosition, bool) {
	switch name {
	case "home":
		return dspHome, true
	case "opposite":
		return dspOpposite, true
	default:
		return 0, false
	}
}

func (d DutyStarPosition) String() string {
	switch d {
	case dspHome:
		return "home"
	case dspOpposite:
		return "opposite"
	default:
		return ""
	}
}

// PillarLevel is the qimen-domain symbol for a four-pillar level.
type PillarLevel uint8

const (
	levelYear  PillarLevel = iota // year
	levelMonth                    // month
	levelDay                      // day
	levelHour                     // hour
)

var pillarLevelNames = [...]string{"year", "month", "day", "hour"}

func (l PillarLevel) String() string {
	if int(l) < len(pillarLevelNames) {
		return pillarLevelNames[l]
	}
	return ""
}

func (l PillarLevel) MarshalJSON() ([]byte, error) { return json.Marshal(l.String()) }

func ParsePillarLevel(name string) (PillarLevel, bool) {
	for i, n := range pillarLevelNames {
		if n == name {
			return PillarLevel(i), true
		}
	}
	return 0, false
}
