package apicontract

import (
	"context"
	"fmt"
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
	Mode   LiuyaoQiguaMode       `json:"fang_shi" jsonschema:"起卦模式：auto 安全随机，coins 传六组铜钱，yaos 传六个爻值。"`
	Rounds LiuyaoCoinRoundsInput `json:"tong_qian_ji_lu,omitempty" jsonschema:"coins 模式：六组三枚铜钱，初爻到上爻。"`
	Yaos   LiuyaoYaoValuesInput  `json:"yao_zhi,omitempty" jsonschema:"yaos 模式：六个爻值，初爻到上爻。"`
}

type LiuyaoQiguaInput struct {
	Origin         LiuyaoCastingRequestInput `json:"origin" jsonschema:"起卦创建方式；只用于创建一次卦象事件。"`
	CastingContext LiuyaoCastingContextInput `json:"qi_gua_qing_jing" jsonschema:"起卦地标准时间、经度与用户确认的占问焦点。"`
}

type LiuyaoAutoCastingRequest struct {
	Mode LiuyaoQiguaMode `json:"fang_shi" jsonschema:"起卦模式；auto 使用安全随机。"`
}

type LiuyaoCoinsCastingRequest struct {
	Mode   LiuyaoQiguaMode       `json:"fang_shi" jsonschema:"起卦模式；coins 要求六组铜钱记录。"`
	Rounds LiuyaoCoinRoundsInput `json:"tong_qian_ji_lu" jsonschema:"六组三枚铜钱，初爻到上爻。"`
}

type LiuyaoYaosCastingRequest struct {
	Mode LiuyaoQiguaMode      `json:"fang_shi" jsonschema:"起卦模式；yaos 要求六个爻值。"`
	Yaos LiuyaoYaoValuesInput `json:"yao_zhi" jsonschema:"六个爻值，初爻到上爻。"`
}

type LiuyaoQiguaAutoForm struct {
	Origin         LiuyaoAutoCastingRequest  `json:"origin" jsonschema:"起卦创建方式。"`
	CastingContext LiuyaoCastingContextInput `json:"qi_gua_qing_jing" jsonschema:"起卦地标准时间、经度与用户确认的占问焦点。"`
}

type LiuyaoQiguaCoinsForm struct {
	Origin         LiuyaoCoinsCastingRequest `json:"origin" jsonschema:"起卦创建方式。"`
	CastingContext LiuyaoCastingContextInput `json:"qi_gua_qing_jing" jsonschema:"起卦地标准时间、经度与用户确认的占问焦点。"`
}

type LiuyaoQiguaYaosForm struct {
	Origin         LiuyaoYaosCastingRequest  `json:"origin" jsonschema:"起卦创建方式。"`
	CastingContext LiuyaoCastingContextInput `json:"qi_gua_qing_jing" jsonschema:"起卦地标准时间、经度与用户确认的占问焦点。"`
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
	Mode           LiuyaoCastMode         `json:"fang_shi" jsonschema:"实际使用的起卦模式。"`
	Order          LiuyaoOrderName        `json:"yao_xu" jsonschema:"爻序；当前固定为初爻到上爻。"`
	CoinConvention *LiuyaoCoinConvention  `json:"tong_qian_yue_ding,omitempty" jsonschema:"coins 模式的铜钱记分约定。"`
	Rounds         LiuyaoCoinRoundsOutput `json:"tong_qian_ji_lu,omitempty" jsonschema:"coins 模式的六组铜钱记录。"`
	Yaos           LiuyaoYaoValuesOutput  `json:"yao_zhi" jsonschema:"六个爻值，初爻到上爻。"`
	DongYao        LiuyaoPositions        `json:"dong_yao" jsonschema:"动爻位置列表，1-6。"`
}

type LiuyaoCoinsCastingReceipt struct {
	Mode           LiuyaoCastMode         `json:"fang_shi" jsonschema:"起卦方式。"`
	Order          LiuyaoOrderName        `json:"yao_xu" jsonschema:"爻序。"`
	CoinConvention *LiuyaoCoinConvention  `json:"tong_qian_yue_ding,omitempty" jsonschema:"铜钱记分约定。"`
	Rounds         LiuyaoCoinRoundsOutput `json:"tong_qian_ji_lu" jsonschema:"六组铜钱记录。"`
	Yaos           LiuyaoYaoValuesOutput  `json:"yao_zhi" jsonschema:"六个爻值。"`
	DongYao        LiuyaoPositions        `json:"dong_yao" jsonschema:"动爻位置列表。"`
}

type LiuyaoYaosCastingReceipt struct {
	Mode    LiuyaoCastMode        `json:"fang_shi" jsonschema:"起卦方式。"`
	Order   LiuyaoOrderName       `json:"yao_xu" jsonschema:"爻序。"`
	Yaos    LiuyaoYaoValuesOutput `json:"yao_zhi" jsonschema:"六个爻值。"`
	DongYao LiuyaoPositions       `json:"dong_yao" jsonschema:"动爻位置列表。"`
}

