package rational

import (
	"math/big"
	"testing"
)

func mustPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: expected panic, got none", name)
		}
	}()
	fn()
}

func TestConstructors(t *testing.T) {
	if !New(1, 2).Equal(New(2, 4)) {
		t.Error("New should reduce to lowest terms")
	}
	if !Zero().IsZero() {
		t.Error("Zero should be zero")
	}
	if !One().IsWhole() {
		t.Error("One should be whole")
	}
	if !FromInt(3).Equal(New(3, 1)) {
		t.Error("FromInt(3) should equal 3/1")
	}
	if got := New(-1, -2).String(); got != "1/2" {
		t.Errorf("New(-1,-2) = %q, want 1/2", got)
	}
	mustPanic(t, "New zero den", func() { New(1, 0) })
}

func TestFromBigRat(t *testing.T) {
	r := big.NewRat(3, 6)
	f := FromBigRat(r)
	if !f.Equal(New(1, 2)) {
		t.Errorf("FromBigRat = %s, want 1/2", f)
	}
	r.Add(r, big.NewRat(1, 1)) // mutate source; f must be unaffected
	if !f.Equal(New(1, 2)) {
		t.Error("FromBigRat did not copy its input")
	}
	if !FromBigRat(nil).IsZero() {
		t.Error("FromBigRat(nil) should be zero")
	}
}

func TestRatAccessorCopies(t *testing.T) {
	f := New(1, 2)
	r := f.Rat()
	r.Add(r, big.NewRat(1, 1)) // mutate the copy
	if !f.Equal(New(1, 2)) {
		t.Error("Rat() must return a copy that cannot mutate the Fraction")
	}
}

func TestZeroValue(t *testing.T) {
	var z Fraction // exercises the rat() nil branch
	if !z.IsZero() {
		t.Error("zero value Fraction should be zero")
	}
	if !z.Add(New(1, 3)).Equal(New(1, 3)) {
		t.Error("zero value should behave as 0 in addition")
	}
	if z.String() != "0" {
		t.Errorf("zero value String = %q, want 0", z.String())
	}
}

func TestArithmetic(t *testing.T) {
	half := New(1, 2)
	third := New(1, 3)
	if got := half.Add(third); !got.Equal(New(5, 6)) {
		t.Errorf("1/2 + 1/3 = %s, want 5/6", got)
	}
	if got := half.Sub(third); !got.Equal(New(1, 6)) {
		t.Errorf("1/2 - 1/3 = %s, want 1/6", got)
	}
	if got := half.Mul(third); !got.Equal(New(1, 6)) {
		t.Errorf("1/2 * 1/3 = %s, want 1/6", got)
	}
	if got := half.Div(third); !got.Equal(New(3, 2)) {
		t.Errorf("(1/2) / (1/3) = %s, want 3/2", got)
	}
	if got := half.Neg(); !got.Equal(New(-1, 2)) {
		t.Errorf("-(1/2) = %s, want -1/2", got)
	}
	if got := New(-1, 2).Abs(); !got.Equal(half) {
		t.Errorf("|-1/2| = %s, want 1/2", got)
	}
	mustPanic(t, "Div by zero", func() { half.Div(Zero()) })
}

func TestArithmeticImmutability(t *testing.T) {
	a := New(1, 2)
	b := New(1, 3)
	_ = a.Add(b)
	_ = a.Mul(b)
	if !a.Equal(New(1, 2)) || !b.Equal(New(1, 3)) {
		t.Error("operations must not mutate their operands")
	}
}

func TestSum(t *testing.T) {
	if !Sum().IsZero() {
		t.Error("Sum of nothing should be zero")
	}
	got := Sum(New(1, 6), New(1, 3), New(1, 2))
	if !got.IsWhole() {
		t.Errorf("1/6 + 1/3 + 1/2 = %s, want 1", got)
	}
}

