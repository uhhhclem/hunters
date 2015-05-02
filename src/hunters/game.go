package hunters

import (
    "fmt"
    "log"
    "strings"
)

type Game struct {
    Output chan *Prompt
    Input chan string
    State State
    Done bool
}

type Prompt struct {
    Message string
    Choices []string
}

// HandleInput displays any outstanding prompt and scans inputs until a
// valid choice is entered
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
    s.Done = true
    fmt.Println("End")
    return nil
}
