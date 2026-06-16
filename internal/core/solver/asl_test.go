package solver

import (
	"math/big"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func bi(n int64) *big.Int { return big.NewInt(n) }

func TestComputeAslExact(t *testing.T) {
	// Husband 1/2 and a single daughter 1/2: base 2, numerators 1 and 1.
	a := ComputeAsl(map[heir.Relation]rational.Fraction{
		heir.Husband:  rational.New(1, 2),
		heir.Daughter: rational.New(1, 2),
	})
	if a.Base.Cmp(bi(2)) != 0 {
		t.Errorf("base = %s, want 2", a.Base)
	}
	if a.Numerators[heir.Husband].Cmp(bi(1)) != 0 || a.Numerators[heir.Daughter].Cmp(bi(1)) != 0 {
		t.Errorf("numerators = %v", a.Numerators)
	}
	if !a.IsExact() || a.NeedsAwl() || a.HasResidue() {
		t.Error("shares should be exact")
	}
}

func TestComputeAslResidue(t *testing.T) {
	// Husband 1/2 and mother 1/3: base 6, numerators 3 and 2, sum 5 < 6.
	a := ComputeAsl(map[heir.Relation]rational.Fraction{
		heir.Husband: rational.New(1, 2),
		heir.Mother:  rational.New(1, 3),
	})
	if a.Base.Cmp(bi(6)) != 0 {
		t.Errorf("base = %s, want 6", a.Base)
	}
	if a.SumNumerators.Cmp(bi(5)) != 0 {
		t.Errorf("sum = %s, want 5", a.SumNumerators)
	}
	if !a.HasResidue() || a.IsExact() || a.NeedsAwl() {
		t.Error("shares should leave a residue")
	}
}

func TestComputeAslAwl(t *testing.T) {
	// Husband 1/2 and two full sisters 2/3: 3/6 + 4/6 = 7/6, over-subscribed.
	a := ComputeAsl(map[heir.Relation]rational.Fraction{
		heir.Husband:    rational.New(1, 2),
		heir.FullSister: rational.New(2, 3),
	})
	if a.Base.Cmp(bi(6)) != 0 {
		t.Errorf("base = %s, want 6", a.Base)
	}
	if a.SumNumerators.Cmp(bi(7)) != 0 {
		t.Errorf("sum = %s, want 7", a.SumNumerators)
	}
	if !a.NeedsAwl() || a.IsExact() || a.HasResidue() {
		t.Error("shares should need awl")
	}
}

func TestComputeAslMixedBase(t *testing.T) {
	// Wife 1/8, mother 1/6, two daughters 2/3: base 24.
	a := ComputeAsl(map[heir.Relation]rational.Fraction{
		heir.Wife:     rational.New(1, 8),
		heir.Mother:   rational.New(1, 6),
		heir.Daughter: rational.New(2, 3),
	})
	if a.Base.Cmp(bi(24)) != 0 {
		t.Errorf("base = %s, want 24", a.Base)
	}
	// 3/24 + 4/24 + 16/24 = 23/24, a residue of 1/24 remains.
	if a.SumNumerators.Cmp(bi(23)) != 0 {
		t.Errorf("sum = %s, want 23", a.SumNumerators)
	}
}

func TestComputeAslEmpty(t *testing.T) {
	a := ComputeAsl(nil)
	if a.Base.Cmp(bi(1)) != 0 {
		t.Errorf("empty base = %s, want 1", a.Base)
	}
	if a.SumNumerators.Sign() != 0 {
		t.Errorf("empty sum = %s, want 0", a.SumNumerators)
	}
}
