package qimen

import "github.com/ml8s/mingligo/pkg/luoshu"

// gongFact is the decoded lookup value for one palace's fixed traditional
// attributes. Keyed by GongIndex (symbol) in gongFactsTable.
type gongFact struct {
	Name       string
	Person     string
	Range      string
	DomainYang string
	DomainYin  string
}

// jiXiongFact is the decoded lookup value for a door/star's inherent
// auspiciousness and classical application. Keyed by DoorIndex/StarIndex.
type jiXiongFact struct {
	JiXiong string
	YingShi string
}

// GongProfile is the boundary view of one palace's fixed factual attributes:
// direction (from luoshu), person/range (from qimen gong facts), and the
// yin/yang-dun inner-outer domain. Values are decoded text — lookup values
// that never enter symbolic computation.
type GongProfile struct {
	Gong      GongIndex `json:"gong"`
	Name      string    `json:"name"`
	Direction string    `json:"direction"`
	Person    string    `json:"person,omitempty"`
	Range     string    `json:"range"`
	Domain    string    `json:"domain"`
}

// MenJiXiong is a door's inherent auspiciousness and classical application.
type MenJiXiong struct {
	Men     DoorIndex `json:"men"`
	JiXiong string    `json:"ji_xiong"`
	YingShi string    `json:"yingshi"`
}

// XingJiXiong is a star's inherent auspiciousness and classical application.
type XingJiXiong struct {
	Xing    StarIndex `json:"xing"`
	JiXiong string    `json:"ji_xiong"`
	YingShi string    `json:"yingshi"`
}

// computeGongProfiles decodes the fixed per-palace facts for the whole board.
// Direction comes from luoshu (shared authority); person/range/domain from
// the qimen gong facts table; inner-outer depends on yin dun.
func computeGongProfiles(yinDun bool) [9]GongProfile {
	var out [9]GongProfile
	for i := 1; i <= 9; i++ {
		g := GongIndex(i)
		f := gongFactsTable[g]
		domain := f.DomainYang
		if yinDun {
			domain = f.DomainYin
		}
		out[i-1] = GongProfile{
			Gong:      g,
			Name:      DecodeGongIndex(g),
			Direction: luoshu.ByNumber(i).Direction,
			Person:    f.Person,
			Range:     f.Range,
			Domain:    domain,
		}
	}
	return out
}

// computeMenJiXiong decodes all eight doors' inherent attributes.
func computeMenJiXiong() []MenJiXiong {
	out := make([]MenJiXiong, 0, 8)
	for d := DoorIndex(1); d <= 8; d++ {
		f := doorJiXiongTable[d]
		out = append(out, MenJiXiong{Men: d, JiXiong: f.JiXiong, YingShi: f.YingShi})
	}
	return out
}

// computeXingJiXiong decodes all nine stars' inherent attributes.
func computeXingJiXiong() []XingJiXiong {
	out := make([]XingJiXiong, 0, 9)
	for s := StarIndex(1); s <= 9; s++ {
		f := starJiXiongTable[s]
		out = append(out, XingJiXiong{Xing: s, JiXiong: f.JiXiong, YingShi: f.YingShi})
	}
	return out
}
