// MingliGo wire DTOs for the MCP contract layer.
//
// These types are ordinary source now, not generated artifacts. Edit them
// directly when the published MingliGo payload changes.
package mcpcontract

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

type qimenScopeSchoolDingjuForm struct {
	Scope        MethodScopeName        `json:"scope" jsonschema:"时间层级。"`
	School       MethodSchoolName       `json:"school" jsonschema:"盘式流派。"`
	DingjuMethod MethodDingjuMethodName `json:"dingju_method" jsonschema:"定局法。"`
}

type qimenScopeSchoolForm struct {
	Scope  MethodScopeName  `json:"scope" jsonschema:"时间层级。"`
	School MethodSchoolName `json:"school" jsonschema:"盘式流派。"`
}

type qimenQuarterTenMethodForm struct {
	Scope       MethodScopeName       `json:"scope" jsonschema:"时间层级。"`
	School      MethodSchoolName      `json:"school" jsonschema:"盘式流派。"`
	QuarterRule MethodQuarterRuleName `json:"quarter_rule" jsonschema:"刻家规则。"`
}

type qimenQuarterTwelveMethodForm struct {
	Scope            MethodScopeName         `json:"scope" jsonschema:"时间层级。"`
	School           MethodSchoolName        `json:"school" jsonschema:"盘式流派。"`
	QuarterRule      MethodQuarterRuleName   `json:"quarter_rule" jsonschema:"刻家规则。"`
	BaseDingjuMethod *MethodDingjuMethodName `json:"base_dingju_method,omitempty" jsonschema:"基准时家定局法；缺省为拆补。"`
	DunSource        *DunSourceName          `json:"dun_source,omitempty" jsonschema:"分遁口径；缺省为领域默认值。"`
	HourBoundary     *HourBoundaryName       `json:"hour_boundary,omitempty" jsonschema:"时辰界；缺省为领域默认值。"`
}

type qimenChartStandardForm struct {
	StandardTime StandardTime         `json:"standard_time" jsonschema:"起局地标准时间。"`
	Longitude    Longitude            `json:"longitude" jsonschema:"起局地经度。"`
	Method       QimenMethodInput     `json:"method" jsonschema:"非金函玉镜奇门方法事实。"`
	YongShen     QimenYongShenSymbols `json:"yong_shen,omitempty" jsonschema:"用神符号组合。"`
	Subject      *NatalSubject        `json:"subject,omitempty" jsonschema:"可选命主；提供时用于年命。"`
}

type qimenChartJinhanForm struct {
	StandardTime StandardTime         `json:"standard_time" jsonschema:"起局地标准时间。"`
	Longitude    Longitude            `json:"longitude" jsonschema:"起局地经度。"`
	Method       qimenScopeSchoolForm `json:"method" jsonschema:"金函玉镜方法事实。"`
}

func qimenMethodFormSchema[T any](options *jsonschema.ForOptions,
	constraints map[string]string, enums map[string][]string,
) (json.RawMessage, error) {
	raw, err := inputSchemaFor[T](options)
	if err != nil {
		return nil, err
	}
	for field, value := range constraints {
		raw, err = withSchemaConst(raw, field, value)
		if err != nil {
			return nil, err
		}
	}
	for field, values := range enums {
		raw, err = withSchemaEnum(raw, field, values...)
		if err != nil {
			return nil, err
		}
	}
	return raw, nil
}

func qimenMethodFormsSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	builders := []func() (json.RawMessage, error){
		func() (json.RawMessage, error) {
			raw, err := qimenMethodFormSchema[qimenScopeSchoolDingjuForm](options,
				map[string]string{"scope": "hour", "school": "zhuanpan"},
				map[string][]string{"dingju_method": {"chaibu", "zhirun", "maoshan"}},
			)
			if err == nil {
				raw, err = withSchemaOptionalDefault(raw, "dingju_method", "chaibu")
			}
			return raw, err
		},
		func() (json.RawMessage, error) {
			raw, err := qimenMethodFormSchema[qimenScopeSchoolDingjuForm](options,
				map[string]string{"scope": "hour", "school": "luoshu_feipan"},
				map[string][]string{"dingju_method": {"chaibu", "zhirun"}},
			)
			if err == nil {
				raw, err = withSchemaOptionalDefault(raw, "dingju_method", "chaibu")
			}
			return raw, err
		},
		func() (json.RawMessage, error) {
			return qimenMethodFormSchema[qimenScopeSchoolForm](options,
				map[string]string{"scope": "hour", "school": "mingfa_feipan"}, nil,
			)
		},
		func() (json.RawMessage, error) {
			raw, err := qimenMethodFormSchema[qimenScopeSchoolDingjuForm](options,
				map[string]string{"scope": "day", "school": "zhuanpan"},
				map[string][]string{"dingju_method": {"chaibu", "zhirun"}},
			)
			if err == nil {
				raw, err = withSchemaOptionalDefault(raw, "dingju_method", "chaibu")
			}
			return raw, err
		},
		func() (json.RawMessage, error) {
			raw, err := qimenMethodFormSchema[qimenScopeSchoolDingjuForm](options,
				map[string]string{"scope": "day", "school": "luoshu_feipan"},
				map[string][]string{"dingju_method": {"chaibu", "zhirun"}},
			)
			if err == nil {
				raw, err = withSchemaOptionalDefault(raw, "dingju_method", "chaibu")
			}
			return raw, err
		},
		func() (json.RawMessage, error) {
			return qimenMethodFormSchema[qimenScopeSchoolForm](options,
				map[string]string{"scope": "month"},
				map[string][]string{"school": {"zhuanpan", "luoshu_feipan"}},
			)
		},
		func() (json.RawMessage, error) {
			return qimenMethodFormSchema[qimenScopeSchoolForm](options,
				map[string]string{"scope": "year"},
				map[string][]string{"school": {"zhuanpan", "luoshu_feipan"}},
			)
		},
		func() (json.RawMessage, error) {
			return qimenMethodFormSchema[qimenQuarterTenMethodForm](options,
				map[string]string{
					"scope": "quarter", "school": "zhuanpan",
					"quarter_rule": "ten_minute_sanyuan",
				}, nil,
			)
		},
		func() (json.RawMessage, error) {
			return qimenMethodFormSchema[qimenQuarterTwelveMethodForm](options,
				map[string]string{
					"scope": "quarter", "school": "zhuanpan",
					"quarter_rule": "twelve_minute_ten_division",
				},
				map[string][]string{"base_dingju_method": {"chaibu", "zhirun", "maoshan"}},
			)
		},
	}
	forms := make([]json.RawMessage, 0, len(builders))
	for _, build := range builders {
		raw, err := build()
		if err != nil {
			return nil, err
		}
		forms = append(forms, raw)
	}
	return json.Marshal(map[string]any{
		"type": "object", "oneOf": forms,
	})
}

func qimenJinhanMethodSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	return qimenMethodFormSchema[qimenScopeSchoolForm](options,
		map[string]string{"scope": "day", "school": "jinhan_yujing"}, nil,
	)
}

func qimenChartOutputSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	standard, err := outputSchemaFor[qimen.ChartView](options)
	if err != nil {
		return nil, fmt.Errorf("standard qimen output schema: %w", err)
	}
	jinhan, err := outputSchemaFor[QimenJinhanChartOutput](options)
	if err != nil {
		return nil, fmt.Errorf("jinhan qimen output schema: %w", err)
	}
	return json.Marshal(map[string]any{
		"type":  "object",
		"oneOf": []json.RawMessage{standard, jinhan},
	})
}

func qimenChartFormsSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	methodForms, err := qimenMethodFormsSchema(options)
	if err != nil {
		return nil, err
	}
	standard, err := inputSchemaFor[qimenChartStandardForm](options)
	if err != nil {
		return nil, err
	}
	standard, err = withSchemaNumberBounds(standard, "longitude", -180, 180)
	if err != nil {
		return nil, err
	}
	standard, err = withSchemaProperty(standard, "method", methodForms)
	if err != nil {
		return nil, err
	}

	jinhanMethod, err := qimenJinhanMethodSchema(options)
	if err != nil {
		return nil, err
	}
	jinhan, err := inputSchemaFor[qimenChartJinhanForm](options)
	if err != nil {
		return nil, err
	}
	jinhan, err = withSchemaNumberBounds(jinhan, "longitude", -180, 180)
	if err != nil {
		return nil, err
	}
	jinhan, err = withSchemaProperty(jinhan, "method", jinhanMethod)
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]any{
		"type":  "object",
		"oneOf": []json.RawMessage{standard, jinhan},
	})
}

// ── qimen.chart output ──

// QimenJinhanChartOutput is the public wire projection of the Golden Mirror
// chart. It intentionally keeps the current direct qimen domain JSON shape.
type QimenJinhanChartOutput = qimen.JinhanChart

// QimenChartOutput is the typed union for the two qimen schools. MarshalJSON
// preserves the current direct standard/Jinhan wire faces.
type QimenChartOutput struct {
	Standard *qimen.ChartView        `json:"-"`
	Jinhan   *QimenJinhanChartOutput `json:"-"`
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
		School MethodSchoolName `json:"school"`
	}
	if err := json.Unmarshal(data, &method); err != nil {
		return err
	}
	if method.School == "jinhan_yujing" {
		return json.Unmarshal(data, &output.Jinhan)
	}
	return json.Unmarshal(data, &output.Standard)
}
