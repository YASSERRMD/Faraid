package solver

import (
	"fmt"
	"sort"

	"github.com/YASSERRMD/Faraid/internal/core/estate"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
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

	enforceSumInvariant(&result)
	result.Derivation = buildDerivation(c, result)

	return result, nil
}

// enforceSumInvariant flags a result whose shares plus any treasury residue do
// not consume the whole estate. The engine never returns a silently incomplete
// result: a configuration it does not fully resolve is marked for review.
func enforceSumInvariant(result *Result) {
	total := result.Residue
	for _, s := range result.Shares {
		total = total.Add(s.Fraction)
	}
	if !total.Equal(rational.One()) && !result.NeedsReview {
		result.NeedsReview = true
		result.ReviewNotes = append(result.ReviewNotes,
			fmt.Sprintf("shares sum to %s, not the whole estate; configuration not fully resolved", total))
	}
}
