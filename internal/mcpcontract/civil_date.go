package mcpcontract

import "github.com/ml8s/mingligo/internal/engine/tianwen"

func civilDateString(date tianwen.CivilDate) DateString {
	return DateString(date.String())
}
