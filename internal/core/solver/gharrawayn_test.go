package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestGharrawaynWithHusband(t *testing.T) {
	shares, ok := Gharrawayn(heir.New().Set(heir.Husband, 1).Set(heir.Father, 1).Set(heir.Mother, 1))
	if !ok {
		t.Fatal("expected gharrawayn")
	}
	want := map[heir.Relation]rational.Fraction{
		heir.Husband: rational.New(1, 2),
		heir.Mother:  rational.New(1, 6),
		heir.Father:  rational.New(1, 3),
	}
	assertShares(t, shares, want)
	// The father takes twice the mother.
	if !shares[heir.Father].Equal(shares[heir.Mother].Mul(rational.FromInt(2))) {
		t.Error("father should be twice the mother")
	}
}

func TestGharrawaynWithWife(t *testing.T) {
	shares, ok := Gharrawayn(heir.New().Set(heir.Wife, 1).Set(heir.Father, 1).Set(heir.Mother, 1))
	if !ok {
		t.Fatal("expected gharrawayn")
	}
	want := map[heir.Relation]rational.Fraction{
		heir.Wife:   rational.New(1, 4),
		heir.Mother: rational.New(1, 4),
		heir.Father: rational.New(1, 2),
	}
	assertShares(t, shares, want)
	if !shares[heir.Father].Equal(shares[heir.Mother].Mul(rational.FromInt(2))) {
		t.Error("father should be twice the mother")
	}
}

func TestGharrawaynNotApplicable(t *testing.T) {
	no := map[string]*heir.Heirs{
		"with descendant": heir.New().Set(heir.Husband, 1).Set(heir.Father, 1).Set(heir.Mother, 1).Set(heir.Son, 1),
		"no father":       heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1),
		"with sibling":    heir.New().Set(heir.Wife, 1).Set(heir.Father, 1).Set(heir.Mother, 1).Set(heir.FullBrother, 1),
	}
	for name, h := range no {
		if _, ok := Gharrawayn(h); ok {
			t.Errorf("%s should not be gharrawayn", name)
		}
	}
}

// assertShares checks that got equals want exactly, including total of one.
func assertShares(t *testing.T, got, want map[heir.Relation]rational.Fraction) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d shares, want %d: %v", len(got), len(want), got)
	}
	total := rational.Zero()
	for r, w := range want {
		if !got[r].Equal(w) {
			t.Errorf("%s = %s, want %s", r, got[r], w)
		}
		total = total.Add(got[r])
	}
	if !total.IsWhole() {
		t.Errorf("shares sum to %s, want 1", total)
	}
}
