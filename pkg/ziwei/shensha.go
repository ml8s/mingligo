package ziwei

// ShenSha is a ziwei-domain symbol for an adjective star (杂曜/神煞).
// Every value has a stable display name and a reverse mapping.
type ShenSha int

const (
	// 年支系
	ssHuaGai ShenSha = iota // 华盖
	ssXianChi               // 咸池
	ssGuChen                // 孤辰
	ssGuaSu                 // 寡宿
	ssPoSui                 // 破碎
	ssFeiLian               // 蜚廉
	ssLongChi               // 龙池
	ssFengGe                // 凤阁
	ssTianKu                // 天哭
	ssTianXu                // 天虚
	ssTianKong              // 天空
	ssTianDe                // 天德
	ssYueDe                 // 月德
	// 月支系
	ssTianYao               // 天姚
	ssTianXing              // 天刑
	ssYinSha                // 阴煞
	ssJieShen               // 解神
	ssTianYue               // 天月
	ssTianWu                // 天巫
	// 时支系
	ssTaiFu                 // 台辅
	ssFengGao               // 封诰
	// 年干系
	ssTianChu               // 天厨
	ssTianGuan              // 天官
	ssTianFu                // 天福
	ssJieLu                 // 截路
	ssKongWang              // 空亡
	// 复合系 / 直接 add
	ssHongLuan              // 红鸾
	ssTianXi                // 天喜
	ssSanTai                // 三台
	ssBaZuo                 // 八座
	ssEnGuang               // 恩光
	ssTianGui               // 天贵
	ssXunKong               // 旬空
	ssNianJie               // 年解
	ssTianCai               // 天才
	ssTianShou              // 天寿
	ssTianShang             // 天伤
	ssTianShi               // 天使
)

var shenShaNames = [...]string{
	"华盖", "咸池", "孤辰", "寡宿", "破碎", "蜚廉", "龙池", "凤阁",
	"天哭", "天虚", "天空", "天德", "月德",
	"天姚", "天刑", "阴煞", "解神", "天月", "天巫",
	"台辅", "封诰",
	"天厨", "天官", "天福", "截路", "空亡",
	"红鸾", "天喜", "三台", "八座", "恩光", "天贵",
	"旬空", "年解", "天才", "天寿", "天伤", "天使",
}

func (s ShenSha) String() string {
	if int(s) >= 0 && int(s) < len(shenShaNames) {
		return shenShaNames[s]
	}
	return ""
}

func ParseShenSha(name string) (ShenSha, bool) {
	for i, n := range shenShaNames {
		if n == name {
			return ShenSha(i), true
		}
	}
	return 0, false
}
