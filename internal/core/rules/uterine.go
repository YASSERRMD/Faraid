package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Uterine sibling shares are prescribed in Quran 4:12. A single uterine sibling
// takes one sixth; two or more share one third equally, with no two to one
// preference between brothers and sisters. Uterine siblings are excluded by any
// inheriting descendant and by the father and the paternal grandfather; that
// exclusion is applied by the blocking stage, and the conditions here also
// guard against it so the rule is correct on its own.
//
// The one third (or one sixth) is the share of the uterine group as a whole.
// When both uterine brothers and sisters are present the pooling of that single
// share across the group is applied during share assembly.
func init() {
	notBlocked := func(c Context) bool {
		return !c.HasInheritingDescendant() &&
			!c.Present(heir.Father) && !c.Present(heir.PaternalGrandfather)
	}
	uterineRules := []FixedShareRule{
		{
			Share:     rational.New(1, 6),
			Condition: "a single uterine sibling; no descendant, father, or grandfather",
			Reference: "Quran 4:12",
			When:      func(c Context) bool { return notBlocked(c) && c.UterineCount() == 1 },
		},
		{
			Share:     rational.New(1, 3),
			Condition: "two or more uterine siblings sharing equally; no descendant, father, or grandfather",
			Reference: "Quran 4:12",
			When:      func(c Context) bool { return notBlocked(c) && c.UterineCount() >= 2 },
		},
	}
	register(heir.UterineBrother, uterineRules...)
	register(heir.UterineSister, uterineRules...)
}
