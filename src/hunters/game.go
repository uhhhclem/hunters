package hunters

import (
	"github.com/uhhhclem/mse/src/interact"
)

type Game struct {
	interact.Game
	Boat   // see boat.go
	Combat // see combat.go
	Done   bool
}

func NewGame() *Game {
	g := &Game{
		Game: *interact.NewGame(),
	}

	g.Boat = Boat{
		Type:       "VIIC",
		ID:         "SS-17",
		Kommandant: "Heinrich Obersdorf",
	}

	g.Game.State = StartState

	return g
}

func (g *Game) Run() {
	for {
		s := g.State
		g.State = handlers[s](g)
		if g.Done {
			g.NextStatus <- nil
			g.NextPrompt <- nil
			g.Ready <- true
			return
		}
		g.Ready <- true
	}
}

const (
	StartState       interact.GameState = "Start"
	MiddleState                         = "Middle"
	EndState                            = "End"
	CombatStartState                    = "CombatStart"
)

type stateHandler func(*Game) interact.GameState

var handlers map[interact.GameState]stateHandler

func init() {
	handlers = map[interact.GameState]stateHandler{
		StartState:       handleStart,
		MiddleState:      handleMiddle,
		EndState:         handleEnd,
		CombatStartState: handleCombatStart,
	}
}

func handleStart(g *Game) interact.GameState {
	g.Log("Start")
	g.NewPrompt("Make a choice:")
	g.AddChoice("Start", "Go to start")
	g.AddChoice("Combat", "Go to combat")
	g.AddChoice("End", "Go to end")
	g.SendPrompt()

	return MiddleState
}

func handleMiddle(g *Game) interact.GameState {
	c := <-g.NextChoice
	switch c.Key {
	case "Start":
		return StartState
	case "Combat":
		return CombatStartState
	}
	return EndState
}

func handleEnd(g *Game) interact.GameState {
	g.Log("End of game.")
	g.Done = true
	return EndState
}

func handleCombatStart(g *Game) interact.GameState {
	g.Log("Got to Combat.")
	return StartState
}
