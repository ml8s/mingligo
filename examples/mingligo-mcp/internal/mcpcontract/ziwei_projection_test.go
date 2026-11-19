package mcpcontract

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestZiweiChartDomainMatchesWire guards that the domain chart's own JSON
// equals the wire DTO shape. If it holds, the handler can return the domain
// value directly (like the other chart/flow tools) with no projection.
func TestZiweiChartDomainMatchesWire(t *testing.T) {
	subject := NatalSubject{
		StandardTime: "2024-02-09T23:30:00-05:00",
		Longitude:    -74.0,
		Gender:       "male",
	}
	chart, err := subject.ziweiChart()
	if err != nil {
		t.Fatal(err)
	}
	wire, err := wireFromDomain[ZiweiChartOutput](chart)
	if err != nil {
		t.Fatal(err)
	}
	domainRaw, _ := json.Marshal(chart)
	wireRaw, _ := json.Marshal(wire)
	var a, b any
	json.Unmarshal(domainRaw, &a)
	json.Unmarshal(wireRaw, &b)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("domain/wire mismatch\ndomain: %s\nwire:   %s", domainRaw, wireRaw)
	}
}
