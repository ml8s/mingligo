package bazi

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// LiuYueJieQi locates the flow month: the observer's standard clock and civil
// selection date, plus the 节 boundary. Values are human-readable.
type LiuYueJieQi struct {
	StandardTime string `json:"standard_time"`
	StartDate    string `json:"start_date"`
	JieqiStart   string `json:"jieqi_start"`
	JieqiEnd     string `json:"jieqi_end"`
	BoundaryDay  bool   `json:"boundary_day"`
}

// LiuYue holds the monthly (流月) analysis output.
type LiuYue struct {
	Year               int                 `json:"year"`
	Month              int                 `json:"month"`
	YueGan             ganzhi.Gan          `json:"gan"`
	YueZhi             ganzhi.Zhi          `json:"zhi"`
	YueGanWuxing       string              `json:"gan_wuxing"`
	ShiShen            string              `json:"shi_shen"`
	ShengRi            int                 `json:"sheng_ri"`
	KeRi               int                 `json:"ke_ri"`
	GanRels            []GanRelation       `json:"gan_rels"`
	ZhiRels            []ZhiRelation       `json:"zhi_rels"`
	DaYunInteraction   zhuInteraction      `json:"da_yun_interaction"`
	LiuNianInteraction zhuInteraction      `json:"liu_nian_interaction"`
	ShenSha            []ShenShaProjection `json:"shen_sha"`
	JieQi              LiuYueJieQi         `json:"jie_qi"`
}

// ComputeLiuYue computes the month pillar for a resolved solar-term month and
// analyzes its interactions with the bazi chart.
func computeLiuYue(bz ganzhi.Bazi, t LiuYueTime, refs flowRefs) (*LiuYue, error) {
	riYuan := bz.Ri.Gan
	bazi := bz.Slice()

	mp := ganzhi.Zhu{Gan: t.YueGan, Zhi: t.YueZhi}

	gen, rest, tgName := flowDayMasterFacts(mp.Gan, riYuan)

	// Month vs bazi: all 4 pillars, consistent with liunian.
	ganRels, zhiRels := analyzeZhuWithBazi(mp, bz)
	var daYunRel, liuNianRel zhuInteraction
	if refs.dayun != nil {
		daYunRel = analyzeFlowAgainstRef(mp, *refs.dayun)
	}
	if refs.liunian != nil {
		liuNianRel = analyzeFlowAgainstRef(mp, *refs.liunian)
	}

	shensha := computeDynamicShenSha(mp.Zhi, bazi[0].Zhi, bazi[2].Zhi, riYuan)

	monthElemStr := ganzhi.DecodeWuxing(ganzhi.GanWuxing(mp.Gan))

	return &LiuYue{
		Year:               t.Date.Year,
		Month:              t.Date.Month,
		YueGan:             mp.Gan,
		YueZhi:             mp.Zhi,
		YueGanWuxing:       monthElemStr,
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
