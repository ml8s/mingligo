<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./assets/logo-dark.svg">
    <img alt="MingliGo logo" src="./assets/logo.svg" width="96">
  </picture>
</p>

<h1 align="center">MingliGo</h1>

<p align="center">
  <strong>A deterministic fact core for Chinese metaphysics in pure Go.</strong>
</p>

<p align="center">
  Accurate time · True solar 时辰 · Explicit schools · Typed chart facts
</p>

<p align="center">
  <a href="https://github.com/ml8s/mingligo/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/ml8s/mingligo/actions/workflows/ci.yml/badge.svg"></a>
  <a href="https://github.com/ml8s/mingligo/actions/workflows/security.yml"><img alt="Security scans" src="https://github.com/ml8s/mingligo/actions/workflows/security.yml/badge.svg"></a>
  <img alt="Pure Go" src="https://img.shields.io/badge/Pure_Go-1.26-00ADD8?style=flat-square&logo=go&logoColor=white">
  <a href="./LICENSE"><img alt="License" src="https://img.shields.io/badge/License-MIT-green?style=flat-square"></a>
</p>

<p align="center">
  <a href="./README.md"><img alt="English" src="https://img.shields.io/badge/docs-English-2563EB?style=flat-square"></a>
  <a href="./docs/i18n/README.zh-CN.md"><img alt="简体中文" src="https://img.shields.io/badge/docs-简体中文-DC2626?style=flat-square"></a>
  <a href="./docs/i18n/README.ja.md"><img alt="日本語" src="https://img.shields.io/badge/docs-日本語-059669?style=flat-square"></a>
  <a href="./docs/i18n/README.es.md"><img alt="Español" src="https://img.shields.io/badge/docs-Español-C2410C?style=flat-square"></a>
  <a href="./docs/i18n/README.pt-BR.md"><img alt="Português (BR)" src="https://img.shields.io/badge/docs-Português_(BR)-65A30D?style=flat-square"></a>
  <a href="./docs/i18n/README.de.md"><img alt="Deutsch" src="https://img.shields.io/badge/docs-Deutsch-334155?style=flat-square"></a>
</p>

MingliGo is a pure Go library for deterministic Chinese metaphysics facts. It
computes calendars, natal charts, divination boards, and feng-shui boards as
typed Go values. The primary product is `pkg/*`; the MCP server is only an
example adapter.

## Accuracy starts with time

Traditional Chinese metaphysics observes the sky, not merely an administrative
clock. 日晷 tells local solar time; it does not follow a national timezone
boundary. MingliGo therefore separates civil time from astronomical time and
uses one shared conversion across domains:

```text
civil birth time
  → IANA timezone / standard-time resolution (DST removed)
    → longitude correction
      → equation of time
        → local apparent solar time / 真太阳时
          → 时辰, day boundary, solar-term boundary, year/month pillars
```

Key properties:

- **One standard-time entry point.** APIs such as
  `bazi.ComputeChartAt(standard, longitude, gender)` accept a civil instant and
  a longitude.
- **Standard time is a contract.** Direct chart APIs expect a `time.Time`
  already expressed in standard time. For ambiguous historical wall clocks,
  `tianwen.ResolveStandard` removes DST and returns boundary candidates.
- **DST is not silently promoted to 命理 time.** The embedded IANA database
  supports the resolver and historical standard-offset handling.
- **Longitude is explicit.** Longitude `0` is Greenwich, not a default Beijing
  meridian.
- **时辰 follows the sun.** Day and hour projections use local apparent solar
  time, including midnight crossings in eastern longitudes and earlier local
  solar hours in western longitudes.
- **Solar terms are instants.** 立春 and the 24 节气 are not approximated as
  calendar dates.
- **Lunar bases are explicit.** Official UTC+8, local-label, and
  local-astronomical projections are separate API choices.

See [Time model](docs/TIME_MODEL.md) for formulas, lunar bases, boundary
resolution, precision, and limitations.

## Library use

```bash
go get github.com/ml8s/mingligo@main
```

```go
package main

import (
  "encoding/json"
  "fmt"
  "time"

  "github.com/ml8s/mingligo/pkg/bazi"
  "github.com/ml8s/mingligo/pkg/ganzhi"
  "github.com/ml8s/mingligo/pkg/tables"
)

func main() {
  // Optional fail-fast startup. Domain APIs also load embedded tables lazily.
  if err := tables.Load(); err != nil {
    panic(err)
  }

  birth := time.Date(1984, 2, 15, 8, 0, 0, 0, time.FixedZone("CST", 8*3600))
  chart := bazi.ComputeChartAt(birth, 116.407526, ganzhi.Male)

  data, err := json.MarshalIndent(chart, "", "  ")
  if err != nil {
    panic(err)
  }
  fmt.Println(string(data))
}
```

The output is a typed domain chart, not an opaque string. Four pillars, hidden
stems, ten gods, relations, 神煞, 大运, 用神 evidence, and policy metadata remain
available to callers.

## Domain packages

### Shared foundations

| Package | Responsibility |
| --- | --- |
| `tianwen` | Solar/lunar astronomy, 节气, true solar time, 时辰, lunar projection |
| `ganzhi` | Heavenly stems, earthly branches, 60 甲子, hidden stems, Nayin, relations |
| `tables` | Fail-closed aggregate loading for embedded domain tables |
| `luoshu` | Shared 洛书九宫 and 后天八卦 table |
| `city` | Embedded city/county/world places with longitude and timezone |

