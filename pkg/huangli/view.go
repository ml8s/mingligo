package huangli

// This file defines the domain-owned, domain-aggregated output view of a
// Huangli day. It groups related facts into历法子域 objects and is the single
// output model that every transport adapter transmits unchanged. jsonschema
// tags are inert metadata read only when a consumer reflects the type.

// CalendarView groups日历/节气 facts.
type CalendarView struct {
	Date      string `json:"date" jsonschema:"公历日期。"`
	JieQi     string `json:"jie_qi" jsonschema:"节气。"`
	JieQiDays int    `json:"jie_qi_days" jsonschema:"距节气天数。"`
	RenYuan   string `json:"ren_yuan" jsonschema:"人元。"`
}

// DirectionsView groups方位神煞 facts.
type DirectionsView struct {
	XiShen      string `json:"xi_shen" jsonschema:"喜神方位。"`
	CaiShen     string `json:"cai_shen" jsonschema:"财神方位。"`
	FuShen      string `json:"fu_shen" jsonschema:"福神方位。"`
	StemTaboo   string `json:"gan_ji" jsonschema:"干忌。"`
	BranchTaboo string `json:"zhi_ji" jsonschema:"支忌。"`
}

// JianChuHuangDaoView groups建除黄道 facts.
type JianChuHuangDaoView struct {
	JianChu  JianChu      `json:"jian_chu" jsonschema:"建除。"`
	HuangDao huangDaoStar `json:"huangdao" jsonschema:"黄道黑道十二神。"`
}

// EventView groups事项适配结论 facts.
type EventView struct {
	Event       string `json:"event,omitempty" jsonschema:"事项。"`
	EventLabel  string `json:"event_label,omitempty" jsonschema:"事项中文标签。"`
	Suitability string `json:"suitability,omitempty" jsonschema:"事项适配结论。"`
	EventReason string `json:"event_reason,omitempty" jsonschema:"事项适配依据。"`
}

// DayView is the domain-aggregated output of a Huangli day.
type DayView struct {
	Calendar        CalendarView        `json:"calendar" jsonschema:"日历与节气。"`
	RiZhu           riZhuInfo           `json:"ri_zhu" jsonschema:"日柱。"`
	Wuxing          string              `json:"wuxing" jsonschema:"日支五行。"`
	Directions      DirectionsView      `json:"directions" jsonschema:"方位神煞。"`
	JianChuHuangDao JianChuHuangDaoView `json:"jian_chu_huangdao" jsonschema:"建除黄道。"`
	Mansion         dayMansion          `json:"xiu" jsonschema:"二十八宿。"`
	ShiChen         []ShiChenFortune    `json:"shi_chen,omitempty" jsonschema:"十二时辰宜忌。"`
	Event           EventView           `json:"event" jsonschema:"事项适配。"`
}

// View produces the domain-aggregated output from the computed day.
func (d Day) View() DayView {
	return DayView{
		Calendar: CalendarView{
			Date: d.Date, JieQi: d.JieQi, JieQiDays: d.JieQiDays, RenYuan: d.RenYuan,
		},
		RiZhu:  d.RiZhu,
		Wuxing: d.Wuxing,
		Directions: DirectionsView{
			XiShen: d.XiShen, CaiShen: d.CaiShen, FuShen: d.FuShen,
			StemTaboo: d.StemTaboo, BranchTaboo: d.BranchTaboo,
		},
		JianChuHuangDao: JianChuHuangDaoView{JianChu: d.JianChu, HuangDao: d.HuangDao},
		Mansion:         d.Mansion,
		ShiChen:         d.ShiChen,
		Event: EventView{
			Event: d.Event, EventLabel: d.EventLabel,
			Suitability: d.Suitability, EventReason: d.EventReason,
		},
	}
}
