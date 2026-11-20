package huangli

import (
	"fmt"
	"slices"
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// Day is the huangli query result for one day.
type ShiChenFortune struct {
	Zhi         string `json:"zhi"`
	Time        string `json:"clock_range"` // e.g. "23:00-01:00"
	HuangDaoStr string `json:"huangdao"`
	JianChu     string `json:"jian_chu"`
	Suitable    bool   `json:"suitable"`
}

type Day struct {
	Date        string           `json:"date"`
	RiZhu       riZhuInfo        `json:"ri_zhu"`
	Wuxing      string           `json:"wuxing"`
	JianChu     string           `json:"jian_chu"`
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
}

// Month holds monthly huangli data.
func renYuanName(ry renYuanSiLing) string {
	if ry.Current == nil {
		return ""
	}
	return ry.Current.GanName
}

// QueryDate returns huangli info for a single date.
func QueryDate(dateStr string, event ...string) (Day, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return Day{}, fmt.Errorf("huangli: parse date %s: %w", dateStr, err)
	}

	dpi := lookupRiZhu(t)
	yueZhi := yueZhuForDate(t).Zhi
	jq := computeJieQiDepth(t.Year(), int(t.Month()), t.Day())
	ry := computeRenYuanSiLingForDate(yueZhi, jq.DaysIn)

	entry := Day{
		Date:        dateStr,
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
		if err := applyEventRule(&entry, event[0]); err != nil {
			return Day{}, err
		}
	}

	return entry, nil
}

func applyEventRule(entry *Day, event string) error {
	if event == "" {
		entry.Suitability = "unspecified"
		return nil
	}
	rule, ok := jianChuCfg.EventRules[event]
	if !ok {
		return fmt.Errorf("huangli: unknown event %q", event)
	}
	entry.Event = event
	entry.EventLabel = rule.Label
	switch {
	case containsJianChu(rule.Suitable, entry.JianChu):
		entry.Suitability = "recommended"
	case containsJianChu(rule.Forbidden, entry.JianChu):
		entry.Suitability = "unsuitable"
	default:
		entry.Suitability = "possible"
	}
	return nil
}

func containsJianChu(values []string, target string) bool {
	return slices.Contains(values, target)
}
