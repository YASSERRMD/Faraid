package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func TestDhawuSingleKindredInherits(t *testing.T) {
	// Only a daughter's son, Hanafi: he takes the whole estate.
	res := DistributeDhawuArham(rational.Zero(), map[DistantKindred]int{DaughtersSon: 1}, Hanafi)
	if !res.Shares[DaughtersSon].IsWhole() {
		t.Errorf("daughter's son = %s, want the whole", res.Shares[DaughtersSon])
	}
	if res.NeedsReview {
		t.Error("a single distant kindred should not need review")
	}
}

func TestDhawuClassPriority(t *testing.T) {
	// A class 1 daughter's daughter excludes a class 3 sister's son.
	res := DistributeDhawuArham(rational.Zero(), map[DistantKindred]int{
		DaughtersDaughter: 1,
		SistersSon:        1,
	}, Hanbali)
	if !res.Shares[DaughtersDaughter].IsWhole() {
		t.Errorf("nearer class should take the whole, got %v", res.Shares)
	}
	if _, ok := res.Shares[SistersSon]; ok {
		t.Error("farther class should be excluded")
	}
}

func TestDhawuWithSpouse(t *testing.T) {
	// A wife takes 1/4; the daughter's son takes the rest.
	res := DistributeDhawuArham(rational.New(1, 4), map[DistantKindred]int{DaughtersSon: 1}, Hanafi)
	if !res.SpouseShare.Equal(rational.New(1, 4)) {
		t.Errorf("spouse = %s, want 1/4", res.SpouseShare)
	}
	if !res.Shares[DaughtersSon].Equal(rational.New(3, 4)) {
		t.Errorf("daughter's son = %s, want 3/4", res.Shares[DaughtersSon])
	}
}

func TestDhawuExcludedToTreasury(t *testing.T) {
	// Maliki excludes distant kindred: the residue passes to the treasury.
	res := DistributeDhawuArham(rational.New(1, 4), map[DistantKindred]int{DaughtersSon: 1}, Maliki)
	if len(res.Shares) != 0 {
		t.Errorf("distant kindred should not inherit under Maliki, got %v", res.Shares)
	}
	if !res.ToTreasury.Equal(rational.New(3, 4)) {
		t.Errorf("treasury = %s, want 3/4", res.ToTreasury)
	}
}

func TestDhawuMultiMemberNeedsReview(t *testing.T) {
	// A daughter's son and daughter's daughter, same class: best-effort 2:1 but
	// flagged for review.
	res := DistributeDhawuArham(rational.Zero(), map[DistantKindred]int{
		DaughtersSon:      1,
		DaughtersDaughter: 1,
	}, Hanafi)
	if !res.Shares[DaughtersSon].Equal(rational.New(2, 3)) || !res.Shares[DaughtersDaughter].Equal(rational.New(1, 3)) {
		t.Errorf("two-to-one split wrong: %v", res.Shares)
	}
	if !res.NeedsReview {
		t.Error("multiple kindred in a class should be flagged for review")
	}
}

func TestDhawuNoKindred(t *testing.T) {
	res := DistributeDhawuArham(rational.New(1, 2), map[DistantKindred]int{}, Hanafi)
	if !res.ToTreasury.Equal(rational.New(1, 2)) {
		t.Errorf("with no kindred the residue goes to the treasury, got %s", res.ToTreasury)
	}
}

func TestDistantKindredMetadata(t *testing.T) {
	if DaughtersSon.Class() != 1 || MaternalGrandfather.Class() != 2 || SistersSon.Class() != 3 || MaternalUncle.Class() != 4 {
		t.Error("class assignment wrong")
	}
	if DaughtersSon.String() != "daughter's son" {
		t.Errorf("label = %q", DaughtersSon.String())
	}
	if DistantKindred(99).String() != "unknown distant kindred" {
		t.Error("unknown label wrong")
	}
}
