package solver

import (
	"math/big"

	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// ShareTable is a solved problem: a base and an integer numerator for each heir,
// keyed by an identity label. It is the form used to chain successive deaths,
// because the heirs of a later death are not relatives of the first deceased
// and so cannot be keyed by relation.
type ShareTable struct {
	Base   *big.Int
	Shares map[string]*big.Int
}

// Fraction returns an heir's share of the whole estate, or zero when absent.
func (s ShareTable) Fraction(label string) rational.Fraction {
	n := s.Shares[label]
	if n == nil {
		return rational.Zero()
	}
	return rational.FromBigRat(new(big.Rat).SetFrac(n, s.Base))
}

// Munasakha combines a first estate with the estate of one of its heirs who
// died before distribution. The argument deceased is that heir's label in
// first, and second is the solved distribution of that heir's own estate. The
// deceased's portion is divided among the second heirs and all shares are
// unified over a common base, the jami'ah. An heir appearing in both estates
// has the two portions added. The second estate's base must be positive.
func Munasakha(first ShareTable, deceased string, second ShareTable) ShareTable {
	d := first.Shares[deceased]
	if d == nil {
		d = new(big.Int)
	}

	g := rational.GCD(d, second.Base)
	firstMult := new(big.Int).Div(second.Base, g)
	secondMult := new(big.Int).Div(d, g)

	result := map[string]*big.Int{}
	for label, num := range first.Shares {
		if label == deceased {
			continue
		}
		addShare(result, label, new(big.Int).Mul(num, firstMult))
	}
	for label, num := range second.Shares {
		addShare(result, label, new(big.Int).Mul(num, secondMult))
	}

	return ShareTable{
		Base:   new(big.Int).Mul(first.Base, firstMult),
		Shares: result,
	}
}

// addShare adds n to the share stored under label, accumulating when an heir
// inherits from more than one estate in the chain.
func addShare(m map[string]*big.Int, label string, n *big.Int) {
	if existing, ok := m[label]; ok {
		existing.Add(existing, n)
		return
	}
	m[label] = n
}
