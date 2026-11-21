package bazhai

// This file defines the domain-owned, domain-aggregated output views of the
// Bazhai tools. The views mirror the domain fields into output contracts with
// jsonschema tags; tags are inert metadata read only when a consumer reflects
// the type for its schema.

// GuaXiangView is a bagua trigram.
type GuaXiangView struct {
	Index   int    `json:"index" jsonschema:"洛书宫数。"`
	Name    string `json:"name" jsonschema:"卦名。"`
	Wuxing  string `json:"wuxing" jsonschema:"卦五行。"`
	YinYang string `json:"yin_yang" jsonschema:"卦阴阳。"`
}

// MingGuaXiangView is the命卦 result.
type MingGuaXiangView struct {
	Gua          GuaXiangView `json:"gua" jsonschema:"命卦。"`
	Group        LifeGroup    `json:"dong_xi_si_ming" jsonschema:"东四命或西四命。"`
	YearBoundary string       `json:"year_boundary" jsonschema:"年界：立春。"`
}

// SiJiSiXiongFangView groups四吉四凶方向.
type SiJiSiXiongFangView struct {
	ShengQi []string `json:"sheng_qi" jsonschema:"生气方。"`
	TianYi  []string `json:"tian_yi" jsonschema:"天医方。"`
	YanNian []string `json:"yan_nian" jsonschema:"延年方。"`
	FuWei   []string `json:"fu_wei" jsonschema:"伏位方。"`
	HuoHai  []string `json:"huo_hai" jsonschema:"祸害方。"`
	WuGui   []string `json:"wu_gui" jsonschema:"五鬼方。"`
	LiuSha  []string `json:"liu_sha" jsonschema:"六煞方。"`
	JueMing []string `json:"jue_ming" jsonschema:"绝命方。"`
}

// NianXingView is one palace's purple-white flying star.
type NianXingView struct {
	GongNum  int    `json:"gong_num" jsonschema:"洛书宫数。"`
	Xing     int    `json:"xing" jsonschema:"紫白星数。"`
	XingName string `json:"xing_name" jsonschema:"紫白星名。"`
	Wuxing   string `json:"wuxing" jsonschema:"星五行。"`
	RuZhong  bool   `json:"ru_zhong" jsonschema:"是否入中。"`
}

// ChuShengNianBaiXingView groups出生年紫白飞星.
type ChuShengNianBaiXingView struct {
	Year    int            `json:"year" jsonschema:"立春年份。"`
	RuZhong string         `json:"ru_zhong" jsonschema:"入中紫白星。"`
	Palaces []NianXingView `json:"gong_wei" jsonschema:"九宫紫白星。"`
}

// BaZhaiMingGuaPanView is the domain-aggregated output of a Bazhai chart.
type BaZhaiMingGuaPanView struct {
	MingGua         MingGuaXiangView        `json:"ming_gua" jsonschema:"命卦。"`
	SiJiSiXiongFang SiJiSiXiongFangView     `json:"si_ji_si_xiong_fang" jsonschema:"四吉四凶方向。"`
	YearStars       ChuShengNianBaiXingView `json:"birth_year_xing" jsonschema:"出生年紫白飞星。"`
}

// MenZhuZaoView groups the门主灶 assessment (bazhai_layout tool output).
type MenZhuZaoView struct {
	DongXiSiMing LifeGroup     `json:"dong_xi_si_ming" jsonschema:"东四命或西四命。"`
	MingGua      string        `json:"ming_gua" jsonschema:"命卦。"`
	Door         DoorStoveInfo `json:"men" jsonschema:"门位判断。"`
	Master       DoorStoveInfo `json:"zhu" jsonschema:"主卧判断。"`
	Stove        DoorStoveInfo `json:"zao" jsonschema:"灶位判断。"`
}

// View produces the domain-aggregated output from the computed chart.
func (c Chart) View() BaZhaiMingGuaPanView {
	stars := ChuShengNianBaiXingView{
		Year: c.YearStars.Year, RuZhong: c.YearStars.RuZhong,
	}
	stars.Palaces = make([]NianXingView, 0, len(c.YearStars.Palaces))
	for _, p := range c.YearStars.Palaces {
		stars.Palaces = append(stars.Palaces, NianXingView{
			GongNum: p.GongNum, Xing: p.Xing, XingName: p.XingName,
			Wuxing: p.Wuxing, RuZhong: p.RuZhong,
		})
	}
	return BaZhaiMingGuaPanView{
		MingGua: MingGuaXiangView{
			Gua: GuaXiangView{
				Index: c.MingGua.Gua.Index, Name: c.MingGua.Gua.Name,
				Wuxing: c.MingGua.Gua.Wuxing, YinYang: c.MingGua.Gua.YinYang,
			},
			Group: c.MingGua.Group, YearBoundary: c.MingGua.YearBoundary,
		},
		SiJiSiXiongFang: SiJiSiXiongFangView{
			ShengQi: c.BaZhaiDirs.ShengQi, TianYi: c.BaZhaiDirs.TianYi,
			YanNian: c.BaZhaiDirs.YanNian, FuWei: c.BaZhaiDirs.FuWei,
			HuoHai: c.BaZhaiDirs.HuoHai, WuGui: c.BaZhaiDirs.WuGui,
			LiuSha: c.BaZhaiDirs.LiuSha, JueMing: c.BaZhaiDirs.JueMing,
		},
		YearStars: stars,
	}
}

// View produces the domain-aggregated layout assessment output.
func (r LayoutResult) View() MenZhuZaoView {
	return MenZhuZaoView{
		DongXiSiMing: r.Group,
		MingGua:      r.MingGua,
		Door:         r.Door,
		Master:       r.Master,
		Stove:        r.Stove,
	}
}
