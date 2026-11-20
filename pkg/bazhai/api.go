// Package bazhai provides八宅风水 computation.
//
// Types
//
//	Chart, MingGua, gua
//
// Functions
//
//	ComputeChart(gender Gender, birthYear int) → Chart
//	ComputeMingGua(gender Gender, birthYear int) → MingGua
package bazhai

import (
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

// ComputeChart computes八宅命盘 from a 立春 birth year and gender.
func ComputeChart(gender ganzhi.Gender, birthYear int) Chart {
	mg := ComputeMingGua(gender, birthYear)
	return Chart{
		MingGua:    mg,
		BaZhaiDirs: baZhaiDirectionsForGua(mg.Gua.Index),
		YearStars:  computeYearStars(birthYear),
	}
}

// ComputeChartAt derives the 立春 year from a standard time and
// computes the 八宅 chart. The 立春 boundary is an absolute instant, so the
// observer's timezone does not change the year.
func ComputeChartAt(gender ganzhi.Gender, standard time.Time) (Chart, error) {
	ly, err := tianwen.LichunYearAt(standard)
	if err != nil {
		return Chart{}, err
	}
	return ComputeChart(gender, ly.Year), nil
}
