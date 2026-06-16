package solver

import (
	"math/big"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/estate"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// solveCase builds and solves a case with the given heirs and estate.
func solveCase(t *testing.T, sex heir.Sex, total int64, h *heir.Heirs, m Madhhab) Result {
	t.Helper()
	r, err := Solve(estate.Case{DeceasedSex: sex, Estate: estate.Estate{Total: total}, Heirs: h}, m)
	if err != nil {
		t.Fatalf("Solve: %v", err)
	}
	return r
}

// share returns the fraction for a relation, or zero if absent.
func share(r Result, rel heir.Relation) rational.Fraction {
	for _, s := range r.Shares {
		if s.Relation == rel {
			return s.Fraction
		}
	}
	return rational.Zero()
}

// assertTotalsToWhole checks the shares and any treasury residue sum to one.
func assertTotalsToWhole(t *testing.T, r Result) {
	t.Helper()
	total := r.Residue
	for _, s := range r.Shares {
		total = total.Add(s.Fraction)
	}
	if !total.IsWhole() {
		t.Errorf("shares plus residue = %s, want 1", total)
	}
}

func TestSolveFuruAndResiduary(t *testing.T) {
	// Husband, son, daughter: husband 1/4, son and daughter split the rest 2:1.
	r := solveCase(t, heir.Female, 0, heir.New().Set(heir.Husband, 1).Set(heir.Son, 1).Set(heir.Daughter, 1), Hanafi)
	if !share(r, heir.Husband).Equal(rational.New(1, 4)) {
		t.Errorf("husband = %s, want 1/4", share(r, heir.Husband))
	}
	if !share(r, heir.Son).Equal(rational.New(1, 2)) {
		t.Errorf("son = %s, want 1/2", share(r, heir.Son))
	}
	if !share(r, heir.Daughter).Equal(rational.New(1, 4)) {
		t.Errorf("daughter = %s, want 1/4", share(r, heir.Daughter))
	}
	assertTotalsToWhole(t, r)
}

func TestSolveFatherAndDaughter(t *testing.T) {
	// Daughter 1/2, father takes 1/6 plus the residue, reaching 1/2.
	r := solveCase(t, heir.Male, 0, heir.New().Set(heir.Father, 1).Set(heir.Daughter, 1), Shafii)
	if !share(r, heir.Daughter).Equal(rational.New(1, 2)) || !share(r, heir.Father).Equal(rational.New(1, 2)) {
		t.Errorf("father/daughter = %s / %s, want 1/2 each", share(r, heir.Father), share(r, heir.Daughter))
	}
}

func TestSolveAwl(t *testing.T) {
	// Husband and two full sisters: base 6 raised to 7.
	r := solveCase(t, heir.Female, 0, heir.New().Set(heir.Husband, 1).Set(heir.FullSister, 2), Maliki)
	if !r.Awl {
		t.Error("expected awl")
	}
	if !share(r, heir.Husband).Equal(rational.New(3, 7)) || !share(r, heir.FullSister).Equal(rational.New(4, 7)) {
		t.Errorf("awl shares = %s / %s, want 3/7 and 4/7", share(r, heir.Husband), share(r, heir.FullSister))
	}
	if r.Base.Cmp(big.NewInt(7)) != 0 {
		t.Errorf("base = %s, want 7", r.Base)
	}
}

func TestSolveRadd(t *testing.T) {
	// Mother and a daughter, no residuary: radd to 1/4 and 3/4.
	r := solveCase(t, heir.Male, 0, heir.New().Set(heir.Mother, 1).Set(heir.Daughter, 1), Hanbali)
	if !r.Radd {
		t.Error("expected radd")
	}
	if !share(r, heir.Mother).Equal(rational.New(1, 4)) || !share(r, heir.Daughter).Equal(rational.New(3, 4)) {
		t.Errorf("radd shares = %s / %s, want 1/4 and 3/4", share(r, heir.Mother), share(r, heir.Daughter))
	}
}

func TestSolveGharrawayn(t *testing.T) {
	r := solveCase(t, heir.Female, 0, heir.New().Set(heir.Husband, 1).Set(heir.Father, 1).Set(heir.Mother, 1), Hanafi)
	if r.SpecialCase != "gharrawayn" {
		t.Errorf("special case = %q, want gharrawayn", r.SpecialCase)
	}
	if !share(r, heir.Husband).Equal(rational.New(1, 2)) ||
		!share(r, heir.Mother).Equal(rational.New(1, 6)) ||
		!share(r, heir.Father).Equal(rational.New(1, 3)) {
		t.Errorf("gharrawayn shares wrong: %v", r.Shares)
	}
}

func TestSolveAkdariyyah(t *testing.T) {
	r := solveCase(t, heir.Female, 0, heir.New().
		Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 1), Shafii)
	if r.SpecialCase != "akdariyyah" {
		t.Errorf("special case = %q, want akdariyyah", r.SpecialCase)
	}
	if !share(r, heir.PaternalGrandfather).Equal(rational.New(8, 27)) || !share(r, heir.FullSister).Equal(rational.New(4, 27)) {
		t.Errorf("akdariyyah shares wrong: %v", r.Shares)
	}
	if r.Base.Cmp(big.NewInt(27)) != 0 {
		t.Errorf("base = %s, want 27", r.Base)
	}
}

