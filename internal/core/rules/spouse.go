package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Spouse shares are prescribed in Quran 4:12. The husband takes one half, or
// one quarter when the deceased leaves an inheriting descendant. The wife (one
// or more wives sharing the slot equally) takes one quarter, or one eighth
// when there is an inheriting descendant.
func init() {
	register(heir.Husband,
		FixedShareRule{
			Share:     rational.New(1, 2),
			Condition: "no inheriting descendant",
			Reference: "Quran 4:12",
			When:      func(c Context) bool { return !c.HasInheritingDescendant() },
		},
		FixedShareRule{
			Share:     rational.New(1, 4),
			Condition: "an inheriting descendant is present",
			Reference: "Quran 4:12",
			When:      func(c Context) bool { return c.HasInheritingDescendant() },
		},
	)
	register(heir.Wife,
		FixedShareRule{
			Share:     rational.New(1, 4),
			Condition: "no inheriting descendant",
			Reference: "Quran 4:12",
			When:      func(c Context) bool { return !c.HasInheritingDescendant() },
		},
		FixedShareRule{
			Share:     rational.New(1, 8),
			Condition: "an inheriting descendant is present",
			Reference: "Quran 4:12",
			When:      func(c Context) bool { return c.HasInheritingDescendant() },
		},
	)
}
