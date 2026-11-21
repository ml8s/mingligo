package mcpcontract

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/bazi"
)

// TestBaziLiuNianDomainMatchesView guards the handler's direct domain output
// against the legacy wire DTO shape.
func TestBaziLiuNianDomainMatchesView(t *testing.T) {
	subject := NatalSubject{StandardTime: "1984-02-15T00:00:00Z", Longitude: 116.4, Gender: "male"}
	core, err := subject.baziCore()
	if err != nil {
		t.Fatal(err)
	}
	standard, _ := time.Parse(time.RFC3339, "2024-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	domain, err := bazi.ComputeLiuNianAt(core, standard, loc)
	if err != nil {
		t.Fatal(err)
	}
	dj, _ := json.Marshal(domain)
	vj, _ := json.Marshal(domain.View())
	var a, b any
	if err := json.Unmarshal(dj, &a); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(vj, &b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("domain/view mismatch\ndomain: %s\nview:   %s", dj, vj)
	}
}
