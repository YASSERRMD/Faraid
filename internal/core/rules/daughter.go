package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Daughter shares are prescribed in Quran 4:11. A single daughter takes one
// half; two or more daughters share two thirds. When a son is present the
// daughters take no fixed share: they become residuary heirs with the son
// (asaba bi-ghayrihi, two to one), which the residuary stage resolves.
func init() {
	register(heir.Daughter,
		FixedShareRule{
			Share:     rational.New(1, 2),
			Condition: "one daughter and no son",
			Reference: "Quran 4:11",
			When: func(c Context) bool {
				return !c.Present(heir.Son) && c.Count(heir.Daughter) == 1
			},
		},
		FixedShareRule{
			Share:     rational.New(2, 3),
			Condition: "two or more daughters and no son",
			Reference: "Quran 4:11",
			When: func(c Context) bool {
				return !c.Present(heir.Son) && c.Count(heir.Daughter) >= 2
			},
		},
	)
}
