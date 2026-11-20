package qimen

import (
	"encoding/json"
	"fmt"
)

// JinhanSpirit is one of the twelve day spirits of the Jinhan Yu Jing school.
// It is a distinct symbol set: these spirits are not the八神/九神 of the other
// qimen schools.
type JinhanSpirit int

const (
	JinhanDiYao    JinhanSpirit = 1 + iota // 地钥
	JinhanLiuZe                            // 六贼
	JinhanWuFu                             // 五符
	JinhanTianCao                          // 天曹
	JinhanDiFu                             // 地符
	JinhanFengBo                           // 风伯
	JinhanLeiGong                          // 雷公
	JinhanYuShi                            // 雨师
	JinhanFengYun                          // 风云
	JinhanTangFu                           // 唐符
	JinhanGuoYin                           // 国印
	JinhanTianGuan                         // 天关
)

var jinhanSpiritNames = [13]string{
	"", "地钥", "六贼", "五符", "天曹", "地符", "风伯", "雷公", "雨师", "风云", "唐符", "国印", "天关",
}

var jinhanSpiritKeys = [13]string{
	"", "diyao", "liuze", "wufu", "tiancao", "difu", "fengbo", "leigong", "yushi", "fengyun", "tangfu", "guoyin", "tianguan",
}

// DecodeJinhanSpirit returns the display name of a Jinhan day spirit.
func DecodeJinhanSpirit(s JinhanSpirit) string {
	if s >= 1 && s <= 12 {
		return jinhanSpiritNames[s]
	}
	return "?"
}

// DecodeJinhanSpiritKey returns the stable key of a Jinhan day spirit.
func DecodeJinhanSpiritKey(s JinhanSpirit) string {
	if s >= 1 && s <= 12 {
		return jinhanSpiritKeys[s]
	}
	return ""
}

// EncodeJinhanSpirit converts a display name to its JinhanSpirit symbol.
func EncodeJinhanSpirit(name string) (JinhanSpirit, error) {
	for s := JinhanSpirit(1); s <= 12; s++ {
		if jinhanSpiritNames[s] == name {
			return s, nil
		}
	}
	return 0, fmt.Errorf("unknown jinhan day spirit: %q", name)
}

func (s JinhanSpirit) MarshalJSON() ([]byte, error) {
	return json.Marshal(DecodeJinhanSpiritKey(s))
}
