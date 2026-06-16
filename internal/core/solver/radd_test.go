package solver

import (
	"testing"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

func sumShares(m map[heir.Relation]rational.Fraction) rational.Fraction {
	s := rational.Zero()
	for _, f := range m {
		s = s.Add(f)
	}
	return s
}

func TestRaddNoSpouse(t *testing.T) {
	// Mother 1/6 and one daughter 1/2, no residuary: surplus returned 1:3.
	res := ApplyRadd(map[heir.Relation]rational.Fraction{
		heir.Mother:   rational.New(1, 6),
		heir.Daughter: rational.New(1, 2),
	})
	if !res.Applied {
		t.Fatal("radd should apply")
	}
	if !res.Shares[heir.Mother].Equal(rational.New(1, 4)) {
		t.Errorf("mother = %s, want 1/4", res.Shares[heir.Mother])
	}
	if !res.Shares[heir.Daughter].Equal(rational.New(3, 4)) {
		t.Errorf("daughter = %s, want 3/4", res.Shares[heir.Daughter])
	}
	if !sumShares(res.Shares).IsWhole() {
		t.Errorf("shares should sum to 1, got %s", sumShares(res.Shares))
	}
}

func TestRaddWithSpouseSingleEligible(t *testing.T) {
	// Wife 1/8 and one daughter 1/2: wife keeps 1/8, daughter takes 7/8.
	res := ApplyRadd(map[heir.Relation]rational.Fraction{
		heir.Wife:     rational.New(1, 8),
		heir.Daughter: rational.New(1, 2),
	})
	if !res.Shares[heir.Wife].Equal(rational.New(1, 8)) {
		t.Errorf("wife = %s, want 1/8", res.Shares[heir.Wife])
	}
	if !res.Shares[heir.Daughter].Equal(rational.New(7, 8)) {
		t.Errorf("daughter = %s, want 7/8", res.Shares[heir.Daughter])
	}
	if !sumShares(res.Shares).IsWhole() {
		t.Errorf("shares should sum to 1, got %s", sumShares(res.Shares))
	}
}

func TestRaddWithSpouseMultipleEligible(t *testing.T) {
	// Wife 1/8, mother 1/6, daughter 1/2. Spouse excluded; mother and daughter
	// share the remaining 7/8 in their 1:3 ratio.
	res := ApplyRadd(map[heir.Relation]rational.Fraction{
		heir.Wife:     rational.New(1, 8),
		heir.Mother:   rational.New(1, 6),
		heir.Daughter: rational.New(1, 2),
	})
	if !res.Shares[heir.Wife].Equal(rational.New(1, 8)) {
		t.Errorf("wife = %s, want 1/8", res.Shares[heir.Wife])
	}
	if !res.Shares[heir.Mother].Equal(rational.New(7, 32)) {
		t.Errorf("mother = %s, want 7/32", res.Shares[heir.Mother])
	}
	if !res.Shares[heir.Daughter].Equal(rational.New(21, 32)) {
		t.Errorf("daughter = %s, want 21/32", res.Shares[heir.Daughter])
	}
	if !sumShares(res.Shares).IsWhole() {
		t.Errorf("shares should sum to 1, got %s", sumShares(res.Shares))
	}
}

func TestRaddOnlySpouseGoesToTreasury(t *testing.T) {
	// Husband alone: keeps 1/2, the surplus passes to the treasury.
	res := ApplyRadd(map[heir.Relation]rational.Fraction{heir.Husband: rational.New(1, 2)})
	if !res.Applied {
		t.Fatal("radd should apply")
	}
	if !res.Shares[heir.Husband].Equal(rational.New(1, 2)) {
		t.Errorf("husband = %s, want 1/2", res.Shares[heir.Husband])
	}
	if !res.Surplus.Equal(rational.New(1, 2)) {
		t.Errorf("surplus = %s, want 1/2", res.Surplus)
	}
}

func TestRaddNotNeeded(t *testing.T) {
	// Shares already consume the whole: no radd.
	res := ApplyRadd(map[heir.Relation]rational.Fraction{
		heir.Husband:  rational.New(1, 2),
		heir.Daughter: rational.New(1, 2),
	})
	if res.Applied {
		t.Error("radd should not apply when shares sum to one")
	}
	if !res.Shares[heir.Daughter].Equal(rational.New(1, 2)) {
		t.Errorf("shares should be unchanged, got %v", res.Shares)
	}
}
