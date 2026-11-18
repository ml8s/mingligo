package mcpcontract

import "sort"

// ToolNames returns the immutable MingliGo tool surface in stable order. The
// descriptions table is the naming source used by typed registration.
func ToolNames() []string {
	names := make([]string, 0, len(toolDescriptions))
	for name := range toolDescriptions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
