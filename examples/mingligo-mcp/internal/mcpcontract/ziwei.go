package mcpcontract

import (
	"context"
	"fmt"
	"time"

	"github.com/ml8s/mingligo/pkg/tianwen"
	"github.com/ml8s/mingligo/pkg/ziwei"
)

func ziweiLunarProjector(loc *time.Location) (tianwen.LunarProjector, error) {
	return tianwen.NewLunarProjector(tianwen.LunarBasisLocalAstronomical, loc)
}

// domain decodes a ziwei.chart wire chart back into the domain chart. The wire
// DTO is a faithful projection of the domain JSON shape, so the round trip is
// lossless; the school/birth year checks keep the "empty chart" failure mode.

// ZiweiChart computes the twelve-palace chart from lunar birth data.
// ZiweiChart computes the canonical complete twelve-palace chart.
func ZiweiChart(_ context.Context, input ZiweiChartInput) (ziwei.Chart, error) {
	chart, err := input.Subject.ziweiChart()
	if err != nil {
		return ziwei.Chart{}, fmt.Errorf("ziwei.chart: %w", err)
	}
	return chart, nil
}

// ZiweiLiupan computes one selected flow layer from the calibrated natal
// subject. The natal chart is prepared once and reused across the requested
// batch.
func ZiweiLiupan(_ context.Context, input ZiweiLiupanInput) (out ZiweiLiupanOutput, err error) {
	if err := validateZiweiLiupan(input); err != nil {
		return ZiweiLiupanOutput{}, fmt.Errorf("ziwei.liupan: %w", err)
	}
	count := seriesCountValue(input.Count)
	chart, err := input.Subject.ziweiChart()
	if err != nil {
		return ZiweiLiupanOutput{}, fmt.Errorf("ziwei.liupan: %w", err)
	}
	current := input

	switch input.Scope {
	case "liu_nian":
		items := make([]ziwei.LiuNian, 0, count)
		for range count {
			item, err := computeZiweiLiuNian(current, chart)
			if err != nil {
				return out, err
			}
			items = append(items, item)
			current, err = advanceZiweiLiuNian(current)
			if err != nil {
				return out, err
			}
		}
		out = ZiweiLiupanOutput{LiuNian: items}
		return out, nil
	case "liu_yue":
		items := make([]ziwei.LiuYue, 0, count)
		for range count {
			item, err := computeZiweiLiuYue(current, chart)
			if err != nil {
				return out, err
			}
			items = append(items, item)
			current, err = advanceZiweiLiuYue(current)
			if err != nil {
				return out, err
			}
		}
		out = ZiweiLiupanOutput{LiuYue: items}
		return out, nil
	case "liu_ri":
		items := make([]ziwei.LiuRi, 0, count)
		for range count {
			item, err := computeZiweiLiuRi(current, chart)
			if err != nil {
				return out, err
			}
			items = append(items, item)
			current, err = advanceZiweiLiuRi(current)
			if err != nil {
				return out, err
			}
		}
		out = ZiweiLiupanOutput{LiuRi: items}
		return out, nil
	case "liu_shi":
		items := make([]ziwei.LiuShi, 0, count)
		for range count {
			item, err := computeZiweiLiuShi(current, chart)
			if err != nil {
				return out, err
			}
			items = append(items, item)
			current, err = advanceZiweiLiuShi(current)
			if err != nil {
				return out, err
			}
		}
		out = ZiweiLiupanOutput{LiuShi: items}
		return out, nil
	default:
		return ZiweiLiupanOutput{}, fmt.Errorf("ziwei.liupan: invalid scope %q", input.Scope)
	}
}

func validateZiweiLiupan(input ZiweiLiupanInput) error {
	if input.StandardTime == "" {
		return fmt.Errorf("standard_time is required")
	}
	switch input.Scope {
	case "liu_nian", "liu_yue", "liu_ri":
		if input.Longitude != nil {
			return fmt.Errorf("longitude is not allowed")
		}
	case "liu_shi":
		if input.Longitude == nil {
			return fmt.Errorf("longitude is required")
		}
		if *input.Longitude < -180 || *input.Longitude > 180 {
			return fmt.Errorf("longitude must be within [-180,180], got %v", *input.Longitude)
		}
	default:
		return fmt.Errorf("invalid scope %q", input.Scope)
	}
	return nil
}

