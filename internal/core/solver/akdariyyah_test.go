package solver

import (
	"math/big"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestAkdariyyahFullSister(t *testing.T) {
	h := heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).
		Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 1)
	shares, ok := Akdariyyah(h)
	if !ok {
		t.Fatal("expected akdariyyah")
	}
	want := map[heir.Relation]rational.Fraction{
		heir.Husband:             rational.New(1, 3),
		heir.Mother:              rational.New(2, 9),
		heir.PaternalGrandfather: rational.New(8, 27),
		heir.FullSister:          rational.New(4, 27),
	}
	total := rational.Zero()
	for r, w := range want {
		if !shares[r].Equal(w) {
			t.Errorf("%s = %s, want %s", r, shares[r], w)
		}
		total = total.Add(shares[r])
	}
	if !total.IsWhole() {
		t.Errorf("shares sum to %s, want 1", total)
	}
	// The canonical base is 27.
	if a := ComputeAsl(shares); a.Base.Cmp(big.NewInt(27)) != 0 {
		t.Errorf("base = %s, want 27", a.Base)
	}
}

func TestAkdariyyahConsanguineSister(t *testing.T) {
	h := heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).
		Set(heir.PaternalGrandfather, 1).Set(heir.ConsanguineSister, 1)
	shares, ok := Akdariyyah(h)
	if !ok {
		t.Fatal("expected akdariyyah with consanguine sister")
	}
	if !shares[heir.ConsanguineSister].Equal(rational.New(4, 27)) {
		t.Errorf("consanguine sister = %s, want 4/27", shares[heir.ConsanguineSister])
	}
}

func TestNotAkdariyyah(t *testing.T) {
	no := map[string]*heir.Heirs{
		"two sisters":    heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 2),
		"no grandfather": heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.FullSister, 1),
		"extra heir":     heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 1).Set(heir.Son, 1),
		"both sisters":   heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 1).Set(heir.ConsanguineSister, 1),
		"no husband":     heir.New().Set(heir.Wife, 1).Set(heir.Mother, 1).Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 1),
	}
	for name, h := range no {
		if IsAkdariyyah(h) {
			t.Errorf("%s should not be akdariyyah", name)
		}
		if _, ok := Akdariyyah(h); ok {
			t.Errorf("%s should not return akdariyyah shares", name)
		}
	}
}
