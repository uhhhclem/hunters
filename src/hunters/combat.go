package hunters

import "tables"

type Combat struct {
    Range string
    Targets []Target
}

type Target struct {
    tables.Target
    Damage int
    Torpedoes []Torpedo
}