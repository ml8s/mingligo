package bazi

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
	"testing"
)

func TestBanHe_19810826(t *testing.T) {
	bz := ganzhi.Bazi{
		Nian: ganzhi.Zhu{Gan: ganzhi.GanXin, Zhi: ganzhi.ZhiYou},
		Yue:  ganzhi.Zhu{Gan: ganzhi.GanBing, Zhi: ganzhi.ZhiShen},
		Ri:   ganzhi.Zhu{Gan: ganzhi.GanBing, Zhi: ganzhi.ZhiZi},
		Shi:  ganzhi.Zhu{Gan: ganzhi.GanJi, Zhi: ganzhi.ZhiChou},
	}
	chart := canonicalChart(bz, ganzhi.Male, 1981)
	detection := DetectRelations(chart.ToBazi())

	// Oracle detects: 酉丑半合金局 (year+hour), 申子半合水局 (month+day)
	if len(detection.SanHePartial) != 2 {
		t.Fatalf("san_he_partial count = %d, want 2 (酉丑半合金 + 申子半合水)", len(detection.SanHePartial))
	}

	foundYouChou := false
	foundShenZi := false
	for _, ph := range detection.SanHePartial {
		branches := ph.Branches
		if len(branches) == 2 {
			if (branches[0] == ganzhi.ZhiYou && branches[1] == ganzhi.ZhiChou) || (branches[0] == ganzhi.ZhiChou && branches[1] == ganzhi.ZhiYou) {
				foundYouChou = true
				if ph.Element != ganzhi.WxJin {
					t.Errorf("酉丑半合 element = %s, want 金", ph.Element)
				}
			}
			if (branches[0] == ganzhi.ZhiShen && branches[1] == ganzhi.ZhiZi) || (branches[0] == ganzhi.ZhiZi && branches[1] == ganzhi.ZhiShen) {
				foundShenZi = true
				if ph.Element != ganzhi.WxShui {
					t.Errorf("申子半合 element = %s, want 水", ph.Element)
				}
			}
		}
	}
	if !foundYouChou {
		t.Error("酉丑半合金 not detected")
	}
	if !foundShenZi {
		t.Error("申子半合水 not detected")
	}
	t.Logf("san_he_partial: %d entries", len(detection.SanHePartial))
	for _, ph := range detection.SanHePartial {
		t.Logf("  %s [%s] element=%s pillars=%v", ph.Name, ph.Branches, ph.Element, ph.Pillars)
	}
}

func TestBanHe_RequiresMiddleBranch(t *testing.T) {
	// 巳丑 missing 旺支酉；子辰 missing 旺支子。二者是拱合候选，不是半合。
	bz := ganzhi.Bazi{
		Nian: ganzhi.Zhu{Gan: ganzhi.GanJia, Zhi: ganzhi.ZhiChen},
		Yue:  ganzhi.Zhu{Gan: ganzhi.GanBing, Zhi: ganzhi.ZhiChou},
		Ri:   ganzhi.Zhu{Gan: ganzhi.GanWu, Zhi: ganzhi.ZhiSi},
		Shi:  ganzhi.Zhu{Gan: ganzhi.GanGeng, Zhi: ganzhi.ZhiWei},
	}
	detection := DetectRelations(canonicalChart(bz, ganzhi.Male, 2000).ToBazi())
	for _, partial := range detection.SanHePartial {
		if zhiSliceEqual(partial.Branches, []ganzhi.Zhi{ganzhi.ZhiSi, ganzhi.ZhiChou}) || zhiSliceEqual(partial.Branches, []ganzhi.Zhi{ganzhi.ZhiZi, ganzhi.ZhiChen}) {
			t.Fatalf("missing-middle pair %v must not be 半合; got %#v", partial.Branches, detection.SanHePartial)
		}
	}
}

func zhiSliceEqual(a, b []ganzhi.Zhi) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
