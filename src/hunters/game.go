package hunters

import (
	"fmt"
	"log"
	"os"
	"strings"

	tb "tables"
)

type Game struct {
	Boat   // see boat.go
	Combat // see combat.go
	Output chan *Prompt
	Input  chan string
	Status chan string
	State  State
	Done   bool
}

type State func(*Game) State

func NewGame() *Game {
	g := &Game{
		Output: make(chan *Prompt),
		Input:  make(chan string),
		Status: make(chan string, 1),
		State:  Start,
	}

	g.Boat = Boat{
		Type:       "VIIC",
		ID:         "SS-17",
		Kommandant: "Heinrich Obersdorf",
	}

	return g
}

// A Prompt is sent to the player to get a Choice.
type Prompt struct {
	Message string
	Choices []string
}

func (g *Game) HandleStatus() {
	for {
		s := <-g.Status
		if s == "" || s == "End" {
			return
		}
		fmt.Println(s)
	}
}

// HandleInput displays any outstanding Prompt and scans inputs until a
// valid Choice is entered
func (g *Game) HandleInput() {
	for {
		p := <-g.Output
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

func (g *Game) GetInput(msg string, choices ...string) string {
	g.Output <- &Prompt{msg, choices}
	return <-g.Input
}

func Start(g *Game) State {
		c := g.GetInput("Start", "middle", "end")
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

func Middle(g *Game) State {
		c := g.GetInput("Middle", "start", "end")
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

func End(g *Game) State {
		g.Output <- nil
		g.Status <- "End"
		g.Done = true
		os.Exit(0)
		return End
}

func CombatStart(g *Game) State {
	_, r := tb.EncounterChartSupplement.Roll(tb.DayOrNight)
	g.Status <- fmt.Sprintf("%s rolled...", r)
	if r == tb.Night {
		return SelectRange
	}
	return SwitchToNightCheck
}

func SelectRange(g *Game) State {
	ranges := []Range{CloseRange, MediumRange, LongRange}
	choices := []string{"C", "M", "L"}

	c := g.GetInput("Select range", choices...)
	g.Range = ranges[getIndex(choices, c)]
	return End
}

func SwitchToNightCheck(g *Game) State {
	if tb.Roll1D6() < 5 {
		g.Day = false
		return SelectRange
	}
	g.Status <- "Switching to night failed."
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
