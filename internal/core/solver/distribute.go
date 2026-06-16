package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// distributeByUnits splits a portion of the estate among the given relations,
// weighting each individual by sex (a male counts as two units, a female as
// one), and returns each relation's share of the whole. Relations with no
// individuals are skipped.
func distributeByUnits(portion rational.Fraction, members []heir.Relation, h *heir.Heirs) map[heir.Relation]rational.Fraction {
	out := map[heir.Relation]rational.Fraction{}
	var totalUnits int64
	for _, r := range members {
		totalUnits += int64(h.Count(r)) * weightOf(r.Sex())
	}
	if totalUnits == 0 || portion.IsZero() {
		return out
	}
	total := rational.FromInt(totalUnits)
	for _, r := range members {
		u := int64(h.Count(r)) * weightOf(r.Sex())
		if u == 0 {
			continue
		}
		out[r] = portion.Mul(rational.FromInt(u)).Div(total)
	}
	return out
}
