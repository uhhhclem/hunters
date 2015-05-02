package main

import "hunters"

func main() {

	g := &hunters.Game{
		Output: make(chan *hunters.Prompt),
		Input:  make(chan string),
	}
	g.State = &hunters.Start{Game: g}
	
	go g.HandleInput()

	for {
		g.State = g.State.Handle()
		if g.Done {
			break
		}
	}
}
