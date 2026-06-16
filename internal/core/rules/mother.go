package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Mother shares are prescribed in Quran 4:11. The mother takes one sixth when
// an inheriting descendant is present or when two or more siblings are present,
// and one third otherwise.
//
// The one third applies to the whole estate except in the two Umariyyatan
// (gharrawayn) configurations detected by IsGharrawayn, where it applies to the
// residue after the spouse. That adjustment is applied by a later phase; the
// fraction recorded here is the base one third.
func init() {
	register(heir.Mother,
		FixedShareRule{
			Share:     rational.New(1, 6),
			Condition: "an inheriting descendant is present, or two or more siblings",
			Reference: "Quran 4:11",
			When: func(c Context) bool {
				return c.HasInheritingDescendant() || c.SiblingCount() >= 2
			},
		},
		FixedShareRule{
			Share:     rational.New(1, 3),
			Condition: "no inheriting descendant and fewer than two siblings",
			Reference: "Quran 4:11",
			When: func(c Context) bool {
				return !c.HasInheritingDescendant() && c.SiblingCount() < 2
			},
		},
	)
}

// IsGharrawayn reports whether the heirs form one of the two Umariyyatan
// (gharrawayn) configurations: a spouse together with both parents, with no
// inheriting descendant and no siblings. In that configuration the mother takes
// one third of the residue remaining after the spouse rather than one third of
// the whole estate, which preserves the father taking twice the mother. The
// adjustment itself is applied by a later phase.
func IsGharrawayn(c Context) bool {
	return (c.Present(heir.Husband) || c.Present(heir.Wife)) &&
		c.Present(heir.Father) && c.Present(heir.Mother) &&
		!c.HasInheritingDescendant() && c.SiblingCount() == 0
}
