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
        d, r := tables.EncounterChartSupplement.Roll(tables.ShipSize)
        fmt.Printf("%d: %s ", d, r)
    }
    fmt.Println()
}