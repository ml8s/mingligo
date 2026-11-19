# Feng-shui packages

MingliGo splits feng shui into directional tables, 八宅, 玄空飞星, and annual
flying stars.

## Shared directional model

`pkg/fengshui` owns the 24 mountains and purple-white star vocabulary.

### 24 mountains

`Mountains24Table` contains all 24 mountains, each spanning 15 degrees:

- index `0` starts at 子 / north `0°`;
- every record has angle, element, yin/yang, trigram, and 三元龙
  (`天元龙`, `地元龙`, `人元龙`);
- angles are direction facts, not an automatic compass-declination correction.

### Purple-white stars

The star table maps numbers 1–9 to names, colors, and elements. The 洛书 flying
order is shared by annual and 玄空 calculations.

## `pkg/bazhai` — 八宅

### 命卦

```go
chart, err := bazhai.ComputeChartAtStandard(standardTime, gender)
```

The standard time is reduced to the exact 立春 year. The birth year and gender
then produce:

- 命卦;
- 东四命 / 西四命 group;
- 大游年 生气, 天医, 延年, 伏位, 祸害, 五鬼, 六煞, and 绝命 directions;
- annual purple-white flying stars;
- year-boundary metadata.

The low-level entry `ComputeMingGua(gender, birthYear)` accepts an already
resolved 立春 year.

### 门主灶

```go
layout, err := bazhai.ComputeLayout(mingGua, doorGua, masterGua, stoveGua)
```

The result evaluates each position for:

- direction;
- trigram;
- element;
- 大游年 star;
- rating;
- 东四/西四 group match.

Invalid trigrams fail closed instead of producing an empty match.

## `pkg/xuankong` — 玄空飞星

```go
chart, err := xuankong.ComputeChartAtStandard(
    standardTime,
    sitMountainIndex,
    faceMountainIndex,
)
```

The standard time determines the 三元运 year by the exact 立春 instant.

The current implementation computes the 下卦正向 board:

- period star flight;
- mountain star from the sitting palace;
- facing star from the facing palace;
- forward/reverse flight from 三元龙 yin/yang;
- no replacement-star path because API input has no replacement-angle degree.

It rejects sitting/facing combinations that are not opposite 24-mountain pairs.

### Board facts

A `xuankong.Chart` includes:

- period and mountain/facing names;
- nine palaces with period, mountain, and facing stars;
- 旺山旺向 / 上山下水 and four-situation facts;
- 伏吟 and 反吟 layers;
- 双星加会;
- 收山出煞.

## `pkg/fengshui` — annual flying stars

```go
board := fengshui.ComputeAnnualFlyingStars(year)
```

The board records:

- year;
- center star;
- `lichun` year boundary;
- nine palace stars with number, name, element, and 入中 status.

The same purple-white implementation is reused by 八宅 annual stars.
