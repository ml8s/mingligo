package xuankong

// This file defines the domain-owned, domain-aggregated output view of an
// XuanKong chart. The domain Chart is already structured by术数子域; the View
// merges the mountain/facing index+name pairs into single named objects and is
// the single output model every transport adapter transmits unchanged.

// MountainView is a 24-mountain坐向 with its index and display name.
type MountainView struct {
	Index int    `json:"index" jsonschema:"二十四山序号(0-23)。"`
	Name  string `json:"name" jsonschema:"山名。"`
}

// ChartView is the domain-aggregated output of an XuanKong chart.
type ChartView struct {
	Yun            SanYuanYun        `json:"yun" jsonschema:"三元九运。"`
	SitMountain    MountainView      `json:"zuo_shan" jsonschema:"坐山。"`
	FaceMountain   MountainView      `json:"xiang_shan" jsonschema:"向首。"`
	GongWei        [9]xuanKongStar   `json:"gong_wei" jsonschema:"九宫运星山星向星。"`
	FourSituation  FourSituation     `json:"four_situation" jsonschema:"四大局。"`
	FuYinLayers    []RecitationLayer `json:"fu_yin_layers" jsonschema:"伏吟层。"`
	FanYinLayers   []RecitationLayer `json:"fan_yin_layers" jsonschema:"反吟层。"`
	ShouShanChuSha shouShanChuSha    `json:"shou_shan_chu_sha" jsonschema:"收山出煞判断。"`
}

// View produces the domain-aggregated output from the computed chart.
func (c Chart) View() ChartView {
	return ChartView{
		Yun:            c.Yun,
		SitMountain:    MountainView{Index: c.SitMountain, Name: c.ZuoShanName},
		FaceMountain:   MountainView{Index: c.FaceMountain, Name: c.XiangShanName},
		GongWei:        c.Palaces,
		FourSituation:  c.FourSituation,
		FuYinLayers:    c.FuYinLayers,
		FanYinLayers:   c.FanYinLayers,
		ShouShanChuSha: c.ShouShanChuSha,
	}
}