func TestComparison(t *testing.T) {
	half := New(1, 2)
	third := New(1, 3)
	if half.Cmp(third) != 1 {
		t.Error("1/2 should compare greater than 1/3")
	}
	if third.Cmp(half) != -1 {
		t.Error("1/3 should compare less than 1/2")
	}
	if half.Cmp(New(2, 4)) != 0 {
		t.Error("1/2 should compare equal to 2/4")
	}
	if !half.Equal(New(2, 4)) || half.Equal(third) {
		t.Error("Equal failed")
	}
	if !third.Less(half) || half.Less(third) {
		t.Error("Less failed")
	}
	if !half.LessEqual(half) || !third.LessEqual(half) || half.LessEqual(third) {
		t.Error("LessEqual failed")
	}
	if !half.Greater(third) || third.Greater(half) {
		t.Error("Greater failed")
	}
	if !half.GreaterEqual(half) || !half.GreaterEqual(third) || third.GreaterEqual(half) {
		t.Error("GreaterEqual failed")
	}
}

func TestPredicates(t *testing.T) {
	if New(-1, 2).Sign() != -1 || Zero().Sign() != 0 || New(1, 2).Sign() != 1 {
		t.Error("Sign failed")
	}
	if !Zero().IsZero() || New(1, 2).IsZero() {
		t.Error("IsZero failed")
	}
	if !One().IsWhole() || New(1, 2).IsWhole() {
		t.Error("IsWhole failed")
	}
	if !FromInt(3).IsInteger() || New(1, 2).IsInteger() {
		t.Error("IsInteger failed")
	}
	if !New(1, 2).IsPositive() || New(-1, 2).IsPositive() || Zero().IsPositive() {
		t.Error("IsPositive failed")
	}
	if !New(-1, 2).IsNegative() || New(1, 2).IsNegative() || Zero().IsNegative() {
		t.Error("IsNegative failed")
	}
}

func TestFormat(t *testing.T) {
	if got := New(1, 2).String(); got != "1/2" {
		t.Errorf("String = %q, want 1/2", got)
	}
	if got := FromInt(3).String(); got != "3" {
		t.Errorf("String = %q, want 3", got)
	}
	if got := Zero().String(); got != "0" {
		t.Errorf("String = %q, want 0", got)
	}
	if got := New(-3, 4).Num(); got.Cmp(big.NewInt(-3)) != 0 {
		t.Errorf("Num = %s, want -3", got)
	}
	if got := New(-3, 4).Den(); got.Cmp(big.NewInt(4)) != 0 {
		t.Errorf("Den = %s, want 4", got)
	}
}

func TestParse(t *testing.T) {
	f, err := Parse("2/3")
	if err != nil || !f.Equal(New(2, 3)) {
		t.Errorf("Parse(2/3) = %s, %v", f, err)
	}
	f, err = Parse("5")
	if err != nil || !f.Equal(FromInt(5)) {
		t.Errorf("Parse(5) = %s, %v", f, err)
	}
	for _, bad := range []string{"abc", "", "1/0"} {
		if _, err := Parse(bad); err == nil {
			t.Errorf("Parse(%q) expected error", bad)
		}
	}
}

func TestGCD(t *testing.T) {
	cases := []struct{ a, b, want int64 }{
		{12, 18, 6},
		{0, 5, 5},
		{5, 0, 5},
		{0, 0, 0},
		{-12, 18, 6},
	}
	for _, c := range cases {
		got := GCD(big.NewInt(c.a), big.NewInt(c.b))
		if got.Cmp(big.NewInt(c.want)) != 0 {
			t.Errorf("GCD(%d,%d) = %s, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestLCM(t *testing.T) {
	cases := []struct{ a, b, want int64 }{
		{4, 6, 12},
		{3, 5, 15},
		{0, 5, 0},
		{5, 0, 0},
		{-4, 6, 12},
	}
	for _, c := range cases {
		got := LCM(big.NewInt(c.a), big.NewInt(c.b))
		if got.Cmp(big.NewInt(c.want)) != 0 {
			t.Errorf("LCM(%d,%d) = %s, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestLCMSlice(t *testing.T) {
	if got := LCMSlice(nil); got.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("LCMSlice(nil) = %s, want 1", got)
	}
	if got := LCMSlice([]*big.Int{big.NewInt(6)}); got.Cmp(big.NewInt(6)) != 0 {
		t.Errorf("LCMSlice([6]) = %s, want 6", got)
	}
	got := LCMSlice([]*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(6), big.NewInt(4)})
	if got.Cmp(big.NewInt(12)) != 0 {
		t.Errorf("LCMSlice(2,3,6,4) = %s, want 12", got)
	}
}
