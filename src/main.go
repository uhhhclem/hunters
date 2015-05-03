package main

import (
	"hunters"
)

func main() {

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
