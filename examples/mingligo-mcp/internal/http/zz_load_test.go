package http

import (
	"os"
	"testing"

	"github.com/ml8s/mingligo/pkg/tables"
)

func TestMain(m *testing.M) {
	if err := tables.Load(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
