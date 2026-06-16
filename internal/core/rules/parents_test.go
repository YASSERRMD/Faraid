package rules

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestMotherShares(t *testing.T) {
	// No descendant, fewer than two siblings: one third.
	mustShare(t, heir.Mother, heir.New().Set(heir.Mother, 1), rational.New(1, 3))
	// One sibling does not reduce the mother.
	mustShare(t, heir.Mother, heir.New().Set(heir.Mother, 1).Set(heir.FullBrother, 1), rational.New(1, 3))
	// An inheriting descendant reduces her to one sixth.
	mustShare(t, heir.Mother, heir.New().Set(heir.Mother, 1).Set(heir.Son, 1), rational.New(1, 6))
	// Two siblings reduce her to one sixth even with no descendant.
	mustShare(t, heir.Mother, heir.New().Set(heir.Mother, 1).Set(heir.FullBrother, 2), rational.New(1, 6))
	// Two siblings of mixed kind also reduce her.
	mustShare(t, heir.Mother, heir.New().Set(heir.Mother, 1).Set(heir.FullSister, 1).Set(heir.UterineBrother, 1), rational.New(1, 6))
}

func TestFatherShares(t *testing.T) {
	// With an inheriting descendant: fixed one sixth.
	mustShare(t, heir.Father, heir.New().Set(heir.Father, 1).Set(heir.Son, 1), rational.New(1, 6))
	mustShare(t, heir.Father, heir.New().Set(heir.Father, 1).Set(heir.Daughter, 1), rational.New(1, 6))
	// With no descendant: pure residuary, no fixed share.
	mustNoShare(t, heir.Father, heir.New().Set(heir.Father, 1))
}

func TestGrandmotherShares(t *testing.T) {
	mustShare(t, heir.MaternalGrandmother, heir.New().Set(heir.MaternalGrandmother, 1), rational.New(1, 6))
	mustShare(t, heir.PaternalGrandmother, heir.New().Set(heir.PaternalGrandmother, 1), rational.New(1, 6))
}

func TestGrandfatherShares(t *testing.T) {
	mustShare(t, heir.PaternalGrandfather, heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.Son, 1), rational.New(1, 6))
	mustNoShare(t, heir.PaternalGrandfather, heir.New().Set(heir.PaternalGrandfather, 1))
}

func TestIsGharrawayn(t *testing.T) {
	yes := []*heir.Heirs{
		heir.New().Set(heir.Husband, 1).Set(heir.Father, 1).Set(heir.Mother, 1),
		heir.New().Set(heir.Wife, 1).Set(heir.Father, 1).Set(heir.Mother, 1),
	}
	for i, h := range yes {
		if !IsGharrawayn(ctx(h)) {
			t.Errorf("case %d should be gharrawayn", i)
		}
	}
	no := map[string]*heir.Heirs{
		"with descendant": heir.New().Set(heir.Husband, 1).Set(heir.Father, 1).Set(heir.Mother, 1).Set(heir.Son, 1),
		"with sibling":    heir.New().Set(heir.Husband, 1).Set(heir.Father, 1).Set(heir.Mother, 1).Set(heir.FullBrother, 1),
		"no spouse":       heir.New().Set(heir.Father, 1).Set(heir.Mother, 1),
		"no father":       heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1),
		"no mother":       heir.New().Set(heir.Husband, 1).Set(heir.Father, 1),
	}
	for name, h := range no {
		if IsGharrawayn(ctx(h)) {
			t.Errorf("%s should not be gharrawayn", name)
		}
	}
}
