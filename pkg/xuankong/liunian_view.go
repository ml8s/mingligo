package xuankong

// Domain-owned output view of the xuankong.liunian series. jsonschema tags
// are inert metadata read only when a consumer reflects the type for its schema.

type LiuNianJieView struct {
	StandardTime string `json:"standard_time" jsonschema:"流年选择器时区标准时间（无 DST）。"`
	LichunYear   int    `json:"lichun_year" jsonschema:"该瞬间所在的立春年。"`
	LichunStart  string `json:"lichun_start" jsonschema:"立春年起始公历日（含）。"`
	LichunEnd    string `json:"lichun_end" jsonschema:"下一立春日（不含）。"`
	BoundaryDay  bool   `json:"boundary_day" jsonschema:"输入瞬间的当地日期是否为立春边界日。"`
}

type AnnualStarView struct {
	GongNum  int    `json:"gong_shu" jsonschema:"洛书宫数。"`
	Xing     int    `json:"xing" jsonschema:"紫白星数。"`
	XingName string `json:"xing_ming" jsonschema:"紫白星名。"`
	Wuxing   string `json:"wuxing" jsonschema:"星五行。"`
	RuZhong  bool   `json:"ru_zhong" jsonschema:"是否入中。"`
}

type HouseOverlayView struct {
	GongNum     int    `json:"gong_shu" jsonschema:"洛书宫数。"`
	Star        string `json:"xing" jsonschema:"流年星名。"`
	PalaceStars string `json:"gong_stars" jsonschema:"宅盘该宫运星/山星/向星。"`
}

type LiuNianOutputView struct {
	Year               int                `json:"nian" jsonschema:"立春年份。"`
	RuZhong            string             `json:"ru_zhong" jsonschema:"入中紫白星。"`
	GongWei            [9]AnnualStarView  `json:"gong_wei" jsonschema:"九宫流年紫白星。"`
	LiuNianXingLinGong []HouseOverlayView `json:"liu_nian_xing_lin_gong,omitempty" jsonschema:"全部流年星落宫对照宅盘该宫三星；不预设星曜固有吉凶。"`
	NianJie            LiuNianJieView     `json:"nian_jie" jsonschema:"玄空流年立春年界。"`
}

func (v LiuNianResult) View() LiuNianOutputView {
	view := LiuNianOutputView{
		Year: v.Year, RuZhong: v.RuZhong,
		GongWei: [9]AnnualStarView{},
		NianJie: LiuNianJieView{
			StandardTime: v.NianJie.StandardTime, LichunYear: v.NianJie.LichunYear,
			LichunStart: v.NianJie.LichunStart, LichunEnd: v.NianJie.LichunEnd,
			BoundaryDay: v.NianJie.BoundaryDay,
		},
	}
	for i, star := range v.GongWei {
		if i >= 9 {
			break
		}
		view.GongWei[i] = AnnualStarView{
			GongNum: star.GongNum, Xing: star.Xing, XingName: star.XingName,
			Wuxing: star.Wuxing, RuZhong: star.RuZhong,
		}
	}
	for _, overlay := range v.HouseOverlay {
		view.LiuNianXingLinGong = append(view.LiuNianXingLinGong, HouseOverlayView(overlay))
	}
	return view
}
