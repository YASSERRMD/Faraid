package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Consanguine sister shares mirror the full sister (Quran 4:176) when no full
// sibling stands above her. One takes one half, two or more share two thirds,
// when there is no descendant, father, full brother, consanguine brother, or
// full sister. When exactly one full sister is present (taking one half), the
// consanguine sisters together take one sixth to complete the group to two
// thirds. Two or more full sisters consume the whole two thirds and leave the
// consanguine sister with no fixed share unless a consanguine brother makes her
// residuary. Exclusion by the paternal grandfather is deferred to the dedicated
// later phase.
func init() {
	eligible := func(c Context) bool {
		return !c.HasInheritingDescendant() &&
			!c.Present(heir.Father) &&
			!c.Present(heir.FullBrother) &&
			!c.Present(heir.ConsanguineBrother)
	}
	register(heir.ConsanguineSister,
		FixedShareRule{
			Share:     rational.New(1, 2),
			Condition: "one consanguine sister; no full sister and no blocking heir",
			Reference: "Quran 4:176 by analogy to the full sister",
			When: func(c Context) bool {
				return eligible(c) && !c.Present(heir.FullSister) && c.Count(heir.ConsanguineSister) == 1
			},
		},
		FixedShareRule{
			Share:     rational.New(2, 3),
			Condition: "two or more consanguine sisters; no full sister and no blocking heir",
			Reference: "Quran 4:176 by analogy to the full sister",
			When: func(c Context) bool {
				return eligible(c) && !c.Present(heir.FullSister) && c.Count(heir.ConsanguineSister) >= 2
			},
		},
		FixedShareRule{
			Share:     rational.New(1, 6),
			Condition: "exactly one full sister present, completing the group to two thirds",
			Reference: "Sunnah, completion to two thirds (takmila)",
			When: func(c Context) bool {
				return eligible(c) && c.Count(heir.FullSister) == 1
			},
		},
	)
}
