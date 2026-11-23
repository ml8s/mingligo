package liuyao

import (
	"encoding/json"
	"fmt"
)

// EncodeYongShen converts a Chinese yongshen name to a YongShen value.
func EncodeYongShen(s string) (YongShen, error) {
	for i, name := range yongShenNames {
		if name == s {
			return YongShen(i), nil
		}
	}
	return -1, fmt.Errorf("unknown yongshen: %q", s)
}

// YaoType classifies a single coin-toss result.
type YaoType int

const (
	LaoYin   YaoType = 6 // 老阴 ⚋→⚊ (6, 三个2 = 三背)
	ShaoYang YaoType = 7 // 少阳 ⚊ (7, 两背一面)
	ShaoYin  YaoType = 8 // 少阴 ⚋ (8, 两面一背)
	LaoYang  YaoType = 9 // 老阳 ⚊→⚋ (9, 三个3 = 三面)
)

func (y YaoType) IsYang() bool     { return y == ShaoYang || y == LaoYang }
func (y YaoType) IsChanging() bool { return y == LaoYang || y == LaoYin }

// LiuQin is the 六亲 classification.
type LiuQin int

const (
	QinFumu    LiuQin = iota // 父母
	QinXiongDi               // 兄弟
	QinGuanGui               // 官鬼
	QinQiCai                 // 妻财
	QinZiSun                 // 子孙
)

var liuQinNames = [5]string{"父母", "兄弟", "官鬼", "妻财", "子孙"}

// DecodeLiuQin decodes a LiuQin symbol into its display text (authoritative mapping).
func DecodeLiuQin(v LiuQin) string {
	if v >= 0 && v <= 4 {
		return liuQinNames[v]
	}
	return "?"
}

func (q LiuQin) MarshalJSON() ([]byte, error) { return json.Marshal(DecodeLiuQin(q)) }

func (q *LiuQin) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("liu_qin: expected string, got %s", string(data))
	}
	for i, n := range liuQinNames {
		if n == s {
			*q = LiuQin(i)
			return nil
		}
	}
	return fmt.Errorf("liu_qin: unknown %q", s)
}

// LiuShou is the 六兽.
type LiuShou int

const (
	ShouQingLong LiuShou = iota
	ShouZhuQue
	ShouGouChen
	ShouTengShe
	ShouBaiHu
	ShouXuanWu
)

var liuShouNames = [6]string{"青龙", "朱雀", "勾陈", "螣蛇", "白虎", "玄武"}

// DecodeLiuShou decodes a LiuShou symbol into its display text (authoritative mapping).
func DecodeLiuShou(v LiuShou) string {
	if v >= 0 && v <= 5 {
		return liuShouNames[v]
	}
	return "?"
}

func (l LiuShou) MarshalJSON() ([]byte, error) { return json.Marshal(DecodeLiuShou(l)) }

func (l *LiuShou) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("liu shou must be a string, got %s", string(data))
	}
	for i, n := range liuShouNames {
		if n == name {
			*l = LiuShou(i)
			return nil
		}
	}
	return fmt.Errorf("unknown liushou: %q", name)
}

// MuKuSource 是入墓来源（日墓/动墓/化墓），域内符号双向映射。
type MuKuSource uint8

const (
	MuKuDay         MuKuSource = iota // day: 随鬼入日墓
	MuKuMoving                        // moving: 随鬼入动墓
	MuKuTransformed                   // transformed: 随鬼化墓
)

var muKuSourceNames = [...]string{"日墓", "动墓", "化墓"}

// DecodeMuKuSource decodes a MuKuSource symbol into its display text (authoritative mapping).
func DecodeMuKuSource(v MuKuSource) string {
	if int(v) < len(muKuSourceNames) {
		return muKuSourceNames[v]
	}
	return ""
}

func (m MuKuSource) MarshalJSON() ([]byte, error) { return json.Marshal(DecodeMuKuSource(m)) }

func EncodeMuKuSource(name string) MuKuSource {
	for i, n := range muKuSourceNames {
		if n == name {
			return MuKuSource(i)
		}
	}
	return MuKuSource(255)
}

// ShiYing 是爻的世应标记（世/应/无）。
type ShiYing uint8

