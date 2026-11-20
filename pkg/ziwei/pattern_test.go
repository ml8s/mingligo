package ziwei

import (
	"testing"
)

func makeGong(idx GongIndex, stars []StarInfo, zaYao []string) Gong {
	return Gong{
		Index: idx,
		Name:  gongLabels[idx],
		Stars: stars,
		ZaYao: encodeShenShaList(zaYao),
	}
}

func makeStar(star StarIndex, name string, siHua string) StarInfo {
	return StarInfo{Star: star, Name: name, IsMajor: star < 14, SiHua: SiHua(siHua)}
}

func hasPattern(patterns []Pattern, name string) bool {
	for _, p := range patterns {
		if p.Name == name {
			return true
		}
	}
	return false
}

func TestJuoHuoYang(t *testing.T) {
	var palaces [12]Gong
	for i := range palaces {
		palaces[i] = makeGong(GongIndex(i), nil, nil)
	}
	palaces[0].Stars = []StarInfo{
		makeStar(JuMen, "巨门", ""),
		makeStar(HuoXing, "火星", ""),
		makeStar(QingYang, "擎羊", ""),
	}
	result := findPatterns(palaces)
	if !hasPattern(result, "巨火羊") {
		t.Errorf("expected 巨火羊 Pattern, got: %v", result)
	}
}

func TestJuoHuoYangNegative(t *testing.T) {
	var palaces [12]Gong
	for i := range palaces {
		palaces[i] = makeGong(GongIndex(i), nil, nil)
	}
	palaces[0].Stars = []StarInfo{
		makeStar(JuMen, "巨门", ""),
		makeStar(HuoXing, "火星", ""),
	}
	result := findPatterns(palaces)
	if hasPattern(result, "巨火羊") {
		t.Errorf("should NOT have 巨火羊 without 擎羊")
	}
}

func TestXingJiJiaYin(t *testing.T) {
	var palaces [12]Gong
	for i := range palaces {
		palaces[i] = makeGong(GongIndex(i), nil, nil)
	}
	palaces[0].Stars = []StarInfo{makeStar(TianXiang, "天相", "")}
	palaces[11].Stars = []StarInfo{makeStar(TianLiang, "天梁", "忌")}
	palaces[1].ZaYao = encodeShenShaList([]string{"天刑"})
	result := findPatterns(palaces)
	if !hasPattern(result, "刑忌夹印") {
		t.Errorf("expected 刑忌夹印 Pattern, got: %v", result)
	}
}

func TestXingJiJiaYinReverse(t *testing.T) {
	var palaces [12]Gong
	for i := range palaces {
		palaces[i] = makeGong(GongIndex(i), nil, nil)
	}
	palaces[0].Stars = []StarInfo{makeStar(TianXiang, "天相", "")}
	palaces[11].ZaYao = encodeShenShaList([]string{"天刑"})
	palaces[1].Stars = []StarInfo{makeStar(TianLiang, "天梁", "忌")}
	result := findPatterns(palaces)
	if !hasPattern(result, "刑忌夹印") {
		t.Errorf("expected 刑忌夹印 (reverse), got: %v", result)
	}
}

func TestXingJiJiaYinNegative(t *testing.T) {
	var palaces [12]Gong
	for i := range palaces {
		palaces[i] = makeGong(GongIndex(i), nil, nil)
	}
	palaces[0].Stars = []StarInfo{makeStar(TianXiang, "天相", "")}
	palaces[11].Stars = []StarInfo{makeStar(TianLiang, "天梁", "")}
	palaces[1].ZaYao = encodeShenShaList([]string{"天刑"})
	result := findPatterns(palaces)
	if hasPattern(result, "刑忌夹印") {
		t.Errorf("should NOT have 刑忌夹印 without 化忌")
	}
}

func TestXingJiJiaYinNoTianXing(t *testing.T) {
	var palaces [12]Gong
	for i := range palaces {
		palaces[i] = makeGong(GongIndex(i), nil, nil)
	}
	palaces[0].Stars = []StarInfo{makeStar(TianXiang, "天相", "")}
	palaces[11].Stars = []StarInfo{makeStar(TianLiang, "天梁", "忌")}
	result := findPatterns(palaces)
	if hasPattern(result, "刑忌夹印") {
		t.Errorf("should NOT have 刑忌夹印 without 天刑")
	}
}

func encodeShenShaList(names []string) []ShenSha {
	out := make([]ShenSha, 0, len(names))
	for _, n := range names { if v, ok := EncodeShenSha(n); ok { out = append(out, v) } }
	return out
}
