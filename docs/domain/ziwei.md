# `pkg/ziwei`

`ziwei` computes 紫微斗数 chart facts and flow projections from a standard
birth instant, longitude, gender, and an explicit lunar calendar basis.

## Time contract

The default entry is:

```go
chart, err := ziwei.ComputeChartAt(standardTime, longitude, gender)
```

It uses the `local_astronomical` lunar basis and derives:

- the lunar date from the standard local date;
- the 时辰 from local apparent solar time.

To select a school-compatible calendar basis explicitly, use:

```go
chart, err := ziwei.ComputeChartAtBasis(
    standardTime,
    longitude,
    gender,
    tianwen.LunarBasisOfficial,
)
```

Available bases are documented in
[`TIME_MODEL.md`](../TIME_MODEL.md#lunar-projection-bases):

- `official`: Beijing civil date and UTC+8 calendar;
- `local_label`: local date label with the UTC+8 calendar table;
- `local_astronomical`: local standard date and locally assigned lunations.

This choice is not cosmetic: midnight, overseas longitudes, and calendar-boundary
cases can change the lunar date used to place the chart.

## Core chart

`ziwei.ComputeChartAt` returns a complete typed chart. The model is intended for
application and interpretation layers, not for string scraping. It includes the
palace arrangement, star placements, brightness, four transformations, and
related chart facts defined by the implementation.

## Cycle APIs

Beyond the natal chart, the package exposes explicit cycle layers:

| Cycle | API family |
|---|---|
| 大限 | `ComputeDaXian` |
| 流年 | `ComputeLiuNian` |
| 流月 | `ComputeLiuYueWith` |
| 流日 | `ComputeLiuRiWith` |
| 流时 | `ComputeLiuShiWith` and `ComputeLiuShiAtTimeWith` |

The `With` APIs accept a `tianwen.LunarProjector` so the caller keeps the same
lunar basis used for the natal chart. Resolver APIs also expose the boundary
time facts rather than silently rounding a date.

## Compatibility

`ComputeBond` compares two subjects' chart facts and returns a typed compatibility
result. Callers should pass consistently prepared chart facts; mixing one chart
computed with `official` lunar basis and another with `local_astronomical` is a
domain-policy mismatch, not a numeric bug.

## Validation

Ziwei fixtures include 100 flow-board golden cases and minor-star brightness
tables cross-checked against `iztro` 2.5.8/2.6.1. See
[`../VALIDATION.md`](../VALIDATION.md#ziwei) for details.
