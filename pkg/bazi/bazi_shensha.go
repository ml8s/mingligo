package bazi

import (
	"slices"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

// Shensha category constants.
const (
	catJi        = "吉"
	catXiong     = "凶"
	catZhongXing = "中性"
)

// ShenShaProjection describes a single shensha hit on a pillar.
type ShenShaProjection struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type ganZhiPair struct {
	gan ganzhi.Gan
	zhi ganzhi.Zhi
}

// Package-level lookup maps, populated by Load() from data/shensha.json.
var (
	taohuaZhiMap   map[ganzhi.Zhi]ganzhi.Zhi
	yimaZhiMap     map[ganzhi.Zhi]ganzhi.Zhi
	huagaiZhiMap   map[ganzhi.Zhi]ganzhi.Zhi
	yangRenLookup  map[ganzhi.Gan]ganzhi.Zhi
	jieshaZhi      map[ganzhi.Zhi]ganzhi.Zhi
	zaishaZhi      map[ganzhi.Zhi]ganzhi.Zhi
	hongluanLookup map[ganzhi.Zhi]ganzhi.Zhi
	tianxiLookup   map[ganzhi.Zhi]ganzhi.Zhi
)

var tianYiLookup map[ganzhi.Gan][]ganzhi.Zhi

var (
	tiandeTargets   map[ganzhi.Zhi][]tianDeTarget
	yuedeGan        map[ganzhi.Zhi]ganzhi.Gan
	jiangxingLookup map[ganzhi.Zhi]ganzhi.Zhi
	jinyuLookup     map[ganzhi.Gan][]ganzhi.Zhi
	taiJiLookup     map[ganzhi.Gan][]ganzhi.Zhi
	tianChuLookup   map[ganzhi.Gan][]ganzhi.Zhi
	fuXingLookup    map[ganzhi.Gan][]ganzhi.Zhi
	guoYinLookup    map[ganzhi.Gan][]ganzhi.Zhi
	yueEnGan        map[ganzhi.Zhi][]ganzhi.Gan
	xueRenLookup    map[ganzhi.Gan]ganzhi.Zhi
	feiRenLookup    map[ganzhi.Gan]ganzhi.Zhi
	wangShenZhi     map[ganzhi.Zhi]ganzhi.Zhi
	deXiuByMonth    map[ganzhi.Zhi]deXiuStems
	yinChaYangCuo   map[int]struct{}
	tongZiSeason    map[int][]ganzhi.Zhi
	tongZiNayin     map[ganzhi.Wuxing][]ganzhi.Zhi
	tianLuoDiWang   map[ganzhi.Zhi]string
	shiEDaBai       map[int]struct{}
)

type deXiuStems struct {
	De  []ganzhi.Gan
	Xiu []ganzhi.Gan
}

// tianDeTarget 天德贵人的匹配目标：天德可为天干型（如正月见丁）或地支型（如二月见申）。
type tianDeTarget struct {
	IsZhi bool
	Gan   ganzhi.Gan
	Zhi   ganzhi.Zhi
}

// mergeShenSha purely appends per-pillar shensha hits. Deduplication stays
// inside each producer (its own marked/seen set); no cross-producer dedup exists.
func mergeShenSha(dst, hits [4][]ShenShaProjection) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	for i := range out {
		out[i] = append(append([]ShenShaProjection{}, dst[i]...), hits[i]...)
	}
	return out
}

