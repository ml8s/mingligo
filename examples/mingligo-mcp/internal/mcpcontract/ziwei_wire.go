// MingliGo wire DTOs for the MCP contract layer.
//
// These types are ordinary source now, not generated artifacts. Edit them
// directly when the published MingliGo payload changes.
package mcpcontract

import "github.com/ml8s/mingligo/pkg/ziwei"

// ── ziwei.chart input ──

// ZiweiChartInput consumes a resolved natal subject. The chart domain derives
// the local astronomical lunar date and LAST shichen from the subject.
type ZiweiChartInput struct {
	Subject NatalSubject `json:"subject" jsonschema:"已归一化的出生主体。"`
}

type ZiweiDaXianItems []ZiweiDaXianItem

type ZiweiDaXianItem struct {
	StartYear int             `json:"start_year" jsonschema:"本限起始公历年"`
	EndYear   int             `json:"end_year" jsonschema:"本限结束公历年"`
	QiSui     int             `json:"qi_sui" jsonschema:"起运虚岁。"`
	ZhiSui    int             `json:"zhi_sui" jsonschema:"止运虚岁。"`
	MingGong  ZiweiPalaceName `json:"ming_gong" jsonschema:"大限命宫。"`
}

// ── ziwei.liupan input ──

type ZiweiLiupanInput struct {
	Scope        ZiweiLiupanScopeName `json:"scope" jsonschema:"流盘层级。"`
	Subject      NatalSubject         `json:"subject" jsonschema:"已归一化的出生主体。"`
	StandardTime StandardTime         `json:"standard_time" jsonschema:"流盘选择起点：local standard time，RFC3339 fixed offset。"`
	Count        *SeriesCount         `json:"count,omitempty" jsonschema:"连续查询数量。"`
	Longitude    *Longitude           `json:"longitude,omitempty" jsonschema:"地点经度；liu_shi 必填（-180..180），真太阳时校正用于时辰。"`
}

// ── ziwei.liupan liu_nian output ──

type ZiweiLiuNianOutput struct {
	MingGong     ZiweiPalaceName           `json:"ming_gong" jsonschema:"命宫。"`
	MingGongName string                    `json:"ming_gong_name" jsonschema:"流年命宫名。"`
	Zhi          ZhiName                   `json:"zhi" jsonschema:"流年地支"`
	SiHua        []ZiweiFlowTransformation `json:"si_hua" jsonschema:"四化事实：星、化、落宫。"`
	GongWei      []ZiweiFlowPalace         `json:"gong_wei" jsonschema:"流年十二宫盘。"`
	NianJie      ZiweiLiuNianResolved      `json:"nian_jie" jsonschema:"流年农历年界。"`
}

type ZiweiLiuNianResolved struct {
	LiuNianStandardTime StandardTime `json:"standard_time" jsonschema:"流年选择器时区标准时间（无 DST）。"`
	LiuNianStartDate    DateString   `json:"start_date" jsonschema:"按选定农历基准取的流年选择日期。"`
	LunarYear           int          `json:"lunar_year" jsonschema:"该日期所在的农历年。"`
	LunarYearStart      DateString   `json:"lunar_year_start" jsonschema:"正月初一公历日（含）。"`
	LunarYearEnd        DateString   `json:"lunar_year_end" jsonschema:"下一正月初一公历日（不含）。"`
	BoundaryDay         bool         `json:"boundary_day" jsonschema:"输入日期是否为正月初一边界日（归入新年）。"`
}

type ZiweiFlowTransformation struct {
	Star   string          `json:"xing" jsonschema:"被化的本命星。"`
	Hua    SiHuaName       `json:"hua" jsonschema:"化禄/化权/化科/化忌。"`
	Palace ZiweiPalaceName `json:"gong" jsonschema:"本命星所在宫。"`
}

