package qimen

import (
	"testing"
)

// TestGongFactsAgainstClassicalTables guards the migrated gong/door/star
// facts against the traditional lookup values they were derived from.
func TestGongFactsAgainstClassicalTables(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatal(err)
	}

	// 八卦人物（《说卦传》）与内外盘定式，阳遁/阴遁。
	wantPerson := map[GongIndex]string{
		GongKan: "中男", GongKun: "老妇", GongZhen: "长男 / 壮男", GongXun: "长女",
		GongZhong: "", GongQian: "老阳", GongDui: "少女", GongGen: "少男 / 童年", GongLi: "中女",
	}
	wantDomainYang := map[GongIndex]string{
		GongKan: "inner", GongKun: "inner", GongZhen: "inner", GongXun: "inner",
		GongZhong: "center", GongQian: "outer", GongDui: "outer", GongGen: "outer", GongLi: "outer",
	}
	yang := computeGongProfiles(false)
	yin := computeGongProfiles(true)
	for g, wantPerson := range wantPerson {
		got := yang[g-1]
		if got.Person != wantPerson {
			t.Errorf("gong %s person = %q, want %q", DecodeGongIndex(g), got.Person, wantPerson)
		}
		if got.Domain != wantDomainYang[g] {
			t.Errorf("gong %s yang domain = %q, want %q", DecodeGongIndex(g), got.Domain, wantDomainYang[g])
		}
		wantYinDomain := "outer"
		if g >= GongQian {
			wantYinDomain = "inner"
		}
		if g == GongZhong {
			wantYinDomain = "center"
		}
		if yin[g-1].Domain != wantYinDomain {
			t.Errorf("gong %s yin domain = %q, want %q", DecodeGongIndex(g), yin[g-1].Domain, wantYinDomain)
		}
	}

	// 门星吉凶（经典固有属性）。
	men := computeMenJiXiong()
	jiByMen := map[DoorIndex]string{}
	for _, m := range men {
		jiByMen[m.Men] = m.JiXiong
	}
	if jiByMen[DoorSheng] != "大吉" || jiByMen[DoorSi] != "大凶" {
		t.Errorf("door ji_xiong wrong: sheng=%q si=%q", jiByMen[DoorSheng], jiByMen[DoorSi])
	}
	xing := computeXingJiXiong()
	jiByXing := map[StarIndex]string{}
	for _, x := range xing {
		jiByXing[x.Xing] = x.JiXiong
	}
	if jiByXing[StarTianPeng] != "凶" || jiByXing[StarTianFu] != "吉" {
		t.Errorf("star ji_xiong wrong: peng=%q fu=%q", jiByXing[StarTianPeng], jiByXing[StarTianFu])
	}
}
