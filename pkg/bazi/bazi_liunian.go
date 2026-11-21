package bazi

import (
	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// LiuNianJie locates the flow year: the observer's standard clock and civil
// selection date, plus the 立春 year boundary. Values are human-readable.
type LiuNianJie struct {
	StandardTime string `json:"standard_time"`
	StartDate    string `json:"start_date"`
	LichunYear   int    `json:"lichun_year"`
	LichunStart  string `json:"lichun_start"`
	LichunEnd    string `json:"lichun_end"`
	BoundaryDay  bool   `json:"boundary_day"`
}

// LiuNian holds the annual (流年) analysis output.
type LiuNian struct {
	Year              int                 `json:"year"`
	NianGan           ganzhi.Gan          `json:"gan"`
	NianZhi           ganzhi.Zhi          `json:"zhi"`
	NianGanWuxing     string              `json:"gan_wuxing"`
	LiuNianNaYin      string              `json:"na_yin"`
	ShiShen           string              `json:"shi_shen"`
	ShengRi           int                 `json:"sheng_ri"`
	KeRi              int                 `json:"ke_ri"`
	NatalInteractions []zhuInteraction    `json:"natal_interactions"`
	DaYunInteractions []zhuInteraction    `json:"da_yun_interactions"`
	ShenSha           []ShenShaProjection `json:"shen_sha"`
	FuYinFanYin       []FuYinFanYin       `json:"fu_yin_fan_yin"`
	NianJie           LiuNianJie          `json:"nian_jie"`
}

// ComputeLiuNian computes the year pillar for a resolved 立春 year and
// analyzes its relationship to the day master.
func computeLiuNian(bz ganzhi.Bazi, t LiuNianTime, currentDaYun *DaYunStep) (*LiuNian, error) {
	riYuan := bz.Ri.Gan
	nianGan, nianZhi := t.NianGan, t.NianZhi

	gen, rest, tgName := flowDayMasterFacts(nianGan, riYuan)

	naYin := ganzhi.NayinLabel(nianGan, nianZhi)

	r := &LiuNian{
		Year:          t.LichunYear,
		NianGan:       nianGan,
		NianZhi:       nianZhi,
		NianGanWuxing: ganzhi.DecodeWuxing(ganzhi.GanWuxing(nianGan)),
		LiuNianNaYin:  naYin,
		ShiShen:       tgName,
		ShengRi:       gen,
		KeRi:          rest,
	}

	// Three-layer analysis when bazi chart and current dayun are available.
	liuNianZhu := ganzhi.Zhu{Gan: nianGan, Zhi: nianZhi}
	r.NatalInteractions = make([]zhuInteraction, 1)
	ganRels, zhiRels := analyzeZhuWithBazi(liuNianZhu, bz)
	r.NatalInteractions[0] = zhuInteraction{
		ZhuLabel: ganzhi.GanName(nianGan) + ganzhi.ZhiName(nianZhi),
		GanRels:  ganRels,
		ZhiRels:  zhiRels,
	}

	if currentDaYun != nil {
		dyZhu := ganzhi.Zhu{Gan: currentDaYun.Gan, Zhi: currentDaYun.Zhi}
		dyGanRels, dyZhiRels := analyzeZhuWithBazi(dyZhu, bz)
		r.DaYunInteractions = []zhuInteraction{{
			ZhuLabel: currentDaYun.ShiShen + "(" + currentDaYun.Name + ")",
			GanRels:  dyGanRels,
			ZhiRels:  dyZhiRels,
		}}
	} else {
		r.DaYunInteractions = []zhuInteraction{}
	}

	r.ShenSha = computeDynamicShenSha(nianZhi, bz.Nian.Zhi, bz.Ri.Zhi, riYuan)
	r.ShenSha = append(r.ShenSha, computeAnnualShenSha(nianZhi, bz)...)
	r.FuYinFanYin = computeFuYinFanYin(liuNianZhu, bz)

	return r, nil
}

func countGenRest(elem, dmElem ganzhi.Wuxing) (gen, rest int) {
	if elem == dmElem {
		return 1, 0
	}
	if ganzhi.Sheng(elem, dmElem) {
		return 1, 0
	}
	if ganzhi.Sheng(dmElem, elem) {
		return 0, 1
	}
	if ganzhi.Ke(elem, dmElem) {
		return 0, 1
	}
	return 1, 0
}
