# Validation and reference sources

This document records how MingliGo is checked, which external sources are used,
where fixtures live, and what a fixture does and does not prove.

## Validation philosophy

MingliGo uses three layers:

1. **Independent rule anchors.** Tests compute expected values from classical
   tables or independent algorithms.
2. **External golden fixtures.** Outputs are compared with a named reference at
   a recorded version or commit.
3. **Boundary and property tests.** DST, longitude, midnight, solar terms,
   ambiguous birth times, invalid inputs, and cross-domain consistency are tested
   explicitly.

A golden fixture is reproducibility evidence at a recorded source version. It is
not a claim that the source is the only valid school, nor an endorsement by the
source.

## Primary scope

The primary tested calendar range is **1900–2100**. Some tables and astronomical
functions may work outside it, but callers should not assume compatibility unless
a fixture or API explicitly documents it.

## Time and astronomy

### Internal astronomical implementation

The solar-longitude chain derives from ShouXing / `lunar-typescript`
`ShouXingUtil` and uses VSOP87D coefficients, nutation, aberration, and ΔT.
Generated coefficient tables live under `pkg/tianwen`.

Known-instant tests cover the 2024 solar terms:

| Term | Tolerance |
|---|---|
| 春分 | ±30 seconds |
| 夏至 | ±30 seconds |
| 秋分 | ±30 seconds |
| 冬至 | ±30 seconds |
| 立春 | ±30 seconds |

A 2024 new-moon test has a ±2-minute tolerance. These are target checks for the
tested instants, not a universal error bound for all history.

### True solar time

Tests cover:

- Beijing, Tokyo, and London LAST goldens;
- longitude correction and equation-of-time behavior;
- Kashgar and Urumqi western-longitude offsets;
- midnight crossings in eastern longitudes;
- `longitude=0` not being replaced by Beijing;
- standard offsets with DST removed.

