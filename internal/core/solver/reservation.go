package solver

import "github.com/YASSERRMD/Faraid/internal/core/rational"

// Scenario is one possible resolution of an uncertainty, as a complete
// distribution that sums to the whole estate. Heirs are keyed by a stable
// identity label, because an uncertain heir may appear as different relations
// across scenarios (for example a male or female intersex heir).
type Scenario struct {
	Name   string
	Shares map[string]rational.Fraction
}

// Reservation is the cautious outcome over several scenarios. Each heir is
// guaranteed the least it receives in any scenario, and the remainder is held
// until the uncertainty resolves. It serves the unborn child (haml), the
// missing person (mafqud), and the intersex heir (khuntha).
type Reservation struct {
	// Guaranteed is each heir's minimum share across all scenarios.
	Guaranteed map[string]rational.Fraction
	// Reserved is the portion held back, equal to the whole minus the sum of
	// the guaranteed shares.
	Reserved rational.Fraction
	// Scenarios are the considered, provisional distributions.
	Scenarios []Scenario
}

// Reserve computes the cautious distribution: every heir is guaranteed the
// minimum it receives in any scenario, and the rest is reserved. The scenarios
// should each be a complete distribution of the estate.
func Reserve(scenarios []Scenario) Reservation {
	labels := map[string]bool{}
	for _, s := range scenarios {
		for l := range s.Shares {
			labels[l] = true
		}
	}

	guaranteed := make(map[string]rational.Fraction, len(labels))
	for l := range labels {
		lowest := rational.One()
		for _, s := range scenarios {
			if share := s.Shares[l]; share.Less(lowest) {
				lowest = share
			}
		}
		guaranteed[l] = lowest
	}

	sum := rational.Zero()
	for _, f := range guaranteed {
		sum = sum.Add(f)
	}

	return Reservation{
		Guaranteed: guaranteed,
		Reserved:   rational.One().Sub(sum),
		Scenarios:  scenarios,
	}
}