// ShenShaProjections computes all shensha for the bazi chart, grouped by pillar.
func ShenShaProjections(bz ganzhi.Bazi, gender ganzhi.Gender) [4][]ShenShaProjection {
	riYuan := bz.Ri.Gan
	yueZhi := bz.Yue.Zhi
	zhus := bz.Slice()
	var out [4][]ShenShaProjection
	zhi := [4]ganzhi.Zhi{zhus[0].Zhi, zhus[1].Zhi, zhus[2].Zhi, zhus[3].Zhi}
	seasonIdx := sanHuiSeasonIndex(yueZhi)
	nianZhi := zhus[0].Zhi
	// yearSanHuiIdx（年支三会组）供孤辰寡宿使用：寅卯辰→0 巳午未→1 申酉戌→2 亥子丑→3。
	// 孤辰寡宿以年支为准（亥子丑人见寅为孤、见戌为寡），与月支三会组不同。
	yearSanHuiIdx := ((int(nianZhi) - 3 + 12) % 12) / 3

	out = mergeShenSha(out, addTianYi(bz, riYuan, zhus[0].Gan))
	out = mergeShenSha(out, addTaiJi(bz, riYuan, zhus[0].Gan))
	out = mergeShenSha(out, addDeXiu(bz, yueZhi))
	out = mergeShenSha(out, addTianChu(bz, riYuan, zhus[0].Gan))
	out = mergeShenSha(out, addFuXing(bz, riYuan))
	out = mergeShenSha(out, addGuoYin(bz, riYuan))
	out = mergeShenSha(out, addWenChang(bz, riYuan))
	out = mergeShenSha(out, addXueTang(bz))
	out = mergeShenSha(out, addLuShen(bz, riYuan))
	out = mergeShenSha(out, addYangRen(bz, riYuan))
	out = mergeShenSha(out, addFeiRen(bz, riYuan))
	out = mergeShenSha(out, addWangShen(bz, nianZhi))
	out = mergeShenSha(out, addYinChaYangCuo(bz))
	out = mergeShenSha(out, addTongZi(bz, seasonIdx))
	out = mergeShenSha(out, addTianDeHe(bz, yueZhi))
	out = mergeShenSha(out, addTianDe(bz, yueZhi))
	out = mergeShenSha(out, addYueDe(bz, yueZhi))
	out = mergeShenSha(out, addTaoHua(bz, zhi))
	out = mergeShenSha(out, addYiMa(bz, zhi))
	out = mergeShenSha(out, addHuaGai(bz, zhi))
	out = mergeShenSha(out, addJiangXing(bz, zhi))
	out = mergeShenSha(out, addJieSha(bz, zhi))
	out = mergeShenSha(out, addZaiSha(bz, zhi))
	out = mergeShenSha(out, addGuChenGuaSu(bz, yearSanHuiIdx))
	out = mergeShenSha(out, addHongLuanTianXi(bz, nianZhi))
	out = mergeShenSha(out, addJinYu(bz, riYuan))
	out = mergeShenSha(out, addCiGuan(bz))
	out = mergeShenSha(out, addYueEn(bz, yueZhi))
	out = mergeShenSha(out, addTianShe(bz, yueZhi))
	out = mergeShenSha(out, addTianLuoDiWang(bz))
	out = mergeShenSha(out, addGouJiao(bz, nianZhi, zhus[0].Gan, gender))
	out = mergeShenSha(out, addYuanChen(bz, nianZhi, zhus[0].Gan, gender))
	out = mergeShenSha(out, addXueRen(bz, riYuan))
	out = mergeShenSha(out, addSiFei(bz, seasonIdx))
	out = mergeShenSha(out, addShiEDaBai(bz))

	return out
}

// addTianYi 天乙贵人：按日干与年干双查（兼容两种流派），同一柱只标注一次（去重）。
func addTianYi(bz ganzhi.Bazi, riYuan, nianGan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	marked := map[int]bool{}
	mark := func(gan ganzhi.Gan) {
		targets, ok := tianYiLookup[gan]
		if !ok {
			return
		}
		zhus := bz.Slice()
		for pi, p := range zhus {
			for _, t := range targets {
				if p.Zhi == t && !marked[pi] {
					marked[pi] = true
					out[pi] = append(out[pi], ShenShaProjection{
						Name: "天乙贵人", Category: catJi, Description: "主贵人相助，逢凶化吉",
					})
				}
			}
		}
	}
	mark(riYuan)
	mark(nianGan)
	return out
}

var wenChangLookup map[ganzhi.Gan][]ganzhi.Zhi

