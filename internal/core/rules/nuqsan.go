package rules

import "github.com/YASSERRMD/Faraid/internal/core/heir"

// hajbNuqsan records the reductions in which one heir lowers another's share
// without excluding it. The fixed-share table applies these reductions through
// its conditions; this list documents them for the derivation and the scholar
// review. The reducers are counted by presence, which is why a sibling reduces
// the mother even when that sibling is itself excluded.
var hajbNuqsan = []NuqsanRule{
	{
		Reduced:   heir.Husband,
		Reducers:  []heir.Relation{heir.Son, heir.Daughter, heir.SonsSon, heir.SonsDaughter},
		Condition: "an inheriting descendant lowers the husband from one half to one quarter",
		Reference: "Quran 4:12",
	},
	{
		Reduced:   heir.Wife,
		Reducers:  []heir.Relation{heir.Son, heir.Daughter, heir.SonsSon, heir.SonsDaughter},
		Condition: "an inheriting descendant lowers the wife from one quarter to one eighth",
		Reference: "Quran 4:12",
	},
	{
		Reduced: heir.Mother,
		Reducers: []heir.Relation{
			heir.Son, heir.Daughter, heir.SonsSon, heir.SonsDaughter,
			heir.FullBrother, heir.FullSister,
			heir.ConsanguineBrother, heir.ConsanguineSister,
			heir.UterineBrother, heir.UterineSister,
		},
		Condition: "an inheriting descendant, or two or more siblings, lower the mother from one third to one sixth",
		Reference: "Quran 4:11",
	},
}

// NuqsanRules returns the documented share reductions. The returned slice must
// not be mutated.
func NuqsanRules() []NuqsanRule {
	return hajbNuqsan
}
