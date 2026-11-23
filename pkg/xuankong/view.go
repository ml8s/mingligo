package xuankong

import "github.com/ml8s/mingligo/pkg/fengshui"

// This file defines the domain-owned, domain-aggregated output view of an
// XuanKong chart. The domain Chart is already structured by术数子域; the View
// merges the mountain/facing index+name pairs into single named objects and is
// the single output model every transport adapter transmits unchanged.

// MountainView is a 24-mountain坐向 with its index and display name.
type MountainView struct {
	Index int    `json:"shan_shu" jsonschema:"二十四山序号(0-23)。"`
	Name  string `json:"shan_ming" jsonschema:"山名。"`
}

// JiaHuiView is one 星曜加会 fact (name + 断语).
type JiaHuiView struct {
	Ming string `json:"ming" jsonschema:"加会名称。"`
	YiJu string `json:"yi_ju" jsonschema:"加会断语。"`
}

// GongWeiView is one palace's 运星/山星/向星 and 加会 facts.
type GongWeiView struct {
	GongShu   int                 `json:"gong_shu" jsonschema:"洛书宫数。"`
	YunXing   fengshui.FlyingStar `json:"yun_xing" jsonschema:"运星。"`
	ShanXing  fengshui.FlyingStar `json:"shan_xing" jsonschema:"山星。"`
	XiangXing fengshui.FlyingStar `json:"xiang_xing" jsonschema:"向星。"`
	JiaHui    JiaHuiView          `json:"jia_hui" jsonschema:"星曜加会。"`
}

// ShouShanChuShaView groups 收山出煞 judgment facts.
type ShouShanChuShaView struct {
	ShouShanOK bool   `json:"shou_shan" jsonschema:"收山是否得宜。"`
	ChuShaOK   bool   `json:"chu_sha" jsonschema:"出煞是否得宜。"`
	PingYu     string `json:"ping_yu" jsonschema:"判断结论。"`
}

// ChartView is the domain-aggregated output of an XuanKong chart.
type ChartView struct {
	Yun            SanYuanYun         `json:"yun" jsonschema:"三元九运。"`
	SitMountain    MountainView       `json:"zuo_shan" jsonschema:"坐山。"`
	FaceMountain   MountainView       `json:"xiang_shan" jsonschema:"向首。"`
	GongWei        [9]GongWeiView     `json:"gong_wei" jsonschema:"九宫运星山星向星。"`
	SiDaJu         FourSituation      `json:"si_da_ju" jsonschema:"四大局。"`
	FuYinLayers    []RecitationLayer  `json:"fu_yin_ceng" jsonschema:"伏吟层。"`
	FanYinLayers   []RecitationLayer  `json:"fan_yin_ceng" jsonschema:"反吟层。"`
	ShouShanChuSha ShouShanChuShaView `json:"shou_shan_chu_sha" jsonschema:"收山出煞判断。"`
}

// View produces the domain-aggregated output from the computed chart.
func (c Chart) View() ChartView {
	var gong [9]GongWeiView
	for i, p := range c.Palaces {
		gong[i] = GongWeiView{
			GongShu: p.PalaceNum, YunXing: p.PeriodStar,
			ShanXing: p.MountainStar, XiangXing: p.FacingStar,
			JiaHui: JiaHuiView{Ming: p.JiaHui.Name, YiJu: p.JiaHui.Meaning},
		}
	}
	return ChartView{
		Yun:          c.Yun,
		SitMountain:  MountainView{Index: c.SitMountain, Name: c.ZuoShanName},
		FaceMountain: MountainView{Index: c.FaceMountain, Name: c.XiangShanName},
		GongWei:      gong,
		SiDaJu:       c.FourSituation,
		FuYinLayers:  c.FuYinLayers,
		FanYinLayers: c.FanYinLayers,
		ShouShanChuSha: ShouShanChuShaView{
			ShouShanOK: c.ShouShanChuSha.ShouShanOK, ChuShaOK: c.ShouShanChuSha.ChuShaOK,
			PingYu: c.ShouShanChuSha.Assessment,
		},
	}
}