func addTaiJi(bz ganzhi.Bazi, riYuan, nianGan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	marked := map[int]bool{}
	mark := func(gan ganzhi.Gan) {
		for pi, p := range bz.Slice() {
			for _, target := range taiJiLookup[gan] {
				if p.Zhi == target && !marked[pi] {
					marked[pi] = true
					out[pi] = append(out[pi], ShenShaProjection{
						Name: "太极贵人", Category: catJi, Description: "主聪慧好学，宜研哲理玄学",
					})
				}
			}
		}
	}
	mark(riYuan)
	mark(nianGan)
	return out
}

func addDeXiu(bz ganzhi.Bazi, yueZhi ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	group, ok := deXiuByMonth[yueZhi]
	if !ok {
		return out
	}
	for pi, p := range bz.Slice() {
		if containsGan(group.De, p.Gan) || containsGan(group.Xiu, p.Gan) {
			out = mergeShenSha(out, addDeXiuEntry(pi))
		}
	}
	return out
}

func containsGan(values []ganzhi.Gan, want ganzhi.Gan) bool {
	return slices.Contains(values, want)
}

func addDeXiuEntry(pillar int) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out[pillar] = append(out[pillar], ShenShaProjection{
		Name: "德秀贵人", Category: catJi, Description: "主禀气清粹，逢凶化吉",
	})
	return out
}

func addTianChu(bz ganzhi.Bazi, riYuan, nianGan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, appendShenShaByGanLookup(bz, riYuan, tianChuLookup, "天厨贵人", catJi, "主食禄福泽"))
	out = mergeShenSha(out, appendShenShaByGanLookup(bz, nianGan, tianChuLookup, "天厨贵人", catJi, "主食禄福泽"))
	return out
}

func addFuXing(bz ganzhi.Bazi, riYuan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, appendShenShaByGanLookup(bz, riYuan, fuXingLookup, "福星贵人", catJi, "主一生福气，平安顺遂"))
	return out
}

func addGuoYin(bz ganzhi.Bazi, riYuan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, appendShenShaByGanLookup(bz, riYuan, guoYinLookup, "国印", catJi, "主持印掌信，宜守正职守"))
	return out
}

func addWenChang(bz ganzhi.Bazi, riYuan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, appendShenShaByGanLookup(bz, riYuan, wenChangLookup, "文昌", catJi, "主学业、文书、才华"))
	return out
}

// addXueTang 采用《三命通会》学堂正位：年命纳音五行的长生支上，
// 再见同一纳音五行的干支，方为学堂。
func addXueTang(bz ganzhi.Bazi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, addNayinStageShenSha(bz, "学堂", "主学业聪颖", true))
	return out
}

func addLuShen(bz ganzhi.Bazi, riYuan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, addChangShengShenSha(bz, riYuan, 3, "禄神", catJi, "日主临官之位，主福禄安康"))
	return out
}

func addChangShengShenSha(bz ganzhi.Bazi, riYuan ganzhi.Gan, stageIdx int, name, cat, desc string) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	stageRow := ganzhi.ChangShengTable[riYuan]
	if len(stageRow) != 12 {
		return out
	}
	for pi, p := range zhus {
		if bn := p.Zhi; bn >= 1 && bn <= 12 && stageRow[stageIdx] == bn {
			out[pi] = append(out[pi], ShenShaProjection{Name: name, Category: cat, Description: desc})
		}
	}
	return out
}

// addCiGuan 采用《三命通会》词馆正位：年命纳音五行的临官支上，
// 再见同一纳音五行的干支，方为词馆。
func addCiGuan(bz ganzhi.Bazi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, addNayinStageShenSha(bz, "词馆", "主文章、口才、文职", false))
	return out
}

func addNayinStageShenSha(bz ganzhi.Bazi, name, desc string, changSheng bool) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	yearElement := ganzhi.NayinWuxingOf(zhus[0].Gan, zhus[0].Zhi)
	if yearElement == 0 {
		return out
	}
	targets := nayinStageBranches(changSheng)
	for pi, p := range zhus {
		if ganzhi.NayinWuxingOf(p.Gan, p.Zhi) != yearElement {
			continue
		}
		if p.Zhi == targets[yearElement] {
			out[pi] = append(out[pi], ShenShaProjection{Name: name, Category: catJi, Description: desc})
		}
	}
	return out
}

