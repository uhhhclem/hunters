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
	ToSink    int
	Torpedoes []Torpedo
}

type Range string

const (
	CloseRange  Range = "Close"
	MediumRange Range = "Medium"
	LongRange   Range = "Long"
)

func toSink(t Target) int {
	switch t.Type {
	case tables.SmallFreighter:
		return 2
	case tables.LargeFreighter, tables.Tanker:
		if t.Tonnage < 10000 {
			return 3
		}
		return 4
	}
	return 5
}
