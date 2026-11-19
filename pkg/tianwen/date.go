package tianwen

import (
	"fmt"
	"time"
)

// Date is a date on the Gregorian calendar, without a time of day or an
// explicit timezone. It is the date-level selector used by flow and almanac
// capabilities. Calendar derivations that need a reference instant use noon
// China Standard Time (see Time).
type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

// NewDate returns the civil date identified by year/month/day. Callers
// that accept untrusted input must call Validate.
func NewDate(year, month, day int) Date {
	return Date{Year: year, Month: month, Day: day}
}

// DateOf returns the calendar date of t in t's location.
func DateOf(t time.Time) Date {
	year, month, day := t.Date()
	return Date{Year: year, Month: int(month), Day: day}
}

// BeijingDateAt returns the Beijing civil date of an instant.
func BeijingDateAt(at time.Time) Date {
	return DateOf(at.In(cstLocation))
}

// TimeIn returns the civil date at 12:00 in the given location.
func (d Date) TimeIn(loc *time.Location) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 12, 0, 0, 0, loc)
}

// SolarTermDate returns the Chinese civil date of a solar-term boundary.
// Solar terms are absolute instants; the date-level projection uses the
// China Standard Time civil date.
func SolarTermDate(year int, longitude float64) Date {
	return DateOf(SolarTermTime(year, longitude).In(cstLocation))
}

// Validate checks that the date exists in the Gregorian calendar. It rejects
// normalization (for example 2026-02-30 must not silently become 2026-03-02).
func (d Date) Validate() error {
	if d.Year < 1 || d.Year > 9999 || d.Month < 1 || d.Month > 12 || d.Day < 1 || d.Day > 31 {
		return fmt.Errorf("invalid civil date %04d-%02d-%02d", d.Year, d.Month, d.Day)
	}
	t := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
	if t.Year() != d.Year || int(t.Month()) != d.Month || t.Day() != d.Day {
		return fmt.Errorf("invalid civil date %04d-%02d-%02d", d.Year, d.Month, d.Day)
	}
	return nil
}

// Time returns a representative instant for the date: 12:00 China Standard
// Time. It is not the day boundary; callers use it for calendar derivations
// that need an instant well inside the civil date.
func (d Date) Time() time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 12, 0, 0, 0, cstLocation)
}

// AddDays returns the date n calendar days after d. n may be negative.
func (d Date) AddDays(n int) Date {
	t := d.Time().AddDate(0, 0, n)
	return DateOf(t)
}

// Before reports whether d is before other in civil date order.
func (d Date) Before(other Date) bool {
	if d.Year != other.Year {
		return d.Year < other.Year
	}
	if d.Month != other.Month {
		return d.Month < other.Month
	}
	return d.Day < other.Day
}

// After reports whether d is after other in civil date order.
func (d Date) After(other Date) bool {
	return other.Before(d)
}

// IsZero reports whether d is the zero value.
func (d Date) IsZero() bool {
	return d == Date{}
}

// String formats the date as YYYY-MM-DD.
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}
