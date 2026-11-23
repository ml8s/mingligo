package mcpcontract

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ml8s/mingligo/pkg/bazi"
)

// TestBaziLiuYueDomainMatchesView guards the handler's direct domain output
// against the legacy wire DTO shape.
func TestBaziLiuYueDomainMatchesView(t *testing.T) {
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
	view := domain.View()
	vj, _ := json.Marshal(view)
	var doc map[string]any
	if err := json.Unmarshal(vj, &doc); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"liu_qi", "ri_zhu_zuo_yong", "guan_xi", "gong_li", "jie_xian"} {
		if _, ok := doc[key]; !ok {
			t.Errorf("view missing命理子域 %s: %s", key, vj)
		}
	}
	liuQi, _ := doc["liu_qi"].(map[string]any)
	if liuQi["shi_shen"] == nil || liuQi["shi_shen"] == "" {
		t.Errorf("view liu_qi.shi_shen empty: %s", vj)
	}
}
