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
	Time               string         `json:"hour_range"`
	ShiGan             ganzhi.Gan     `json:"gan"`
	ShiZhi             ganzhi.Zhi     `json:"zhi"`
	ShiShen            string         `json:"shi_shen"`
	ShengRi            int            `json:"sheng_ri"`
	KeRi               int            `json:"ke_ri"`
	GanRels            []GanRelation  `json:"gan_rels"`
	ZhiRels            []ZhiRelation  `json:"zhi_rels"`
	DaYunInteraction   zhuInteraction `json:"da_yun_interaction"`
	LiuNianInteraction zhuInteraction `json:"liu_nian_interaction"`
	LiuRiInteraction   zhuInteraction `json:"liu_ri_interaction"`
	ShiChen            LiuShiChen     `json:"shi_chen"`
}

// computeLiuShi analyzes a resolved LAST hour pillar against the bazi chart.
func computeLiuShi(bz ganzhi.Bazi, t LiuShiTime, refs flowRefs) (*LiuShi, error) {
	riYuan := bz.Ri.Gan
	shiGan := t.ShiGan
	shiZhi := t.ShiZhi

	gen, rest, tgName := flowDayMasterFacts(shiGan, riYuan)
	// Hour vs bazi: all 4 pillars, consistent with liunian.
	hourZhu := ganzhi.Zhu{Gan: shiGan, Zhi: shiZhi}
	ganRels, zhiRels := analyzeZhuWithBazi(hourZhu, bz)
	var daYunRel, liuNianRel, liuRiRel zhuInteraction
	if refs.dayun != nil {
		daYunRel = analyzeFlowAgainstRef(hourZhu, *refs.dayun)
	}
	if refs.liunian != nil {
		liuNianRel = analyzeFlowAgainstRef(hourZhu, *refs.liunian)
	}
	if refs.liuri != nil {
		liuRiRel = analyzeFlowAgainstRef(hourZhu, *refs.liuri)
	}

	return &LiuShi{
		Time:               ganzhi.HourRanges[int(shiZhi)-1],
		ShiGan:             shiGan,
		ShiZhi:             shiZhi,
		ShiShen:            tgName,
		ShengRi:            gen,
		KeRi:               rest,
		GanRels:            ganRels,
		ZhiRels:            zhiRels,
		DaYunInteraction:   daYunRel,
		LiuNianInteraction: liuNianRel,
		LiuRiInteraction:   liuRiRel,
	}, nil
}
