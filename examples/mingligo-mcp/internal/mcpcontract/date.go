package mcpcontract

import "github.com/ml8s/mingligo/pkg/tianwen"

func dateString(date tianwen.Date) DateString {
	return DateString(date.String())
}
