package solver

import "github.com/YASSERRMD/Faraid/internal/core/heir"

// ResiduaryMember is a relation that shares in the residue, with its asaba type.
type ResiduaryMember struct {
	Relation heir.Relation
	Type     AsabaType
}

// Residuary describes the heirs who take the residue after the fixed shares.
type Residuary struct {
	Members []ResiduaryMember
}

// Found reports whether any heir takes the residue.
func (r Residuary) Found() bool { return len(r.Members) > 0 }

// biGhayrihiCounterpart maps a residuary male to the female of his own level who
// becomes residuary through him, sharing two to one. Only these four males have
// such a counterpart.
var biGhayrihiCounterpart = map[heir.Relation]heir.Relation{
	heir.Son:                heir.Daughter,
	heir.SonsSon:            heir.SonsDaughter,
	heir.FullBrother:        heir.FullSister,
	heir.ConsanguineBrother: heir.ConsanguineSister,
}

// ClassifyResiduary determines which surviving heirs take the residue and how.
// It must be given the set that survives blocking. The nearest male agnate
// takes the residue, joined by the female of his level as bi-ghayrihi when she
// is present. With no male agnate, a full or consanguine sister becomes
// residuary alongside a daughter or son's daughter (ma'a ghayrihi). When
// neither applies, there is no residuary heir.
func ClassifyResiduary(h *heir.Heirs) Residuary {
	if m, ok := highestAgnate(h); ok {
		members := []ResiduaryMember{{Relation: m, Type: AsabaBiNafsihi}}
		if f, ok := biGhayrihiCounterpart[m]; ok && h.Present(f) {
			members = append(members, ResiduaryMember{Relation: f, Type: AsabaBiGhayrihi})
		}
		return Residuary{Members: members}
	}

	if h.Present(heir.Daughter) || h.Present(heir.SonsDaughter) {
		switch {
		case h.Present(heir.FullSister):
			return Residuary{Members: []ResiduaryMember{{Relation: heir.FullSister, Type: AsabaMaaGhayrihi}}}
		case h.Present(heir.ConsanguineSister):
			return Residuary{Members: []ResiduaryMember{{Relation: heir.ConsanguineSister, Type: AsabaMaaGhayrihi}}}
		}
	}

	return Residuary{}
}
