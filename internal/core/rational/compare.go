package rational

import "math/big"

// Cmp compares f and o and returns -1 if f < o, 0 if f == o, and +1 if f > o.
func (f Fraction) Cmp(o Fraction) int {
	return f.rat().Cmp(o.rat())
}

// Equal reports whether f and o represent the same value.
func (f Fraction) Equal(o Fraction) bool {
	return f.Cmp(o) == 0
}

// Less reports whether f < o.
func (f Fraction) Less(o Fraction) bool {
	return f.Cmp(o) < 0
}

// LessEqual reports whether f <= o.
func (f Fraction) LessEqual(o Fraction) bool {
	return f.Cmp(o) <= 0
}

// Greater reports whether f > o.
func (f Fraction) Greater(o Fraction) bool {
	return f.Cmp(o) > 0
}

// GreaterEqual reports whether f >= o.
func (f Fraction) GreaterEqual(o Fraction) bool {
	return f.Cmp(o) >= 0
}

// Sign returns -1 if f < 0, 0 if f == 0, and +1 if f > 0.
func (f Fraction) Sign() int {
	return f.rat().Sign()
}

// IsZero reports whether f equals 0.
func (f Fraction) IsZero() bool {
	return f.rat().Sign() == 0
}

// IsWhole reports whether f equals 1, that is, the whole estate.
func (f Fraction) IsWhole() bool {
	return f.rat().Cmp(big.NewRat(1, 1)) == 0
}

// IsInteger reports whether f has denominator 1.
func (f Fraction) IsInteger() bool {
	return f.rat().IsInt()
}

// IsPositive reports whether f > 0.
func (f Fraction) IsPositive() bool {
	return f.rat().Sign() > 0
}

// IsNegative reports whether f < 0.
func (f Fraction) IsNegative() bool {
	return f.rat().Sign() < 0
}
