package hunters

import (
	"github.com/uhhhclem/mse/src/interact"
)

const (
	EncounterStartState interact.GameState = "EncounterStart"
)

func init() {
	go func() {
		handlerRegistry <- map[interact.GameState]stateHandler{
			EncounterStartState: handleEncounterStart,
		}
	}()
}

type Encounter struct {
}

func handleEncounterStart(g *Game) interact.GameState {
	g.Log(string(EncounterStartState))
	return CombatStartState
}
