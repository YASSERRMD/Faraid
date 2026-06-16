package rules

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

// anyHirman reports whether any total-exclusion rule for r fires under the
// given heirs. The full resolution engine is built in the next phase; this
// helper exercises the rule data directly.
func anyHirman(r heir.Relation, h *heir.Heirs) bool {
	c := ctx(h)
	for _, rule := range HirmanRules(r) {
		if rule.When(c) {
			return true
		}
	}
	return false
}

func TestLatticeValid(t *testing.T) {
	if err := ValidateLattice(); err != nil {
		t.Fatalf("blocking lattice should be valid and acyclic: %v", err)
	}
}

func TestDescendantBlocking(t *testing.T) {
	son := heir.New().Set(heir.Son, 1)
	for _, r := range []heir.Relation{
		heir.SonsSon, heir.SonsDaughter, heir.FullBrother, heir.FullSister,
		heir.ConsanguineBrother, heir.ConsanguineSister, heir.UterineBrother, heir.UterineSister,
	} {
		if !anyHirman(r, son) {
			t.Errorf("a son should exclude %s", r)
		}
	}
	// Any descendant excludes uterine siblings, even a daughter.
	if !anyHirman(heir.UterineBrother, heir.New().Set(heir.Daughter, 1)) {
		t.Error("a daughter should exclude uterine siblings")
	}
	// Two daughters exclude the son's daughter only without a son's son.
	if !anyHirman(heir.SonsDaughter, heir.New().Set(heir.Daughter, 2)) {
		t.Error("two daughters should exclude the son's daughter")
	}
	if anyHirman(heir.SonsDaughter, heir.New().Set(heir.Daughter, 2).Set(heir.SonsSon, 1)) {
		t.Error("a son's son should rescue the son's daughter from exclusion")
	}
}

func TestAscendantBlocking(t *testing.T) {
	father := heir.New().Set(heir.Father, 1)
	for _, r := range []heir.Relation{
		heir.PaternalGrandfather, heir.PaternalGrandmother,
		heir.FullBrother, heir.ConsanguineSister, heir.UterineBrother,
	} {
		if !anyHirman(r, father) {
			t.Errorf("the father should exclude %s", r)
		}
	}
	mother := heir.New().Set(heir.Mother, 1)
	if !anyHirman(heir.MaternalGrandmother, mother) || !anyHirman(heir.PaternalGrandmother, mother) {
		t.Error("the mother should exclude the grandmothers")
	}
	if !anyHirman(heir.UterineSister, heir.New().Set(heir.PaternalGrandfather, 1)) {
		t.Error("the grandfather should exclude uterine siblings")
	}
}

func TestSiblingBlocking(t *testing.T) {
	if !anyHirman(heir.ConsanguineBrother, heir.New().Set(heir.FullBrother, 1)) {
		t.Error("a full brother should exclude the consanguine brother")
	}
	// Two full sisters exclude the consanguine sister unless a consanguine brother saves her.
	if !anyHirman(heir.ConsanguineSister, heir.New().Set(heir.FullSister, 2)) {
		t.Error("two full sisters should exclude the consanguine sister")
	}
	if anyHirman(heir.ConsanguineSister, heir.New().Set(heir.FullSister, 2).Set(heir.ConsanguineBrother, 1)) {
		t.Error("a consanguine brother should rescue the consanguine sister")
	}
	// A full sister who is residuary with a daughter excludes the consanguine siblings.
	resid := heir.New().Set(heir.FullSister, 1).Set(heir.Daughter, 1)
	if !anyHirman(heir.ConsanguineSister, resid) || !anyHirman(heir.ConsanguineBrother, resid) {
		t.Error("a residuary full sister should exclude the consanguine siblings")
	}
}

func TestNuqsanRules(t *testing.T) {
	rules := NuqsanRules()
	if len(rules) == 0 {
		t.Fatal("expected nuqsan rules")
	}
	found := false
	for _, r := range rules {
		if r.Reduced == heir.Mother {
			found = true
		}
	}
	if !found {
		t.Error("expected a nuqsan rule reducing the mother")
	}
}

func TestValidateLatticeCatchesErrors(t *testing.T) {
	goodBlock := func() BlockRule {
		return BlockRule{
			Blocked:   heir.SonsSon,
			Blockers:  []heir.Relation{heir.Son},
			Condition: "c",
			Reference: "r",
			When:      func(Context) bool { return true },
		}
	}
	goodNuqsan := NuqsanRule{
		Reduced:   heir.Mother,
		Reducers:  []heir.Relation{heir.Son},
		Condition: "c",
		Reference: "r",
	}

	mutate := func(f func(b *BlockRule)) map[heir.Relation][]BlockRule {
		b := goodBlock()
		f(&b)
		return map[heir.Relation][]BlockRule{b.Blocked: {b}}
	}

	hirmanCases := map[string]map[heir.Relation][]BlockRule{
		"unknown blocked": {heir.Relation(999): {goodBlock()}},
		"mismatched key":  {heir.Son: {goodBlock()}},
		"nil when":        mutate(func(b *BlockRule) { b.When = nil }),
		"empty condition": mutate(func(b *BlockRule) { b.Condition = "" }),
		"empty reference": mutate(func(b *BlockRule) { b.Reference = "" }),
		"no blockers":     mutate(func(b *BlockRule) { b.Blockers = nil }),
		"unknown blocker": mutate(func(b *BlockRule) { b.Blockers = []heir.Relation{heir.Relation(999)} }),
	}
	for name, h := range hirmanCases {
		if err := validateLattice(h, nil); err == nil {
			t.Errorf("hirman %s: expected error", name)
		}
	}

	mutateN := func(f func(n *NuqsanRule)) []NuqsanRule {
		n := goodNuqsan
		f(&n)
		return []NuqsanRule{n}
	}
	nuqsanCases := map[string][]NuqsanRule{
		"unknown reduced": mutateN(func(n *NuqsanRule) { n.Reduced = heir.Relation(999) }),
		"empty condition": mutateN(func(n *NuqsanRule) { n.Condition = "" }),
		"empty reference": mutateN(func(n *NuqsanRule) { n.Reference = "" }),
		"no reducers":     mutateN(func(n *NuqsanRule) { n.Reducers = nil }),
		"unknown reducer": mutateN(func(n *NuqsanRule) { n.Reducers = []heir.Relation{heir.Relation(999)} }),
	}
	for name, n := range nuqsanCases {
		if err := validateLattice(nil, n); err == nil {
			t.Errorf("nuqsan %s: expected error", name)
		}
	}

	// A cyclic hirman set is rejected.
	cyclic := map[heir.Relation][]BlockRule{
		heir.Son:    {{Blocked: heir.Son, Blockers: []heir.Relation{heir.Father}, Condition: "c", Reference: "r", When: func(Context) bool { return true }}},
		heir.Father: {{Blocked: heir.Father, Blockers: []heir.Relation{heir.Son}, Condition: "c", Reference: "r", When: func(Context) bool { return true }}},
	}
	if err := validateLattice(cyclic, nil); err == nil {
		t.Error("expected a cycle error")
	}
}

func TestDetectCycle(t *testing.T) {
	acyclic := map[heir.Relation][]heir.Relation{heir.Son: {heir.Father}}
	if err := detectCycle(acyclic); err != nil {
		t.Errorf("acyclic graph reported a cycle: %v", err)
	}
	cyclic := map[heir.Relation][]heir.Relation{
		heir.Son:    {heir.Father},
		heir.Father: {heir.Son},
	}
	if err := detectCycle(cyclic); err == nil {
		t.Error("expected a cycle to be detected")
	}
}
