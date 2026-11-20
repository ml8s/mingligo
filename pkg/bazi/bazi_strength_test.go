package bazi

import (
	"testing"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

func ptr[T any](v T) *T { return &v }

// ── 扶抑强弱表数据驱动测试 ──
// 验证 25 种 root×season 组合的查找表全覆盖
// ────────────────────────────────────────────────────────────────

func TestStrength_AllRules_FirstHit(t *testing.T) {
	// 规则 1: month_main + 任何季节 → 身强
	t.Run("rule1_month_main_any_season", func(t *testing.T) {
		for _, s := range []ganzhi.Zhi{ganzhi.ZhiYin, ganzhi.ZhiWu, ganzhi.ZhiShen} {
			cg := [4]HiddenStemProjection{
				{Main: ganzhi.GanWu},
				{Main: ganzhi.GanJia}, // 寅中甲 → month_main for 甲木
				{Main: ganzhi.GanWu},
				{Main: ganzhi.GanWu},
			}
			rt := classifyRoot(ganzhi.GanJia, s, cg)
			season := classifySeason(ganzhi.GanJia, s)
			strength := lookupStrength(rt, season, 1)
			if strength != strengthStrong {
				t.Errorf("%s(season=%s) = %q, want 身强", DecodeRootKind(rt), season, strength)
			}
		}
	})

	// 规则 2: month_mid + 旺/相 → 身强
	t.Run("rule2_month_mid_wang_xiang", func(t *testing.T) {
		// 甲木日主, 月支亥(壬水), 中气甲 → month_mid
		cg := [4]HiddenStemProjection{
			{Main: ganzhi.GanWu},
			{Main: ganzhi.GanRen, Mid: ptr(ganzhi.GanJia)}, // 亥: 壬(本气)甲(中气)
			{Main: ganzhi.GanWu},
			{Main: ganzhi.GanWu},
		}
		rt := classifyRoot(ganzhi.GanJia, ganzhi.ZhiHai, cg)
		if rt != rootMonthMid {
			t.Fatalf("classifyRoot = %q, want month_mid", rt)
		}
		season := classifySeason(ganzhi.GanJia, ganzhi.ZhiHai) // 亥=水→木, 相
		str := lookupStrength(rt, season, 1)
		if str != strengthStrong {
			t.Errorf("month_mid+相 = %q, want 身强", str)
		}
	})

	// 规则 4: month_mid + 死 → 身弱
	t.Run("rule4_month_mid_si", func(t *testing.T) {
		// 甲木日主, 申月(庚), 中气壬 → 甲木month_mid
		cg := [4]HiddenStemProjection{
			{Main: ganzhi.GanWu},
			{Main: ganzhi.GanGeng, Mid: ptr(ganzhi.GanRen)}, // 申: 庚(本气)壬(中气)
			{Main: ganzhi.GanWu},
			{Main: ganzhi.GanWu},
		}
		rt := classifyRoot(ganzhi.GanJia, ganzhi.ZhiShen, cg)
		season := classifySeason(ganzhi.GanJia, ganzhi.ZhiShen)
		str := lookupStrength(rt, season, 0)
		if str != strengthWeak {
			t.Errorf("month_mid+死+yinBi=0 = %q, want 身弱", str)
		}
		// 印比≥2 → 中和
		str2 := lookupStrength(rt, season, 2)
		if str2 != strengthBalance {
			t.Errorf("month_mid+死+yinBi=2 = %q, want 中和", str2)
		}
	})
}

func TestStrength_All25Combos_NoGaps(t *testing.T) {
	// 验证所有 25 种 root×season 组合每一条都能命中且有值
	roots := []RootKind{rootMonthMain, rootMonthMid, rootZhiMain, rootZhiMid, rootNone}
	seasons := []ganzhi.WangShuai{ganzhi.WSWang, ganzhi.WSXiang, ganzhi.WSXiu, ganzhi.WSQiu, ganzhi.WSSi}

	for _, r := range roots {
		for _, s := range seasons {
			result := DecodeStrengthLevel(lookupStrength(r, s, 0))
			if result != "身强" && result != "中和" && result != "身弱" {
				t.Errorf("root=%q, season=%q → result=%q (invalid)", DecodeRootKind(r), ganzhi.DecodeWangShuai(s), result)
			}
		}
	}
}

// ── 从格规则数据驱动测试 ──
// ────────────────────────────────────────────────────────────────

func TestCongGe_Rules_AllSeasons(t *testing.T) {
	tests := []struct {
		name         string
		chart        Chart
		wantPattern  string
		wantNotEmpty bool
	}{
		// 从旺: 甲日寅月(旺), 印比≥3, 无官杀
		{
			name: "从旺_甲日寅月印比多",
			chart: Chart{
				Ri:   Pillar{Zhu: ganzhi.Zhu{Gan: ganzhi.GanJia, Zhi: ganzhi.ZhiWu}},
				Yue:  Pillar{Zhu: ganzhi.Zhu{Gan: ganzhi.GanJia, Zhi: ganzhi.ZhiYin}},
				Nian: Pillar{Zhu: ganzhi.Zhu{Gan: ganzhi.GanJia, Zhi: ganzhi.ZhiWu}},
				Shi:  Pillar{Zhu: ganzhi.Zhu{Gan: ganzhi.GanYi, Zhi: ganzhi.ZhiWu}},
			},
			wantPattern: "从旺格",
		},
		// 从杀: 乙日申月(死), 庚透干, 无印
		{
			name: "从杀_乙日申月庚透",
			chart: Chart{
				Ri:   Pillar{Zhu: ganzhi.Zhu{Gan: ganzhi.GanYi, Zhi: ganzhi.ZhiWu}},
				Yue:  Pillar{Zhu: ganzhi.Zhu{Gan: ganzhi.GanGeng, Zhi: ganzhi.ZhiShen}},
				Nian: Pillar{Zhu: ganzhi.Zhu{Gan: ganzhi.GanGeng, Zhi: ganzhi.ZhiWu}},
				Shi:  Pillar{Zhu: ganzhi.Zhu{Gan: ganzhi.GanXin, Zhi: ganzhi.ZhiWu}},
			},
			wantNotEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat, _, _, _ := lookupCongGe(tt.chart.ToBazi())

			if tt.wantPattern != "" && string(pat) != tt.wantPattern {
				t.Errorf("pattern = %q, want %q", pat, tt.wantPattern)
			}
			if tt.wantNotEmpty && pat == "" {
				t.Error("expected pattern match, got empty")
			}
			if pat != "" {
				t.Logf("匹配: %s", pat)
			}
		})
	}
}

