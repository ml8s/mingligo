package mcpcontract

import (
	"encoding/json"
	"testing"
)

// TestZiweiChartViewMatchesOutputSchema guards that the domain-owned ChartView
// serializes to a shape satisfying the schema reflected from that same type.
func TestZiweiChartViewMatchesOutputSchema(t *testing.T) {
	registry := testContractServer()
	subject := NatalSubject{
		StandardTime: "2024-02-09T23:30:00-05:00",
		Longitude:    -74.0,
		Gender:       "male",
	}
	chart, err := subject.ziweiChart()
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(chart.View())
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.ValidateOutput("ziwei_chart", raw); err != nil {
		t.Fatalf("view output contract: %v", err)
	}
}
