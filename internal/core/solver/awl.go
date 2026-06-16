package solver

import (
	"math/big"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

// AwlResult records whether awl was applied and how the base changed, for the
// derivation.
type AwlResult struct {
	Applied      bool
	OriginalBase *big.Int
	AdjustedBase *big.Int
}

// ApplyAwl handles over-subscription. When the assigned shares sum to more than
// the base, the base is raised to that sum so the whole is shared in the same
// proportions, which lowers every heir. The numerators are unchanged; they are
// now expressed over the larger base. The classical bases that admit awl are 6
// (to 7, 8, 9, or 10), 12 (to 13, 15, or 17), and 24 (to 27).
func ApplyAwl(a Asl) (Asl, AwlResult) {
	if !a.NeedsAwl() {
		return a, AwlResult{
			OriginalBase: new(big.Int).Set(a.Base),
			AdjustedBase: new(big.Int).Set(a.Base),
		}
	}

	nums := make(map[heir.Relation]*big.Int, len(a.Numerators))
	for r, n := range a.Numerators {
		nums[r] = new(big.Int).Set(n)
	}
	adjusted := Asl{
		Base:          new(big.Int).Set(a.SumNumerators),
		Numerators:    nums,
		SumNumerators: new(big.Int).Set(a.SumNumerators),
	}
	return adjusted, AwlResult{
		Applied:      true,
		OriginalBase: new(big.Int).Set(a.Base),
		AdjustedBase: new(big.Int).Set(adjusted.Base),
	}
}
