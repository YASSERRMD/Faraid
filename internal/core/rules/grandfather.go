package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Paternal grandfather shares follow the father in the father's absence: a
// fixed one sixth when an inheriting descendant is present, plus the residue as
// a residuary heir when there is no male descendant, and pure residuary when
// there is no descendant. The grandfather is excluded by the father, which the
// blocking stage applies.
//
// The grandfather competing with siblings (jadd wa ikhwa) involves a comparison
// between a fixed fraction and sharing with the siblings (muqasama) that
// diverges by school. That is implemented in a dedicated later phase; this rule
// covers only the base one sixth with a descendant.
func init() {
	register(heir.PaternalGrandfather,
		FixedShareRule{
			Share:     rational.New(1, 6),
			Condition: "an inheriting descendant is present",
			Reference: "Quran 4:11 by analogy to the father",
			When:      func(c Context) bool { return c.HasInheritingDescendant() },
		},
	)
}
