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
	g.State = &hunters.CombatStart{g}
	
	go g.HandleInput()
	//go g.HandleStatus()

	for {
		g.Dump()
		select {
			case s := <- g.Status:
			  fmt.Println(s)
			default:
			  break
		}
		if g.State == nil {
			time.Sleep(50 * time.Millisecond)
			break
		}
		g.State = g.State.Handle()
	}
}
