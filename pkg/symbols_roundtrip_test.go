package pkg_test

import (
	"encoding/json"
	"testing"

	"github.com/ml8s/mingligo/pkg/bazi"
	"github.com/ml8s/mingligo/pkg/liuyao"
	"github.com/ml8s/mingligo/pkg/qimen"
	"github.com/ml8s/mingligo/pkg/ziwei"
)

func symbolRoundtrip[T any](t *testing.T, text string) {
	t.Helper()
	quoted, err := json.Marshal(text)
	if err != nil {
		t.Fatal(err)
	}
	var v T
	if err := json.Unmarshal(quoted, &v); err != nil {
		t.Fatalf("unmarshal %s: %v", text, err)
	}
	got, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(quoted) {
		t.Errorf("roundtrip %s: got %s", quoted, got)
	}
}

func TestSymbolJSONRoundtrip(t *testing.T) {
	t.Run("ziwei.ShenSha", func(t *testing.T) { symbolRoundtrip[ziwei.ShenSha](t, "华盖") })
	t.Run("ziwei.Brightness", func(t *testing.T) { symbolRoundtrip[ziwei.Brightness](t, "庙") })
	t.Run("ziwei.FlowStar", func(t *testing.T) { symbolRoundtrip[ziwei.FlowStar](t, "月禄") })
	t.Run("liuyao.MuKuSource", func(t *testing.T) { symbolRoundtrip[liuyao.MuKuSource](t, "日墓") })
	t.Run("liuyao.ShiYing", func(t *testing.T) { symbolRoundtrip[liuyao.ShiYing](t, "世") })
	t.Run("liuyao.DayRelationType", func(t *testing.T) { symbolRoundtrip[liuyao.DayRelationType](t, "冲") })
	t.Run("liuyao.FeiFuRelation", func(t *testing.T) { symbolRoundtrip[liuyao.FeiFuRelation](t, "飞来生伏") })
	t.Run("qimen.DunSource", func(t *testing.T) { symbolRoundtrip[qimen.DunSource](t, "hour_branch") })
	t.Run("qimen.HourBoundary", func(t *testing.T) { symbolRoundtrip[qimen.HourBoundary](t, "late_zi") })
	t.Run("qimen.PillarLevel", func(t *testing.T) { symbolRoundtrip[qimen.PillarLevel](t, "year") })
	t.Run("bazi.NatalRelationKind", func(t *testing.T) { symbolRoundtrip[bazi.NatalRelationKind](t, "zhi_liu_he") })
}
