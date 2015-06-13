package hunters

import (
	"fmt"
)

type DamageState int

const (
	DamageStateOK         DamageState = 0
	DamageStateDamaged                = 1
	DamageStateInoperable             = 2
)

var damageStateNames = map[DamageState]string{
	DamageStateOK:         "OK",
	DamageStateDamaged:    "damaged",
	DamageStateInoperable: "inoperable",
}

func (d DamageState) String() string {
	return damageStateNames[d]
}

type EquipmentName string

const (
	Hydrophones     EquipmentName = "Hydrophones"
	DivePlanes                    = "Dive Planes"
	Periscope                     = "Periscope"
	FlakGun                       = "Flak Gun"
	DeckGun                       = "Deck Gun"
	FwdTorpedoDoors               = "Forward Torpedo Doors"
	AftTorpedoDoors               = "Aft Torpedo Doors"
	FuelTanks                     = "Fuel Tanks"
	Radio                         = "Radio"
	Batteries                     = "Batteries"
)

type Damage struct {
	Sunk      bool
	Equipment map[EquipmentName]DamageState
}

// testEquipmentDamaged returns a gameTest for use with Game.getDrm().
func (g *Game) testEquipmentDamaged(n EquipmentName) gameTest {
	return func() (bool, string) {
		state := g.Boat.Damage.Equipment[n]
		damaged := (state != DamageStateOK)
		return damaged, fmt.Sprintf("%s %s", n, state)
	}
}

func (g *Game) RollEscortAirAttackDamage() {
	//TODO
}
