# ADR 0004: Domain-owned aggregated output, thin protocol adapters

## Status

Accepted.

## Context

MCP handlers originally rebuilt domain results into wire DTOs with hand-written
projection functions. For most tools these projections were pure type packing
(enum → readable string, value → pointer) with no 命理 logic. Some tools also
performed real domain aggregation in the adapter: grouping parallel arrays into
one object (e.g. a 六爻 line), encoding names, and separating analysis facts.
This duplicated domain knowledge into the transport layer and would force a
second set of projections when a REST adapter is added.

## Decision

1. Domain packages own the complete, readable, domain-aggregated output shape.
   Grouping related facts into one object is domain knowledge and lives in
   `pkg/**`.
2. Adapters (MCP, and later REST) transmit the domain result unchanged. Wire
   DTOs exist only to carry JSON Schema tags (`jsonschema`) and enum
   constraints; they are not a second data model.
3. Presentation-only concerns (generated IDs, display labels, enum dropdowns,
   rendering) stay in the adapter/frontend, not the domain.
4. Time values use `time.Time` with symbol names distinguishing standard time,
   civil time, local mean time, and solar time. Calendar dates use
   `tianwen.Date`. Protocol schema tags do not enter the domain.
5. Field names follow the domain vocabulary (pinyin terms such as `gong`,
   `ri_chen_gan`, `dong_yao`). Migrating a wire contract's English names is a
   separate, explicitly-reviewed change.

## Consequences

- Handlers return domain values directly; projection helpers are removed.
- Every tool keeps a domain-vs-wire equivalence test so schema drift is caught.
- Adding a REST adapter needs no new projection layer.
- liuyao, whose adapter view is coupled to MCP enum types, migrates in two
  steps: first move view/analysis aggregation into domain while keeping the
  existing English contract, then rename fields as a reviewed migration.