### Charts and cycles

| Package | Facts computed |
| --- | --- |
| `bazi` | 四柱, 藏干, 十神, 神煞, 合会冲刑破害, 大运/流年/流月/流日/流时/小运, 胎元/命宫/身宫, 扶抑/调候/格局 evidence |
| `ziwei` | 紫微命盘, lunar basis selection, 大限/流年/流月/流日/流时, 合盘 |

### Divination and calendars

| Package | Facts computed |
| --- | --- |
| `liuyao` | Coin/value casting receipts, 纳甲装卦, 世应, 六亲, 六神, 用神, 旺衰, patterns, 应期 candidates |
| `qimen` | Hour/day/quarter/month/year boards; rotating, flying, Ming-fa, and Golden Mirror schools |
| `huangli` | 日柱, 建除, 黄道/黑道, direction gods, Peng Zu taboos, 28 mansions, per-时辰 facts |

### Feng shui

| Package | Facts computed |
| --- | --- |
| `bazhai` | 命卦, 东四/西四 group, 大游年 directions, 门主灶, annual flying stars |
| `xuankong` | 玄空飞星 board, 三元运, four situations, 伏吟/反吟, 双星加会, 收山出煞 |
| `fengshui` | 24 mountains, purple-white stars, annual flying-star board |

Detailed per-domain contracts live under [`docs/domain`](docs/domain/README.md).

## Explicit rules, not hidden choices

Chinese metaphysics has real school differences. MingliGo does not collapse them
into an undocumented default. Where a convention matters, it is represented as a
field or typed method:

- 晚子时 day-pillar behavior is explicit.
- 八字 神煞 records its year/day dual-reference policy.
- 奇门 method is a closed catalog of scope × school × 定局 combinations.
- 紫微 lunar projection chooses an explicit basis.
- Evidence and final judgment are separated; the core emits candidates and
  structural evidence rather than advice.

This keeps interpretation layers and products free to choose a school while
still knowing exactly which rule produced a fact.

## Validation and golden sources

MingliGo is tested against both independent rule anchors and external reference
fixtures. Representative coverage includes:

- Lunar and 建月 consensus cases from `lunar-python` 1.4.8 and `sxtwl` 2.0.7.
- 192 Bazi solar-term boundary cases requiring three-oracle consensus.
- 115 lunar-typescript-derived Bazi golden cases for pillars, 大运, and 三元.
- Qimen vectors from `atopx/qimen`, `potuo/feipan-qimen`,
  `deminzhang/qimen-go`, and `kentang2017/kinqimen`.
- 400 八宅命卦 cases across 1900–2099 and both genders.
- 108 玄空下卦正向 combinations across nine 元运.

The full provenance, license, fixture paths, tolerances, and known limits are in
[`docs/VALIDATION.md`](docs/VALIDATION.md).

## Architecture

| Boundary | Responsibility |
| --- | --- |
| `pkg/*` | Deterministic domain facts, embedded tables, and validation |
| `examples/mingligo-mcp` | Production-quality example MCP adapter |
| `contracts/mcp-tools.json` | Frozen MCP tool-name catalog |

The fact core never imports MCP, transport, product, deployment, or judgment
product code. Domain packages can be consumed directly by Go services, CLIs, or
an adapter.

## Example MCP adapter

The optional adapter publishes 17 tools with the official Go MCP SDK. Typed wire
DTOs are registered through `mcp.AddTool`; handlers map domain results to wire
contracts.

```bash
make run        # HTTP server at :8081, endpoint /mcp
make test       # Go and MCP contract tests
make boundary   # reject MCP/internal imports inside pkg/*
make smoke-mcp  # build and smoke-test the container image
```

Use `--stdio` for the same tool surface over standard input/output.

## Documentation

- Documentation map: [`docs/README.md`](docs/README.md)
- Time model: [`docs/TIME_MODEL.md`](docs/TIME_MODEL.md)
- Domain packages: [`docs/domain/README.md`](docs/domain/README.md)
- Domain overview in Chinese: [`docs/domain/README.zh-CN.md`](docs/domain/README.zh-CN.md)
- Validation sources: [`docs/VALIDATION.md`](docs/VALIDATION.md)
- MCP transport: [Streamable HTTP](https://modelcontextprotocol.io/specification/2025-06-18/basic/transports)
- MCP decision: [Official Go MCP SDK typed adapter](docs/decisions/0003-official-mcp-adapter.md)
- Geographic attribution: [`pkg/city/ATTRIBUTION.txt`](pkg/city/ATTRIBUTION.txt)

## Development

```bash
make test       # all Go tests
make check      # gofmt, lint, vet, and fact-core boundary checks
make gate       # check + test + build
```

The primary tested calendar range is 1900–2100. Outside that range, callers
should not assume historical calendar, timezone, or domain-table compatibility.

## License and attribution

Code is MIT-licensed. Embedded datasets and reference-derived fixtures keep
their own licenses and attribution requirements. See [`LICENSE`](./LICENSE),
[`docs/VALIDATION.md`](docs/VALIDATION.md), and
[`pkg/city/ATTRIBUTION.txt`](pkg/city/ATTRIBUTION.txt).
