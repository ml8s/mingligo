package qimen

import (
	"testing"
	"time"
)

func TestJinhanSpiritSymbols(t *testing.T) {
	names := []string{"地钥", "六贼", "五符", "天曹", "地符", "风伯", "雷公", "雨师", "风云", "唐符", "国印", "天关"}
	keys := []string{"diyao", "liuze", "wufu", "tiancao", "difu", "fengbo", "leigong", "yushi", "fengyun", "tangfu", "guoyin", "tianguan"}
	if err := Load(); err != nil {
		t.Fatal(err)
	}
	for i, name := range names {
		spirit := JinhanSpirit(i + 1)
		if got := DecodeJinhanSpirit(spirit); got != name {
			t.Errorf("DecodeJinhanSpirit(%d) = %q, want %q", spirit, got, name)
		}
		if got := DecodeJinhanSpiritKey(spirit); got != keys[i] {
			t.Errorf("DecodeJinhanSpiritKey(%d) = %q, want %q", spirit, got, keys[i])
		}
		parsed, err := EncodeJinhanSpirit(name)
		if err != nil || parsed != spirit {
			t.Errorf("EncodeJinhanSpirit(%q) = %d, %v; want %d", name, parsed, err, spirit)
		}
	}
}

func TestJinhanDaySpiritsCarrySymbol(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatal(err)
	}
	chart := ComputeJinhanChartAt(time.Date(2024, 6, 29, 12, 0, 0, 0, time.UTC), 116.4)
	for i, day := range chart.Board.DaySpirits {
		if day.ShenSymbol == "" {
			t.Errorf("day spirit %d lacks shen_symbol", i)
		}
		parsed, err := EncodeJinhanSpirit(day.Shen)
		if err != nil || DecodeJinhanSpiritKey(parsed) != day.ShenSymbol {
			t.Errorf("day spirit %d: shen=%q symbol=%q mismatch (%v)", i, day.Shen, day.ShenSymbol, err)
		}
	}
}
