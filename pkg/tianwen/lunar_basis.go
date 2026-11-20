package tianwen

import (
	"time"
)

// LunarProjector combines the basis choice with the date-assignment reference.
// It is the only place where a timezone participates in lunar conversion.
type LunarProjector struct {
	basis LunarBasis
	ref   *LunarReference
	loc   *time.Location
}

// NewLunarProjector builds a projector. A nil location falls back to UTC+8.
func NewLunarProjector(basis LunarBasis, loc *time.Location) (LunarProjector, error) {
	parsed, err := ParseLunarBasis(string(basis))
	if err != nil {
		return LunarProjector{}, err
	}
	if loc == nil {
		loc = cstLocation
	}
	ref := ShanghaiLunarReference()
	if parsed == LunarBasisLocalAstronomical {
		ref = NewLunarReference(loc)
	}
	return LunarProjector{basis: parsed, ref: ref, loc: loc}, nil
}

// LocalLabelProjector returns the official UTC+8 calendar keyed by the local
// civil date label. It is the fallback for date-based callers without an
// explicit basis.
func LocalLabelProjector() LunarProjector {
	return LunarProjector{
		basis: LunarBasisLocalLabel,
		ref:   ShanghaiLunarReference(),
		loc:   cstLocation,
	}
}

// SelectorDate is the civil date used to look up the lunar date.
func (p LunarProjector) SelectorDate(at time.Time) Date {
	if p.basis == LunarBasisOfficial {
		return BeijingDateAt(at)
	}
	return StandardDateAt(at, p.loc)
}

// SelectorInstant is a representative instant for a selector date in the
// basis's date frame.
func (p LunarProjector) SelectorInstant(date Date) time.Time {
	if p.basis == LunarBasisOfficial {
		return date.Time()
	}
	return StandardNoon(date, p.loc)
}

// DayStart is midnight on the civil date in the basis's date frame.
func (p LunarProjector) DayStart(date Date) time.Time {
	if p.basis == LunarBasisOfficial {
		return time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, cstLocation)
	}
	return StandardMidnight(date, p.loc)
}

// SolarDateToLunar converts an already-selected civil date.
func (p LunarProjector) SolarDateToLunar(date Date) LunarDate {
	return p.ref.SolarDateToLunar(date)
}

// LunarAt projects an instant to the lunar date under this basis.
func (p LunarProjector) LunarAt(at time.Time) LunarDate {
	return p.ref.SolarDateToLunar(p.SelectorDate(at))
}

// DateFromLunar converts a lunar date under this basis.
func (p LunarProjector) DateFromLunar(ld LunarDate) (Date, error) {
	return p.ref.DateFromLunar(ld)
}

// ValidateLunarDate validates a lunar date under this basis.
func (p LunarProjector) ValidateLunarDate(date LunarDate) (LunarMonthInfo, error) {
	return p.ref.ValidateLunarDate(date)
}

// NextLunarMonth returns the next lunation under this basis.
func (p LunarProjector) NextLunarMonth(year, month int, leap bool) (LunarMonthInfo, error) {
	return p.ref.NextLunarMonth(year, month, leap)
}
