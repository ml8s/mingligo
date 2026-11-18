package main

import (
	"strings"
	"testing"

	"github.com/ml8s/mingligo/examples/mingligo-mcp/internal/observability"
)

func TestTypedToolCallRecordsMetrics(t *testing.T) {
	observability.ResetMetrics()
	sess, _ := startMCPServer(t)
	callToolStructured(t, sess, "time_now", map[string]any{
		"timezone": "Asia/Shanghai",
	})

	text := string(observability.RenderMetrics())
	for _, fragment := range []string{
		`domain="time",tool="time_now",status="success"`,
		"mingligo_tool_request_duration_seconds",
	} {
		if !strings.Contains(text, fragment) {
			t.Fatalf("metrics missing %q:\n%s", fragment, text)
		}
	}
}
