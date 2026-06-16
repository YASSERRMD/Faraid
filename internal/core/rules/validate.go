package rules

import (
	"fmt"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// canonicalFurud is the set of prescribed fractions recognized in the Quran:
// 1/2, 1/4, 1/8, 2/3, 1/3, and 1/6. Every fixed-share rule must assign one of
// these.
var canonicalFurud = []rational.Fraction{
	rational.New(1, 2),
	rational.New(1, 4),
	rational.New(1, 8),
	rational.New(2, 3),
	rational.New(1, 3),
	rational.New(1, 6),
}

func isCanonicalFurud(f rational.Fraction) bool {
	for _, v := range canonicalFurud {
		if f.Equal(v) {
			return true
		}
	}
	return false
}

// validateTable checks that every rule in the given table is well-formed: a
// known relation, at least one rule, and each rule carrying a condition
// description, a source reference, a non-nil predicate, and a canonical furud.
func validateTable(table map[heir.Relation][]FixedShareRule) error {
	for r, rules := range table {
		if !r.Valid() {
			return fmt.Errorf("rules: table has unknown relation %d", int(r))
		}
		if len(rules) == 0 {
			return fmt.Errorf("rules: relation %s has no rules", r)
		}
		for i, rule := range rules {
			switch {
			case rule.When == nil:
				return fmt.Errorf("rules: %s rule %d has a nil condition", r, i)
			case rule.Condition == "":
				return fmt.Errorf("rules: %s rule %d has an empty condition description", r, i)
			case rule.Reference == "":
				return fmt.Errorf("rules: %s rule %d has an empty source reference", r, i)
			case !isCanonicalFurud(rule.Share):
				return fmt.Errorf("rules: %s rule %d has non-canonical share %s", r, i, rule.Share)
			}
		}
	}
	return nil
}

// ValidateTable validates the package fixed-share table. It is exercised by the
// tests so a malformed table fails the build.
func ValidateTable() error {
	return validateTable(fixedShareTable)
}
