package hunters

import (
	"strings"

	"tables"

	"github.com/uhhhclem/mse/src/interact"
)

type Game struct {
	interact.Game
	Year      int
	Month     int
	Boat      // see boat.go
	Patrol    // see patrol.go
	Encounter // see encounter.go
	Combat    // see combat.go
	Done      bool
}

func NewGame() *Game {
	g := &Game{
		Game:  *interact.NewGame(),
		Year:  1939,
		Month: 9,
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

type gameTest func() (bool, string)

func (g *Game) getDrm(mod int, t gameTest) tables.DRM {
	result, name := t()
	if result {
		return tables.DRM{mod, name}
	}
	return tables.DRM{0, name}
}

func (g *Game) LogWithDRMs(desc string, drms []tables.DRM) {
	app := make([]string, 0)
	for _, d := range drms {
		if d.Mod == 0 {
			continue
		}
		app = append(app, d.String())
	}
	if len(app) > 0 {
		g.Logf("%s (Mods: %s)", desc, strings.Join(app, ","))
		return
	}
	g.Logf("%s (Mods: none)", desc)
}
