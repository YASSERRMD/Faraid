package rules

import "github.com/YASSERRMD/Faraid/internal/core/heir"

// Ascendant exclusions. The father excludes his own father (the paternal
// grandfather), his own mother (the paternal grandmother), and every class of
// sibling. The mother excludes all grandmothers, maternal and paternal alike.
// The paternal grandfather, standing in for the father, excludes the uterine
// siblings; his competition with full and consanguine siblings (jadd wa ikhwa)
// diverges by school and is handled in a dedicated later phase.
func init() {
	registerBlock(blockOnPresence(heir.Father, "the father excludes his ascendants and the siblings",
		heir.PaternalGrandfather, heir.PaternalGrandmother,
		heir.FullBrother, heir.FullSister,
		heir.ConsanguineBrother, heir.ConsanguineSister,
		heir.UterineBrother, heir.UterineSister,
	)...)

	registerBlock(blockOnPresence(heir.Mother, "the mother excludes the grandmothers",
		heir.MaternalGrandmother, heir.PaternalGrandmother,
	)...)

	registerBlock(blockOnPresence(heir.PaternalGrandfather, "the grandfather excludes uterine siblings",
		heir.UterineBrother, heir.UterineSister,
	)...)
}
