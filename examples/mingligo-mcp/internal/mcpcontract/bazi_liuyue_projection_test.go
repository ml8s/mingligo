package mcpcontract

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/bazi"
)

// TestBaziLiuYueDomainMatchesWire guards the handler's direct domain output
// against the legacy wire DTO shape.
func TestBaziLiuYueDomainMatchesWire(t *testing.T) {
	subject := NatalSubject{StandardTime: "1984-02-15T00:00:00Z", Longitude: 116.4, Gender: "male"}
	core, err := subject.baziCore()
	if err != nil {
		t.Fatal(err)
	}
	standard, _ := time.Parse(time.RFC3339, "2024-06-15T12:00:00+08:00")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	domain, err := bazi.ComputeLiuYueAt(core, standard, loc)
	if err != nil {
		t.Fatal(err)
	}
	wire, err := wireFromDomain[BaziLiuYueOutput](*domain)
	if err != nil {
		t.Fatal(err)
	}
	dj, _ := json.Marshal(domain)
	wj, _ := json.Marshal(wire)
	var a, b any
	json.Unmarshal(dj, &a)
	json.Unmarshal(wj, &b)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("domain/wire mismatch\ndomain: %s\nwire:   %s", dj, wj)
	}
}
