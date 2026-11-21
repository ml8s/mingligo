package mcpcontract

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ml8s/mingligo/pkg/bazi"
)

// TestBaziBondDomainMatchesWire guards that the domain Bond's own JSON equals
// the wire DTO shape (nested side objects). The handler now returns the domain
// value directly.
func TestBaziBondDomainMatchesWire(t *testing.T) {
	registry := testContractServer()
	input := BaziBondInput{
		A: BaziBondSideInput{Subject: NatalSubject{StandardTime: "1984-02-15T00:00:00Z", Longitude: 116.4, Gender: "male"}},
		B: BaziBondSideInput{Subject: NatalSubject{StandardTime: "1985-06-20T10:00:00Z", Longitude: 116.4, Gender: "female"}},
	}
	raw, err := registry.Call(context.Background(), "bazi_bond", mustJSON(t, input))
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.ValidateOutput("bazi_bond", raw); err != nil {
		t.Fatalf("output contract: %v", err)
	}
	coreA, err := input.A.Subject.baziCore()
	if err != nil {
		t.Fatal(err)
	}
	coreB, err := input.B.Subject.baziCore()
	if err != nil {
		t.Fatal(err)
	}
	domain := bazi.ComputeBond(coreA, coreB)
	domainRaw, _ := json.Marshal(domain)
	var a, b any
	if err := json.Unmarshal(raw, &a); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(domainRaw, &b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("domain/wire mismatch\ndomain: %s\nwire:   %s", domainRaw, raw)
	}
}
