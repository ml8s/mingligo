package mcpcontract

// This file is MingliGo's typed schema catalogue: the wire vocabulary that
// schemas derive from and helpers shared by registration. Each tool's DTOs
// live beside the tool they describe.

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/ml8s/mingligo/pkg/bazi"
	"github.com/ml8s/mingligo/pkg/ganzhi"
	"github.com/ml8s/mingligo/pkg/liuyao"
	"github.com/ml8s/mingligo/pkg/qimen"
)

// buildStaticSchemaOptions builds wire-level constraints that Go struct tags
// cannot express through this library. Field-level descriptions still come from
// the DTO's `jsonschema` tags.
func buildStaticSchemaOptions() map[reflect.Type]*jsonschema.Schema {
	return map[reflect.Type]*jsonschema.Schema{
		// A few legacy wire fields carry more than one JSON type (a pillar 六破
		// reference is a name in some charts and an index in others). Those fields
		// keep their exact JSON value instead of being reshaped by a Go type.
		reflect.TypeFor[json.RawMessage](): {},
		reflect.TypeFor[StandardTime]():    {Type: "string", Format: "date-time"},
		reflect.TypeFor[DateString]():      {Type: "string", Format: "date"},
		reflect.TypeFor[Longitude]():       {Type: "number", Minimum: jsonschema.Ptr(float64(-180)), Maximum: jsonschema.Ptr(float64(180))},
		reflect.TypeFor[GenderName]():      {Type: "string", Enum: []any{"male", "female"}},
		reflect.TypeFor[CityName]():        {Type: "string", MinLength: jsonschema.Ptr(1), MaxLength: jsonschema.Ptr(128)},
		reflect.TypeFor[CityLevelName]():   {Type: "string", Enum: []any{"city", "county", "world"}},
		reflect.TypeFor[HuangliEvent](): {
			Type: "string",
			Enum: []any{
				"wedding", "engage", "opening", "sign", "move", "travel", "build",
				"exam", "medical", "sacrifice", "cleaning", "renovation",
				"bed_install", "income", "funeral",
			},
		},
		reflect.TypeFor[HuangliCount](): {Type: "integer", Minimum: jsonschema.Ptr(float64(1)), Maximum: jsonschema.Ptr(float64(30))},
		reflect.TypeFor[SeriesCount](): {
			Type: "integer", Minimum: jsonschema.Ptr(float64(1)), Maximum: jsonschema.Ptr(float64(MaxSeriesCount)),
		},
		reflect.TypeFor[GuaNumber](): {Type: "integer", Minimum: jsonschema.Ptr(float64(1)), Maximum: jsonschema.Ptr(float64(9))},
		reflect.TypeFor[MountainIndex](): {
			Type: "integer", Minimum: jsonschema.Ptr(float64(0)), Maximum: jsonschema.Ptr(float64(23)),
		},
		reflect.TypeFor[XuankongYear](): {
			Type: "integer", Minimum: jsonschema.Ptr(float64(1864)), Maximum: jsonschema.Ptr(float64(2200)),
		},
		reflect.TypeFor[FlyingStarNumber](): {
			Type: "integer", Minimum: jsonschema.Ptr(float64(1)), Maximum: jsonschema.Ptr(float64(9)),
		},
		reflect.TypeFor[BazhaiBirthYear](): {
			Type: "integer", Minimum: jsonschema.Ptr(float64(1900)), Maximum: jsonschema.Ptr(float64(2100)),
		},
		reflect.TypeFor[EightGuaName](): {
			Type: "string", Enum: []any{"坎", "坤", "震", "巽", "乾", "兑", "艮", "离"},
		},
		reflect.TypeFor[ElementName](): {
			// Absent candidates are null, not an ambiguous empty-string sentinel.
			Type: "string", Enum: []any{"水", "土", "木", "金", "火"},
		},
		reflect.TypeFor[StarColorName](): {
			Type: "string", Enum: []any{"白", "黑", "碧", "绿", "黄", "赤", "紫"},
		},
		reflect.TypeFor[MountainName](): {
			Type: "string", Enum: []any{
				"子", "癸", "丑", "艮", "寅", "甲", "卯", "乙", "辰", "巽", "巳", "丙",
				"午", "丁", "未", "坤", "申", "庚", "酉", "辛", "戌", "乾", "亥", "壬",
			},
		},
		reflect.TypeFor[XuankongYuanName](): {Type: "string", Enum: []any{"上元", "中元", "下元"}},
		reflect.TypeFor[YunName](): {
			Type: "string", Enum: []any{
				"一运", "二运", "三运", "四运", "五运", "六运", "七运", "八运", "九运",
			},
		},
		reflect.TypeFor[XuankongSituationName](): {
			Type: "string", Enum: []any{
				"旺山旺向", "双星会坐", "双星会向", "上山下水", "未入四大局",
			},
		},
		reflect.TypeFor[XuankongFuYinLayer](): {
			Type: "string", Enum: []any{"运盘伏吟", "山星伏吟", "向星伏吟", "全盘伏吟"},
		},
		reflect.TypeFor[XuankongFanYinLayer](): {
			Type: "string", Enum: []any{"运盘反吟", "山星反吟", "向星反吟"},
		},
		reflect.TypeFor[XuankongAssessment](): {
			Type: "string", Enum: []any{
				"收山得宜、拨水入零堂，理气合局（完整收山出煞须验实际砂水）",
				"收山得宜（坐宫山星当令）；拨水入零堂未得，向首零神不见水",
				"拨水入零堂得宜（向宫向星零神）；收山未得，坐宫山星不当令",
				"收山、拨水入零堂俱未得，宜择时改向并验实际砂水",
			},
		},
		reflect.TypeFor[YinYangName](): {Type: "string", Enum: []any{"阳", "阴"}},
		reflect.TypeFor[DirectionName](): {
			Type: "string", Enum: []any{"北", "东北", "东", "东南", "南", "西南", "西", "西北"},
		},
		reflect.TypeFor[LifeGroupName]():    {Type: "string", Enum: []any{"东四命", "西四命"}},
		reflect.TypeFor[YearBoundaryName](): {Type: "string", Enum: []any{"lichun"}},
		reflect.TypeFor[FlyingStarName](): {
			Type: "string", Enum: []any{
				"一白贪狼", "二黑巨门", "三碧禄存", "四绿文曲", "五黄廉贞",
				"六白武曲", "七赤破军", "八白左辅", "九紫右弼",
			},
		},
		reflect.TypeFor[YouxingName](): {
			Type: "string", Enum: []any{"生气", "天医", "延年", "伏位", "祸害", "五鬼", "六煞", "绝命"},
		},
		reflect.TypeFor[BazhaiRatingName](): {
			Type: "string", Enum: []any{"大吉", "吉", "平", "凶", "大凶"},
		},
		reflect.TypeFor[TrigramGroupName](): {Type: "string", Enum: []any{"东四卦", "西四卦"}},
		reflect.TypeFor[BazhaiMatchName]():  {Type: "string", Enum: []any{"吉", "凶"}},
		reflect.TypeFor[LiuyaoQiguaMode]():  {Type: "string", Enum: []any{"auto", "coins", "yaos"}},
		reflect.TypeFor[LiuyaoCastMode]():   {Type: "string", Enum: []any{"coins", "yaos"}},
		reflect.TypeFor[LiuyaoOrderName]():  {Type: "string", Enum: []any{"bottom_up"}},
		reflect.TypeFor[LiuyaoCoinSymbol](): {
			Type: "string", Enum: []any{"正", "反"},
		},
		reflect.TypeFor[LiuyaoHeadsScore](): integerConstSchema(int(liuyao.CoinHeadsScore), "正面记 3。"),
		reflect.TypeFor[LiuyaoTailsScore](): integerConstSchema(int(liuyao.CoinTailsScore), "反面记 2。"),
		reflect.TypeFor[LiuyaoYaoValue](): {
			Type: "integer", Minimum: jsonschema.Ptr(float64(6)), Maximum: jsonschema.Ptr(float64(9)),
		},
		reflect.TypeFor[LiuyaoOptionalYaoValue](): {
			Type: "integer", Minimum: jsonschema.Ptr(float64(0)), Maximum: jsonschema.Ptr(float64(9)),
		},
		reflect.TypeFor[LiuyaoPosition](): {
			Type: "integer", Minimum: jsonschema.Ptr(float64(1)), Maximum: jsonschema.Ptr(float64(6)),
		},
		reflect.TypeFor[LiuyaoOptionalPosition](): {
			Type: "integer", Minimum: jsonschema.Ptr(float64(0)), Maximum: jsonschema.Ptr(float64(6)),
		},
		reflect.TypeFor[LiuyaoYaoLabel](): {
			Type: "string", Enum: []any{"老阴", "少阳", "少阴", "老阳"},
		},
		reflect.TypeFor[ganzhi.Gan](): {
			Types: []string{"null", "string"}, Enum: []any{nil, "甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"},
		},
		reflect.TypeFor[ganzhi.Zhi](): {
			Type: "string", Enum: []any{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"},
		},
		reflect.TypeFor[ganzhi.Wuxing](): {
			Type: "string", Enum: []any{"木", "火", "土", "金", "水"},
		},
		reflect.TypeFor[ganzhi.WangShuai](): {
			Type: "string", Enum: []any{"旺", "相", "休", "囚", "死"},
		},
		reflect.TypeFor[ganzhi.ShiShen](): {
			Type: "string", Enum: []any{"比肩", "劫财", "食神", "伤官", "偏财", "正财", "七杀", "正官", "偏印", "正印"},
		},
		reflect.TypeFor[bazi.NatalRelationKind](): {
			Type: "string", Enum: []any{"zhi_liu_he", "liu_chong", "liu_hai", "liu_po", "an_he", "liu_xing", "gan_he", "gan_he_candidate", "gan_chong", "gan_chong_candidate", "san_he", "san_he_partial", "san_hui"},
		},
		reflect.TypeFor[bazi.PillarIndex](): {
			Type: "string", Enum: []any{"nian", "yue", "ri", "shi"},
		},
		reflect.TypeFor[qimen.JinhanSpirit](): {
			Type: "string", Enum: []any{"diyao", "liuze", "wufu", "tiancao", "difu", "fengbo", "leigong", "yushi", "fengyun", "tangfu", "guoyin", "tianguan"},
		},
		reflect.TypeFor[qimen.GongIndex](): {
			Type: "string", Enum: []any{"坎", "坤", "震", "巽", "中", "乾", "兑", "艮", "离"},
		},
		reflect.TypeFor[qimen.DoorIndex](): {
			Type: "string", Enum: []any{"中门", "休门", "生门", "伤门", "杜门", "景门", "死门", "惊门", "开门"},
		},
		reflect.TypeFor[qimen.StarIndex](): {
			Type: "string", Enum: []any{"天蓬", "天芮", "天冲", "天辅", "天禽", "天心", "天柱", "天任", "天英"},
		},
		reflect.TypeFor[qimen.DunSource](): {
			Type: "string", Enum: []any{"hour_branch", "solar_term"},
		},
		reflect.TypeFor[qimen.HourBoundary](): {
			Type: "string", Enum: []any{"late_zi", "zi_zheng"},
		},
		reflect.TypeFor[qimen.PillarLevel](): {
			Type: "string", Enum: []any{"year", "month", "day", "hour"},
		},
		reflect.TypeFor[qimen.ElementRelation](): {
			Type: "string", Enum: []any{"same", "generates", "controls", "generated_by", "controlled_by"},
		},
		reflect.TypeFor[qimen.YingQiMatchKind](): {
			Type: "string", Enum: []any{"same", "opposite"},
		},
		reflect.TypeFor[GanName](): {
			Type: "string", Enum: []any{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"},
		},
		reflect.TypeFor[ZhiName](): {
			Type: "string", Enum: []any{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"},
		},
		reflect.TypeFor[WangShuaiName](): {
			Type: "string", Enum: []any{"旺", "相", "休", "囚", "死"},
		},
		reflect.TypeFor[LiuyaoYongShenName](): {
			Type: "string", Enum: []any{"父母", "兄弟", "官鬼", "妻财", "子孙", "世爻", "应爻"},
		},
		reflect.TypeFor[QimenYongShenSymbols](): {
			Type:        "array",
			Items:       &jsonschema.Schema{Type: "string"},
			MinItems:    jsonschema.Ptr(1),
			UniqueItems: true,
		},
	}
}

// integerConstSchema expresses an invariant as JSON Schema const, not as a
// one-value enum. The latter is for real alternatives and caused canonical
// fast-validator numeric mismatch.
func integerConstSchema(value int, description string) *jsonschema.Schema {
	constValue := any(value)
	return &jsonschema.Schema{
		Type:        "integer",
		Const:       &constValue,
		Description: description,
	}
}

var (
	schemaOptionsOnce sync.Once
	schemaOptions     *jsonschema.ForOptions
	schemaOptionsErr  error
)

func cachedSchemaOptions() (*jsonschema.ForOptions, error) {
	schemaOptionsOnce.Do(func() {
		if err := liuyao.Load(); err != nil {
			schemaOptionsErr = fmt.Errorf("load liuyao tables for MCP contracts: %w", err)
			return
		}
		schemaOptions = buildSchemaOptions(liuyao.HexagramNames())
	})
	return schemaOptions, schemaOptionsErr
}

func buildSchemaOptions(liuyaoHexagrams []string) *jsonschema.ForOptions {
	options := &jsonschema.ForOptions{TypeSchemas: buildStaticSchemaOptions()}
	registerWireEnums(options)
	registerSixItemArray[LiuyaoCoinRoundInput, LiuyaoCoinRoundsInput](options)
	registerSixItemArray[LiuyaoYaoValue, LiuyaoYaoValuesInput](options)
	registerSixItemArray[LiuyaoCoinRoundOutput, LiuyaoCoinRoundsOutput](options)
	registerSixItemArray[LiuyaoYaoValue, LiuyaoYaoValuesOutput](options)
	daXianItem, err := jsonschema.For[ZiweiDaXianItem](options)
	if err != nil {
		panic(fmt.Sprintf("derive ZiweiDaXianItem schema: %v", err))
	}
	options.TypeSchemas[reflect.TypeFor[ZiweiDaXianItems]()] = &jsonschema.Schema{
		Type:     "array",
		Items:    daXianItem,
		MinItems: jsonschema.Ptr(12),
		MaxItems: jsonschema.Ptr(12),
	}
	hexagramEnum := make([]any, 0, len(liuyaoHexagrams))
	for _, name := range liuyaoHexagrams {
		hexagramEnum = append(hexagramEnum, name)
	}
	options.TypeSchemas[reflect.TypeFor[LiuyaoHexagramName]()] = &jsonschema.Schema{
		Type: "string", Enum: hexagramEnum,
	}
	return options
}

// registerSixItemArray registers a fixed-length six-element array schema for a
// DTO slice type, deriving the item contract from the element DTO. Go slices
// cannot express the length constraint on their own, and the item schema stays
// DTO-owned instead of being spelled out here.
func registerSixItemArray[Item, Slice any](options *jsonschema.ForOptions) {
	items, err := jsonschema.For[Item](options)
	if err != nil {
		panic(fmt.Sprintf("derive %T item schema: %v", *new(Item), err))
	}
	options.TypeSchemas[reflect.TypeFor[Slice]()] = &jsonschema.Schema{
		Type: "array", MinItems: jsonschema.Ptr(6), MaxItems: jsonschema.Ptr(6), Items: items,
	}
}
