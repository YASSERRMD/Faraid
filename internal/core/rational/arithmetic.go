package rational

import "math/big"

// Add returns f + o.
func (f Fraction) Add(o Fraction) Fraction {
	return Fraction{r: new(big.Rat).Add(f.rat(), o.rat())}
}

// Sub returns f - o.
func (f Fraction) Sub(o Fraction) Fraction {
	return Fraction{r: new(big.Rat).Sub(f.rat(), o.rat())}
}

// Mul returns f * o.
func (f Fraction) Mul(o Fraction) Fraction {
	return Fraction{r: new(big.Rat).Mul(f.rat(), o.rat())}
}

// Div returns f / o. It panics if o is zero.
func (f Fraction) Div(o Fraction) Fraction {
	if o.rat().Sign() == 0 {
		panic("rational: division by zero")
	}
	return Fraction{r: new(big.Rat).Quo(f.rat(), o.rat())}
}

// Neg returns -f.
func (f Fraction) Neg() Fraction {
	return Fraction{r: new(big.Rat).Neg(f.rat())}
}

// Abs returns the absolute value of f.
func (f Fraction) Abs() Fraction {
	return Fraction{r: new(big.Rat).Abs(f.rat())}
}

// Sum returns the sum of all the given fractions. Sum of no fractions is zero.
func Sum(fs ...Fraction) Fraction {
	acc := new(big.Rat)
	for _, f := range fs {
		acc.Add(acc, f.rat())
	}
	return Fraction{r: acc}
}