func nayinStageBranches(changSheng bool) map[ganzhi.Wuxing]ganzhi.Zhi {
	if changSheng {
		return map[ganzhi.Wuxing]ganzhi.Zhi{
			ganzhi.WxMu:   ganzhi.ZhiHai,
			ganzhi.WxHuo:  ganzhi.ZhiYin,
			ganzhi.WxTu:   ganzhi.ZhiShen,
			ganzhi.WxJin:  ganzhi.ZhiSi,
			ganzhi.WxShui: ganzhi.ZhiShen,
		}
	}
	return map[ganzhi.Wuxing]ganzhi.Zhi{
		ganzhi.WxMu:   ganzhi.ZhiYin,
		ganzhi.WxHuo:  ganzhi.ZhiSi,
		ganzhi.WxTu:   ganzhi.ZhiHai,
		ganzhi.WxJin:  ganzhi.ZhiShen,
		ganzhi.WxShui: ganzhi.ZhiHai,
	}
}

func addYangRen(bz ganzhi.Bazi, riYuan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	for pi, p := range zhus {
		if yangRenLookup[riYuan] == p.Zhi {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "羊刃", Category: catXiong, Description: "日干帝旺/刃位，主刚强果断，但易冲动",
			})
		}
	}
	return out
}

func addFeiRen(bz ganzhi.Bazi, riYuan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	for pi, p := range zhus {
		if feiRenLookup[riYuan] == p.Zhi {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "飞刃", Category: catXiong, Description: "羊刃冲位，主动荡急变",
			})
		}
	}
	return out
}

func addWangShen(bz ganzhi.Bazi, nianZhi ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	target, ok := wangShenZhi[nianZhi]
	if !ok {
		return out
	}
	for pi, p := range bz.Slice() {
		if p.Zhi == target {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "亡神", Category: catXiong, Description: "主心机深虑，事多暗耗",
			})
		}
	}
	return out
}

func addYinChaYangCuo(bz ganzhi.Bazi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	if _, ok := yinChaYangCuo[ganzhi.SixtyCycleIndex(bz.Ri.Gan, bz.Ri.Zhi)]; ok {
		out[2] = append(out[2], ShenShaProjection{
			Name: "阴差阳错", Category: catXiong, Description: "主婚缘人事阴差阳错，多生错过",
		})
	}
	return out
}

func addTongZi(bz ganzhi.Bazi, seasonIdx int) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	checkBranches := func(targets []ganzhi.Zhi) bool {
		if len(targets) == 0 {
			return false
		}
		for _, p := range zhus[2:] {
			for _, target := range targets {
				if p.Zhi == target {
					return true
				}
			}
		}
		return false
	}
	if checkBranches(tongZiSeason[seasonIdx]) {
		out = mergeShenSha(out, addTongZiEntry(2))
		return out
	}
	yearElement := ganzhi.NayinWuxingOf(zhus[0].Gan, zhus[0].Zhi)
	if checkBranches(tongZiNayin[yearElement]) {
		out = mergeShenSha(out, addTongZiEntry(2))
	}
	return out
}

func addTongZiEntry(pillar int) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out[pillar] = append(out[pillar], ShenShaProjection{
		Name: "童子煞", Category: catZhongXing, Description: "传统取象主幼年多病、姻缘迟滞",
	})
	return out
}

func addTianDeHe(bz ganzhi.Bazi, yueZhi ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	targets, ok := tiandeTargets[yueZhi]
	if !ok {
		return out
	}
	zhus := bz.Slice()
	for _, target := range targets {
		for pi, p := range zhus {
			hit := false
			if target.IsZhi {
				hit = ganzhi.IsZhiHe(target.Zhi, p.Zhi)
			} else {
				hit = ganzhi.IsGanHe(target.Gan, p.Gan)
			}
			if hit {
				out[pi] = append(out[pi], ShenShaProjection{
					Name: "天德合", Category: catJi, Description: "天德所合，主贵人助力",
				})
			}
		}
	}
	return out
}

