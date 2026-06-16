package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
	"github.com/YASSERRMD/Faraid/internal/core/rules"
)

// MushtarakaView selects how a school treats the shared case (mushtaraka, also
// called himariyyah): a husband, the mother, two or more uterine siblings, and
// full brothers, where the full brothers would otherwise take nothing.
type MushtarakaView int

const (
	// MushtarakaShare is the view of the Maliki and Shafi'i schools: the full
	// siblings join the uterine siblings in their one third, sharing it equally
	// per head through their common mother.
	MushtarakaShare MushtarakaView = iota
	// MushtarakaNoShare is the view of the Hanafi and Hanbali schools: the full
	// brothers take nothing, since the fixed shares consume the whole estate.
	MushtarakaNoShare
)

// String returns a label for the view.
func (v MushtarakaView) String() string {
	if v == MushtarakaNoShare {
		return "no-share"
	}
	return "share"
}

// IsMushtaraka reports whether the heirs form the shared case: a husband, the
// mother, two or more uterine siblings, and at least one full brother, with no
// descendant and no father. In that configuration the husband, mother, and
// uterine siblings consume the whole estate, leaving the full brothers, who are
// residuary heirs, with nothing.
func IsMushtaraka(h *heir.Heirs) bool {
	ctx := rules.Context{Heirs: h}
	return h.Present(heir.Husband) &&
		h.Present(heir.Mother) &&
		(h.Count(heir.UterineBrother)+h.Count(heir.UterineSister)) >= 2 &&
		h.Present(heir.FullBrother) &&
		!ctx.HasInheritingDescendant() &&
		!h.Present(heir.Father)
}

// Mushtaraka returns the shares for the shared case under the given view, or ok
// false when the heirs are not that configuration.
func Mushtaraka(h *heir.Heirs, view MushtarakaView) (map[heir.Relation]rational.Fraction, bool) {
	if !IsMushtaraka(h) {
		return nil, false
	}

	shares := map[heir.Relation]rational.Fraction{
		heir.Husband: rational.New(1, 2),
		heir.Mother:  rational.New(1, 6),
	}
	pool := rational.New(1, 3)

	members := []heir.Relation{heir.UterineBrother, heir.UterineSister}
	if view == MushtarakaShare {
		members = append(members, heir.FullBrother, heir.FullSister)
	}
	for r, f := range distributeEqually(pool, members, h) {
		shares[r] = f
	}
	return shares, true
}

// distributeEqually splits a portion equally per head among the given
// relations, regardless of sex, and returns each relation's share of the whole.
func distributeEqually(portion rational.Fraction, members []heir.Relation, h *heir.Heirs) map[heir.Relation]rational.Fraction {
	out := map[heir.Relation]rational.Fraction{}
	var heads int64
	for _, r := range members {
		heads += int64(h.Count(r))
	}
	if heads == 0 || portion.IsZero() {
		return out
	}
	total := rational.FromInt(heads)
	for _, r := range members {
		c := int64(h.Count(r))
		if c == 0 {
			continue
		}
		out[r] = portion.Mul(rational.FromInt(c)).Div(total)
	}
	return out
}
