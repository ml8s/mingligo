package qimen

import (
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

type solarTermInfo struct {
	JieQi tianwen.JieQi
	Time  time.Time
}

func (s solarTermInfo) name() string { return tianwen.DecodeJieQi(s.JieQi) }
func (s solarTermInfo) index() int   { return int(s.JieQi) }

// DetermineJuShu 用 termTime（Gregorian time）判定节气边界，用 chartTime（真太阳时壁钟）
// 做时辰内的分钟分法，两者与 bz 的时柱口径保持一致。
func DetermineJuShu(termTime, chartTime time.Time, bz ganzhi.Bazi, method Method) juShu {
	actual := currentSolarTermInfo(termTime)
	switch scopeMethods[method.Scope].MethodSource {
	case "quarter_rule":
		variant := quarterRules[method.QuarterRule]
		if method.QuarterRule == QuarterTwelveMinuteTenDivision {
			base := DetermineJuShu(termTime, chartTime, bz, Method{
				Scope: variant.BaseScope, School: variant.BaseSchool, Dingju: method.BaseDingju,
			})
			quarter := newTwelveMinuteQuarterInfo(chartTime, bz.Shi, !base.YinDun, method)
			ju := twelveMinuteShiftJu(base.Number, quarter.Index, quarter.HourYangDun)
			return juShu{
				Method: method, Number: ju, YinDun: !quarter.HourYangDun,
				Yuan: base.Yuan, JieQi: actual.JieQi, YongJuJieQi: base.YongJuJieQi,
				Quarter: &quarter,
			}
		}
		quarter := newTenMinuteQuarterInfo(chartTime, bz.Shi)
		ju := quarter.Sanyuan.YangJu
		if !quarter.HourYangDun {
			ju = quarter.Sanyuan.YinJu
		}
		return juShu{
			Method: method, Number: ju, YinDun: !quarter.HourYangDun,
			Yuan: quarter.Sanyuan.Yuan, JieQi: actual.JieQi,
			YongJuJieQi: variant.YongJuJieQi, Quarter: &quarter,
		}
	case "month_cycle":
		group := monthGroups[bz.Nian.Zhi]
		monthSeq := floorMod(int(bz.Yue.Zhi)-int(monthFirstBranch), 12)
		ju := floorMod(group.LeadingJu-1-monthSeq, 9) + 1
		return juShu{
			Method: method, Number: ju, YinDun: true, Yuan: group.Yuan,
			JieQi: actual.JieQi, YongJuJieQi: "月家定局",
		}
	case "year_cycle":
		cycleYear := actual.Time.Year()
		if actual.index() < 3 {
			cycleYear--
		}
		years := cycleYear - yearAnchor
		element := floorMod((years-floorMod(years, yearCycle))/yearCycle, 3)
		yuan := Yuan(element)
		return juShu{
			Method: method, Number: yearJuByYuan[yuan], YinDun: true, Yuan: yuan,
			JieQi: actual.JieQi, YongJuJieQi: "年家定局",
		}
	default:
		yongJu := actual.JieQi
		var state ZhiRunState
		if method.Dingju == DingjuZhiRun {
			yongJu, state = determineZhiRunTerm(termTime, actual)
		}
		yuanIdx := int(determineYuan(bz.Ri))
		if method.Dingju == DingjuMaoShan {
			yuanIdx = determineMaoShanYuan(termTime, actual)
		}
		entry := solarTermJuTable[yongJu]
		return juShu{
			Method: method,
			Number: entry[yuanIdx], YinDun: entry[3] == 0,
			Yuan: Yuan(yuanIdx), JieQi: actual.JieQi,
			YongJuJieQi: tianwen.DecodeJieQi(yongJu), ZhiRunState: state,
		}
	}
}

func determineMaoShanYuan(t time.Time, actual solarTermInfo) int {
	anchor := actual.Time.Truncate(time.Hour)
	if actual.Time.Hour()%2 == 0 {
		anchor = anchor.Add(-time.Hour)
	}
	elapsed := t.Sub(anchor)
	if elapsed < 0 {
		elapsed = 0
	}
	shichen := int(elapsed / (2 * time.Hour))
	yuan := shichen / maoShanShichenPerYuan
	if yuan >= maoShanMaxYuan {
		return maoShanMaxYuan - 1
	}
	return yuan
}

func currentSolarTermInfo(t time.Time) solarTermInfo {
	result := solarTermInfo{JieQi: tianwen.JQDongZhi}
	var latest time.Time
	if previousWinter := tianwen.SolarTermTime(t.Year()-1, 270); !t.Before(previousWinter) {
		result = solarTermInfoAt(t.Year()-1, int(tianwen.JQDongZhi))
		latest = previousWinter
	}
	for jq := tianwen.JieQi(0); jq < 24; jq++ {
		term := tianwen.SolarTermTime(t.Year(), jq.Longitude())
		if !t.Before(term) && term.After(latest) {
			result = solarTermInfo{JieQi: jq, Time: term}
			latest = term
		}
	}
	return result
}

func solarTermInfoAt(year, index int) solarTermInfo {
	jq := tianwen.JieQi(index)
	return solarTermInfo{JieQi: jq, Time: tianwen.SolarTermTime(year, jq.Longitude())}
}

func solarTermSequence(info solarTermInfo) int {
	year := info.Time.Year()
	if info.index() == 0 {
		year++
	}
	return year*24 + info.index()
}

func determineZhiRunTerm(t time.Time, actual solarTermInfo) (tianwen.JieQi, ZhiRunState) {
	day := qimenDayNumber(t)
	leader := day - floorMod(day, zhiRunLeaderDays)
	upcomingLongitude := 90.0
	if actual.index() >= 12 {
		upcomingLongitude = 270
	}
	upcoming := nextSolarTerm(t, upcomingLongitude)
	currentLongitude := 270.0
	if upcomingLongitude == 270 {
		currentLongitude = 90
	}
	current := latestSolarTerm(t, currentLongitude)
	previous := solarTermInfoAt(current.Time.Year()-1, 12-current.index())

	for _, anchor := range []solarTermInfo{upcoming, current, previous} {
		anchorLeader := adoptingZhiRunLeader(anchor, t.Location())
		if anchorLeader > leader {
			continue
		}
		rawInterval := (leader - anchorLeader) / zhiRunLeaderDays
		interval := min(rawInterval, zhiRunMaxIntervals-1)
		workingSequence := solarTermSequence(anchor) + interval
		workingIndex := floorMod(workingSequence, 24)
		working := tianwen.JieQi(workingIndex)
		return working, zhiRunState(t, actual, working, upcoming, anchor, leader, rawInterval, workingSequence)
	}
	return actual.JieQi, zsBehind
}

func nextSolarTerm(t time.Time, longitude float64) solarTermInfo {
	year := t.Year()
	term := tianwen.SolarTermTime(year, longitude)
	if term.Before(t) {
		year++
		term = tianwen.SolarTermTime(year, longitude)
	}
	jq, _ := tianwen.JieQiOfLongitude(longitude)
	return solarTermInfo{JieQi: jq, Time: term}
}

func latestSolarTerm(t time.Time, longitude float64) solarTermInfo {
	year := t.Year()
	term := tianwen.SolarTermTime(year, longitude)
	if term.After(t) {
		year--
		term = tianwen.SolarTermTime(year, longitude)
	}
	jq, _ := tianwen.JieQiOfLongitude(longitude)
	return solarTermInfo{JieQi: jq, Time: term}
}

func adoptingZhiRunLeader(term solarTermInfo, location *time.Location) int {
	day := qimenDayNumber(term.Time.In(location))
	lead := floorMod(day, zhiRunLeaderDays)
	leader := day - lead
	if lead > zhiRunThreshold {
		leader += zhiRunLeaderDays
	}
	return leader
}

func zhiRunState(t time.Time, actual solarTermInfo, working tianwen.JieQi, upcoming, anchor solarTermInfo, leader, rawInterval, workingSequence int) ZhiRunState {
	workingTerm := solarTermInfoFromSequence(workingSequence)
	workingDay := qimenDayNumber(workingTerm.Time.In(t.Location()))
	actualSequence := solarTermSequence(actual)
	if workingDay == leader && workingSequence == actualSequence && !workingTerm.Time.After(t) {
		return zsExact
	}
	if qimenDayNumber(upcoming.Time.In(t.Location())) == qimenDayNumber(t) {
		return zsBehind
	}
	if rawInterval > zhiRunMaxIntervals-1 && zhiRunRepeatTerms[working] &&
		actual.JieQi == working && qimenDayNumber(anchor.Time.In(t.Location())) != qimenDayNumber(t) {
		return zsRepeat
	}
	if workingSequence > actualSequence {
		return zsAhead
	}
	return zsBehind
}

func solarTermInfoFromSequence(sequence int) solarTermInfo {
	sequenceYear := floorDiv(sequence, 24)
	index := floorMod(sequence, 24)
	occurrenceYear := sequenceYear
	if index == 0 {
		occurrenceYear--
	}
	return solarTermInfoAt(occurrenceYear, index)
}

func qimenDayNumber(t time.Time) int {
	year, month, day := t.Date()
	unixDay := int(time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Unix() / 86400)
	result := unixDay - 10963
	if t.Hour() >= 23 {
		result++
	}
	return result
}

func floorMod(value, modulus int) int {
	return ((value % modulus) + modulus) % modulus
}

func floorDiv(value, modulus int) int {
	return (value - floorMod(value, modulus)) / modulus
}

func determineYuan(dayZhu ganzhi.Zhu) int {
	dayIdx := ganzhi.SixtyCycleIndex(dayZhu.Gan, dayZhu.Zhi)
	fuTou := ganzhi.SixtyToZhu(dayIdx - dayIdx%5)
	return fuTouYuanTable[fuTou.Zhi]
}

func fuTou(dayZhu ganzhi.Zhu) ganzhi.Zhu {
	dayIdx := ganzhi.SixtyCycleIndex(dayZhu.Gan, dayZhu.Zhi)
	return ganzhi.SixtyToZhu(dayIdx - dayIdx%5)
}
