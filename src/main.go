package main

import "hunters"

func main() {
    
    g := &hunters.Game{}
    g.State = &hunters.Start{}
    for {
        g.State = g.State.Handle()
        if g.State == nil {
            break
        }
    }
}