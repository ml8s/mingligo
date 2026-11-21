package mcpcontract

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"

	"github.com/ml8s/mingligo/pkg/liuyao"
)

// LiuyaoCoinRoundInput is one raw three-coin outcome. The domain owns totaling,
// labels, changing flags, and the normalized receipt.
type LiuyaoCoinRoundInput [3]LiuyaoCoinSymbol

type LiuyaoCoinRoundsInput []LiuyaoCoinRoundInput

type LiuyaoYaoValuesInput []LiuyaoYaoValue

// LiuyaoCastingRequestInput is the domain-owned source fact for one casting.
// It is the request form of LiuyaoCastingOutput.
type LiuyaoCastingRequestInput struct {
	Mode   LiuyaoQiguaMode       `json:"mode" jsonschema:"起卦模式：auto 安全随机，coins 传六组铜钱，yaos 传六个爻值。"`
	Rounds LiuyaoCoinRoundsInput `json:"rounds,omitempty" jsonschema:"coins 模式：六组三枚铜钱，初爻到上爻。"`
	Yaos   LiuyaoYaoValuesInput  `json:"yaos,omitempty" jsonschema:"yaos 模式：六个爻值，初爻到上爻。"`
}

type LiuyaoQiguaInput struct {
	Origin         LiuyaoCastingRequestInput `json:"origin" jsonschema:"起卦创建方式；只用于创建一次卦象事件。"`
	CastingContext LiuyaoCastingContextInput `json:"casting_context" jsonschema:"起卦地标准时间、经度与用户确认的占问焦点。"`
}

type LiuyaoAutoCastingRequest struct {
	Mode LiuyaoQiguaMode `json:"mode" jsonschema:"起卦模式；auto 使用安全随机。"`
}

type LiuyaoCoinsCastingRequest struct {
	Mode   LiuyaoQiguaMode       `json:"mode" jsonschema:"起卦模式；coins 要求六组铜钱记录。"`
	Rounds LiuyaoCoinRoundsInput `json:"rounds" jsonschema:"六组三枚铜钱，初爻到上爻。"`
}

type LiuyaoYaosCastingRequest struct {
	Mode LiuyaoQiguaMode      `json:"mode" jsonschema:"起卦模式；yaos 要求六个爻值。"`
	Yaos LiuyaoYaoValuesInput `json:"yaos" jsonschema:"六个爻值，初爻到上爻。"`
}

type LiuyaoQiguaAutoForm struct {
	Origin         LiuyaoAutoCastingRequest  `json:"origin" jsonschema:"起卦创建方式。"`
	CastingContext LiuyaoCastingContextInput `json:"casting_context" jsonschema:"起卦地标准时间、经度与用户确认的占问焦点。"`
}

type LiuyaoQiguaCoinsForm struct {
	Origin         LiuyaoCoinsCastingRequest `json:"origin" jsonschema:"起卦创建方式。"`
	CastingContext LiuyaoCastingContextInput `json:"casting_context" jsonschema:"起卦地标准时间、经度与用户确认的占问焦点。"`
}

type LiuyaoQiguaYaosForm struct {
	Origin         LiuyaoYaosCastingRequest  `json:"origin" jsonschema:"起卦创建方式。"`
	CastingContext LiuyaoCastingContextInput `json:"casting_context" jsonschema:"起卦地标准时间、经度与用户确认的占问焦点。"`
}

func liuyaoQiguaFormsSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	modes := []string{"auto", "coins", "yaos"}
	forms := make([]json.RawMessage, 0, len(modes))
	for _, mode := range modes {
		var raw json.RawMessage
		var err error
		nested, err := liuyaoOriginBranchSchema(options, mode)
		if err != nil {
			return nil, err
		}
		switch mode {
		case "auto":
			raw, err = inputSchemaFor[LiuyaoQiguaAutoForm](options)
		case "coins":
			raw, err = inputSchemaFor[LiuyaoQiguaCoinsForm](options)
		case "yaos":
			raw, err = inputSchemaFor[LiuyaoQiguaYaosForm](options)
		}
		if err != nil {
			return nil, err
		}
		raw, err = withSchemaProperty(raw, "origin", nested)
		if err != nil {
			return nil, err
		}
		castingContext, err := liuyaoCastingContextSchema(options)
		if err != nil {
			return nil, err
		}
		raw, err = withSchemaProperty(raw, "casting_context", castingContext)
		if err != nil {
			return nil, err
		}
		forms = append(forms, raw)
	}
	return json.Marshal(map[string]any{
		"type":  "object",
		"oneOf": forms,
	})
}