func TestSolveMushtarakaSharingVsNot(t *testing.T) {
	h := heir.New().Set(heir.Husband, 1).Set(heir.Mother, 1).Set(heir.UterineBrother, 2).Set(heir.FullBrother, 1)
	// Maliki shares: the full brother gets 1/9.
	maliki := solveCase(t, heir.Female, 0, h, Maliki)
	if !share(maliki, heir.FullBrother).Equal(rational.New(1, 9)) {
		t.Errorf("Maliki full brother = %s, want 1/9", share(maliki, heir.FullBrother))
	}
	// Hanafi does not share: the full brother gets nothing.
	hanafi := solveCase(t, heir.Female, 0, h, Hanafi)
	if !share(hanafi, heir.FullBrother).IsZero() {
		t.Errorf("Hanafi full brother = %s, want 0", share(hanafi, heir.FullBrother))
	}
}

func TestSolveJaddWaIkhwa(t *testing.T) {
	h := heir.New().Set(heir.PaternalGrandfather, 1).Set(heir.FullBrother, 2)
	// Zayd (Maliki): grandfather takes one third, brothers two thirds.
	zayd := solveCase(t, heir.Male, 0, h, Maliki)
	if !share(zayd, heir.PaternalGrandfather).Equal(rational.New(1, 3)) {
		t.Errorf("Zayd grandfather = %s, want 1/3", share(zayd, heir.PaternalGrandfather))
	}
	// Abu Hanifa: grandfather excludes the brothers and takes everything.
	hanafi := solveCase(t, heir.Male, 0, h, Hanafi)
	if !share(hanafi, heir.PaternalGrandfather).IsWhole() {
		t.Errorf("Hanafi grandfather = %s, want the whole", share(hanafi, heir.PaternalGrandfather))
	}
}

func TestSolveBlocking(t *testing.T) {
	// A son excludes both the son's son and the full brother.
	r := solveCase(t, heir.Female, 0, heir.New().
		Set(heir.Husband, 1).Set(heir.Son, 1).Set(heir.SonsSon, 1).Set(heir.FullBrother, 1), Hanafi)
	if share(r, heir.FullBrother).Sign() != 0 || share(r, heir.SonsSon).Sign() != 0 {
		t.Error("son's son and full brother should be excluded")
	}
	if len(r.Excluded) != 2 {
		t.Errorf("expected two excluded heirs, got %v", r.Excluded)
	}
	// Excluded is sorted in canonical relation order.
	if r.Excluded[0] != heir.SonsSon || r.Excluded[1] != heir.FullBrother {
		t.Errorf("excluded order = %v, want [son's son, full brother]", r.Excluded)
	}
}

