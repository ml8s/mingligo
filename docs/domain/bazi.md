# `pkg/bazi`

`bazi` computes the deterministic facts needed by a 八字 chart and its time
cycles. It is the largest chart domain and the reference implementation of
MingliGo's standard-time API style.

## Time contract

The preferred entry point is:

```go
chart := bazi.ComputeChartAt(standardTime, longitude, gender)
```

Internally:

- year and month pillars use the civil/standard instant and solar-term bounds;
- true solar time is derived from longitude and equation of time;
- day and hour pillars use true solar time;
- gender affects 大运 direction and gender-dependent 神煞 projections, not the
  four pillars.

The lower-level API is:

```go
chart := bazi.ComputeChart(civilTime, trueSolarTime, gender)
```

It is useful when a caller has already produced LAST or must replay an explicit
historical conversion.

## Chart model

A `bazi.Chart` contains typed pillar values. Each pillar carries:

- the stem/branch pair;
- day-master trend;
- 旬 and 旬空;
- 纳音;
- 藏干 projection;
- ten-god projections for transparent and hidden stems;
- self-sitting and 十二长生 trend;
- 神煞 projections;
- void and special-pillar flags.

The chart also contains:

- 大运 direction, start distance, and steps;
- 胎元, 命宫, 身宫;
- 宫甲;
- 纳音 relations;
- 十二长生 positions;
- zodiac;
- relations such as 合, 会, 冲, 刑, 破, 害, and stem combinations;
- ten-god states, element states, and Lu roots;
- 官杀 and 财库 conclusions;
- 扶抑, 调候, and 格局 evidence.

## Time-cycle APIs

After computing a natal core, MingliGo can derive:

| Cycle | API |
|---|---|
| 大运 | Included in `Chart` |
| 流年 | `ComputeLiuNianAtTime` |
| 流月 | `ComputeLiuYueAtTime` |
| 流日 | `ComputeLiuRi` |
| 流时 | `ComputeLiuShi` |
| 小运 | `ComputeXiaoYun` |
| 合盘 | `ComputeBond` |

Each flow has explicit resolver APIs, such as
`ResolveLiuYueTimeAt(standard, location)` and
`ResolveLiuShiTime(standard, location, longitude)`, so callers can inspect the
boundary instants instead of guessing.

## Evidence, not final judgment

八字 schools disagree about 用神 and weight. MingliGo therefore separates:

1. **扶抑 candidates** — strength evidence and favorable/unfavorable element
   candidates;
2. **调候 candidates** — climate-adjustment evidence;
3. **格局 evidence** — month-order structure, transparency, rooting, and pattern
   evidence.

The core does not pretend that these always agree. A judgment product can consume
the evidence and choose its school.

## Policies

### Year, month, day, and hour boundaries

- Year boundary: exact 立春.
- Month boundary: exact 节.
- Day boundary: true solar day.
- Late 子时: current output declares its policy explicitly. Do not infer a new
  convention from a rendered string.

### 神煞 references

The chart records a dual-reference 神煞 school policy. Selected 神煞 may consider
both year and day references, and the output includes the policy instead of
hiding it.

## Validation

Major validation layers include:

- `lunar-typescript` (6tail) authoritative golden cases for pillars, 大运, and
  三元;
- 192 three-oracle solar-term boundary cases from `lunar-python`, `sxtwl`, and
  `bazi-calculator`;
- an external true-solar 八字 oracle fixture;
- independent table and boundary tests for hidden stems, Nayin, ten gods,
  神煞, and 时辰;
- regional boundary tests for longitude, midnight crossing, DST-free standard
  time, and historical offsets.

See [`../VALIDATION.md`](../VALIDATION.md) for exact fixture paths and limits.