type LiuyaoHeadsScore int
type LiuyaoTailsScore int

type LiuyaoCoinConvention struct {
	Heads LiuyaoHeadsScore `json:"正" jsonschema:"正面记 3。"`
	Tails LiuyaoTailsScore `json:"反" jsonschema:"反面记 2。"`
}

type LiuyaoCoinRoundsOutput []LiuyaoCoinRoundOutput

type LiuyaoCoinRoundOutput struct {
	Position LiuyaoPosition      `json:"position" jsonschema:"爻位，初爻到上爻 1-6。"`
	Coins    [3]LiuyaoCoinSymbol `json:"coins" jsonschema:"三枚铜钱记录，正/反按调用顺序保留。"`
	Value    LiuyaoYaoValue      `json:"yao_value" jsonschema:"三币合计爻值 6-9。"`
	Label    LiuyaoYaoLabel      `json:"yao_label" jsonschema:"老阴/少阳/少阴/老阳。"`
	Changing bool                `json:"is_dong_yao" jsonschema:"是否动爻。"`
}

type LiuyaoYaoValuesOutput []LiuyaoYaoValue
type LiuyaoPositions []LiuyaoPosition

type LiuyaoCastingOutput struct {
	Mode           LiuyaoCastMode         `json:"mode" jsonschema:"实际使用的起卦模式。"`
	Order          LiuyaoOrderName        `json:"order" jsonschema:"爻序；当前固定为初爻到上爻。"`
	CoinConvention *LiuyaoCoinConvention  `json:"coin_convention,omitempty" jsonschema:"coins 模式的铜钱记分约定。"`
	Rounds         LiuyaoCoinRoundsOutput `json:"rounds,omitempty" jsonschema:"coins 模式的六组铜钱记录。"`
	Yaos           LiuyaoYaoValuesOutput  `json:"yaos" jsonschema:"六个爻值，初爻到上爻。"`
	DongYao        LiuyaoPositions        `json:"dong_yao" jsonschema:"动爻位置列表，1-6。"`
}

type LiuyaoCoinsCastingReceipt struct {
	Mode           LiuyaoCastMode         `json:"mode" jsonschema:"起卦方式。"`
	Order          LiuyaoOrderName        `json:"order" jsonschema:"爻序。"`
	CoinConvention *LiuyaoCoinConvention  `json:"coin_convention,omitempty" jsonschema:"铜钱记分约定。"`
	Rounds         LiuyaoCoinRoundsOutput `json:"rounds" jsonschema:"六组铜钱记录。"`
	Yaos           LiuyaoYaoValuesOutput  `json:"yaos" jsonschema:"六个爻值。"`
	DongYao        LiuyaoPositions        `json:"dong_yao" jsonschema:"动爻位置列表。"`
}

type LiuyaoYaosCastingReceipt struct {
	Mode    LiuyaoCastMode        `json:"mode" jsonschema:"起卦方式。"`
	Order   LiuyaoOrderName       `json:"order" jsonschema:"爻序。"`
	Yaos    LiuyaoYaoValuesOutput `json:"yaos" jsonschema:"六个爻值。"`
	DongYao LiuyaoPositions       `json:"dong_yao" jsonschema:"动爻位置列表。"`
}

