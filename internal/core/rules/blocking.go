package rules

import "github.com/YASSERRMD/Faraid/internal/core/heir"

// BlockRule records one way an heir can be totally excluded (hajb hirman). When
// its condition holds, the Blocked relation inherits nothing. Blockers lists
// the relations the rule depends on, used for documentation and for the
// acyclicity check; Condition and Reference feed the audit trail.
type BlockRule struct {
	Blocked   heir.Relation
	Blockers  []heir.Relation
	Condition string
	Reference string
	When      func(Context) bool
}

// NuqsanRule documents a reduction (hajb nuqsan) in which the presence of one
// heir lowers another's share without excluding it. The reduction itself is
// applied by the fixed-share table; this record exists for the derivation and
// the scholar review.
type NuqsanRule struct {
	Reduced   heir.Relation
	Reducers  []heir.Relation
	Condition string
	Reference string
}

// hajbHirman maps each excludable relation to the rules that can exclude it.
var hajbHirman = map[heir.Relation][]BlockRule{}

// registerBlock adds total-exclusion rules to the lattice.
func registerBlock(rules ...BlockRule) {
	for _, r := range rules {
		hajbHirman[r.Blocked] = append(hajbHirman[r.Blocked], r)
	}
}

// blockOnPresence builds rules that exclude each blocked relation whenever the
// single blocker is present.
func blockOnPresence(blocker heir.Relation, reference string, blocked ...heir.Relation) []BlockRule {
	rules := make([]BlockRule, 0, len(blocked))
	for _, b := range blocked {
		rules = append(rules, BlockRule{
			Blocked:   b,
			Blockers:  []heir.Relation{blocker},
			Condition: blocker.String() + " is present",
			Reference: reference,
			When:      func(c Context) bool { return c.Present(blocker) },
		})
	}
	return rules
}

// HirmanRules returns the total-exclusion rules registered for a relation. The
// returned slice must not be mutated.
func HirmanRules(r heir.Relation) []BlockRule {
	return hajbHirman[r]
}
