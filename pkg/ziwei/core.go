package ziwei

import "github.com/ml8s/mingligo/pkg/ganzhi"

// Core is the narrow natal fact set consumed by derived Ziwei methods.
// It intentionally excludes palace projections; derived methods must not take
// the complete chart.
type Core struct {
	Gender    ganzhi.Gender
	NianGan   Gan
	JuShu     JuShu
	BirthYear int
}

// Core projects the chart to the narrow facts consumed by derived methods.
func (c Chart) Core() Core {
	return Core{
		Gender:    c.Gender,
		NianGan:   c.NianGan,
		JuShu:     c.JuShu,
		BirthYear: c.BirthYear,
	}
}

// FlowContext carries only the natal facts required by flow month/day/hour
// and annual projection. It prevents derived Ziwei methods from consuming the
// complete Chart.
type FlowContext struct {
	Core            Core
	BirthLunarMonth int
	BirthLunarDay   int
	BirthIsLeap     bool
	ShiZhi          Zhi
	MingGongZhi     Zhi
	SiHuaPalaces    map[StarIndex]GongIndex
}

// FlowContext extracts the narrow facts used by flow methods.
func (c Chart) FlowContext() FlowContext {
	siHuaPalaces := make(map[StarIndex]GongIndex)
	for _, palace := range c.GongWei {
		for _, star := range palace.Stars {
			if _, exists := siHuaPalaces[star.Star]; !exists {
				siHuaPalaces[star.Star] = palace.Index
			}
		}
	}
	return FlowContext{
		Core:            c.Core(),
		BirthLunarMonth: c.BirthLunarMonth,
		BirthLunarDay:   c.LunarDay,
		BirthIsLeap:     c.BirthIsLeap,
		ShiZhi:          c.ShiZhi,
		MingGongZhi:     c.GongWei[c.MingGong].Zhi,
		SiHuaPalaces:    siHuaPalaces,
	}
}

// BondPalaceFacts carries only the palace facts consumed by bond output.
type BondPalaceFacts struct {
	Palace     string      `json:"gong"`
	MajorStars []StarIndex `json:"major_xing"`
}

// BondSideFacts carries only the natal palace facts consumed by bond.
type BondSideFacts struct {
	MingGong BondPalaceFacts `json:"-"`
	FuQiGong BondPalaceFacts `json:"-"`
	ZiNvGong BondPalaceFacts `json:"-"`
}

// BondSideFacts extracts the palace facts consumed by bond.
func (c Chart) BondSideFacts() BondSideFacts {
	facts := func(idx int) BondPalaceFacts {
		majorStars := make([]StarIndex, 0, len(c.GongWei[idx].Stars))
		for _, star := range c.GongWei[idx].Stars {
			if star.IsMajor {
				majorStars = append(majorStars, star.Star)
			}
		}
		return BondPalaceFacts{Palace: c.GongWei[idx].Name, MajorStars: majorStars}
	}
	return BondSideFacts{
		MingGong: facts(int(c.MingGong)),
		FuQiGong: facts(2),
		ZiNvGong: facts(3),
	}
}
