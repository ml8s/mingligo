package fengshui

import (
	"os"
	"testing"

	"github.com/ml8s/mingligo/pkg/luoshu"
)

func TestMain(m *testing.M) {
	if err := luoshu.Load(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
