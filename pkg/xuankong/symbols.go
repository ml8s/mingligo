package xuankong

import "fmt"

// Yuan is the closed set of the three 60-year 元 of 三元九运.
type Yuan string

const (
	YuanShang Yuan = "上元"
	YuanZhong Yuan = "中元"
	YuanXia   Yuan = "下元"
)

func DecodeYuan(y Yuan) string { return string(y) }

func EncodeYuan(s string) (Yuan, error) {
	for _, y := range [...]Yuan{YuanShang, YuanZhong, YuanXia} {
		if string(y) == s {
			return y, nil
		}
	}
	return "", fmt.Errorf("xuankong: unknown yuan %q", s)
}

// YunName is the closed set of the nine 20-year 运 names.
type YunName string

const (
	Yun1 YunName = "一运"
	Yun2 YunName = "二运"
	Yun3 YunName = "三运"
	Yun4 YunName = "四运"
	Yun5 YunName = "五运"
	Yun6 YunName = "六运"
	Yun7 YunName = "七运"
	Yun8 YunName = "八运"
	Yun9 YunName = "九运"
)

var yunNames = [...]YunName{"", Yun1, Yun2, Yun3, Yun4, Yun5, Yun6, Yun7, Yun8, Yun9}

func DecodeYunName(y YunName) string { return string(y) }

func EncodeYunName(s string) (YunName, error) {
	for _, y := range yunNames {
		if y != "" && string(y) == s {
			return y, nil
		}
	}
	return "", fmt.Errorf("xuankong: unknown yun %q", s)
}

// RecitationLayer names a 伏吟/反吟 finding on one flying-star layer (or the
// whole chart). The value is the stable display text.
type RecitationLayer string

const (
	RecitationPeriodFuYin    RecitationLayer = "运盘伏吟"
	RecitationMountainFuYin  RecitationLayer = "山星伏吟"
	RecitationFacingFuYin    RecitationLayer = "向星伏吟"
	RecitationAllFuYin       RecitationLayer = "全盘伏吟"
	RecitationPeriodFanYin   RecitationLayer = "运盘反吟"
	RecitationMountainFanYin RecitationLayer = "山星反吟"
	RecitationFacingFanYin   RecitationLayer = "向星反吟"
	RecitationAllFanYin      RecitationLayer = "全盘反吟"
)

func DecodeRecitationLayer(r RecitationLayer) string { return string(r) }

func EncodeRecitationLayer(s string) (RecitationLayer, error) {
	for _, r := range [...]RecitationLayer{
		RecitationPeriodFuYin, RecitationMountainFuYin, RecitationFacingFuYin, RecitationAllFuYin,
		RecitationPeriodFanYin, RecitationMountainFanYin, RecitationFacingFanYin, RecitationAllFanYin,
	} {
		if string(r) == s {
			return r, nil
		}
	}
	return "", fmt.Errorf("xuankong: unknown recitation layer %q", s)
}
