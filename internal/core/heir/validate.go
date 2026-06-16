package heir

import (
	"errors"
	"fmt"
)

// Validate checks the heir set for structural impossibilities and returns the
// first problem found, or nil when the set is internally consistent. It does
// not apply any legal ruling: it only rejects inputs that cannot describe a
// real family. For example, a present father does not block a present paternal
// grandfather here; that is a blocking rule resolved by the solver, not an
// impossibility, so both may appear together.
func (h *Heirs) Validate() error {
	for r, count := range h.counts {
		if !r.Valid() {
			return fmt.Errorf("heir: unknown relation %d", int(r))
		}
		if count < 0 {
			return fmt.Errorf("heir: negative count %d for %s", count, r)
		}
		if max := r.MaxCount(); max > 0 && count > max {
			return fmt.Errorf("heir: %s count %d exceeds maximum %d", r, count, max)
		}
	}
	if h.Present(Husband) && h.Present(Wife) {
		return errors.New("heir: a case cannot have both a husband and a wife")
	}
	return nil
}
