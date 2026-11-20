package xuankong

import (
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/tianwen"
)

func BenchmarkComputeChart(b *testing.B) {
	st := tianwen.ApparentSolarFromCivil(
		time.Date(2020, 6, 1, 12, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		116.4, 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ComputeChartAt(st, 0, 12); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComputeLiuNian(b *testing.B) {
	st := tianwen.ApparentSolarFromCivil(
		time.Date(2020, 6, 1, 12, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		116.4, 8)
	chart, err := ComputeChartAt(st, 0, 12)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ComputeLiuNian(2027, &chart)
	}
}