func addTianDe(bz ganzhi.Bazi, yueZhi ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	targets, ok := tiandeTargets[yueZhi]
	if !ok {
		return out
	}
	for _, tgt := range targets {
		for pi, p := range zhus {
			hit := false
			if tgt.IsZhi {
				hit = p.Zhi == tgt.Zhi
			} else {
				hit = p.Gan == tgt.Gan
			}
			if hit {
				out[pi] = append(out[pi], ShenShaProjection{
					Name: "天德", Category: catJi, Description: "天德贵人，主福泽深厚，化险为夷",
				})
			}
		}
	}
	return out
}

func addYueDe(bz ganzhi.Bazi, yueZhi ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	targetStem, ok := yuedeGan[yueZhi]
	if !ok {
		return out
	}
	for pi, p := range zhus {
		if p.Gan == targetStem {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "月德", Category: catJi, Description: "月德贵人，主月令之德，人缘佳",
			})
		}
	}
	return out
}

// addTriadShenSha 三合类神煞（驿马/桃花/华盖/劫煞/灾煞/将星）：
// 以年支与日支为参考（兼容两种流派），同一柱同一神煞只标注一次（去重）。
func addTriadShenSha(bz ganzhi.Bazi, zhi [4]ganzhi.Zhi, lookup map[ganzhi.Zhi]ganzhi.Zhi, name, cat, desc string) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	marked := map[int]bool{}
	for _, refIdx := range []int{0, 2} { // year & day zhi
		if tb, ok := lookup[zhi[refIdx]]; ok {
			for pi, p := range zhus {
				if p.Zhi == tb && !marked[pi] {
					marked[pi] = true
					out[pi] = append(out[pi], ShenShaProjection{Name: name, Category: cat, Description: desc})
				}
			}
		}
	}
	return out
}

func addTaoHua(bz ganzhi.Bazi, zhi [4]ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, addTriadShenSha(bz, zhi, taohuaZhiMap, "桃花", catZhongXing, "主异性缘佳，浪漫多情"))
	return out
}

func addYiMa(bz ganzhi.Bazi, zhi [4]ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, addTriadShenSha(bz, zhi, yimaZhiMap, "驿马", catZhongXing, "主动荡、奔波、迁移"))
	return out
}

func addHuaGai(bz ganzhi.Bazi, zhi [4]ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, addTriadShenSha(bz, zhi, huagaiZhiMap, "华盖", catZhongXing, "主孤独清高，聪明好学，有艺术天赋"))
	return out
}

func addJiangXing(bz ganzhi.Bazi, zhi [4]ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, addTriadShenSha(bz, zhi, jiangxingLookup, "将星", catJi, "主领导才能，有权威"))
	return out
}

func addJieSha(bz ganzhi.Bazi, zhi [4]ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, addTriadShenSha(bz, zhi, jieshaZhi, "劫煞", catXiong, "主破财、意外、是非"))
	return out
}

func addZaiSha(bz ganzhi.Bazi, zhi [4]ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, addTriadShenSha(bz, zhi, zaishaZhi, "灾煞", catXiong, "主灾祸、疾病、横事"))
	return out
}

func addGuChenGuaSu(bz ganzhi.Bazi, seasonIdx int) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	guchenBranches := [4]ganzhi.Zhi{6, 9, 12, 3}
	guasuBranches := [4]ganzhi.Zhi{2, 5, 8, 11}
	for pi, p := range zhus {
		if p.Zhi == guchenBranches[seasonIdx] {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "孤辰", Category: catXiong, Description: "主性格孤僻，晚婚或婚姻不顺",
			})
		}
		if p.Zhi == guasuBranches[seasonIdx] {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "寡宿", Category: catXiong, Description: "主孤独寂寞，夫妻缘薄",
			})
		}
	}
	return out
}

