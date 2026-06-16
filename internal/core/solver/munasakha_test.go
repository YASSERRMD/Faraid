package solver

import (
	"math/big"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func st(base int64, shares map[string]int64) ShareTable {
	m := make(map[string]*big.Int, len(shares))
	for k, v := range shares {
		m[k] = big.NewInt(v)
	}
	return ShareTable{Base: big.NewInt(base), Shares: m}
}

func wantShares(t *testing.T, got ShareTable, base int64, shares map[string]int64) {
	t.Helper()
	if got.Base.Cmp(big.NewInt(base)) != 0 {
		t.Errorf("base = %s, want %d", got.Base, base)
	}
	for label, n := range shares {
		if got.Shares[label] == nil || got.Shares[label].Cmp(big.NewInt(n)) != 0 {
			t.Errorf("%s = %v, want %d", label, got.Shares[label], n)
		}
	}
	if len(got.Shares) != len(shares) {
		t.Errorf("got %d heirs, want %d: %v", len(got.Shares), len(shares), got.Shares)
	}
}

func TestMunasakhaSecondBaseDivides(t *testing.T) {
	// Deceased C has 2 of 4; the second base 2 divides it cleanly.
	first := st(4, map[string]int64{"A": 1, "B": 1, "C": 2})
	second := st(2, map[string]int64{"D": 1, "E": 1})
	got := Munasakha(first, "C", second)
	wantShares(t, got, 4, map[string]int64{"A": 1, "B": 1, "D": 1, "E": 1})
}

func TestMunasakhaNeedsScaling(t *testing.T) {
	// Deceased C has 2 of 4; the second base 3 is coprime, so the base scales.
	first := st(4, map[string]int64{"A": 2, "C": 2})
	second := st(3, map[string]int64{"D": 1, "E": 1, "F": 1})
	got := Munasakha(first, "C", second)
	wantShares(t, got, 12, map[string]int64{"A": 6, "D": 2, "E": 2, "F": 2})
}

func TestMunasakhaOverlappingHeir(t *testing.T) {
	// A inherits from both estates; the portions add.
	first := st(4, map[string]int64{"A": 1, "C": 3})
	second := st(2, map[string]int64{"A": 1, "D": 1})
	got := Munasakha(first, "C", second)
	// firstMult 2, secondMult 3: A = 1*2 + 1*3 = 5, D = 1*3 = 3, base 8.
	wantShares(t, got, 8, map[string]int64{"A": 5, "D": 3})
	if !got.Fraction("A").Equal(rational.New(5, 8)) {
		t.Errorf("A fraction = %s, want 5/8", got.Fraction("A"))
	}
}

func TestMunasakhaThreeDeaths(t *testing.T) {
	// Chain three estates: B dies after the first, then D dies after the second.
	first := st(2, map[string]int64{"A": 1, "B": 1})
	second := st(2, map[string]int64{"C": 1, "D": 1})
	step1 := Munasakha(first, "B", second)
	wantShares(t, step1, 4, map[string]int64{"A": 2, "C": 1, "D": 1})

	third := st(2, map[string]int64{"E": 1, "F": 1})
	step2 := Munasakha(step1, "D", third)
	wantShares(t, step2, 8, map[string]int64{"A": 4, "C": 2, "E": 1, "F": 1})
}

func TestMunasakhaMissingDeceased(t *testing.T) {
	// A deceased label not present contributes nothing to the second heirs.
	first := st(2, map[string]int64{"A": 1, "B": 1})
	second := st(2, map[string]int64{"C": 1, "D": 1})
	got := Munasakha(first, "Z", second)
	// d = 0, so the second heirs get nothing and the first is unchanged.
	if got.Fraction("C").Sign() != 0 {
		t.Errorf("missing deceased should give second heirs nothing, got %s", got.Fraction("C"))
	}
	if !got.Fraction("A").Equal(rational.New(1, 2)) {
		t.Errorf("A = %s, want 1/2", got.Fraction("A"))
	}
	// An entirely absent label has a zero share.
	if !got.Fraction("nobody").IsZero() {
		t.Errorf("absent label should be zero, got %s", got.Fraction("nobody"))
	}
}
