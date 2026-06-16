package rational

import (
	"fmt"
	"math/big"
)

// String returns the canonical textual form of f: "n/d" in lowest terms, or
// "n" when f is an integer, or "0" when f is zero.
func (f Fraction) String() string {
	return f.rat().RatString()
}

// Num returns a copy of the numerator in lowest terms. Its sign carries the
// sign of the fraction.
func (f Fraction) Num() *big.Int {
	return new(big.Int).Set(f.rat().Num())
}

// Den returns a copy of the denominator in lowest terms. It is always
// positive.
func (f Fraction) Den() *big.Int {
	return new(big.Int).Set(f.rat().Denom())
}

// Parse reads a fraction from its textual form. It accepts "n/d" and "n"
// forms. It returns an error if s is not a valid fraction or has a zero
// denominator.
func Parse(s string) (Fraction, error) {
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		return Fraction{}, fmt.Errorf("rational: cannot parse %q as a fraction", s)
	}
	return Fraction{r: r}, nil
}
