package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Son's daughter shares mirror the daughter when no daughter stands above her.
// One son's daughter takes one half, two or more share two thirds. When there
// is exactly one daughter, the son's daughters together take one sixth to
// complete the female group to two thirds (takmila). Two or more daughters
// consume the whole two thirds and leave the son's daughter with no fixed
// share, unless a son's son makes her residuary. A son, or a son's son, also
// leaves her with no fixed share. Those residuary and blocking outcomes are
// resolved by their own stages; here the conditions simply yield no furud.
func init() {
	register(heir.SonsDaughter,
		FixedShareRule{
			Share:     rational.New(1, 2),
			Condition: "one son's daughter; no son, no daughter, no son's son",
			Reference: "Quran 4:11 by analogy to the daughter",
			When: func(c Context) bool {
				return !c.Present(heir.Son) && !c.Present(heir.SonsSon) &&
					!c.Present(heir.Daughter) && c.Count(heir.SonsDaughter) == 1
			},
		},
		FixedShareRule{
			Share:     rational.New(2, 3),
			Condition: "two or more son's daughters; no son, no daughter, no son's son",
			Reference: "Quran 4:11 by analogy to the daughter",
			When: func(c Context) bool {
				return !c.Present(heir.Son) && !c.Present(heir.SonsSon) &&
					!c.Present(heir.Daughter) && c.Count(heir.SonsDaughter) >= 2
			},
		},
		FixedShareRule{
			Share:     rational.New(1, 6),
			Condition: "exactly one daughter present, completing the group to two thirds; no son, no son's son",
			Reference: "Sunnah, the ruling of Ibn Mas'ud (completion to two thirds)",
			When: func(c Context) bool {
				return !c.Present(heir.Son) && !c.Present(heir.SonsSon) &&
					c.Count(heir.Daughter) == 1
			},
		},
	)
}
