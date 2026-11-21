package mcpcontract

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ml8s/mingligo/pkg/bazi"
)

func TestBaziChartSubjectProducesStableChart(t *testing.T) {
	registry := testContractServer()
	for _, subject := range []NatalSubject{
		{StandardTime: "1990-05-20T11:49:00+08:00", Longitude: 116.4, Gender: "male"},
		{StandardTime: "1984-02-16T08:00:00+08:00", Longitude: 116.4, Gender: "female"},
		{StandardTime: "2000-09-20T08:30:00+08:00", Longitude: 118.5, Gender: "female"},
	} {
		raw, err := registry.Call(context.Background(), "bazi_chart", mustJSON(t, BaziChartInput{Subject: subject}))
		if err != nil {
			t.Fatalf("bazi.chart %s: %v", subject.StandardTime, err)
		}
		var first, second any
		if err := json.Unmarshal(raw, &first); err != nil {
			t.Fatal(err)
		}
		raw, err = registry.Call(context.Background(), "bazi_chart", mustJSON(t, BaziChartInput{Subject: subject}))
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(raw, &second); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(first, second) {
			t.Fatalf("subject %s produced unstable chart", subject.StandardTime)
		}
	}
}

func TestBaziChartSubjectRejectsInvalidLongitude(t *testing.T) {
	registry := testContractServer()
	subject := NatalSubject{StandardTime: "1990-05-20T12:49:00+08:00", Longitude: 500, Gender: "male"}
	if _, err := registry.Call(context.Background(), "bazi_chart", mustJSON(t, BaziChartInput{Subject: subject})); err == nil {
		t.Fatal("invalid subject longitude was accepted")
	}
}

func TestBaziChartViewMatchesOutputSchema(t *testing.T) {
	registry := testContractServer()
	subject := NatalSubject{StandardTime: "1990-05-20T12:49:00+08:00", Longitude: 116.4, Gender: "male"}
	chart, err := subject.baziChart()
	if err != nil {
		t.Fatal(err)
	}
	for name, secondary := range map[string]*bazi.TiaoHouAvailability{
		"present": chart.TiaoHou.Secondary,
		"absent":  nil,
	} {
		t.Run(name, func(t *testing.T) {
			value := chart
			value.TiaoHou.Secondary = secondary
			// The handler transmits the domain-owned ChartView; validate that
			// its serialized form satisfies the schema reflected from that type.
			raw, err := json.Marshal(value.View())
			if err != nil {
				t.Fatal(err)
			}
			if err := registry.ValidateOutput("bazi_chart", raw); err != nil {
				t.Fatalf("view output contract: %v", err)
			}
		})
	}
}
