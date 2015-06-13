package hunters

import (
	"fmt"

	t "tables"

	"github.com/uhhhclem/mse/src/interact"
)

const (
	EncounterStartState       interact.GameState = "EncounterStart"
	EncounterEndState                            = "EncounterEnd"
	EncounterCombatState                         = "EncounterCombat"
	EncounterAircraftState                       = "EncounterAircraft"
	EncounterRandomEventState                    = "EncounterRandomEvent"
)

func init() {
	go func() {
		handlerRegistry <- map[interact.GameState]stateHandler{
			EncounterStartState:       handleEncounterStart,
			EncounterEndState:         handleEncounterEnd,
			EncounterCombatState:      handleEncounterCombat,
			EncounterAircraftState:    handleEncounterAircraft,
			EncounterRandomEventState: handleEncounterRandomEvent,
		}
	}()
}

type Encounter struct {
	Type string
}

func handleEncounterStart(g *Game) interact.GameState {
	p := &g.Patrol
	e := &g.Encounter
	tn := p.TravelBoxes[p.TravelBox]

	g.Logf("Determining encounter for travel box %d (%s)", p.TravelBox+1, tn)
	// TODO:  this seems hacky
	if tn == t.BritishIslesM {
		tn = t.BritishIsles
	}
	roll, result := t.EncounterChart.Roll(t.TableName(tn))

	e.Type = string(result)
	if roll == 12 && !p.Natural12EncounterRoll {
		e.Type = t.RandomEvent
		p.Natural12EncounterRoll = true
	}

	enc := e.Type
	if e.Type == "" {
		enc = "no encounter"
	}
	g.Logf("Rolling on Encounter Chart [E1]:  rolled %d (%s)", roll, enc)

	return EncounterCombatState
}

func handleEncounterEnd(g *Game) interact.GameState {
	g.Log(string(EncounterEndState))
	p := &g.Patrol
	p.TravelBox++
	if p.TravelBox < len(p.TravelBoxes) {
		return EncounterStartState
	}
	// TODO:  return to port
	return EndState
}

func handleEncounterCombat(g *Game) interact.GameState {
	e := &g.Encounter
	g.Combat = Combat{}
	c := &g.Combat

	switch e.Type {
	case "":
		return EncounterEndState
	case t.Aircraft:
		return EncounterAircraftState
	case t.CapitalShip:
		c.Escort = true
		g.rollTargetFromRoster(t.CapitalShipTargetRoster)
	case t.Convoy:
		c.Escort = true
		for i := 0; i < 4; i++ {
			g.rollTarget()
		}
	case t.RandomEvent:
		return EncounterRandomEventState
	case t.Ship:
		g.rollTarget()
	case t.ShipPlusEscort:
		c.Escort = true
		g.rollTarget()
	case t.TwoShips:
		g.rollTarget()
		g.rollTarget()
	case t.TwoShipsPlusEscort:
		c.Escort = true
		g.rollTarget()
		g.rollTarget()
	case t.Tanker:
		g.rollTargetFromRoster(t.TankerTargetRoster)
	default:
		g.Logf("Unknown encounter type: %s", e.Type)
		return EndState
	}

	return CombatStartState
}

func (g *Game) rollTargetFromRoster(tr t.TargetRoster) {
	roll, target := tr.Roll()
	g.Logf("Roll for target: %d (%s)", roll, target)
	c := &g.Combat

	ct := Target{
		Target: target,
		Number: len(c.Targets) + 1,
	}
	c.Targets = append(c.Targets, ct)
}

func (g *Game) rollTarget() {
	roll, size := t.EncounterChartSupplement.Roll(t.ShipSize)
	g.Logf("Roll for ship size: %d (%s)", roll, size)

	switch size {
	case t.SmallFreighter:
		g.rollTargetFromRoster(t.SmallFreighterTargetRoster)
	case t.LargeFreighter:
		g.rollTargetFromRoster(t.LargeFreighterTargetRoster)
	case t.Tanker:
		g.rollTargetFromRoster(t.TankerTargetRoster)
	}

}

func handleEncounterAircraft(g *Game) interact.GameState {
	g.Log(string(EncounterAircraftState))
	one := t.Result("1 Attack + 1 Crew Injury")
	two := t.Result("2 Attacks + 1 Crew Injury")
	none := t.Result("Successful crash dive (no air attack)")
	tb := t.Table{
		1: two,
		2: one,
		3: one,
		4: one,
		5: one,
		6: none,
	}
	drms := []t.DRM{
		g.getDrm(-1, g.testAllCrewKIAOrSW),
		g.getDrm(-1, g.testGreenCrew),
		g.getDrm(1, g.testEliteCrew),
		g.getDrm(-1, g.testEquipmentDamaged(DivePlanes)),
		g.getDrm(-1, g.testTypeVIIDOrIX),
		g.getDrm(1, g.testYear(1939)),
		g.getDrm(-1, g.testYear(1942)),
		g.getDrm(-2, g.testYear(1943)),
		g.getDrm(-1, g.testMOrAPatrolAssignment),
	}
	natural, modified := t.Roll2D6WithDRMs(drms, 1, t.NoMaxRoll)
	result, ok := tb[modified]
	if !ok {
		result = none
	}
	g.LogWithDRMs("Aircraft encounter ", drms)
	g.Logf("Rolled %d (modified to %d):  %s", natural, modified, result)

	switch result {
	case none:
		return EncounterEndState
	case one:
		g.Combat.AircraftAttacks = 1
	case two:
		g.Combat.AircraftAttacks = 2
	}

	return CombatAircraftAttackState
}

//TODO
func (g *Game) testAllCrewKIAOrSW() (bool, string) {
	return false, "All crew KIA or SW"
}

//TODO
func (g *Game) testGreenCrew() (bool, string) {
	return false, "Green crew"
}

func (g *Game) testEliteCrew() (bool, string) {
	//TODO
	return false, "Elite Crew"
}

func (g *Game) testTypeVIIDOrIX() (bool, string) {
	return g.Boat.Type == t.TypeVIID || g.Boat.Type == t.TypeIX, "Type VIID or IX"
}

func (g *Game) testYear(year int) gameTest {
	return func() (bool, string) {
		return g.Year == year, fmt.Sprintf("%d", year)
	}
}

//TODO
func (g *Game) testMOrAPatrolAssignment() (bool, string) {
	return false, "(M) or (A) patrol assignment"
}

func handleEncounterRandomEvent(g *Game) interact.GameState {
	// TODO: implement this
	g.Log(string(EncounterRandomEventState))
	return EncounterEndState
}
