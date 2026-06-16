package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// JaddView selects how a school treats the grandfather competing with the
// brothers and sisters (jadd wa ikhwa).
type JaddView int

const (
	// JaddZayd is the view of Zayd ibn Thabit, adopted by the Maliki, Shafi'i,
	// and Hanbali schools: the grandfather competes with the full and
	// consanguine siblings and takes the better of sharing with them
	// (muqasama), one third of the remainder, or one sixth of the estate.
	JaddZayd JaddView = iota
	// JaddAbuHanifa is the view of Abu Hanifa: the grandfather stands in for
	// the father and excludes the siblings entirely.
	JaddAbuHanifa
)

// String returns a label for the view.
func (v JaddView) String() string {
	if v == JaddAbuHanifa {
		return "abu-hanifa"
	}
	return "zayd"
}

// JaddResult holds the outcome of the grandfather-with-siblings computation,
// with all shares expressed as fractions of the whole estate.
type JaddResult struct {
	GrandfatherShare rational.Fraction
	SiblingShares    map[heir.Relation]rational.Fraction
	SiblingsExcluded bool
	// Method names which option the grandfather took under the Zayd view.
	Method string
	// NeedsReview marks a configuration whose ruling this engine does not yet
	// resolve, so the result must not be treated as authoritative.
	NeedsReview bool
	ReviewNote  string
}

// GrandfatherWithSiblings computes the grandfather and sibling shares from the
// portion available to them after the fixed-share heirs (available is the whole
// estate when there are no fixed-share heirs). The caller applies this only
// when the grandfather competes with full or consanguine siblings and no father
// is present.
func GrandfatherWithSiblings(available rational.Fraction, h *heir.Heirs, view JaddView) JaddResult {
	if view == JaddAbuHanifa {
		return JaddResult{
			GrandfatherShare: available,
			SiblingShares:    map[heir.Relation]rational.Fraction{},
			SiblingsExcluded: true,
			Method:           "exclusion",
		}
	}

	fullB := h.Count(heir.FullBrother)
	fullS := h.Count(heir.FullSister)
	consB := h.Count(heir.ConsanguineBrother)
	consS := h.Count(heir.ConsanguineSister)
	units := int64(2*(fullB+consB) + fullS + consS)

	if units == 0 {
		return JaddResult{
			GrandfatherShare: available,
			SiblingShares:    map[heir.Relation]rational.Fraction{},
			Method:           "no competing siblings",
		}
	}

	// The grandfather takes the best of three options. The consanguine siblings
	// are counted in the muqasama divisor even when the full siblings will then
	// absorb their portion.
	muqasama := available.Mul(rational.New(2, 2+units))
	thirdOfRemainder := available.Mul(rational.New(1, 3))
	sixthOfEstate := rational.New(1, 6)

	gf, method := muqasama, "muqasama"
	if thirdOfRemainder.Greater(gf) {
		gf, method = thirdOfRemainder, "one third of the remainder"
	}
	if sixthOfEstate.Greater(gf) {
		gf, method = sixthOfEstate, "one sixth of the estate"
	}

	siblingsPortion := available.Sub(gf)
	if siblingsPortion.IsNegative() {
		siblingsPortion = rational.Zero()
	}

	res := JaddResult{GrandfatherShare: gf, Method: method, SiblingShares: map[heir.Relation]rational.Fraction{}}

	hasFull := fullB > 0 || fullS > 0
	hasCons := consB > 0 || consS > 0
	switch {
	case hasFull && hasCons && fullB == 0:
		res.NeedsReview = true
		res.ReviewNote = "grandfather with full sisters and consanguine siblings needs manual review"
	case hasFull:
		res.SiblingShares = distributeByUnits(siblingsPortion, []heir.Relation{heir.FullBrother, heir.FullSister}, h)
	default:
		res.SiblingShares = distributeByUnits(siblingsPortion, []heir.Relation{heir.ConsanguineBrother, heir.ConsanguineSister}, h)
	}
	return res
}
