package ziwei

import "github.com/ml8s/mingligo/pkg/ganzhi"

func daXianQiSui(ju JuShu) int {
	return int(ju)
}

func ComputeDaXian(core Core) [12]DaXianStep {
	forward := isDaXianForward(core.Gender, core.NianGan)
	qiSui := daXianQiSui(core.JuShu)
	var steps [12]DaXianStep
	pos := GongIndex(0)
	// 大限起于出生年：虚岁 QiSui → 公历年 = birthYear + QiSui − 1（虚岁 1 岁 = 出生年）。
	startY := core.BirthYear + qiSui - 1
	for i := 0; i < 12; i++ {
		steps[i] = DaXianStep{
			QiSui:     qiSui + i*10,
			ZhiSui:    qiSui + i*10 + 9,
			StartYear: startY + i*10,
			EndYear:   startY + i*10 + 9,
			MingGong:  pos,
		}
		if forward {
			pos = (pos + 11) % 12 // 经典顺行→逆Mingli序往后走
		} else {
			pos = (pos + 1) % 12 // 经典逆行→顺Mingli序往前走
		}
	}
	return steps
}

func isDaXianForward(gender ganzhi.Gender, nianGan Gan) bool {
	isYang := int(nianGan)%2 == 1
	isMale := gender == Male
	// 阳男阴女→顺行（与长生方向相反）
	return (isMale && isYang) || (!isMale && !isYang)
}