func addHongLuanTianXi(bz ganzhi.Bazi, nianZhi ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	if target, ok := hongluanLookup[nianZhi]; ok {
		for pi, p := range zhus {
			if p.Zhi == target {
				out[pi] = append(out[pi], ShenShaProjection{
					Name: "红鸾", Category: catJi, Description: "主婚喜、恋爱、添丁",
				})
			}
		}
	}
	if target, ok := tianxiLookup[nianZhi]; ok {
		for pi, p := range zhus {
			if p.Zhi == target {
				out[pi] = append(out[pi], ShenShaProjection{
					Name: "天喜", Category: catJi, Description: "主喜庆之事，婚恋吉兆",
				})
			}
		}
	}
	return out
}

func addJinYu(bz ganzhi.Bazi, riYuan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	out = mergeShenSha(out, appendShenShaByGanLookup(bz, riYuan, jinyuLookup, "金舆", catJi, "主财运、车辆、出行顺利"))
	return out
}

func addYueEn(bz ganzhi.Bazi, yueZhi ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	targets, ok := yueEnGan[yueZhi]
	if !ok {
		return out
	}
	for _, ts := range targets {
		for pi, p := range zhus {
			if p.Gan == ts {
				out[pi] = append(out[pi], ShenShaProjection{
					Name: "月恩", Category: catJi, Description: "月令之恩，主福佑加持",
				})
			}
		}
	}
	return out
}

func addTianShe(bz ganzhi.Bazi, yueZhi ganzhi.Zhi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	season := sanHuiSeasonIndex(yueZhi)
	tianSheChecks := [4]ganZhiPair{{5, 3}, {1, 7}, {5, 9}, {1, 1}} // 戊寅, 甲午, 戊申, 甲子
	if season >= 0 && season < 4 {
		pair := tianSheChecks[season]
		if zhus[2].Gan == pair.gan && zhus[2].Zhi == pair.zhi {
			out[2] = append(out[2], ShenShaProjection{
				Name: "天赦", Category: catJi, Description: "天赦日出生，主逢凶化吉，宽恕赦免",
			})
		}
	}
	return out
}

func addTianLuoDiWang(bz ganzhi.Bazi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	yearElement := ganzhi.NayinWuxingOf(zhus[0].Gan, zhus[0].Zhi)
	for pi, p := range zhus {
		label, ok := tianLuoDiWang[p.Zhi]
		if ok && ((label == "天罗" && yearElement == ganzhi.WxHuo) ||
			(label == "地网" && (yearElement == ganzhi.WxShui || yearElement == ganzhi.WxTu))) {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: label, Category: catXiong, Description: "主运势阻滞，有志难伸",
			})
		}
	}
	return out
}

// sanHuiSeasonIndex returns 寅卯辰→0、巳午未→1、申酉戌→2、亥子丑→3.
func sanHuiSeasonIndex(zhi ganzhi.Zhi) int {
	return ((int(zhi) - int(ganzhi.ZhiYin) + 12) % 12) / 3
}

func addGouJiao(bz ganzhi.Bazi, nianZhi ganzhi.Zhi, nianGan ganzhi.Gan, gender ganzhi.Gender) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	ahead := zhiOffset(nianZhi, 3)
	behind := zhiOffset(nianZhi, -3)
	gouShen, jiaoShen := ahead, behind
	if !yangMaleOrYinFemale(nianGan, gender) {
		gouShen, jiaoShen = behind, ahead
	}
	for pi, p := range zhus {
		if p.Zhi == gouShen {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "勾神", Category: catXiong, Description: "主纠缠牵连，是非官讼",
			})
		}
		if p.Zhi == jiaoShen {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "绞神", Category: catXiong, Description: "主受困被缚，身不由己",
			})
		}
	}
	return out
}

func addYuanChen(bz ganzhi.Bazi, nianZhi ganzhi.Zhi, nianGan ganzhi.Gan, gender ganzhi.Gender) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	ycBranch := yuanChenZhi(nianZhi, nianGan, gender)
	for pi, p := range zhus {
		if p.Zhi == ycBranch {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "元辰", Category: catXiong, Description: "主波折反复，好事多磨",
			})
		}
	}
	return out
}

