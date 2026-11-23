package liuyao

import "fmt"

// YongShenFromMatter resolves a占问事项 (and the optional perspective for
// relationship questions) into the用神 category. This is a domain rule: it
// belongs to liuyao, not to any transport adapter.
//
// Rationale:
//   - self/health/travel/general/family: 自占以世爻为用
//   - other: 无尊卑称呼、未深交、泛指对象取应爻
//   - relationship: 男占女以妻财为主，女占男以官鬼为主
//   - parents/home/study: 父母为长辈、文书、房屋
//   - children: 子孙；siblings/competition: 兄弟
//   - wealth/lost_item: 妻财为钱物、所失物之主
//   - career/legal: 官鬼为功名官府、问官
func YongShenFromMatter(matter, perspective string) (YongShen, error) {
	if matter == "relationship" {
		switch perspective {
		case "male":
			return YongQiCai, nil
		case "female":
			return YongGuanGui, nil
		default:
			return -1, fmt.Errorf("relationship requires perspective=male/female")
		}
	}
	if perspective != "" {
		return -1, fmt.Errorf("perspective is only valid for relationship")
	}
	switch matter {
	case "self", "health", "travel", "general", "family":
		return YongShiYao, nil
	case "other":
		return YongYingYao, nil
	case "parents", "home", "study":
		return YongFumu, nil
	case "children":
		return YongZiSun, nil
	case "siblings", "competition":
		return YongXiongDi, nil
	case "wealth", "lost_item":
		return YongQiCai, nil
	case "career", "legal":
		return YongGuanGui, nil
	default:
		return -1, fmt.Errorf("unsupported matter %q", matter)
	}
}
