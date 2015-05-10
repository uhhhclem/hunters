package hunters

import (
	"log"
	"os"
	"text/template"

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

var boatTxt = `=================================================================

{{.ID}} (Type {{.Type}})
Kmdt: {{.Kommandant}}

Forward Tubes: {{range .Forward.Tubes}}{{.Number}}: {{.Torpedo}} {{end}}
      Reloads: {{.Forward.SteamReloads}} Steam, {{.Forward.ElectricReloads}} Electric
         Ammo:
    Aft Tubes: {{range .Aft.Tubes}}{{.Number}}: {{.Torpedo}} {{end}}
      Reloads: {{.Aft.SteamReloads}} Steam, {{.Aft.ElectricReloads}} Electric

`

var combatTxt = `---------------------------------------------------------------

{{if .Day}}Day{{else}}Night{{end}}, {{if .Surface}}Surface{{else}}Submerged{{end}}, {{.Range}} Range
{{if .Escort}}Escorted{{else}}Unescorted{{end}}

Targets: 
{{range .Targets}}  {{.Number}}: {{.ShipID}} ({{.Type}}, {{.Tonnage}}T)
{{end}}
`

var boatT, combatT *template.Template

func init() {
	boatT = template.Must(template.New("boat").Parse(boatTxt))
	combatT = template.Must(template.New("combat").Parse(combatTxt))
}

func (g *Game) Dump() {
	if err := boatT.Execute(os.Stdout, g.Boat); err != nil {
		log.Fatal(err)
	}
	if err := combatT.Execute(os.Stdout, g.Combat); err != nil {
		log.Fatal(err)
	}
}
