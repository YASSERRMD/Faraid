package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

func TestResolveBlockingBySon(t *testing.T) {
	h := heir.New().
		Set(heir.Husband, 1).
		Set(heir.Mother, 1).
		Set(heir.Father, 1).
		Set(heir.Son, 1).
		Set(heir.SonsSon, 1).
		Set(heir.FullBrother, 1).
		Set(heir.UterineSister, 1)
	res := ResolveBlocking(h)

	for _, r := range []heir.Relation{heir.SonsSon, heir.FullBrother, heir.UterineSister} {
		if !res.IsExcluded(r) {
			t.Errorf("%s should be excluded by the son", r)
		}
		if res.Surviving.Present(r) {
			t.Errorf("%s should not be in the surviving set", r)
		}
	}
	for _, r := range []heir.Relation{heir.Husband, heir.Mother, heir.Father, heir.Son} {
		if res.IsExcluded(r) {
			t.Errorf("%s should not be excluded", r)
		}
		if !res.Surviving.Present(r) {
			t.Errorf("%s should survive", r)
		}
	}
	// The recorded exclusion carries its blockers and reference.
	ex := res.Excluded[heir.FullBrother]
	if len(ex.By) == 0 || ex.Reference == "" || ex.Reason == "" {
		t.Errorf("exclusion record incomplete: %+v", ex)
	}
}

func TestResolveBlockingByFatherAndMother(t *testing.T) {
	h := heir.New().
		Set(heir.Father, 1).
		Set(heir.Mother, 1).
		Set(heir.PaternalGrandfather, 1).
		Set(heir.PaternalGrandmother, 1).
		Set(heir.MaternalGrandmother, 1).
		Set(heir.FullSister, 1)
	res := ResolveBlocking(h)

	for _, r := range []heir.Relation{
		heir.PaternalGrandfather, heir.PaternalGrandmother, heir.MaternalGrandmother, heir.FullSister,
	} {
		if !res.IsExcluded(r) {
			t.Errorf("%s should be excluded", r)
		}
	}
	if res.Surviving.Present(heir.PaternalGrandfather) {
		t.Error("grandfather should be excluded by the father")
	}
}

func TestReductionForWithoutNuqsanRule(t *testing.T) {
	// A relation with no documented reduction yields a bare record.
	r := reductionFor(heir.Son)
	if r.Relation != heir.Son || r.Reference != "" || r.Reason != "" {
		t.Errorf("expected a bare reduction for a relation without a nuqsan rule, got %+v", r)
	}
}

func TestResolveBlockingCountsPreserved(t *testing.T) {
	h := heir.New().Set(heir.Daughter, 3).Set(heir.Mother, 1)
	res := ResolveBlocking(h)
	if res.Surviving.Count(heir.Daughter) != 3 {
		t.Errorf("surviving daughters = %d, want 3", res.Surviving.Count(heir.Daughter))
	}
}

func TestReductions(t *testing.T) {
	// Husband reduced by a descendant.
	rs := Reductions(heir.New().Set(heir.Husband, 1).Set(heir.Son, 1))
	if len(rs) != 1 || rs[0].Relation != heir.Husband || rs[0].Reference == "" {
		t.Errorf("expected husband reduction, got %+v", rs)
	}
	// Mother reduced by a descendant.
	rs = Reductions(heir.New().Set(heir.Mother, 1).Set(heir.Daughter, 1))
	if len(rs) != 1 || rs[0].Relation != heir.Mother {
		t.Errorf("expected mother reduction, got %+v", rs)
	}
	// Mother reduced by two siblings, with a wife also reduced is not the case here.
	rs = Reductions(heir.New().Set(heir.Mother, 1).Set(heir.FullBrother, 2))
	if len(rs) != 1 || rs[0].Relation != heir.Mother {
		t.Errorf("expected mother reduction by siblings, got %+v", rs)
	}
	// Wife and mother both reduced by a descendant.
	rs = Reductions(heir.New().Set(heir.Wife, 1).Set(heir.Mother, 1).Set(heir.Son, 1))
	if len(rs) != 2 {
		t.Errorf("expected two reductions, got %+v", rs)
	}
	// No reductions when alone.
	if rs := Reductions(heir.New().Set(heir.Mother, 1)); len(rs) != 0 {
		t.Errorf("expected no reductions, got %+v", rs)
	}
	if rs := Reductions(heir.New().Set(heir.Husband, 1)); len(rs) != 0 {
		t.Errorf("expected no reductions, got %+v", rs)
	}
}
