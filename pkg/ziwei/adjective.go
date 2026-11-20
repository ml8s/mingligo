package ziwei

type adjStarResult map[ShenSha]int

func computeAdjectiveStars(nianZhi, shiZhi Zhi, lunarMonth, lunarDay int, nianGan Gan) adjStarResult {
	r := make(adjStarResult)
	add := func(name ShenSha, zhiIdx int) { r[name] = zhiIdx }

	// A 年支系
	hl := (4 - int(nianZhi) + 12) % 12
	add(ssHongLuan, hl)
	add(ssTianXi, (hl+6)%12)
	for name, table := range nianStars {
		add(name, table[nianZhi-1])
	}

	// B 月支系
	for name, table := range yueStars {
		add(name, table[lunarMonth-1])
	}

	// C 时支系
	for name, table := range shiStars {
		add(name, table[shiZhi-1])
	}

	// D 日系
	zf := (lunarMonth + 3) % 12
	yf := (11 - lunarMonth + 12) % 12
	wc := (11 - int(shiZhi) + 12) % 12
	wq := (int(shiZhi) + 3) % 12
	add(ssSanTai, (zf+lunarDay-1)%12)
	add(ssBaZuo, ((yf-lunarDay+1)%12+12)%12)
	add(ssEnGuang, (wc+lunarDay-2+12)%12)
	add(ssTianGui, (wq+lunarDay-2+12)%12)

	// E 年干系
	for name, table := range ganStars {
		add(name, table[nianGan-1])
	}

	// soul/body for PALACE conversion
	soulIdx := (lunarMonth - 1 - earthlyIdxTable[shiZhi-1] + 12) % 12
	bodyIdx := (lunarMonth - 1 + earthlyIdxTable[shiZhi-1] + 12) % 12

	// F 复合系 — 年系旬空 (iztro公式)
	add(ssXunKong, xunKongNianXi(nianZhi, nianGan))
	add(ssNianJie, nianJiePos(nianZhi))

	add(ssTianCai, anXingIdxToZhiIdx((soulIdx+earthlyIdxTable[nianZhi-1])%12))
	add(ssTianShou, anXingIdxToZhiIdx((bodyIdx+earthlyIdxTable[nianZhi-1])%12))
	add(ssTianShang, anXingIdxToZhiIdx((5+soulIdx)%12))
	add(ssTianShi, anXingIdxToZhiIdx((7+soulIdx)%12))

	return r
}

// anXingIdxToZhiIdx 定义见 coord.go（统一坐标转换）

func nianJiePos(nianZhi Zhi) int {
	e := earthlyIdxTable[nianZhi-1]
	return (10 - e + 12) % 12
}

// xunKongNianXi computes 旬空 using year-gan-zhi formula (iztro年系).
func xunKongNianXi(nianZhi Zhi, nianGan Gan) int {
	yinIdx := (int(nianZhi) - 3 + 12) % 12 // EARTHLY_BRANCHES→寅0系
	xk := yinIdx + 9 - (int(nianGan) - 1) + 1
	xk = (xk + 12) % 12
	if (int(nianZhi)-1)%2 != xk%2 { // 年支阴阳 vs 索引阴阳
		xk = (xk + 1) % 12
	}
	return (2 + xk) % 12 // anXingIdxToZhiIdx
}
