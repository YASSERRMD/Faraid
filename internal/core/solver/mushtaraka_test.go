package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestMushtarakaShare(t *testing.T) {
	// Husband, mother, two uterine brothers, one full brother. Sharing view:
	// the third is split equally among the three maternal siblings.
	h := heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).
		Set(heir.UterineBrother, 2).Set(heir.FullBrother, 1)
	shares, ok := Mushtaraka(h, MushtarakaShare)
	if !ok {
		t.Fatal("expected mushtaraka")
	}
	if !shares[heir.Husband].Equal(rational.New(1, 2)) || !shares[heir.Mother].Equal(rational.New(1, 6)) {
		t.Errorf("husband/mother wrong: %v", shares)
	}
	if !shares[heir.UterineBrother].Equal(rational.New(2, 9)) {
		t.Errorf("uterine = %s, want 2/9", shares[heir.UterineBrother])
	}
	if !shares[heir.FullBrother].Equal(rational.New(1, 9)) {
		t.Errorf("full brother = %s, want 1/9", shares[heir.FullBrother])
	}
	if !sumShares(shares).IsWhole() {
		t.Errorf("shares sum to %s, want 1", sumShares(shares))
	}
}

func TestMushtarakaNoShare(t *testing.T) {
	// Same heirs, non-sharing view: the full brother takes nothing.
	h := heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).
		Set(heir.UterineBrother, 2).Set(heir.FullBrother, 1)
	shares, ok := Mushtaraka(h, MushtarakaNoShare)
	if !ok {
		t.Fatal("expected mushtaraka")
	}
	if !shares[heir.UterineBrother].Equal(rational.New(1, 3)) {
		t.Errorf("uterine = %s, want 1/3", shares[heir.UterineBrother])
	}
	if _, ok := shares[heir.FullBrother]; ok {
		t.Error("full brother should take nothing")
	}
	if !sumShares(shares).IsWhole() {
		t.Errorf("shares sum to %s, want 1", sumShares(shares))
	}
}

func TestMushtarakaMixedSexes(t *testing.T) {
	// One of each maternal sibling; sharing view divides equally, no two to one.
	h := heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).
		Set(heir.UterineBrother, 1).Set(heir.UterineSister, 1).
		Set(heir.FullBrother, 1).Set(heir.FullSister, 1)
	shares, _ := Mushtaraka(h, MushtarakaShare)
	for _, r := range []heir.Relation{heir.UterineBrother, heir.UterineSister, heir.FullBrother, heir.FullSister} {
		if !shares[r].Equal(rational.New(1, 12)) {
			t.Errorf("%s = %s, want 1/12 (equal per head)", r, shares[r])
		}
	}
	if !sumShares(shares).IsWhole() {
		t.Errorf("shares sum to %s, want 1", sumShares(shares))
	}
}

func TestNotMushtaraka(t *testing.T) {
	no := map[string]*heir.Heirs{
		"wife not husband": heir.New().Set(heir.Wife, 1).Set(heir.Mother, 1).Set(heir.UterineBrother, 2).Set(heir.FullBrother, 1),
		"one uterine":      heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.UterineBrother, 1).Set(heir.FullBrother, 1),
		"no full brother":  heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.UterineBrother, 2),
		"with son":         heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.UterineBrother, 2).Set(heir.FullBrother, 1).Set(heir.Son, 1),
	}
	for name, h := range no {
		if IsMushtaraka(h) {
			t.Errorf("%s should not be mushtaraka", name)
		}
		if _, ok := Mushtaraka(h, MushtarakaShare); ok {
			t.Errorf("%s should not return mushtaraka shares", name)
		}
	}
}

func TestDistributeEquallyEmpty(t *testing.T) {
	h := heir.New().Set(heir.UterineBrother, 2)
	// Zero portion yields nothing.
	if out := distributeEqually(rational.Zero(), []heir.Relation{heir.UterineBrother}, h); len(out) != 0 {
		t.Errorf("zero portion should distribute nothing, got %v", out)
	}
	// No members yields nothing.
	if out := distributeEqually(rational.New(1, 3), nil, h); len(out) != 0 {
		t.Errorf("no members should distribute nothing, got %v", out)
	}
}

func TestMushtarakaViewString(t *testing.T) {
	if MushtarakaShare.String() != "share" || MushtarakaNoShare.String() != "no-share" {
		t.Error("view labels wrong")
	}
}
