package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
	"github.com/YASSERRMD/Faraid/internal/core/rules"
)

// oneThird and twoThirds split the residue that remains after the spouse in the
// gharrawayn cases.
var (
	oneThird  = rational.New(1, 3)
	twoThirds = rational.New(2, 3)
)

// Gharrawayn returns the shares for the two Umariyyatan (gharrawayn) cases, or
// ok false when the heirs are not a gharrawayn configuration.
//
// In these cases (a spouse with both parents, no descendant, and no siblings)
// the mother takes one third of the residue remaining after the spouse rather
// than one third of the whole estate, and the father takes the rest. This keeps
// the father at twice the mother:
//
//	husband 1/2, mother 1/6, father 1/3
//	wife    1/4, mother 1/4, father 1/2
func Gharrawayn(h *heir.Heirs) (map[heir.Relation]rational.Fraction, bool) {
	ctx := rules.Context{Heirs: h}
	if !rules.IsGharrawayn(ctx) {
		return nil, false
	}

	spouse, _ := presentSpouse(h)
	spouseShare, _, _ := rules.FixedShare(spouse, ctx)
	remainder := rational.One().Sub(spouseShare)

	return map[heir.Relation]rational.Fraction{
		spouse:      spouseShare,
		heir.Mother: remainder.Mul(oneThird),
		heir.Father: remainder.Mul(twoThirds),
	}, true
}
