package solver

import (
	"sort"

	"github.com/YASSERRMD/Faraid/internal/core/estate"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

// Solve computes the full distribution of a case under a school. It settles the
// estate, computes each heir's share of the distributable remainder, derives
// the base of the problem, applies tashih so each slot divides into whole
// parts, and converts shares to exact amounts.
func Solve(c estate.Case, m Madhhab) (Result, error) {
	dist, err := c.Estate.PreDistribute()
	if err != nil {
		return Result{}, err
	}
	if err := c.Validate(); err != nil {
		return Result{}, err
	}

	shares, meta := computeShares(c.Heirs, m)

	result := Result{
		Madhhab:       m.Name,
		Distributable: dist.Distributable,
		SpecialCase:   meta.specialCase,
		Awl:           meta.awl,
		Radd:          meta.radd,
		Residue:       meta.residue,
		NeedsReview:   meta.needsReview,
		ReviewNotes:   meta.notes,
	}

	// Order the heirs canonically for a deterministic result.
	rels := make([]heir.Relation, 0, len(shares))
	for r := range shares {
		rels = append(rels, r)
	}
	sort.Slice(rels, func(i, j int) bool { return rels[i] < rels[j] })

	// Derive the base and apply tashih so every slot divides into whole parts.
	asl := ComputeAsl(shares)
	groups := make([]TashihGroup, len(rels))
	for i, r := range rels {
		groups[i] = TashihGroup{Numerator: asl.Numerators[r], Heads: int64(c.Heirs.Count(r))}
	}
	tash := Tashih(asl.Base, groups)
	result.Base = tash.Base

	for i, r := range rels {
		f := shares[r]
		result.Shares = append(result.Shares, HeirShare{
			Relation: r,
			Count:    c.Heirs.Count(r),
			Fraction: f,
			Parts:    tash.Numerators[i],
			Amount:   f.Mul(dist.Distributable),
		})
	}

	for r := range ResolveBlocking(c.Heirs).Excluded {
		result.Excluded = append(result.Excluded, r)
	}
	sort.Slice(result.Excluded, func(i, j int) bool { return result.Excluded[i] < result.Excluded[j] })

	result.Derivation = buildDerivation(c, result)

	return result, nil
}
