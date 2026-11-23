package ziwei

// Star categories are stable domain classifications. They replace the former
// factor-only palace_facts projection.
// StarLei is a stable star-category symbol.
type StarLei string

const (
	StarLeiZhuXing StarLei = "主星"
	StarLeiFuXing          = "辅星"
	StarLeiLiuJi           = "流吉"
	StarLeiWenXing StarLei = "文星"
	StarLeiShaXing StarLei = "煞星"
)

var (
	mainStars = map[StarIndex]bool{
		ZiWei: true, TianJi: true, TaiYang: true, WuQu: true,
		TianTong: true, LianZhen: true, TianFu: true, TaiYin: true,
		TanLang: true, JuMen: true, TianXiang: true, TianLiang: true,
		QiSha: true, PoJun: true,
	}
	auspiciousStars = map[StarIndex]bool{
		TianKui: true, TianYue: true, ZuoFu: true,
		YouBi: true, WenChang: true, WenQu: true,
	}
	literaryStars      = map[StarIndex]bool{WenChang: true, WenQu: true}
	palaceMaleficStars = map[StarIndex]bool{
		QingYang: true, TuoLuo: true, HuoXing: true,
		LingXing: true, DiKong: true, DiJie: true,
	}
)

func starLei(star StarIndex, isMajor bool) []StarLei {
	values := make([]StarLei, 0, 4)
	if isMajor {
		values = append(values, StarLeiZhuXing)
	}
	if auspiciousStars[star] {
		values = append(values, StarLeiLiuJi)
	}
	if literaryStars[star] {
		values = append(values, StarLeiWenXing)
	}
	if palaceMaleficStars[star] {
		values = append(values, StarLeiShaXing)
	}
	if len(values) == 0 {
		values = append(values, StarLeiFuXing)
	}
	return values
}

func newStarInfo(star StarIndex, isMajor bool) StarInfo {
	return StarInfo{Star: star, Name: starName(star), IsMajor: isMajor, Lei: starLei(star, isMajor)}
}
