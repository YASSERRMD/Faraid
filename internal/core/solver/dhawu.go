package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// DistantKindred identifies a distant-kindred (dhawu al-arham) relation: a
// relative who is neither a fixed-share heir nor a residuary heir. They inherit
// only when no fixed-share heir other than a spouse and no residuary heir
// exists, and only under schools that admit them.
type DistantKindred int

const (
	// Class 1: descendants of daughters and of son's daughters.
	DaughtersSon DistantKindred = iota
	DaughtersDaughter
	// Class 2: false (non-agnatic) ascendants.
	MaternalGrandfather
	// Class 3: non-agnatic descendants of siblings.
	SistersSon
	SistersDaughter
	// Class 4: aunts and non-agnatic uncles.
	MaternalUncle
	PaternalAunt
	MaternalAunt
)

type dkInfo struct {
	name  string
	class int
	sex   heir.Sex
}

var dkTable = map[DistantKindred]dkInfo{
	DaughtersSon:        {"daughter's son", 1, heir.Male},
	DaughtersDaughter:   {"daughter's daughter", 1, heir.Female},
	MaternalGrandfather: {"maternal grandfather", 2, heir.Male},
	SistersSon:          {"sister's son", 3, heir.Male},
	SistersDaughter:     {"sister's daughter", 3, heir.Female},
	MaternalUncle:       {"maternal uncle", 4, heir.Male},
	PaternalAunt:        {"paternal aunt", 4, heir.Female},
	MaternalAunt:        {"maternal aunt", 4, heir.Female},
}

// Class returns the inheritance class (1 to 4) of the distant kindred, nearer
// classes excluding farther ones.
func (d DistantKindred) Class() int { return dkTable[d].class }

// Sex returns the sex of the distant kindred.
func (d DistantKindred) Sex() heir.Sex { return dkTable[d].sex }

// String returns the label for the distant kindred.
func (d DistantKindred) String() string {
	if info, ok := dkTable[d]; ok {
		return info.name
	}
	return "unknown distant kindred"
}

// DhawuResult is the outcome of a distant-kindred distribution.
type DhawuResult struct {
	// SpouseShare is the spouse's fixed share, kept ahead of the distribution.
	SpouseShare rational.Fraction
	// Shares are the distant kindred who inherit and their shares of the whole.
	Shares map[DistantKindred]rational.Fraction
	// ToTreasury is the residue passed to the public treasury when the school
	// does not admit distant kindred.
	ToTreasury rational.Fraction
	// NeedsReview marks distributions whose detailed ruling this engine does
	// not yet resolve.
	NeedsReview bool
	Note        string
}

// DistributeDhawuArham distributes the estate when only distant kindred, and
// possibly a spouse, inherit. The spouse keeps the given fixed share; the
// residue goes to the public treasury under schools that exclude distant
// kindred, or to the nearest class of distant kindred under schools that admit
// them, split two to one by sex within that class.
func DistributeDhawuArham(spouseShare rational.Fraction, kindred map[DistantKindred]int, m Madhhab) DhawuResult {
	residue := rational.One().Sub(spouseShare)
	res := DhawuResult{SpouseShare: spouseShare, Shares: map[DistantKindred]rational.Fraction{}}

	present := make([]DistantKindred, 0, len(kindred))
	for dk, c := range kindred {
		if c > 0 {
			present = append(present, dk)
		}
	}

	if len(present) == 0 || m.DhawuArham == DhawuArhamExcluded {
		res.ToTreasury = residue
		return res
	}

	nearest := 5
	for _, dk := range present {
		if dk.Class() < nearest {
			nearest = dk.Class()
		}
	}

	var members []DistantKindred
	var units int64
	for _, dk := range present {
		if dk.Class() == nearest {
			members = append(members, dk)
			units += int64(kindred[dk]) * weightOf(dk.Sex())
		}
	}

	total := rational.FromInt(units)
	for _, dk := range members {
		u := int64(kindred[dk]) * weightOf(dk.Sex())
		res.Shares[dk] = residue.Mul(rational.FromInt(u)).Div(total)
	}

	if len(members) > 1 {
		res.NeedsReview = true
		res.Note = "distribution among multiple distant kindred of one class needs review"
	}
	return res
}
