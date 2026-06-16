package rules

import "github.com/YASSERRMD/Faraid/internal/core/heir"

// Sibling exclusions among the collateral line. A full brother excludes the
// consanguine brother and consanguine sister. Two or more full sisters take the
// whole two thirds and exclude the consanguine sister unless a consanguine
// brother makes her residuary. A lone full sister who inherits as residuary
// alongside a female descendant (asaba ma'a ghayrihi) takes the entire residue
// and excludes the consanguine siblings.
func init() {
	registerBlock(blockOnPresence(heir.FullBrother, "a full brother excludes the consanguine siblings",
		heir.ConsanguineBrother, heir.ConsanguineSister,
	)...)

	registerBlock(BlockRule{
		Blocked:   heir.ConsanguineSister,
		Blockers:  []heir.Relation{heir.FullSister},
		Condition: "two or more full sisters and no consanguine brother",
		Reference: "the full sisters take the whole two thirds",
		When: func(c Context) bool {
			return c.Count(heir.FullSister) >= 2 && !c.Present(heir.ConsanguineBrother)
		},
	})

	// A full sister becomes residuary with a daughter or son's daughter only
	// when no son, son's son, or father excludes her and no full brother makes
	// her residuary with him. In that state she consumes the residue and
	// excludes the consanguine siblings.
	fullSisterIsResiduary := func(c Context) bool {
		return c.Present(heir.FullSister) &&
			(c.Present(heir.Daughter) || c.Present(heir.SonsDaughter)) &&
			!c.Present(heir.Son) && !c.Present(heir.SonsSon) &&
			!c.Present(heir.Father) && !c.Present(heir.FullBrother)
	}
	registerBlock(
		BlockRule{
			Blocked:   heir.ConsanguineBrother,
			Blockers:  []heir.Relation{heir.FullSister},
			Condition: "a full sister inherits as residuary with a female descendant",
			Reference: "asaba ma'a ghayrihi takes the residue",
			When:      fullSisterIsResiduary,
		},
		BlockRule{
			Blocked:   heir.ConsanguineSister,
			Blockers:  []heir.Relation{heir.FullSister},
			Condition: "a full sister inherits as residuary with a female descendant",
			Reference: "asaba ma'a ghayrihi takes the residue",
			When:      fullSisterIsResiduary,
		},
	)
}
