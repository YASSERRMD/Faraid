package solver

import (
	"math/big"

	"github.com/YASSERRMD/Faraid/internal/core/derivation"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// HeirShare is one heir slot's outcome.
type HeirShare struct {
	Relation heir.Relation
	Count    int
	// Fraction is the slot's share of the whole distributable estate.
	Fraction rational.Fraction
	// Parts is the slot's numerator over the result base, after tashih.
	Parts *big.Int
	// Amount is the slot's exact share of the distributable estate in the
	// smallest currency unit.
	Amount rational.Fraction
}

// Result is the full outcome of solving a case under a school.
type Result struct {
	Madhhab string
	// Distributable is the estate remaining after funeral, debts, and bequests.
	Distributable rational.Fraction
	// Base is the asl al-mas'ala after awl, radd, and tashih.
	Base *big.Int
	// Shares are the surviving heirs, in canonical relation order.
	Shares []HeirShare
	// Excluded are the relations totally excluded by blocking.
	Excluded []heir.Relation
	// SpecialCase names a special case applied, if any.
	SpecialCase string
	// Awl and Radd report whether those adjustments were applied.
	Awl  bool
	Radd bool
	// Residue is any portion passing to the public treasury.
	Residue rational.Fraction
	// NeedsReview marks a result that this engine does not fully resolve.
	NeedsReview bool
	ReviewNotes []string
	// Derivation is the ordered, step-by-step audit trail of the result.
	Derivation *derivation.Derivation
}

// shareMeta carries the non-share outcome of the share computation.
type shareMeta struct {
	specialCase string
	awl, radd   bool
	residue     rational.Fraction
	needsReview bool
	notes       []string
}
