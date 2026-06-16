package solver

import "github.com/YASSERRMD/Faraid/internal/core/heir"

// AsabaType classifies how a residuary heir takes the residue.
type AsabaType uint8

const (
	// AsabaBiNafsihi is a residuary in their own right: a male agnate.
	AsabaBiNafsihi AsabaType = iota + 1
	// AsabaBiGhayrihi is a female made residuary by a male of her own level,
	// sharing two to one with him.
	AsabaBiGhayrihi
	// AsabaMaaGhayrihi is a full or consanguine sister made residuary by the
	// presence of a daughter or son's daughter.
	AsabaMaaGhayrihi
)

// String returns a label for the asaba type.
func (a AsabaType) String() string {
	switch a {
	case AsabaBiNafsihi:
		return "bi-nafsihi"
	case AsabaBiGhayrihi:
		return "bi-ghayrihi"
	case AsabaMaaGhayrihi:
		return "maa-ghayrihi"
	default:
		return "none"
	}
}

// asabaPriority lists the residuary heirs in their own right (asaba
// bi-nafsihi) from nearest to farthest. The nearest present agnate takes the
// residue and excludes those after it: the sonship line, then the fatherhood
// line, then the brotherhood line, then the uncle line.
//
// The grandfather is placed before the brothers, which matches the view that
// he excludes them. The competing view in which he shares with them (jadd wa
// ikhwa) is a school divergence handled in a dedicated later phase.
var asabaPriority = []heir.Relation{
	heir.Son,
	heir.SonsSon,
	heir.Father,
	heir.PaternalGrandfather,
	heir.FullBrother,
	heir.ConsanguineBrother,
	heir.FullBrothersSon,
	heir.ConsanguineBrothersSon,
	heir.FullPaternalUncle,
	heir.ConsanguinePaternalUncle,
	heir.FullPaternalUnclesSon,
	heir.ConsanguinePaternalUnclesSon,
}

// highestAgnate returns the nearest residuary-in-their-own-right heir present,
// if any.
func highestAgnate(h *heir.Heirs) (heir.Relation, bool) {
	for _, m := range asabaPriority {
		if h.Present(m) {
			return m, true
		}
	}
	return heir.RelationInvalid, false
}