type LiuyaoMatterFocusInput struct {
	Kind        string `json:"kind" jsonschema:"焦点类型；matter 表示按问事类别推导。"`
	Matter      string `json:"matter" jsonschema:"问事类别；relationship 需同时提供 perspective。"`
	Perspective string `json:"perspective,omitempty" jsonschema:"仅 relationship 使用：问卦人性别，male 男占女，female 女占男。"`
}

type LiuyaoYongShenFocusInput struct {
	Kind     string             `json:"kind" jsonschema:"焦点类型；yong_shen 表示显式用神。"`
	YongShen LiuyaoYongShenName `json:"yong_shen" jsonschema:"高级显式用神。"`
}

type LiuyaoFocusInput struct {
	Kind        string `json:"kind"`
	Matter      string `json:"matter,omitempty"`
	Perspective string `json:"perspective,omitempty"`
	YongShen    string `json:"yong_shen,omitempty"`
}

type LiuyaoCastingContextInput struct {
	StandardTime StandardTime     `json:"standard_time" jsonschema:"起卦地标准时间（无 DST），RFC3339 带 offset。"`
	Longitude    Longitude        `json:"longitude" jsonschema:"起卦地经度（-180..180）；真太阳时校正用于时辰。"`
	Focus        LiuyaoFocusInput `json:"focus" jsonschema:"占问焦点；起卦时确定性解析为用神。"`
}

// LiuyaoCastingContext is the immutable event context. It carries only the
// resolved yong_shen so chart/ask never re-interpret the user-facing focus.
type LiuyaoCastingContext = LiuyaoCastingContextEventInput

type LiuyaoCastingEvent struct {
	Casting        LiuyaoCastingOutput  `json:"casting" jsonschema:"一次起卦的完整源事实。"`
	CastingContext LiuyaoCastingContext `json:"casting_context" jsonschema:"起卦地标准时间、经度与解析出的用神。"`
}

// LiuyaoCastingEventInput is the canonical minimal event created by qigua.
type LiuyaoCastingEventInput struct {
	Casting        LiuyaoCastingOutput            `json:"casting" jsonschema:"一次起卦的完整源事实。"`
	CastingContext LiuyaoCastingContextEventInput `json:"casting_context" jsonschema:"起卦地标准时间、经度与解析出的用神。"`
}

type LiuyaoCastingContextEventInput struct {
	StandardTime StandardTime       `json:"standard_time" jsonschema:"起卦地标准时间（无 DST），RFC3339 带 offset。"`
	Longitude    Longitude          `json:"longitude" jsonschema:"起卦地经度（-180..180）；真太阳时校正用于时辰。"`
	YongShen     LiuyaoYongShenName `json:"yong_shen" jsonschema:"起卦时解析出的用神；chart/ask 唯一取用依据。"`
}

func liuyaoMatterFocusSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := inputSchemaFor[LiuyaoMatterFocusInput](options)
	if err != nil {
		return nil, err
	}
	raw, err = withSchemaConst(raw, "kind", "matter")
	if err != nil {
		return nil, err
	}
	raw, err = withSchemaEnum(
		raw, "matter",
		"self", "other", "parents", "children", "siblings",
		"wealth", "career", "study", "health", "legal", "home", "travel", "lost_item",
	)
	if err != nil {
		return nil, err
	}
	// perspective is only valid on the relationship branch.
	return withoutSchemaProperty(raw, "perspective")
}

func liuyaoRelationshipFocusSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := inputSchemaFor[LiuyaoMatterFocusInput](options)
	if err != nil {
		return nil, err
	}
	raw, err = withSchemaConst(raw, "kind", "matter")
	if err != nil {
		return nil, err
	}
	raw, err = withSchemaConst(raw, "matter", "relationship")
	if err != nil {
		return nil, err
	}
	raw, err = withSchemaEnum(raw, "perspective", "male", "female")
	if err != nil {
		return nil, err
	}
	return withSchemaRequiredFields(raw, "kind", "matter", "perspective")
}

func liuyaoYongShenFocusSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := inputSchemaFor[LiuyaoYongShenFocusInput](options)
	if err != nil {
		return nil, err
	}
	return withSchemaConst(raw, "kind", "yong_shen")
}

