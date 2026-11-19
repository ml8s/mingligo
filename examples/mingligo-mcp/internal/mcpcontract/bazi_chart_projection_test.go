package mcpcontract

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ml8s/mingligo/pkg/bazi"
)

// TestBaziChartDomainMatchesWire guards that the domain chart's own JSON equals
// the wire DTO shape (including the explicit all-null secondary closure). The
// handler now returns the domain value directly.
func TestBaziChartDomainMatchesWire(t *testing.T) {
	subject := NatalSubject{StandardTime: "1990-05-20T12:49:00+08:00", Longitude: 116.4, Gender: "male"}
	chart, err := subject.baziChart()
	if err != nil {
		t.Fatal(err)
	}
	for name, secondary := range map[string]*bazi.TiaoHouAvailability{
		"present": chart.TiaoHou.Secondary,
		"absent":  nil,
	} {
		v := chart
		v.TiaoHou.Secondary = secondary
		wire, err := wireFromDomain[BaziChartOutput](v)
		if err != nil {
			t.Fatal(err)
		}
		domainRaw, _ := json.Marshal(v)
		wireRaw, _ := json.Marshal(wire)
		var a, b any
		json.Unmarshal(domainRaw, &a)
		json.Unmarshal(wireRaw, &b)
		if !reflect.DeepEqual(a, b) {
			t.Errorf("[%s] domain/wire mismatch\ndomain: %s\nwire:   %s", name, domainRaw, wireRaw)
		}
	}
}
