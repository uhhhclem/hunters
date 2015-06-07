package hunters

import (
	"github.com/uhhhclem/mse/src/interact"
	"tables"
)

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

const (
	CombatStartState    interact.GameState = "CombatStart"
	CombatShipSizeState                    = "CombatShipSize"
)

func init() {
	go func() {
		handlerRegistry <- handlerMap{
			CombatStartState:    handleCombatStart,
			CombatShipSizeState: handleCombatShipSize,
		}
	}()
}

func handleCombatStart(g *Game) interact.GameState {
	c := &g.Combat
	roll := tables.Roll1D6()
	c.Day = roll < 4
	g.Logf("Roll for Day/Night: %d - %s", roll,
		map[bool]string{true: "Day", false: "Night"}[c.Day])

	return CombatShipSizeState
}

func handleCombatShipSize(g *Game) interact.GameState {
	g.Log(CombatShipSizeState)
	return EndState
}
