package bazi

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// LiuRiJie locates the flow day: the observer's standard clock and standard
// civil date. Values are already human-readable for direct transmission.
type LiuRiJie struct {
	StandardTime string `json:"standard_time"`
	StandardDate string `json:"standard_date"`
}

// LiuRi holds the daily flow (流日) analysis: today's pillar
// and its interactions with the bazi chart, current dayun, and current liunian.
type LiuRi struct {
	Date               string              `json:"date"`
	RiGan              ganzhi.Gan          `json:"gan"`
	RiZhi              ganzhi.Zhi          `json:"zhi"`
	DayNaYin           string              `json:"na_yin"`
	ShiShen            string              `json:"shi_shen"`
	ShengRi            int                 `json:"sheng_ri"`
	KeRi               int                 `json:"ke_ri"`
	GanRels            []GanRelation       `json:"gan_rels"`
	ZhiRels            []ZhiRelation       `json:"zhi_rels"`
	DaYunInteraction   zhuInteraction      `json:"da_yun_interaction"`
	LiuNianInteraction zhuInteraction      `json:"liu_nian_interaction"`
	ShenSha            []ShenShaProjection `json:"shen_sha"`
	RiJie              LiuRiJie            `json:"ri_jie"`
}

// ComputeLiuRi computes the day pillar for the resolved date and its full
// interactions with the bazi chart, current dayun, and current liunian.
func computeLiuRi(bz ganzhi.Bazi, t LiuRiTime, daYunZhu *ganzhi.Zhu, liuNianZhu *ganzhi.Zhu) (*LiuRi, error) {
	riYuan := bz.Ri.Gan
	bazi := bz.Slice()

	dp := ganzhi.Zhu{Gan: t.RiGan, Zhi: t.RiZhi}
	gen, rest, tgName := flowDayMasterFacts(dp.Gan, riYuan)

	// Day vs bazi (gan + zhi relations) — all 4 pillars, consistent with liunian.
	ganRels, zhiRels := analyzeZhuWithBazi(dp, bz)

	// Day vs dayun / liunian: two-layer (gan+zhi) via the shared pure function.
	var daYunRel, liuNianRel zhuInteraction
	if daYunZhu != nil {
		daYunRel = analyzeFlowAgainstRef(dp, *daYunZhu)
	}
	if liuNianZhu != nil {
		liuNianRel = analyzeFlowAgainstRef(dp, *liuNianZhu)
	}

	// Na yin.
	naYin := ganzhi.NayinLabel(dp.Gan, dp.Zhi)

	// Daily shensha: day gan/zhi vs bazi.
	var shensha []ShenShaProjection
	// 天乙贵人 on day gan.
	if targets, ok := tianYiLookup[dp.Gan]; ok {
		for _, tb := range targets {
			for _, np := range bazi {
				if np.Zhi == tb {
					shensha = append(shensha, ShenShaProjection{Name: "天乙贵人", Category: ShenShaCategoryJi, Description: "流日天乙贵人日"})
				}
			}
		}
	}
	// 文昌 on day gan.
	if targets, ok := wenChangLookup[dp.Gan]; ok {
		for _, tb := range targets {
			for _, np := range bazi {
				if np.Zhi == tb {
					shensha = append(shensha, ShenShaProjection{Name: "文昌", Category: ShenShaCategoryJi, Description: "流日文昌日，利学业文书"})
				}
			}
		}
	}
	// 驿马/桃花/华盖 from year zhi triad → day zhi check.
	yBranch := bazi[0].Zhi
	triadMaps := []struct {
		m    map[ganzhi.Zhi]ganzhi.Zhi
		name ShenShaName
		cat  ShenShaCategory
		desc string
	}{
		{yimaZhiMap, "驿马", ShenShaCategoryZhongXing, "流日驿马，动象"},
		{taohuaZhiMap, "桃花", ShenShaCategoryZhongXing, "流日桃花，异性缘佳"},
		{huagaiZhiMap, "华盖", ShenShaCategoryZhongXing, "流日华盖，宜静思"},
	}
	for _, tm := range triadMaps {
		if tb, ok := tm.m[yBranch]; ok && dp.Zhi == tb {
			shensha = append(shensha, ShenShaProjection{Name: tm.name, Category: tm.cat, Description: tm.desc})
		}
	}

	if shensha == nil {
		shensha = []ShenShaProjection{}
	}
	return &LiuRi{
		Date:               t.Date.String(),
		RiGan:              dp.Gan,
		RiZhi:              dp.Zhi,
		DayNaYin:           naYin,
		ShiShen:            tgName,
		ShengRi:            gen,
		KeRi:               rest,
		GanRels:            ganRels,
		ZhiRels:            zhiRels,
		DaYunInteraction:   daYunRel,
		LiuNianInteraction: liuNianRel,
		ShenSha:            shensha,
	}, nil
}
