package solver

import (
	"math/big"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Asl is the asl al-mas'ala: the common base of a problem and each fixed
// share expressed as an integer numerator over that base.
type Asl struct {
	// Base is the least common denominator of all assigned fixed shares.
	Base *big.Int
	// Numerators maps each heir to its share numerator over Base.
	Numerators map[heir.Relation]*big.Int
	// SumNumerators is the total of all numerators, compared with Base to
	// detect over-subscription (awl) or a residue.
	SumNumerators *big.Int
}

// ComputeAsl derives the base from the assigned fixed shares and expresses each
// share as a numerator over it. The result is deterministic regardless of map
// iteration order, since the base is an LCM and the numerators are independent.
func ComputeAsl(shares map[heir.Relation]rational.Fraction) Asl {
	dens := make([]*big.Int, 0, len(shares))
	for _, f := range shares {
		dens = append(dens, f.Den())
	}
	base := rational.LCMSlice(dens)
	baseF := rational.FromBigRat(new(big.Rat).SetInt(base))

	numerators := make(map[heir.Relation]*big.Int, len(shares))
	sum := new(big.Int)
	for r, f := range shares {
		num := f.Mul(baseF).Num() // exact integer, since base is a multiple of the denominator
		numerators[r] = num
		sum.Add(sum, num)
	}

	return Asl{Base: base, Numerators: numerators, SumNumerators: sum}
}

// NeedsAwl reports whether the assigned shares over-subscribe the estate, so
// the base must be raised proportionally (awl).
func (a Asl) NeedsAwl() bool { return a.SumNumerators.Cmp(a.Base) > 0 }

// HasResidue reports whether the assigned shares leave a remainder, which goes
// to the residuary heirs or, failing that, returns by radd.
func (a Asl) HasResidue() bool { return a.SumNumerators.Cmp(a.Base) < 0 }

// IsExact reports whether the assigned shares consume exactly the whole estate.
func (a Asl) IsExact() bool { return a.SumNumerators.Cmp(a.Base) == 0 }
