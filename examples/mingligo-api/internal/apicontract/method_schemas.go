package apicontract

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/ml8s/mingligo/pkg/bazhai"
	"github.com/ml8s/mingligo/pkg/bazi"
	"github.com/ml8s/mingligo/pkg/liuyao"
	"github.com/ml8s/mingligo/pkg/xuankong"
	"github.com/ml8s/mingligo/pkg/ziwei"
)

type MethodSchema struct {
	Name   string          `json:"name"`
	Input  json.RawMessage `json:"input"`
	Output json.RawMessage `json:"output"`
}

func inSchema[T any](o *jsonschema.ForOptions) (json.RawMessage, error) { return inputSchemaFor[T](o) }
func outSchema[T any](o *jsonschema.ForOptions) (json.RawMessage, error) {
	return outputSchemaFor[T](o)
}

type schemaPair struct {
	input  func(*jsonschema.ForOptions) (json.RawMessage, error)
	output func(*jsonschema.ForOptions) (json.RawMessage, error)
}

var methodSchemaRegistry = []struct {
	name string
	pair schemaPair
}{
	{"time_now", schemaPair{inSchema[TimeNowInput], outSchema[TimeNowOutput]}},
	{"subject_candidates", schemaPair{inSchema[TianwenSubjectCandidatesInput], outSchema[TianwenSubjectCandidatesOutput]}},
	{"city_info", schemaPair{inSchema[CityInfoInput], outSchema[CityInfoOutput]}},
	{"huangli_days", schemaPair{inSchema[HuangliDaysInput], outSchema[huangliDaysResult]}},
	{"bazi_chart", schemaPair{inSchema[BaziChartInput], outSchema[bazi.ChartView]}},
	{"bazi_bond", schemaPair{inSchema[BaziBondInput], outSchema[bazi.BondView]}},
	{"bazi_liupan", schemaPair{inSchema[BaziLiupanInput], outSchema[BaziLiupanOutput]}},
	{"bazhai_chart", schemaPair{inSchema[BazhaiChartInput], outSchema[bazhai.BaZhaiMingGuaPanView]}},
	{"bazhai_layout", schemaPair{inSchema[BazhaiLayoutInput], outSchema[bazhai.MenZhuZaoView]}},
	{"xuankong_chart", schemaPair{inSchema[XuankongChartInput], outSchema[xuankong.ChartView]}},
	{"xuankong_liunian", schemaPair{inSchema[XuankongLiunianInput], outSchema[XuankongLiuNianSeries]}},
	{"liuyao_qigua", schemaPair{inSchema[LiuyaoQiguaInput], outSchema[liuyao.QiguaView]}},
	{"liuyao_chart", schemaPair{inSchema[LiuyaoChartInput], outSchema[LiuyaoChartResultOutput]}},
	{"ziwei_chart", schemaPair{inSchema[ZiweiChartInput], outSchema[ziwei.ChartView]}},
	{"ziwei_liupan", schemaPair{inSchema[ZiweiLiupanInput], outSchema[ZiweiLiupanOutput]}},
	{"ziwei_bond", schemaPair{inSchema[ZiweiBondInput], outSchema[ziwei.BondView]}},
	{"qimen_chart", schemaPair{inSchema[QimenChartInput], qimenChartOutputSchema}},
}

func MethodSchemas() ([]MethodSchema, error) {
	opts, err := cachedSchemaOptions()
	if err != nil {
		return nil, err
	}
	out := make([]MethodSchema, 0, len(methodSchemaRegistry))
	for _, m := range methodSchemaRegistry {
		in, err := m.pair.input(opts)
		if err != nil {
			return nil, err
		}
		o, err := m.pair.output(opts)
		if err != nil {
			return nil, err
		}
		out = append(out, MethodSchema{Name: m.name, Input: in, Output: o})
	}
	return out, nil
}

func MethodNames() []string {
	out := make([]string, 0, len(methodSchemaRegistry))
	for _, m := range methodSchemaRegistry {
		out = append(out, m.name)
	}
	return out
}
