package main

import (
	"os"
	"testing"

	"github.com/ml8s/mingligo/pkg/tables"
)

func TestMain(m *testing.M) {
	if err := tables.Load(); err != nil {
		panic(err)
	}
	if err := engineContracts.MaterializeSchemas(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
