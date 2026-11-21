package bazhai

// This file defines the domain-owned, domain-aggregated output views of the
// Bazhai tools. The views mirror the domain fields into output contracts with
// jsonschema tags; tags are inert metadata read only when a consumer reflects
// the type for its schema.

// GuaView is a bagua trigram.
type GuaView struct {
	Index   int    `json:"index" jsonschema:"洛书宫数。"`
	Name    string `json:"name" jsonschema:"卦名。"`
	Wuxing  string `json:"wuxing" jsonschema:"卦五行。"`
	YinYang string `json:"yin_yang" jsonschema:"卦阴阳。"`
}

// MingGuaView is the命卦 result.
type MingGuaView struct {
	Gua          GuaView   `json:"gua" jsonschema:"命卦。"`
	Group        LifeGroup `json:"group" jsonschema:"东四命或西四命。"`
	YearBoundary string    `json:"year_boundary" jsonschema:"年界：立春。"`
}

// DirectionsView groups四吉四凶方向.
type DirectionsView struct {
	ShengQi []string `json:"sheng_qi" jsonschema:"生气方。"`
	TianYi  []string `json:"tian_yi" jsonschema:"天医方。"`
	YanNian []string `json:"yan_nian" jsonschema:"延年方。"`
	FuWei   []string `json:"fu_wei" jsonschema:"伏位方。"`
	HuoHai  []string `json:"huo_hai" jsonschema:"祸害方。"`
	WuGui   []string `json:"wu_gui" jsonschema:"五鬼方。"`
	LiuSha  []string `json:"liu_sha" jsonschema:"六煞方。"`
	JueMing []string `json:"jue_ming" jsonschema:"绝命方。"`
}

// AnnualStarView is one palace's purple-white flying star.
type AnnualStarView struct {
	GongNum  int    `json:"gong_num" jsonschema:"洛书宫数。"`
	Xing     int    `json:"xing" jsonschema:"紫白星数。"`
	XingName string `json:"xing_name" jsonschema:"紫白星名。"`
	Wuxing   string `json:"wuxing" jsonschema:"星五行。"`
	RuZhong  bool   `json:"ru_zhong" jsonschema:"是否入中。"`
}

// YearStarsView groups出生年紫白飞星.
type YearStarsView struct {
	Year         int              `json:"year" jsonschema:"立春年份。"`
	RuZhong      string           `json:"ru_zhong" jsonschema:"入中紫白星。"`
	YearBoundary string           `json:"year_boundary" jsonschema:"年界：立春。"`
	Palaces      []AnnualStarView `json:"gong_wei" jsonschema:"九宫紫白星。"`
}

// ChartView is the domain-aggregated output of a Bazhai chart.
type ChartView struct {
	MingGua    MingGuaView    `json:"ming_gua" jsonschema:"命卦。"`
	BaZhaiDirs DirectionsView `json:"ba_zhai_dirs" jsonschema:"四吉四凶方向。"`
	YearStars  YearStarsView  `json:"birth_year_xing" jsonschema:"出生年紫白飞星。"`
}

// LayoutView groups the门主灶 assessment (bazhai_layout tool output).
type LayoutView struct {
	Group   LifeGroup     `json:"group" jsonschema:"东四命或西四命。"`
	MingGua string        `json:"ming_gua" jsonschema:"命卦。"`
	Door    DoorStoveInfo `json:"men" jsonschema:"门位判断。"`
	Master  DoorStoveInfo `json:"zhu" jsonschema:"主卧判断。"`
	Stove   DoorStoveInfo `json:"zao" jsonschema:"灶位判断。"`
}

// View produces the domain-aggregated output from the computed chart.
func (c Chart) View() ChartView {
	stars := YearStarsView{
		Year: c.YearStars.Year, RuZhong: c.YearStars.RuZhong,
		YearBoundary: c.YearStars.YearBoundary,
	}
	stars.Palaces = make([]AnnualStarView, 0, len(c.YearStars.Palaces))
	for _, p := range c.YearStars.Palaces {
		stars.Palaces = append(stars.Palaces, AnnualStarView{
			GongNum: p.GongNum, Xing: p.Xing, XingName: p.XingName,
			Wuxing: p.Wuxing, RuZhong: p.RuZhong,
		})
	}
	return ChartView{
		MingGua: MingGuaView{
			Gua: GuaView{
				Index: c.MingGua.Gua.Index, Name: c.MingGua.Gua.Name,
				Wuxing: c.MingGua.Gua.Wuxing, YinYang: c.MingGua.Gua.YinYang,
			},
			Group: c.MingGua.Group, YearBoundary: c.MingGua.YearBoundary,
		},
		BaZhaiDirs: DirectionsView{
			ShengQi: c.BaZhaiDirs.ShengQi, TianYi: c.BaZhaiDirs.TianYi,
			YanNian: c.BaZhaiDirs.YanNian, FuWei: c.BaZhaiDirs.FuWei,
			HuoHai: c.BaZhaiDirs.HuoHai, WuGui: c.BaZhaiDirs.WuGui,
			LiuSha: c.BaZhaiDirs.LiuSha, JueMing: c.BaZhaiDirs.JueMing,
		},
		YearStars: stars,
	}
}

// View produces the domain-aggregated layout assessment output.
func (r LayoutResult) View() LayoutView {
	return LayoutView(r)
}
