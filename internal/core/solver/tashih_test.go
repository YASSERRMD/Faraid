package solver

import (
	"math/big"
	"testing"
)

func TestTashihSingleGroup(t *testing.T) {
	// Three daughters share a numerator of 2 over base 3: indivisible.
	res := Tashih(big.NewInt(3), []TashihGroup{{Numerator: big.NewInt(2), Heads: 3}})
	if !res.Applied || res.Factor.Cmp(big.NewInt(3)) != 0 {
		t.Errorf("factor = %s, want 3", res.Factor)
	}
	if res.Base.Cmp(big.NewInt(9)) != 0 {
		t.Errorf("base = %s, want 9", res.Base)
	}
	if res.Numerators[0].Cmp(big.NewInt(6)) != 0 {
		t.Errorf("numerator = %s, want 6", res.Numerators[0])
	}
}

func TestTashihTwoGroupsCoprime(t *testing.T) {
	// Heads 2 and 3 are coprime, so the factor is their product 6.
	res := Tashih(big.NewInt(6), []TashihGroup{
		{Numerator: big.NewInt(1), Heads: 2},
		{Numerator: big.NewInt(2), Heads: 3},
	})
	if res.Factor.Cmp(big.NewInt(6)) != 0 {
		t.Errorf("factor = %s, want 6", res.Factor)
	}
	if res.Base.Cmp(big.NewInt(36)) != 0 {
		t.Errorf("base = %s, want 36", res.Base)
	}
	// Each group now divides evenly: 6/2 = 3 and 12/3 = 4.
	if res.Numerators[0].Cmp(big.NewInt(6)) != 0 || res.Numerators[1].Cmp(big.NewInt(12)) != 0 {
		t.Errorf("numerators = %v, want [6 12]", res.Numerators)
	}
}

func TestTashihTawafuq(t *testing.T) {
	// Heads 6 with numerator 4 share a factor of 2, so the group needs 3.
	res := Tashih(big.NewInt(8), []TashihGroup{{Numerator: big.NewInt(4), Heads: 6}})
	if res.Factor.Cmp(big.NewInt(3)) != 0 {
		t.Errorf("factor = %s, want 3", res.Factor)
	}
	// 12 / 6 = 2 parts each.
	if res.Numerators[0].Cmp(big.NewInt(12)) != 0 {
		t.Errorf("numerator = %s, want 12", res.Numerators[0])
	}
}

func TestTashihNotNeeded(t *testing.T) {
	// A single head, a zero share, and an already-divisible group: no change.
	res := Tashih(big.NewInt(6), []TashihGroup{
		{Numerator: big.NewInt(3), Heads: 1},
		{Numerator: big.NewInt(0), Heads: 5},
		{Numerator: big.NewInt(4), Heads: 2},
	})
	if res.Applied {
		t.Error("no correction should be needed")
	}
	if res.Factor.Cmp(big.NewInt(1)) != 0 || res.Base.Cmp(big.NewInt(6)) != 0 {
		t.Errorf("base/factor changed: base %s factor %s", res.Base, res.Factor)
	}
}

func TestTashihEmpty(t *testing.T) {
	res := Tashih(big.NewInt(6), nil)
	if res.Applied || res.Base.Cmp(big.NewInt(6)) != 0 {
		t.Errorf("empty tashih should not change the base, got %s", res.Base)
	}
}
