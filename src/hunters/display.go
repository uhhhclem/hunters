package hunters

import (
    "os"
    "log"
    "text/template"
)



func (g *Game) LoadTestData() {
    g.Forward = TorpedoSection{
        Capacity: 4,
        SteamReloads: 8,
        ElectricReloads: 4,
    }
    g.Aft = TorpedoSection{
        Capacity: 1,
        SteamReloads: 2,
    }
    initTubes(&g.Forward, 4)
    initTubes(&g.Aft, 1)
}

func initTubes(s *TorpedoSection, cap int) {
    s.Capacity = cap
    s.Tubes = make([]Tube, cap, cap)
    for i := range s.Tubes {
        s.Tubes[i] = Tube{Number: i+1, Torpedo: EmptyTube}
    }
}

var boatTxt = `{{.ID}} (Type {{.Type}})
Kmdt: {{.Kommandant}}

Forward Tubes: {{range .Forward.Tubes}}{{.Number}}: {{.Torpedo}} {{end}}
      Reloads: {{.Forward.SteamReloads}} Steam, {{.Forward.ElectricReloads}} Electric
         Ammo:
    Aft Tubes: {{range .Aft.Tubes}}{{.Number}}: {{.Torpedo}} {{end}}
      Reloads: {{.Aft.SteamReloads}} Steam, {{.Aft.ElectricReloads}} Electric

`

var combatTxt = ``

var boatT, combatT *template.Template

func init() {
    boatT = template.Must(template.New("boat").Parse(boatTxt))
    combatT = template.Must(template.New("combat").Parse(combatTxt))
}

func (g *Game) Dump() {
    if err := boatT.Execute(os.Stdout, g.Boat); err != nil {
        log.Fatal(err)
    }
}