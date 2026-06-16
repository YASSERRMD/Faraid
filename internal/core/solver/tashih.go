package solver

import (
	"math/big"

	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// TashihGroup is a share that must split evenly: a numerator over the base that
// is divided among Heads equal units. For a residuary group of mixed sex the
// heads are weighted (a male counts as two units), so the group still divides
// into equal units.
type TashihGroup struct {
	Numerator *big.Int
	Heads     int64
}

// TashihResult holds the corrected base and per-group numerators.
type TashihResult struct {
	// Applied reports whether a correction was needed.
	Applied bool
	// Factor is the multiplier applied to the base and to every numerator.
	Factor *big.Int
	// Base is the corrected base of the problem.
	Base *big.Int
	// Numerators are the corrected group numerators, in the order given.
	Numerators []*big.Int
}

// Tashih corrects the problem so every group's numerator divides evenly among
// its heads. For each group that does not divide it computes heads divided by
// the gcd of heads and numerator, then multiplies the base and all numerators
// by the least common multiple of those values. Ratios are preserved and every
// individual then receives an integer number of parts.
func Tashih(base *big.Int, groups []TashihGroup) TashihResult {
	factor := big.NewInt(1)
	for _, g := range groups {
		if g.Heads <= 1 || g.Numerator.Sign() == 0 {
			continue
		}
		heads := big.NewInt(g.Heads)
		gcd := rational.GCD(heads, g.Numerator)
		reduced := new(big.Int).Div(heads, gcd)
		factor = rational.LCM(factor, reduced)
	}

	newBase := new(big.Int).Mul(base, factor)
	newNums := make([]*big.Int, len(groups))
	for i, g := range groups {
		newNums[i] = new(big.Int).Mul(g.Numerator, factor)
	}

	return TashihResult{
		Applied:    factor.Cmp(big.NewInt(1)) > 0,
		Factor:     factor,
		Base:       newBase,
		Numerators: newNums,
	}
}
