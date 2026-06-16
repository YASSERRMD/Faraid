package rules

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// FixedShareRule assigns a prescribed fraction (furud) to an heir when its
// condition holds. The rules registered for a relation are evaluated in order,
// and the first whose When returns true applies. Each rule records a human
// readable condition and a source reference so the result can be explained and
// audited.
type FixedShareRule struct {
	Share     rational.Fraction
	Condition string
	Reference string
	When      func(Context) bool
}

// fixedShareTable maps each fixed-share heir to its ordered rules. It is
// populated by per heir-group init functions through register.
var fixedShareTable = map[heir.Relation][]FixedShareRule{}

// register appends fixed-share rules for a relation. It is called from the
// init function of each heir-group file.
func register(r heir.Relation, rules ...FixedShareRule) {
	fixedShareTable[r] = append(fixedShareTable[r], rules...)
}

// FixedShare returns the prescribed fraction for the relation under the given
// context, the rule that matched, and whether any rule matched. When no rule
// matches, the relation has no fixed share in this context: it is either
// blocked or a residuary heir, which is resolved elsewhere in the pipeline.
func FixedShare(r heir.Relation, ctx Context) (rational.Fraction, FixedShareRule, bool) {
	for _, rule := range fixedShareTable[r] {
		if rule.When(ctx) {
			return rule.Share, rule, true
		}
	}
	return rational.Zero(), FixedShareRule{}, false
}

// Rules returns the fixed-share rules registered for a relation, for
// inspection, documentation, and the scholar review pack. The returned slice
// must not be mutated.
func Rules(r heir.Relation) []FixedShareRule {
	return fixedShareTable[r]
}
