package qimen

// JinhanChartView is the domain-owned output view of the 金函玉镜 daily chart.
// It is a separate branch from ChartView (method=jinhan_yujing).
type JinhanChartView struct {
	Method JinhanChartMethodView `json:"pai_pan_fang_fa" jsonschema:"金函玉镜起局方法维度。"`
	Board  JinhanBoard           `json:"jiu_gong_pan" jsonschema:"金函玉镜九宫盘与十二地神。"`
}

// View maps the computed Jinhan chart into its tagged output contract. The
// mapping is explicit so output keys never depend on domain json tags.
func (c JinhanChart) View() JinhanChartView {
	return JinhanChartView{Method: jinhanChartMethodView(c.Method), Board: c.Board}
}
