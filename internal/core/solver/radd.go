package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// RaddResult holds the shares after applying radd.
type RaddResult struct {
	// Applied reports whether the shares were under-subscribed, putting the
	// case in the radd situation.
	Applied bool
	// Shares are the adjusted shares of the whole estate.
	Shares map[heir.Relation]rational.Fraction
	// Surplus is the portion not returned to any heir. It is non-zero only
	// when the sole heir is a spouse, in which case the surplus passes to the
	// public treasury rather than to the spouse.
	Surplus rational.Fraction
}

// ApplyRadd returns the surplus proportionally to the non-spouse fixed-share
// heirs when the shares under-subscribe the estate. It assumes the caller has
// already established that there is no residuary heir; with a residuary heir
// the surplus goes there instead.
//
// The spouse is excluded from radd in the view of the majority of the schools:
// the spouse keeps the prescribed share and the rest is shared among the other
// fixed-share heirs in proportion to their shares. When the only heir is a
// spouse, the surplus passes to the public treasury.
func ApplyRadd(shares map[heir.Relation]rational.Fraction) RaddResult {
	total := rational.Zero()
	spouseShare := rational.Zero()
	for r, f := range shares {
		total = total.Add(f)
		if r.IsSpouse() {
			spouseShare = spouseShare.Add(f)
		}
	}

	if total.GreaterEqual(rational.One()) {
		return RaddResult{Applied: false, Shares: copyShares(shares), Surplus: rational.Zero()}
	}

	target := rational.One().Sub(spouseShare)
	eligibleSum := total.Sub(spouseShare)
	out := make(map[heir.Relation]rational.Fraction, len(shares))

	if eligibleSum.IsZero() {
		for r, f := range shares {
			out[r] = f
		}
		return RaddResult{Applied: true, Shares: out, Surplus: target}
	}

	factor := target.Div(eligibleSum)
	for r, f := range shares {
		if r.IsSpouse() {
			out[r] = f
		} else {
			out[r] = f.Mul(factor)
		}
	}
	return RaddResult{Applied: true, Shares: out, Surplus: rational.Zero()}
}

// copyShares returns a shallow copy of a share map.
func copyShares(m map[heir.Relation]rational.Fraction) map[heir.Relation]rational.Fraction {
	out := make(map[heir.Relation]rational.Fraction, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
