package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rules"
)

// Exclusion records that a relation was totally excluded (hajb hirman) and why.
type Exclusion struct {
	Relation  heir.Relation
	By        []heir.Relation
	Reason    string
	Reference string
}

// BlockingResult is the outcome of resolving total exclusion over a case.
type BlockingResult struct {
	// Excluded maps each excluded relation to the rule that excluded it.
	Excluded map[heir.Relation]Exclusion
	// Surviving holds the heirs who continue to inherit, with their counts.
	Surviving *heir.Heirs
}

// IsExcluded reports whether the relation was totally excluded.
func (b BlockingResult) IsExcluded(r heir.Relation) bool {
	_, ok := b.Excluded[r]
	return ok
}

// ResolveBlocking applies the hajb hirman lattice to the heirs present in the
// case and returns the excluded relations and the surviving set. A relation is
// excluded when any of its total-exclusion rules fires.
//
// Conditions are evaluated against the heirs present, not a partially reduced
// set. The lattice is constructed so that whenever a blocker is itself
// excluded, a higher heir excludes the same targets, which keeps the principle
// that an excluded heir does not exclude others (al-mahjub la yahjub) intact
// without a separate fixpoint pass.
func ResolveBlocking(h *heir.Heirs) BlockingResult {
	ctx := rules.Context{Heirs: h}
	excluded := map[heir.Relation]Exclusion{}
	surviving := heir.New()

	for _, r := range h.Relations() {
		if ex, blocked := firstFiringRule(r, ctx); blocked {
			excluded[r] = ex
		} else {
			surviving.Set(r, h.Count(r))
		}
	}

	return BlockingResult{Excluded: excluded, Surviving: surviving}
}

// firstFiringRule returns the first total-exclusion rule for r that fires under
// the context, if any.
func firstFiringRule(r heir.Relation, ctx rules.Context) (Exclusion, bool) {
	for _, rule := range rules.HirmanRules(r) {
		if rule.When(ctx) {
			return Exclusion{
				Relation:  r,
				By:        rule.Blockers,
				Reason:    rule.Condition,
				Reference: rule.Reference,
			}, true
		}
	}
	return Exclusion{}, false
}
