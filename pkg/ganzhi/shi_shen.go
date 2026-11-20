package ganzhi

// ShiShenName returns the Chinese name for a ten god type.
func ShiShenName(tg ShiShen) string { return tg.String() }

// ShiShenFromGan returns the ShiShen for another gan relative to the day master.
func ShiShenFromGan(riYuan, other Gan) ShiShen {
	dmElem := GanWuxing(riYuan)
	otherElem := GanWuxing(other)
	dmYY := GanYinYang(riYuan)
	otherYY := GanYinYang(other)
	return ShiShenType(dmElem, dmYY, otherElem, otherYY)
}

// ShiShenType classifies the ten-god relationship between day master and another gan.
func ShiShenType(dmElem Wuxing, dmYY YinYang, otherElem Wuxing, otherYY YinYang) ShiShen {
	switch {
	case dmElem == otherElem:
		if dmYY == otherYY {
			return ShiShenBiJian
		}
		return ShiShenJieCai
	case Sheng(dmElem, otherElem):
		if dmYY == otherYY {
			return ShiShenShiShen
		}
		return ShiShenShangGuan
	case Sheng(otherElem, dmElem):
		if dmYY == otherYY {
			return ShiShenPianYin
		}
		return ShiShenZhengYin
	case Ke(dmElem, otherElem):
		if dmYY == otherYY {
			return ShiShenPianCai
		}
		return ShiShenZhengCai
	default:
		if dmYY == otherYY {
			return ShiShenQiSha
		}
		return ShiShenZhengGuan
	}
}
