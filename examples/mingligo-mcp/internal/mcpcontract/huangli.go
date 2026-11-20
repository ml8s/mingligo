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

type HuangliDaysOutput struct {
	Days []HuangliDayOutput `json:"days" jsonschema:"逐日结果数组，长度 = count。"`
}

type HuangliRiZhuOutput struct {
	Gan   string `json:"gan" jsonschema:"日干。"`
	Zhi   string `json:"zhi" jsonschema:"日支。"`
	NaYin string `json:"na_yin" jsonschema:"日柱纳音。"`
}

type HuangliHuangDaoOutput struct {
	Index    int    `json:"index" jsonschema:"黄道/黑道序号。"`
	Name     string `json:"name" jsonschema:"黄道/黑道神名。"`
	Path     string `json:"path" jsonschema:"黄道或黑道。"`
	Sequence int    `json:"sequence" jsonschema:"十二神循环位置。"`
}

type HuangliMansionOutput struct {
	Index    int    `json:"index" jsonschema:"二十八宿序号。"`
	Name     string `json:"name" jsonschema:"二十八宿名称。"`
	DongWu   string `json:"dong_wu" jsonschema:"七曜动物。"`
	Wuxing   string `json:"wuxing" jsonschema:"宿五行。"`
	Group    string `json:"group" jsonschema:"四象分组。"`
	GroupIdx int    `json:"group_idx" jsonschema:"四象内序号。"`
}

type HuangliShiChenOutput struct {
	Zhi        string `json:"zhi" jsonschema:"时辰地支。"`
	ClockRange string `json:"clock_range" jsonschema:"时辰钟点区间。"`
	HuangDao   string `json:"huangdao" jsonschema:"该时辰黄道/黑道神。"`
	JianChu    string `json:"jian_chu" jsonschema:"该时辰建除。"`
	Suitable   bool   `json:"suitable" jsonschema:"该时辰是否宜。"`
}

// HuangliDayOutput is the schema carrier for one almanac day. Its field names
// match the domain Day JSON; the handler returns the domain value directly.
type HuangliDayOutput struct {
	Date        string                 `json:"date" jsonschema:"日期 YYYY-MM-DD。"`
	RiZhu       HuangliRiZhuOutput     `json:"ri_zhu" jsonschema:"日柱信息。"`
	Wuxing      string                 `json:"wuxing" jsonschema:"日支五行。"`
	JianChu     string                 `json:"jian_chu" jsonschema:"建除十二神。"`
	HuangDao    HuangliHuangDaoOutput  `json:"huangdao" jsonschema:"当日黄道/黑道。"`
	XiShen      string                 `json:"xi_shen" jsonschema:"喜神方位。"`
	CaiShen     string                 `json:"cai_shen" jsonschema:"财神方位。"`
	FuShen      string                 `json:"fu_shen" jsonschema:"福神方位。"`
	StemTaboo   string                 `json:"gan_ji" jsonschema:"彭祖日干忌。"`
	BranchTaboo string                 `json:"zhi_ji" jsonschema:"彭祖日支忌。"`
	Mansion     HuangliMansionOutput   `json:"xiu" jsonschema:"当日二十八宿。"`
	JieQi       string                 `json:"jie_qi" jsonschema:"节气。"`
	JieQiDays   int                    `json:"jie_qi_days" jsonschema:"进入节气天数。"`
	RenYuan     string                 `json:"ren_yuan" jsonschema:"人元司令天干。"`
	ShiChen     []HuangliShiChenOutput `json:"shi_chen,omitempty" jsonschema:"十二时辰宜忌。"`
	Event       string                 `json:"event,omitempty" jsonschema:"查询事项。"`
	EventLabel  string                 `json:"event_label,omitempty" jsonschema:"事项中文标签。"`
	Suitability string                 `json:"suitability,omitempty" jsonschema:"事项适配结论。"`
	EventReason string                 `json:"event_reason,omitempty" jsonschema:"事项适配依据。"`
}

// huangliDaysResult is the direct domain output; HuangliDaysOutput is only the
// published schema carrier. The field names match the wire contract.
type huangliDaysResult struct {
	Days []huangli.Day `json:"days"`
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

	days := make([]huangli.Day, 0, count)
	for i := 0; i < count; i++ {
		date := start.AddDays(i).String()
		entry, err := huangli.QueryDate(date, string(input.Event))
		if err != nil {
			return huangliDaysResult{}, fmt.Errorf("huangli.days: %w", err)
		}
		days = append(days, entry)
	}
	return huangliDaysResult{Days: days}, nil
}

