package ziwei

import (
	"encoding/json"
	"fmt"
)

// ShenSha is a ziwei-domain symbol for an adjective star (杂曜/神煞).
// Every value has a stable display name and a reverse mapping.
type ShenSha int

const (
	// 年支系
	ssHuaGai   ShenSha = iota // 华盖
	ssXianChi                 // 咸池
	ssGuChen                  // 孤辰
	ssGuaSu                   // 寡宿
	ssPoSui                   // 破碎
	ssFeiLian                 // 蜚廉
	ssLongChi                 // 龙池
	ssFengGe                  // 凤阁
	ssTianKu                  // 天哭
	ssTianXu                  // 天虚
	ssTianKong                // 天空
	ssTianDe                  // 天德
	ssYueDe                   // 月德
	// 月支系
	ssTianYao  // 天姚
	ssTianXing // 天刑
	ssYinSha   // 阴煞
	ssJieShen  // 解神
	ssTianYue  // 天月
	ssTianWu   // 天巫
	// 时支系
	ssTaiFu   // 台辅
	ssFengGao // 封诰
	// 年干系
	ssTianChu  // 天厨
	ssTianGuan // 天官
	ssTianFu   // 天福
	ssJieLu    // 截路
	ssKongWang // 空亡
	// 复合系 / 直接 add
	ssHongLuan  // 红鸾
	ssTianXi    // 天喜
	ssSanTai    // 三台
	ssBaZuo     // 八座
	ssEnGuang   // 恩光
	ssTianGui   // 天贵
	ssXunKong   // 旬空
	ssNianJie   // 年解
	ssTianCai   // 天才
	ssTianShou  // 天寿
	ssTianShang // 天伤
	ssTianShi   // 天使
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

// DecodeShenSha decodes a ShenSha symbol into its display text (authoritative mapping).
func DecodeShenSha(v ShenSha) string {
	if int(v) >= 0 && int(v) < len(shenShaNames) {
		return shenShaNames[v]
	}
	return ""
}

func (s ShenSha) MarshalJSON() ([]byte, error) { return []byte(`"` + DecodeShenSha(s) + `"`), nil }

func EncodeShenSha(name string) (ShenSha, bool) {
	for i, n := range shenShaNames {
		if n == name {
			return ShenSha(i), true
		}
	}
	return 0, false
}

// FlowStarKind is the period-independent kind of a flow star (禄/羊/陀/…),
// defined once in the ziwei domain with a stable rank for display ordering.
type FlowStarKind int

const (
	fskNone  FlowStarKind = iota
	fskLu                 // 禄
	fskYang               // 羊
	fskTuo                // 陀
	fskKui                // 魁
	fskYue                // 钺
	fskMa                 // 马
	fskLuan               // 鸾
	fskXi                 // 喜
	fskChang              // 昌
	fskQu                 // 曲
)

var flowStarKindNames = [...]string{"", "禄", "羊", "陀", "魁", "钺", "马", "鸾", "喜", "昌", "曲"}

// DecodeFlowStarKind decodes a FlowStarKind symbol into its display text (authoritative mapping).
func DecodeFlowStarKind(v FlowStarKind) string {
	if int(v) >= 0 && int(v) < len(flowStarKindNames) {
		return flowStarKindNames[v]
	}
	return ""
}

func EncodeFlowStarKind(name string) (FlowStarKind, bool) {
	for i, n := range flowStarKindNames {
		if n != "" && n == name {
			return FlowStarKind(i), true
		}
	}
	return fskNone, false
}

// Brightness level: 0=庙 1=旺 2=利 3=平 4=陷.
type Brightness int

const (
	Miao Brightness = iota
	Wang
	De
	Li
	Ping
	Xian
	Bu
)

var brightnessNames = [7]string{"庙", "旺", "得", "利", "平", "陷", "不"}

// DecodeBrightness decodes a Brightness symbol into its display text (authoritative mapping).
func DecodeBrightness(v Brightness) string { return brightnessNames[v] }

// GongIndex identifies one of 12 palaces (0=命宫 … 11=父母).
type GongIndex int

// DecodeGongIndex decodes a GongIndex symbol into its display text (authoritative mapping).
func DecodeGongIndex(v GongIndex) string {
	if int(v) < 0 || int(v) >= len(gongLabels) {
		return ""
	}
	return gongLabels[v]
}

func (p GongIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(DecodeGongIndex(p))
}

func (p *GongIndex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("宫名(gong_wei[].name) must be a string (e.g. \"命宫\"), got %s", string(data))
	}
	for i, name := range gongLabels {
		if name == s {
			*p = GongIndex(i)
			return nil
		}
	}
	return fmt.Errorf("unknown gong name: %q", s)
}

// StarIndex enumerates all stars (main + minor).
type StarIndex int

// 14 main stars.
const (
	ZiWei     StarIndex = iota // 0  紫微
	TianJi                     // 1  天机
	TaiYang                    // 2  太阳
	WuQu                       // 3  武曲
	TianTong                   // 4  天同
	LianZhen                   // 5  廉贞
	TianFu                     // 6  天府
	TaiYin                     // 7  太阴
	TanLang                    // 8  贪狼
	JuMen                      // 9  巨门
	TianXiang                  // 10 天相
	TianLiang                  // 11 天梁
	QiSha                      // 12 七杀
	PoJun                      // 13 破军
)

