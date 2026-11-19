# Foundations

The foundation packages define the shared language used by every chart and
board: time, astronomy, stems/branches, palace tables, and places.

## `tianwen` — astronomy and calendar time

### Responsibilities

- Convert civil/standard time to local apparent solar time (LAST).
- Derive 时辰 and the solar-day boundary.
- Compute stems/branches for year, month, day, and hour.
- Compute 节气 instants and 建月.
- Project civil and lunar dates with explicit lunar bases.
- Resolve ambiguous date/时辰/longitude-range inputs into candidates.

### Core API

| Function | Purpose |
|---|---|
| `ApparentSolarFromStandard(t, longitude)` | Convert standard time to LAST using the timestamp's standard offset |
| `ShichenFromStandard(t, longitude)` | Derive the 时辰 branch from LAST |
| `ComputeBazi(civilTime, solarTime)` | Compute year/month pillars from civil time and day/hour pillars from LAST |
| `SolarTermTime(year, targetLongitude)` | Compute a solar-term instant |
| `JianYue(t)` | Compute the solar-term month branch |
| `NewLunarProjector(basis, location)` | Build a lunar projection under an explicit basis |
| `ResolveStandard(input)` | Normalize an ambiguous birth-time input into candidates |

### Important distinctions

- Direct compute entry points expect standard time; `ResolveStandard` is the
  front door for DST, vague time, ranges, and boundary ambiguity.
- Year and month pillars use civil/standard instants and solar-term boundaries.
- Day and hour pillars use LAST.
- Lunar calendar projection uses a separately selected basis.
- LAST wall values are minute-precision in the current public conversion; see
  [`../TIME_MODEL.md`](../TIME_MODEL.md#precision-and-current-limits).

## `ganzhi` — stems, branches, and shared relations

### Responsibilities

`ganzhi` owns the primitive domain vocabulary and shared relational tables:

- 天干 and 地支 typed enums;
- 五行, yin/yang, and zodiac mappings;
- 60 甲子;
- 天干五合, 地支六合, 三合, 三会;
- 六冲, 相刑, 六害, 破, 暗合, and stem clashes;
- 藏干;
- 纳音;
- 十二长生;
- 人元司令分野.

### Design

These tables are loaded from embedded JSON and validated before use. Other
domains import `ganzhi` instead of redefining the same terms. This prevents
divergent hidden-stem tables, Nayin tables, or relation definitions.

## `tables` — fail-closed startup barrier

Domain loaders are lazy and idempotent. `tables.Load` provides one explicit
entry point for services that prefer to fail before accepting requests:

```go
if err := tables.Load(); err != nil {
    return fmt.Errorf("initialize mingligo: %w", err)
}
```

It initializes `ganzhi`, `bazi`, `ziwei`, `liuyao`, `qimen`, `huangli`,
`bazhai`, and `xuankong`. It is optional for direct library use.

## `luoshu` — 洛书九宫

`luoshu` is the shared table for:

- palace numbers 1–9;
- palace names and directions;
- trigram names;
- five elements;
- yin/yang.

Divination and feng-shui packages use this table instead of maintaining private
palace arrays.

## `city` — embedded places

`city.ResolveCandidates(query)` searches an embedded table without network
access. The generated dataset includes:

- 2,844 China county-level records;
- 337 city records;
- 88 selected world cities.

Each record can contain longitude, latitude, timezone, localized names,
administrative region, aliases, and a bounding box.

The resolver is deliberately conservative:

- exact matches rank first;
- ambiguous names return multiple candidates;
- it never silently chooses one 朝阳区;
- empty results mean the embedded table has no match.

Geographic sources and license notices are documented in
[`pkg/city/ATTRIBUTION.txt`](../../pkg/city/ATTRIBUTION.txt).