func parseZiweiLiuNianTime(input ZiweiLiupanInput) (ziwei.LiuNianTime, tianwen.Date, time.Time, *time.Location, tianwen.LunarProjector, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return ziwei.LiuNianTime{}, tianwen.Date{}, time.Time{}, nil, tianwen.LunarProjector{}, err
	}
	projector, err := ziweiLunarProjector(loc)
	if err != nil {
		return ziwei.LiuNianTime{}, tianwen.Date{}, time.Time{}, nil, tianwen.LunarProjector{}, err
	}
	date := projector.SelectorDate(standard)
	resolved, err := ziwei.ResolveLiuNianTimeWith(projector, date)
	return resolved, date, standard, loc, projector, err
}

func computeZiweiLiuNian(input ZiweiLiupanInput, chart ziwei.Chart) (ziwei.LiuNian, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return ziwei.LiuNian{}, fmt.Errorf("ziwei.liupan liu_nian: %w", err)
	}
	result, err := ziwei.ComputeLiuNianAt(chart.FlowContext(), standard, loc)
	if err != nil {
		return ziwei.LiuNian{}, fmt.Errorf("ziwei.liupan liu_nian: %w", err)
	}
	return result, nil
}

func advanceZiweiLiuNian(input ZiweiLiupanInput) (ZiweiLiupanInput, error) {
	resolved, _, _, _, projector, err := parseZiweiLiuNianTime(input)
	if err != nil {
		return input, fmt.Errorf("ziwei.liupan liu_nian: %w", err)
	}
	next, err := ziwei.NextLiuNianTimeWith(projector, resolved)
	if err != nil {
		return input, fmt.Errorf("ziwei.liupan liu_nian: %w", err)
	}
	input.StandardTime = standardTimeValue(projector.SelectorInstant(next.Start))
	return input, nil
}

func parseZiweiLiuYueTime(input ZiweiLiupanInput) (ziwei.LiuYueTime, tianwen.Date, time.Time, *time.Location, tianwen.LunarProjector, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return ziwei.LiuYueTime{}, tianwen.Date{}, time.Time{}, nil, tianwen.LunarProjector{}, err
	}
	projector, err := ziweiLunarProjector(loc)
	if err != nil {
		return ziwei.LiuYueTime{}, tianwen.Date{}, time.Time{}, nil, tianwen.LunarProjector{}, err
	}
	date := projector.SelectorDate(standard)
	resolved, err := ziwei.ResolveLiuYueTimeWith(projector, date)
	return resolved, date, standard, loc, projector, err
}

func computeZiweiLiuYue(input ZiweiLiupanInput, chart ziwei.Chart) (ziwei.LiuYue, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return ziwei.LiuYue{}, fmt.Errorf("ziwei.liupan liu_yue: %w", err)
	}
	result, err := ziwei.ComputeLiuYueAt(chart.FlowContext(), standard, loc)
	if err != nil {
		return ziwei.LiuYue{}, fmt.Errorf("ziwei.liupan liu_yue: %w", err)
	}
	return result, nil
}

func advanceZiweiLiuYue(input ZiweiLiupanInput) (ZiweiLiupanInput, error) {
	resolved, _, _, _, projector, err := parseZiweiLiuYueTime(input)
	if err != nil {
		return input, fmt.Errorf("ziwei.liupan liu_yue: %w", err)
	}
	next, err := ziwei.NextLiuYueTimeWith(projector, resolved)
	if err != nil {
		return input, fmt.Errorf("ziwei.liupan liu_yue: %w", err)
	}
	input.StandardTime = standardTimeValue(projector.SelectorInstant(next.Date))
	return input, nil
}

