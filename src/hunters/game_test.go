package hunters

import (
	"testing"
)

func TestGetDrm(t *testing.T) {
	g := NewGame()

	d := g.getDrm(1, g.testEquipmentDamaged(Hydrophones))
	if got, want := d.String(), " 0 Hydrophones OK"; got != want {
		t.Errorf("got: %q, want %q", got, want)
	}

	g.Damage.Equipment[Hydrophones] = DamageStateDamaged

	d = g.getDrm(1, g.testEquipmentDamaged(Hydrophones))
	if got, want := d.String(), "+1 Hydrophones damaged"; got != want {
		t.Errorf("got: %q, want %q", got, want)
	}
}