func liuyaoOriginBranchSchema(options *jsonschema.ForOptions, mode string) (json.RawMessage, error) {
	var raw json.RawMessage
	var err error
	switch mode {
	case "auto":
		raw, err = inputSchemaFor[LiuyaoAutoCastingRequest](options)
	case "coins":
		raw, err = inputSchemaFor[LiuyaoCoinsCastingRequest](options)
	case "yaos":
		raw, err = inputSchemaFor[LiuyaoYaosCastingRequest](options)
	default:
		return nil, fmt.Errorf("unsupported liuyao origin mode %q", mode)
	}
	if err != nil {
		return nil, err
	}
	return withSchemaConst(raw, "mode", mode)
}

func patchLiuyaoEvent(event map[string]any, options *jsonschema.ForOptions) (map[string]any, error) {
	raw, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	casting, err := liuyaoCastingReceiptSchema(options)
	if err != nil {
		return nil, err
	}
	raw, err = withSchemaProperty(raw, "casting", casting)
	if err != nil {
		return nil, err
	}
	context, err := liuyaoEventCastingContextSchema(options)
	if err != nil {
		return nil, err
	}
	raw, err = withSchemaProperty(raw, "casting_context", context)
	if err != nil {
		return nil, err
	}
	var patched map[string]any
	if err := json.Unmarshal(raw, &patched); err != nil {
		return nil, err
	}
	return patched, nil
}

func patchLiuyaoEventInput(event map[string]any, options *jsonschema.ForOptions) (map[string]any, error) {
	raw, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	casting, err := liuyaoCastingReceiptSchema(options)
	if err != nil {
		return nil, err
	}
	raw, err = withSchemaProperty(raw, "casting", casting)
	if err != nil {
		return nil, err
	}
	context, err := inputSchemaFor[LiuyaoCastingContextEventInput](options)
	if err != nil {
		return nil, err
	}
	raw, err = withSchemaProperty(raw, "casting_context", context)
	if err != nil {
		return nil, err
	}
	var patched map[string]any
	if err := json.Unmarshal(raw, &patched); err != nil {
		return nil, err
	}
	return patched, nil
}

func withLiuyaoEventSchema(document map[string]any, options *jsonschema.ForOptions) (map[string]any, error) {
	properties, ok := document["properties"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema has no properties")
	}
	event, ok := properties["event"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema has no event properties")
	}
	patched, err := patchLiuyaoEvent(event, options)
	if err != nil {
		return nil, err
	}
	properties["event"] = patched
	return document, nil
}

func withLiuyaoEventInputSchema(document map[string]any, options *jsonschema.ForOptions) (map[string]any, error) {
	properties, ok := document["properties"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema has no properties")
	}
	event, ok := properties["event"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema has no event properties")
	}
	patched, err := patchLiuyaoEventInput(event, options)
	if err != nil {
		return nil, err
	}
	properties["event"] = patched
	return document, nil
}

func liuyaoFocusSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	matter, err := liuyaoMatterFocusSchema(options)
	if err != nil {
		return nil, err
	}
	relationship, err := liuyaoRelationshipFocusSchema(options)
	if err != nil {
		return nil, err
	}
	yongShen, err := liuyaoYongShenFocusSchema(options)
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]any{
		"oneOf": []json.RawMessage{matter, relationship, yongShen},
	})
}

func liuyaoCastingContextSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := inputSchemaFor[LiuyaoCastingContextInput](options)
	if err != nil {
		return nil, err
	}
	focus, err := liuyaoFocusSchema(options)
	if err != nil {
		return nil, err
	}
	return withSchemaProperty(raw, "focus", focus)
}

func liuyaoEventCastingContextSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := inputSchemaFor[LiuyaoCastingContext](options)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func liuyaoCastingReceiptSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	coins, err := inputSchemaFor[LiuyaoCoinsCastingReceipt](options)
	if err != nil {
		return nil, err
	}
	coins, err = withSchemaConst(coins, "mode", "coins")
	if err != nil {
		return nil, err
	}
	coins, err = withSchemaConst(coins, "order", "bottom_up")
	if err != nil {
		return nil, err
	}
	yaos, err := inputSchemaFor[LiuyaoYaosCastingReceipt](options)
	if err != nil {
		return nil, err
	}
	yaos, err = withSchemaConst(yaos, "mode", "yaos")
	if err != nil {
		return nil, err
	}
	yaos, err = withSchemaConst(yaos, "order", "bottom_up")
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]any{"oneOf": []json.RawMessage{coins, yaos}})
}

func liuyaoChartSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := inputSchemaFor[LiuyaoChartInput](options)
	if err != nil {
		return nil, err
	}
	var document map[string]any
	if err := json.Unmarshal(raw, &document); err != nil {
		return nil, err
	}
	document, err = withLiuyaoEventInputSchema(document, options)
	if err != nil {
		return nil, err
	}
	return json.Marshal(document)
}

func liuyaoStrictArraySchema[T any](options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := outputSchemaFor[T](options)
	if err != nil {
		return nil, err
	}
	var document map[string]any
	if err := json.Unmarshal(raw, &document); err != nil {
		return nil, err
	}
	strictArrays(document)
	if _, ok := document["oneOf"]; !ok {
		return json.Marshal(document)
	}
	branches, ok := document["oneOf"].([]any)
	if !ok {
		return nil, fmt.Errorf("schema oneOf is not an array")
	}
	for _, branch := range branches {
		strictArrays(branch)
	}
	return json.Marshal(document)
}

func liuyaoQiguaOutputSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := outputSchemaFor[liuyao.QiguaView](options)
	if err != nil {
		return nil, err
	}
	var document map[string]any
	if err := json.Unmarshal(raw, &document); err != nil {
		return nil, err
	}
	document, err = withLiuyaoEventSchema(document, options)
	if err != nil {
		return nil, err
	}
	strictArrays(document)
	return json.Marshal(document)
}

func liuyaoChartOutputSchema(options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := liuyaoStrictArraySchema[LiuyaoChartResultOutput](options)
	if err != nil {
		return nil, err
	}
	var document map[string]any
	if err := json.Unmarshal(raw, &document); err != nil {
		return nil, err
	}
	document, err = withLiuyaoEventSchema(document, options)
	if err != nil {
		return nil, err
	}
	return json.Marshal(document)
}

func strictArrays(value any) {
	switch node := value.(type) {
	case map[string]any:
		if types, ok := node["type"].([]any); ok && len(types) == 2 {
			hasNull, hasArray := false, false
			for _, kind := range types {
				if kind == "null" {
					hasNull = true
				}
				if kind == "array" {
					hasArray = true
				}
			}
			if hasNull && hasArray {
				node["type"] = "array"
			}
		}
		for _, child := range node {
			strictArrays(child)
		}
	case []any:
		for _, child := range node {
			strictArrays(child)
		}
	}
}

type LiuyaoChartResultOutput struct {
	Event LiuyaoCastingEvent `json:"event" jsonschema:"本次排盘原样回传的不可变卦象事件。"`
	Chart liuyao.ChartView   `json:"chart" jsonschema:"六爻盘面统一视图；analysis 一并置于 chart 内。"`
}

// liuyaoMatterYongShen is the single source of truth for matter -> yong_shen.
// rationale:
//   - self/health/travel: 自占以世爻为用；疾病以世为身，官鬼为病神；出行吉凶看世。
//   - other: 无尊卑称呼、未深交、泛指对象取应爻。
//   - relationship: 男占女以妻财为主，女占男以官鬼为主；由 perspective 决定。
//   - parents/home: 父母为墙城宅舍、文书。
//   - children/siblings: 六亲分类。
//   - wealth/lost_item: 财爻为钱物、所失物之主。
//   - career/legal: 官鬼为功名官府、问官。
//   - study: 父爻为学业文书，官鬼只作名次辅助。
//
// Legacy matters remain readable for historical events only.
var liuyaoMatterYongShen = map[string]LiuyaoYongShenName{
	"self":        "世爻",
	"other":       "应爻",
	"parents":     "父母",
	"children":    "子孙",
	"siblings":    "兄弟",
	"wealth":      "妻财",
	"career":      "官鬼",
	"study":       "父母",
	"health":      "世爻",
	"legal":       "官鬼",
	"home":        "父母",
	"travel":      "世爻",
	"lost_item":   "妻财",
	"general":     "世爻",
	"family":      "世爻",
	"competition": "兄弟",
}

