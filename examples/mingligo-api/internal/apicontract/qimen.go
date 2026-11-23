package apicontract

import (
	"context"
	"fmt"

	"github.com/ml8s/mingligo/pkg/qimen"
)

// QimenYongShenSymbols is the 用神 symbol combination accepted by qimen.chart.
// It is not a plain []string because the legacy params schema declared
// `minItems: 1` and `uniqueItems: true`; the named type carries those
// constraints through the registry's TypeSchemas (an empty or duplicated list
// used to be a schema error, and silently ignoring an empty list would drop the
// 用神 analysis the caller asked for).
type QimenYongShenSymbols []string

// QimenChart computes a qimen chart for an explicitly selected method.
func QimenChart(_ context.Context, input QimenChartInput) (QimenChartOutput, error) {
	at, _, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return QimenChartOutput{}, fmt.Errorf("qimen.chart: %w", err)
	}
	if float64(input.Longitude) < -180 || float64(input.Longitude) > 180 {
		return QimenChartOutput{}, fmt.Errorf("qimen.chart: longitude must be within [-180,180], got %v", float64(input.Longitude))
	}
	method, err := qimen.ResolveMethodWithOptions(
		string(input.Method.Scope), string(input.Method.School), string(input.Method.DingjuMethod),
		string(input.Method.QuarterRule), string(input.Method.BaseDingjuMethod),
		string(input.Method.DunSource), string(input.Method.HourBoundary),
	)
	if err != nil {
		return QimenChartOutput{}, fmt.Errorf("qimen.chart: %w", err)
	}
	if method.School == qimen.SchoolJinhanYuJing {
		if len(input.YongShen) > 0 || input.Subject != nil {
			return QimenChartOutput{}, fmt.Errorf("qimen.chart: jinhan_yujing does not support yong_shen or subject")
		}
		jinhan := qimen.ComputeJinhanChartAt(at, float64(input.Longitude)).View()
		return QimenChartOutput{Jinhan: &jinhan}, nil
	}
	var birthDate qimen.BirthDate
	if input.Subject != nil {
		birthAt, err := input.Subject.standardTime()
		if err != nil {
			return QimenChartOutput{}, fmt.Errorf("qimen.chart: subject: %w", err)
		}
		birthDate = qimen.BirthDate{Time: birthAt, Precision: qimen.BirthDateMoment, Has: true}
	}
	if len(input.YongShen) > 0 || birthDate.Has {
		symbols := make([]qimen.YongShenSymbol, 0, len(input.YongShen))
		for _, name := range input.YongShen {
			symbol, err := qimen.EncodeYongShen(name)
			if err != nil {
				return QimenChartOutput{}, fmt.Errorf("qimen.chart: %w", err)
			}
			symbols = append(symbols, symbol)
		}
		result, err := qimen.ComputeChartWithYongShenAndMethodAt(at, float64(input.Longitude), symbols, birthDate, method)
		if err != nil {
			return QimenChartOutput{}, fmt.Errorf("qimen.chart: %w", err)
		}
		view := result.View()
		return QimenChartOutput{Standard: &view}, nil
	}
	result, err := qimen.ComputeChartWithMethodAt(at, float64(input.Longitude), method)
	if err != nil {
		return QimenChartOutput{}, fmt.Errorf("qimen.chart: %w", err)
	}
	view := result.View()
	return QimenChartOutput{Standard: &view}, nil
}
