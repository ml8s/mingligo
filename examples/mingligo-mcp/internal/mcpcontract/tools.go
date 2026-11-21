package mcpcontract

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/ml8s/mingligo/examples/mingligo-mcp/internal/observability"
	"github.com/ml8s/mingligo/pkg/bazi"
)

// registerTyped adapts a pure domain handler to the official SDK typed
// handler. The SDK owns schema validation and protocol errors; this adapter
// owns timing and structured-only content.
func registerTyped[In, Out any](
	server *mcp.Server,
	tool mcp.Tool,
	annotations *mcp.ToolAnnotations,
	handler func(context.Context, In) (Out, error),
) {
	if annotations != nil {
		tool.Annotations = annotations
	}
	mcp.AddTool(server, &tool, func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		input In,
	) (*mcp.CallToolResult, Out, error) {
		started := time.Now()
		output, err := handler(ctx, input)
		status := "success"
		if err != nil {
			status = "error"
		}
		defer func() {
			observability.ObserveTool(toolDomain(tool.Name), tool.Name, status, time.Since(started).Seconds())
		}()
		if err != nil {
			return nil, zeroOutput[Out](), err
		}
		// Existing clients consume structuredContent only. Keep the response
		// compact by supplying an explicit, empty content slice.
		return &mcp.CallToolResult{Content: []mcp.Content{}}, output, nil
	})
}

