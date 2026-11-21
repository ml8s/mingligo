package huangli

import (
	"time"

	"github.com/ml8s/mingligo/pkg/tianwen"
)

func mustDate(s string) tianwen.Date {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return tianwen.DateOf(t)
}