func liuyaoResolveFocus(focus LiuyaoFocusInput) (LiuyaoYongShenName, error) {
	switch focus.Kind {
	case "matter":
		if focus.YongShen != "" {
			return "", fmt.Errorf("matter focus does not accept yong_shen")
		}
		if focus.Matter == "relationship" {
			if focus.Perspective != "male" && focus.Perspective != "female" {
				return "", fmt.Errorf("relationship focus requires perspective=male/female")
			}
			if focus.Perspective == "male" {
				return "妻财", nil
			} else {
				return "官鬼", nil
			}
		}
		if focus.Perspective != "" {
			return "", fmt.Errorf("perspective is only valid for relationship")
		}
		name, ok := liuyaoMatterYongShen[focus.Matter]
		if !ok {
			return "", fmt.Errorf("unsupported matter focus %q", focus.Matter)
		}
		return name, nil
	case "yong_shen":
		if focus.Matter != "" || focus.Perspective != "" {
			return "", fmt.Errorf("yong_shen focus does not accept matter or perspective")
		}
		if focus.YongShen == "" {
			return "", fmt.Errorf("yong_shen focus requires yong_shen")
		}
		if _, err := liuyao.EncodeYongShen(focus.YongShen); err != nil {
			return "", fmt.Errorf("invalid yong_shen focus: %w", err)
		}
		return LiuyaoYongShenName(focus.YongShen), nil
	default:
		return "", fmt.Errorf("focus.kind must be matter/yong_shen")
	}
}

func LiuyaoQigua(_ context.Context, input LiuyaoQiguaInput) (liuyao.QiguaView, error) {
	yongShen, err := liuyaoResolveFocus(input.CastingContext.Focus)
	if err != nil {
		return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: invalid focus: %w", err)
	}
	if _, _, err := parseStandardTime(input.CastingContext.StandardTime); err != nil {
		return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: %w", err)
	}
	if float64(input.CastingContext.Longitude) < -180 || float64(input.CastingContext.Longitude) > 180 {
		return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: longitude must be within [-180,180], got %v", float64(input.CastingContext.Longitude))
	}
	var casting liuyao.Casting
	request := input.Origin
	switch request.Mode {
	case "auto":
		if request.Rounds != nil || request.Yaos != nil {
			return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: auto mode does not accept rounds or yaos")
		}
		casting, err = liuyao.SecureQigua()
	case "coins":
		if request.Yaos != nil {
			return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: coins mode accepts rounds only")
		}
		if len(request.Rounds) != 6 {
			return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: casting.coins requires exactly 6 rounds, got %d", len(request.Rounds))
		}
		rounds := make([][]liuyao.CoinFace, 0, 6)
		for _, round := range request.Rounds {
			coins := make([]liuyao.CoinFace, 0, 3)
			for _, coin := range round {
				coins = append(coins, liuyao.CoinFace(coin))
			}
			rounds = append(rounds, coins)
		}
		casting, err = liuyao.NewCoinsCasting(rounds)
	case "yaos":
		if request.Rounds != nil {
			return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: yaos mode accepts yaos only")
		}
		if len(request.Yaos) != 6 {
			return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: casting.yaos requires exactly 6 values, got %d", len(request.Yaos))
		}
		var values [6]int
		for i, value := range request.Yaos {
			values[i] = int(value)
		}
		casting, err = liuyao.NewValuesCasting(values)
	default:
		return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: casting.mode unsupported %q", string(request.Mode))
	}
	if err != nil {
		return liuyao.QiguaView{}, fmt.Errorf("liuyao_qigua: %w", err)
	}
	var result liuyao.QiguaView
	result.Event.Casting = casting.QiguaCastingView()
	result.Event.CastingContext = liuyao.QiguaCastingContextView{
		StandardTime: string(input.CastingContext.StandardTime),
		Longitude:    float64(input.CastingContext.Longitude),
		YongShen:     string(yongShen),
	}
	return result, nil
}

