package qimen

import (
	"time"

	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/luoshu"
	"github.com/ml8s/mingligo/pkg/tianwen"
)

// Chart bundles a complete奇门盘 with all analysis layers.
type Chart struct {
	Method           ChartMethod        `json:"method"`
	Board            Board              `json:"pan"`
	GanInteractions  []GanInteraction   `json:"gan_interaction"`
	MenInteractions  []MenInteraction   `json:"men_interaction"`
	XingInteractions []XingInteraction  `json:"xing_interaction"`
	XingGongWuXing   []XingGongWuXing   `json:"xing_gong_wuxing"`
	MenPo            []GongIndex        `json:"men_po"`
	MenZhi           []GongIndex        `json:"men_zhi"`
	Patterns         []Pattern          `json:"ge_ju"`
	YingQi           YingQi             `json:"ying_qi"`
	DayStemPalace    StemPalace         `json:"ri_gan_gong"`
	HourStemPalace   StemPalace         `json:"shi_gan_gong"`
	RiShiRelation    RiShiRelation      `json:"ri_shi_relation"`
	KongWangAffected []AffectedSymbol   `json:"kong_wang_affected"`
	MaXingAffected   []AffectedSymbol   `json:"ma_xing_affected"`
	GongWangShuai    []GongWangShuai    `json:"gong_wang_shuai"`
	DutyStarPalace   GongIndex          `json:"zhi_fu_xing_gong"`    // 值符星落宫（排盘固有）
	DutyDoorPalace   GongIndex          `json:"zhi_shi_men_gong"`    // 值使门落宫（排盘固有）
	YongShen         *YongShenResult    `json:"yong_shen,omitempty"` // 用神领域对象（求测人+事象用神）
	Specialized      SpecializedContext `json:"special_facts"`
}

// StemPalace is the canonical view of one stem's placement: the primary
// heaven-first palace, all heaven/earth occurrences, and seasonal state.
type StemPalace struct {
	Palace    GongIndex     `json:"gong"`
	Facts     []GanGongFact `json:"facts"`
	WangShuai GongWangShuai `json:"wang_shuai"`
}

// GanGongFact records every palace layer occupied by a stem. The primary
// primary palace keeps the heaven-first rule; facts preserve
// the earth-stem palace that the single value previously discarded.
type GanGongFact struct {
	Palace GongIndex `json:"gong"`
	Layer  string    `json:"layer"` // heaven / earth
}

type RiShiRelation struct {
	Zhu      string `json:"zhu"`
	Ke       string `json:"ke"`
	Relation string `json:"relation"`
	Name     string `json:"name"`
}

type AffectedSymbol struct {
	Symbol string     `json:"symbol"`
	Role   string     `json:"role"`
	Branch ganzhi.Zhi `json:"zhi"`
	Gong   GongIndex  `json:"gong"`
}

// GongWangShuai 表示宫位五行在当前月令下的旺相休囚死状态。
type GongWangShuai struct {
	Gong      GongIndex        `json:"gong"`
	Wuxing    ganzhi.Wuxing    `json:"wuxing"`
	WangShuai ganzhi.WangShuai `json:"wang_shuai"`
}

// computeChart computes a qimen chart with the base timing focus.
func computeChart(bz ganzhi.Bazi, termTime, chartTime time.Time, method Method) Chart {
	chart := buildChart(bz, termTime, chartTime, method)
	chart.YingQi = computeYingQi(chart, chartTime, baseYingQiFocuses(chart))
	return chart
}

// computeChartWithYongShen computes the chart and projects timing through the
// selected symbols without first building and discarding a base timing result.
func computeChartWithYongShen(
	bz ganzhi.Bazi, termTime, chartTime time.Time, method Method,
	syms []YongShenSymbol, birthDate BirthDate,
) Chart {
	chart := buildChart(bz, termTime, chartTime, method)
	var yongShen YongShenResult
	if birthDate.Has {
		yongShen = computeYongShenWithBirth(chart, syms, birthDate)
	} else {
		yongShen = computeYongShenSymbols(chart, syms)
	}
	chart.YongShen = &yongShen
	chart.YingQi = computeYingQi(chart, chartTime, yingQiFocuses(chart, yongShen))
	return chart
}

