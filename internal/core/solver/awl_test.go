package solver

import (
	"math/big"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// aslWith builds an Asl with the given base and one numerator per unit so the
// sum equals total, for exercising awl across the classical bases.
func aslWith(base, total int64) Asl {
	return Asl{
		Base:          big.NewInt(base),
		Numerators:    map[heir.Relation]*big.Int{heir.Husband: big.NewInt(total)},
		SumNumerators: big.NewInt(total),
	}
}

func TestApplyAwlClassicalBases(t *testing.T) {
	cases := []struct{ base, sum int64 }{
		{6, 7}, {6, 8}, {6, 9}, {6, 10},
		{12, 13}, {12, 15}, {12, 17},
		{24, 27},
	}
	for _, c := range cases {
		adjusted, res := ApplyAwl(aslWith(c.base, c.sum))
		if !res.Applied {
			t.Errorf("base %d sum %d: awl should apply", c.base, c.sum)
		}
		if adjusted.Base.Cmp(big.NewInt(c.sum)) != 0 {
			t.Errorf("base %d sum %d: adjusted base = %s, want %d", c.base, c.sum, adjusted.Base, c.sum)
		}
		if !adjusted.IsExact() {
			t.Errorf("base %d sum %d: adjusted should be exact", c.base, c.sum)
		}
		if res.OriginalBase.Cmp(big.NewInt(c.base)) != 0 {
			t.Errorf("original base = %s, want %d", res.OriginalBase, c.base)
		}
	}
}

func TestApplyAwlNotNeeded(t *testing.T) {
	a := aslWith(6, 5) // under-subscribed, no awl
	adjusted, res := ApplyAwl(a)
	if res.Applied {
		t.Error("awl should not apply when under-subscribed")
	}
	if adjusted.Base.Cmp(big.NewInt(6)) != 0 {
		t.Errorf("base should be unchanged, got %s", adjusted.Base)
	}
}

func TestApplyAwlFromShares(t *testing.T) {
	// Classic 6 to 7: husband 1/2 and two full sisters 2/3.
	a := ComputeAsl(map[heir.Relation]rational.Fraction{
		heir.Husband:    rational.New(1, 2),
		heir.FullSister: rational.New(2, 3),
	})
	adjusted, res := ApplyAwl(a)
	if !res.Applied || adjusted.Base.Cmp(big.NewInt(7)) != 0 {
		t.Errorf("6 to 7 awl failed: base %s applied %v", adjusted.Base, res.Applied)
	}
	if adjusted.Numerators[heir.Husband].Cmp(big.NewInt(3)) != 0 ||
		adjusted.Numerators[heir.FullSister].Cmp(big.NewInt(4)) != 0 {
		t.Errorf("numerators changed unexpectedly: %v", adjusted.Numerators)
	}

	// The minbariyya, base 24 to 27: wife 1/8, two daughters 2/3, father 1/6, mother 1/6.
	m := ComputeAsl(map[heir.Relation]rational.Fraction{
		heir.Wife:     rational.New(1, 8),
		heir.Daughter: rational.New(2, 3),
		heir.Father:   rational.New(1, 6),
		heir.Mother:   rational.New(1, 6),
	})
	adjustedM, resM := ApplyAwl(m)
	if !resM.Applied || adjustedM.Base.Cmp(big.NewInt(27)) != 0 {
		t.Errorf("minbariyya 24 to 27 failed: base %s applied %v", adjustedM.Base, resM.Applied)
	}
}
