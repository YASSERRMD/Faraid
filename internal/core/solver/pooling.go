package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// sumFractions returns the sum of all the share fractions.
func sumFractions(m map[heir.Relation]rational.Fraction) rational.Fraction {
	s := rational.Zero()
	for _, f := range m {
		s = s.Add(f)
	}
	return s
}

// hasNonSpouseFurud reports whether any fixed-share heir other than a spouse is
// present, which decides whether radd has eligible recipients.
func hasNonSpouseFurud(m map[heir.Relation]rational.Fraction) bool {
	for r := range m {
		if !r.IsSpouse() {
			return true
		}
	}
	return false
}

// isCompetingSibling reports whether the relation competes with the grandfather
// in the jadd wa ikhwa case (uterine siblings are excluded by the grandfather).
func isCompetingSibling(r heir.Relation) bool {
	switch r {
	case heir.FullBrother, heir.FullSister, heir.ConsanguineBrother, heir.ConsanguineSister:
		return true
	default:
		return false
	}
}

// poolGrandmothers pools the grandmother fixed share. The fixed-share table
// assigns one sixth to each grandmother slot, but grandmothers at the same
// level share a single one sixth, so when more than one is present the one
// sixth is split equally between them.
func poolGrandmothers(fixed map[heir.Relation]rational.Fraction) {
	var present []heir.Relation
	for _, r := range []heir.Relation{heir.MaternalGrandmother, heir.PaternalGrandmother} {
		if _, ok := fixed[r]; ok {
			present = append(present, r)
		}
	}
	if len(present) <= 1 {
		return
	}
	each := rational.New(1, 6).Div(rational.FromInt(int64(len(present))))
	for _, r := range present {
		fixed[r] = each
	}
}

// poolUterine pools the uterine sibling share. The fixed-share table assigns the
// whole uterine share to each uterine slot, so when both brothers and sisters
// are present the single share is split equally per head between the two slots.
func poolUterine(fixed map[heir.Relation]rational.Fraction, h *heir.Heirs) {
	groupB, hasB := fixed[heir.UterineBrother]
	_, hasS := fixed[heir.UterineSister]
	if !hasB || !hasS {
		return
	}
	group := groupB // both slots hold the same group share
	uB := int64(h.Count(heir.UterineBrother))
	uS := int64(h.Count(heir.UterineSister))
	total := rational.FromInt(uB + uS)
	fixed[heir.UterineBrother] = group.Mul(rational.FromInt(uB)).Div(total)
	fixed[heir.UterineSister] = group.Mul(rational.FromInt(uS)).Div(total)
}
