package ziwei

import "testing"

func TestBrightnessGroupOf(t *testing.T) {
	cases := []struct {
		brightness Brightness
		want       BrightnessGroup
	}{
		{Miao, brightnessGroupMiaoWang},
		{Wang, brightnessGroupMiaoWang},
		{De, brightnessGroupPing},
		{Li, brightnessGroupPing},
		{Ping, brightnessGroupPing},
		{Xian, brightnessGroupLuoXian},
	}
	for _, c := range cases {
		if got := brightnessGroupOf(c.brightness); got != c.want {
			t.Errorf("brightness %s: group = %q, want %q", DecodeBrightness(c.brightness), got, c.want)
		}
	}
}

func TestMajorStarStateSymbols(t *testing.T) {
	if majorStarNone != "无主星" || majorStarSingle != "唯一主星" || majorStarMulti != "多主星" {
		t.Error("major star state values changed")
	}
}
