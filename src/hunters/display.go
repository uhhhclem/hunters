package hunters

import (
	"tables"
)

func (g *Game) LoadTestData() {
	// initialize the Boat
	g.Forward = TorpedoSection{
		Capacity:        4,
		SteamReloads:    8,
		ElectricReloads: 4,
	}
	g.Aft = TorpedoSection{
		Capacity:     1,
		SteamReloads: 2,
	}
	initTubes(&g.Forward, 4)
	initTubes(&g.Aft, 1)

	// initialize the Combat
	g.Combat = Combat{
		Day:     true,
		Surface: false,
		Range:   CloseRange,
		Escort:  true,
		Targets: make([]Target, 3),
	}
	var t = g.Combat.Targets
	t[0] = makeTarget(tables.SmallFreighterTargetRoster, 1)
	t[1] = makeTarget(tables.LargeFreighterTargetRoster, 2)
	t[2] = makeTarget(tables.TankerTargetRoster, 3)
	t[2].Damage = 1
	for i := range g.Combat.Targets {
		t := &g.Combat.Targets[i]
		t.ToSink = toSink(*t)
	}
}

func makeTarget(r tables.TargetRoster, n int) Target {
	_, t := r.Roll()
	return Target{
		Target:    t,
		Number:    n,
		Torpedoes: make([]Torpedo, 10),
	}
}
func initTubes(s *TorpedoSection, cap int) {
	s.Capacity = cap
	s.Tubes = make([]Tube, cap, cap)
	for i := range s.Tubes {
		s.Tubes[i] = Tube{Number: i + 1, Torpedo: EmptyTube}
	}
}
