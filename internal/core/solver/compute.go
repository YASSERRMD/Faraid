package solver

import (
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/rational"
	"github.com/YASSERRMD/Faraid/internal/core/rules"
)

// computeShares produces each heir's share of the whole distributable estate
// and the non-share metadata. It dispatches the special cases first, then runs
// the normal pipeline: blocking, fixed shares, residuary distribution, and the
// awl or radd adjustment.
func computeShares(h *heir.Heirs, m Madhhab) (map[heir.Relation]rational.Fraction, shareMeta) {
	meta := shareMeta{residue: rational.Zero()}

	if s, ok := Akdariyyah(h); ok {
		meta.specialCase = "akdariyyah"
		return s, meta
	}
	if IsMushtaraka(h) {
		s, _ := Mushtaraka(h, m.MushtarakaView)
		meta.specialCase = "mushtaraka"
		return s, meta
	}
	if s, ok := Gharrawayn(h); ok {
		meta.specialCase = "gharrawayn"
		return s, meta
	}

	surviving := ResolveBlocking(h).Surviving
	ctx := rules.Context{Heirs: h}

	if isJaddWaIkhwa(surviving, ctx) {
		return jaddShares(surviving, ctx, m, meta)
	}

	fixed := assignFurud(surviving, ctx)
	sumFixed := sumFractions(fixed)

	if sumFixed.Greater(rational.One()) {
		meta.awl = true
		out := make(map[heir.Relation]rational.Fraction, len(fixed))
		for r, f := range fixed {
			out[r] = f.Div(sumFixed)
		}
		return out, meta
	}

	residue := rational.One().Sub(sumFixed)
	if residue.IsPositive() {
		if resid := ClassifyResiduary(surviving); resid.Found() {
			for r, f := range DistributeResidue(residue, resid, surviving) {
				fixed[r] = fixed[r].Add(f)
			}
			return fixed, meta
		}
		if hasNonSpouseFurud(fixed) {
			meta.radd = true
			rr := ApplyRadd(fixed)
			meta.residue = rr.Surplus
			return rr.Shares, meta
		}
		// Only a spouse, or no heir at all: the residue passes to the treasury.
		// Distant kindred, when present, are handled by DistributeDhawuArham.
		meta.residue = residue
		return fixed, meta
	}

	return fixed, meta
}

// assignFurud assigns the fixed share to each surviving fixed-share heir and
// pools the grandmother and uterine shares.
func assignFurud(surviving *heir.Heirs, ctx rules.Context) map[heir.Relation]rational.Fraction {
	fixed := map[heir.Relation]rational.Fraction{}
	for _, r := range surviving.Relations() {
		if share, _, ok := rules.FixedShare(r, ctx); ok {
			fixed[r] = share
		}
	}
	poolGrandmothers(fixed)
	poolUterine(fixed, surviving)
	return fixed
}

// isJaddWaIkhwa reports whether the grandfather competes with full or
// consanguine siblings and no descendant is present. A surviving grandfather
// already implies no father, since the father excludes him.
func isJaddWaIkhwa(surviving *heir.Heirs, ctx rules.Context) bool {
	if !surviving.Present(heir.PaternalGrandfather) || ctx.HasInheritingDescendant() {
		return false
	}
	return surviving.Present(heir.FullBrother) || surviving.Present(heir.FullSister) ||
		surviving.Present(heir.ConsanguineBrother) || surviving.Present(heir.ConsanguineSister)
}

// jaddShares allocates the grandfather and sibling shares after the other
// fixed-share heirs take theirs.
func jaddShares(surviving *heir.Heirs, ctx rules.Context, m Madhhab, meta shareMeta) (map[heir.Relation]rational.Fraction, shareMeta) {
	meta.specialCase = "jadd wa ikhwa"

	fixed := map[heir.Relation]rational.Fraction{}
	for _, r := range surviving.Relations() {
		if r == heir.PaternalGrandfather || isCompetingSibling(r) {
			continue
		}
		if share, _, ok := rules.FixedShare(r, ctx); ok {
			fixed[r] = share
		}
	}
	poolGrandmothers(fixed)

	available := rational.One().Sub(sumFractions(fixed))
	jadd := GrandfatherWithSiblings(available, surviving, m.GrandfatherView)
	fixed[heir.PaternalGrandfather] = jadd.GrandfatherShare
	for r, f := range jadd.SiblingShares {
		fixed[r] = f
	}
	if jadd.NeedsReview {
		meta.needsReview = true
		meta.notes = append(meta.notes, jadd.ReviewNote)
	}
	return fixed, meta
}
