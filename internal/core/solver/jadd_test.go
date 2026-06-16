package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestJaddMuqasamaWithFullBrothers(t *testing.T) {
	// Grandfather and one full brother share equally (muqasama).
	res := GrandfatherWithSiblings(rational.One(), heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 1), JaddZayd)
	if res.Method != "muqasama" || !res.GrandfatherShare.Equal(rational.New(1, 2)) {
		t.Errorf("gf = %s via %s, want 1/2 muqasama", res.GrandfatherShare, res.Method)
	}
	if !res.SiblingShares[heir.FullBrother].Equal(rational.New(1, 2)) {
		t.Errorf("brother = %s, want 1/2", res.SiblingShares[heir.FullBrother])
	}

	// Two full brothers: muqasama equals one third, grandfather takes one third.
	res = GrandfatherWithSiblings(rational.One(), heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 2), JaddZayd)
	if !res.GrandfatherShare.Equal(rational.New(1, 3)) {
		t.Errorf("gf with two brothers = %s, want 1/3", res.GrandfatherShare)
	}

	// Three full brothers: one third beats muqasama (2/8).
	res = GrandfatherWithSiblings(rational.One(), heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 3), JaddZayd)
	if res.Method != "one third of the remainder" || !res.GrandfatherShare.Equal(rational.New(1, 3)) {
		t.Errorf("gf with three brothers = %s via %s, want 1/3 third", res.GrandfatherShare, res.Method)
	}
}

func TestJaddMuqasamaWithFullSister(t *testing.T) {
	// Grandfather counts as a brother, so two to one against one sister.
	res := GrandfatherWithSiblings(rational.One(), heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 1), JaddZayd)
	if !res.GrandfatherShare.Equal(rational.New(2, 3)) {
		t.Errorf("gf = %s, want 2/3", res.GrandfatherShare)
	}
	if !res.SiblingShares[heir.FullSister].Equal(rational.New(1, 3)) {
		t.Errorf("sister = %s, want 1/3", res.SiblingShares[heir.FullSister])
	}
}

func TestJaddWithFixedHeirs(t *testing.T) {
	// After a mother takes 1/6, the grandfather and one full brother split the
	// remaining 5/6 equally: 5/12 each.
	h := heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 1)
	res := GrandfatherWithSiblings(rational.New(5, 6), h, JaddZayd)
	if !res.GrandfatherShare.Equal(rational.New(5, 12)) {
		t.Errorf("gf = %s, want 5/12", res.GrandfatherShare)
	}
	if !res.SiblingShares[heir.FullBrother].Equal(rational.New(5, 12)) {
		t.Errorf("brother = %s, want 5/12", res.SiblingShares[heir.FullBrother])
	}
}

func TestJaddSixthMinimum(t *testing.T) {
	// Heavy fixed shares leave 1/3; with four brothers both muqasama (1/15) and
	// a third of the remainder (1/9) fall below the 1/6 guarantee.
	h := heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 4)
	res := GrandfatherWithSiblings(rational.New(1, 3), h, JaddZayd)
	if res.Method != "one sixth of the estate" || !res.GrandfatherShare.Equal(rational.New(1, 6)) {
		t.Errorf("gf = %s via %s, want 1/6 sixth", res.GrandfatherShare, res.Method)
	}
}

func TestJaddConsanguineOnly(t *testing.T) {
	res := GrandfatherWithSiblings(rational.One(), heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.ConsanguineBrother, 1), JaddZayd)
	if !res.GrandfatherShare.Equal(rational.New(1, 2)) || !res.SiblingShares[heir.ConsanguineBrother].Equal(rational.New(1, 2)) {
		t.Errorf("consanguine muqasama wrong: gf %s brother %s", res.GrandfatherShare, res.SiblingShares[heir.ConsanguineBrother])
	}
}

func TestJaddAbuHanifaExcludes(t *testing.T) {
	h := heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 2)
	res := GrandfatherWithSiblings(rational.One(), h, JaddAbuHanifa)
	if !res.SiblingsExcluded || !res.GrandfatherShare.IsWhole() {
		t.Errorf("abu hanifa: gf = %s excluded = %v, want whole and excluded", res.GrandfatherShare, res.SiblingsExcluded)
	}
	if len(res.SiblingShares) != 0 {
		t.Errorf("siblings should get nothing, got %v", res.SiblingShares)
	}
}

func TestJaddMixedNeedsReview(t *testing.T) {
	// Full sister plus consanguine sibling with no full brother is flagged.
	h := heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 1).Set(heir.ConsanguineBrother, 1)
	res := GrandfatherWithSiblings(rational.One(), h, JaddZayd)
	if !res.NeedsReview || res.ReviewNote == "" {
		t.Error("mixed full sister and consanguine should need review")
	}
}

func TestJaddFullBrotherExcludesConsanguine(t *testing.T) {
	// A full brother present: consanguine siblings counted against the
	// grandfather but excluded from the share.
	h := heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 1).Set(heir.ConsanguineBrother, 1)
	res := GrandfatherWithSiblings(rational.One(), h, JaddZayd)
	if res.NeedsReview {
		t.Error("full brother present should not need review")
	}
	// units = 2 (gf-as-brother is 2) vs 2 (full) + 2 (cons) = 4 sibling units; muqasama 2/6 = 1/3.
	if !res.GrandfatherShare.Equal(rational.New(1, 3)) {
		t.Errorf("gf = %s, want 1/3", res.GrandfatherShare)
	}
	if _, ok := res.SiblingShares[heir.ConsanguineBrother]; ok {
		t.Error("consanguine brother should be excluded from the share")
	}
	if !res.SiblingShares[heir.FullBrother].Equal(rational.New(2, 3)) {
		t.Errorf("full brother = %s, want 2/3", res.SiblingShares[heir.FullBrother])
	}
}

func TestJaddSixthExceedsAvailable(t *testing.T) {
	// When fixed shares leave less than 1/6, the grandfather's 1/6 guarantee
	// exceeds the available portion (an over-subscribed case for awl), and the
	// siblings are left with nothing here.
	h := heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 1)
	res := GrandfatherWithSiblings(rational.New(1, 12), h, JaddZayd)
	if !res.GrandfatherShare.Equal(rational.New(1, 6)) {
		t.Errorf("gf = %s, want the 1/6 guarantee", res.GrandfatherShare)
	}
	if len(res.SiblingShares) != 0 {
		t.Errorf("siblings should get nothing, got %v", res.SiblingShares)
	}
}

func TestJaddNoSiblings(t *testing.T) {
	res := GrandfatherWithSiblings(rational.New(1, 2), heir.New().Set(heir.PaternalGrandfather, 1), JaddZayd)
	if !res.GrandfatherShare.Equal(rational.New(1, 2)) {
		t.Errorf("gf alone = %s, want the whole available 1/2", res.GrandfatherShare)
	}
}

func TestJaddViewString(t *testing.T) {
	if JaddZayd.String() != "zayd" || JaddAbuHanifa.String() != "abu-hanifa" {
		t.Error("view labels wrong")
	}
}