// LiuyaoChartInput consumes the immutable minimal event created by liuyao_qigua.
type LiuyaoChartInput struct {
	Event *LiuyaoCastingEventInput `json:"event" jsonschema:"liuyao_qigua 返回的不可变卦象事件。"`
}

func LiuyaoChart(_ context.Context, input LiuyaoChartInput) (LiuyaoChartResultOutput, error) {
	if input.Event == nil {
		return LiuyaoChartResultOutput{}, fmt.Errorf("liuyao.chart: event is required")
	}
	casting, err := liuyaoCastingDomain(input.Event.Casting)
	if err != nil {
		return LiuyaoChartResultOutput{}, fmt.Errorf("liuyao.chart: invalid casting: %w", err)
	}
	standard, _, err := parseStandardTime(input.Event.CastingContext.StandardTime)
	if err != nil {
		return LiuyaoChartResultOutput{}, fmt.Errorf("liuyao.chart: %w", err)
	}
	if float64(input.Event.CastingContext.Longitude) < -180 || float64(input.Event.CastingContext.Longitude) > 180 {
		return LiuyaoChartResultOutput{}, fmt.Errorf("liuyao.chart: longitude must be within [-180,180], got %v", float64(input.Event.CastingContext.Longitude))
	}
	ys, err := liuyao.EncodeYongShen(string(input.Event.CastingContext.YongShen))
	if err != nil {
		return LiuyaoChartResultOutput{}, fmt.Errorf("liuyao.chart: %w", err)
	}
	chart := liuyao.ComputeChartAt(standard, float64(input.Event.CastingContext.Longitude), ys, casting.Yaos)
	view := chart.View()
	return LiuyaoChartResultOutput{
		Event: LiuyaoCastingEvent{
			Casting:        input.Event.Casting,
			CastingContext: input.Event.CastingContext,
		},
		Chart: view,
	}, nil
}

func liuyaoCastingDomain(output LiuyaoCastingOutput) (liuyao.Casting, error) {
	if len(output.Yaos) != 6 {
		return liuyao.Casting{}, fmt.Errorf("yaos requires exactly 6 values, got %d", len(output.Yaos))
	}
	casting := liuyao.Casting{
		Mode:    string(output.Mode),
		Order:   string(output.Order),
		Yaos:    [6]int{},
		DongYao: make([]int, 0, len(output.DongYao)),
	}
	for i, value := range output.Yaos {
		casting.Yaos[i] = int(value)
	}
	if output.CoinConvention != nil {
		casting.CoinConvention = &liuyao.CoinConvention{
			Heads: liuyao.CoinScore(output.CoinConvention.Heads),
			Tails: liuyao.CoinScore(output.CoinConvention.Tails),
		}
	}
	for _, round := range output.Rounds {
		coins := make([]liuyao.CoinFace, 0, 3)
		for _, coin := range round.Coins {
			coins = append(coins, liuyao.CoinFace(coin))
		}
		casting.Rounds = append(casting.Rounds, liuyao.CoinRound{
			Position: int(round.Position), Coins: coins, Value: liuyao.YaoType(round.Value),
			Label: string(round.Label), Changing: round.Changing,
		})
	}
	for _, position := range output.DongYao {
		casting.DongYao = append(casting.DongYao, int(position))
	}
	if err := casting.Validate(); err != nil {
		return liuyao.Casting{}, err
	}
	return casting, nil
}