type ZiweiFlowPalace struct {
	Zhi       ZhiName         `json:"zhi" jsonschema:"地支。"`
	Name      ZiweiPalaceName `json:"name" jsonschema:"宫名。"`
	Stars     []string        `json:"xing" jsonschema:"流曜。"`
	IsLiuMing bool            `json:"is_liu_ming" jsonschema:"是否流盘命宫。"`
}

// ── ziwei.liupan liu_yue output ──

type ZiweiLunarDate struct {
	Year  int  `json:"year" jsonschema:"农历年。"`
	Month int  `json:"month" jsonschema:"农历月。"`
	Day   int  `json:"day" jsonschema:"农历日。"`
	Leap  bool `json:"leap" jsonschema:"是否闰月。"`
}

type ZiweiYueJie struct {
	LiuYueStandardTime StandardTime                            `json:"standard_time" jsonschema:"流月选择器时区标准时间（无 DST）。"`
	LiuYueStartDate    DateString                              `json:"start_date" jsonschema:"按选定农历基准取的流月选择日期。"`
	Lunar              ZiweiLunarDate                          `json:"lunar" jsonschema:"该日期对应的农历日期。"`
	LunarMonthStart    DateString                              `json:"lunar_month_start" jsonschema:"农历月初一公历日（含）。"`
	LunarMonthEnd      DateString                              `json:"lunar_month_end" jsonschema:"下一农历月初一公历日（不含）。"`
	BoundaryDay        bool                                    `json:"boundary_day" jsonschema:"输入日期是否为农历月初一边界日（归入新月）。"`
	CalendarPeriod     ZiweiLiuYueResolvedPeriodCalendarPeriod `json:"calendar_period" jsonschema:"流月所在日历区间。"`
	FlowMonth          ZiweiLiuYueResolvedPeriodFlowMonth      `json:"flow_month" jsonschema:"流月农历年月。"`
}

type ZiweiLiuYueOutput struct {
	YueJie       ZiweiYueJie               `json:"yue_jie" jsonschema:"流月农历月界与月令。"`
	MingGong     ZiweiPalaceName           `json:"ming_gong" jsonschema:"命宫。"`
	MingGongName string                    `json:"ming_gong_name" jsonschema:"流月命宫名。"`
	Zhi          ZhiName                   `json:"zhi" jsonschema:"流月地支"`
	SiHua        []ZiweiFlowTransformation `json:"si_hua" jsonschema:"四化事实：星、化、落宫。"`
	GongWei      []ZiweiFlowPalace         `json:"gong_wei" jsonschema:"流月十二宫盘。"`
}

type ZiweiLiuYueResolvedPeriodCalendarPeriod struct {
	Year     int                    `json:"year" jsonschema:"年。"`
	Month    int                    `json:"month" jsonschema:"月。"`
	Leap     bool                   `json:"leap" jsonschema:"是否闰月。"`
	Half     CalendarPeriodHalfName `json:"half" jsonschema:"农历半月；whole 为整月，闰月按 first / second 拆分。"`
	DayStart int                    `json:"day_start" jsonschema:"流月日段起始农历日序。"`
	DayEnd   int                    `json:"day_end" jsonschema:"流月日段结束农历日序。"`
}

type ZiweiLiuYueResolvedPeriodFlowMonth struct {
	Year  int  `json:"year" jsonschema:"年。"`
	Month int  `json:"month" jsonschema:"月。"`
	Leap  bool `json:"leap" jsonschema:"是否闰月。"`
}

type ZiweiLiuRiJie struct {
	LiuRiStandardTime StandardTime   `json:"standard_time" jsonschema:"流日选择器时区标准时间（无 DST）。"`
	LiuRiStandardDate DateString     `json:"standard_date" jsonschema:"流日选择器的时区标准日（无 DST）。"`
	Lunar             ZiweiLunarDate `json:"lunar" jsonschema:"该日期对应的农历日期。"`
}

// ── ziwei.liupan liu_ri output ──

