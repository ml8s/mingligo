package ziwei

import (
	"fmt"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

func hongLuanPos(zhi Zhi) int {
	// 红鸾：卯(年支1)起子(0)，顺数至生年支
	return (4 - int(zhi) + 12) % 12
}

// ── 流月 ──

// LiuYueJie locates the flow lunar month: the observer's standard clock and
// civil selection date, the lunar date, the month boundary, and the calendar
// interval that disambiguates leap months. Values are human-readable.
type LiuYueJie struct {
	StandardTime    string            `json:"standard_time"`
	StartDate       string            `json:"start_date"`
	Lunar           tianwen.LunarDate `json:"lunar"`
	LunarMonthStart string            `json:"lunar_month_start"`
	LunarMonthEnd   string            `json:"lunar_month_end"`
	BoundaryDay     bool              `json:"boundary_day"`
	CalendarPeriod  LunarMonthPeriod  `json:"calendar_period"`
	FlowMonth       LunarMonth        `json:"flow_month"`
}

type LiuYue struct {
	MingGong   GongIndex            `json:"ming_gong"`
	LiuMingZhi Zhi                  `json:"liu_ming_zhi"`
	Zhi        Zhi                  `json:"zhi"`
	SiHua      []FlowTransformation `json:"si_hua"`
	GongWei    [12]FlowPalace       `json:"gong_wei"`
	YueJie     LiuYueJie            `json:"yue_jie"`
}

// FlowMonthPeriod separates calendar truth from the Ziwei flow-month choice.
type FlowMonthPeriod struct {
	CalendarPeriod LunarMonthPeriod `json:"calendar_period"`
	FlowMonth      LunarMonth       `json:"flow_month"`
}

type LunarMonthPeriod struct {
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Leap     bool   `json:"leap"`
	Half     string `json:"half"`
	DayStart int    `json:"day_start"`
	DayEnd   int    `json:"day_end"`
}

type LunarMonth struct {
	Year  int  `json:"year"`
	Month int  `json:"month"`
	Leap  bool `json:"leap"`
}

// ComputeLiuYue computes the flow month chart.
// target IS a complete resolved lunar date. The day is required because the
// iztro-compatible rule splits a leap month at day 15; no half/fixLeap flag
// is exposed to callers.
func ComputeLiuYueWith(ctx FlowContext, target tianwen.LunarDate, p tianwen.LunarProjector) (LiuYue, error) {
	period, err := resolveLiuYuePeriodWith(p, target)
	if err != nil {
		return LiuYue{}, err
	}
	return computeLiuYueWithPeriod(ctx, target, period)
}

func computeLiuYueWithPeriod(ctx FlowContext, target tianwen.LunarDate, period FlowMonthPeriod) (LiuYue, error) {
	flowYear, flowMonth := period.FlowMonth.Year, period.FlowMonth.Month
	flowYearGan, _ := yearGanZhi(flowYear)
	// 流月干支 = 目标月干支（五虎遁：年干 + 农历月）——真相源，与命主无关
	yueZhi := zhiIdxToZhi((flowMonth + 1) % 12) // 正月寅起：农历1月→zhiIdx 2(寅)
	monthGan := yueGanByWuHuDun(flowYearGan, yueZhi)
	stars := liuMinorStars(monthGan, yueZhi)

	// 盘起点 = monthlyIndex（iztro 公式，命主相关）
	mi := computeMonthlyIndex(ctx, flowTarget{
		Year:        target.Year,
		LunarMonth:  target.Month,
		LunarDay:    target.Day,
		IsLeapMonth: target.Leap,
	})
	short := liuMinorShort("月")
	starByAnXingIdx := make(map[int][]FlowStar)
	for sid, z := range stars {
		anXingIdx := zhiIdxToAnXingIdx(zhiToZhiIdx(z))
		name, ok := short[sid]
		if !ok {
			name = starName(sid)
		}
		starByAnXingIdx[anXingIdx] = append(starByAnXingIdx[anXingIdx], flowStarOf(sid, name))
	}
	return LiuYue{
		MingGong:   0,
		LiuMingZhi: flowMingZhi(mi),
		Zhi:        yueZhi,
		SiHua:      newFlowTransformations(computeSiHua(monthGan), ctx.SiHuaPalaces),
		GongWei:    buildFlowPalaces(mi, monthGan, starByAnXingIdx),
	}, nil
}

// ComputeLiuYueAt computes the flow month directly from a standard time. It
// builds the lunar projector, resolves the month boundary, and fills the YueJie
// facts from the same standard time.
func ComputeLiuYueAt(ctx FlowContext, standard time.Time, loc *time.Location) (LiuYue, error) {
	p, err := tianwen.NewLunarProjector(tianwen.LunarBasisLocalAstronomical, loc)
	if err != nil {
		return LiuYue{}, err
	}
	date := p.SelectorDate(standard)
	resolved, err := ResolveLiuYueTimeWith(p, date)
	if err != nil {
		return LiuYue{}, err
	}
	period, err := resolveLiuYuePeriodWith(p, resolved.Lunar)
	if err != nil {
		return LiuYue{}, err
	}
	result, err := computeLiuYueWithPeriod(ctx, resolved.Lunar, period)
	if err != nil {
		return LiuYue{}, err
	}
	if loc == nil {
		loc = time.UTC
	}
	result.YueJie = LiuYueJie{
		StandardTime:    standard.In(loc).Format(time.RFC3339Nano),
		StartDate:       tianwen.StandardDateAt(standard, loc).String(),
		Lunar:           resolved.Lunar,
		LunarMonthStart: resolved.Start.String(),
		LunarMonthEnd:   resolved.End.String(),
		BoundaryDay:     resolved.BoundaryDay,
		CalendarPeriod:  period.CalendarPeriod,
		FlowMonth:       period.FlowMonth,
	}
	return result, nil
}

func resolveLiuYuePeriodWith(p tianwen.LunarProjector, target tianwen.LunarDate) (FlowMonthPeriod, error) {
	meta, err := p.ValidateLunarDate(target)
	if err != nil {
		return FlowMonthPeriod{}, fmt.Errorf("resolve flow month: %w", err)
	}

	calendar := LunarMonthPeriod{
		Year:     target.Year,
		Month:    target.Month,
		Leap:     target.Leap,
		Half:     "whole",
		DayStart: 1,
		DayEnd:   meta.DayCount,
	}
	flowMonth := LunarMonth{Year: target.Year, Month: target.Month, Leap: target.Leap}
	if target.Leap {
		if target.Day <= 15 {
			calendar.Half = "first"
			calendar.DayEnd = 15
		} else {
			calendar.Half = "second"
			calendar.DayStart = 16
			next, err := p.NextLunarMonth(target.Year, target.Month, true)
			if err != nil {
				return FlowMonthPeriod{}, fmt.Errorf("resolve flow month: %w", err)
			}
			flowMonth = LunarMonth{Year: next.Year, Month: next.Month, Leap: next.Leap}
		}
	}
	return FlowMonthPeriod{CalendarPeriod: calendar, FlowMonth: flowMonth}, nil
}

// liuMinorStars returns the ten flow minor stars as typed star symbols. The
// period-specific display name (月禄/日羊/...) is added only at the output edge.
func liuMinorStars(gan Gan, zhi Zhi) map[StarIndex]Zhi {
	chg, qu := liuChangQuByGan(gan)
	toZhi := func(zhiIdx int) Zhi { return Zhi(zhiIdx + 1) }
	return map[StarIndex]Zhi{
		LuCun:    toZhi(luCunPos(gan)),
		QingYang: toZhi(qingYangPos(gan)),
		TuoLuo:   toZhi(tuoLuoPos(gan)),
		TianKui:  toZhi(tianKuiPos(gan)),
		TianYue:  toZhi(tianYuePos(gan)),
		TianMa:   toZhi(tianMaPos(zhi)),
		HongLuan: toZhi(hongLuanPos(zhi)),
		TianXi:   toZhi((hongLuanPos(zhi) + 6) % 12),
		WenChang: toZhi(chg),
		WenQu:    toZhi(qu),
	}
}

// liuMinorShort maps a flow star symbol to its period-prefixed display name.
func liuMinorShort(prefix string) map[StarIndex]string {
	return map[StarIndex]string{
		LuCun: prefix + "禄", QingYang: prefix + "羊", TuoLuo: prefix + "陀",
		TianKui: prefix + "魁", TianYue: prefix + "钺", TianMa: prefix + "马",
		HongLuan: prefix + "鸾", TianXi: prefix + "喜",
		WenChang: prefix + "昌", WenQu: prefix + "曲",
	}
}

// ── 流日 ──

// LiuRiJie locates the flow day: the observer's standard clock and civil date,
// plus the lunar date. Values are human-readable.
type LiuRiJie struct {
	StandardTime string            `json:"standard_time"`
	StandardDate string            `json:"standard_date"`
	Lunar        tianwen.LunarDate `json:"lunar"`
}

type LiuRi struct {
	MingGong   GongIndex            `json:"ming_gong"`
	LiuMingZhi Zhi                  `json:"liu_ming_zhi"`
	Zhi        Zhi                  `json:"zhi"`
	SiHua      []FlowTransformation `json:"si_hua"`
	GongWei    [12]FlowPalace       `json:"gong_wei"`
	RiJie      LiuRiJie             `json:"ri_jie"`
}

// ComputeLiuRiAt computes the flow day directly from a standard time. It
// builds the lunar projector and resolves the lunar day internally, then fills
// the RiJie facts from the same standard time.
func ComputeLiuRiAt(ctx FlowContext, standard time.Time, loc *time.Location) (LiuRi, error) {
	p, err := tianwen.NewLunarProjector(tianwen.LunarBasisLocalAstronomical, loc)
	if err != nil {
		return LiuRi{}, err
	}
	resolved, err := ResolveLiuRiTimeAtWith(p, standard, loc)
	if err != nil {
		return LiuRi{}, err
	}
	result, err := ComputeLiuRiWith(ctx, resolved.Lunar, p)
	if err != nil {
		return LiuRi{}, err
	}
	result.RiJie = LiuRiJie{
		StandardTime: standard.In(loc).Format(time.RFC3339Nano),
		StandardDate: resolved.StandardDate.String(),
		Lunar:        resolved.Lunar,
	}
	return result, nil
}

// ComputeLiuRi computes the flow day chart using iztro's dailyIndex formula.
// lunarMonth/lunarDay ARE lunar (domain-native). Day pillar requires Gregorian:
// internally converts lunar→solar via the projector inverse.
func ComputeLiuRiWith(ctx FlowContext, target tianwen.LunarDate, p tianwen.LunarProjector) (LiuRi, error) {
	if _, err := p.ValidateLunarDate(target); err != nil {
		return LiuRi{}, fmt.Errorf("validate flow day: %w", err)
	}

	// 盘起点 = dailyIndex（直接用农历）
	di := computeDailyIndex(ctx, flowTarget{
		Year:        target.Year,
		LunarMonth:  target.Month,
		LunarDay:    target.Day,
		IsLeapMonth: target.Leap,
	})
	// 流日干支 = 目标日干支（需公历——农历转公历）
	date, err := p.DateFromLunar(target)
	if err != nil {
		return LiuRi{}, fmt.Errorf("validate flow day: %w", err)
	}
	zhu := tianwen.RiZhu(date.Time())
	riGan, riZhi := Gan(zhu.Gan), Zhi(zhu.Zhi)
	stars := liuMinorStars(riGan, riZhi)
	short := liuMinorShort("日")
	starByAnXingIdx := make(map[int][]FlowStar)
	for sid, z := range stars {
		anXingIdx := zhiIdxToAnXingIdx(zhiToZhiIdx(z))
		name, ok := short[sid]
		if !ok {
			name = starName(sid)
		}
		starByAnXingIdx[anXingIdx] = append(starByAnXingIdx[anXingIdx], flowStarOf(sid, name))
	}
	return LiuRi{
		MingGong:   0,
		LiuMingZhi: flowMingZhi(di),
		Zhi:        riZhi,
		SiHua:      newFlowTransformations(computeSiHua(riGan), ctx.SiHuaPalaces),
		GongWei:    buildFlowPalaces(di, riGan, starByAnXingIdx),
	}, nil
}

// ── 流时 ──

// LiuShiChen locates the flow hour: the observer's standard clock and civil
// date, the lunar date, and the true-solar 时辰. Values are human-readable.
type LiuShiChen struct {
	StandardTime string            `json:"standard_time"`
	StandardDate string            `json:"standard_date"`
	Lunar        tianwen.LunarDate `json:"lunar"`
	Shichen      string            `json:"shichen"`
}

type LiuShi struct {
	MingGong   GongIndex            `json:"ming_gong"`
	LiuMingZhi Zhi                  `json:"liu_ming_zhi"`
	Zhi        Zhi                  `json:"zhi"`
	SiHua      []FlowTransformation `json:"si_hua"`
	GongWei    [12]FlowPalace       `json:"gong_wei"`
	ShiChen    LiuShiChen           `json:"shi_chen"`
}

// ComputeLiuShiAt computes the flow hour directly from a standard time. It
// builds the lunar projector and resolves the true-solar 时辰 internally, then
// fills the ShiChen facts from the same standard time.
func ComputeLiuShiAt(ctx FlowContext, standard time.Time, loc *time.Location, longitude float64) (LiuShi, error) {
	p, err := tianwen.NewLunarProjector(tianwen.LunarBasisLocalAstronomical, loc)
	if err != nil {
		return LiuShi{}, err
	}
	resolved, err := ResolveLiuShiTimeWith(p, standard, loc, longitude)
	if err != nil {
		return LiuShi{}, err
	}
	result, err := ComputeLiuShiAtTimeWith(ctx, resolved, p)
	if err != nil {
		return LiuShi{}, err
	}
	if loc == nil {
		loc = time.UTC
	}
	result.ShiChen = LiuShiChen{
		StandardTime: standard.In(loc).Format(time.RFC3339Nano),
		StandardDate: resolved.StandardDate.String(),
		Lunar:        resolved.Lunar,
		Shichen:      ganzhi.DecodeZhi(resolved.Shichen),
	}
	return result, nil
}

// ComputeLiuShi computes the flow hour chart using iztro's hourlyIndex formula.
// lunarMonth/lunarDay ARE lunar (domain-native); shiZhi is the target hour zhi.
func ComputeLiuShiWith(ctx FlowContext, target tianwen.LunarDate, shiZhi Zhi, p tianwen.LunarProjector) (LiuShi, error) {
	if _, err := p.ValidateLunarDate(target); err != nil {
		return LiuShi{}, fmt.Errorf("validate flow hour: %w", err)
	}

	// 盘起点 = hourlyIndex（直接用农历）
	hi := computeHourlyIndex(ctx, flowTarget{
		Year:        target.Year,
		LunarMonth:  target.Month,
		LunarDay:    target.Day,
		IsLeapMonth: target.Leap,
		ShiZhi:      shiZhi,
	})
	// 流时干支：日干 + 五鼠遁（需公历——农历转公历）
	date, err := p.DateFromLunar(target)
	if err != nil {
		return LiuShi{}, fmt.Errorf("validate flow hour: %w", err)
	}
	zhu := tianwen.RiZhu(date.Time())
	riGan := Gan(zhu.Gan)
	shiGan := shiGanCalc(riGan, shiZhi)
	stars := liuMinorStars(shiGan, shiZhi)
	short := liuMinorShort("时")
	starByAnXingIdx := make(map[int][]FlowStar)
	for sid, z := range stars {
		anXingIdx := zhiIdxToAnXingIdx(zhiToZhiIdx(z))
		name, ok := short[sid]
		if !ok {
			name = starName(sid)
		}
		starByAnXingIdx[anXingIdx] = append(starByAnXingIdx[anXingIdx], flowStarOf(sid, name))
	}
	return LiuShi{
		MingGong:   0,
		LiuMingZhi: flowMingZhi(hi),
		Zhi:        shiZhi,
		SiHua:      newFlowTransformations(computeSiHua(shiGan), ctx.SiHuaPalaces),
		GongWei:    buildFlowPalaces(hi, shiGan, starByAnXingIdx),
	}, nil
}

func shiGanCalc(riGan Gan, shiZhi Zhi) Gan {
	// 五鼠遁：甲己→甲子, 乙庚→丙子, 丙辛→戊子, 丁壬→庚子, 戊癸→壬子
	ziGan := (int(riGan)-1)*2%10 + 1
	return Gan(((ziGan-1+int(shiZhi)-1)%10+10)%10 + 1)
}

// 流昌流曲：天干定位（iztro算法）, 返回zhiIdx
func liuChangQuByGan(gan Gan) (changZhiIdx, quZhiIdx int) {
	table := [10][2]int{
		{3, 7}, {4, 6}, {6, 4}, {7, 3}, {6, 4},
		{7, 3}, {9, 1}, {10, 0}, {0, 10}, {1, 9},
	}
	idx := int(gan) - 1
	return (table[idx][0] + 2) % 12, (table[idx][1] + 2) % 12
}

// yueGanByWuHuDun computes the month gan via 五虎遁 (year gan + month zhi).
func yueGanByWuHuDun(nianGan Gan, yueZhi Zhi) Gan {
	// 寅月干: 甲己→丙, 乙庚→戊, 丙辛→庚, 丁壬→壬, 戊癸→甲
	base := [10]Gan{3, 5, 7, 9, 1, 3, 5, 7, 9, 1} // 寅月干（甲=1..癸=10）
	offset := int(yueZhi) - 3                     // 寅=0 偏移
	if offset < 0 {
		offset += 12
	}
	return Gan(((int(base[nianGan-1]) - 1 + offset) % 10) + 1)
}
