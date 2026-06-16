package estate

import "errors"

// Estate holds the monetary facts of a case in the smallest currency unit (for
// example fils or cents). All amounts are non-negative integers.
type Estate struct {
	// Total is the gross estate before any deduction.
	Total int64
	// Funeral is the funeral expense, settled first.
	Funeral int64
	// Debts are outstanding debts, settled after the funeral.
	Debts int64
	// Bequests is the wasiyyah requested, settled after debts and subject to
	// the one third cap.
	Bequests int64
	// HeirsConsentToExcessBequest records whether the heirs consent to
	// bequests beyond one third. Without consent, bequests are capped at one
	// third of the estate remaining after funeral and debts.
	HeirsConsentToExcessBequest bool
}

// validate reports the first problem with the monetary inputs.
func (e Estate) validate() error {
	if e.Total < 0 || e.Funeral < 0 || e.Debts < 0 || e.Bequests < 0 {
		return errors.New("estate: amounts must be non-negative")
	}
	return nil
}