func buildChart(bz ganzhi.Bazi, termTime, chartTime time.Time, method Method) Chart {
	ju := DetermineJuShu(termTime, chartTime, bz, method)
	p := ArrangePan(ju, bz)

	riGanP := findGanPalaceIdx(p, resolveJiaDunGan(bz.Ri.Gan, bz.Ri.Zhi))
	shiGanP := findGanPalaceIdx(p, resolveJiaDunGan(bz.Shi.Gan, bz.Shi.Zhi))
	riGanFacts := findGanGongFacts(p, resolveJiaDunGan(bz.Ri.Gan, bz.Ri.Zhi))
	shiGanFacts := findGanGongFacts(p, resolveJiaDunGan(bz.Shi.Gan, bz.Shi.Zhi))
	leadP := leadPillarPalace(p)
	wangShuaiList := palaceWangShuai(p, bz.Yue.Zhi)

	chart := Chart{
		Method:           buildChartMethod(ju, bz),
		Board:            p,
		GanInteractions:  computeGanInteractions(p),
		MenInteractions:  computeMenInteractions(p),
		XingInteractions: computeXingInteractions(p),
		XingGongWuXing:   computeXingGongWuXing(p),
		MenPo:            findMenPo(p),
		MenZhi:           findMenZhi(p),
		Patterns:         findPatterns(p),
		DayStemPalace:    newStemPalace(riGanP, riGanFacts, wangShuaiAt(wangShuaiList, riGanP)),
		HourStemPalace:   newStemPalace(shiGanP, shiGanFacts, wangShuaiAt(wangShuaiList, shiGanP)),
		RiShiRelation:    analyzeRiShiRelation(riGanP, shiGanP),
		KongWangAffected: affectedVoidSymbols(p, riGanP, shiGanP, leadP),
		MaXingAffected:   affectedHorseSymbols(p, riGanP, shiGanP, leadP),
		GongWangShuai:    wangShuaiList,
		DutyStarPalace:   findStarPalace(p, p.DutyStar),
		DutyDoorPalace:   findDoorPalaceIdx(p, p.DutyDoor),
	}
	chart.Specialized = computeSpecializedContext(chart, bz)
	return chart
}

func buildChartMethod(ju juShu, bz ganzhi.Bazi) ChartMethod {
	leadZhu := leadPillar(bz, ju)
	xunIdx := ganzhi.SixtyCycleIndex(leadZhu.Gan, leadZhu.Zhi) / 10 * 10
	xunShou := ganzhi.SixtyToZhu(xunIdx)
	fu := fuTou(bz.Ri)
	duty := findDuty(leadZhu, placeDiPan(ju.Number, ju.YinDun))
	xunShouInfo := newPillarInfo(xunShou)
	xunShouInfo.LiuYi = ganzhi.GanName(findXunShou(leadZhu))
	scope := scopeMethods[ju.Method.Scope]
	school := schoolMethods[ju.Method.School]
	var tianQin *TianQinRule
	if ju.Method.School == SchoolZhuanPan {
		rule := tianQinRule
		tianQin = &rule
	}
	var quarter *ChartQuarter
	if ju.Quarter != nil {
		quarter = &ChartQuarter{
			Index:             ju.Quarter.Index,
			MinutesPerQuarter: ju.Quarter.Minutes,
			LeadPillarMode:    ju.Quarter.LeadPillarMode,
			DunSource:         ju.Quarter.DunSource,
			HourBoundary:      ju.Quarter.HourBoundary,
		}
	}
	return ChartMethod{
		Scope: scope.Scope, School: school.School,
		DingjuMethod: string(ju.Method.Dingju), QuarterRule: string(ju.Method.QuarterRule),
		BaseDingjuMethod: string(ju.Method.BaseDingju),
		MethodSource:     scope.MethodSource, DayBoundary: qimenDayBoundary,
		SpiritMode: school.SpiritMode, Geometry: school.Geometry,
		StarMode: school.StarMode, DoorMode: school.DoorMode,
		SpiritFlight: school.SpiritFlight,
		JieQi:        tianwen.DecodeJieQi(ju.JieQi), YongJuJieQi: ju.YongJuJieQi, Yuan: DecodeYuan(ju.Yuan),
		ZhiRunState:     DecodeZhiRunState(ju.ZhiRunState),
		TianQin:         tianQin,
		LeadPillar:      newPillarInfo(leadZhu),
		LeadXunShou:     xunShouInfo,
		LeadXunShouGong: duty.Palace,
		DayFuTou:        newPillarInfo(fu),
		Quarter:         quarter,
	}
}

func newPillarInfo(zhu ganzhi.Zhu) PillarInfo {
	return PillarInfo{
		Name: ganzhi.GanName(zhu.Gan) + ganzhi.ZhiName(zhu.Zhi),
		Gan:  ganzhi.GanName(zhu.Gan),
		Zhi:  ganzhi.ZhiName(zhu.Zhi),
	}
}

func affectedVoidSymbols(p Board, dayPalace, hourPalace, leadPalace GongIndex) []AffectedSymbol {
	result := []AffectedSymbol{}
	for _, void := range p.KongWang {
		for _, focus := range focusSymbols(dayPalace, hourPalace, leadPalace) {
			if void.Gong == focus.palace {
				result = append(result, AffectedSymbol{
					Symbol: focus.name, Role: focus.role, Branch: void.Branch, Gong: void.Gong,
				})
			}
		}
	}
	return result
}

