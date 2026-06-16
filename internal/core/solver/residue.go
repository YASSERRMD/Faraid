package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// weightOf returns the residue weight of one individual: a male counts as two,
// a female as one, which encodes the two to one preference. Within a single sex
// the weights cancel to an equal split.
func weightOf(s heir.Sex) int64 {
	if s == heir.Male {
		return 2
	}
	return 1
}

// DistributeResidue allocates the residue among the residuary members and
// returns each member slot's share of the whole estate. A male individual
// takes twice a female individual; members of the same sex share equally. The
// counts come from the heir set. An empty result is returned when there is no
// residue, no residuary heir, or the members have no individuals.
func DistributeResidue(residue rational.Fraction, r Residuary, h *heir.Heirs) map[heir.Relation]rational.Fraction {
	out := map[heir.Relation]rational.Fraction{}
	if residue.IsZero() || !r.Found() {
		return out
	}

	var totalUnits int64
	units := make(map[heir.Relation]int64, len(r.Members))
	for _, m := range r.Members {
		u := int64(h.Count(m.Relation)) * weightOf(m.Relation.Sex())
		units[m.Relation] = u
		totalUnits += u
	}
	if totalUnits == 0 {
		return out
	}

	total := rational.FromInt(totalUnits)
	for _, m := range r.Members {
		out[m.Relation] = residue.Mul(rational.FromInt(units[m.Relation])).Div(total)
	}
	return out
}
