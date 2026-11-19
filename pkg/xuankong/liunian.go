package xuankong

import (
	"fmt"
	"time"

	"github.com/ml8s/mingligo/pkg/fengshui"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

// ComputeLiuNianAt computes the flow-year flying stars directly from a standard
// time. It resolves the 立春 boundary internally and fills the NianJie facts.
func ComputeLiuNianAt(standard time.Time, loc *time.Location, chart *Chart) (LiuNianResult, error) {
	ly, err := tianwen.LichunYearAt(standard)
	if err != nil {
		return LiuNianResult{}, fmt.Errorf("computeLiuNian: lichun: %w", err)
	}
	result := ComputeLiuNian(ly.Year, chart)
	if loc == nil {
		loc = time.UTC
	}
	startDate := tianwen.StandardDateAt(ly.Start, loc)
	result.NianJie = LiuNianJie{
		StandardTime: standard.In(loc).Format(time.RFC3339Nano),
		LichunYear:   ly.Year,
		LichunStart:  startDate.String(),
		LichunEnd:    tianwen.StandardDateAt(ly.End, loc).String(),
		BoundaryDay:  tianwen.StandardDateAt(standard, loc) == startDate,
	}
	return result, nil
}

// LiuNianResult is the 流年飞星 result for a year, optionally overlaid on a
// 宅盘 (chart). 全部为确定性计算：流年飞星盘来自共享 fengshui 实现，
// house_overlay 把每个流年星落宫对照宅盘该宫三星；不预设星曜固有吉凶。
// LiuNianJie locates the flow year: the observer's standard clock and the
// 立春 year boundary. Values are human-readable.
type LiuNianJie struct {
	StandardTime string `json:"standard_time"`
	LichunYear   int    `json:"lichun_year"`
	LichunStart  string `json:"lichun_start"`
	LichunEnd    string `json:"lichun_end"`
	BoundaryDay  bool   `json:"boundary_day"`
}

type LiuNianResult struct {
	Year         int                         `json:"year"`
	RuZhong      string                      `json:"ru_zhong"`
	GongWei      []fengshui.AnnualFlyingStar `json:"gong_wei"`
	HouseOverlay []HouseOverlay              `json:"house_overlay,omitempty"`
	NianJie      LiuNianJie                  `json:"nian_jie"`
}

// HouseOverlay highlights a 流年凶星落宫 and the 宅盘 stars in that palace.
type HouseOverlay struct {
	GongNum     int    `json:"gong_num"`
	Star        string `json:"xing"`         // 流年星名（如 五黄廉贞）
	PalaceStars string `json:"gong_stars"` // 宅盘该宫 运星/山星/向星
}

// ComputeLiuNian computes the annual flying-star board for a year.
// chart 可选：nil 时只返回流年飞星盘；非 nil 时叠加宅盘全量流年星落宫对照。
func ComputeLiuNian(year int, chart *Chart) LiuNianResult {
	board := fengshui.ComputeAnnualFlyingStars(year)
	res := LiuNianResult{
		Year:    board.Year,
		RuZhong: board.RuZhong,
		GongWei: board.GongWei,
	}

	if chart == nil {
		return res
	}
	for _, s := range board.GongWei {
		if s.GongNum < 1 || s.GongNum > 9 {
			continue
		}
		pal := chart.Palaces[s.GongNum-1]
		res.HouseOverlay = append(res.HouseOverlay, HouseOverlay{
			GongNum:     s.GongNum,
			Star:        s.XingName,
			PalaceStars: fmt.Sprintf("%s/%s/%s", pal.PeriodStar.Name, pal.MountainStar.Name, pal.FacingStar.Name),
		})
	}
	return res
}
