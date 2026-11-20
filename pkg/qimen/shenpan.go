package qimen

// placeShenPan arranges the spirit plate around the outer ring.
// The leading spirit follows the duty star; Yang Dun proceeds clockwise and
// Yin Dun counter-clockwise.
func placeShenPan(yinDun bool, tianStarPos GongIndex) [9]SpiritIndex {
	var spirits [9]SpiritIndex
	start := outerRingIndex(tianStarPos)
	if start < 0 {
		return spirits
	}
	for step, si := 0, 0; si < 8; step++ {
		ringIdx := (start + step) % 8
		if yinDun {
			ringIdx = (start - step + 8) % 8
		}
		spirits[int(outerRing[ringIdx])-1] = spiritOrder[si]
		si++
	}
	return spirits
}

func spiritDisplayName(spirit SpiritIndex, yinDun bool, school School) string {
	if school == SchoolMingFaFeiPan && int(spirit) < len(mingfaSpiritNames) &&
		mingfaSpiritNames[int(spirit)] != "" {
		return mingfaSpiritNames[int(spirit)]
	}
	if school == SchoolLuoShuFeiPan && int(spirit) < len(flySpiritNames) && flySpiritNames[spirit] != "" {
		return flySpiritNames[spirit]
	}
	if yinDun {
		return spirit.YinName()
	}
	return spirit.YangName()
}

// spiritIsYinXuanWu reports whether the slot is the 玄武 (朱雀 under 阴遁)
// in the given school's naming convention. Only schools that keep the default
// 阴遁 naming map 朱雀→玄武 qualify; schools with their own spirit-name table
// override the slot name and are excluded.
func spiritIsYinXuanWu(spirit SpiritIndex, yinDun bool, school School) bool {
	if !yinDun || spirit != SpiritZhuQue {
		return false
	}
	if school == SchoolMingFaFeiPan && spiritSlotOverridden(mingfaSpiritNames, spirit) {
		return false
	}
	if school == SchoolLuoShuFeiPan && spiritSlotOverridden(flySpiritNames, spirit) {
		return false
	}
	return true
}

func spiritSlotOverridden(names [10]string, spirit SpiritIndex) bool {
	return int(spirit) >= 0 && int(spirit) < len(names) && names[spirit] != ""
}
