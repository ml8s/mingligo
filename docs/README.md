# Documentation map

MingliGo documentation has three layers:

1. **README** introduces the library and its accuracy model.
2. **Time model** defines the calendar and astronomy contract.
3. **Domain documents** define what each `pkg/*` package computes and which
   conventions are explicit.

## Core documents

| Document | Purpose |
|---|---|
| [`TIME_MODEL.md`](TIME_MODEL.md) | Civil time, standard time, true solar time, 时辰, lunar bases, precision |
| [`domain/README.md`](domain/README.md) | Domain-package architecture and contracts |
| [`domain/README.zh-CN.md`](domain/README.zh-CN.md) | Chinese domain-package overview for wiki writing |
| [`VALIDATION.md`](VALIDATION.md) | Golden fixtures, reference projects, licenses, known limits |
| [`decisions/0003-official-mcp-adapter.md`](decisions/0003-official-mcp-adapter.md) | Why MCP is a typed example adapter, not the core |

## Domain documents

| Document | Scope |
|---|---|
| [`domain/foundations.md`](domain/foundations.md) | `tianwen`, `ganzhi`, `tables`, `luoshu`, `city` |
| [`domain/bazi.md`](domain/bazi.md) | 四柱, cycles, strength/pattern evidence, relations, 神煞 |
| [`domain/ziwei.md`](domain/ziwei.md) | 紫微斗数 chart and flow layers |
| [`domain/liuyao.md`](domain/liuyao.md) | 六爻 casting and 纳甲 board |
| [`domain/qimen.md`](domain/qimen.md) | 奇门遁甲 scopes, schools, and 定局 methods |
| [`domain/huangli.md`](domain/huangli.md) | 黄历 facts and event queries |
| [`domain/fengshui.md`](domain/fengshui.md) | 八宅, 玄空, 24 mountains, annual flying stars |

## Writing conventions

- English is canonical for implementation docs.
- Chinese domain terms are kept when they are the precise terminology.
- "Fact" means a deterministic domain output; it is not life advice.
- Rule names, policy names, and school names are public API whenever callers
  must reproduce or audit a result.
- A reference fixture proves agreement with a source at a recorded version; it
  does not imply endorsement by that source.
