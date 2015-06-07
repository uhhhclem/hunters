package hunters

import (
	t "tables"

	"github.com/uhhhclem/mse/src/interact"
)

const (
	EncounterStartState  interact.GameState = "EncounterStart"
	EncounterEndState                       = "EncounterEnd"
	EncounterCombatState                    = "EncounterCombat"
)

func init() {
	go func() {
		handlerRegistry <- map[interact.GameState]stateHandler{
			EncounterStartState:  handleEncounterStart,
			EncounterEndState:    handleEncounterEnd,
			EncounterCombatState: handleEncounterCombat,
		}
	}()
}

type Encounter struct {
}

func handleEncounterStart(g *Game) interact.GameState {
	p := &g.Patrol
	tn := p.TravelBoxes[p.TravelBox]

	g.Logf("Determining encounter for travel box %d (%s)", p.TravelBox+1, tn)
	roll, result := t.EncounterChart.Roll(t.TableName(tn))
	s := result
	if s == "" {
		s = "no encounter"
	}
	g.Logf("Rolling on Encounter Chart [E1]:  rolled %d (%s)", roll, s)

	if result == "" {
		return EncounterEndState
	}
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
	//TODO gotta implement this too
	g.Log(string(EncounterCombatState))
	return EncounterEndState
}
