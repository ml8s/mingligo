package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/ml8s/mingligo/examples/mingligo-mcp/internal/mcpcontract"
)

// serverIdentity is runtime configuration, not a build-time brand. Neutral
// defaults keep Mingli images reusable; deployment injects hosted identities.
type serverIdentity struct {
	Name        string
	Title       string
	Description string
	WebsiteURL  string
}

func loadServerIdentity() serverIdentity {
	return serverIdentity{
		Name:        envOr("MINGLIGO_SERVER_NAME", "local/mingligo-mcp"),
		Title:       envOr("MINGLIGO_SERVER_TITLE", "MingliGo MCP"),
		Description: "MingliGo MCP — deterministic metaphysics computation and auditable divination events: bazi, ziwei, liuyao, qimen, huangli, fengshui.",
		WebsiteURL:  envOr("MINGLIGO_SERVER_URL", "https://github.com/ml8s"),
	}
}

// newMCPServer builds a standard MCP server exposing every engine tool.
func newMCPServer(version string, logger *slog.Logger) (*mcp.Server, error) {
	identity := loadServerIdentity()
	impl := &mcp.Implementation{
		Name:        identity.Name,
		Title:       identity.Title,
		Description: identity.Description,
		Version:     version,
		WebsiteURL:  identity.WebsiteURL,
	}
	opts := &mcp.ServerOptions{
		Instructions: "Deterministic Chinese metaphysics computation and auditable divination events. Discover available capabilities with tools/list before use. When required, first resolve time, location, or calendar facts; then use the discovered chart, calendar, divination, or fengshui capability. Auto Liuyao casting creates a new event on each call; coins/yaos receipts reproduce one event. Pass returned artifacts unchanged, interpret only returned facts, and do not invent missing fields. Conclusions are conditioned readings from a traditional culture perspective, not medical, legal or investment advice.",
		Logger:       logger,
		// 2026-07-28 is the preferred stateless surface. Standard clients still
		// discover this endpoint through the initialize handshake, whose newest
		// era is 2025-11-25; advertising only 2026 made the SDK return that
		// handshake fallback and then reject its required header on the next call.
		SupportedProtocolVersions: []string{"2025-11-25", "2026-07-28"},
		// The deterministic tool set is immutable for the life of an image.
		// Advertising listChanged makes modern clients open a long-lived
		// subscriptions/listen stream, which this stateless surface does not
		// need and cannot keep alive behind its request timeout.
		Capabilities: &mcp.ServerCapabilities{
			Tools: &mcp.ToolCapabilities{},
		},
	}
	s := mcp.NewServer(impl, opts)

	if err := mcpcontract.RegisterTools(s); err != nil {
		return nil, err
	}
	return s, nil
}

func runStdio(version string, logger *slog.Logger) error {
	s, err := newMCPServer(version, logger)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return s.Run(ctx, &mcp.StdioTransport{})
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
