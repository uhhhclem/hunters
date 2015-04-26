package main

import (
    "fmt"
    "math/rand"
    "tables"
    "time"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < 10; i++ {
        d, r := tables.TankerTargetRoster.Roll()
        fmt.Printf("%d: %s (%d) ", d, r.ShipID, r.Tonnage)
    }
    fmt.Println()
}