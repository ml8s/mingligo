# `pkg/liuyao`

`liuyao` computes 六爻纳甲 boards from an auditable casting and a standard
event time.

## Casting is a source fact

A six-line casting is not thrown away after deriving a hexagram. It is retained
as a receipt.

MingliGo supports two explicit modes.

### Coin rounds

```go
rounds := [][3]liuyao.CoinFace{
    {"反", "反", "反"}, // example only; six rounds are required
}
casting, err := liuyao.NewCoinsCasting(rounds)
```

Rules:

- exactly six rounds;
- exactly three coins per round;
- bottom line first;
- heads `正` scores 3, tails `反` scores 2;
- sums `6`, `7`, `8`, `9` map to 老阴, 少阳, 少阴, 老阳;
- old lines are moving lines;
- the receipt records the convention and normalized rounds.

### Line values

Callers that already have values use:

```go
casting, err := liuyao.NewValuesCasting([6]int{7, 8, 8, 9, 7, 8})
```

The receipt recomputes moving lines and rejects inconsistent source facts.

## Standard-time entry

```go
chart := liuyao.ComputeChartAt(
    standardTime,
    longitude,
    liuyao.YongQiCai,
    casting.Yaos,
)
```

The time model is shared with all other domains:

- year and month pillars come from civil/standard time and solar terms;
- day and hour pillars come from true solar time.

## Board facts

A `liuyao.Chart` computes and keeps:

- original and changed lines;
- hexagram names and palace;
- 纳甲 stems/branches;
- 五行;
- 六亲;
- 世应;
- 六神;
- 日辰 relations;
- 月建 strength;
- 旬空, 月破, and 墓库;
- hidden lines and 伏神;
- 用神 position and aggregate state;
- 原神/忌神 movement relations;
- patterns such as 旬空, 月破, 飞伏, 进退神, 六冲/六合, 反吟/伏吟, 随鬼入墓,
  独发/独静, and 用神两现;
- timing candidates.

## Pattern semantics

Patterns are not returned as vague labels. For example, 旬空 and 月破 record
whether the condition is true or false based on support such as 旺相, moving
lines, moving-line support, and day support. This lets interpretation layers
apply their own rules without reversing a collapsed label.

## Randomness

The domain board is deterministic. Randomness is isolated to casting helpers
such as secure 起卦; once a casting exists, board computation has no hidden
random behavior.

## Validation

The 京房八宫 core is covered by a complete 64-hexagram golden table. Additional
golden and oracle tests cover 装卦, conflicts, moving lines, timing, lifecycle,
and calendar boundaries. See [`../VALIDATION.md`](../VALIDATION.md#liuyao).