func addXueRen(bz ganzhi.Bazi, riYuan ganzhi.Gan) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	for pi, p := range zhus {
		if xueRenLookup[riYuan] == p.Zhi {
			out[pi] = append(out[pi], ShenShaProjection{
				Name: "血刃", Category: catXiong, Description: "主意外血光，手术外伤",
			})
		}
	}
	return out
}

func addSiFei(bz ganzhi.Bazi, seasonIdx int) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	siFeiZhus := [4][]ganZhiPair{
		{{7, 9}, {8, 10}}, {{9, 1}, {10, 12}}, {{1, 3}, {2, 4}}, {{3, 7}, {4, 6}},
	}
	if seasonIdx < 0 || seasonIdx >= 4 {
		return out
	}
	for pi, p := range zhus {
		for _, pair := range siFeiZhus[seasonIdx] {
			if p.Gan == pair.gan && p.Zhi == pair.zhi {
				out[pi] = append(out[pi], ShenShaProjection{
					Name: "四废", Category: catXiong, Description: "四季废日，主事业阻滞，有志难伸",
				})
			}
		}
	}
	return out
}

func addShiEDaBai(bz ganzhi.Bazi) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	if _, ok := shiEDaBai[ganzhi.SixtyCycleIndex(zhus[2].Gan, zhus[2].Zhi)]; ok {
		out[2] = append(out[2], ShenShaProjection{
			Name: "十恶大败", Category: catXiong, Description: "日柱十恶大败日，主财库不聚，须谨慎理财",
		})
	}
	return out
}

func appendShenShaByGanLookup(bz ganzhi.Bazi, s ganzhi.Gan, lookup map[ganzhi.Gan][]ganzhi.Zhi, name, cat, desc string) [4][]ShenShaProjection {
	var out [4][]ShenShaProjection
	zhus := bz.Slice()
	targets, ok := lookup[s]
	if !ok {
		return out
	}
	for pi, p := range zhus {
		for _, t := range targets {
			if p.Zhi == t {
				out[pi] = append(out[pi], ShenShaProjection{Name: name, Category: cat, Description: desc})
			}
		}
	}
	return out
}

// VoidPalaces returns pillar indices whose zhi fall in the void (空亡)
// of the day pillar's 旬.
func VoidPalaces(bz ganzhi.Bazi) []int {
	sbIdx := ganzhi.SixtyCycleIndex(bz.Ri.Gan, bz.Ri.Zhi)
	xunIdx := sbIdx / 10

	voidPairs := [6][2]ganzhi.Zhi{
		{11, 12}, {9, 10}, {7, 8}, {5, 6}, {3, 4}, {1, 2},
	}
	v1, v2 := voidPairs[xunIdx][0], voidPairs[xunIdx][1]

	var hits []int
	for pi, p := range bz.Slice() {
		b := p.Zhi
		if b == v1 || b == v2 {
			hits = append(hits, pi)
		}
	}
	return hits
}

