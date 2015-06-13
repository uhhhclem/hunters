package tables

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

const (
	AdditionalRoundOfCombat = "Additional Round Of Combat"
	Aircraft                = "Aircraft"
	Arctic                  = "Arctic"
	Atlantic                = "Atlantic"
	BayOfBiscay             = "Bay Of Biscay"
	BritishIsles            = "British Isles"
	BritishIslesM           = "British Isles (M)"
	CapitalShip             = "Capital Ship"
	Caribbean               = "Caribbean"
	Convoy                  = "Convoy"
	Day                     = "Day"
	DayOrNight              = "Day Or Night"
	Escort                  = "Escort"
	Gibraltar               = "Gibraltar"
	LargeFreighter          = "Large Freighter"
	NorthAmerica            = "North America"
	Norway                  = "Norway"
	Night                   = "Night"
	Period1939ToMar1940     = "1939 - Mar 1940"
	RandomEvent             = "Random Event"
	Ship                    = "Ship"
	ShipSize                = "ShipSize"
	ShipPlusEscort          = "Ship + Escort"
	SmallFreighter          = "Small Freighter"
	SpanishCoast            = "Spanish Coast"
	SpecialMissions         = "Special Missions"
	Tanker                  = "Tanker"
	TheMed                  = "The Med"
	Transit                 = "Transit"
	TypeVIIB                = "VIIB"
	TypeVIIC                = "VIIC"
	TypeVIICFlak            = "VIIC Flak"
	TypeVIIA                = "VIIA"
	TypeVIID                = "VIID"
	TypeIX                  = "IX"
	TwoShips                = "Two Ships"
	TwoShipsPlusEscort      = "Two Ships + Escort"
	SwitchToNight           = "Switch To Night"
	WestAfricanCoast        = "West African Coast"
)

type Roller func() int

type TableName string

type Result string

type Table map[int]Result

type Chart struct {
	RollDice Roller
	Tables   map[TableName]Table
}

type Target struct {
	Type    string
	Tonnage int
	ShipID  string
}

func (t Target) String() string {
	return fmt.Sprintf("%s: %s (%dT", t.Type, t.ShipID, t.Tonnage)
}

type TargetRoster []Target

var EncounterChart Chart = Chart{
	RollDice: Roll2D6,
	Tables: map[TableName]Table{
		Transit: {
			2:  Aircraft,
			3:  Aircraft,
			12: Ship,
		},
		Arctic: {
			2:  CapitalShip,
			3:  Ship,
			6:  Convoy,
			7:  Convoy,
			8:  Convoy,
			12: Aircraft,
		},
		Atlantic: {
			2:  CapitalShip,
			3:  Ship,
			6:  Convoy,
			7:  Convoy,
			9:  Convoy,
			12: Convoy,
		},
		BritishIsles: {
			2:  CapitalShip,
			4:  Ship,
			5:  ShipPlusEscort,
			8:  Ship,
			10: Convoy,
			12: Aircraft,
		},
		Caribbean: {
			2:  Aircraft,
			4:  Ship,
			6:  TwoShipsPlusEscort,
			8:  Ship,
			9:  Tanker,
			10: Tanker,
			12: Aircraft,
		},
		WestAfricanCoast: {
			2:  Aircraft,
			4:  Ship,
			6:  TwoShipsPlusEscort,
			8:  Ship,
			9:  Tanker,
			10: Tanker,
			12: Aircraft,
		},
		AdditionalRoundOfCombat: {
			0: Escort,
			1: Escort,
			2: Escort,
			3: Escort,
			4: Aircraft,
			5: Aircraft,
		},
		BayOfBiscay: {
			0: Aircraft,
			1: Aircraft,
			2: Aircraft,
			3: Aircraft,
			4: Aircraft,
		},
	},
}

var EncounterChartSupplement = Chart{
	RollDice: Roll1D6,
	Tables: map[TableName]Table{
		DayOrNight: {
			1: Day,
			2: Day,
			3: Day,
			4: Night,
			5: Night,
			6: Night,
		},
		SwitchToNight: {
			1: Night,
			2: Night,
			3: Night,
			4: Night,
			5: Day,
			6: Day,
		},
		ShipSize: {
			1: SmallFreighter,
			2: SmallFreighter,
			3: SmallFreighter,
			4: LargeFreighter,
			5: LargeFreighter,
			6: Tanker,
		},
	},
}

