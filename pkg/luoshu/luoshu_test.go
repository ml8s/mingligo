package luoshu

import "testing"

func TestPalaceTableAnchors(t *testing.T) {
	// 洛书九宫核心锚点：宫号/卦名/方位/五行。
	anchors := []struct {
		num       int
		name      string
		direction string
		wuxing    string
	}{
		{1, "坎", "北", "水"},
		{2, "坤", "西南", "土"},
		{3, "震", "东", "木"},
		{4, "巽", "东南", "木"},
		{5, "中", "中", "土"},
		{6, "乾", "西北", "金"},
		{7, "兑", "西", "金"},
		{8, "艮", "东北", "土"},
		{9, "离", "南", "火"},
	}
	for _, a := range anchors {
		p := PalaceTable[a.num]
		if p.Number != a.num || p.Name != a.name || p.Direction != a.direction || p.Element.String() != a.wuxing {
			t.Errorf("palace %d = %+v, want %s/%s/%s", a.num, p, a.name, a.direction, a.wuxing)
		}
	}
}

func TestTrigramRoundTrip(t *testing.T) {
	for name, num := range map[string]int{
		"坎": 1, "坤": 2, "震": 3, "巽": 4, "中": 5,
		"乾": 6, "兑": 7, "艮": 8, "离": 9,
	} {
		trigram, ok := EncodeTrigram(name)
		if !ok && name != "中" { // 中宫不是八卦
			t.Errorf("EncodeTrigram(%q) failed", name)
			continue
		}
		if ok && TrigramNumber(trigram) != num {
			t.Errorf("TrigramNumber(%q) = %d, want %d", name, TrigramNumber(trigram), num)
		}
	}
}
