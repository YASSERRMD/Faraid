package solver

import (
	"math/rand"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/estate"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// propertySeed is recorded so any failure is reproducible.
const propertySeed = 20260616

// randomCase builds a valid random case: a deceased sex, a random subset of
// heirs with counts within their limits, the matching spouse only, and a random
// estate.
func randomCase(rng *rand.Rand) estate.Case {
	sex := heir.Male
	if rng.Intn(2) == 0 {
		sex = heir.Female
	}
	h := heir.New()
	for _, r := range heir.AllRelations() {
		if sex == heir.Male && r == heir.Husband {
			continue
		}
		if sex == heir.Female && r == heir.Wife {
			continue
		}
		if rng.Intn(3) != 0 {
			continue
		}
		count := 1 + rng.Intn(3)
		if max := r.MaxCount(); max > 0 {
			count = 1 + rng.Intn(max)
		}
		h.Set(r, count)
	}
	return estate.Case{
		DeceasedSex: sex,
		Estate:      estate.Estate{Total: int64(rng.Intn(100000))},
		Heirs:       h,
	}
}

func sameResult(a, b Result) bool {
	if len(a.Shares) != len(b.Shares) || a.Base.Cmp(b.Base) != 0 {
		return false
	}
	for i := range a.Shares {
		if a.Shares[i].Relation != b.Shares[i].Relation || !a.Shares[i].Fraction.Equal(b.Shares[i].Fraction) {
			return false
		}
	}
	return true
}

func TestPropertyInvariants(t *testing.T) {
	rng := rand.New(rand.NewSource(propertySeed))
	t.Logf("property seed = %d", propertySeed)

	const cases = 500
	resolved, review := 0, 0

	for i := 0; i < cases; i++ {
		c := randomCase(rng)
		if c.Validate() != nil {
			continue
		}
		for _, m := range Madhahib() {
			r, err := Solve(c, m)
			if err != nil {
				t.Fatalf("solve error: %v", err)
			}

			// No negative share or amount.
			for _, s := range r.Shares {
				if s.Fraction.IsNegative() {
					t.Fatalf("negative share %s for %s", s.Fraction, s.Relation)
				}
				if s.Amount.IsNegative() {
					t.Fatalf("negative amount %s for %s", s.Amount, s.Relation)
				}
			}

			// The base is always positive.
			if r.Base.Sign() <= 0 {
				t.Fatalf("non-positive base %s", r.Base)
			}

			// Excluded heirs receive nothing.
			for _, ex := range r.Excluded {
				if share(r, ex).Sign() != 0 {
					t.Fatalf("excluded heir %s has a nonzero share", ex)
				}
			}

			// Determinism: the same input yields the same output.
			if again, _ := Solve(c, m); !sameResult(r, again) {
				t.Fatalf("nondeterministic result for %v under %s", c.Heirs.Relations(), m.Name)
			}

			if r.NeedsReview {
				review++
				continue
			}
			resolved++

			// A fully resolved result consumes exactly the whole estate.
			total := r.Residue
			for _, s := range r.Shares {
				total = total.Add(s.Fraction)
			}
			if !total.Equal(rational.One()) {
				t.Fatalf("resolved shares sum to %s, want 1, for %v under %s", total, c.Heirs.Relations(), m.Name)
			}
		}
	}

	t.Logf("resolved=%d needs-review=%d", resolved, review)
	if resolved == 0 {
		t.Fatal("expected some fully resolved cases")
	}
}

func TestEnforceSumInvariant(t *testing.T) {
	// Shares that fall short of the whole are flagged for review.
	short := &Result{Shares: []HeirShare{{Fraction: rational.New(1, 2)}}}
	enforceSumInvariant(short)
	if !short.NeedsReview || len(short.ReviewNotes) == 0 {
		t.Error("short shares should be flagged for review")
	}
	// Shares that consume the whole are not flagged.
	whole := &Result{Shares: []HeirShare{{Fraction: rational.New(1, 2)}, {Fraction: rational.New(1, 2)}}}
	enforceSumInvariant(whole)
	if whole.NeedsReview {
		t.Error("complete shares should not be flagged")
	}
	// A result already under review is not flagged twice.
	already := &Result{Shares: []HeirShare{{Fraction: rational.New(1, 3)}}, NeedsReview: true}
	enforceSumInvariant(already)
	if len(already.ReviewNotes) != 0 {
		t.Error("an already-reviewed result should not gain a sum note")
	}
}

func TestPropertyBlockedHeirsZero(t *testing.T) {
	// A son present excludes the son's son, all siblings, and lower kin, who
	// must then receive nothing across every school.
	h := heir.New().Set(heir.Son, 1).Set(heir.SonsSon, 1).
		Set(heir.FullBrother, 2).Set(heir.UterineSister, 1).Set(heir.ConsanguineBrother, 1)
	for _, m := range Madhahib() {
		r, err := Solve(estate.Case{DeceasedSex: heir.Male, Heirs: h}, m)
		if err != nil {
			t.Fatal(err)
		}
		for _, blocked := range []heir.Relation{heir.SonsSon, heir.FullBrother, heir.UterineSister, heir.ConsanguineBrother} {
			if share(r, blocked).Sign() != 0 {
				t.Errorf("%s: %s should be blocked, got %s", m.Name, blocked, share(r, blocked))
			}
		}
		if !share(r, heir.Son).IsWhole() {
			t.Errorf("%s: son should take the whole, got %s", m.Name, share(r, heir.Son))
		}
	}
}
