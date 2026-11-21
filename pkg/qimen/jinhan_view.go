package qimen

import "encoding/json"

// JinhanChartView is the domain-owned output view of the 金函玉镜 daily chart.
// It is a separate branch from ChartView (method=jinhan_yujing).
type JinhanChartView struct {
	Method JinhanChartMethod `json:"method" jsonschema:"金函玉镜起局方法维度。"`
	Board  JinhanBoard       `json:"pan" jsonschema:"金函玉镜九宫盘与十二地神。"`
}

// View maps the computed Jinhan chart into its tagged output contract.
func (c JinhanChart) View() JinhanChartView {
	raw, err := json.Marshal(c)
	if err != nil {
		return JinhanChartView{}
	}
	var view JinhanChartView
	if err := json.Unmarshal(raw, &view); err != nil {
		return JinhanChartView{}
	}
	return view
}
