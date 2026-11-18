# MingliGo official MCP migration plan

## Goal

Make MingliGo a trustworthy Go fact library and make its example MCP server a
standards-compliant official Go SDK application.

## Non-negotiables

- `pkg/**` is a pure domain/facts layer: no MCP, product, counsel, Skill,
  deployment, or transport imports.
- `examples/mingligo-mcp` is a production-quality adapter, not a throwaway
  demo.
- Every MCP tool has explicit wire DTOs.
- No custom schema/handler registry replaces the official SDK.
- No `json.RawMessage` MCP handler return type.
- Existing tool names and public schemas stay compatible unless a migration
  review explicitly records the change.
- Every schema/handler mismatch fails the request and CI.

## Target architecture

```text
pkg/ganzhi, pkg/bazi, pkg/ziwei, ...
  Pure domain facts and invariants.

examples/mingligo-mcp/internal/mcpcontract
  Wire DTOs, domain -> wire mapping, and official SDK typed registration.
  No server lifecycle enters this package.

examples/mingligo-mcp/main.go
  Composition root: config, logging, observability, server lifecycle.
```

## Work packages

### WP1 — Freeze the compatibility baseline

- Generate and commit `tools/list` schema snapshots for all 17 tools.
- Store representative request/output golden fixtures.
- Record known fixes found by output validation (for example, missing
  `bazi_chart.gong_jia[].relation_kind`).

### WP2 — Introduce the official adapter

- Add `mcpserver.RegisterTools(*mcp.Server) error`.
- Use `mcp.AddTool[In, Out]` for every tool.
- Keep only explicit schemas for:
  - discriminated input unions;
  - current named-array output compatibility roots;
  - frozen fields whose constraints cannot be expressed with struct tags.
- Remove the custom Registry from the production path.

### WP3 — Convert raw handlers to wire DTOs

Convert these handlers first:

- `bazi_chart`
- `bazi_bond`
- `bazi_liupan`
- `ziwei_chart`
- `ziwei_liupan`
- `ziwei_bond`
- `qimen_chart`

Each must return a typed wire DTO. Mapping is explicit; if a temporary JSON
roundtrip is used, it must be marked as migration debt and covered by a
contract test.

### WP4 — Tests and quality gates

- Unit tests for domain packages.
- Contract tests for wire DTO mapping.
- In-process MCP client tests over the SDK transport.
- Error-branch tests, not only happy paths.
- Schema compatibility tests against the frozen baseline.
- `go vet`, `gofumpt`, `golangci-lint`, `go test ./... -race`.

### WP5 — Example quality

- Standard SDK HTTP and stdio lifecycle.
- Config only at the composition root.
- Observability middleware outside domain handlers.
- Graceful shutdown and context propagation.

## Completion criteria

- No custom registry or `RawTool` remains.
- Every registered tool uses official typed `mcp.AddTool`.
- Every wire DTO has explicit field ownership and mapping.
- Official SDK validates input and output on every request.
- All CI gates pass.
- Old and new tool schemas are compared and documented.

## Completion status

Complete for the MingliGo repository. Official SDK registration, typed wire
DTOs, in-process transport tests, lint, race, full tests, the aggregate HTTP
smoke, and the pre/post schema review have passed. Cross-repository release and
deployment work is intentionally out of scope for this boundary.