// Minor stars (0.6).
const (
	LuCun StarIndex = iota + 14
	TianKui
	TianYue
	ZuoFu
	YouBi
	WenChang
	WenQu
	QingYang
	TuoLuo
	TianMa
	HuoXing
	LingXing
	DiKong
	DiJie
	HongLuan
	TianXi
)

var starNames = map[StarIndex]string{
	ZiWei: "紫微", TianJi: "天机", TaiYang: "太阳", WuQu: "武曲",
	TianTong: "天同", LianZhen: "廉贞", TianFu: "天府", TaiYin: "太阴",
	TanLang: "贪狼", JuMen: "巨门", TianXiang: "天相", TianLiang: "天梁",
	QiSha: "七杀", PoJun: "破军",
	LuCun: "禄存", TianKui: "天魁", TianYue: "天钺",
	ZuoFu: "左辅", YouBi: "右弼", WenChang: "文昌", WenQu: "文曲",
	QingYang: "擎羊", TuoLuo: "陀罗", TianMa: "天马",
	HuoXing: "火星", LingXing: "铃星", DiKong: "地空", DiJie: "地劫",
	HongLuan: "红鸾", TianXi: "天喜",
}

// starName returns the Chinese name of a star.
func starName(s StarIndex) string { return starNames[s] }

// DecodeStarIndex decodes a StarIndex symbol into its display text (authoritative mapping).
func DecodeStarIndex(v StarIndex) string {
	return starNames[v]
}

func (s StarIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(DecodeStarIndex(s))
}

func (s *StarIndex) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("星名(xing_yao) must be a string, got %s", string(data))
	}
	return s.fromName(name)
}

// MarshalText 使 StarIndex 作为 map key 时序列化为星名（Go json 的 map key 只认 TextMarshaler）。
func (s StarIndex) MarshalText() ([]byte, error) {
	name, ok := starNames[s]
	if !ok {
		return []byte(""), nil
	}
	return []byte(name), nil
}

// UnmarshalText 使 StarIndex 作为 map key 时可从星名反序列化。
func (s *StarIndex) UnmarshalText(b []byte) error {
	return s.fromName(string(b))
}

func (s *StarIndex) fromName(name string) error {
	for k, n := range starNames {
		if n == name {
			*s = k
			return nil
		}
	}
	return fmt.Errorf("unknown star: %q", name)
}

// JuShu is the five-element bureau number (2/3/4/5/6).
type JuShu int

const (
	JuWater JuShu = 2 // 水二局
	JuWood  JuShu = 3
	JuMetal JuShu = 4 // 金四局
	JuEarth JuShu = 5 // 土五局
	JuFire  JuShu = 6
)

var juShuCodes = map[JuShu]string{
	JuWater: "water_2",
	JuWood:  "wood_3",
	JuMetal: "metal_4",
	JuEarth: "earth_5",
	JuFire:  "fire_6",
}

// DecodeJuShu decodes a JuShu symbol into its display text (authoritative mapping).
func DecodeJuShu(v JuShu) string {
	return juShuCodes[v]
}

func (j JuShu) MarshalJSON() ([]byte, error) {
	return json.Marshal(DecodeJuShu(j))
}

func (j *JuShu) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ju_shu must be a string (e.g. \"水二局\"), got %s", string(data))
	}
	if s == "" {
		*j = 0
		return nil
	}
	for k, name := range juShuCodes {
		if name == s {
			*j = k
			return nil
		}
	}
	for k, name := range juShuNames {
		if name == s {
			*j = k
			return nil
		}
	}
	return fmt.Errorf("unknown ju: %q", s)
}

var juShuNames = map[JuShu]string{
	JuWater: "水二局",
	JuWood:  "木三局",
	JuMetal: "金四局",
	JuEarth: "土五局",
	JuFire:  "火六局",
}

// juShuName returns the Chinese name of a bureau.
func juShuName(j JuShu) string { return juShuNames[j] }

// SiHua is one of the four transformations.
type SiHua string

const (
	HuaLu   SiHua = "禄"
	HuaQuan SiHua = "权"
	HuaKe   SiHua = "科"
	HuaJi   SiHua = "忌"
)

// FlowStar is a flow-period star symbol: a stable kind plus its display name.
// It keeps the 命理 kind symbolic while carrying the period-prefixed text.
type FlowStar struct {
	Kind FlowStarKind
	Name string
}

func (f FlowStar) MarshalJSON() ([]byte, error) {
	return []byte(`"` + f.Name + `"`), nil
}

// flowStarKindFor maps a star symbol to its flow-star kind.
func flowStarKindFor(sid StarIndex) FlowStarKind {
	switch sid {
	case LuCun:
		return fskLu
	case QingYang:
		return fskYang
	case TuoLuo:
		return fskTuo
	case TianKui:
		return fskKui
	case TianYue:
		return fskYue
	case TianMa:
		return fskMa
	case HongLuan:
		return fskLuan
	case TianXi:
		return fskXi
	case WenChang:
		return fskChang
	case WenQu:
		return fskQu
	default:
		return fskNone
	}
}

func flowStarOf(sid StarIndex, name string) FlowStar {
	return FlowStar{Kind: flowStarKindFor(sid), Name: name}
}

// ChangShengStage, BoShiStage, JiangQianStage, SuiQianStage are the four
// twelve-stage cycles attached to palaces. Values are the stable display names.
type ChangShengStage string
type BoShiStage string
type JiangQianStage string
type SuiQianStage string
