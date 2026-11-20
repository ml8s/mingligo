package bazi

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
	"testing"
)

func TestAnHe(t *testing.T) {
	// 暗合对: 寅丑, 卯申, 午亥, 子戌
	tests := []struct {
		name  string
		zhi   [4]ganzhi.Zhi
		wantA ganzhi.Zhi
		wantB ganzhi.Zhi
	}{
		{"子戌暗合", [4]ganzhi.Zhi{ganzhi.ZhiZi, ganzhi.ZhiChou, ganzhi.ZhiXu, ganzhi.ZhiChen}, ganzhi.ZhiZi, ganzhi.ZhiXu},
		{"寅丑暗合", [4]ganzhi.Zhi{ganzhi.ZhiYin, ganzhi.ZhiChou, ganzhi.ZhiZi, ganzhi.ZhiWu}, ganzhi.ZhiYin, ganzhi.ZhiChou},
		{"卯申暗合", [4]ganzhi.Zhi{ganzhi.ZhiZi, ganzhi.ZhiShen, ganzhi.ZhiMao, ganzhi.ZhiWu}, ganzhi.ZhiMao, ganzhi.ZhiShen},
		{"午亥暗合", [4]ganzhi.Zhi{ganzhi.ZhiZi, ganzhi.ZhiChou, ganzhi.ZhiWu, ganzhi.ZhiHai}, ganzhi.ZhiWu, ganzhi.ZhiHai},
		{"无暗合", [4]ganzhi.Zhi{ganzhi.ZhiZi, ganzhi.ZhiChou, ganzhi.ZhiChen, ganzhi.ZhiSi}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bz := ganzhi.Bazi{
				Nian: ganzhi.Zhu{Gan: ganzhi.GanJia, Zhi: tt.zhi[0]},
				Yue:  ganzhi.Zhu{Gan: ganzhi.GanBing, Zhi: tt.zhi[1]},
				Ri:   ganzhi.Zhu{Gan: ganzhi.GanWu, Zhi: tt.zhi[2]},
				Shi:  ganzhi.Zhu{Gan: ganzhi.GanGeng, Zhi: tt.zhi[3]},
			}
			detection := DetectRelations(canonicalChart(bz, ganzhi.Male, 2000).ToBazi())

			if tt.wantA == 0 {
				if len(detection.AnHe) > 0 {
					t.Errorf("expected no an_he, got %d", len(detection.AnHe))
				}
				return
			}
			found := false
			for _, ah := range detection.AnHe {
				if (ah.ZhiA == tt.wantA && ah.ZhiB == tt.wantB) || (ah.ZhiA == tt.wantB && ah.ZhiB == tt.wantA) {
					found = true
				}
			}
			if !found {
				t.Errorf("暗合 %s%s not found in an_he (%d entries)", tt.wantA, tt.wantB, len(detection.AnHe))
			}
		})
	}
}
