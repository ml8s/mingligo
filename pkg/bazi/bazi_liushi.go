package bazi

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// LiuShiChen locates the flow hour: the observer's standard clock, standard
// civil date, and the true-solar 时辰 that defines the hour pillar. Values are
// already human-readable so the domain output can be transmitted directly.
type LiuShiChen struct {
	StandardTime string `json:"standard_time"`
	StandardDate string `json:"standard_date"`
	Shichen      string `json:"shichen"`
}

// LiuShi holds the hourly flow (流时) analysis: the current two-hour period's
// pillar and its interactions with the bazi chart.
type LiuShi struct {
	Time    string        `json:"hour_range"`
	ShiGan  ganzhi.Gan    `json:"gan"`
	ShiZhi  ganzhi.Zhi    `json:"zhi"`
	ShiShen string        `json:"shi_shen"`
	GanRels []GanRelation `json:"gan_rels"`
	ZhiRels []ZhiRelation `json:"zhi_rels"`
	ShiChen LiuShiChen    `json:"shi_chen"`
}

// computeLiuShi analyzes a resolved LAST hour pillar against the bazi chart.
func computeLiuShi(bz ganzhi.Bazi, t LiuShiTime) (*LiuShi, error) {
	riYuan := bz.Ri.Gan
	shiGan := t.ShiGan
	shiZhi := t.ShiZhi

	tgName := ganzhi.ShiShenFromGan(riYuan, shiGan)
	// Hour vs bazi: all 4 pillars, consistent with liunian.
	ganRels, zhiRels := analyzeZhuWithBazi(ganzhi.Zhu{Gan: shiGan, Zhi: shiZhi}, bz)

	return &LiuShi{
		Time:    ganzhi.HourRanges[int(shiZhi)-1],
		ShiGan:  shiGan,
		ShiZhi:  shiZhi,
		ShiShen: ganzhi.DecodeShiShen(tgName),
		GanRels: ganRels,
		ZhiRels: zhiRels,
	}, nil
}
