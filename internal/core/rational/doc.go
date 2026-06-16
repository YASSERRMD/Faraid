// Package rational provides exact rational arithmetic for inheritance shares.
//
// Shares are fractions of the estate and are represented by the Fraction type,
// a small immutable wrapper over math/big.Rat, so the legal core never relies
// on floating point. Every operation returns a new Fraction and never mutates
// its operands, which keeps calculations deterministic and safe to share.
//
// The package also exposes integer helpers (GCD, LCM) used by the solver to
// compute asl al-mas'ala, the least common denominator of all assigned shares.
package rational
