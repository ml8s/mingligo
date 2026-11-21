// Wire enum types used by MingliGo tool contracts. Edit value sets directly
// when a domain grows.

package mcpcontract

import (
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

// registerWireEnum binds a generated wire enum to its JSON Schema values.
func registerWireEnum[T ~string](options *jsonschema.ForOptions, values ...string) {
	enum := make([]any, 0, len(values))
	for _, value := range values {
		enum = append(enum, value)
	}
	options.TypeSchemas[reflect.TypeFor[T]()] = &jsonschema.Schema{Type: "string", Enum: enum}
}

func registerWireEnums(options *jsonschema.ForOptions) {
	registerWireEnum[CalendarKindName](options, "solar", "lunar")
	registerWireEnum[BaziLiupanScopeName](options, "liu_nian", "liu_yue", "liu_ri", "liu_shi", "xiao_yun")
	registerWireEnum[DunSourceName](options, "hour_branch", "solar_term")
	registerWireEnum[HourBoundaryName](options, "late_zi", "zi_zheng")
	registerWireEnum[MethodDingjuMethodName](options, "chaibu", "maoshan", "none", "zhirun")
	registerWireEnum[MethodQuarterRuleName](options, "ten_minute_sanyuan", "twelve_minute_ten_division")
	registerWireEnum[MethodSchoolName](options, "jinhan_yujing", "luoshu_feipan", "mingfa_feipan", "zhuanpan")
	registerWireEnum[MethodScopeName](options, "day", "hour", "month", "quarter", "year")
	registerWireEnum[ZiweiPalaceName](options, "仆役", "兄弟", "命宫", "夫妻", "子女", "官禄", "父母", "田宅", "疾厄", "福德", "财帛", "迁移")
}

type BaziLiupanScopeName string

type DunSourceName string

type HourBoundaryName string

type MethodDingjuMethodName string

type MethodQuarterRuleName string

type MethodSchoolName string

type MethodScopeName string

type ZiweiPalaceName string

type ZiweiLiupanScopeName string
