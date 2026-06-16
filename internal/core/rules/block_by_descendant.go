package rules

import "github.com/YASSERRMD/Faraid/internal/core/heir"

// Descendant exclusions. A son, being the nearest agnatic descendant, excludes
// the son's son and son's daughter below him and every class of sibling. A
// son's son (in the absence of a son) likewise excludes the siblings. Uterine
// siblings are excluded by any inheriting descendant, male or female, so a
// daughter or son's daughter excludes them too. Two or more daughters consume
// the daughters' two thirds and exclude the son's daughter unless a son's son
// makes her residuary.
func init() {
	registerBlock(blockOnPresence(heir.Son, "a son excludes lower descendants and all siblings",
		heir.SonsSon, heir.SonsDaughter,
		heir.FullBrother, heir.FullSister,
		heir.ConsanguineBrother, heir.ConsanguineSister,
		heir.UterineBrother, heir.UterineSister,
	)...)

	registerBlock(blockOnPresence(heir.SonsSon, "a son's son excludes the siblings",
		heir.FullBrother, heir.FullSister,
		heir.ConsanguineBrother, heir.ConsanguineSister,
		heir.UterineBrother, heir.UterineSister,
	)...)

	registerBlock(blockOnPresence(heir.Daughter, "an inheriting descendant excludes uterine siblings",
		heir.UterineBrother, heir.UterineSister,
	)...)

	registerBlock(blockOnPresence(heir.SonsDaughter, "an inheriting descendant excludes uterine siblings",
		heir.UterineBrother, heir.UterineSister,
	)...)

	registerBlock(BlockRule{
		Blocked:   heir.SonsDaughter,
		Blockers:  []heir.Relation{heir.Daughter, heir.SonsSon},
		Condition: "two or more daughters and no son's son",
		Reference: "the daughters take the whole two thirds",
		When: func(c Context) bool {
			return c.Count(heir.Daughter) >= 2 && !c.Present(heir.SonsSon)
		},
	})
}
