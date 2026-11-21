package mcpcontract

import (
	"context"
	"fmt"

	"github.com/ml8s/mingligo/pkg/huangli"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

// HuangliDaysInput selects a contiguous date range and optional life event.
type HuangliDaysInput struct {
	Event        HuangliEvent `json:"event" jsonschema:"事项类型。"`
	StandardTime StandardTime `json:"standard_time" jsonschema:"查询起始日标准时间（无 DST），RFC3339 带 offset。"`
	Count        HuangliCount `json:"count" jsonschema:"黄历查询天数，最多30。"`
}

type huangliDaysResult struct {
	Days []huangli.DayView `json:"days" jsonschema:"逐日黄历视图数组，长度 = count。"`
}

// HuangliDays returns almanac facts for consecutive dates.
func HuangliDays(_ context.Context, input HuangliDaysInput) (huangliDaysResult, error) {
	standard, loc, err := parseStandardTime(input.StandardTime)
	if err != nil {
		return huangliDaysResult{}, fmt.Errorf("huangli.days: %w", err)
	}
	start := tianwen.StandardDateAt(standard, loc)
	count := int(input.Count)
	// 边界由 HuangliCount 的 schema 约束保证（1-30）；此处不再静默修正。

	days := make([]huangli.DayView, 0, count)
	for i := 0; i < count; i++ {
		date := start.AddDays(i)
		entry, err := huangli.QueryDate(date, string(input.Event))
		if err != nil {
			return huangliDaysResult{}, fmt.Errorf("huangli.days: %w", err)
		}
		days = append(days, entry.View())
	}
	return huangliDaysResult{Days: days}, nil
}
