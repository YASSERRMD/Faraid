package solver

import (
	"fmt"

	"github.com/YASSERRMD/Faraid/internal/core/derivation"
	"github.com/YASSERRMD/Faraid/internal/core/estate"
	"github.com/YASSERRMD/Faraid/internal/core/rules"
)

// buildDerivation reconstructs the ordered audit trail for a solved result. For
// the normal pipeline it records the blocking, fixed shares, and residuary
// heirs with their references; for a special case it records the named case and
// the final shares.
func buildDerivation(c estate.Case, res Result) *derivation.Derivation {
	d := derivation.New()
	d.Add(derivation.Step{
		Stage:    derivation.StageEstate,
		Detail:   "distributable estate after funeral, debts, and bequests",
		Fraction: res.Distributable,
	})

	if res.SpecialCase != "" {
		d.Add(derivation.Step{Stage: derivation.StageSpecialCase, Detail: res.SpecialCase})
	}

	blocking := ResolveBlocking(c.Heirs)
	for _, r := range res.Excluded {
		ex := blocking.Excluded[r]
		d.Add(derivation.Step{Stage: derivation.StageBlocking, Relation: r, Detail: ex.Reason, Reference: ex.Reference})
	}

	if res.SpecialCase == "" {
		ctx := rules.Context{Heirs: c.Heirs}
		surviving := blocking.Surviving
		for _, r := range surviving.Relations() {
			if _, rule, ok := rules.FixedShare(r, ctx); ok {
				d.Add(derivation.Step{
					Stage:     derivation.StageFixedShare,
					Relation:  r,
					Detail:    rule.Condition,
					Reference: rule.Reference,
					Fraction:  rule.Share,
				})
			}
		}
		if resid := ClassifyResiduary(surviving); resid.Found() {
			for _, mem := range resid.Members {
				d.Add(derivation.Step{Stage: derivation.StageResiduary, Relation: mem.Relation, Detail: "residuary " + mem.Type.String()})
			}
		}
	}

	if res.Awl {
		d.Add(derivation.Step{Stage: derivation.StageAwl, Detail: "shares over-subscribed; base raised to the sum of shares"})
	}
	if res.Radd {
		d.Add(derivation.Step{Stage: derivation.StageRadd, Detail: "surplus returned to the eligible fixed-share heirs"})
	}
	if res.Base != nil {
		d.Add(derivation.Step{Stage: derivation.StageAsl, Detail: fmt.Sprintf("base of the problem is %s", res.Base)})
	}

	for _, s := range res.Shares {
		d.Add(derivation.Step{Stage: derivation.StageResult, Relation: s.Relation, Detail: "final share", Fraction: s.Fraction})
	}
	if res.Residue.Sign() > 0 {
		d.Add(derivation.Step{Stage: derivation.StageResult, Detail: "residue to the public treasury", Fraction: res.Residue})
	}
	return d
}
