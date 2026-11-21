package fengshui

import (
	"testing"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/luoshu"
)

func Test24MountainsYinYang(t *testing.T) {
	// Build index of 天元龙 yin-yang by trigram (symbol keys).
	tianYinYang := make(map[luoshu.Trigram]ganzhi.YinYang)
	for _, m := range Mountains24Table {
		if m.YuanLong == YuanLongTian {
			tianYinYang[m.Trigram] = m.YinYang
		}
	}

	// Verify 天元龙 follows 四正阴/四隅阳.
	for _, m := range Mountains24Table {
		if m.YuanLong != YuanLongTian {
			continue
		}
		var expected ganzhi.YinYang
		switch m.Trigram {
		case luoshu.TrigramKan, luoshu.TrigramZhen, luoshu.TrigramLi, luoshu.TrigramDui:
			expected = ganzhi.Yin // 四正卦天元=阴
		default:
			expected = ganzhi.Yang // 四隅卦(乾坤艮巽)天元=阳
		}
		if m.YinYang != expected {
			t.Errorf("%s(%s 天元): yinyang=%v, want %v", m.Name, m.Trigram, m.YinYang, expected)
		}
	}

	// Verify 人元龙 = 天元龙 (same yin-yang).
	for _, m := range Mountains24Table {
		if m.YuanLong != YuanLongRen {
			continue
		}
		expected := tianYinYang[m.Trigram]
		if m.YinYang != expected {
			t.Errorf("%s(%s 人元): yinyang=%v, want %v (应同天元)", m.Name, m.Trigram, m.YinYang, expected)
		}
	}

	// Verify 地元龙 = 天元龙 reversed.
	for _, m := range Mountains24Table {
		if m.YuanLong != YuanLongDi {
			continue
		}
		tianYY := tianYinYang[m.Trigram]
		expected := !tianYY
		if m.YinYang != expected {
			t.Errorf("%s(%s 地元): yinyang=%v, want %v (应反天元)", m.Name, m.Trigram, m.YinYang, expected)
		}
	}
}

// TestStarTable verifies紫白飞星 data integrity.
func TestStarTable(t *testing.T) {
	// Verify 1-9 all present and non-empty.
	for i := 1; i <= 9; i++ {
		s := StarTable[i]
		if s.Number != i {
			t.Errorf("star %d: number = %d", i, s.Number)
		}
		if s.Name == "" {
			t.Errorf("star %d: empty name", i)
		}
		if s.Color == "" {
			t.Errorf("star %d: empty color", i)
		}
	}

}

// TestPalaceTable verifies九宫 data integrity.
func TestPalaceTable(t *testing.T) {
	palaceNames := map[int]string{
		1: "坎", 2: "坤", 3: "震", 4: "巽", 5: "中",
		6: "乾", 7: "兑", 8: "艮", 9: "离",
	}
	for i := 1; i <= 9; i++ {
		p := luoshu.PalaceTable[i]
		if p.Number != i {
			t.Errorf("palace %d: number = %d", i, p.Number)
		}
		if p.Name != palaceNames[i] {
			t.Errorf("palace %d: name = %s, want %s", i, p.Name, palaceNames[i])
		}
	}
}

// TestLuoshuFlyOrder verifies the飞星 order skips center.
func TestLuoshuFlyOrder(t *testing.T) {
	// The 8-cell fly order: 6,7,8,9,1,2,3,4
	expected := [8]int{6, 7, 8, 9, 1, 2, 3, 4}
	if LuoshuFlyOrder != expected {
		t.Errorf("LuoshuFlyOrder = %v, want %v", LuoshuFlyOrder, expected)
	}
}
