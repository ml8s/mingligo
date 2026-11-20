package qimen

import "testing"

func TestDecodeSpiritKey_StableAcrossSchools(t *testing.T) {
	want := map[SpiritIndex]string{
		SpiritZhiFu: "zhifu", SpiritTengShe: "tengshe", SpiritTaiYin: "taiyin",
		SpiritLiuHe: "liuhe", SpiritGouChen: "gouchen", SpiritZhuQue: "zhuque",
		SpiritJiuDi: "jiudi", SpiritJiuTian: "jiutian", SpiritTaiChang: "taichang",
	}
	for spirit, key := range want {
		if got := DecodeSpiritKey(spirit); got != key {
			t.Errorf("DecodeSpiritKey(%d) = %q, want %q", spirit, got, key)
		}
	}
	// Slot identity stays constant regardless of school/yin-dun display name.
	if SpiritGouChen.YangName() != "勾陈" || SpiritGouChen.YinName() != "白虎" {
		t.Error("gouchen display names changed")
	}
	if DecodeSpiritKey(SpiritGouChen) != "gouchen" {
		t.Error("gouchen key must not follow display name (baihu)")
	}
}
