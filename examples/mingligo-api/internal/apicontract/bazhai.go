package apicontract

import (
	"context"
	"fmt"

	"github.com/ml8s/mingligo/pkg/bazhai"
	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/luoshu"
)

// BazhaiSubjectInput is the minimal calibrated identity used by Bazhai. It does
// not carry longitude because Bazhai does not use local apparent solar time.
type BazhaiSubjectInput struct {
	StandardTime StandardTime `json:"standard_time" jsonschema:"出生年代表标准时间（无 DST），RFC3339 带 offset；用于立春年界。"`
	Gender       GenderName   `json:"gender" jsonschema:"性别。"`
}

// BazhaiChartInput consumes the Bazhai subject.
type BazhaiChartInput struct {
	Subject BazhaiSubjectInput `json:"subject" jsonschema:"八宅命卦主体。"`
}

type BazhaiLayoutInput struct {
	MingGua   EightGuaName `json:"ming_gua" jsonschema:"命卦。"`
	DoorGua   EightGuaName `json:"men_gua" jsonschema:"门卦。"`
	MasterGua EightGuaName `json:"master_gua" jsonschema:"主卧卦。"`
	StoveGua  EightGuaName `json:"stove_gua" jsonschema:"灶卦。"`
}

func BazhaiChart(_ context.Context, input BazhaiChartInput) (bazhai.BaZhaiMingGuaPanView, error) {
	standard, _, err := parseStandardTime(input.Subject.StandardTime)
	if err != nil {
		return bazhai.BaZhaiMingGuaPanView{}, fmt.Errorf("bazhai.chart: %w", err)
	}
	gender, err := ganzhi.EncodeGender(string(input.Subject.Gender))
	if err != nil {
		return bazhai.BaZhaiMingGuaPanView{}, fmt.Errorf("bazhai.chart: %w", err)
	}
	chart, err := bazhai.ComputeChartAt(gender, standard)
	if err != nil {
		return bazhai.BaZhaiMingGuaPanView{}, fmt.Errorf("bazhai.chart: %w", err)
	}
	return chart.View(), nil
}

func BazhaiLayout(_ context.Context, input BazhaiLayoutInput) (bazhai.MenZhuZaoView, error) {
	mingGua, ok1 := luoshu.EncodeTrigram(string(input.MingGua))
	doorGua, ok2 := luoshu.EncodeTrigram(string(input.DoorGua))
	masterGua, ok3 := luoshu.EncodeTrigram(string(input.MasterGua))
	stoveGua, ok4 := luoshu.EncodeTrigram(string(input.StoveGua))
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return bazhai.MenZhuZaoView{}, fmt.Errorf("bazhai.layout: invalid trigram input")
	}
	layout, err := bazhai.ComputeLayout(mingGua, doorGua, masterGua, stoveGua)
	if err != nil {
		return bazhai.MenZhuZaoView{}, err
	}
	return layout.View(), nil
}
