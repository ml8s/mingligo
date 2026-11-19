package tianwen

import (
	"testing"
	"time"
)

func dateFromLunarForTest(t *testing.T, lunar LunarDate) Date {
	t.Helper()
	date, err := ShanghaiLunarReference().DateFromLunar(lunar)
	if err != nil {
		t.Fatalf("DateFromLunar(%+v): %v", lunar, err)
	}
	return date
}

func lunarDateAtForTest(ref *LunarReference, at time.Time) LunarDate {
	return ref.SolarDateToLunar(ref.dateOf(at))
}