Current limitation: the public LAST wall-time conversion is minute-precision.
See [`TIME_MODEL.md`](TIME_MODEL.md#precision-and-current-limits).

### Lunar and 建月 consensus

File: `pkg/tianwen/testdata/multi_oracle_golden.json`

- 99 cases;
- years 2024–2027;
- fields: lunar date and 建月;
- sources: `lunar-python` 1.4.8 and `sxtwl` 2.0.7;
- policy: all oracles must agree.

Known limits recorded in the fixture:

- date-level solar-term disagreement cases are not hard-consensus candidates;
- LAST is tested by separate solar goldens;
- day pillars are covered by Bazi calendar goldens and domain oracle tables.

## Bazi

### Lunar-typescript golden set

Files:

```text
pkg/bazi/testdata/bazi_golden.json
pkg/bazi/testdata/bazi_golden2.json
pkg/bazi/testdata/bazi_golden3.json
pkg/bazi/testdata/bazi_golden4.json
pkg/bazi/testdata/bazi_golden5.json
pkg/bazi/testdata/bazi_golden6.json
```

- 115 cases total;
- source: `lunar-typescript` (6tail), recorded in the test comment;
- fields include 四柱, 大运 direction/start/steps, and selected 三元 fields;
- tests separately enable 命宫, 身宫, 胎元, and minute-level 起运 checks by file;
- overseas cases allow bounded year/month tolerance because VSOP87 precise terms
  and the reference's approximate term time can interact at hour boundaries.

### Three-oracle solar-term boundary matrix

File: `pkg/bazi/testdata/multi_oracle_golden.json`

- 192 cases;
- years 2024–2027;
- all 24 terms, with one case before and one after each term;
- consensus sources:
  - `lunar-python` 1.4.8;
  - `sxtwl` 2.0.7;
  - `bazi-calculator` 1.0.0;
- all three must agree on the four pillars.

Known limits recorded in the fixture:

- `bazi-calculator` uses date-level solar terms; anchors are placed within ±25
  hours of the `lunar-python` term time;
- late-子 cases are moved off `23:00`; the exact late-子 policy is covered by
  lunar-typescript goldens;
- `bazi-calculator` is authoritative only for its declared 2000–2099 range.

### External true-solar oracle case

File: `pkg/bazi/testdata/external_8char_uniapp_19810826.json`

Source is recorded as:

```text
axbug/8CHAR-Uni-App + api.app.yxbug.cn/api/8char/get-info
```

This is a single detailed case for a 1981 Jiamusi male birth. It records a
true-solar-time read, 四柱, hidden stems, ten gods, day-master trend, 胎元,
旬/空亡, self-sitting, 神煞, and 大运 start. License terms for that external
source should be reviewed before reuse outside this repository's regression
fixture.

### Flow golden

File: `pkg/bazi/testdata/liunian_golden.json`

- 7 cases;
- verifies a 1984 natal chart against named 流年 years, 十神, and 神煞.

## Ziwei

### Flow golden

File: `pkg/ziwei/testdata/flow_golden.json`

- 100 cases;
- verifies 流月, 流日, and flow-star placements.

### Minor-star brightness

File: `pkg/ziwei/testdata/minor_brightness_golden.json`

Source: `iztro minor_star_brightness`, baseline 2.5.8, with checks synchronized
against 2.6.1 for major-star brightness changes. The fixture explicitly records
that 文昌/文曲 did not change from 2.5.8 to 2.6.1.

## Qimen

Qimen has the densest external-reference matrix. Method metadata, provenance,
and hashes are embedded in fixtures.

| File | Source | Vectors / cases | Scope |
|---|---|---:|---|
| `chart_golden.json` | Internal canonical regression board | 1 board | Hour / 转盘 / 拆补 full board |
| `external_rotating_golden.json` | `atopx/qimen` MIT, commit `8eb06d0` | 11 | Rotating-board charts |
| `method_matrix_golden.json` | `atopx/qimen` MIT, commit `8eb06d0` | 12 | Cross-school/method matrix |
| `zhirun_golden.json` | `atopx/qimen` MIT, commit `8eb06d0` | 16 定局, 6 states, 2 day boundaries | 置闰 and state transitions |
| `mingfa_golden.json` | `potuo/feipan-qimen` MIT, commit `e4fe69d` | 1 | 鸣法飞盘 |
| `maoshan_golden.json` | `deminzhang/qimen-go` MIT, commit `4d3f58f` | 7 | 茅山定局 |
| `jinhan_golden.json` | `kentang2017/kinqimen`, MIT declared by README | 9 | 金函玉镜 |
| `quarter_golden.json` | proprietary web execution facts | 5 | 十分钟三元刻家 |
| `twelve_minute_quarter_golden.json` | Horosa AGPL-3.0 execution facts | 8 | 十二分钟十分局刻家 |
| `pattern_golden.json` | Curated classical pattern definitions | 25 | Pattern definitions |

Provenance rules recorded in the fixtures:

- MIT/MIT-derived sources include commit, source-file hash, and license hash
  where applicable;
- proprietary and AGPL sources are used as fixed execution facts only;
- no proprietary or AGPL source code or data table is copied into MingliGo.

The `kinqimen` fixture additionally records its classical domain source,
《奇门遁甲秘笈大全·诸葛武侯行兵遁甲金函玉镜卷一》.

## Liuyao

### 京房八宫 core

File: `pkg/liuyao/testdata/core_64_golden.json`

- 8 palaces × 8 hexagrams = 64 cases;
- verifies palace sequence, 本宫, 一世 through 五世, 游魂, 归魂, 世应 positions,
  and related core facts.

### Board golden

File: `pkg/liuyao/testdata/chart_golden.json`

A full sample board verifies 纳甲, 六亲, 世应, 六神, 长生, and related line facts.
Additional oracle and property tests cover conflicts, moving lines, timing,
lifecycle, and calendar boundaries.

## Huangli

File: `pkg/huangli/testdata/query_date_matrix_golden.json`

- 90 cases;
- 15 event types;
- all 12 建除 states;
- day pillar and 节月 sources: `lunar-python` 1.4.8 and `sxtwl` 2.0.7;
- rule basis: 《协纪辨方书》 建除 and 青龙黄道起例.

Dates with date-level solar-term disagreement between oracles are skipped rather
than forced into consensus.

## Bazhai

### Complete chart golden

File: `pkg/bazhai/testdata/chart_golden.json`

A full chart golden verifies 命卦, 大游年 directions, and annual stars.

### 命卦 cycle matrix

File: `pkg/bazhai/testdata/minggua_cycle_golden.json`

- 400 cases;
- 1900–2099;
- male and female cases;
- verifies 命卦 and 东四/西四 group;
- independent 命理 anchor assertions run even when regenerating golden files.

## Xuankong

File: `pkg/xuankong/testdata/flying_star_matrix_golden.json`

- rule basis: 《沈氏玄空学》下卦正向;
- 9 periods;
- 12 sitting/facing pairs per period;
- 108 combinations;
- verifies period, mountain, and facing stars.

The current implementation intentionally covers lower-plate forward direction.
Replacement-star boards are not implied by this fixture.

## Reference project index

| Name | Role in MingliGo |
|---|---|
| `lunar-python` | Calendar/pillar/lunar/建月 consensus source |
| `sxtwl` | Calendar/pillar/lunar/建月 consensus source |
| `lunar-typescript` / ShouXingUtil | Astronomical implementation basis and authoritative Bazi golden source |
| `bazi-calculator` | Limited four-pillar consensus source for declared cases/range |
| `atopx/qimen` | Rotating, method-matrix, and 置闰 golden vectors |
| `potuo/feipan-qimen` | 鸣法飞盘 golden vector |
| `deminzhang/qimen-go` | 茅山定局 golden vectors |
| `kentang2017/kinqimen` | 金函玉镜 golden vectors |
| `iztro` | Ziwei brightness baseline |
| `axbug/8CHAR-Uni-App` | External 八字 oracle case |

For links, commits, hashes, and licenses, inspect the `source` metadata inside
each fixture. Keep provenance with the fixture when adding or updating vectors.

## Regenerating goldens

Some tests support:

```bash
UPDATE_GOLDEN=1 go test ./pkg/<domain>
```

Never use regeneration to bless an unexplained diff. Before updating:

1. identify the domain rule or source-version change;
2. run independent anchor assertions;
3. inspect the full JSON diff;
4. update source version, commit, hash, and known limits;
5. run the full package and boundary tests.

## Adding a new reference

A new reference fixture should include:

- source name and URL;
- version, commit, or artifact identifier;
- source file or execution artifact;
- SHA-256 hash;
- license and license hash when available;
- exact fields used;
- known school differences and tolerances;
- rule basis;
- prohibition or permission status for copying source/data.

Execution facts from proprietary software may be used only when the fixture
explicitly records that usage and no source code or data table is copied.
