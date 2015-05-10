package main

import (
	"fmt"
	"math/rand"
	"time"
	
	"hunters"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	g := hunters.NewGame()
	g.LoadTestData()
	g.State = hunters.CombatStart
	
	go g.HandleInput()
	//go g.HandleStatus()

	for {
		printStatus: for {
			select {
				case s := <- g.Status:
				  fmt.Println(s)
				default:
				  break printStatus
			}
		}
		g.State = g.State(g)
	}
}