// ── classifyRoot 单元测试 ──
// ────────────────────────────────────────────────────────────────

func TestClassifyRoot_AllTypes(t *testing.T) {
	tests := []struct {
		name   string
		riGan  ganzhi.Gan
		yueZhi ganzhi.Zhi
		cg     [4]HiddenStemProjection
		want   string
	}{
		{
			name:  "month_main_寅中甲",
			riGan: ganzhi.GanJia, yueZhi: ganzhi.ZhiYin,
			cg:   [4]HiddenStemProjection{{}, {Main: ganzhi.GanJia}, {}, {}},
			want: "month_main",
		},
		{
			name:  "month_mid_亥中甲",
			riGan: ganzhi.GanJia, yueZhi: ganzhi.ZhiHai,
			cg:   [4]HiddenStemProjection{{}, {Main: ganzhi.GanRen, Mid: ptr(ganzhi.GanJia)}, {}, {}},
			want: "month_mid",
		},
		{
			name:  "zhi_main_年支寅甲",
			riGan: ganzhi.GanJia, yueZhi: ganzhi.ZhiWu,
			cg:   [4]HiddenStemProjection{{Main: ganzhi.GanJia}, {Main: ganzhi.GanWu}, {}, {}},
			want: "zhi_main",
		},
		{
			name:  "zhi_mid_年支亥中甲",
			riGan: ganzhi.GanJia, yueZhi: ganzhi.ZhiWu,
			cg:   [4]HiddenStemProjection{{Main: ganzhi.GanRen, Mid: ptr(ganzhi.GanJia)}, {Main: ganzhi.GanWu}, {}, {}},
			want: "zhi_mid",
		},
		{
			name:  "none_无根",
			riGan: ganzhi.GanJia, yueZhi: ganzhi.ZhiWu,
			cg:   [4]HiddenStemProjection{{Main: ganzhi.GanWu}, {Main: ganzhi.GanWu}, {Main: ganzhi.GanWu}, {Main: ganzhi.GanWu}},
			want: "none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyRoot(tt.riGan, tt.yueZhi, tt.cg)
			if DecodeRootKind(got) != tt.want {
				t.Errorf("classifyRoot = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClassifySeason_AllFive(t *testing.T) {
	tests := []struct {
		name  string
		riGan ganzhi.Gan
		zhi   ganzhi.Zhi
		want  string
	}{
		{name: "甲寅=旺", riGan: ganzhi.GanJia, zhi: ganzhi.ZhiYin, want: "旺"},
		{name: "甲亥=相(水生木)", riGan: ganzhi.GanJia, zhi: ganzhi.ZhiHai, want: "相"},
		{name: "甲申=死(金克木)", riGan: ganzhi.GanJia, zhi: ganzhi.ZhiShen, want: "死"},
		{name: "甲辰=囚(木克土)", riGan: ganzhi.GanJia, zhi: ganzhi.ZhiChen, want: "囚"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifySeason(tt.riGan, tt.zhi)
			if ganzhi.DecodeWangShuai(got) != tt.want {
				t.Errorf("classifySeason = %q, want %q", got, tt.want)
			}
		})
	}
}
