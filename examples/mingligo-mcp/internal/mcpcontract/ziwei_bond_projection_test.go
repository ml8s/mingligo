package mcpcontract

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ml8s/mingligo/pkg/ziwei"
)

func TestZiweiBondDomainMatchesView(t *testing.T) {
	registry := testContractServer()
	input := ZiweiBondInput{
		A: ZiweiBondSide{Subject: NatalSubject{StandardTime: "1984-02-15T00:00:00Z", Longitude: 116.4, Gender: "male"}},
		B: ZiweiBondSide{Subject: NatalSubject{StandardTime: "1985-06-20T10:00:00Z", Longitude: 116.4, Gender: "female"}},
	}
	raw, err := registry.Call(context.Background(), "ziwei_bond", mustJSON(t, input))
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.ValidateOutput("ziwei_bond", raw); err != nil {
		t.Fatalf("output contract: %v", err)
	}
	chartA, err := input.A.Subject.ziweiChart()
	if err != nil {
		t.Fatal(err)
	}
	chartB, err := input.B.Subject.ziweiChart()
	if err != nil {
		t.Fatal(err)
	}
	domain := ziwei.ComputeBond(chartA.BondSideFacts(), chartB.BondSideFacts())
	viewRaw, _ := json.Marshal(domain.View())
	var a, b any
	json.Unmarshal(raw, &a)
	json.Unmarshal(viewRaw, &b)
	if !reflect.DeepEqual(a, b) {
		t.Errorf("domain/view mismatch\ndomain: %s\nwire:   %s", viewRaw, raw)
	}
}
