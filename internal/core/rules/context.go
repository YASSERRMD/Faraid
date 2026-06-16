package rules

import "github.com/YASSERRMD/Faraid/internal/core/heir"

// Context provides the facts that a fixed-share condition needs in order to
// select the correct fraction. It is a read-only view over the heirs that are
// effectively present. The solver evaluates conditions against the set that
// survives blocking, so a condition sees only heirs who actually inherit.
type Context struct {
	Heirs *heir.Heirs
}

// Present reports whether the relation occupies at least one slot.
func (c Context) Present(r heir.Relation) bool { return c.Heirs.Present(r) }

// Count returns the number of heirs in the relation slot.
func (c Context) Count(r heir.Relation) int { return c.Heirs.Count(r) }

// HasInheritingDescendant reports whether any child or agnatic grandchild is
// present: a son, daughter, son's son, or son's daughter. The presence of such
// a descendant reduces the spouse share and, with siblings, the mother share.
func (c Context) HasInheritingDescendant() bool {
	return c.Present(heir.Son) || c.Present(heir.Daughter) ||
		c.Present(heir.SonsSon) || c.Present(heir.SonsDaughter)
}
