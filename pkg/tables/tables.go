// Package tables owns engine data-table loading. Every embedded table must be
// loaded before a tool computes; the engine MCP server calls Load from main
// before serving, and the contract tests call it from TestMain.
package tables

import (
	"github.com/ml8s/mingligo/pkg/bazhai"
	"github.com/ml8s/mingligo/pkg/bazi"
	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/huangli"
	"github.com/ml8s/mingligo/pkg/liuyao"
	"github.com/ml8s/mingligo/pkg/qimen"
	"github.com/ml8s/mingligo/pkg/xuankong"
	"github.com/ml8s/mingligo/pkg/ziwei"
)

// Load initializes every embedded engine data table. It is idempotent and must
// be called before any tool runs. Each domain Load pulls in its own
// dependencies, so listing the leaves and domains is enough.
func Load() error {
	for _, load := range []func() error{
		ganzhi.Load,
		bazi.Load,
		ziwei.Load,
		liuyao.Load,
		qimen.Load,
		huangli.Load,
		bazhai.Load,
		xuankong.Load,
	} {
		if err := load(); err != nil {
			return err
		}
	}
	return nil
}
