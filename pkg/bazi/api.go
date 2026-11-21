// Package bazi provides 八字 computation.
//
// Types
//
//	Chart, Core,
//	DaYunStep, DaYun,
//	Bond,
//	LiuNian, LiuYue, LiuRi, LiuShi,
//	XiaoYunZhu,
//	FuYiResult, TiaoHouResult, GeJuResult
package bazi

import (
	"fmt"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

// computeChartFromTimes is the internal chart primitive for already-projected times.
// (年/月柱) and true-solar time (日/时柱).
func computeChartFromTimes(gt time.Time, st time.Time, gender ganzhi.Gender) Chart {
	if err := Load(); err != nil {
		panic(fmt.Errorf("bazi: data tables are not available: %w", err))
	}
	return ComputeChartFromCore(computeCore(gt, st, tianwen.ComputeBazi(gt, st), gender))
}

// ComputeChartAt is the standard-time entry for new callers.
func ComputeChartAt(
	standard time.Time,
	longitude float64,
	gender ganzhi.Gender,
) Chart {
	solar := tianwen.ApparentSolarFromStandard(standard, longitude)
	return computeChartFromTimes(standard, solar, gender)
}

// ComputeCoreAt is the standard-time entry for derived-method callers.
func ComputeCoreAt(
	standard time.Time,
	longitude float64,
	gender ganzhi.Gender,
) Core {
	solar := tianwen.ApparentSolarFromStandard(standard, longitude)
	return computeCoreFromTimes(standard, solar, gender)
}

// computeCoreFromTimes is the internal core primitive for already-projected times.
func computeCoreFromTimes(gt time.Time, st time.Time, gender ganzhi.Gender) Core {
	return computeCore(gt, st, tianwen.ComputeBazi(gt, st), gender)
}

// ComputeChartFromCore fills the complete chart from a narrow natal core.
func ComputeChartFromCore(core Core) Chart {
	bz := core.ToBazi()
	analysis := AnalyzePillars(bz, core.Gender)
	chart := assembleChart(core, bz, analysis)
	chart = withDaYunRoots(chart, bz)
	chart = withChartExtra(chart, core)
	chart = withInteractions(chart, bz)
	chart = withYongshenEvidence(chart, core, analysis)
	chart = withDerivedFacts(chart)
	chart.ZiShiRule = "晚子时(23:00-24:00)日柱不变、时柱按次日日干起（lunar 约定）"
	return chart
}

// applyChartExtra derives the supplementary 三元 / 胎元 / 命宫 / 纳音 projections.
func withChartExtra(chart Chart, core Core) Chart {
	chart = withChartExtraFields(chart, core)
	chart.Zodiac = ganzhi.ZhiAnimal(core.ToBazi().Nian.Zhi)
	return chart
}

// applyInteractions derives the complete stem/branch interaction groups.
func withInteractions(chart Chart, bz ganzhi.Bazi) Chart {
	detection := DetectRelations(bz)
	chart.GuanXi = buildNatalRelationsFromDetection(bz, detection)
	return chart
}

// applyYongshenEvidence derives FuYi/TiaoHou candidates and GeJu structure.
func withYongshenEvidence(chart Chart, core Core, analysis PillarAnalysis) Chart {
	chart.FuYi, chart.TiaoHou, chart.GeJu = computeYongShenEvidenceFromAnalysis(core, analysis, chart.GuanXi)
	chart.ShenShaSchool = ShenShaSchool{
		DualReference: []string{"天乙贵人", "桃花", "驿马", "华盖", "将星", "劫煞", "灾煞"},
		Policy:        "union_of_year_and_day_references",
	}
	return chart
}

// applyDerivedFacts derives the read-only facts consumed by judgment layers.
func withDerivedFacts(chart Chart) Chart {
	var elementStates map[ganzhi.Wuxing]ElementState
	chart.TenGodStates, elementStates = computeTenGodStates(chart)
	chart.ElementStates = elementStateList(elementStates)
	chart.LuRoots = computeLuRoots(chart)
	facts := computeChartConclusions(chart)
	chart.GuanSha = GuanSha{QuQing: facts.OfficerKillingCleaned}
	chart.CaiKu = CaiKu{Xian: facts.WealthTombPresent, XingRuKu: facts.WealthStarInTomb}
	chart.GeJu.GeShenTou = facts.PatternGodTransparent
	chart.SpousePalaceState = spousePalaceState(chart)
	chart.DayBranchType = dayBranchType(chart)
	return chart
}

func withDaYunRoots(chart Chart, bz ganzhi.Bazi) Chart {
	if chart.DaYun == nil {
		return chart
	}
	natalHidden := natalHiddenStems(bz)
	daYun := *chart.DaYun
	daYun.Steps = append([]DaYunStep(nil), chart.DaYun.Steps...)
	for index := range daYun.Steps {
		step := &daYun.Steps[index]
		step.RootRefs = daYunRootRefs(ganzhi.Zhu{Gan: step.Gan, Zhi: step.Zhi}, natalHidden)
		rooted := len(step.RootRefs) > 0
		step.Rooted = &rooted
	}
	chart.DaYun = &daYun
	return chart
}

// ComputeLiuNianAtTime computes the year pillar from a resolved 立春 year.
// ComputeLiuNianAt computes the flow year directly from a standard time. It
// resolves the 立春 year boundary internally and fills the NianJie location
// facts, so callers only pass the standard time.
func ComputeLiuNianAt(core Core, standard time.Time, loc *time.Location) (*LiuNian, error) {
	t, err := ResolveLiuNianTimeAt(standard, loc)
	if err != nil {
		return nil, fmt.Errorf("computeLiuNian: resolve: %w", err)
	}
	cd := dayunStepForYear(core.DaYun, core.BirthYear, t.LichunYear)
	result, err := computeLiuNian(core.ToBazi(), t, cd)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		loc = time.UTC
	}
	result.NianJie = LiuNianJie{
		StandardTime: standard.In(loc).Format(time.RFC3339Nano),
		StartDate:    tianwen.StandardDateAt(standard, loc).String(),
		LichunYear:   t.LichunYear,
		LichunStart:  t.Start.String(),
		LichunEnd:    t.End.String(),
		BoundaryDay:  t.BoundaryDay,
	}
	return result, nil
}

