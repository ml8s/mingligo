package huangli

import (
	"fmt"
	"slices"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

// Day is the huangli query result for one day.
type ShiChenFortune struct {
	Zhi              string  `json:"zhi"`
	Time             string  `json:"clock_range"` // e.g. "23:00-01:00"
	HuangDaoName     string  `json:"huang_dao_name"`
	JianChu          JianChu `json:"jian_chu"`
	HuangDaoSuitable bool    `json:"huang_dao_suitable"`
}

type Day struct {
	Date        string           `json:"date"`
	RiZhu       riZhuInfo        `json:"ri_zhu"`
	Wuxing      string           `json:"wuxing"`
	JianChu     JianChu          `json:"jian_chu"`
	HuangDao    huangDaoStar     `json:"huangdao"`
	XiShen      string           `json:"xi_shen"`
	CaiShen     string           `json:"cai_shen"`
	FuShen      string           `json:"fu_shen"`
	StemTaboo   string           `json:"gan_ji"`
	BranchTaboo string           `json:"zhi_ji"`
	Mansion     dayMansion       `json:"xiu"`
	JieQi       string           `json:"jie_qi"`
	JieQiDays   int              `json:"jie_qi_days"`
	RenYuan     string           `json:"ren_yuan"`
	ShiChen     []ShiChenFortune `json:"shi_chen,omitempty"`
	Event       string           `json:"event,omitempty"`
	EventLabel  string           `json:"event_label,omitempty"`
	Suitability string           `json:"suitability,omitempty"`
	EventReason string           `json:"event_reason,omitempty"`
}

// Month holds monthly huangli data.
func renYuanName(ry renYuanSiLing) string {
	if ry.Current == nil {
		return ""
	}
	return ry.Current.GanName
}

// QueryDate returns huangli info for a single civil date.
func QueryDate(date tianwen.Date, event ...string) (Day, error) {
	if err := date.Validate(); err != nil {
		return Day{}, fmt.Errorf("huangli: %w", err)
	}
	t := date.Time()

	dpi := lookupRiZhu(t)
	yueZhi := yueZhuForDate(t).Zhi
	jq := computeJieQiDepth(t.Year(), int(t.Month()), t.Day())
	ry := computeRenYuanSiLingForDate(yueZhi, jq.DaysIn)

	entry := Day{
		Date:        date.String(),
		RiZhu:       dpi,
		JianChu:     lookupJianChu(t),
		HuangDao:    huangDaoForDay(yueZhi, dpi.Zhi),
		XiShen:      xiShenDirection(dpi.Gan),
		CaiShen:     caiShenDirection(dpi.Gan),
		FuShen:      fuShenDirection(dpi.Gan),
		StemTaboo:   pengZuGanTaboo(dpi.Gan),
		BranchTaboo: pengZuZhiTaboo(dpi.Zhi),
		Wuxing:      ganzhi.DecodeWuxing(ganzhi.ZhiWuxing(dpi.Zhi)),
		Mansion:     mansionForDay(ganzhi.Zhu{Gan: dpi.Gan, Zhi: dpi.Zhi}),
		JieQi:       jq.TermName,
		JieQiDays:   jq.DaysIn,
		RenYuan:     renYuanName(ry),
	}

	entry.ShiChen = computeShiChen(dpi.Zhi, yueZhi, entry.JianChu)
	if len(event) > 0 {
		var err error
		entry, err = withEventRule(entry, event[0])
		if err != nil {
			return Day{}, err
		}
	}

	return entry, nil
}

func withEventRule(entry Day, event string) (Day, error) {
	if event == "" {
		entry.Suitability = "unspecified"
		entry.EventReason = "未指定事项，仅列黄历事实。"
		return entry, nil
	}
	rule, ok := jianChuCfg.EventRules[event]
	if !ok {
		return Day{}, fmt.Errorf("huangli: unknown event %q", event)
	}
	entry.Event = event
	entry.EventLabel = rule.Label
	switch {
	case containsJianChu(rule.Suitable, entry.JianChu):
		entry.Suitability = "recommended"
		entry.EventReason = "建除「" + DecodeJianChu(entry.JianChu) + "」适合" + rule.Label + "。"
	case containsJianChu(rule.Forbidden, entry.JianChu):
		entry.Suitability = "unsuitable"
		entry.EventReason = "建除「" + DecodeJianChu(entry.JianChu) + "」忌" + rule.Label + "。"
	default:
		entry.Suitability = "possible"
		entry.EventReason = "建除「" + DecodeJianChu(entry.JianChu) + "」无明确适配或冲突。"
	}
	return entry, nil
}

func containsJianChu(values []JianChu, target JianChu) bool {
	return slices.Contains(values, target)
}
