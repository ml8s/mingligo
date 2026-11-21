package ziwei

import "encoding/json"

// Domain-owned output views of the ziwei.liupan series. jsonschema tags are
// inert metadata read only when a consumer reflects the type for its schema.

type LiuNianOutputView struct {
	MingGong   GongIndex                `json:"ming_gong" jsonschema:"命宫。"`
	LiuMingZhi Zhi                      `json:"liu_ming_zhi" jsonschema:"流命所在地支（与 is_liu_ming 宫一致）。"`
	Gan        Gan                      `json:"period_gan" jsonschema:"流期天干（与 zhi 构成流期干支；非流宫宫干）。"`
	Zhi        Zhi                      `json:"zhi" jsonschema:"流年地支"`
	SiHua      []FlowTransformationView `json:"si_hua" jsonschema:"四化事实：星、化、落宫（落宫为本命宫序坐标）。"`
	GongWei    []FlowPalaceView         `json:"gong_wei" jsonschema:"流年十二宫盘。"`
	NianJie    LiuNianResolvedView      `json:"nian_jie" jsonschema:"流年农历年界。"`
}

type LiuNianResolvedView struct {
	LiuNianStandardTime string `json:"standard_time" jsonschema:"流年选择器时区标准时间（无 DST）。"`
	LiuNianStartDate    string `json:"start_date" jsonschema:"按选定农历基准取的流年选择日期。"`
	LunarYear           int    `json:"lunar_year" jsonschema:"该日期所在的农历年。"`
	LunarYearStart      string `json:"lunar_year_start" jsonschema:"正月初一公历日（含）。"`
	LunarYearEnd        string `json:"lunar_year_end" jsonschema:"下一正月初一公历日（不含）。"`
	BoundaryDay         bool   `json:"boundary_day" jsonschema:"输入日期是否为正月初一边界日（归入新年）。"`
}

type FlowTransformationView struct {
	Star   StarIndex `json:"xing" jsonschema:"被化的本命星。"`
	Hua    SiHua     `json:"hua" jsonschema:"化禄/化权/化科/化忌。"`
	Palace GongIndex `json:"gong" jsonschema:"四化落宫：本命宫序坐标（与本命盘 GongWei 对应，非流盘宫名）。"`
}

type FlowPalaceView struct {
	Zhi       Zhi       `json:"zhi" jsonschema:"地支。"`
	Gan       Gan       `json:"gan" jsonschema:"宫干（五虎遁）。"`
	Name      GongIndex `json:"name" jsonschema:"流盘宫名（以流命为命宫，非本命宫序）。"`
	Stars     []string  `json:"xing" jsonschema:"流曜。"`
	IsLiuMing bool      `json:"is_liu_ming" jsonschema:"是否流盘命宫。"`
}

type LunarDateView struct {
	Year  int  `json:"year" jsonschema:"农历年。"`
	Month int  `json:"month" jsonschema:"农历月。"`
	Day   int  `json:"day" jsonschema:"农历日。"`
	Leap  bool `json:"leap" jsonschema:"是否闰月。"`
}

type YueJieView struct {
	LiuYueStandardTime string                                 `json:"standard_time" jsonschema:"流月选择器时区标准时间（无 DST）。"`
	LiuYueStartDate    string                                 `json:"start_date" jsonschema:"按选定农历基准取的流月选择日期。"`
	Lunar              LunarDateView                          `json:"lunar" jsonschema:"该日期对应的农历日期。"`
	LunarMonthStart    string                                 `json:"lunar_month_start" jsonschema:"农历月初一公历日（含）。"`
	LunarMonthEnd      string                                 `json:"lunar_month_end" jsonschema:"下一农历月初一公历日（不含）。"`
	BoundaryDay        bool                                   `json:"boundary_day" jsonschema:"输入日期是否为农历月初一边界日（归入新月）。"`
	CalendarPeriod     LiuYueResolvedPeriodCalendarPeriodView `json:"calendar_period" jsonschema:"流月所在日历区间。"`
	FlowMonth          LiuYueResolvedPeriodFlowMonthView      `json:"flow_month" jsonschema:"流月农历年月。"`
}

type LiuYueOutputView struct {
	YueJie     YueJieView               `json:"yue_jie" jsonschema:"流月农历月界与月令。"`
	MingGong   GongIndex                `json:"ming_gong" jsonschema:"命宫。"`
	LiuMingZhi Zhi                      `json:"liu_ming_zhi" jsonschema:"流命所在地支（与 is_liu_ming 宫一致）。"`
	Gan        Gan                      `json:"period_gan" jsonschema:"流期天干（与 zhi 构成流期干支；非流宫宫干）。"`
	Zhi        Zhi                      `json:"zhi" jsonschema:"流月地支"`
	SiHua      []FlowTransformationView `json:"si_hua" jsonschema:"四化事实：星、化、落宫（落宫为本命宫序坐标）。"`
	GongWei    []FlowPalaceView         `json:"gong_wei" jsonschema:"流月十二宫盘。"`
}

