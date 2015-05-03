package hunters

import "tables"

type Combat struct {
    Day bool
    Surface bool
    Range string
    Escort bool
    Targets []Target
}

type Target struct {
    tables.Target
    Damage int
    Torpedoes []Torpedo
}