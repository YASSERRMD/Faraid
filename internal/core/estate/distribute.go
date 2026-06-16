package estate

import "github.com/YASSERRMD/Faraid/internal/core/rational"

// Distribution is the result of settling an estate before faraid. All amounts
// are exact rationals in the smallest currency unit.
type Distribution struct {
	// AfterFuneralAndDebts is the estate remaining after funeral and debts,
	// never negative.
	AfterFuneralAndDebts rational.Fraction
	// BequestCap is one third of AfterFuneralAndDebts, the classical thuluth.
	BequestCap rational.Fraction
	// EffectiveBequests is the bequest actually paid after applying the cap.
	EffectiveBequests rational.Fraction
	// BequestsClamped reports whether the requested bequest was reduced to
	// honor the one third cap. This is only possible without heir consent.
	BequestsClamped bool
	// Distributable is the remainder allocated to heirs by faraid.
	Distributable rational.Fraction
}

var oneThird = rational.New(1, 3)

// PreDistribute settles the estate in the classical order (funeral, debts,
// bequests) and returns the distributable remainder together with the
// breakdown. It returns an error only when the inputs are invalid.
func (e Estate) PreDistribute() (Distribution, error) {
	if err := e.validate(); err != nil {
		return Distribution{}, err
	}

	net := e.Total - e.Funeral - e.Debts
	if net < 0 {
		net = 0
	}
	netF := rational.FromInt(net)
	bequestCap := netF.Mul(oneThird)

	// A bequest can never exceed what remains after funeral and debts.
	requested := rational.FromInt(e.Bequests)
	if requested.Greater(netF) {
		requested = netF
	}

	effective := requested
	clamped := false
	if !e.HeirsConsentToExcessBequest && requested.Greater(bequestCap) {
		effective = bequestCap
		clamped = true
	}

	return Distribution{
		AfterFuneralAndDebts: netF,
		BequestCap:           bequestCap,
		EffectiveBequests:    effective,
		BequestsClamped:      clamped,
		Distributable:        netF.Sub(effective),
	}, nil
}
