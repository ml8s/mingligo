package mcpcontract

// normalizeProjectionKeys recursively renames keys in a decoded JSON value so
// that a domain payload can be compared in full against its View payload after
// an intentional key rename. It does not weaken the assertion: every leaf and
// every array element is still compared by reflect.DeepEqual in the caller.
func normalizeProjectionKeys(value any, renames map[string]string) any {
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, child := range v {
			newKey := key
			if renamed, ok := renames[key]; ok {
				newKey = renamed
			}
			out[newKey] = normalizeProjectionKeys(child, renames)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, child := range v {
			out[i] = normalizeProjectionKeys(child, renames)
		}
		return out
	default:
		return value
	}
}

// foldLiuQi merges the domain's flat liu_qi_gan/liu_qi_zhi into the View's
// liu_qi object before key normalization.
func foldLiuQi(value any) any {
	obj, ok := value.(map[string]any)
	if !ok {
		return value
	}
	gan, hasGan := obj["liu_qi_gan"]
	zhi, hasZhi := obj["liu_qi_zhi"]
	if !hasGan && !hasZhi {
		return value
	}
	out := make(map[string]any, len(obj))
	for k, v := range obj {
		if k == "liu_qi_gan" || k == "liu_qi_zhi" {
			continue
		}
		out[k] = v
	}
	out["liu_qi"] = map[string]any{"gan": gan, "zhi": zhi}
	return out
}
