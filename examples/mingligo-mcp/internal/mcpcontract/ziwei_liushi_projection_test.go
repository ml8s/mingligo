package mcpcontract

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/ziwei"
)

// TestZiweiLiuShiDomainMatchesView guards the handler's direct domain output
// against the legacy wire DTO shape.
func TestZiweiLiuShiDomainMatchesView(t *testing.T) {
	subject := NatalSubject{
		StandardTime: "1984-02-15T00:00:00Z",
		Longitude:    116.4,
		Gender:       "male",
	}
	chart, err := subject.ziweiChart()
	if err != nil {
		t.Fatal(err)
	}
	standard, _ := time.Parse(time.RFC3339, "2030-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	domain, err := ziwei.ComputeLiuShiAt(chart.FlowContext(), standard, loc, 116.4)
	if err != nil {
		t.Fatal(err)
	}
	dj, _ := json.Marshal(domain)
	vj, _ := json.Marshal(domain.View())
	var a, b any
	json.Unmarshal(dj, &a)
	json.Unmarshal(vj, &b)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("domain/view mismatch\ndomain: %s\nview:   %s", dj, vj)
	}
}
