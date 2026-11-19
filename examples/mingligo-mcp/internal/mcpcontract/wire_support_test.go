package mcpcontract

import "encoding/json"

// wireFromDomain converts a domain value into its typed wire DTO through JSON.
// It is retained only as the oracle in projection equivalence tests: production
// handlers no longer round-trip through JSON, they project domain values
// directly. The comparison guards those projections against the legacy shape.
func wireFromDomain[Out any](value any) (Out, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		var zero Out
		return zero, err
	}
	var out Out
	if err := json.Unmarshal(raw, &out); err != nil {
		var zero Out
		return zero, err
	}
	return out, nil
}