func TestSolveTashihAndAmounts(t *testing.T) {
	// Husband and two sons over an estate of 800: base raised to 8 by tashih,
	// husband 200, each son 300.
	r := solveCase(t, heir.Female, 800, heir.New().Set(heir.Husband, 1).Set(heir.Son, 2), Hanafi)
	if r.Base.Cmp(big.NewInt(8)) != 0 {
		t.Errorf("base = %s, want 8 after tashih", r.Base)
	}
	if !share(r, heir.Husband).Equal(rational.New(1, 4)) {
		t.Errorf("husband = %s, want 1/4", share(r, heir.Husband))
	}
	for _, s := range r.Shares {
		if s.Relation == heir.Husband && !s.Amount.Equal(rational.FromInt(200)) {
			t.Errorf("husband amount = %s, want 200", s.Amount)
		}
		if s.Relation == heir.Son && !s.Amount.Equal(rational.FromInt(600)) {
			t.Errorf("son slot amount = %s, want 600", s.Amount)
		}
	}
}

func TestSolveDeterminism(t *testing.T) {
	h := heir.New().Set(heir.Wife, 1).Set(heir.Mother, 1).Set(heir.Son, 2).Set(heir.Daughter, 3)
	a := solveCase(t, heir.Male, 1000, h, Hanbali)
	b := solveCase(t, heir.Male, 1000, h, Hanbali)
	if len(a.Shares) != len(b.Shares) || a.Base.Cmp(b.Base) != 0 {
		t.Fatal("results differ in shape")
	}
	for i := range a.Shares {
		if a.Shares[i].Relation != b.Shares[i].Relation || !a.Shares[i].Fraction.Equal(b.Shares[i].Fraction) {
			t.Errorf("nondeterministic at %d: %v vs %v", i, a.Shares[i], b.Shares[i])
		}
	}
	assertTotalsToWhole(t, a)
}

func TestSolveOnlySpouseToTreasury(t *testing.T) {
	// A husband alone keeps 1/2; the rest passes to the treasury.
	r := solveCase(t, heir.Female, 0, heir.New().Set(heir.Husband, 1), Hanafi)
	if !share(r, heir.Husband).Equal(rational.New(1, 2)) {
		t.Errorf("husband = %s, want 1/2", share(r, heir.Husband))
	}
	if !r.Residue.Equal(rational.New(1, 2)) {
		t.Errorf("treasury residue = %s, want 1/2", r.Residue)
	}
	assertTotalsToWhole(t, r)
}

func TestSolveGrandmotherPooling(t *testing.T) {
	// Both grandmothers share one sixth; the son takes the rest.
	r := solveCase(t, heir.Male, 0, heir.New().
		Set(heir.Son, 1).Set(heir.MaternalGrandmother, 1).Set(heir.PaternalGrandmother, 1), Hanafi)
	if !share(r, heir.MaternalGrandmother).Equal(rational.New(1, 12)) {
		t.Errorf("maternal grandmother = %s, want 1/12 (pooled)", share(r, heir.MaternalGrandmother))
	}
	if !share(r, heir.Son).Equal(rational.New(5, 6)) {
		t.Errorf("son = %s, want 5/6", share(r, heir.Son))
	}
	assertTotalsToWhole(t, r)
}

func TestSolveUterinePooling(t *testing.T) {
	// Mother and one uterine brother and sister: the uterine third is pooled,
	// then radd raises all to 1/3 each.
	r := solveCase(t, heir.Male, 0, heir.New().
		Set(heir.Mother, 1).Set(heir.UterineBrother, 1).Set(heir.UterineSister, 1), Hanafi)
	if !share(r, heir.Mother).Equal(rational.New(1, 3)) {
		t.Errorf("mother = %s, want 1/3 (pooling then radd)", share(r, heir.Mother))
	}
	assertTotalsToWhole(t, r)
}

