package mcpcontract

import (
	"context"
	"fmt"
	"time"

	"github.com/ml8s/mingligo/pkg/bazi"
)

// BaziChartInput consumes a resolved natal subject. The chart domain derives
// year/month from solar-term events and day/hour from LAST + 子时.
type BaziChartInput struct {
	Subject NatalSubject `json:"subject" jsonschema:"schema 合法的出生主体；可来自用户明确事实或用户确认后的 subject_candidates 候选。"`
}

// BaziChartCompute returns the canonical, domain-aggregated natal chart view.
func BaziChart(_ context.Context, input BaziChartInput) (bazi.ChartView, error) {
	chart, err := input.Subject.baziChart()
	if err != nil {
		return bazi.ChartView{}, fmt.Errorf("bazi.chart: %w", err)
	}
	return chart.View(), nil
}

// ── typed handlers for the remaining bazi tools ──

// BaziBond returns cross-chart facts for two bazi charts.
func BaziBond(_ context.Context, input BaziBondInput) (bazi.BondView, error) {
	coreA, err := input.A.Subject.baziCore()
	if err != nil {
		return bazi.BondView{}, fmt.Errorf("bazi.bond: chart a: %w", err)
	}
	coreB, err := input.B.Subject.baziCore()
	if err != nil {
		return bazi.BondView{}, fmt.Errorf("bazi.bond: chart b: %w", err)
	}
	return bazi.ComputeBond(coreA, coreB).View(), nil
}

// BaziLiupan computes one selected dynamic time layer from the calibrated natal
// subject. The five layers share one core preparation; their calendar advance
// rules remain domain-specific.
func BaziLiupan(_ context.Context, input BaziLiupanInput) (out BaziLiupanOutput, err error) {
	if err := validateBaziLiupan(input); err != nil {
		return out, fmt.Errorf("bazi.liupan: %w", err)
	}
	count := seriesCountValue(input.Count)
	core, err := input.Subject.baziCore()
	if err != nil {
		return out, fmt.Errorf("bazi.liupan: %w", err)
	}
	current := input

	switch input.Scope {
	case "liu_nian":
		items := make([]bazi.LiuNian, 0, count)
		for range count {
			item, err := computeBaziLiuNian(current, core)
			if err != nil {
				return out, err
			}
			items = append(items, item)
			current, err = advanceBaziLiuNian(current)
			if err != nil {
				return out, err
			}
		}
		out = BaziLiupanOutput{LiuNian: items}
		return out, nil
	case "liu_yue":
		items := make([]bazi.LiuYue, 0, count)
		for range count {
			item, err := computeBaziLiuYue(current, core)
			if err != nil {
				return out, err
			}
			items = append(items, item)
			current, err = advanceBaziLiuYue(current)
			if err != nil {
				return out, err
			}
		}
		out = BaziLiupanOutput{LiuYue: items}
		return out, nil
	case "liu_ri":
		items := make([]bazi.LiuRi, 0, count)
		for range count {
			item, err := computeBaziLiuRi(current, core)
			if err != nil {
				return out, err
			}
			items = append(items, item)
			current, err = advanceBaziLiuRi(current)
			if err != nil {
				return out, err
			}
		}
		out = BaziLiupanOutput{LiuRi: items}
		return out, nil
	case "liu_shi":
		items := make([]bazi.LiuShi, 0, count)
		for range count {
			item, err := computeBaziLiuShi(current, core)
			if err != nil {
				return out, err
			}
			items = append(items, item)
			current, err = advanceBaziLiuShi(current)
			if err != nil {
				return out, err
			}
		}
		out = BaziLiupanOutput{LiuShi: items}
		return out, nil
	case "xiao_yun":
		out = BaziLiupanOutput{XiaoYun: bazi.ComputeXiaoYun(core, count)}
		return out, nil
	default:
		return BaziLiupanOutput{}, fmt.Errorf("bazi.liupan: invalid scope %q", input.Scope)
	}
}

func validateBaziLiupan(input BaziLiupanInput) error {
	if input.Longitude != nil && (*input.Longitude < -180 || *input.Longitude > 180) {
		return fmt.Errorf("longitude must be within [-180,180], got %v", *input.Longitude)
	}
	switch input.Scope {
	case "liu_nian", "liu_yue", "liu_ri", "liu_shi":
		if input.StandardTime == "" {
			return fmt.Errorf("standard_time is required")
		}
		if input.Scope == "liu_ri" || input.Scope == "liu_shi" {
			if input.Longitude == nil {
				return fmt.Errorf("longitude is required")
			}
		} else if input.Longitude != nil {
			return fmt.Errorf("longitude is not allowed")
		}
	case "xiao_yun":
		if input.StandardTime != "" {
			return fmt.Errorf("standard_time is not allowed")
		}
		if input.Longitude != nil {
			return fmt.Errorf("longitude is not allowed")
		}
	default:
		return fmt.Errorf("invalid scope %q", input.Scope)
	}
	return nil
}

func parseBaziLiuNianTime(input BaziLiupanInput) (bazi.LiuNianTime, *time.Location, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return bazi.LiuNianTime{}, nil, err
	}
	resolved, err := bazi.ResolveLiuNianTimeAt(standard, loc)
	return resolved, loc, err
}

