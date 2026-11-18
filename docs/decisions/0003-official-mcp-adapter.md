# ADR 0003: Use the official Go MCP SDK typed adapter

## Status

Accepted.

## Context

MingliGo previously had a custom registry around the official Go MCP SDK. The
original migration kept `RawTool` handlers returning `json.RawMessage`, custom
schema slimming, and manually maintained overrides. This weakened the core
guarantee that a tool's Go handler, published schema, and serialized result
share one contract.

The project's public surface must be a trustworthy standard-library-quality Go
API. Domain code remains pure; MCP transport is an application adapter.

## Decision

MingliGo will use the official Go MCP SDK's typed registration API:

```go
mcp.AddTool[In, Out](server, tool, handler)
```

where:

- `In` is the tool's wire input DTO;
- `Out` is the tool's wire output DTO;
- schemas are inferred by the SDK, with explicit schemas only for true unions,
  discriminated forms, or named domain-array roots;
- handlers map pure domain results to wire DTOs;
- `json.RawMessage` tool handlers are prohibited.

## Architecture

```text
pkg/<domain>           pure domain facts; no MCP, product, or transport imports
examples/mingligo-mcp/internal/mcpcontract
                       wire DTOs and explicit tool-to-domain adapters
examples/mingligo-mcp  composition root and standard MCP adapter
```

The example is the canonical, standards-compliant MCP adapter. It is
production-quality code, not a toy.

## Constraints

- Public routes and tool names remain frozen.
- The existing schema shape is the review baseline. Deliberate corrections are
  documented rather than silently introduced.
- Official SDK validates input and output on every call.
- `json.RawMessage` may cross internal boundaries only after an explicit wire
  DTO mapping, and never as a registered MCP handler output.
- No product identity, Skill, counsel, paywall, or deployment logic enters
  MingliGo.

## Migration plan

1. Freeze the current `tools/list` schemas and golden outputs.
2. Replace the custom registry with official `mcp.AddTool` registrations.
3. Convert every `RawTool` to a typed wire DTO handler.
4. Keep explicit schemas only for unions and named-array compatibility roots.
5. Add in-process MCP client tests covering every tool and representative
   error branch.
6. Add CI gates: gofumpt, go vet, golangci-lint, go test, race, and MCP smoke.
7. Compare old and new tools/list schemas before merging; record intentional
   fixes such as missing fields rather than silently changing contracts.

## Migration review

The migration is complete at the MingliGo boundary:

- Production registers through `RegisterTools` and `mcp.AddTool[In, Out]`;
  the custom Registry and raw handlers are deleted.
- All 17 tools use typed wire DTOs. The aggregate `/mcp` surface remains 17
  tools, and the local HTTP smoke passes.
- A pre/post `tools/list` review found the same names and semantically equal
  schemas for 15 tools.
- `qimen_chart` is the sole deliberate schema correction: its output is now an
  explicit standard/Jinhan `oneOf` instead of a flattened union, and the
  twelve-minute quarter's domain-defaulted options are optional. Runtime tests
  cover both standard and Jinhan charts and all quarter-rule paths.
- The remaining public-text correction removes a counsel-product reference
  from `bazi_chart.fu_yi`; the factual fields are unchanged.
- MingliGo no longer embeds hosted web origins. Browser deployments must set
  `MINGLIGO_ALLOWED_ORIGINS`; server-to-server and local use need none.
- The legacy consumer-owned projection contract was removed; MingliGo tests
  now assert its own MCP schemas and domain fixtures directly.