const (
	ShiYingNone ShiYing = iota
	ShiYingShi          // 世
	ShiYingYing         // 应
)

// DecodeShiYing decodes a ShiYing symbol into its display text (authoritative mapping).
func DecodeShiYing(v ShiYing) string {
	switch v {
	case ShiYingShi:
		return "世"
	case ShiYingYing:
		return "应"
	default:
		return ""
	}
}

func (s ShiYing) MarshalJSON() ([]byte, error) { return json.Marshal(DecodeShiYing(s)) }

func EncodeShiYing(name string) ShiYing {
	switch name {
	case "世":
		return ShiYingShi
	case "应":
		return ShiYingYing
	default:
		return ShiYingNone
	}
}

// guaIndex identifies one of the 64 hexagrams.
type guaIndex int // 0-63, upper trigram 0-7, lower trigram 0-7

// DecodeGuaIndex decodes a guaIndex symbol into its hexagram name.
func DecodeGuaIndex(g guaIndex) string {
	if int(g) < 0 || int(g) >= len(zhouyiTable) {
		return ""
	}
	return zhouyiTable[g].Name
}

func (g guaIndex) Name() string { return DecodeGuaIndex(g) }

func (g guaIndex) MarshalJSON() ([]byte, error) {
	return json.Marshal(DecodeGuaIndex(g))
}

func (g *guaIndex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("卦名(ben_gua) must be a string, got %s", string(data))
	}
	for i, gc := range zhouyiTable {
		if gc.Name == s {
			*g = guaIndex(i)
			return nil
		}
	}
	return fmt.Errorf("unknown gua: %q", s)
}

// DayRelationType 是日建与爻的单一关系类型（冲/合/生/扶/克/平）。
type DayRelationType uint8

const (
	DayChong DayRelationType = iota // 冲
	DayHe                           // 合
	DaySheng                        // 生
	DayFu                           // 扶
	DayKe                           // 克
	DayPing                         // 平
)

var dayRelationNames = [...]string{"冲", "合", "生", "扶", "克", "平"}

// DecodeDayRelationType decodes a DayRelationType symbol into its display text (authoritative mapping).
func DecodeDayRelationType(v DayRelationType) string {
	if int(v) < len(dayRelationNames) {
		return dayRelationNames[v]
	}
	return ""
}

func (r DayRelationType) MarshalJSON() ([]byte, error) { return json.Marshal(DecodeDayRelationType(r)) }

func EncodeDayRelationType(name string) DayRelationType {
	for i, n := range dayRelationNames {
		if n == name {
			return DayRelationType(i)
		}
	}
	return DayRelationType(255)
}

// FeiFuRelation 是飞神与伏神的五行关系（飞伏域内符号，双向映射）。
type FeiFuRelation uint8

const (
	FeiFuBiHe    FeiFuRelation = iota // 飞伏比和
	FeiFuShengFu                      // 飞来生伏
	FeiFuKeFu                         // 飞来克伏
	FeiFuFuSheng                      // 伏来生飞
	FeiFuFuKe                         // 伏来克飞
	FeiFuNone                         // 无
)

var feiFuNames = [...]string{"飞伏比和", "飞来生伏", "飞来克伏", "伏来生飞", "伏来克飞", ""}

// DecodeFeiFuRelation decodes a FeiFuRelation symbol into its display text (authoritative mapping).
func DecodeFeiFuRelation(v FeiFuRelation) string {
	if int(v) < len(feiFuNames) {
		return feiFuNames[v]
	}
	return ""
}

func (r FeiFuRelation) MarshalJSON() ([]byte, error) { return json.Marshal(DecodeFeiFuRelation(r)) }

func EncodeFeiFuRelation(name string) FeiFuRelation {
	for i, n := range feiFuNames {
		if n == name {
			return FeiFuRelation(i)
		}
	}
	return FeiFuNone
}

// DongYaoRelationType 动爻与用神的关系类型（枚举，正交化）。
type DongYaoRelationType string

