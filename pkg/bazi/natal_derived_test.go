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

func TestSpousePalaceState_AllStatesAndPriority(t *testing.T) {
	ganzhi.Load()
	// Each relation kind involving the day pillar maps to its state.
	cases := []struct {
		kind NatalRelationKind
		want PalaceState
	}{
		{relationKindLiuChong, palaceStateChong},
		{relationKindZhiLiuHe, palaceStateHe},
		{relationKindLiuXing, palaceStateXing},
		{relationKindLiuHai, palaceStateHai},
	}
	for _, c := range cases {
		chart := Chart{GuanXi: []NatalRelationFact{{
			Kind:    c.kind,
			Members: []NatalRelationMember{{Pillar: pillarRi}},
		}}}
		if got := spousePalaceState(chart); got != c.want {
			t.Errorf("kind %s: state = %q, want %q", DecodeNatalRelationKind(c.kind), got, c.want)
		}
	}
	// No relation involving day pillar → 静.
	chart := Chart{GuanXi: []NatalRelationFact{{
		Kind:    relationKindLiuChong,
		Members: []NatalRelationMember{{Pillar: pillarNian}},
	}}}
	if got := spousePalaceState(chart); got != palaceStateQuiet {
		t.Errorf("no day relation: state = %q, want 静", got)
	}
	// Priority: 冲 wins when multiple relations involve the day pillar.
	chart = Chart{GuanXi: []NatalRelationFact{
		{Kind: relationKindZhiLiuHe, Members: []NatalRelationMember{{Pillar: pillarRi}}},
		{Kind: relationKindLiuChong, Members: []NatalRelationMember{{Pillar: pillarRi}}},
	}}
	if got := spousePalaceState(chart); got != palaceStateChong {
		t.Errorf("冲+合 priority: state = %q, want 冲", got)
	}
}
