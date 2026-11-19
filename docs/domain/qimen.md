# `pkg/qimen`

`qimen` computes 奇门遁甲 boards across multiple scopes, schools, and 定局
methods. Its method catalog is closed and validated, so callers cannot silently
combine incompatible rules.

## Method axes

### Scope

| Scope | Identifier |
|---|---|
| 时家 | `hour` |
| 刻家 | `quarter` |
| 日家 | `day` |
| 月家 | `month` |
| 年家 | `year` |

### Plate school

| School | Identifier |
|---|---|
| 转盘 | `zhuanpan` |
| 洛书飞盘 | `luoshu_feipan` |
| 鸣法飞盘 | `mingfa_feipan` |
| 金函玉镜 | `jinhan_yujing` |

### 定局 method

| Method | Identifier |
|---|---|
| 拆补 | `chaibu` |
| 置闰 | `zhirun` |
| 茅山 | `maoshan` |
| 无定局 | `none` |

Quarter scopes additionally expose named 刻家 rules. The catalog validates
allowed combinations and contains 17 valid method combinations.

## API style

Resolve a typed method first:

```go
method, err := qimen.ResolveMethod(
    qimen.ScopeHour,
    qimen.SchoolZhuanPan,
    qimen.DingjuChaiBu,
    "",
)
```

Then compute:

```go
chart, err := qimen.ComputeChartWithMethodAt(
    standardTime,
    longitude,
    method,
)
```

The lower-level `ComputeChart(gt, st)` entry exists for compatibility and uses
the default hour/转盘 path. New callers should prefer the explicit method API.

The separate Golden Mirror entry is:

```go
jinhan := qimen.ComputeJinhanChartAt(standardTime, longitude)
```

## Time contract

Qimen mixes time dimensions deliberately:

- 节气 and 定局 comparisons use absolute civil instants;
- day/hour pillars follow MingliGo's solar-day and LAST rules;
- quarter rules can expose additional hour-boundary policy;
- charts retain method metadata such as solar term, 元, lead pillar, duty star,
  duty door, and boundary behavior.

This separation prevents a longitude-only day/hour correction from silently
changing a solar-term determination.

## Board facts

Depending on method and scope, the result can include:

- yin/yang 遁 and 局数;
- 元 and lead/符头 pillars;
- duty star and duty door;
- 空亡 and 旺衰;
- nine palaces with earth plate, heaven plate, stars, doors, and spirits;
- 天禽寄宫 rules;
- pattern detections and structured evidence;
- 用神 palace results for question symbols and year-life stems.

Method metadata remains on the chart so a rendered board can be reproduced and
audited.

## Validation

Qimen has one of the strongest external-fixture matrices in the repository:

- external rotating-board vectors from `atopx/qimen`;
- 置闰 vectors and state transitions from `atopx/qimen`;
- method-matrix vectors from `atopx/qimen`;
- Ming-fa vectors from `potuo/feipan-qimen`;
- 茅山 vectors from `deminzhang/qimen-go`;
- 金函玉镜 vectors from `kentang2017/kinqimen`;
- quarter vectors and twelve-minute quarter vectors with recorded provenance;
- curated classical pattern vectors.

See [`../VALIDATION.md`](../VALIDATION.md#qimen) for exact counts and licenses.
