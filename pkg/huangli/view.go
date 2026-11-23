package huangli

import "github.com/ml8s/mingligo/pkg/ganzhi"

// This file defines the domain-owned, domain-aggregated output view of a
// Huangli day. Every published type is exported and every key is a命理拼音
// term; computation types are mapped explicitly and never leak into the wire.

// LiFaView groups日历/节气 facts.
type LiFaView struct {
	Date      string `json:"date" jsonschema:"公历日期。"`
	JieQi     string `json:"jie_qi" jsonschema:"节气。"`
	JieQiDays int    `json:"jie_qi_ri_shu" jsonschema:"距节气天数。"`
}

// RenYuanView groups 人元司令 facts. 人元分野是月令司事（命理域），不是历法。
type RenYuanView struct {
	DangLing string `json:"dang_ling" jsonschema:"当前司事天干（人元司令）。"`
}

// RiZhuView is the day pillar.
type RiZhuView struct {
	Gan    ganzhi.Gan `json:"gan" jsonschema:"日干。"`
	Zhi    ganzhi.Zhi `json:"zhi" jsonschema:"日支。"`
	NaYin  string     `json:"na_yin" jsonschema:"纳音。"`
	Wuxing string     `json:"wuxing" jsonschema:"日支五行。"`
}

// FangWeiShenShaView groups方位神煞 facts.
type FangWeiShenShaView struct {
	XiShen  string `json:"xi_shen" jsonschema:"喜神方位。"`
	CaiShen string `json:"cai_shen" jsonschema:"财神方位。"`
	FuShen  string `json:"fu_shen" jsonschema:"福神方位。"`
}

// PengZuBaiJiView groups 彭祖百忌 facts (day-stem/branch prohibitions).
type PengZuBaiJiView struct {
	StemTaboo   string `json:"gan_ji" jsonschema:"干忌。"`
	BranchTaboo string `json:"zhi_ji" jsonschema:"支忌。"`
}

// HuangDaoView is the黄道/黑道十二神 for the day.
type HuangDaoView struct {
	Xu    int    `json:"xu" jsonschema:"十二神序号。"`
	Ming  string `json:"ming" jsonschema:"神名。"`
	Dao   string `json:"dao" jsonschema:"黄道或黑道。"`
	XuWei int    `json:"xu_wei" jsonschema:"十二神cycle位置。"`
}

// JianChuHuangDaoView groups建除黄道 facts.
type JianChuHuangDaoView struct {
	JianChu  JianChu      `json:"jian_chu" jsonschema:"建除。"`
	HuangDao HuangDaoView `json:"huang_dao" jsonschema:"黄道黑道十二神。"`
}

// XiuView is one of the 二十八宿.
type XiuView struct {
	Xu     int    `json:"xu" jsonschema:"宿序号。"`
	Ming   string `json:"ming" jsonschema:"宿名。"`
	DongWu string `json:"dong_wu" jsonschema:"动物。"`
	Wuxing string `json:"wuxing" jsonschema:"五行。"`
	Zu     string `json:"zu" jsonschema:"七宿组。"`
	ZuXu   int    `json:"zu_xu" jsonschema:"组内序号。"`
}

// YiJiView groups事项宜忌 facts with显式 domain keys.
type YiJiView struct {
	ShiXiang     string `json:"shi_xiang,omitempty" jsonschema:"事项。"`
	ShiXiangMing string `json:"shi_xiang_ming,omitempty" jsonschema:"事项中文名称。"`
	ShiYi        string `json:"shi_yi,omitempty" jsonschema:"事项适配结论。"`
	YiJu         string `json:"yi_ju,omitempty" jsonschema:"事项适配依据。"`
}

// DayView is the domain-aggregated output of a Huangli day.
type DayView struct {
	LiFa            LiFaView            `json:"li_fa" jsonschema:"历法与节气。"`
	RenYuan         RenYuanView         `json:"ren_yuan" jsonschema:"人元司令（命理域）。"`
	RiZhu           RiZhuView           `json:"ri_zhu" jsonschema:"日柱。"`
	FangWeiShenSha  FangWeiShenShaView  `json:"fang_wei_shen_sha" jsonschema:"方位神煞。"`
	PengZuBaiJi     PengZuBaiJiView     `json:"peng_zu_bai_ji" jsonschema:"彭祖百忌。"`
	JianChuHuangDao JianChuHuangDaoView `json:"jian_chu_huang_dao" jsonschema:"建除黄道。"`
	Xiu             XiuView             `json:"xiu" jsonschema:"二十八宿。"`
	ShiChen         []ShiChenFortune    `json:"shi_chen,omitempty" jsonschema:"十二时辰宜忌。"`
	YiJi            YiJiView            `json:"yi_ji" jsonschema:"事项宜忌。"`
}

// View produces the domain-aggregated output from the computed day.
func (d Day) View() DayView {
	return DayView{
		LiFa: LiFaView{
			Date: d.Date, JieQi: d.JieQi, JieQiDays: d.JieQiDays,
		},
		RenYuan: RenYuanView{DangLing: d.RenYuan},
		RiZhu: RiZhuView{
			Gan: d.RiZhu.Gan, Zhi: d.RiZhu.Zhi, NaYin: d.RiZhu.NaYin, Wuxing: d.Wuxing,
		},
		FangWeiShenSha: FangWeiShenShaView{
			XiShen: d.XiShen, CaiShen: d.CaiShen, FuShen: d.FuShen,
		},
		PengZuBaiJi: PengZuBaiJiView{StemTaboo: d.StemTaboo, BranchTaboo: d.BranchTaboo},
		JianChuHuangDao: JianChuHuangDaoView{
			JianChu: d.JianChu,
			HuangDao: HuangDaoView{
				Xu: d.HuangDao.Index, Ming: d.HuangDao.Name,
				Dao: d.HuangDao.Path, XuWei: d.HuangDao.Sequence,
			},
		},
		Xiu: XiuView{
			Xu: d.Mansion.Index, Ming: d.Mansion.Name, DongWu: d.Mansion.DongWu,
			Wuxing: d.Mansion.Element, Zu: d.Mansion.Group, ZuXu: d.Mansion.GroupIdx,
		},
		ShiChen: d.ShiChen,
		YiJi: YiJiView{
			ShiXiang: d.Event, ShiXiangMing: d.EventLabel,
			ShiYi: d.Suitability, YiJu: d.EventReason,
		},
	}
}
