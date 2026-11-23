package mcpcontract

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/ml8s/mingligo/pkg/liuyao"
)

var scanOutputVerbose = os.Getenv("SCAN_OUTPUT") == "1"

// TestScanOutputForMachineCodes is a diagnostic: it calls every tool with a
// representative input and lists every string leaf containing ASCII letters,
// so left-over machine codes in the output can be found exhaustively.
func TestScanOutputForMachineCodes(t *testing.T) {
	if !scanOutputVerbose {
		t.Skip("set SCAN_OUTPUT=1 to run the diagnostic")
	}
	reg := testContractServer()
	calls := []struct {
		name  string
		input any
	}{
		{"time_now", map[string]any{"timezone": "Asia/Shanghai"}},
		{"subject_candidates", civilResolveInput(t, "1990-05-20T03:49:00Z", "Asia/Shanghai", 116.4)},
		{"city_info", CityInfoInput{City: "北京"}},
		{"huangli_days", HuangliDaysInput{StandardTime: "2024-06-15T12:00:00+08:00", Count: HuangliCount(1), Event: "sign"}},
		{"bazi_chart", BaziChartInput{Subject: testNatalSubject(t, standardInputFromInstant(t, "1990-05-20T03:49:00Z", "Asia/Shanghai"), "male")}},
		{"bazhai_chart", BazhaiChartInput{Subject: BazhaiSubjectInput{StandardTime: "1984-06-15T12:00:00+08:00", Gender: "male"}}},
		{"bazhai_layout", BazhaiLayoutInput{MingGua: "乾", DoorGua: "坤", MasterGua: "艮", StoveGua: "兑"}},
		{"xuankong_chart", XuankongChartInput{StandardTime: standardInputFromInstant(t, "2026-07-31T04:00:00Z", "Asia/Shanghai"), ZuoShan: 2, XiangShan: 14}},
		{"xuankong_liunian", XuankongLiunianInput{StandardTime: standardInputFromInstant(t, "2026-06-15T04:00:00Z", "Asia/Shanghai"), ZuoShan: 2, XiangShan: 14}},
		{"liuyao_qigua", LiuyaoQiguaInput{Origin: LiuyaoCastingRequestInput{Mode: "yaos", Yaos: LiuyaoYaoValuesInput{7, 8, 7, 8, 7, 8}}, CastingContext: testCastingContext(t)}},
	}
	// liuyao_chart needs the qigua event
	qiguaRaw, err := reg.Call(context.Background(), "liuyao_qigua", mustJSON(t, calls[9].input))
	if err != nil {
		t.Fatal(err)
	}
	var qv liuyao.QiguaView
	decode(t, qiguaRaw, &qv)
	calls = append(calls, struct {
		name  string
		input any
	}{"liuyao_chart", LiuyaoChartInput{Event: liuyaoCastingEventInput(castingEventFromView(t, qv.Event))}})

	bondA := testNatalSubject(t, standardInputFromInstant(t, "1990-05-20T03:49:00Z", "Asia/Shanghai"), "male")
	bondB := testNatalSubject(t, standardInputFromInstant(t, "1992-08-12T09:00:00Z", "Asia/Shanghai"), "female")
	calls = append(calls,
		struct {
			name  string
			input any
		}{"bazi_bond", map[string]any{"a": map[string]any{"subject": bondA}, "b": map[string]any{"subject": bondB}}},
		struct {
			name  string
			input any
		}{"ziwei_bond", map[string]any{"a": map[string]any{"subject": bondA}, "b": map[string]any{"subject": bondB}}},
		struct {
			name  string
			input any
		}{"bazi_liupan", map[string]any{"scope": "liu_nian", "standard_time": "2026-06-15T12:00:00+08:00", "subject": bondA}},
		struct {
			name  string
			input any
		}{"ziwei_liupan", map[string]any{"scope": "liu_nian", "standard_time": "2026-06-15T12:00:00+08:00", "subject": bondA}},
		struct {
			name  string
			input any
		}{"ziwei_chart", map[string]any{"subject": bondA}},
		struct {
			name  string
			input any
		}{"qimen_chart", map[string]any{"standard_time": "2024-06-15T12:00:00+08:00", "longitude": 116.4, "method": map[string]any{"scope": "hour", "school": "zhuanpan", "dingju_method": "chaibu"}}},
	)

	ascii := regexp.MustCompile(`[A-Za-z]`)
	// known-legitimate tokens (protocol enums, timestamps, ids) are reported
	// too, but grouped so left-over codes stand out.
	seen := map[string][]string{}
	var walk func(path string, v any)
	walk = func(path string, v any) {
		switch x := v.(type) {
		case map[string]any:
			for k, child := range x {
				walk(path+"."+k, child)
			}
		case []any:
			for i, child := range x {
				walk(path+"["+itoa(i)+"]", child)
			}
		case string:
			if ascii.MatchString(x) {
				seen[x] = append(seen[x], path)
			}
		}
	}
	for _, c := range calls {
		raw, err := reg.Call(context.Background(), c.name, mustJSON(t, c.input))
		if err != nil {
			t.Logf("SKIP %s: %v", c.name, err)
			continue
		}
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			t.Fatal(err)
		}
		walk(c.name, v)
	}
	vals := make([]string, 0, len(seen))
	for v := range seen {
		vals = append(vals, v)
	}
	sort.Strings(vals)
	for _, v := range vals {
		paths := seen[v]
		sort.Strings(paths)
		t.Logf("VALUE %q  (%d)  %s", v, len(paths), strings.Join(paths[:min(3, len(paths))], ", "))
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
