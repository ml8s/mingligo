package mcpcontract

import "encoding/json"

// wireFromDomain converts a domain value into its typed wire DTO through JSON.
// The wire DTOs are faithful projections of the domain JSON shapes; output
// schema checks are validated by contract tests.
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
