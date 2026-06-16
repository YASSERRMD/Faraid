package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rules"
)

// Reduction records that an heir takes the reduced (nuqsan) variant of its
// fixed share because of the heirs present. The reduced fraction itself is
// produced by the fixed-share table; this marks why it is reduced for the
// derivation.
type Reduction struct {
	Relation  heir.Relation
	By        []heir.Relation
	Reason    string
	Reference string
}

// Reductions reports which present heirs take the reduced variant of their
// fixed share. The spouse is reduced by an inheriting descendant; the mother is
// reduced by an inheriting descendant or by two or more siblings.
func Reductions(h *heir.Heirs) []Reduction {
	ctx := rules.Context{Heirs: h}
	var out []Reduction

	if spouse, ok := presentSpouse(h); ok && ctx.HasInheritingDescendant() {
		out = append(out, reductionFor(spouse))
	}
	if h.Present(heir.Mother) && (ctx.HasInheritingDescendant() || ctx.SiblingCount() >= 2) {
		out = append(out, reductionFor(heir.Mother))
	}
	return out
}

// presentSpouse returns the spouse relation present in the case, if any.
func presentSpouse(h *heir.Heirs) (heir.Relation, bool) {
	switch {
	case h.Present(heir.Husband):
		return heir.Husband, true
	case h.Present(heir.Wife):
		return heir.Wife, true
	default:
		return heir.RelationInvalid, false
	}
}

// reductionFor builds a Reduction from the documented nuqsan rule for the
// relation, so the description and reference have a single source.
func reductionFor(r heir.Relation) Reduction {
	for _, n := range rules.NuqsanRules() {
		if n.Reduced == r {
			return Reduction{Relation: r, By: n.Reducers, Reason: n.Condition, Reference: n.Reference}
		}
	}
	return Reduction{Relation: r}
}
