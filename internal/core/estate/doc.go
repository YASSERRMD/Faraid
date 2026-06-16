// Package estate models the input to an inheritance computation and processes
// the estate before the fixed shares are applied.
//
// The estate is settled in the classical order: funeral expenses first, then
// debts, then bequests (wasiyyah). Bequests to non-heirs are capped at one
// third of the estate remaining after funeral and debts unless the heirs
// consent to more. The amount left after these deductions is the distributable
// remainder that the solver allocates by faraid.
//
// Monetary inputs are integers in the smallest currency unit, while computed
// remainders are exact rationals so the one third cap stays precise.
package estate
