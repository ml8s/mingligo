package fengshui

import (
	"os"
	"testing"

	"github.com/ml8s/mingligo/internal/engine/luoshu"
)

func TestMain(m *testing.M) {
	if err := luoshu.Load(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
