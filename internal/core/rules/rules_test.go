package rules

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func ctx(h *heir.Heirs) Context { return Context{Heirs: h} }

func mustShare(t *testing.T, r heir.Relation, h *heir.Heirs, want rational.Fraction) {
	t.Helper()
	got, _, ok := FixedShare(r, ctx(h))
	if !ok {
		t.Fatalf("%s: expected a fixed share, got none", r)
	}
	if !got.Equal(want) {
		t.Errorf("%s share = %s, want %s", r, got, want)
	}
}

func TestHusbandShare(t *testing.T) {
	// No descendant: 1/2.
	mustShare(t, heir.Husband, heir.New().Set(heir.Husband, 1), rational.New(1, 2))
	// With a son: 1/4.
	mustShare(t, heir.Husband, heir.New().Set(heir.Husband, 1).Set(heir.Son, 1), rational.New(1, 4))
	// With a daughter (still a descendant): 1/4.
	mustShare(t, heir.Husband, heir.New().Set(heir.Husband, 1).Set(heir.Daughter, 1), rational.New(1, 4))
}

func TestWifeShare(t *testing.T) {
	// No descendant: 1/4 for the wife slot.
	mustShare(t, heir.Wife, heir.New().Set(heir.Wife, 1), rational.New(1, 4))
	// With a son's daughter (an inheriting descendant): 1/8.
	mustShare(t, heir.Wife, heir.New().Set(heir.Wife, 2).Set(heir.SonsDaughter, 1), rational.New(1, 8))
}

func TestNoFixedShareForResiduaryOnly(t *testing.T) {
	// A son has no fixed share; he is a residuary heir.
	if _, _, ok := FixedShare(heir.Son, ctx(heir.New().Set(heir.Son, 1))); ok {
		t.Error("son should have no fixed share")
	}
}

func TestRules(t *testing.T) {
	if len(Rules(heir.Husband)) != 2 {
		t.Errorf("husband should have 2 rules, got %d", len(Rules(heir.Husband)))
	}
	if len(Rules(heir.Son)) != 0 {
		t.Errorf("son should have 0 rules, got %d", len(Rules(heir.Son)))
	}
}

func TestContextHelpers(t *testing.T) {
	c := ctx(heir.New().Set(heir.SonsSon, 1))
	if !c.HasInheritingDescendant() {
		t.Error("son's son is an inheriting descendant")
	}
	if c.Count(heir.SonsSon) != 1 || !c.Present(heir.SonsSon) {
		t.Error("context count/present wrong")
	}
	empty := ctx(heir.New())
	if empty.HasInheritingDescendant() {
		t.Error("empty heirs have no descendant")
	}
}

func TestValidateTableOK(t *testing.T) {
	if err := ValidateTable(); err != nil {
		t.Errorf("package table should be valid: %v", err)
	}
}

func TestValidateTableCatchesErrors(t *testing.T) {
	good := FixedShareRule{
		Share:     rational.New(1, 2),
		Condition: "x",
		Reference: "y",
		When:      func(Context) bool { return true },
	}
	cases := map[string]map[heir.Relation][]FixedShareRule{
		"unknown relation": {heir.Relation(999): {good}},
		"no rules":         {heir.Husband: {}},
		"nil predicate":    {heir.Husband: {{Share: rational.New(1, 2), Condition: "x", Reference: "y"}}},
		"empty condition":  {heir.Husband: {{Share: rational.New(1, 2), Reference: "y", When: func(Context) bool { return true }}}},
		"empty reference":  {heir.Husband: {{Share: rational.New(1, 2), Condition: "x", When: func(Context) bool { return true }}}},
		"non-canonical":    {heir.Husband: {{Share: rational.New(3, 7), Condition: "x", Reference: "y", When: func(Context) bool { return true }}}},
	}
	for name, table := range cases {
		if err := validateTable(table); err == nil {
			t.Errorf("%s: expected validation error", name)
		}
	}
}
