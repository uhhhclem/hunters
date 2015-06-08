package hunters

import (
	"github.com/uhhhclem/mse/src/interact"
	t "tables"
)

type Patrol struct {
	Assignment             string
	TravelBoxes            []string
	TravelBox              int
	Natural12EncounterRoll bool
}

const (
	PatrolStartState interact.GameState = "PatrolStartState"
)

func init() {
	go func() {
		handlerRegistry <- handlerMap{
			PatrolStartState: handlePatrolStart,
		}
	}()
}

type patrolDefinition struct {
	transit   int
	onStation int
}

var patrols = map[string]map[string]patrolDefinition{
	t.TypeVIIB: {
		t.Atlantic:      patrolDefinition{2, 4},
		t.BritishIsles:  patrolDefinition{2, 4},
		t.BritishIslesM: patrolDefinition{2, 4},
		t.NorthAmerica:  patrolDefinition{4, 3},
		t.SpanishCoast:  patrolDefinition{2, 4},
		t.Norway:        patrolDefinition{2, 4},
		t.Arctic:        patrolDefinition{2, 4},
		t.TheMed:        patrolDefinition{2, 4},
	},
}

func handlePatrolStart(g *Game) interact.GameState {
	g.Log(string(PatrolStartState))
	p := &g.Patrol

	tn := t.Period1939ToMar1940
	roll, assignment := t.UBoatPatrolAssignmentChart.Roll(t.TableName(tn))
	// TODO: add rerolls for VIID/VIIC Flak and Mediterranean

	switch g.Boat.Type {
	case t.TypeVIIB, t.TypeVIIC:
		if assignment == t.WestAfricanCoast || assignment == t.Caribbean {
			assignment = t.Atlantic
		}
	case t.TypeVIID:
		if assignment == t.WestAfricanCoast {
			assignment = t.Atlantic
		}
		if assignment == t.BritishIsles {
			assignment = t.BritishIslesM
		}
		// TODO:  reassignments for IX
	}

	g.Logf("Patrol assigment:  rolled %d (%s)", roll, assignment)

	p.Assignment = string(assignment)
	pd := patrols[g.Boat.Type][p.Assignment]
	var tb []string

	for i := 0; i < pd.transit; i++ {
		tb = append(tb, string(t.Transit))
	}
	for i := 0; i < pd.onStation; i++ {
		tb = append(tb, p.Assignment)
	}
	for i := 0; i < pd.transit; i++ {
		tb = append(tb, string(t.Transit))
	}
	p.TravelBoxes = tb
	g.Logf("Travel boxes: %s", p.TravelBoxes)

	return EncounterStartState
}
