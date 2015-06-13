package hunters

import (
	"fmt"
	"tables"

	"github.com/uhhhclem/mse/src/interact"
)

type Game struct {
	interact.Game
	Boat      // see boat.go
	Patrol    // see patrol.go
	Encounter // see encounter.go
	Combat    // see combat.go
	Done      bool
}

func NewGame() *Game {
	g := &Game{
		Game: *interact.NewGame(),
	}

	g.Boat = Boat{
		Type:       tables.TypeVIIB,
		ID:         "SS-17",
		Kommandant: "Heinrich Obersdorf",
		Damage: Damage{
			Equipment: make(map[EquipmentName]DamageState),
		},
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

type stateHandler func(*Game) interact.GameState

type handlerMap map[interact.GameState]stateHandler

var (
	handlers        = make(handlerMap)
	handlerRegistry = make(chan handlerMap)
)

func init() {
	go func() {
		handlerRegistry <- map[interact.GameState]stateHandler{
			StartState:  handleStart,
			MiddleState: handleMiddle,
			EndState:    handleEnd,
		}
	}()
	go func() {
		for {
			h := <-handlerRegistry
			for k, v := range h {
				handlers[k] = v
			}
		}
	}()
}

const (
	StartState  interact.GameState = "Start"
	MiddleState                    = "Middle"
	EndState                       = "End"
)

func handleStart(g *Game) interact.GameState {
	g.Log("Start")
	g.NewPrompt("Make a choice:")
	g.AddChoice("Start", "Go to start")
	g.AddChoice("Patrol", "Go to patrol assignment")
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
	case "Patrol":
		return PatrolStartState
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

type drm struct {
	mod  int
	desc string
}

func (d drm) String() string {
	sign := " "
	if d.mod > 0 {
		sign = "+"
	}
	if d.mod < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%d %s", sign, d.mod, d.desc)
}

type gameTest func() (bool, string)

func (g *Game) getDrm(mod int, t gameTest) drm {
	result, name := t()
	if result {
		return drm{mod, name}
	}
	return drm{0, name}
}
