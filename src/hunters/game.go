package hunters

import (
    "fmt"
    "log"
    "strings"
    
    tb "tables"
)

type Game struct {
    Boat        // see boat.go
    Combat      // see combat.go
    Output chan *Prompt
    Input chan string
    Status chan string
    State State
    Done bool
}

func NewGame() *Game {
    g := &Game{
    	Output: make(chan *Prompt),
		Input:  make(chan string),
		Status: make(chan string, 1),
	}
	g.State = &Start{Game: g}
	
	g.Boat = Boat{
	    Type: "VIIC",
	    ID: "SS-17",
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
        s := <- g.Status
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
        p := <- g.Output
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

type State interface {
    Handle() State
}

type Start struct {
    *Game
}

func (g *Game) GetInput(msg string, choices... string) string {
    g.Output <- &Prompt{msg, choices}
    return <- g.Input
}

func (s *Start) Handle() State {
    c := s.Game.GetInput("Start", "start", "middle", "end")
    switch c {
        case "start":
          return &Start{s.Game}
        case "end":
          return &End{s.Game}
        case "middle":
          return &Middle{s.Game}
        default:
          log.Fatalf("Unknown input: %s\n", c)
          return nil
    }
}

type Middle struct {
    *Game
}

func (s *Middle) Handle() State {
    c := s.Game.GetInput("Middle", "start", "end")
    switch c {
        case "start":
          return &Start{s.Game}
        case "end":
          return &End{s.Game}
        default:
          log.Fatalf("Unknown input: %s\n", c)
          return nil
    }
}

type End struct{
    *Game
}

func (s *End) Handle() State {
    s.Output <- nil
    s.Status <- "End"
    s.Done = true
    return nil
}

type CombatStart struct {*Game}

func (s *CombatStart) Handle() State {
    _, r := tb.EncounterChartSupplement.Roll(tb.DayOrNight)
    if r == tb.Night {
        return &SelectRange{s.Game}
    }
    return &SwitchToNightCheck{s.Game}
}

type SelectRange struct {*Game}

func (s *SelectRange) Handle() State {
    ranges := []Range{CloseRange, MediumRange, LongRange}
    choices := []string{"C", "M", "L"}
    
    g := s.Game
    c := g.GetInput("Select range", choices...)
    g.Range = ranges[getIndex(choices, c)]
    return &End{g}
}

type SwitchToNightCheck struct {*Game}

func (s *SwitchToNightCheck) Handle() State {
    if tb.Roll1D6() < 5 {
        s.Day = false
        return &SelectRange{s.Game}
    }
    s.Status <- "Switching to night failed."
    return &End{s.Game}
}

func getIndex(choices []string, c string) int {
    for i, choice := range choices {
        if strings.ToLower(choice) == strings.ToLower(c) {
            return i
        }
    }
    return -1
}