type LiuYueResolvedPeriodCalendarPeriodView struct {
	Year     int    `json:"year" jsonschema:"年。"`
	Month    int    `json:"month" jsonschema:"月。"`
	Leap     bool   `json:"leap" jsonschema:"是否闰月。"`
	Half     string `json:"half" jsonschema:"农历半月；whole 为整月，闰月按 first / second 拆分。"`
	DayStart int    `json:"day_start" jsonschema:"流月日段起始农历日序。"`
	DayEnd   int    `json:"day_end" jsonschema:"流月日段结束农历日序。"`
}

type LiuYueResolvedPeriodFlowMonthView struct {
	Year  int  `json:"year" jsonschema:"年。"`
	Month int  `json:"month" jsonschema:"月。"`
	Leap  bool `json:"leap" jsonschema:"是否闰月。"`
}

type LiuRiJieView struct {
	LiuRiStandardTime string        `json:"standard_time" jsonschema:"流日选择器时区标准时间（无 DST）。"`
	LiuRiStandardDate string        `json:"standard_date" jsonschema:"流日选择器的时区标准日（无 DST）。"`
	Lunar             LunarDateView `json:"lunar" jsonschema:"该日期对应的农历日期。"`
}

type LiuRiOutputView struct {
	RiJie      LiuRiJieView             `json:"ri_jie" jsonschema:"流日日界。"`
	MingGong   GongIndex                `json:"ming_gong" jsonschema:"命宫。"`
	LiuMingZhi Zhi                      `json:"liu_ming_zhi" jsonschema:"流命所在地支（与 is_liu_ming 宫一致）。"`
	Gan        Gan                      `json:"period_gan" jsonschema:"流期天干（与 zhi 构成流期干支；非流宫宫干）。"`
	Zhi        Zhi                      `json:"zhi" jsonschema:"流日地支"`
	SiHua      []FlowTransformationView `json:"si_hua" jsonschema:"四化事实：星、化、落宫（落宫为本命宫序坐标）。"`
	GongWei    []FlowPalaceView         `json:"gong_wei" jsonschema:"流日十二宫盘。"`
}

type LiuShiChenView struct {
	LiuShiStandardTime string        `json:"standard_time" jsonschema:"流时选择器时区标准时间（无 DST）。"`
	LiuShiStandardDate string        `json:"standard_date" jsonschema:"流时选择器的时区标准日（无 DST）。"`
	Lunar              LunarDateView `json:"lunar" jsonschema:"该日期对应的农历日期。"`
	Shichen            Zhi           `json:"shichen" jsonschema:"真太阳时时辰；流时干支按当日日柱起五鼠遁（iztro 口径，晚子时不换日干），与八字 HourPillar 的晚子时次日口径不同。"`
}

type LiuShiOutputView struct {
	MingGong   GongIndex                `json:"ming_gong" jsonschema:"命宫。"`
	LiuMingZhi Zhi                      `json:"liu_ming_zhi" jsonschema:"流命所在地支（与 is_liu_ming 宫一致）。"`
	Gan        Gan                      `json:"period_gan" jsonschema:"流期天干（与 zhi 构成流期干支；非流宫宫干）。"`
	Zhi        Zhi                      `json:"zhi" jsonschema:"流时地支"`
	SiHua      []FlowTransformationView `json:"si_hua" jsonschema:"四化事实：星、化、落宫（落宫为本命宫序坐标）。"`
	GongWei    []FlowPalaceView         `json:"gong_wei" jsonschema:"流时十二宫盘。"`
	ShiChen    LiuShiChenView           `json:"shi_chen" jsonschema:"流时时辰定位。"`
}

func liupanView[T any](src any, dst *T) T {
	raw, err := json.Marshal(src)
	if err != nil {
		return *dst
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return *dst
	}
	return *dst
}

func (v LiuNian) View() LiuNianOutputView { return liupanView(v, &LiuNianOutputView{}) }
func (v LiuYue) View() LiuYueOutputView   { return liupanView(v, &LiuYueOutputView{}) }
func (v LiuRi) View() LiuRiOutputView     { return liupanView(v, &LiuRiOutputView{}) }
func (v LiuShi) View() LiuShiOutputView   { return liupanView(v, &LiuShiOutputView{}) }
