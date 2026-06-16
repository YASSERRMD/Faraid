package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestDistributeResidueSonsAndDaughters(t *testing.T) {
	// Residue of the whole shared by one son and one daughter, two to one.
	h := heir.New().Set(heir.Son, 1).Set(heir.Daughter, 1)
	out := DistributeResidue(rational.One(), ClassifyResiduary(h), h)
	if !out[heir.Son].Equal(rational.New(2, 3)) || !out[heir.Daughter].Equal(rational.New(1, 3)) {
		t.Errorf("son/daughter split = %v", out)
	}

	// Two sons and one daughter: units 2*2 + 1 = 5, sons 4/5, daughter 1/5.
	h = heir.New().Set(heir.Son, 2).Set(heir.Daughter, 1)
	out = DistributeResidue(rational.One(), ClassifyResiduary(h), h)
	if !out[heir.Son].Equal(rational.New(4, 5)) || !out[heir.Daughter].Equal(rational.New(1, 5)) {
		t.Errorf("two sons one daughter = %v", out)
	}
}

func TestDistributeResidueWithFixedShares(t *testing.T) {
	// Husband (1/4 with descendant), son, and daughter. Residue 3/4 split 2:1.
	h := heir.New().Set(heir.Husband, 1).Set(heir.Son, 1).Set(heir.Daughter, 1)
	residue := rational.New(3, 4)
	out := DistributeResidue(residue, ClassifyResiduary(h), h)
	if !out[heir.Son].Equal(rational.New(1, 2)) {
		t.Errorf("son = %s, want 1/2", out[heir.Son])
	}
	if !out[heir.Daughter].Equal(rational.New(1, 4)) {
		t.Errorf("daughter = %s, want 1/4", out[heir.Daughter])
	}
	// Husband 1/4 + son 1/2 + daughter 1/4 = 1.
	total := rational.New(1, 4).Add(out[heir.Son]).Add(out[heir.Daughter])
	if !total.IsWhole() {
		t.Errorf("total = %s, want 1", total)
	}
}

func TestDistributeResiduePureAgnate(t *testing.T) {
	// Father alone takes the whole residue.
	h := heir.New().Set(heir.Father, 1)
	out := DistributeResidue(rational.New(1, 3), ClassifyResiduary(h), h)
	if !out[heir.Father].Equal(rational.New(1, 3)) {
		t.Errorf("father = %s, want 1/3", out[heir.Father])
	}
}

func TestDistributeResidueBrotherSister(t *testing.T) {
	h := heir.New().Set(heir.FullBrother, 1).Set(heir.FullSister, 1)
	out := DistributeResidue(rational.One(), ClassifyResiduary(h), h)
	if !out[heir.FullBrother].Equal(rational.New(2, 3)) || !out[heir.FullSister].Equal(rational.New(1, 3)) {
		t.Errorf("brother/sister = %v", out)
	}
}

func TestDistributeResidueMaaGhayrihi(t *testing.T) {
	// One daughter takes 1/2 as fixed; the full sister takes the 1/2 residue.
	h := heir.New().Set(heir.Daughter, 1).Set(heir.FullSister, 2)
	out := DistributeResidue(rational.New(1, 2), ClassifyResiduary(h), h)
	if !out[heir.FullSister].Equal(rational.New(1, 2)) {
		t.Errorf("full sisters slot = %s, want 1/2", out[heir.FullSister])
	}
}

func TestDistributeResidueEmptyCases(t *testing.T) {
	h := heir.New().Set(heir.Son, 1)
	// Zero residue yields nothing.
	if out := DistributeResidue(rational.Zero(), ClassifyResiduary(h), h); len(out) != 0 {
		t.Errorf("zero residue should distribute nothing, got %v", out)
	}
	// No residuary heir yields nothing.
	none := heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1)
	if out := DistributeResidue(rational.One(), ClassifyResiduary(none), none); len(out) != 0 {
		t.Errorf("no residuary should distribute nothing, got %v", out)
	}
	// A residuary member absent from the heir set has no units.
	r := Residuary{Members: []ResiduaryMember{{Relation: heir.Son, Type: AsabaBiNafsihi}}}
	if out := DistributeResidue(rational.One(), r, heir.New()); len(out) != 0 {
		t.Errorf("absent members should distribute nothing, got %v", out)
	}
}