func ComputeLiuNianAtTime(core Core, t LiuNianTime) (*LiuNian, error) {
	cd := dayunStepForYear(core.DaYun, core.BirthYear, t.LichunYear)
	return computeLiuNian(core.ToBazi(), t, cd)
}

// ComputeLiuYueAt computes the flow month directly from a standard time. It
// resolves the solar-term boundary internally and fills the JieQi facts.
func ComputeLiuYueAt(core Core, standard time.Time, loc *time.Location) (*LiuYue, error) {
	t, err := ResolveLiuYueTimeAt(standard, loc)
	if err != nil {
		return nil, fmt.Errorf("computeLiuYue: resolve: %w", err)
	}
	nianForRefs, err := ResolveLiuNianTimeAt(standard, loc)
	if err != nil {
		return nil, fmt.Errorf("computeLiuYue: refs nian: %w", err)
	}
	dayForRefs := tianwen.RiZhu(t.Date.Time())
	refs, err := resolveFlowRefs(core, nianForRefs, ganzhi.Zhu{Gan: dayForRefs.Gan, Zhi: dayForRefs.Zhi})
	if err != nil {
		return nil, fmt.Errorf("computeLiuYue: refs: %w", err)
	}
	result, err := computeLiuYue(core.ToBazi(), t, refs)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		loc = time.UTC
	}
	result.JieQi = LiuYueJieQi{
		StandardTime: standard.In(loc).Format(time.RFC3339Nano),
		StartDate:    t.Date.String(),
		JieqiStart:   t.Start.String(),
		JieqiEnd:     t.End.String(),
		BoundaryDay:  t.BoundaryDay,
	}
	return result, nil
}

// ComputeLiuYueAtTime computes the month pillar from a resolved solar-term month.
func ComputeLiuYueAtTime(core Core, t LiuYueTime) (*LiuYue, error) {
	nianForRefs, err := ResolveLiuNianTime(tianwen.NewDate(t.LichunYear, 6, 15))
	if err != nil {
		return nil, err
	}
	dayForRefs := tianwen.RiZhu(t.Date.Time())
	refs, err := resolveFlowRefs(core, nianForRefs, ganzhi.Zhu{Gan: dayForRefs.Gan, Zhi: dayForRefs.Zhi})
	if err != nil {
		return nil, err
	}
	return computeLiuYue(core.ToBazi(), t, refs)
}

// ComputeLiuRiAt computes the flow day directly from a standard time. It
// resolves the true-solar day boundary internally and fills the RiJie location
// facts, so callers only pass the standard time.
func ComputeLiuRiAt(core Core, standard time.Time, loc *time.Location, longitude float64) (*LiuRi, error) {
	t, err := ResolveLiuRiTimeAt(standard, loc, longitude)
	if err != nil {
		return nil, fmt.Errorf("computeLiuRi: resolve: %w", err)
	}
	// 流日的流年/大运按该日所属的立春年取，而不是公历年。
	nianTime, err := ResolveLiuNianTime(t.Date)
	if err != nil {
		return nil, fmt.Errorf("computeLiuRi: liunian time: %w", err)
	}
	dzZhu, lnZhu, err := liuRiDayunLiunian(core, nianTime)
	if err != nil {
		return nil, err
	}
	result, err := computeLiuRi(core.ToBazi(), t, dzZhu, lnZhu)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		loc = time.UTC
	}
	result.RiJie = LiuRiJie{
		StandardTime: t.Standard.In(loc).Format(time.RFC3339Nano),
		StandardDate: t.StandardDate.String(),
	}
	return result, nil
}

