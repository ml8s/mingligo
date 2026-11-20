package bazi

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
	"testing"
)

func TestDayBranchType_TombBranch(t *testing.T) {
	ganzhi.Load()
	for _, zhi := range []ganzhi.Zhi{ganzhi.ZhiChen, ganzhi.ZhiXu, ganzhi.ZhiChou, ganzhi.ZhiWei} {
		chart := Chart{Ri: Pillar{Zhu: ganzhi.Zhu{Gan: ganzhi.GanJia, Zhi: zhi}}}
		if got := dayBranchType(chart); got != "墓库" {
			t.Errorf("zhi %s: dayBranchType = %q, want 墓库", ganzhi.DecodeZhi(zhi), got)
		}
	}
}

func TestDayBranchType_ShenShaPriority(t *testing.T) {
	ganzhi.Load()
	chart := Chart{Ri: Pillar{
		Zhu:     ganzhi.Zhu{Gan: ganzhi.GanJia, Zhi: ganzhi.ZhiZi},
		ShenSha: []ShenShaProjection{{Name: "桃花"}, {Name: "驿马"}},
	}}
	if got := dayBranchType(chart); got != "桃花" {
		t.Errorf("with 桃花+驿马: got %q, want 桃花", got)
	}
}
