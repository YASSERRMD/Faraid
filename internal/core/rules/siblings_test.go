package rules

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestUterineShares(t *testing.T) {
	mustShare(t, heir.UterineBrother, heir.New().Set(heir.UterineBrother, 1), rational.New(1, 6))
	mustShare(t, heir.UterineSister, heir.New().Set(heir.UterineSister, 1), rational.New(1, 6))
	mustShare(t, heir.UterineBrother, heir.New().Set(heir.UterineBrother, 2), rational.New(1, 3))
	// One brother and one sister: the group takes one third (pooled at assembly).
	mixed := heir.New().Set(heir.UterineBrother, 1).Set(heir.UterineSister, 1)
	mustShare(t, heir.UterineBrother, mixed, rational.New(1, 3))
	mustShare(t, heir.UterineSister, mixed, rational.New(1, 3))
	// Blocked by any descendant (even a daughter), by the father, by the grandfather.
	mustNoShare(t, heir.UterineBrother, heir.New().Set(heir.UterineBrother, 1).Set(heir.Daughter, 1))
	mustNoShare(t, heir.UterineBrother, heir.New().Set(heir.UterineBrother, 1).Set(heir.Father, 1))
	mustNoShare(t, heir.UterineBrother, heir.New().Set(heir.UterineBrother, 1).Set(heir.PaternalGrandfather, 1))
}

func TestFullSisterShares(t *testing.T) {
	mustShare(t, heir.FullSister, heir.New().Set(heir.FullSister, 1), rational.New(1, 2))
	mustShare(t, heir.FullSister, heir.New().Set(heir.FullSister, 2), rational.New(2, 3))
	// Full brother makes her residuary; descendant or father blocks her.
	mustNoShare(t, heir.FullSister, heir.New().Set(heir.FullSister, 1).Set(heir.FullBrother, 1))
	mustNoShare(t, heir.FullSister, heir.New().Set(heir.FullSister, 1).Set(heir.Son, 1))
	mustNoShare(t, heir.FullSister, heir.New().Set(heir.FullSister, 1).Set(heir.Father, 1))
	// With a daughter she is residuary (ma'a ghayrihi), no fixed share.
	mustNoShare(t, heir.FullSister, heir.New().Set(heir.FullSister, 1).Set(heir.Daughter, 1))
}

func TestConsanguineSisterShares(t *testing.T) {
	mustShare(t, heir.ConsanguineSister, heir.New().Set(heir.ConsanguineSister, 1), rational.New(1, 2))
	mustShare(t, heir.ConsanguineSister, heir.New().Set(heir.ConsanguineSister, 2), rational.New(2, 3))
	// One full sister (1/2) plus consanguine sister: she completes with 1/6.
	completion := heir.New().Set(heir.FullSister, 1).Set(heir.ConsanguineSister, 1)
	mustShare(t, heir.FullSister, completion, rational.New(1, 2))
	mustShare(t, heir.ConsanguineSister, completion, rational.New(1, 6))
	// Two full sisters consume 2/3; consanguine sister blocked with no brother.
	mustNoShare(t, heir.ConsanguineSister, heir.New().Set(heir.FullSister, 2).Set(heir.ConsanguineSister, 1))
	// A consanguine brother or full brother removes her fixed share.
	mustNoShare(t, heir.ConsanguineSister, heir.New().Set(heir.ConsanguineSister, 1).Set(heir.ConsanguineBrother, 1))
	mustNoShare(t, heir.ConsanguineSister, heir.New().Set(heir.ConsanguineSister, 1).Set(heir.FullBrother, 1))
}
