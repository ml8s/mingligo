# `pkg/huangli`

`huangli` computes one-day 黄历 facts from a civil date string. It is a
date-level domain: unlike birth charts, it does not take a longitude or derive a
true-solar 时辰 for its primary day selection.

## API

```go
day, err := huangli.QueryDate("2026-06-28")
```

To assess an event:

```go
day, err := huangli.QueryDate("2026-06-28", "wedding")
```

The event identifier is validated against the embedded event table. The result
records the event, label, suitability, and reason.

## Day facts

`huangli.Day` includes:

- civil date;
- 日柱;
- day-element 五行;
- 建除;
- 黄道/黑道 star and path;
- 喜神, 财神, and 福神 directions;
- 彭祖百忌 stem and branch taboos;
- 28 lunar mansion;
- current 节气 and days within it;
- 人元司令分野;
- twelve 时辰 entries with fortune facts;
- optional event assessment.

## Rule basis

Core rules come from:

- 干支 day and solar-term month calculations;
- 建除 rules;
- 青龙/黄道 rules;
- Peng Zu stem/branch taboos;
- 28 mansions;
- 人元司令分野;
- embedded event-to-建除 mappings.

## Validation

A 90-case golden matrix verifies `QueryDate` across 15 events and all 12 建除
states. Calendar pillars and 建月 use `lunar-python` 1.4.8 and `sxtwl` 2.0.7
consensus; event rules follow 《协纪辨方书》. See
[`../VALIDATION.md`](../VALIDATION.md#huangli).
