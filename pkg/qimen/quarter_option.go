package qimen

import "encoding/json"

// DunSource is the qimen-domain symbol for how 遁阴阳 is selected in the
// twelve-minute quarter method.
type DunSource uint8

const (
	dsHourBranch DunSource = iota // hour_branch
	dsSolarTerm                    // solar_term
)

var dunSourceNames = [...]string{"hour_branch", "solar_term"}

func (d DunSource) String() string {
	if int(d) < len(dunSourceNames) {
		return dunSourceNames[d]
	}
	return ""
}

func (d DunSource) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

func ParseDunSource(name string) (DunSource, bool) {
	for i, n := range dunSourceNames {
		if n == name {
			return DunSource(i), true
		}
	}
	return 0, false
}

// HourBoundary is the qimen-domain symbol for the hour-start boundary used by
// the twelve-minute quarter method.
type HourBoundary uint8

const (
	hbLateZi  HourBoundary = iota // late_zi
	hbZiZheng                     // zi_zheng
)

var hourBoundaryNames = [...]string{"late_zi", "zi_zheng"}

func (h HourBoundary) String() string {
	if int(h) < len(hourBoundaryNames) {
		return hourBoundaryNames[h]
	}
	return ""
}

func (h HourBoundary) MarshalJSON() ([]byte, error) { return json.Marshal(h.String()) }

func ParseHourBoundary(name string) (HourBoundary, bool) {
	for i, n := range hourBoundaryNames {
		if n == name {
			return HourBoundary(i), true
		}
	}
	return 0, false
}
