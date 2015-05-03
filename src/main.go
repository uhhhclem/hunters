package main

import (
	"math/rand"
	"time"
	
	"hunters"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	g := hunters.NewGame()
	g.LoadTestData()
	
	go g.HandleInput()

	for {
		g.Dump()
		g.State = g.State.Handle()
		if g.Done {
			break
		}
	}
}
