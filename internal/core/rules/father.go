package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Father shares are prescribed in Quran 4:11. The father takes a fixed one
// sixth when an inheriting descendant is present. When the descendant is female
// only (daughters or son's daughters, with no son or son's son) he also takes
// the residue as a residuary heir, which the residuary stage adds. When there
// is no descendant at all he is a pure residuary heir and has no fixed share,
// so no rule matches here.
func init() {
	register(heir.Father,
		FixedShareRule{
			Share:     rational.New(1, 6),
			Condition: "an inheriting descendant is present",
			Reference: "Quran 4:11",
			When:      func(c Context) bool { return c.HasInheritingDescendant() },
		},
	)
}