// computeDynamicShenSha computes shensha triggered by an external zhi against the bazi chart.
func computeDynamicShenSha(b ganzhi.Zhi, nianZhi, riBranch ganzhi.Zhi, riYuan ganzhi.Gan) []ShenShaProjection {
	var result []ShenShaProjection
	seen := map[string]bool{}
	add := func(name, cat, desc string) {
		if !seen[name] {
			seen[name] = true
			result = append(result, ShenShaProjection{Name: name, Category: cat, Description: desc})
		}
	}
	// 三合局系神煞（桃花/驿马/华盖/劫煞/灾煞）——年支+日支双查（《三命通会》年支/日支桃花驿马）
	for _, rb := range []ganzhi.Zhi{nianZhi, riBranch} {
		if tb, ok := taohuaZhiMap[rb]; ok && tb == b {
			add("桃花", catZhongXing, "流运桃花，异性缘佳")
		}
		if tb, ok := yimaZhiMap[rb]; ok && tb == b {
			add("驿马", catZhongXing, "流运驿马，动象奔波")
		}
		if tb, ok := huagaiZhiMap[rb]; ok && tb == b {
			add("华盖", catZhongXing, "流运华盖，宜静思")
		}
		if js, ok := jieshaZhi[rb]; ok && js == b {
			add("劫煞", catXiong, "流运劫煞，防破财是非")
		}
		if zs, ok := zaishaZhi[rb]; ok && zs == b {
			add("灾煞", catXiong, "流运灾煞，防意外灾祸")
		}
	}
	// 红鸾/天喜——仅年支查（红鸾属年支体系，日支不取）
	if hl, ok := hongluanLookup[nianZhi]; ok && hl == b {
		add("红鸾", catJi, "流运红鸾，主婚喜添丁")
	}
	if tx, ok := tianxiLookup[nianZhi]; ok && tx == b {
		add("天喜", catJi, "流运天喜，喜庆之事")
	}
	// 天乙贵人/羊刃——按日干
	if targets, ok := tianYiLookup[riYuan]; ok {
		for _, t := range targets {
			if t == b {
				add("天乙贵人", catJi, "流运天乙贵人，有贵人相助")
				break
			}
		}
	}
	if yr, ok := yangRenLookup[riYuan]; ok && yr == b {
		add("羊刃", catXiong, "流运羊刃，防冲动冲突")
	}

	if result == nil {
		return []ShenShaProjection{}
	}
	return result
}

// computeAnnualShenSha 值年神煞（《协纪辨方书》——按太岁/流年支查表，命局四柱逢煞支即应）。
// 病符=太岁后1辰、丧门=前2辰、吊客=后2辰、大耗（岁破）=对冲。白虎查表有版本争议，不做。
func computeAnnualShenSha(b ganzhi.Zhi, bz ganzhi.Bazi) []ShenShaProjection {
	annual := []struct {
		name string
		zhi  ganzhi.Zhi
		desc string
	}{
		{"病符", zhiOffset(b, -1), "流运病符临命，主病灾"},
		{"丧门", zhiOffset(b, +2), "流运丧门临命，主孝服/丧事"},
		{"吊客", zhiOffset(b, -2), "流运吊客临命，主吊丧/孝服"},
		{"大耗", zhiOffset(b, +6), "流运大耗临命，主破财大耗"},
	}
	var result []ShenShaProjection
	seen := map[string]bool{}
	for _, a := range annual {
		for _, p := range bz.Slice() {
			if p.Zhi == a.zhi && !seen[a.name] {
				seen[a.name] = true
				result = append(result, ShenShaProjection{Name: a.name, Category: catXiong, Description: a.desc})
				break
			}
		}
	}
	if result == nil {
		return []ShenShaProjection{}
	}
	return result
}

// zhiOffset 地支循环偏移（z 后移 n 位，n 可为负）。
func zhiOffset(z ganzhi.Zhi, n int) ganzhi.Zhi {
	return ganzhi.Zhi((int(z)-1+n+120)%12 + 1)
}

func yuanChenZhi(nianZhi ganzhi.Zhi, nianGan ganzhi.Gan, gender ganzhi.Gender) ganzhi.Zhi {
	for _, p := range ganzhi.ChongPairs {
		if p.A == nianZhi {
			return yuanChenOffset(p.B, nianGan, gender)
		}
		if p.B == nianZhi {
			return yuanChenOffset(p.A, nianGan, gender)
		}
	}
	return 0
}

// yuanChenOffset applies the gender-aware offset to the clash zhi.
// 阳男阴女取冲前一位；阴男阳女取冲后一位（《三命通会·论元辰》）。
func yuanChenOffset(chongZhi ganzhi.Zhi, nianGan ganzhi.Gan, gender ganzhi.Gender) ganzhi.Zhi {
	if yangMaleOrYinFemale(nianGan, gender) {
		return chongZhi%12 + 1
	}
	return (chongZhi-2+12)%12 + 1
}

func yangMaleOrYinFemale(nianGan ganzhi.Gan, gender ganzhi.Gender) bool {
	yangYear := ganzhi.GanYinYang(nianGan) == ganzhi.Yang
	return yangYear == (gender == ganzhi.Male)
}
