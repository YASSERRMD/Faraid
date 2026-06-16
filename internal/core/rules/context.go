package rules

import "github.com/YASSERRMD/Faraid/internal/core/heir"

// Context provides the facts that a fixed-share condition needs in order to
// select the correct fraction. It is a read-only view over the heirs present
// in the case.
//
// Conditions read presence, not post-blocking survival, because some share
// adjustments (hajb nuqsan) depend on heirs who are themselves blocked. For
// example, two or more siblings reduce the mother from one third to one sixth
// even when those siblings are excluded by the father. Total exclusion (hajb
// hirman) is computed separately and decides which present heir actually
// receives the fraction this package assigns.
type Context struct {
	Heirs *heir.Heirs
}

// Present reports whether the relation occupies at least one slot.
func (c Context) Present(r heir.Relation) bool { return c.Heirs.Present(r) }

// Count returns the number of heirs in the relation slot.
func (c Context) Count(r heir.Relation) int { return c.Heirs.Count(r) }

// HasInheritingDescendant reports whether any child or agnatic grandchild is
// present: a son, daughter, son's son, or son's daughter. Whenever such a
// descendant is present at least one of them inherits, so presence is
// sufficient. The presence of such a descendant reduces the spouse share and
// the mother share.
func (c Context) HasInheritingDescendant() bool {
	return c.Present(heir.Son) || c.Present(heir.Daughter) ||
		c.Present(heir.SonsSon) || c.Present(heir.SonsDaughter)
}

// SiblingCount returns the total number of siblings of any kind present: full,
// consanguine, and uterine, brothers and sisters alike. Two or more siblings
// reduce the mother from one third to one sixth, and this reduction applies
// even when those siblings are themselves blocked from inheriting.
func (c Context) SiblingCount() int {
	return c.Count(heir.FullBrother) + c.Count(heir.FullSister) +
		c.Count(heir.ConsanguineBrother) + c.Count(heir.ConsanguineSister) +
		c.Count(heir.UterineBrother) + c.Count(heir.UterineSister)
}

// UterineCount returns the number of uterine (maternal) siblings present,
// brothers and sisters together. They inherit a single share equally regardless
// of sex: one sixth for a single uterine sibling, one third shared for two or
// more.
func (c Context) UterineCount() int {
	return c.Count(heir.UterineBrother) + c.Count(heir.UterineSister)
}