func parseZiweiLiuRiTime(input ZiweiLiupanInput) (ziwei.LiuRiTime, tianwen.LunarProjector, time.Time, *time.Location, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return ziwei.LiuRiTime{}, tianwen.LunarProjector{}, time.Time{}, nil, err
	}
	projector, err := ziweiLunarProjector(loc)
	if err != nil {
		return ziwei.LiuRiTime{}, tianwen.LunarProjector{}, time.Time{}, nil, err
	}
	resolved, err := ziwei.ResolveLiuRiTimeAtWith(projector, standard, loc)
	return resolved, projector, standard, loc, err
}

func computeZiweiLiuRi(input ZiweiLiupanInput, chart ziwei.Chart) (ziwei.LiuRi, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return ziwei.LiuRi{}, fmt.Errorf("ziwei.liupan liu_ri: %w", err)
	}
	result, err := ziwei.ComputeLiuRiAt(chart.FlowContext(), standard, loc)
	if err != nil {
		return ziwei.LiuRi{}, fmt.Errorf("ziwei.liupan liu_ri: %w", err)
	}
	return result, nil
}

func advanceZiweiLiuRi(input ZiweiLiupanInput) (ZiweiLiupanInput, error) {
	resolved, projector, _, loc, err := parseZiweiLiuRiTime(input)
	if err != nil {
		return input, fmt.Errorf("ziwei.liupan liu_ri: %w", err)
	}
	next, err := ziwei.NextLiuRiTimeWith(projector, resolved)
	if err != nil {
		return input, fmt.Errorf("ziwei.liupan liu_ri: %w", err)
	}
	input.StandardTime = standardTimeValue(next.Standard.In(loc))
	return input, nil
}

func parseZiweiLiuShiTime(input ZiweiLiupanInput) (ziwei.LiuShiTime, tianwen.LunarProjector, time.Time, *time.Location, error) {
	if input.Longitude == nil {
		return ziwei.LiuShiTime{}, tianwen.LunarProjector{}, time.Time{}, nil, fmt.Errorf("longitude is required")
	}
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return ziwei.LiuShiTime{}, tianwen.LunarProjector{}, time.Time{}, nil, err
	}
	projector, err := ziweiLunarProjector(loc)
	if err != nil {
		return ziwei.LiuShiTime{}, tianwen.LunarProjector{}, time.Time{}, nil, err
	}
	resolved, err := ziwei.ResolveLiuShiTimeWith(projector, standard, loc, float64(*input.Longitude))
	return resolved, projector, standard, loc, err
}

func computeZiweiLiuShi(input ZiweiLiupanInput, chart ziwei.Chart) (ziwei.LiuShi, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return ziwei.LiuShi{}, fmt.Errorf("ziwei.liupan liu_shi: %w", err)
	}
	result, err := ziwei.ComputeLiuShiAt(chart.FlowContext(), standard, loc, float64(*input.Longitude))
	if err != nil {
		return ziwei.LiuShi{}, fmt.Errorf("ziwei.liupan liu_shi: %w", err)
	}
	return result, nil
}

func advanceZiweiLiuShi(input ZiweiLiupanInput) (ZiweiLiupanInput, error) {
	resolved, projector, _, loc, err := parseZiweiLiuShiTime(input)
	if err != nil {
		return input, fmt.Errorf("ziwei.liupan liu_shi: %w", err)
	}
	next, err := ziwei.NextLiuShiTimeWith(projector, resolved)
	if err != nil {
		return input, fmt.Errorf("ziwei.liupan liu_shi: %w", err)
	}
	input.StandardTime = standardTimeValue(next.Standard.In(loc))
	return input, nil
}

// ZiweiBond returns raw cross-chart facts for two charts.
func ZiweiBond(_ context.Context, input ZiweiBondInput) (ziwei.Bond, error) {
	chartA, err := input.A.Subject.ziweiChart()
	if err != nil {
		return ziwei.Bond{}, fmt.Errorf("ziwei.bond: chart a: %w", err)
	}
	chartB, err := input.B.Subject.ziweiChart()
	if err != nil {
		return ziwei.Bond{}, fmt.Errorf("ziwei.bond: chart b: %w", err)
	}
	return ziwei.ComputeBond(chartA.BondSideFacts(), chartB.BondSideFacts()), nil
}