func TestSolveMaaGhayrihi(t *testing.T) {
	// A daughter takes 1/2; the full sister becomes residuary with her.
	r := solveCase(t, heir.Male, 0, heir.New().Set(heir.Daughter, 1).Set(heir.FullSister, 1), Hanafi)
	if !share(r, heir.Daughter).Equal(rational.New(1, 2)) || !share(r, heir.FullSister).Equal(rational.New(1, 2)) {
		t.Errorf("maa-ghayrihi shares = %s / %s, want 1/2 each", share(r, heir.Daughter), share(r, heir.FullSister))
	}
}

func TestSolveJaddConsanguineAndOtherFurud(t *testing.T) {
	// Grandfather, mother, and one full brother (Zayd): mother 1/3, then the
	// grandfather and brother split 2/3 equally.
	r := solveCase(t, heir.Male, 0, heir.New().
		Set(heir.PaternalGrandfather, 1).Set(heir.Mother, 1).Set(heir.FullBrother, 1), Maliki)
	if !share(r, heir.Mother).Equal(rational.New(1, 3)) {
		t.Errorf("mother = %s, want 1/3", share(r, heir.Mother))
	}
	if !share(r, heir.PaternalGrandfather).Equal(rational.New(1, 3)) {
		t.Errorf("grandfather = %s, want 1/3", share(r, heir.PaternalGrandfather))
	}

	// Consanguine brother competes too.
	c := solveCase(t, heir.Male, 0, heir.New().
		Set(heir.PaternalGrandfather, 1).Set(heir.ConsanguineBrother, 1), Maliki)
	if !share(c, heir.PaternalGrandfather).Equal(rational.New(1, 2)) {
		t.Errorf("grandfather with consanguine brother = %s, want 1/2", share(c, heir.PaternalGrandfather))
	}
}

func TestSolveJaddNeedsReview(t *testing.T) {
	// Grandfather with a full sister and a consanguine brother is flagged.
	r := solveCase(t, heir.Male, 0, heir.New().
		Set(heir.PaternalGrandfather, 1).Set(heir.FullSister, 1).Set(heir.ConsanguineBrother, 1), Maliki)
	if !r.NeedsReview || len(r.ReviewNotes) == 0 {
		t.Error("grandfather with full sister and consanguine should need review")
	}
}

func TestSolveExactFixedShares(t *testing.T) {
	// Husband and one full sister consume the estate exactly, with no residue,
	// awl, or radd.
	r := solveCase(t, heir.Female, 0, heir.New().Set(heir.Husband, 1).Set(heir.FullSister, 1), Hanafi)
	if r.Awl || r.Radd || r.Residue.Sign() != 0 {
		t.Errorf("expected an exact split, got awl=%v radd=%v residue=%s", r.Awl, r.Radd, r.Residue)
	}
	if !share(r, heir.Husband).Equal(rational.New(1, 2)) || !share(r, heir.FullSister).Equal(rational.New(1, 2)) {
		t.Errorf("exact shares = %s / %s, want 1/2 each", share(r, heir.Husband), share(r, heir.FullSister))
	}
}

func TestSolveInvalidCase(t *testing.T) {
	_, err := Solve(estate.Case{DeceasedSex: heir.Male, Heirs: heir.New().Set(heir.Husband, 1)}, Hanafi)
	if err == nil {
		t.Error("a deceased male with a husband should be rejected")
	}
}

func TestSolveNegativeEstate(t *testing.T) {
	_, err := Solve(estate.Case{DeceasedSex: heir.Male, Estate: estate.Estate{Total: -100}, Heirs: heir.New().Set(heir.Son, 1)}, Hanafi)
	if err == nil {
		t.Error("a negative estate should be rejected")
	}
}
