package qimen

// Yuan is the qimen-domain symbol for the upper/middle/lower 元 of a
// solar-term or fu-tou cycle.
type Yuan uint8

const (
	YuanShang Yuan = iota // 上元
	YuanZhong             // 中元
	YuanXia               // 下元
)

var yuanLabels = [...]string{"上元", "中元", "下元"}

func (y Yuan) String() string {
	if int(y) < len(yuanLabels) {
		return yuanLabels[y]
	}
	return ""
}

func ParseYuan(name string) (Yuan, bool) {
	for i, n := range yuanLabels {
		if n == name {
			return Yuan(i), true
		}
	}
	return 0, false
}

// ZhiRunState is the qimen-domain symbol for the 置闰 state of a dingju.
type ZhiRunState uint8

const (
	zsNone ZhiRunState = iota // 无置闰状态（非置闰法）
	zsExact                   // exact 正授
	zsAhead                   // ahead 超神
	zsBehind                  // behind 接气
	zsRepeat                  // repeat 置闰
)

var zhiRunStateNamesByKey = [...]string{"", "exact", "ahead", "behind", "repeat"}
var zhiRunStateLabels = [...]string{"", "正授", "超神", "接气", "置闰"}

func (s ZhiRunState) String() string {
	if int(s) < len(zhiRunStateLabels) {
		return zhiRunStateLabels[s]
	}
	return ""
}

func (s ZhiRunState) Key() string {
	if int(s) < len(zhiRunStateNamesByKey) {
		return zhiRunStateNamesByKey[s]
	}
	return ""
}

func ParseZhiRunState(key string) (ZhiRunState, bool) {
	for i, n := range zhiRunStateNamesByKey {
		if n == key && key != "" {
			return ZhiRunState(i), true
		}
	}
	return 0, false
}
