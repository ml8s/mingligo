# Domain packages

MingliGo is a fact core. Each `pkg/*` package computes deterministic Chinese
metaphysics facts from explicit inputs and rules. The packages do not provide
life advice, product UX, transport, or persistence.

## Layered model

```text
Foundations
  tianwen, ganzhi, tables, luoshu, city
        ↓
Chart and cycle domains
  bazi, ziwei
        ↓
Divination and calendar domains
  liuyao, qimen, huangli
        ↓
Feng-shui domains
  bazhai, xuankong, fengshui
```

Foundations are the only shared language for time and stems/branches. Chart
domains may consume foundations, but they never depend on transport or product
code.

## Package contract summary

| Package | Primary entry | Result model |
|---|---|---|
| `tianwen` | `ComputeBazi`, `ApparentSolarFromStandard`, `ShichenFromStandard`, `ResolveStandard` | Pillars, LAST, lunar dates, boundary candidates |
| `ganzhi` | Typed enums, parsers, relation tables | Gan/Zhi/Wuxing, hidden stems, Nayin, relations |
| `tables` | `Load` | Explicit startup barrier for embedded tables |
| `luoshu` | Shared palace table and loaders | 洛书九宫 / 后天八卦 facts |
| `city` | `ResolveCandidates` | Embedded place candidates with longitude/timezone |
| `bazi` | `ComputeChartAt`, `ComputeCoreAt`, cycle APIs | Complete 八字 chart and flow projections |
| `ziwei` | `ComputeChartAt`, `ComputeChartAtBasis`, flow APIs | Complete 紫微 chart and flow layers |
| `liuyao` | `NewCoinsCasting`, `NewValuesCasting`, `ComputeChartAt` | Auditable casting and 六爻 board |
| `qimen` | `ResolveMethod`, `ComputeChartWithMethodAt` | Typed 奇门 method and board |
| `huangli` | `QueryDate` | One-day 黄历 facts and event assessment |
| `bazhai` | `ComputeChartAtStandard`, `ComputeLayout` | 命卦, 大游年, 门主灶 |
| `xuankong` | `ComputeChartAtStandard`, `ComputeSanYuanYun` | 玄空飞星 board |
| `fengshui` | `ComputeAnnualFlyingStars`, mountain/star tables | 24 mountains and annual stars |

## Shared invariants

### Explicit time

Domains that depend on birth or event time use the common time model in
[`../TIME_MODEL.md`](../TIME_MODEL.md). They do not invent their own DST or
true-solar-time rules. Direct `Compute*At` APIs expect a standard-time
`time.Time`; use `tianwen.ResolveStandard` when DST, ambiguous wall-clock input,
时辰-only input, or longitude ranges require boundary resolution.

### Embedded data, fail-closed loading

Domain packages embed their tables and load them lazily. `tables.Load` is an
optional aggregate startup barrier. Loaders validate the embedded JSON/YAML
shapes and fail closed; they do not fall back to partial tables.

### Policy is public

Where schools differ, the package exposes the chosen rule:

- Bazi exposes 晚子时 and dual-reference 神煞 policy.
- Ziwei exposes lunar projection basis.
- Qimen exposes scope, school, and 定局 method.
- Liuyao casting exposes coin convention and moving lines.
- Feng-shui boards expose period, boundary, and plate policy.

### Facts versus interpretation

Some packages emit evidence such as 扶抑 candidates, 调候 candidates, 格局
structure, 用神 candidates, patterns, or timing candidates. These are domain
facts and auditable intermediate layers. They are not guarantees, counsel, or a
final reading.

### Determinism

Domain entry points do not call the network and do not use random behavior.
Randomness belongs to the caller, except for explicitly named secure casting
helpers. A fixed input and fixed rule set produce the same result.

## Compatibility and scope

- Primary tested calendar range is 1900–2100.
- Longitude is east-positive and must be in `[-180, 180]`.
- Longitude `0` is a valid Greenwich value.
- Timezone names are IANA names.
- Public JSON field names are compatibility surfaces.
- Reference fixtures prove agreement at a recorded source version, not universal
  correctness across all schools.

See [`../VALIDATION.md`](../VALIDATION.md) for exact sources and limits.

## Domain index

| Document | Package scope |
|---|---|
| [`foundations.md`](foundations.md) | `tianwen`, `ganzhi`, `tables`, `luoshu`, `city` |
| [`bazi.md`](bazi.md) | 八字 and time cycles |
| [`ziwei.md`](ziwei.md) | 紫微斗数 |
| [`liuyao.md`](liuyao.md) | 六爻 |
| [`qimen.md`](qimen.md) | 奇门遁甲 |
| [`huangli.md`](huangli.md) | 黄历 |
| [`fengshui.md`](fengshui.md) | 八宅, 玄空, 24 mountains, annual stars |
