// Package qimen computes Qimen charts across hour, quarter, day, month, and
// year scopes with rotating, flying, and Golden Mirror plates.
package qimen

import (
	"fmt"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

// Method is a closed Qimen method combination.
type Method struct {
	Scope        Scope
	School       School
	Dingju       DingjuMethod
	QuarterRule  QuarterRule
	BaseDingju   DingjuMethod
	DunSource    DunSource
	HourBoundary HourBoundary
}

// EncodeMethod validates a method combination against the catalog.
func EncodeMethod(scope, school, dingjuMethod, quarterRule string) (Method, error) {
	return EncodeMethodWithOptions(scope, school, dingjuMethod, quarterRule, "", "", "")
}

// validateMethodFields validates an already-typed method combination against
// the catalog. It is the single source of method validation; string parsing
// lives only in EncodeMethodWithOptions at the input boundary.
func validateMethodFields(scope Scope, school School, dingju DingjuMethod, quarter QuarterRule, baseDingju DingjuMethod, dunSource DunSource, hourBoundary HourBoundary) (Method, error) {
	scopeEntry, ok := scopeMethods[scope]
	if !ok {
		return Method{}, fmt.Errorf("unknown scope %q", scope)
	}
	schoolEntry, ok := schoolMethods[school]
	if !ok {
		return Method{}, fmt.Errorf("unknown school %q", school)
	}
	if !schoolEntry.allowsScope(scope) {
		return Method{}, fmt.Errorf("school %q does not support scope %q", school, scope)
	}
	result := Method{Scope: scope, School: school}
	switch scope {
	case ScopeQuarter:
		if dingju != "" {
			return Method{}, fmt.Errorf("scope %q does not support dingju_method", scope)
		}
		if quarter == "" {
			return Method{}, fmt.Errorf("scope %q requires quarter_rule", scope)
		}
		if err := validateQuarterRule(quarter); err != nil {
			return Method{}, err
		}
		if !scopeEntry.allowsQuarterRule(quarter) || !schoolEntry.allowsQuarterRule(quarter) {
			return Method{}, fmt.Errorf("school %q does not support quarter rule %q", school, quarter)
		}
		result.QuarterRule = quarter
		if quarter == QuarterTwelveMinuteTenDivision {
			base := baseDingju
			if base == "" {
				base = DingjuChaiBu
			}
			if err := validateDingju(base); err != nil {
				return Method{}, fmt.Errorf("invalid base dingju method: %w", err)
			}
			if !quarterRules[quarter].BaseDingjuMethods[base] {
				return Method{}, fmt.Errorf("twelve-minute quarter does not support base dingju method %q", base)
			}
			if dunSource == 0 {
				dunSource = quarterRules[quarter].DefaultDunSource
			}
			if !quarterRules[quarter].DunSources[dunSource] {
				return Method{}, fmt.Errorf("twelve-minute quarter does not support dun source %q", dunSource)
			}
			if hourBoundary == 0 {
				hourBoundary = quarterRules[quarter].DefaultHourBoundary
			}
			if _, exists := quarterRules[quarter].HourBoundaries[hourBoundary]; !exists {
				return Method{}, fmt.Errorf("twelve-minute quarter does not support hour boundary %q", hourBoundary)
			}
			result.BaseDingju, result.DunSource, result.HourBoundary = base, dunSource, hourBoundary
		} else if baseDingju != "" || dunSource != 0 || hourBoundary != 0 {
			return Method{}, fmt.Errorf("quarter rule %q does not support Kejia options", quarter)
		}
	default:
		if quarter != "" {
			return Method{}, fmt.Errorf("scope %q does not support quarter_rule", scope)
		}
		if baseDingju != "" || dunSource != 0 || hourBoundary != 0 {
			return Method{}, fmt.Errorf("scope %q does not support Kejia options", scope)
		}
		if scopeEntry.allowsNoDingjuSchool(school) {
			if dingju != "" {
				return Method{}, fmt.Errorf("school %q must omit dingju_method", school)
			}
			result.Dingju = DingjuNone
		} else if len(scopeEntry.DingjuMethods) == 0 {
			if dingju != "" {
				return Method{}, fmt.Errorf("scope %q does not support dingju_method", scope)
			}
		} else {
			if err := validateDingju(dingju); err != nil {
				return Method{}, err
			}
			if !scopeEntry.allowsDingju(dingju) {
				return Method{}, fmt.Errorf("scope %q does not support dingju method %q", scope, dingju)
			}
			result.Dingju = dingju
		}
		if result.Dingju == DingjuNone {
			if !scopeEntry.allowsNoDingjuSchool(school) || !schoolEntry.allowsDingju(DingjuNone) {
				return Method{}, fmt.Errorf("school %q must omit dingju_method", school)
			}
		} else if result.Dingju != "" && !schoolEntry.allowsDingju(result.Dingju) {
			return Method{}, fmt.Errorf("school %q does not support dingju method %q", school, result.Dingju)
		}
	}
	return result, nil
}

func validateQuarterRule(r QuarterRule) error {
	switch r {
	case QuarterTenMinuteSanYuan, QuarterTwelveMinuteTenDivision:
		return nil
	default:
		return fmt.Errorf("unknown quarter rule %q", r)
	}
}

func validateDingju(d DingjuMethod) error {
	switch d {
	case DingjuChaiBu, DingjuZhiRun, DingjuMaoShan, DingjuNone:
		return nil
	default:
		return fmt.Errorf("unknown dingju method %q", d)
	}
}

// EncodeMethodWithOptions validates a method and its Kejia-specific options.
func EncodeMethodWithOptions(
	scope, school, dingjuMethod, quarterRule, baseDingjuMethod, dunSource, hourBoundary string,
) (Method, error) {
	parsedScope, err := EncodeScope(scope)
	if err != nil {
		return Method{}, err
	}
	parsedSchool, err := EncodeSchool(school)
	if err != nil {
		return Method{}, err
	}
	parsedDingju := DingjuMethod(dingjuMethod)
	if dingjuMethod != "" {
		parsedDingju, err = EncodeDingjuMethod(dingjuMethod)
		if err != nil {
			return Method{}, err
		}
	}
	parsedQuarter := QuarterRule(quarterRule)
	if quarterRule != "" {
		parsedQuarter, err = EncodeQuarterRule(quarterRule)
		if err != nil {
			return Method{}, err
		}
	}
	parsedBase := DingjuMethod(baseDingjuMethod)
	if baseDingjuMethod != "" {
		parsedBase, err = EncodeDingjuMethod(baseDingjuMethod)
		if err != nil {
			return Method{}, err
		}
	}
	var parsedDun DunSource
	if dunSource != "" {
		var ok bool
		parsedDun, ok = EncodeDunSource(dunSource)
		if !ok {
			return Method{}, fmt.Errorf("unknown dun source %q", dunSource)
		}
	}
	var parsedBoundary HourBoundary
	if hourBoundary != "" {
		var ok bool
		parsedBoundary, ok = EncodeHourBoundary(hourBoundary)
		if !ok {
			return Method{}, fmt.Errorf("unknown hour boundary %q", hourBoundary)
		}
	}
	return validateMethodFields(
		parsedScope, parsedSchool, parsedDingju, parsedQuarter, parsedBase, parsedDun, parsedBoundary,
	)
}

// ResolveMethod fills omitted API dimensions from the catalog and validates the
// resulting closed combination.
func ResolveMethod(scope, school, dingjuMethod, quarterRule string) (Method, error) {
	if scope == "" {
		scope = DecodeScope(defaultScope)
	}
	if school == "" {
		school = DecodeSchool(defaultSchool)
	}
	parsedScope, err := EncodeScope(scope)
	if err != nil {
		return Method{}, err
	}
	parsedSchool, err := EncodeSchool(school)
	if err != nil {
		return Method{}, err
	}
	entry := scopeMethods[parsedScope]
	if parsedScope == ScopeQuarter {
		if quarterRule == "" {
			quarterRule = DecodeQuarterRule(entry.DefaultQuarterRule)
		}
		return EncodeMethod(scope, school, "", quarterRule)
	}
	if dingjuMethod == "" && entry.allowsNoDingjuSchool(parsedSchool) {
		return EncodeMethod(scope, school, "", "")
	}
	if dingjuMethod == "" && entry.DefaultDingju != "" {
		dingjuMethod = DecodeDingjuMethod(entry.DefaultDingju)
	}
	return EncodeMethod(scope, school, dingjuMethod, "")
}

func ResolveMethodWithOptions(
	scope, school, dingjuMethod, quarterRule, baseDingjuMethod, dunSource, hourBoundary string,
) (Method, error) {
	resolved, err := ResolveMethod(scope, school, dingjuMethod, quarterRule)
	if err != nil {
		return Method{}, err
	}
	// resolved already carries typed Scope/School/QuarterRule/Dingju; keep them
	// symbolic instead of Decode→Encode roundtripping. DingjuNone means "omitted"
	// at this validation boundary, so pass it as the zero value.
	dingju := resolved.Dingju
	if dingju == DingjuNone {
		dingju = ""
	}
	var base DingjuMethod
	if baseDingjuMethod != "" {
		var err error
		base, err = EncodeDingjuMethod(baseDingjuMethod)
		if err != nil {
			return Method{}, err
		}
	}
	var dun DunSource
	if dunSource != "" {
		var ok bool
		dun, ok = EncodeDunSource(dunSource)
		if !ok {
			return Method{}, fmt.Errorf("unknown dun source %q", dunSource)
		}
	}
	var boundary HourBoundary
	if hourBoundary != "" {
		var ok bool
		boundary, ok = EncodeHourBoundary(hourBoundary)
		if !ok {
			return Method{}, fmt.Errorf("unknown hour boundary %q", hourBoundary)
		}
	}
	return validateMethodFields(
		resolved.Scope, resolved.School, dingju, resolved.QuarterRule, base, dun, boundary,
	)
}

// computeChartWithMethod computes a chart for a catalog method and rejects
// combinations absent from the catalog.
func computeChartWithMethod(gt time.Time, st time.Time, method Method) (Chart, error) {
	parsed, err := validateMethodFields(
		method.Scope, method.School, method.Dingju, method.QuarterRule,
		method.BaseDingju, method.DunSource, method.HourBoundary,
	)
	if err != nil {
		return Chart{}, fmt.Errorf("invalid qimen method: %w", err)
	}
	if parsed.School == SchoolJinhanYuJing {
		return Chart{}, fmt.Errorf("jinhan_yujing has a distinct result model; use ComputeJinhanChart")
	}
	bz, termTime, chartTime := chartInputs(gt, st)
	return computeChart(bz, termTime, chartTime, parsed), nil
}

// chartInputs 返回 (四柱, 节气判定时刻, 盘面时刻)：
//   - 节气判定用 Gregorian time（节气是绝对瞬时，换局/置闰/茅山元与出生 Gregorian time同帧比较）；
//   - 盘面分钟分法（十二分钟/十分钟）用真太阳时壁钟，与 bz.Shi（真太阳时辰）一致。
func chartInputs(gt time.Time, st time.Time) (ganzhi.Bazi, time.Time, time.Time) {
	bz := tianwen.ComputeBazi(gt, st)
	if qimenDayBoundary == "late_zi_rolls_day_pillar" && st.Hour() >= 23 {
		dayIndex := ganzhi.SixtyCycleIndex(bz.Ri.Gan, bz.Ri.Zhi) + 1
		bz.Ri = ganzhi.SixtyToZhu(dayIndex)
	}
	return bz, gt, st
}

// computeChartWithYongShenAndMethod aggregates selected Yong Shen under a
// catalog method and rejects combinations absent from the catalog.
func computeChartWithYongShenAndMethod(
	gt time.Time, st time.Time, syms []YongShenSymbol, birthDate BirthDate, method Method,
) (Chart, error) {
	parsed, err := validateMethodFields(
		method.Scope, method.School, method.Dingju, method.QuarterRule,
		method.BaseDingju, method.DunSource, method.HourBoundary,
	)
	if err != nil {
		return Chart{}, err
	}
	if parsed.School == SchoolJinhanYuJing {
		return Chart{}, fmt.Errorf("jinhan_yujing has a distinct result model; use ComputeJinhanChart")
	}
	bz, termTime, chartTime := chartInputs(gt, st)
	return computeChartWithYongShen(bz, termTime, chartTime, parsed, syms, birthDate), nil
}