func affectedHorseSymbols(p Board, dayPalace, hourPalace, leadPalace GongIndex) []AffectedSymbol {
	result := []AffectedSymbol{}
	for _, focus := range focusSymbols(dayPalace, hourPalace, leadPalace) {
		if p.MaXing.Gong == focus.palace {
			result = append(result, AffectedSymbol{
				Symbol: focus.name, Role: focus.role, Branch: p.MaXing.Branch, Gong: p.MaXing.Gong,
			})
		}
	}
	return result
}

type focusSymbol struct {
	name   string
	role   string
	palace GongIndex
}

func focusSymbols(dayPalace, hourPalace, leadPalace GongIndex) []focusSymbol {
	return []focusSymbol{
		{name: "日干", role: "pillar", palace: dayPalace},
		{name: "时干", role: "pillar", palace: hourPalace},
		{name: "主柱", role: "lead", palace: leadPalace},
	}
}

func palaceWangShuai(p Board, yueZhi ganzhi.Zhi) []GongWangShuai {
	result := []GongWangShuai{}
	for i := range p.GongWei {
		wuxing := luoshu.PalaceTable[i+1].Element
		wangShuai := ganzhi.WangShuaiOf(wuxing, yueZhi)
		result = append(result, GongWangShuai{
			Gong:      GongIndex(i + 1),
			Wuxing:    wuxing,
			WangShuai: wangShuai,
		})
	}
	return result
}

func wangShuaiAt(list []GongWangShuai, palace GongIndex) GongWangShuai {
	if palace < GongKan || palace > GongLi || int(palace) > len(list) {
		return GongWangShuai{}
	}
	return list[palace-1]
}

func newStemPalace(palace GongIndex, facts []GanGongFact, wangShuai GongWangShuai) StemPalace {
	return StemPalace{Palace: palace, Facts: facts, WangShuai: wangShuai}
}

func leadPillarPalace(p Board) GongIndex {
	return findGanPalaceIdx(p, resolveJiaDunGan(p.LeadGan, p.LeadZhi))
}

// findGanPalaceIdx 求干（求测人日干/时干/用神干）的落宫。
// 奇门领域规则：用神落宫以天盘为核心判断依据（天盘主求测人/所问事之当下状态，
// 地盘反映过去/潜藏）。故优先取天盘干所在宫；无天盘则取地盘干所在宫。
// 甲遁：甲不露，遁于六仪（由 resolveJiaDunGan 先转成六仪）。
func findGanPalaceIdx(p Board, g ganzhi.Gan) GongIndex {
	for i := range p.GongWei {
		for _, item := range p.GongWei[i].TianPan {
			if item.Gan == g {
				return GongIndex(i + 1)
			}
		}
	}
	for i := range p.GongWei {
		if p.GongWei[i].DiPanGan == g {
			return GongIndex(i + 1)
		}
	}
	return 0
}

func findGanGongFacts(p Board, g ganzhi.Gan) []GanGongFact {
	facts := []GanGongFact{}
	for i := range p.GongWei {
		for _, item := range p.GongWei[i].TianPan {
			if item.Gan == g {
				facts = append(facts, GanGongFact{Palace: GongIndex(i + 1), Layer: "heaven"})
			}
		}
	}
	for i := range p.GongWei {
		if p.GongWei[i].DiPanGan == g {
			facts = append(facts, GanGongFact{Palace: GongIndex(i + 1), Layer: "earth"})
		}
	}
	return facts
}

func findStarPalace(p Board, star StarIndex) GongIndex {
	for i, palace := range p.GongWei {
		for _, item := range palace.TianPan {
			if item.Star == star {
				return GongIndex(i + 1)
			}
		}
	}
	return 0
}

func analyzeRiShiRelation(dayPalace, hourPalace GongIndex) RiShiRelation {
	relation := wrSame
	if dayPalace > 0 && hourPalace > 0 {
		dayElement := palaceWuxing(dayPalace)
		hourElement := palaceWuxing(hourPalace)
		if dayElement == hourElement {
			relation = wrSame
		} else if ganzhi.Sheng(dayElement, hourElement) {
			relation = wrRiGeneratesShi
		} else if ganzhi.Sheng(hourElement, dayElement) {
			relation = wrShiGeneratesRi
		} else if ganzhi.Ke(dayElement, hourElement) {
			relation = wrRiControlsShi
		} else {
			relation = wrShiControlsRi
		}
	}
	return RiShiRelation{
		Zhu: "day_stem_palace", Ke: "hour_stem_palace",
		Relation: DecodeWuxingRelation(relation), Name: riShiRelations[relation],
	}
}

// findDoorPalaceIdx finds which gong a door resides in.
func findDoorPalaceIdx(p Board, d DoorIndex) GongIndex {
	for i, pg := range p.GongWei {
		if pg.Door == d {
			return GongIndex(i + 1)
		}
	}
	return 0
}
