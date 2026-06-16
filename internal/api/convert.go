package api

import (
	"fmt"

	"github.com/YASSERRMD/Faraid/internal/core/estate"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
	"github.com/YASSERRMD/Faraid/internal/core/solver"
)

// toCaseInput maps a request to a core case (deceased sex, estate, heirs)
// without resolving a madhhab, so it can be solved across one or all schools.
func toCaseInput(req solveRequest) (estate.Case, error) {
	var sex heir.Sex
	switch req.DeceasedSex {
	case "male":
		sex = heir.Male
	case "female":
		sex = heir.Female
	default:
		return estate.Case{}, fmt.Errorf("invalid deceasedSex %q", req.DeceasedSex)
	}

	h := heir.New()
	for name, count := range req.Heirs {
		r, ok := heir.ParseRelation(name)
		if !ok {
			return estate.Case{}, fmt.Errorf("unknown heir %q", name)
		}
		h.Set(r, count)
	}

	return estate.Case{
		DeceasedSex: sex,
		Estate: estate.Estate{
			Total:                       req.Estate.Total,
			Funeral:                     req.Estate.Funeral,
			Debts:                       req.Estate.Debts,
			Bequests:                    req.Estate.Bequests,
			HeirsConsentToExcessBequest: req.Estate.HeirsConsentToExcessBequest,
		},
		Heirs: h,
	}, nil
}

// toCase maps a solve request to a core case and a madhhab profile.
func toCase(req solveRequest) (estate.Case, solver.Madhhab, error) {
	c, err := toCaseInput(req)
	if err != nil {
		return estate.Case{}, solver.Madhhab{}, err
	}
	m, ok := solver.MadhhabByName(req.Madhhab)
	if !ok {
		return estate.Case{}, solver.Madhhab{}, fmt.Errorf("unknown madhhab %q", req.Madhhab)
	}
	return c, m, nil
}

// toSolveResult maps a core result to its response form.
func toSolveResult(r solver.Result) solveResultDTO {
	out := solveResultDTO{
		Madhhab:       r.Madhhab,
		Distributable: r.Distributable.String(),
		Base:          r.Base.Int64(),
		SpecialCase:   r.SpecialCase,
		Awl:           r.Awl,
		Radd:          r.Radd,
		NeedsReview:   r.NeedsReview,
		ReviewNotes:   r.ReviewNotes,
	}
	if r.Residue.Sign() > 0 {
		out.Residue = r.Residue.String()
	}
	for _, s := range r.Shares {
		out.Shares = append(out.Shares, heirShareDTO{
			Relation: s.Relation.String(),
			Count:    s.Count,
			Fraction: s.Fraction.String(),
			Parts:    s.Parts.Int64(),
			Amount:   s.Amount.String(),
		})
	}
	for _, ex := range r.Excluded {
		out.Excluded = append(out.Excluded, ex.String())
	}
	if r.Derivation != nil {
		for _, step := range r.Derivation.Steps {
			d := derivationStepDTO{
				Stage:     string(step.Stage),
				Detail:    step.Detail,
				Reference: step.Reference,
			}
			if step.Relation.Valid() {
				d.Relation = step.Relation.String()
			}
			if !step.Fraction.IsZero() {
				d.Fraction = step.Fraction.String()
			}
			out.Derivation = append(out.Derivation, d)
		}
	}
	return out
}
