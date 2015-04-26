package tables

import "math/rand"

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
	Night                   = "Night"
	Period1939ToMar1940     = "1939 - Mar 1940"
	Ship                    = "Ship"
	ShipSize                = "ShipSize"
	ShipPlusEscort          = "Ship + Escort"
	SmallFreighter          = "Small Freighter"
	SpanishCoast            = "Spanish Coast"
	SpecialMissions         = "Special Missions"
	Tanker                  = "Tanker"
	Transit                 = "Transit"
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

func Roll1D6() int {
	return rand.Intn(6) + 1
}

func Roll2D6() int {
	return rand.Intn(6) + rand.Intn(6) + 2
}

func init() {
	EncounterChart.Tables[Gibraltar] = EncounterChart.Tables[AdditionalRoundOfCombat]
	EncounterChart.Tables[SpecialMissions] = EncounterChart.Tables[BayOfBiscay]
}
