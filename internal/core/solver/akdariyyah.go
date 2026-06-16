package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// akdariyyahSister returns the lone sister relation when the heirs form the
// akdariyyah configuration: a husband, the mother, the paternal grandfather,
// and exactly one sister (full or consanguine) with no other heir.
func akdariyyahSister(h *heir.Heirs) (heir.Relation, bool) {
	if !h.Present(heir.Husband) || !h.Present(heir.Mother) || !h.Present(heir.PaternalGrandfather) {
		return heir.RelationInvalid, false
	}
	if len(h.Relations()) != 4 {
		return heir.RelationInvalid, false
	}
	switch {
	case h.Count(heir.FullSister) == 1 && !h.Present(heir.ConsanguineSister):
		return heir.FullSister, true
	case h.Count(heir.ConsanguineSister) == 1 && !h.Present(heir.FullSister):
		return heir.ConsanguineSister, true
	default:
		return heir.RelationInvalid, false
	}
}

// IsAkdariyyah reports whether the heirs form the akdariyyah configuration.
func IsAkdariyyah(h *heir.Heirs) bool {
	_, ok := akdariyyahSister(h)
	return ok
}

// Akdariyyah returns the canonical base-27 shares for the akdariyyah, or ok
// false when the heirs are not that configuration. It is a Zayd-view anomaly:
// the grandfather would otherwise leave the sister nothing, so her half is
// assigned, the problem is raised by awl, and the grandfather and sister then
// pool and re-divide their parts two to one, yielding the base of 27.
//
//	husband 1/3, mother 2/9, grandfather 8/27, sister 4/27
//
// Under the Abu Hanifa view the grandfather simply excludes the sister, so this
// case does not arise.
func Akdariyyah(h *heir.Heirs) (map[heir.Relation]rational.Fraction, bool) {
	sister, ok := akdariyyahSister(h)
	if !ok {
		return nil, false
	}
	return map[heir.Relation]rational.Fraction{
		heir.Husband:             rational.New(1, 3),
		heir.Mother:              rational.New(2, 9),
		heir.PaternalGrandfather: rational.New(8, 27),
		sister:                   rational.New(4, 27),
	}, true
}
