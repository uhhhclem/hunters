package hunters

import "tables"

type Combat struct {
	Day     bool
	Surface bool
	Range   Range
	Escort  bool
	Targets []Target
}

type Target struct {
	tables.Target
	Number    int
	Damage    int
	Torpedoes []Torpedo
}

type Range string

const (
	CloseRange  Range = "Close"
	MediumRange Range = "Medium"
	LongRange   Range = "Long"
)
