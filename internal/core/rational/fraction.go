package rational

import "math/big"

// Fraction is an exact rational number used to represent a share of an estate.
//
// It wraps math/big.Rat. The zero value is valid and equals 0. Fraction is
// immutable: every method returns a new Fraction and never mutates its
// receiver or arguments, so values may be copied and shared freely.
type Fraction struct {
	r *big.Rat
}

// rat returns the underlying big.Rat, treating the zero value as 0. The
// returned value must only be read, never mutated, by callers inside this
// package; values that escape the package are always defensively copied.
func (f Fraction) rat() *big.Rat {
	if f.r == nil {
		return new(big.Rat)
	}
	return f.r
}

// New returns the fraction num/den in lowest terms. It panics if den is zero,
// which is always a programming error in this domain.
func New(num, den int64) Fraction {
	if den == 0 {
		panic("rational: zero denominator")
	}
	return Fraction{r: big.NewRat(num, den)}
}

// Zero returns the fraction 0.
func Zero() Fraction {
	return Fraction{r: new(big.Rat)}
}

// One returns the fraction 1, representing the whole estate.
func One() Fraction {
	return Fraction{r: big.NewRat(1, 1)}
}

// FromInt returns the fraction n/1.
func FromInt(n int64) Fraction {
	return Fraction{r: big.NewRat(n, 1)}
}

// FromBigRat returns a Fraction with the same value as r. The input is copied,
// so later mutation of r does not affect the returned Fraction. A nil r yields
// the zero fraction.
func FromBigRat(r *big.Rat) Fraction {
	if r == nil {
		return Fraction{}
	}
	return Fraction{r: new(big.Rat).Set(r)}
}

// Rat returns a copy of the value as a *big.Rat. Mutating the result does not
// affect the Fraction.
func (f Fraction) Rat() *big.Rat {
	return new(big.Rat).Set(f.rat())
}
