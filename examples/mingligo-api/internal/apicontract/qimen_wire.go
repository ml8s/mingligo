// MingliGo wire DTOs for the MCP contract layer.
//
// These types are ordinary source now, not generated artifacts. Edit them
// directly when the published MingliGo payload changes.
package apicontract

import (
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/ml8s/mingligo/pkg/qimen"
)

// ── qimen.chart input ──

// QimenMethodInput is the closed method fact consumed by qimen.chart. Its
// published schema is narrowed by qimenChartFormsSchema to the catalog's valid
// combinations; the domain resolver remains the final authority.
type QimenMethodInput struct {
	Scope            MethodScopeName        `json:"scope" jsonschema:"时间层级；与 school 及定局/刻家参数共同构成一个方法分支。"`
	School           MethodSchoolName       `json:"school" jsonschema:"盘式流派；与 scope 组合决定可用参数。"`
	DingjuMethod     MethodDingjuMethodName `json:"dingju_method,omitempty" jsonschema:"定局法；仅 hour/day 分支携带。"`
	QuarterRule      MethodQuarterRuleName  `json:"quarter_rule,omitempty" jsonschema:"刻家规则；仅 quarter 分支携带。"`
	BaseDingjuMethod MethodDingjuMethodName `json:"base_dingju_method,omitempty" jsonschema:"基准时家定局法；仅十二分钟十分局分支携带。"`
	DunSource        DunSourceName          `json:"dun_source,omitempty" jsonschema:"分遁口径；仅十二分钟十分局分支携带。"`
	HourBoundary     HourBoundaryName       `json:"hour_boundary,omitempty" jsonschema:"时辰界；仅十二分钟十分局分支携带。"`
}

// QimenChartInput takes the Qimen event, a closed method object, and optional
// chart-analysis facts. Event, method, and focus stay orthogonal.
type QimenChartInput struct {
	StandardTime StandardTime         `json:"standard_time" jsonschema:"起局地标准时间（无 DST），RFC3339 带 offset；节气换局/置闰按此时刻判定。"`
	Longitude    Longitude            `json:"longitude" jsonschema:"起局地经度（-180..180）；真太阳时校正用于时柱与分钟分法。"`
	Method       QimenMethodInput     `json:"method" jsonschema:"奇门方法事实；scope、盘式流派和条件参数按目录组合传入。"`
	YongShen     QimenYongShenSymbols `json:"yong_shen,omitempty" jsonschema:"用神符号组合；甲不露盘，不能作为独立符号"`
	Subject      *NatalSubject        `json:"subject,omitempty" jsonschema:"可选命主；提供时用于年命。"`
}

// ── qimen.chart output ──

// QimenChartOutput is the typed union for the two qimen schools. MarshalJSON
// preserves the current direct standard/Jinhan wire faces.
type QimenChartOutput struct {
	Standard *qimen.ChartView       `json:"-"`
	Jinhan   *qimen.JinhanChartView `json:"-"`
}

func (output QimenChartOutput) MarshalJSON() ([]byte, error) {
	switch {
	case output.Standard != nil:
		return json.Marshal(output.Standard)
	case output.Jinhan != nil:
		return json.Marshal(output.Jinhan)
	default:
		return nil, fmt.Errorf("qimen chart output is empty")
	}
}

func (output *QimenChartOutput) UnmarshalJSON(data []byte) error {
	var method struct {
		PaiPanFangFa struct {
			School MethodSchoolName `json:"school"`
		} `json:"pai_pan_fang_fa"`
	}
	if err := json.Unmarshal(data, &method); err != nil {
		return err
	}
	if method.PaiPanFangFa.School == "jinhan_yujing" {
		return json.Unmarshal(data, &output.Jinhan)
	}
	return json.Unmarshal(data, &output.Standard)
}

func qimenChartOutputSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	standard, err := outputSchemaFor[qimen.ChartView](options)
	if err != nil {
		return nil, fmt.Errorf("standard qimen output schema: %w", err)
	}
	jinhan, err := outputSchemaFor[qimen.JinhanChartView](options)
	if err != nil {
		return nil, fmt.Errorf("jinhan qimen output schema: %w", err)
	}
	return json.Marshal(map[string]any{
		"type":  "object",
		"oneOf": []json.RawMessage{standard, jinhan},
	})
}
