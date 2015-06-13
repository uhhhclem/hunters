package hunters

import (
	t "tables"

	"github.com/uhhhclem/mse/src/interact"
)

type Combat struct {
	AircraftAttacks int
	Day             bool
	Surface         bool
	Range           Range
	Escort          bool
	Targets         []Target
}

type Target struct {
	t.Target
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

func toSink(tg Target) int {
	switch tg.Type {
	case t.SmallFreighter:
		return 2
	case t.LargeFreighter, t.Tanker:
		if tg.Tonnage < 10000 {
			return 3
		}
		return 4
	}
	return 5
}

const (
	CombatStartState             interact.GameState = "CombatStart"
	CombatShipSizeState                             = "CombatShipSize"
	CombatAircraftAttackState                       = "CombatAircraftAttack"
	CombatAircraftAttackEndState                    = "CombatAircraftAttack"
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
	roll := t.Roll1D6()
	c.Day = roll < 4
	g.Logf("Roll for Day/Night: %d - %s", roll,
		map[bool]string{true: "Day", false: "Night"}[c.Day])

	return CombatShipSizeState
}

func handleCombatShipSize(g *Game) interact.GameState {
	g.Log(CombatShipSizeState)
	return EndState
}

func handleCombatAircraftAttack(g *Game) interact.GameState {
	g.Logf(CombatAircraftAttackState)

	if damaged, _ := g.testEquipmentDamaged(FlakGun)(); !damaged {
		shotDown := t.Result("Shot Down")
		damaged := t.Result("Damaged")
		miss := t.Result("Miss")
		tb := t.Table{
			3: shotDown,
			4: damaged,
			5: damaged,
			6: miss,
		}
		drms := []t.DRM{
			g.getDrm(1, g.testTypeVIIA),
			g.getDrm(-1, g.testTypeIX), //TODO: deal with 2 flak guns on type IX
			g.getDrm(-1, g.testVeteranCrew),
			g.getDrm(-1, g.testEliteCrew),
			g.getDrm(-2, g.testTypeVIICFlak),
		}
		natural, modified := t.Roll2D6WithDRMs(drms, t.NoMinRoll, t.NoMaxRoll)
		result, ok := tb[modified]
		if !ok {
			result = miss
			if modified < 3 {
				result = shotDown
			}
		}
		g.LogWithDRMs("Flak Attack vs. Aircraft", drms)
		g.Logf("Rolled %d (modified to %d): %s", natural, modified, result)

		if result == shotDown {
			g.Combat.AircraftAttacks = 0
			return EncounterEndState
		}
		if result == damaged {
			return CombatAircraftAttackEndState
		}
	}

	g.RollEscortAirAttackDamage()
	if g.Boat.Damage.Sunk {
		return EndState // TODO
	}

	return CombatAircraftAttackEndState
}

func handleCombatAircraftAttackEnd(g *Game) interact.GameState {
	g.Combat.AircraftAttacks -= 1
	if g.Combat.AircraftAttacks > 0 {
		return CombatAircraftAttackState
	}
	return EncounterEndState
}

func (g *Game) testTypeVIIA() (bool, string) {
	return g.Boat.Type == t.TypeVIIA, "Type VIIA"
}

func (g *Game) testTypeVIICFlak() (bool, string) {
	return g.Boat.Type == t.TypeVIICFlak, "Type VIIC Flak"
}

func (g *Game) testTypeIX() (bool, string) {
	return g.Boat.Type == t.TypeIX, "Type IX"
}

func (g *Game) testVeteranCrew() (bool, string) {
	return false, "Veteran crew" //TODO
}
