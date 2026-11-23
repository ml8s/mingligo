package liuyao

// QiguaView is the domain-owned output contract for one immutable casting
// event. The adapter patches enum constraints onto the casting and context
// properties at schema-build time; this type is the single handler output.
type QiguaView struct {
	Event QiguaEventView `json:"qi_gua" jsonschema:"一次起卦创建的不可变卦象事件。"`
}

type QiguaEventView struct {
	Casting        QiguaCastingView        `json:"qi_gua_ji_lu" jsonschema:"一次起卦的完整源事实。"`
	CastingContext QiguaCastingContextView `json:"qi_gua_qing_jing" jsonschema:"起卦地标准时间、经度与解析出的用神。"`
}

type QiguaCastingView struct {
	Mode           string          `json:"fang_shi" jsonschema:"实际使用的起卦模式。"`
	Order          string          `json:"yao_xu" jsonschema:"爻序；当前固定为初爻到上爻。"`
	CoinConvention *CoinConvention `json:"tong_qian_yue_ding,omitempty" jsonschema:"coins 模式的铜钱记分约定。"`
	Rounds         []CoinRound     `json:"tong_qian_ji_lu,omitempty" jsonschema:"coins 模式的六组铜钱记录。"`
	Yaos           [6]int          `json:"yao_zhi" jsonschema:"六个爻值，初爻到上爻。"`
	DongYao        []int           `json:"dong_yao" jsonschema:"动爻位置列表，1-6。"`
}

type QiguaCastingContextView struct {
	StandardTime string  `json:"shi_jian" jsonschema:"起卦地标准时间（无 DST），RFC3339 带 offset。"`
	Longitude    float64 `json:"jing_du" jsonschema:"起卦地经度（-180..180）；真太阳时校正用于时辰。"`
	QuYong       string  `json:"qu_yong" jsonschema:"起卦时确定的取用类别（如妻财/官鬼）；成卦后的用神爻见 chart.yong_shen.gua_ming。"`
}

// QiguaCastingView maps the pure casting value into its tagged output view.
// Mapping is explicit so the View key names never depend on domain json tags.
func (c Casting) QiguaCastingView() QiguaCastingView {
	return QiguaCastingView(c)
}