type ZiweiLiuRiOutput struct {
	RiJie        ZiweiLiuRiJie             `json:"ri_jie" jsonschema:"流日日界。"`
	MingGong     ZiweiPalaceName           `json:"ming_gong" jsonschema:"命宫。"`
	MingGongName string                    `json:"ming_gong_name" jsonschema:"流日命宫名。"`
	Zhi          ZhiName                   `json:"zhi" jsonschema:"流日地支"`
	SiHua        []ZiweiFlowTransformation `json:"si_hua" jsonschema:"四化事实：星、化、落宫。"`
	GongWei      []ZiweiFlowPalace         `json:"gong_wei" jsonschema:"流日十二宫盘。"`
}

// ── ziwei.liupan liu_shi output ──

type ZiweiLiuShiChen struct {
	LiuShiStandardTime StandardTime   `json:"standard_time" jsonschema:"流时选择器时区标准时间（无 DST）。"`
	LiuShiStandardDate DateString     `json:"standard_date" jsonschema:"流时选择器的时区标准日（无 DST）。"`
	Lunar              ZiweiLunarDate `json:"lunar" jsonschema:"该日期对应的农历日期。"`
	Shichen            ZhiName        `json:"shichen" jsonschema:"真太阳时时辰。"`
}

type ZiweiLiuShiOutput struct {
	MingGong     ZiweiPalaceName           `json:"ming_gong" jsonschema:"命宫。"`
	MingGongName string                    `json:"ming_gong_name" jsonschema:"流时命宫名。"`
	Zhi          ZhiName                   `json:"zhi" jsonschema:"流时地支"`
	SiHua        []ZiweiFlowTransformation `json:"si_hua" jsonschema:"四化事实：星、化、落宫。"`
	GongWei      []ZiweiFlowPalace         `json:"gong_wei" jsonschema:"流时十二宫盘。"`
	ShiChen      ZiweiLiuShiChen           `json:"shi_chen" jsonschema:"流时时辰定位。"`
}

// ── ziwei.liupan union output ──

// ZiweiLiupanOutput carries exactly one selected flow layer. The explicit
// output schema keeps the current single-root wire face while allowing one
// typed handler to serve the four scope forms.
type ZiweiLiupanOutput struct {
	LiuNian []ziwei.LiuNian `json:"liu_nian,omitempty" jsonschema:"流年数组，长度 = count。"`
	LiuYue  []ziwei.LiuYue  `json:"liu_yue,omitempty" jsonschema:"流月数组，长度 = count。"`
	LiuRi   []ziwei.LiuRi   `json:"liu_ri,omitempty" jsonschema:"流日数组，长度 = count。"`
	LiuShi  []ziwei.LiuShi  `json:"liu_shi,omitempty" jsonschema:"流时数组，长度 = count。"`
}

// ── ziwei.bond input ──

type ZiweiBondSide struct {
	Subject NatalSubject `json:"subject" jsonschema:"已归一化的出生主体。"`
}

type ZiweiBondInput struct {
	A ZiweiBondSide `json:"a" jsonschema:"甲方：已归一化的出生主体。"`
	B ZiweiBondSide `json:"b" jsonschema:"乙方：已归一化的出生主体。"`
}

// ── ziwei.bond output ──

type ZiweiBondOutput struct {
	MingGong ZiweiBondPalacePair `json:"ming_gong" jsonschema:"双方本命宫。"`
	FuQiGong ZiweiBondPalacePair `json:"fu_qi_gong" jsonschema:"双方夫妻宫。"`
	ZiNvGong ZiweiBondPalacePair `json:"zi_nv_gong" jsonschema:"双方子女宫。"`
}

type ZiweiBondPalaceFacts struct {
	Palace     string   `json:"gong"`
	MajorStars []string `json:"major_xing"`
}

type ZiweiBondPalacePair struct {
	A      ZiweiBondPalaceFacts `json:"a"`
	B      ZiweiBondPalaceFacts `json:"b"`
	Status string               `json:"status"`
}

// ── ziwei.chart output ──
