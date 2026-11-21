package liuyao

import "encoding/json"

// QiguaView is the domain-owned output contract for one immutable casting
// event. The adapter patches enum constraints onto the casting and context
// properties at schema-build time; this type is the single handler output.
type QiguaView struct {
	Event QiguaEventView `json:"event" jsonschema:"一次起卦创建的不可变卦象事件。"`
}

type QiguaEventView struct {
	Casting        QiguaCastingView        `json:"casting" jsonschema:"一次起卦的完整源事实。"`
	CastingContext QiguaCastingContextView `json:"casting_context" jsonschema:"起卦地标准时间、经度与解析出的用神。"`
}

type QiguaCastingView struct {
	Mode           string          `json:"mode" jsonschema:"实际使用的起卦模式。"`
	Order          string          `json:"order" jsonschema:"爻序；当前固定为初爻到上爻。"`
	CoinConvention *CoinConvention `json:"coin_convention,omitempty" jsonschema:"coins 模式的铜钱记分约定。"`
	Rounds         []CoinRound     `json:"rounds,omitempty" jsonschema:"coins 模式的六组铜钱记录。"`
	Yaos           [6]int          `json:"yaos" jsonschema:"六个爻值，初爻到上爻。"`
	DongYao        []int           `json:"dong_yao" jsonschema:"动爻位置列表，1-6。"`
}

type QiguaCastingContextView struct {
	StandardTime string  `json:"standard_time" jsonschema:"起卦地标准时间（无 DST），RFC3339 带 offset。"`
	Longitude    float64 `json:"longitude" jsonschema:"起卦地经度（-180..180）；真太阳时校正用于时辰。"`
	YongShen     string  `json:"yong_shen" jsonschema:"起卦时解析出的用神；chart/ask 唯一取用依据。"`
}

// QiguaCastingView maps the pure casting value into its tagged output view.
func (c Casting) QiguaCastingView() QiguaCastingView {
	raw, err := json.Marshal(c)
	if err != nil {
		return QiguaCastingView{}
	}
	var view QiguaCastingView
	if err := json.Unmarshal(raw, &view); err != nil {
		return QiguaCastingView{}
	}
	return view
}