const (
	RelationShengYong DongYaoRelationType = "生用"  // 动爻生用神
	RelationKeYong    DongYaoRelationType = "克用"  // 动爻克用神
	RelationBiHe      DongYaoRelationType = "比和"  // 动爻与用神五行相同
	RelationChongYong DongYaoRelationType = "冲用"  // 动爻地支冲用神
	RelationShengYuan DongYaoRelationType = "生原神" // 动爻生原神
	RelationKeYuan    DongYaoRelationType = "克原神" // 动爻克原神
	RelationShengJi   DongYaoRelationType = "生忌神" // 动爻生忌神
	RelationKeJi      DongYaoRelationType = "克忌神" // 动爻克忌神
)

// PatternType 特殊格局类型
type PatternType string

const (
	PatternXunKong    PatternType = "旬空"   // 旬空
	PatternYuePo      PatternType = "月破"   // 月破
	PatternFeiFu      PatternType = "飞伏"   // 飞伏
	PatternJinTui     PatternType = "进退"   // 进退神
	PatternChongHe    PatternType = "冲合"   // 六冲/六合
	PatternFanYin     PatternType = "反吟"   // 动爻变支六冲
	PatternFuYin      PatternType = "伏吟"   // 动爻变支相同
	PatternSuiGuiRuMu PatternType = "随鬼入墓" // 随鬼入墓
	PatternDuFa       PatternType = "独发"   // 独发
	PatternDuJing     PatternType = "独静"   // 独静
	PatternLiangXian  PatternType = "两现"   // 用神两现
)

type ForceRole string

const (
	RoleYongShen ForceRole = "用神"
	RoleYuanShen ForceRole = "原神"
	RoleJiShen   ForceRole = "忌神"
	RoleChouShen ForceRole = "仇神"
)

// YongShen categorizes what the user is asking about.
type YongShen int

const (
	YongFumu    YongShen = iota // 父母: 长辈、文书、房屋
	YongXiongDi                 // 兄弟: 朋友、同事、竞争
	YongGuanGui                 // 官鬼: 工作、官运、疾病
	YongQiCai                   // 妻财: 财运、妻子、物品
	YongZiSun                   // 子孙: 子女、健康、宠物
	YongShiYao                  // 世爻: 自身、求问人
	YongYingYao                 // 应爻: 对方、外部环境
)

var yongShenNames = [7]string{"父母", "兄弟", "官鬼", "妻财", "子孙", "世爻", "应爻"}

// DecodeYongShen decodes a YongShen symbol into its display text.
func DecodeYongShen(v YongShen) string { return yongShenNames[v] }

// TargetRole is the role a line plays for a branch-relation fact.
// Its value is the stable code used in output.
type TargetRole string

const (
	targetYongShen TargetRole = "用神"
	targetWorld    TargetRole = "世爻"
	targetOther    TargetRole = "应爻"
	targetMoving   TargetRole = "动爻"
)

// LifeStage is one of the twelve life-cycle stages (长生十二宫).
// Backed by ganzhi's ChangShengTable index; value is the display name.
type LifeStage string

// PatternSubType is the specific kind of a 格局 pattern. Its value IS the
// stable display text; the set is closed by the compute* producers.
type PatternSubType string

// BranchRelationType is a zhi-zhi relation (六冲/六合/六害/自刑/成组刑).
type BranchRelationType string

func (v *MuKuSource) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	decoded := EncodeMuKuSource(name)
	if DecodeMuKuSource(decoded) != name {
		return fmt.Errorf("unknown mu-ku source %q", name)
	}
	*v = decoded
	return nil
}

func (v *ShiYing) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	decoded := EncodeShiYing(name)
	if DecodeShiYing(decoded) != name {
		return fmt.Errorf("unknown shi-ying %q", name)
	}
	*v = decoded
	return nil
}

func (v *DayRelationType) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	decoded := EncodeDayRelationType(name)
	if DecodeDayRelationType(decoded) != name {
		return fmt.Errorf("unknown day relation %q", name)
	}
	*v = decoded
	return nil
}

func (v *FeiFuRelation) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	decoded := EncodeFeiFuRelation(name)
	if DecodeFeiFuRelation(decoded) != name {
		return fmt.Errorf("unknown fei-fu relation %q", name)
	}
	*v = decoded
	return nil
}
