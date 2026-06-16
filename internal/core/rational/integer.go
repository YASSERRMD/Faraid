package rational

import "math/big"

// GCD returns the greatest common divisor of a and b as a non-negative integer.
// The sign of the inputs is ignored. GCD(0, 0) is 0.
func GCD(a, b *big.Int) *big.Int {
	return new(big.Int).GCD(nil, nil, new(big.Int).Abs(a), new(big.Int).Abs(b))
}

// LCM returns the least common multiple of a and b as a non-negative integer.
// The sign of the inputs is ignored. If either input is zero, LCM is 0.
func LCM(a, b *big.Int) *big.Int {
	if a.Sign() == 0 || b.Sign() == 0 {
		return big.NewInt(0)
	}
	g := GCD(a, b)
	res := new(big.Int).Div(new(big.Int).Abs(a), g)
	res.Mul(res, new(big.Int).Abs(b))
	return res
}

// LCMSlice returns the least common multiple of all the given integers. The
// LCM of an empty slice is 1, the multiplicative identity, so it can seed a
// running denominator. This is used to compute asl al-mas'ala from the
// denominators of all assigned shares.
func LCMSlice(nums []*big.Int) *big.Int {
	if len(nums) == 0 {
		return big.NewInt(1)
	}
	acc := new(big.Int).Abs(nums[0])
	for _, n := range nums[1:] {
		acc = LCM(acc, n)
	}
	return acc
}
