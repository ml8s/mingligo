package tianwen

import (
	"os"
	"testing"

	"github.com/ml8s/mingligo/pkg/ganzhi"
)

func TestMain(m *testing.M) {
	if err := ganzhi.Load(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
