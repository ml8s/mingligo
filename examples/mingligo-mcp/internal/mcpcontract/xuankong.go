package mcpcontract

import (
	"context"
	"fmt"
	"time"

	"github.com/ml8s/mingligo/pkg/tianwen"
	"github.com/ml8s/mingligo/pkg/xuankong"
)

type XuankongChartInput struct {
	StandardTime StandardTime  `json:"standard_time" jsonschema:"宅运起盘地标准日期时间（无 DST），RFC3339 带 offset。"`
	ZuoShan      MountainIndex `json:"zuo_shan" jsonschema:"坐山，二十四山 index。"`
	XiangShan    MountainIndex `json:"xiang_shan" jsonschema:"朝向；必须与坐山相对。"`
}

// XuankongChart is the chart tool's public output. It is not a reusable method
// input; callers pass standard_time and zuo/xiang again for another fact.
type XuankongLiunianInput struct {
	StandardTime StandardTime  `json:"standard_time" jsonschema:"流年选择地标准时间（无 DST），RFC3339 带 offset。"`
	ZuoShan      MountainIndex `json:"zuo_shan" jsonschema:"坐山，二十四山 index。"`
	XiangShan    MountainIndex `json:"xiang_shan" jsonschema:"朝向；必须与坐山相对。"`
	LiuNianCount *SeriesCount  `json:"count,omitempty" jsonschema:"连续查询年数。"`
}

type XuankongLiuNianJie struct {
	StandardTime StandardTime `json:"standard_time" jsonschema:"流年选择器时区标准时间（无 DST）。"`
	LichunYear   int          `json:"lichun_year" jsonschema:"该瞬间所在的立春年。"`
	LichunStart  DateString   `json:"lichun_start" jsonschema:"立春年起始公历日（含）。"`
	LichunEnd    DateString   `json:"lichun_end" jsonschema:"下一立春日（不含）。"`
	BoundaryDay  bool         `json:"boundary_day" jsonschema:"输入瞬间的当地日期是否为立春边界日。"`
}

type XuankongAnnualStarOutput struct {
	GongNum  GuaNumber        `json:"gong_num" jsonschema:"洛书宫数。"`
	Xing     FlyingStarNumber `json:"xing" jsonschema:"紫白星数。"`
	XingName FlyingStarName   `json:"xing_name" jsonschema:"紫白星名。"`
	Wuxing   ElementName      `json:"wuxing" jsonschema:"星五行。"`
	RuZhong  bool             `json:"ru_zhong" jsonschema:"是否入中。"`
}

type XuankongHouseOverlayOutput struct {
	GongNum     GuaNumber      `json:"gong_num" jsonschema:"洛书宫数。"`
	Star        FlyingStarName `json:"xing" jsonschema:"流年星名。"`
	PalaceStars string         `json:"gong_stars" jsonschema:"宅盘该宫运星/山星/向星。"`
}

type XuankongLiunianOutput struct {
	Year         XuankongYear                 `json:"year" jsonschema:"立春年份。"`
	RuZhong      FlyingStarName               `json:"ru_zhong" jsonschema:"入中紫白星。"`
	GongWei      [9]XuankongAnnualStarOutput  `json:"gong_wei" jsonschema:"九宫流年紫白星。"`
	HouseOverlay []XuankongHouseOverlayOutput `json:"house_overlay,omitempty" jsonschema:"全部流年星落宫对照宅盘该宫三星；不预设星曜固有吉凶。"`
	NianJie      XuankongLiuNianJie           `json:"nian_jie" jsonschema:"玄空流年立春年界。"`
}

type XuankongLiuNianSeries struct {
	LiuNian []xuankong.LiuNianResult `json:"liu_nian" jsonschema:"玄空流年数组，长度 = count。"`
}

func XuankongChartCompute(_ context.Context, input XuankongChartInput) (xuankong.ChartView, error) {
	at, _, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return xuankong.ChartView{}, fmt.Errorf("xuankong.chart: %w", err)
	}
	chart, err := xuankongComputeChart(at, input.ZuoShan, input.XiangShan)
	if err != nil {
		return xuankong.ChartView{}, fmt.Errorf("xuankong.chart: %w", err)
	}
	return chart.View(), nil
}

func XuankongLiunian(_ context.Context, input XuankongLiunianInput) (XuankongLiuNianSeries, error) {
	items, err := runSeries(input, seriesCountValue(input.LiuNianCount), MaxSeriesCount,
		func(in XuankongLiunianInput) (xuankong.LiuNianResult, error) {
			return computeXuankongLiunian(in)
		},
		func(in XuankongLiunianInput) (XuankongLiunianInput, error) {
			at, loc, err := parseStandardTime(in.StandardTime)
			if err != nil {
				return in, fmt.Errorf("xuankong.liunian: %w", err)
			}
			ly, err := tianwen.LichunYearAt(at)
			if err != nil {
				return in, fmt.Errorf("xuankong.liunian: %w", err)
			}
			in.StandardTime = standardTimeValue(ly.End.In(loc))
			return in, nil
		},
	)
	if err != nil {
		return XuankongLiuNianSeries{}, err
	}
	return XuankongLiuNianSeries{LiuNian: items}, nil
}

func computeXuankongLiunian(input XuankongLiunianInput) (xuankong.LiuNianResult, error) {
	at, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return xuankong.LiuNianResult{}, fmt.Errorf("xuankong.liunian: %w", err)
	}
	chart, err := xuankongComputeChart(at, input.ZuoShan, input.XiangShan)
	if err != nil {
		return xuankong.LiuNianResult{}, fmt.Errorf("xuankong.liunian: %w", err)
	}
	result, err := xuankong.ComputeLiuNianAt(at, loc, &chart)
	if err != nil {
		return xuankong.LiuNianResult{}, fmt.Errorf("xuankong.liunian: %w", err)
	}
	return result, nil
}

func xuankongComputeChart(at time.Time, zuoShan MountainIndex, xiangShan MountainIndex) (xuankong.Chart, error) {
	if at.UTC().Year() < 1864 || at.UTC().Year() > 2200 {
		return xuankong.Chart{}, fmt.Errorf("year %d out of range (1864-2200)", at.UTC().Year())
	}
	if (int(zuoShan)+12)%24 != int(xiangShan) {
		return xuankong.Chart{}, fmt.Errorf(
			"xiang_shan must be opposite zuo_shan by 12 mountains, got %d/%d",
			zuoShan, xiangShan,
		)
	}
	return xuankong.ComputeChartAt(at, int(zuoShan), int(xiangShan))
}
