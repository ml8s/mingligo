package ziwei

// soulStar returns the 命主 (soul star) name for a given ming gong zhi.
func soulStar(mingZhi Zhi) StarIndex {
	switch mingZhi {
	case 1:
		return TanLang
	case 2:
		return JuMen
	case 3:
		return LuCun
	case 4:
		return WenQu
	case 5:
		return LianZhen
	case 6:
		return WuQu
	case 7:
		return PoJun
	case 8:
		return WuQu
	case 9:
		return LianZhen
	case 10:
		return WenQu
	case 11:
		return LuCun
	case 12:
		return JuMen
	}
	return -1
}

// bodyStar returns the 身主 (body star) name for a given birth year zhi.
func bodyStar(nianZhi Zhi) StarIndex {
	switch nianZhi {
	case 1, 7:
		return HuoXing
	case 2, 8:
		return TianXiang
	case 3, 9:
		return TianLiang
	case 4, 10:
		return TianTong
	case 5, 11:
		return WenChang
	case 6, 12:
		return TianJi
	}
	return -1
}

// yuanGongPalace returns the gong index of 来因宫 (original gong).
// iztro rule: the gong whose heavenly gan equals the birth year's heavenly gan.
func yuanGongPalace(palaces [12]Gong, nianGan Gan) GongIndex {
	for i, p := range palaces {
		if p.Gan == nianGan && p.Zhi != 1 && p.Zhi != 2 { // 排除子丑
			return GongIndex(i)
		}
	}
	return 0
}
