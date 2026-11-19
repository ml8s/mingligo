# MingliGo time model

Time is the foundation of Chinese metaphysics computation. MingliGo treats
calendar, astronomy, and school conventions as separate, explicit facts instead
of collapsing them into one opaque timestamp.

## The main conversion chain

```text
civil / wall-clock time
  → standard local time (DST removed)
    → longitude correction
      → equation of time
        → local apparent solar time (LAST, 真太阳时)
          → 时辰, solar day, solar-term boundary, year/month pillars
```

This chain is shared by domain packages. Domain APIs normally expose
`Compute*At(standard time.Time, longitude float64, ...)` so callers do not have
to implement astronomy themselves.

The direct `Compute*At` APIs expect the caller to provide a standard-time
`time.Time`. When input may be a DST wall clock, a vague 时辰, a range, or a
longitude range, call `ResolveStandard` first; it performs DST removal and
boundary resolution before handing a standard instant to a chart API.

## Civil time and standard time

A civil timestamp may include DST. A 命理 chart should use the jurisdiction's
standard offset for that instant, not the temporary DST offset.

`tianwen.StandardLocalTime` projects an instant to the zone's non-DST standard
offset. If the instant itself is in DST, nearby instants are probed to recover
the zone's ordinary offset.

Implementation properties:

- Timezone names must be IANA names such as `Asia/Shanghai`; aliases such as
  `CST` are rejected.
- The IANA database is embedded through `time/tzdata`, so historical offsets do
  not depend on the host OS.
- Historical non-DST offsets, such as Asia/Kuala_Lumpur before 1982, are handled
  by the timezone database.

## True solar time

Chinese local apparent solar time (LAST, 真太阳时) is derived from two
corrections:

```text
LAST minutes =
    local wall-clock minutes
  + 4 × (longitude − timezone standard meridian)
  + equation of time
```

Where:

- longitude is degrees east, from `-180` to `180`;
- the timezone standard meridian is the standard UTC offset in hours × 15;
- the equation of time uses the Meeus equation-of-time algorithm.

For example, at UTC+8 the standard meridian is `120°E`. At longitude `87.6°E`,
the longitude correction is about `-129.6` minutes. Local noon can therefore
occur near solar `09:50`. Conversely, at longitude `134°E`, `23:30` civil time
can project past local solar midnight into the next solar day.

### 时辰

The twelve 时辰 are computed on LAST, not on administrative local time:

```text
branch index = floor(((LAST minutes + 60) mod 1440) / 120)
```

Thus 子时 spans `23:00–01:00` in local solar time. Because the projection may
cross midnight, the solar day used for 日柱 can differ from the civil calendar
date.

## Solar terms and year/month boundaries

Solar terms are absolute astronomical instants. MingliGo does not treat a 节气
as a whole civil date.

- 年柱 uses the exact 立春 instant.
- 月柱 uses the exact 节 boundary.
- 日柱 and 时柱 use local apparent solar time.

This separation matters near boundaries: a few minutes around 立春 can change
both the year pillar and month pillar, while longitude changes the day and hour
dimensions.

## Lunar projection bases

The Chinese lunisolar calendar can be keyed in different ways. MingliGo exposes
the basis instead of silently choosing one:

| Basis | Meaning | Use |
|---|---|---|
| `official` | Beijing civil date with the UTC+8 calendar | Official Chinese calendar alignment |
| `local_label` | Local civil date label, but UTC+8 calendar table | Compatibility with date-label school conventions |
| `local_astronomical` | Local standard date, with new moons and 中气 assigned in that zone | Globally aligned locality; default for new APIs |

`tianwen.NewLunarProjector(basis, location)` is the only place where a timezone
participates in lunar conversion. The projector converts civil dates, validates
lunar months and leap months, and projects an instant to a lunar date.

## Ambiguous birth-time resolution

Historical records often give only a date, a 时辰, or a range. The input may also
be close to a boundary or the longitude of a large administrative area may span
more than one 时辰.

`tianwen.ResolveStandard` accepts:

- solar date plus clock time;
- solar date plus true-solar 时辰;
- solar date plus a same-day time range;
- lunar date plus a clock time.

It returns candidates with reasons such as:

- user input;
- 时辰 boundary;
- solar-day boundary;
- solar-term crossing;
- historical time system or DST.

It deliberately refuses to collapse ambiguity into one silent answer. A large
longitude range can legitimately produce more than one candidate chart time.

## Precision and current limits

### Astronomical instants

Solar-longitude computation uses a Julian day with hours, minutes, seconds, and
fractional seconds. The implementation derives from the ShouXing / lunar
VSOP87D solar-longitude chain and includes nutation, aberration, and ΔT.

Tests compare known 2024 solar terms with a tolerance of ±30 seconds and a known
new moon with a tolerance of ±2 minutes. Do not extrapolate these tolerances
into a blanket claim for every historical date.

### True-solar-time wall value

The equation-of-time implementation targets sub-minute precision, but the
current public `ApparentSolarFromCivil` wall-time conversion uses the input's
hour and minute and writes zero seconds. This is sufficient for 时辰
classification, whose cells are two hours wide, but callers must not treat the
returned LAST value as a second-level wall-time representation.

Changing this contract is a compatibility change: the timestamp fields, golden
JSON, and boundary tests must be updated together.

### Calendar range

Primary domain tests emphasize 1900–2100. Some astronomical formulas can be
evaluated more broadly, but historical calendar behavior, supported tables, and
external-reference consensus should not be assumed outside the declared primary
range.

## Caller rules

1. Keep the user's civil timestamp and timezone separate from the longitude.
2. Use an IANA timezone name; never use an ambiguous abbreviation.
3. Pass longitude east-positive, including `0` for Greenwich.
4. Store the chosen lunar basis and school policy alongside the chart.
5. Use `ResolveStandard` when an input date/时辰/longitude range can cross a
   boundary.
6. Persist enough policy metadata to reproduce the chart.

## Common pitfalls

- Treating Beijing time as universal solar time.
- Applying DST to a 命理 standard-time chart.
- Replacing longitude `0` with a default city longitude.
- Comparing LAST wall time with civil wall time as though they were the same
  scale.
- Using UTC+8 lunar dates for every locality when the caller chose a local
  astronomical basis.
- Assuming a whole administrative region has one 时辰.