func liuRiDayunLiunian(core Core, nianTime LiuNianTime) (*ganzhi.Zhu, *ganzhi.Zhu, error) {
	ds := dayunStepForYear(core.DaYun, core.BirthYear, nianTime.LichunYear)
	var dzZhu *ganzhi.Zhu
	if ds != nil {
		dzZhu = &ganzhi.Zhu{Gan: ds.Gan, Zhi: ds.Zhi}
	}
	ln, err := ComputeLiuNianAtTime(core, nianTime)
	if err != nil {
		return nil, nil, fmt.Errorf("computeLiuRi: liunian: %w", err)
	}
	var lnZhu *ganzhi.Zhu
	if ln != nil {
		lnZhu = &ganzhi.Zhu{Gan: ln.NianGan, Zhi: ln.NianZhi}
	}
	return dzZhu, lnZhu, nil
}

// ComputeLiuRi computes the day pillar and its interactions with the bazi chart.
func ComputeLiuRi(core Core, t LiuRiTime) (*LiuRi, error) {
	nianTime, err := ResolveLiuNianTime(t.Date)
	if err != nil {
		return nil, fmt.Errorf("computeLiuRi: liunian time: %w", err)
	}
	dzZhu, lnZhu, err := liuRiDayunLiunian(core, nianTime)
	if err != nil {
		return nil, err
	}
	return computeLiuRi(core.ToBazi(), t, dzZhu, lnZhu)
}

// ComputeLiuShi computes the hour pillar and its interactions with the bazi chart.
func ComputeLiuShi(core Core, t LiuShiTime) (*LiuShi, error) {
	refs, err := liuShiRefs(core, t)
	if err != nil {
		return nil, err
	}
	return computeLiuShi(core.ToBazi(), t, refs)
}

// ComputeLiuShiAt computes the flow hour directly from a standard time. It
// resolves the true-solar hour boundary internally and fills the ShiChen
// location facts, so callers only pass the standard time and need not
// pre-resolve a LiuShiTime or assemble boundary fields themselves.
func ComputeLiuShiAt(core Core, standard time.Time, loc *time.Location, longitude float64) (*LiuShi, error) {
	t, err := ResolveLiuShiTime(standard, loc, longitude)
	if err != nil {
		return nil, fmt.Errorf("computeLiuShi: resolve: %w", err)
	}
	refs, err := liuShiRefs(core, t)
	if err != nil {
		return nil, fmt.Errorf("computeLiuShi: refs: %w", err)
	}
	result, err := computeLiuShi(core.ToBazi(), t, refs)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		loc = time.UTC
	}
	result.ShiChen = LiuShiChen{
		StandardTime: t.Standard.In(loc).Format(time.RFC3339Nano),
		StandardDate: t.StandardDate.String(),
		Shichen:      ganzhi.DecodeZhi(t.Shichen),
	}
	return result, nil
}

// dayunStepForYear returns the DaYun step governing the given year, based on
// 虚岁 = year - birthYear + 1 matching the step's qi_sui/zhi_sui (虚岁) interval.
// Returns nil when: dy is nil / empty, birthYear unknown (0), year earlier than
// birth year, not yet in 大运, or past
// all steps.
func dayunStepForYear(dy *DaYun, birthYear, year int) *DaYunStep {
	if dy == nil || len(dy.Steps) == 0 || birthYear <= 0 || year < birthYear {
		return nil
	}
	for i := range dy.Steps {
		if year >= dy.Steps[i].StartYear && year <= dy.Steps[i].EndYear {
			return &dy.Steps[i]
		}
	}
	return nil
}

// ComputeXiaoYun computes the minor fortune (小运) by multiple schools.
func ComputeXiaoYun(core Core, maxAge int) []XiaoYunSet {
	return computeXiaoYun(core.ToBazi(), core.Gender, maxAge)
}

func liuShiRefs(core Core, t LiuShiTime) (flowRefs, error) {
	dayDate := t.LastDate
	if t.LateZi {
		dayDate = dayDate.AddDays(1)
	}
	dayPillar := tianwen.RiZhu(dayDate.Time())
	// 大运/流年按真太阳日所属立春年（与流日同权威）。
	nianForRefs, err := ResolveLiuNianTime(dayDate)
	if err != nil {
		return flowRefs{}, err
	}
	return resolveFlowRefs(core, nianForRefs, ganzhi.Zhu{Gan: dayPillar.Gan, Zhi: dayPillar.Zhi})
}
