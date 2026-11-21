package mcpcontract

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ml8s/mingligo/pkg/xuankong"
)

func TestXuankongLiuNianDomainMatchesView(t *testing.T) {
	result, err := computeXuankongLiunian(XuankongLiunianInput{
		StandardTime: "2026-06-15T12:00:00+08:00",
		ZuoShan:      2,
		XiangShan:    14,
	})
	if err != nil {
		t.Fatal(err)
	}
	var domain, view any
	dj, _ := json.Marshal(result)
	vj, _ := json.Marshal(xuankong.LiuNianResult(result).View())
	json.Unmarshal(dj, &domain)
	json.Unmarshal(vj, &view)
	if !reflect.DeepEqual(domain, view) {
		t.Errorf("domain/view mismatch\ndomain: %s\nview:   %s", dj, vj)
	}
}
