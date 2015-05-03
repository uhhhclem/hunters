package main

import (
	"hunters"
)

func main() {

	g := hunters.NewGame()
	
	go g.HandleInput()

	for {
		g.State = g.State.Handle()
		if g.Done {
			break
		}
	}
}
