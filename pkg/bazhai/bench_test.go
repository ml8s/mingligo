package bazhai

import (
	"testing"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

func BenchmarkComputeChart(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ComputeChart(ganzhi.Male, 1984)
	}
}

func BenchmarkComputeLayout(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ComputeLayout("兑", "乾", "坤", "艮"); err != nil {
			b.Fatal(err)
		}
	}
}