// RegisterTools registers every MingliGo fact tool with the official Go MCP
// SDK. Registration is typed: wire DTOs define the contract, SDK validates
// input and output, and domain handlers never touch transport details.
func RegisterTools(server *mcp.Server) error {
	options, err := cachedSchemaOptions()
	if err != nil {
		return fmt.Errorf("build schema options: %w", err)
	}

	timeNowInput, timeNowOutput, err := typedSchemas[TimeNowInput, TimeNowOutput](options)
	if err != nil {
		return fmt.Errorf("time_now schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{Name: "time_now", Description: toolDescription("time_now"), InputSchema: timeNowInput, OutputSchema: timeNowOutput},
		toolAnnotations("time_now"), TimeNow)

	huangliInput, huangliOutput, err := typedSchemas[HuangliDaysInput, HuangliDaysOutput](options)
	if err != nil {
		return fmt.Errorf("huangli_days schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{Name: "huangli_days", Description: toolDescription("huangli_days"), InputSchema: huangliInput, OutputSchema: huangliOutput},
		toolAnnotations("huangli_days"), HuangliDays)

	subjectInput := mustSchema("subject_candidates input", func() (json.RawMessage, error) {
		return tianwenResolveFormsSchema(options)
	})
	subjectOutput, err := slimmedOutputSchema[TianwenSubjectCandidatesOutput](options)
	if err != nil {
		return fmt.Errorf("subject_candidates output schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{
			Name:         "subject_candidates",
			Description:  toolDescription("subject_candidates"),
			InputSchema:  subjectInput,
			OutputSchema: subjectOutput,
		},
		toolAnnotations("subject_candidates"), TianwenSubjectCandidates)

	cityInput, cityOutput, err := typedSchemas[CityInfoInput, CityInfoOutput](options)
	if err != nil {
		return fmt.Errorf("city_info schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{Name: "city_info", Description: toolDescription("city_info"), InputSchema: cityInput, OutputSchema: cityOutput},
		toolAnnotations("city_info"), CityInfo)

	baziChartInput, baziChartOutput, err := typedSchemas[BaziChartInput, bazi.ChartView](options)
	if err != nil {
		return fmt.Errorf("bazi_chart schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{
			Name:         "bazi_chart",
			Description:  toolDescription("bazi_chart"),
			InputSchema:  baziChartInput,
			OutputSchema: baziChartOutput,
		},
		toolAnnotations("bazi_chart"), BaziChart)

	bazhaiChartInput, bazhaiChartOutput, err := typedSchemas[BazhaiChartInput, BazhaiChartOutput](options)
	if err != nil {
		return fmt.Errorf("bazhai_chart schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{Name: "bazhai_chart", Description: toolDescription("bazhai_chart"), InputSchema: bazhaiChartInput, OutputSchema: bazhaiChartOutput},
		toolAnnotations("bazhai_chart"), BazhaiChart)

	bazhaiLayoutInput, bazhaiLayoutOutput, err := typedSchemas[BazhaiLayoutInput, BazhaiLayoutOutput](options)
	if err != nil {
		return fmt.Errorf("bazhai_layout schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{Name: "bazhai_layout", Description: toolDescription("bazhai_layout"), InputSchema: bazhaiLayoutInput, OutputSchema: bazhaiLayoutOutput},
		toolAnnotations("bazhai_layout"), BazhaiLayout)

	xuankongChartInput, xuankongChartOutput, err := typedSchemas[XuankongChartInput, XuankongChart](options)
	if err != nil {
		return fmt.Errorf("xuankong_chart schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{
			Name:         "xuankong_chart",
			Description:  toolDescription("xuankong_chart"),
			InputSchema:  xuankongChartInput,
			OutputSchema: xuankongChartOutput,
		},
		toolAnnotations("xuankong_chart"), XuankongChartCompute)

	xuankongLiunianInputRaw, err := inputSchemaFor[XuankongLiunianInput](options)
	if err != nil {
		return fmt.Errorf("xuankong_liunian input schema: %w", err)
	}
	xuankongLiunianInput, err := slimSchemaRaw(xuankongLiunianInputRaw)
	if err != nil {
		return fmt.Errorf("xuankong_liunian input schema: %w", err)
	}
	xuankongLiunianOutput := mustSchema("xuankong_liunian output", func() (json.RawMessage, error) {
		return seriesOutputSchema[XuankongLiunianOutput](options, "liu_nian")
	})
	registerTyped(server,
		mcp.Tool{
			Name:         "xuankong_liunian",
			Description:  toolDescription("xuankong_liunian"),
			InputSchema:  xuankongLiunianInput,
			OutputSchema: xuankongLiunianOutput,
		},
		toolAnnotations("xuankong_liunian"), XuankongLiunian)

	liuyaoQiguaInput := mustSchema("liuyao_qigua input", func() (json.RawMessage, error) {
		return liuyaoQiguaFormsSchema(options)
	})
	liuyaoQiguaOutput := mustSchema("liuyao_qigua output", func() (json.RawMessage, error) {
		return liuyaoQiguaOutputSchema(options)
	})
	registerTyped(server,
		mcp.Tool{
			Name:         "liuyao_qigua",
			Description:  toolDescription("liuyao_qigua"),
			InputSchema:  liuyaoQiguaInput,
			OutputSchema: liuyaoQiguaOutput,
		},
		toolAnnotations("liuyao_qigua"), LiuyaoQigua)

	liuyaoChartInput := mustSchema("liuyao_chart input", func() (json.RawMessage, error) {
		return liuyaoChartSchema(options)
	})
	liuyaoChartOutput := mustSchema("liuyao_chart output", func() (json.RawMessage, error) {
		return liuyaoChartOutputSchema(options)
	})
	registerTyped(server,
		mcp.Tool{
			Name:         "liuyao_chart",
			Description:  toolDescription("liuyao_chart"),
			InputSchema:  liuyaoChartInput,
			OutputSchema: liuyaoChartOutput,
		},
		toolAnnotations("liuyao_chart"), LiuyaoChart)

	baziBondInput, baziBondOutput, err := typedSchemas[BaziBondInput, BaziBondOutput](options)
	if err != nil {
		return fmt.Errorf("bazi_bond schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{Name: "bazi_bond", Description: toolDescription("bazi_bond"), InputSchema: baziBondInput, OutputSchema: baziBondOutput},
		toolAnnotations("bazi_bond"), BaziBond)

	baziLiupanInput := mustSchema("bazi_liupan input", func() (json.RawMessage, error) {
		return baziLiupanFormsSchema(options)
	})
	baziLiupanOutput := mustSchema("bazi_liupan output", func() (json.RawMessage, error) {
		return baziLiupanOutputSchema(options)
	})
	registerTyped(server,
		mcp.Tool{
			Name:         "bazi_liupan",
			Description:  toolDescription("bazi_liupan"),
			InputSchema:  baziLiupanInput,
			OutputSchema: baziLiupanOutput,
		},
		toolAnnotations("bazi_liupan"), BaziLiupan)

	ziweiChartInput, ziweiChartOutput, err := typedSchemas[ZiweiChartInput, ZiweiChartOutput](options)
	if err != nil {
		return fmt.Errorf("ziwei_chart schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{
			Name:         "ziwei_chart",
			Description:  toolDescription("ziwei_chart"),
			InputSchema:  ziweiChartInput,
			OutputSchema: ziweiChartOutput,
		},
		toolAnnotations("ziwei_chart"), ZiweiChart)

	ziweiLiupanInput := mustSchema("ziwei_liupan input", func() (json.RawMessage, error) {
		return ziweiLiupanFormsSchema(options)
	})
	ziweiLiupanOutput := mustSchema("ziwei_liupan output", func() (json.RawMessage, error) {
		return ziweiLiupanOutputSchema(options)
	})
	registerTyped(server,
		mcp.Tool{
			Name:         "ziwei_liupan",
			Description:  toolDescription("ziwei_liupan"),
			InputSchema:  ziweiLiupanInput,
			OutputSchema: ziweiLiupanOutput,
		},
		toolAnnotations("ziwei_liupan"), ZiweiLiupan)

	ziweiBondInput, ziweiBondOutput, err := typedSchemas[ZiweiBondInput, ZiweiBondOutput](options)
	if err != nil {
		return fmt.Errorf("ziwei_bond schema: %w", err)
	}
	registerTyped(server,
		mcp.Tool{Name: "ziwei_bond", Description: toolDescription("ziwei_bond"), InputSchema: ziweiBondInput, OutputSchema: ziweiBondOutput},
		toolAnnotations("ziwei_bond"), ZiweiBond)

	qimenInput := mustRawSchema("qimen_chart input", func() (json.RawMessage, error) {
		return qimenChartFormsSchema(options)
	})
	qimenOutput := mustRawSchema("qimen_chart output", func() (json.RawMessage, error) {
		return qimenChartOutputSchema(options)
	})
	registerTyped(server,
		mcp.Tool{
			Name:         "qimen_chart",
			Description:  toolDescription("qimen_chart"),
			InputSchema:  qimenInput,
			OutputSchema: qimenOutput,
		},
		toolAnnotations("qimen_chart"), QimenChart)

	return nil
}

// typedSchemas derives and slims both public schemas from wire DTOs.
func typedSchemas[In, Out any](options *jsonschema.ForOptions) (input, output json.RawMessage, err error) {
	inputRaw, err := inputSchemaFor[In](options)
	if err != nil {
		return nil, nil, err
	}
	input, err = slimSchemaRaw(inputRaw)
	if err != nil {
		return nil, nil, err
	}
	outputSchema, err := outputSchemaFor[Out](options)
	if err != nil {
		return nil, nil, err
	}
	output, err = slimSchemaRaw(outputSchema)
	if err != nil {
		return nil, nil, err
	}
	return input, output, nil
}

// slimmedOutputSchema explicitly publishes DTO-derived output schemas even
// when a custom input union requires separate construction.
func slimmedOutputSchema[Out any](options *jsonschema.ForOptions) (json.RawMessage, error) {
	raw, err := outputSchemaFor[Out](options)
	if err != nil {
		return nil, err
	}
	return slimSchemaRaw(raw)
}

func mustSchema(name string, build func() (json.RawMessage, error)) json.RawMessage {
	raw, err := build()
	if err != nil {
		panic(fmt.Sprintf("%s schema: %v", name, err))
	}
	raw, err = slimSchemaRaw(raw)
	if err != nil {
		panic(fmt.Sprintf("%s schema: %v", name, err))
	}
	return raw
}

// mustRawSchema publishes composition roots verbatim. Slimming after
// composition can rewrite or break references, so qimen keeps its explicit
// oneOf roots untouched.
func mustRawSchema(name string, build func() (json.RawMessage, error)) json.RawMessage {
	raw, err := build()
	if err != nil {
		panic(fmt.Sprintf("%s schema: %v", name, err))
	}
	return raw
}

func toolDescription(name string) string {
	if description, ok := toolDescriptions[name]; ok {
		return description
	}
	panic(fmt.Sprintf("missing MingliGo tool description: %s", name))
}

func toolAnnotations(name string) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		ReadOnlyHint:   true,
		IdempotentHint: name != "liuyao_qigua",
	}
}

func toolDomain(name string) string {
	if domain, _, found := strings.Cut(name, "_"); found {
		return domain
	}
	return "other"
}

func zeroOutput[Out any]() Out {
	var zero Out
	return zero
}
