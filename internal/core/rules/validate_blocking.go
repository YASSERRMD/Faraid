package rules

import (
	"fmt"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

// ValidateLattice checks that every blocking and reduction rule is well-formed
// and that the total-exclusion lattice is acyclic. A cycle would mean two heirs
// each exclude the other, which is never valid.
func ValidateLattice() error {
	return validateLattice(hajbHirman, hajbNuqsan)
}

// validateLattice is the testable core of ValidateLattice, operating on the
// given rule sets rather than the package globals.
func validateLattice(hajbHirman map[heir.Relation][]BlockRule, hajbNuqsan []NuqsanRule) error {
	edges := map[heir.Relation][]heir.Relation{}

	for blocked, rules := range hajbHirman {
		if !blocked.Valid() {
			return fmt.Errorf("rules: hirman has unknown blocked relation %d", int(blocked))
		}
		for i, r := range rules {
			switch {
			case r.Blocked != blocked:
				return fmt.Errorf("rules: hirman rule for %s is filed under %s", r.Blocked, blocked)
			case r.When == nil:
				return fmt.Errorf("rules: hirman rule %d for %s has a nil condition", i, blocked)
			case r.Condition == "":
				return fmt.Errorf("rules: hirman rule %d for %s has an empty condition", i, blocked)
			case r.Reference == "":
				return fmt.Errorf("rules: hirman rule %d for %s has an empty reference", i, blocked)
			case len(r.Blockers) == 0:
				return fmt.Errorf("rules: hirman rule %d for %s lists no blockers", i, blocked)
			}
			for _, b := range r.Blockers {
				if !b.Valid() {
					return fmt.Errorf("rules: hirman rule %d for %s has unknown blocker %d", i, blocked, int(b))
				}
				edges[b] = append(edges[b], r.Blocked)
			}
		}
	}

	for _, n := range hajbNuqsan {
		switch {
		case !n.Reduced.Valid():
			return fmt.Errorf("rules: nuqsan has unknown reduced relation %d", int(n.Reduced))
		case n.Condition == "":
			return fmt.Errorf("rules: nuqsan rule for %s has an empty condition", n.Reduced)
		case n.Reference == "":
			return fmt.Errorf("rules: nuqsan rule for %s has an empty reference", n.Reduced)
		case len(n.Reducers) == 0:
			return fmt.Errorf("rules: nuqsan rule for %s lists no reducers", n.Reduced)
		}
		for _, b := range n.Reducers {
			if !b.Valid() {
				return fmt.Errorf("rules: nuqsan rule for %s has unknown reducer %d", n.Reduced, int(b))
			}
		}
	}

	return detectCycle(edges)
}

// detectCycle reports a cycle in the directed graph of blocker to blocked
// edges, if any, using a depth-first three-color walk.
func detectCycle(edges map[heir.Relation][]heir.Relation) error {
	const (
		white = 0 // unvisited
		gray  = 1 // on the current path
		black = 2 // fully explored
	)
	color := map[heir.Relation]int{}

	var visit func(n heir.Relation) error
	visit = func(n heir.Relation) error {
		color[n] = gray
		for _, m := range edges[n] {
			switch color[m] {
			case gray:
				return fmt.Errorf("rules: blocking lattice has a cycle through %s and %s", n, m)
			case white:
				if err := visit(m); err != nil {
					return err
				}
			}
		}
		color[n] = black
		return nil
	}

	for n := range edges {
		if color[n] == white {
			if err := visit(n); err != nil {
				return err
			}
		}
	}
	return nil
}