var UBoatPatrolAssignmentChart = Chart{
	RollDice: Roll2D6,
	Tables: map[TableName]Table{
		Period1939ToMar1940: {
			2:  SpanishCoast,
			3:  BritishIsles,
			4:  BritishIsles,
			5:  BritishIslesM,
			6:  BritishIsles,
			7:  BritishIsles,
			8:  BritishIsles,
			9:  BritishIslesM,
			10: BritishIsles,
			11: BritishIsles,
			12: WestAfricanCoast,
		},
	},
}

func (c Chart) Roll(tn TableName) (int, Result) {
	t := c.Tables[tn]
	r := c.RollDice()
	return r, t[r]
}

func populateTargetRoster(data string, targetType string) TargetRoster {
	lines := strings.Split(data, "\n")
	r := make(TargetRoster, len(lines))
	for i, line := range lines {
		s := strings.Split(line, ",")
		tons, err := strconv.Atoi(s[0])
		if err != nil {
			log.Fatal(err)
		}
		r[i] = Target{
			Type:    targetType,
			Tonnage: tons,
			ShipID:  s[1],
		}
	}
	return r
}

func (t TargetRoster) Roll() (int, Target) {
	r := rand.Intn(len(t))
	return r, t[r]
}

func Roll1D6() int {
	return rand.Intn(6) + 1
}

func Roll2D6() int {
	return rand.Intn(6) + rand.Intn(6) + 2
}

type DRM struct {
	Mod  int
	Name string
}

func (d DRM) String() string {
	sign := " "
	if d.Mod > 0 {
		sign = "+"
	}
	if d.Mod < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%d %s", sign, d.Mod, d.Name)
}

const (
	NoMinRoll int = -9999
	NoMaxRoll     = 9999
)

func Roll2D6WithDRMs(drms []DRM, min, max int) (natural, modified int) {
	natural = Roll2D6()
	modified = natural
	for _, m := range drms {
		modified += m.Mod
	}
	if min != NoMinRoll && modified < min {
		modified = min
	}
	if max != NoMaxRoll && modified > max {
		modified = max
	}
	return
}

func init() {
	EncounterChart.Tables[Gibraltar] = EncounterChart.Tables[AdditionalRoundOfCombat]
	EncounterChart.Tables[SpecialMissions] = EncounterChart.Tables[BayOfBiscay]

	SmallFreighterTargetRoster = populateTargetRoster(smallFreighterTargets, SmallFreighter)
	LargeFreighterTargetRoster = populateTargetRoster(largeFreighterTargets, LargeFreighter)
	TankerTargetRoster = populateTargetRoster(tankerTargets, Tanker)
	CapitalShipTargetRoster = populateTargetRoster(capitalShipTargets, CapitalShip)
}

var (
	SmallFreighterTargetRoster TargetRoster
	LargeFreighterTargetRoster TargetRoster
	TankerTargetRoster         TargetRoster
	CapitalShipTargetRoster    TargetRoster
	NorthAmericaTargetRoster   map[string]TargetRoster
)

var smallFreighterTargets = `1800,Bosnia
4100,Rio Claro
1800,Gartavon
4800,RoyalSceptre,
4500,Blairlogie
4900,Firby
4000,Avemore
5000,Kafiristan
1000,Truro
2700,Akenside`

var largeFreighterTargets = `12300,Sultan Star
5300,SS Browning
7200,Manaar
5200,Fanad Head
5500,Kennebec
7000,Louisiana
9200,Lochavon
10000,Bretagne
5400,Vermont
7200,City Of Mandalay`

var tankerTargets = `9400,Inverliffey
10000,Regent Tiger
8500,British Influence
8800,Cheyenne
14000,Emile-Miguet
11000,Arne Kjode
7400,San Alberto
5000,Vaclite,
5200,Chastine Maersk
8000,Ceronia`

var capitalShipTargets = `22000,CV Ark Royal*
29100,BB Royal Oak*
18600,CV Courageous
10000,CA Belfast
31100,BB Barnham*
34000,BB Nelson*
31300,BB Malaya*
22600,CV Eagle*
12800,CVE Avenger
11000,CVE Audacity`