func computeBaziLiuNian(input BaziLiupanInput, core bazi.Core) (bazi.LiuNian, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return bazi.LiuNian{}, fmt.Errorf("bazi.liupan liu_nian: %w", err)
	}
	result, err := bazi.ComputeLiuNianAt(core, standard, loc)
	if err != nil {
		return bazi.LiuNian{}, fmt.Errorf("bazi.liupan liu_nian: %w", err)
	}
	return *result, nil
}

func advanceBaziLiuNian(input BaziLiupanInput) (BaziLiupanInput, error) {
	resolved, loc, err := parseBaziLiuNianTime(input)
	if err != nil {
		return input, fmt.Errorf("bazi.liupan liu_nian: %w", err)
	}
	next, err := bazi.NextLiuNianTimeAt(resolved, loc)
	if err != nil {
		return input, fmt.Errorf("bazi.liupan liu_nian: %w", err)
	}
	input.StandardTime = standardTimeValue(next.StartInstant.In(loc))
	return input, nil
}

func parseBaziLiuYueTime(input BaziLiupanInput) (bazi.LiuYueTime, *time.Location, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return bazi.LiuYueTime{}, nil, err
	}
	resolved, err := bazi.ResolveLiuYueTimeAt(standard, loc)
	return resolved, loc, err
}

func computeBaziLiuYue(input BaziLiupanInput, core bazi.Core) (bazi.LiuYue, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return bazi.LiuYue{}, fmt.Errorf("bazi.liupan liu_yue: %w", err)
	}
	result, err := bazi.ComputeLiuYueAt(core, standard, loc)
	if err != nil {
		return bazi.LiuYue{}, fmt.Errorf("bazi.liupan liu_yue: %w", err)
	}
	return *result, nil
}

func advanceBaziLiuYue(input BaziLiupanInput) (BaziLiupanInput, error) {
	resolved, loc, err := parseBaziLiuYueTime(input)
	if err != nil {
		return input, fmt.Errorf("bazi.liupan liu_yue: %w", err)
	}
	next, err := bazi.NextLiuYueTimeAt(resolved, loc)
	if err != nil {
		return input, fmt.Errorf("bazi.liupan liu_yue: %w", err)
	}
	input.StandardTime = standardTimeValue(next.StartInstant.In(loc))
	return input, nil
}

func parseBaziLiuRiTime(input BaziLiupanInput) (bazi.LiuRiTime, *time.Location, error) {
	if input.Longitude == nil {
		return bazi.LiuRiTime{}, nil, fmt.Errorf("longitude is required")
	}
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return bazi.LiuRiTime{}, nil, err
	}
	resolved, err := bazi.ResolveLiuRiTimeAt(standard, loc, float64(*input.Longitude))
	return resolved, loc, err
}

func computeBaziLiuRi(input BaziLiupanInput, core bazi.Core) (bazi.LiuRi, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return bazi.LiuRi{}, fmt.Errorf("bazi.liupan liu_ri: %w", err)
	}
	result, err := bazi.ComputeLiuRiAt(core, standard, loc, float64(*input.Longitude))
	if err != nil {
		return bazi.LiuRi{}, fmt.Errorf("bazi.liupan liu_ri: %w", err)
	}
	return *result, nil
}

func advanceBaziLiuRi(input BaziLiupanInput) (BaziLiupanInput, error) {
	resolved, loc, err := parseBaziLiuRiTime(input)
	if err != nil {
		return input, fmt.Errorf("bazi.liupan liu_ri: %w", err)
	}
	next, err := bazi.NextLiuRiTime(resolved)
	if err != nil {
		return input, fmt.Errorf("bazi.liupan liu_ri: %w", err)
	}
	input.StandardTime = standardTimeValue(next.Standard.In(loc))
	return input, nil
}

func parseBaziLiuShiTime(input BaziLiupanInput) (bazi.LiuShiTime, *time.Location, error) {
	if input.Longitude == nil {
		return bazi.LiuShiTime{}, nil, fmt.Errorf("longitude is required")
	}
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return bazi.LiuShiTime{}, nil, err
	}
	resolved, err := bazi.ResolveLiuShiTime(standard, loc, float64(*input.Longitude))
	return resolved, loc, err
}

func computeBaziLiuShi(input BaziLiupanInput, core bazi.Core) (bazi.LiuShi, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return bazi.LiuShi{}, fmt.Errorf("bazi.liupan liu_shi: %w", err)
	}
	result, err := bazi.ComputeLiuShiAt(core, standard, loc, float64(*input.Longitude))
	if err != nil {
		return bazi.LiuShi{}, fmt.Errorf("bazi.liupan liu_shi: %w", err)
	}
	return *result, nil
}

func advanceBaziLiuShi(input BaziLiupanInput) (BaziLiupanInput, error) {
	resolved, loc, err := parseBaziLiuShiTime(input)
	if err != nil {
		return input, fmt.Errorf("bazi.liupan liu_shi: %w", err)
	}
	next, err := bazi.NextLiuShiTime(resolved)
	if err != nil {
		return input, fmt.Errorf("bazi.liupan liu_shi: %w", err)
	}
	input.StandardTime = standardTimeValue(next.Standard.In(loc))
	return input, nil
}
