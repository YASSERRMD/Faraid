package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

// asMembers turns a Residuary into a comparable map of relation to type.
func asMembers(r Residuary) map[heir.Relation]AsabaType {
	m := map[heir.Relation]AsabaType{}
	for _, mem := range r.Members {
		m[mem.Relation] = mem.Type
	}
	return m
}

func TestResiduaryBiNafsihiAndBiGhayrihi(t *testing.T) {
	// Son with a daughter: son in his own right, daughter through him.
	got := asMembers(ClassifyResiduary(heir.New().Set(heir.Son, 1).Set(heir.Daughter, 2)))
	if got[heir.Son] != AsabaBiNafsihi || got[heir.Daughter] != AsabaBiGhayrihi || len(got) != 2 {
		t.Errorf("son with daughters = %v", got)
	}
	// Son's son with son's daughter.
	got = asMembers(ClassifyResiduary(heir.New().Set(heir.SonsSon, 1).Set(heir.SonsDaughter, 1)))
	if got[heir.SonsSon] != AsabaBiNafsihi || got[heir.SonsDaughter] != AsabaBiGhayrihi {
		t.Errorf("son's son with son's daughter = %v", got)
	}
	// Full brother with full sister.
	got = asMembers(ClassifyResiduary(heir.New().Set(heir.FullBrother, 1).Set(heir.FullSister, 1)))
	if got[heir.FullBrother] != AsabaBiNafsihi || got[heir.FullSister] != AsabaBiGhayrihi {
		t.Errorf("full brother with full sister = %v", got)
	}
}

func TestResiduaryPriority(t *testing.T) {
	// Son outranks father, grandfather, and brothers; only the son (with his
	// daughter) takes the residue.
	got := asMembers(ClassifyResiduary(heir.New().
		Set(heir.Son, 1).Set(heir.Father, 1).Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 1)))
	if _, ok := got[heir.Son]; !ok || len(got) != 1 {
		t.Errorf("son should be the sole residuary, got %v", got)
	}
	// Father (no son) is the residuary; a daughter beside him is a fixed-share
	// heir, not bi-ghayrihi.
	got = asMembers(ClassifyResiduary(heir.New().Set(heir.Father, 1).Set(heir.Daughter, 1)))
	if got[heir.Father] != AsabaBiNafsihi || len(got) != 1 {
		t.Errorf("father should be the sole residuary, got %v", got)
	}
	// Grandfather outranks the full brother (Hanafi view; refined in Phase 20).
	got = asMembers(ClassifyResiduary(heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 1)))
	if _, ok := got[heir.PaternalGrandfather]; !ok || len(got) != 1 {
		t.Errorf("grandfather should outrank the full brother, got %v", got)
	}
	// Son's son outranks the full brother.
	got = asMembers(ClassifyResiduary(heir.New().Set(heir.SonsSon, 1).Set(heir.FullBrother, 1)))
	if _, ok := got[heir.SonsSon]; !ok || len(got) != 1 {
		t.Errorf("son's son should outrank the full brother, got %v", got)
	}
	// A lone uncle takes the residue.
	got = asMembers(ClassifyResiduary(heir.New().Set(heir.FullPaternalUncle, 1)))
	if got[heir.FullPaternalUncle] != AsabaBiNafsihi {
		t.Errorf("uncle should be residuary, got %v", got)
	}
}

func TestResiduaryMaaGhayrihi(t *testing.T) {
	// Full sister with a daughter becomes residuary alongside her.
	got := asMembers(ClassifyResiduary(heir.New().Set(heir.FullSister, 1).Set(heir.Daughter, 1)))
	if got[heir.FullSister] != AsabaMaaGhayrihi || len(got) != 1 {
		t.Errorf("full sister should be maa-ghayrihi, got %v", got)
	}
	// Consanguine sister with a son's daughter, no full sister.
	got = asMembers(ClassifyResiduary(heir.New().Set(heir.ConsanguineSister, 1).Set(heir.SonsDaughter, 1)))
	if got[heir.ConsanguineSister] != AsabaMaaGhayrihi {
		t.Errorf("consanguine sister should be maa-ghayrihi, got %v", got)
	}
}

func TestResiduaryNone(t *testing.T) {
	// Husband and mother only: no residuary heir, the residue returns by radd.
	if r := ClassifyResiduary(heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1)); r.Found() {
		t.Errorf("expected no residuary, got %v", r.Members)
	}
	// A lone sister with no descendant and no agnate is a fixed-share heir.
	if r := ClassifyResiduary(heir.New().Set(heir.FullSister, 1)); r.Found() {
		t.Errorf("expected no residuary, got %v", r.Members)
	}
}

func TestAsabaTypeString(t *testing.T) {
	cases := map[AsabaType]string{
		AsabaBiNafsihi:   "bi-nafsihi",
		AsabaBiGhayrihi:  "bi-ghayrihi",
		AsabaMaaGhayrihi: "maa-ghayrihi",
		AsabaType(0):     "none",
	}
	for a, want := range cases {
		if got := a.String(); got != want {
			t.Errorf("AsabaType(%d).String() = %q, want %q", a, got, want)
		}
	}
}
