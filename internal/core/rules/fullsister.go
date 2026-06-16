package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Full sister shares are prescribed in Quran 4:176. A single full sister takes
// one half; two or more share two thirds. She receives a fixed share only when
// no descendant, no father, and no full brother are present: a male descendant
// or the father excludes her, a daughter or son's daughter makes her residuary
// alongside them (ma'a ghayrihi), and a full brother makes her residuary with
// him (bi-ghayrihi). Those residuary and blocking outcomes are resolved by their
// own stages.
//
// Exclusion by the paternal grandfather (jadd wa ikhwa) diverges by school and
// is handled in a dedicated later phase, so it is not guarded here.
func init() {
	eligible := func(c Context) bool {
		return !c.HasInheritingDescendant() &&
			!c.Present(heir.Father) && !c.Present(heir.FullBrother)
	}
	register(heir.FullSister,
		FixedShareRule{
			Share:     rational.New(1, 2),
			Condition: "one full sister; no descendant, father, or full brother",
			Reference: "Quran 4:176",
			When:      func(c Context) bool { return eligible(c) && c.Count(heir.FullSister) == 1 },
		},
		FixedShareRule{
			Share:     rational.New(2, 3),
			Condition: "two or more full sisters; no descendant, father, or full brother",
			Reference: "Quran 4:176",
			When:      func(c Context) bool { return eligible(c) && c.Count(heir.FullSister) >= 2 },
		},
	)
}
