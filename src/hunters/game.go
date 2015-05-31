package hunters

import (
	"fmt"
	"log"
	"os"
	"strings"

	tb "tables"
	
	"github.com/uhhhclem/mse/src/interact"
)

type Game struct {
	interact.Game
	Boat   // see boat.go
	Combat // see combat.go
	PromptMessages chan *PromptMessage
	Input  chan string
	StatusMessages chan *StatusMessage
	State  State
	Done   bool
}

func NewGame() *Game {
	g := &Game{
		PromptMessages: make(chan *PromptMessage),
		Input:  make(chan string),
		StatusMessages: make(chan *StatusMessage),
		State:  Start,
	}

	g.Boat = Boat{
		Type:       "VIIC",
		ID:         "SS-17",
		Kommandant: "Heinrich Obersdorf",
	}

	return g
}

// A PromptMessage is sent to the player to get a Choice.
type PromptMessage struct {
	State string
	Message string
	Choices []string
}

// A StatusMessage is sent to the player to report on otherwise opaque processes.
type StatusMessage struct {
	State string
	Message string
}

func (g *Game) Status(msg string) {
	g.StatusMessages <- &StatusMessage{string(g.State), msg}
}

func (g *Game) Statusf(f string, a... interface{}) {
	g.Status(fmt.Sprintf(f, a))
}

// ChangeState progresses to the next state and executes its handler.
func (g *Game) ChangeState() {
	h, ok := handlers[g.State]
	if !ok {
		log.Fatalf("No handler defined for state %s", g.State)
	}
	g.State = h(g)
}

// HandleInput displays any outstanding Prompt and scans inputs until a
// valid Choice is entered.
func (g *Game) HandleInput() {
	for {
		p := <-g.PromptMessages
		if p == nil {
			break
		}
		choices := make(map[string]bool)
		for _, c := range p.Choices {
			choices[strings.ToLower(c)] = true
		}

		var c string
		for {
			fmt.Printf("%s %s : ", p.Message, p.Choices)
			if _, err := fmt.Scanf("%s", &c); err != nil {
				fmt.Print(err)
				continue
			}
			c = strings.ToLower(c)
			if choices[c] {
				break
			}
		}
		g.Input <- c
	}
}

func (g *Game) getInput(msg string, choices ...string) string {
	g.PromptMessages <- &PromptMessage{string(g.State), msg, choices}
	// TODO:  handle timeouts and reprompting
	return <-g.Input
}

func start(g *Game) State {
	c := g.getInput("Start", "middle", "end")
	switch c {
	case "end":
		return End
	case "middle":
		return Middle
	default:
		log.Fatalf("Unknown input: %s\n", c)
		return End
	}
}

func middle(g *Game) State {
	c := g.getInput("Middle", "start", "end")
	switch c {
	case "start":
		return Start
	case "end":
		return End
	default:
		log.Fatalf("Unknown input: %s\n", c)
		return End
	}
}

func end(g *Game) State {
	g.PromptMessages <- nil
	g.StatusMessages <- nil
	g.Done = true
	os.Exit(0)
	return End
}

func combatStart(g *Game) State {
	_, r := tb.EncounterChartSupplement.Roll(tb.DayOrNight)
	g.Statusf("%s rolled...", r)
	if r == tb.Night {
		return SelectRange
	}
	return SwitchToNightCheck
}

func selectRange(g *Game) State {
	ranges := []Range{CloseRange, MediumRange, LongRange}
	choices := []string{"C", "M", "L"}

	c := g.getInput("Select range", choices...)
	g.Range = ranges[getIndex(choices, c)]
	return End
}

func switchToNightCheck(g *Game) State {
	if tb.Roll1D6() < 5 {
		g.Day = false
		return SelectRange
	}
	g.Status("Switching to night failed.")
	return End
}

func getIndex(choices []string, c string) int {
	for i, choice := range choices {
		if strings.ToLower(choice) == strings.ToLower(c) {
			return i
		}
	}
	return -1
}

const (
	Start              State = "Start"
	Middle                   = "Middle"
	End                      = "End"
	CombatStart              = "CombatStart"
	SelectRange              = "SelectRange"
	SwitchToNightCheck       = "SwitchToNightCheck"
)

type State string

type stateHandler func(*Game) State

var handlers map[State]stateHandler

func init() {
	handlers = map[State]stateHandler{
		Start:              start,
		Middle:             middle,
		End:                end,
		CombatStart:        combatStart,
		SelectRange:        selectRange,
		SwitchToNightCheck: switchToNightCheck,
	}
}
