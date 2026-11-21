package mcpcontract

import (
	"strings"
)

// fixtureTool derives the tool name from a fixture file name such as
// "bazi_chart_1984_liunian" or "qimen_hour_zhuanpan_default".
func fixtureTool(name string) string {
	domain, _, ok := strings.Cut(name, "_")
	if !ok {
		return ""
	}
	if domain == "qimen" {
		return "qimen_chart"
	}
	rest := name[len(domain)+1:]
	tokens := strings.Split(rest, "_")
	for index := len(tokens) - 1; index >= 0; index-- {
		switch tokens[index] {
		case "liunian", "liuyue", "liuri", "liushi", "xiaoyun", "bond", "chart":
			switch tokens[index] {
			case "liunian", "liuyue", "liuri", "liushi", "xiaoyun":
				return domain + "_liupan"
			default:
				return domain + "_" + tokens[index]
			}
		}
	}
	return ""
}

// diffJSON lists structural differences between two decoded JSON documents.
