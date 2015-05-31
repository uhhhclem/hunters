package hunters

import (
	"log"

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
		if s == EndState {
			g.NextStatus <- nil
			g.NextPrompt <- nil
			return
		}
		g.State = handlers[s](g)
		g.Ready <- true
	}
}

const (
	StartState  interact.GameState = "Start"
	MiddleState                    = "Middle"
	EndState                       = "End"
)

type stateHandler func(*Game) interact.GameState

var handlers map[interact.GameState]stateHandler

func init() {
	handlers = map[interact.GameState]stateHandler{
		StartState:  handleStart,
		MiddleState: handleMiddle,
		EndState:    handleEnd,
	}
}

func handleStart(g *Game) interact.GameState {
	g.Log("Start")
	g.NewPrompt("Make a choice:")
	g.AddChoice("Start", "Go to start")
	g.AddChoice("End", "Go to end")
	g.SendPrompt()

	return MiddleState
}

func handleMiddle(g *Game) interact.GameState {
	g.Log("Middle")
	c := <-g.NextChoice
	g.Logf("You chose %s: %s", c.Key, c.Name)
	if c.Key == "Start" {
		return StartState
	}
	return EndState
}

func handleEnd(g *Game) interact.GameState {
	return EndState
}