type LiuyaoMatterFocusInput struct {
	Kind        string `json:"kind" jsonschema:"焦点类型；matter 表示按问事类别推导。"`
	Matter      string `json:"matter" jsonschema:"问事类别；relationship 需同时提供 perspective。"`
	Perspective string `json:"perspective,omitempty" jsonschema:"仅 relationship 使用：问卦人性别，male 男占女，female 女占男。"`
}

type LiuyaoYongShenFocusInput struct {
	Kind   string             `json:"kind" jsonschema:"焦点类型；qu_yong 表示显式取用类别。"`
	QuYong LiuyaoYongShenName `json:"qu_yong" jsonschema:"高级显式取用类别。"`
}

type LiuyaoFocusInput struct {
	Kind        string `json:"kind"`
	Matter      string `json:"matter,omitempty"`
	Perspective string `json:"perspective,omitempty"`
	QuYong      string `json:"qu_yong,omitempty"`
}

type LiuyaoCastingContextInput struct {
	StandardTime StandardTime     `json:"shi_jian" jsonschema:"起卦地标准时间（无 DST），RFC3339 带 offset。"`
	Longitude    Longitude        `json:"jing_du" jsonschema:"起卦地经度（-180..180）；真太阳时校正用于时辰。"`
	Focus        LiuyaoFocusInput `json:"focus" jsonschema:"占问焦点；起卦时确定性解析为用神。"`
}

// LiuyaoCastingContext is the immutable event context. It carries only the
// resolved yong_shen so liuyao_chart never re-interpret the user-facing focus.
type LiuyaoCastingContext = LiuyaoCastingContextEventInput

type LiuyaoCastingEvent struct {
	Casting        LiuyaoCastingOutput  `json:"qi_gua_ji_lu" jsonschema:"一次起卦的完整源事实。"`
	CastingContext LiuyaoCastingContext `json:"qi_gua_qing_jing" jsonschema:"起卦地标准时间、经度与解析出的用神。"`
}

// LiuyaoCastingEventInput is the canonical minimal event created by qigua.
type LiuyaoCastingEventInput struct {
	Casting        LiuyaoCastingOutput            `json:"qi_gua_ji_lu" jsonschema:"一次起卦的完整源事实。"`
	CastingContext LiuyaoCastingContextEventInput `json:"qi_gua_qing_jing" jsonschema:"起卦地标准时间、经度与解析出的用神。"`
}

type LiuyaoCastingContextEventInput struct {
	StandardTime StandardTime       `json:"shi_jian" jsonschema:"起卦地标准时间（无 DST），RFC3339 带 offset。"`
	Longitude    Longitude          `json:"jing_du" jsonschema:"起卦地经度（-180..180）；真太阳时校正用于时辰。"`
	QuYong       LiuyaoYongShenName `json:"qu_yong" jsonschema:"起卦时确定的取用类别；成卦后的用神爻见 chart.yong_shen.gua_ming。"`
}

// perspective is only valid on the relationship branch.

type LiuyaoChartResultOutput struct {
	Event LiuyaoCastingEvent `json:"qi_gua" jsonschema:"本次排盘原样回传的不可变卦象事件。"`
	Chart liuyao.ChartView   `json:"chart" jsonschema:"六爻盘面统一视图。"`
}

// liuyaoResolveFocus only validates the wire shape and delegates the命理
// matter→用神 rule to the domain. The MCP layer owns no命理判断.
func liuyaoResolveFocus(focus LiuyaoFocusInput) (LiuyaoYongShenName, error) {
	switch focus.Kind {
	case "matter":
		if focus.QuYong != "" {
			return "", fmt.Errorf("matter focus does not accept qu_yong")
		}
		yongShen, err := liuyao.YongShenFromMatter(focus.Matter, focus.Perspective)
		if err != nil {
			return "", fmt.Errorf("matter focus: %w", err)
		}
		return LiuyaoYongShenName(liuyao.DecodeYongShen(yongShen)), nil
	case "qu_yong":
		if focus.Matter != "" || focus.Perspective != "" {
			return "", fmt.Errorf("qu_yong focus does not accept matter or perspective")
		}
		if focus.QuYong == "" {
			return "", fmt.Errorf("qu_yong focus requires qu_yong")
		}
		if _, err := liuyao.EncodeYongShen(focus.QuYong); err != nil {
			return "", fmt.Errorf("invalid qu_yong focus: %w", err)
		}
		return LiuyaoYongShenName(focus.QuYong), nil
	default:
		return "", fmt.Errorf("focus.kind must be matter/qu_yong")
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
		QuYong:       string(yongShen),
	}
	return result, nil
}

// LiuyaoChartInput consumes the immutable minimal event created by liuyao_qigua.
type LiuyaoChartInput struct {
	Event *LiuyaoCastingEventInput `json:"qi_gua" jsonschema:"liuyao_qigua 返回的不可变卦象事件。"`
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
	ys, err := liuyao.EncodeYongShen(string(input.Event.CastingContext.QuYong))
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
