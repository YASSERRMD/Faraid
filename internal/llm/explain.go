package llm

import (
	"context"
	"regexp"

	"github.com/YASSERRMD/Faraid/internal/core/rational"
)

// Explanation is the trial-tier prose explanation of a result.
type Explanation struct {
	// Text is the prose explanation, or the raw derivation when the guard
	// rejected the prose.
	Text string
	// Consistent reports whether the prose passed the numeric guard.
	Consistent bool
	// Experimental is always true: this output is non-authoritative.
	Experimental bool
}

var fractionRe = regexp.MustCompile(`\d+/\d+`)

// Explain asks the model to rephrase the derivation as readable prose, then
// runs a numeric consistency guard: every fraction mentioned in the prose must
// match, by value, a fraction the engine produced. On any drift it discards the
// prose and returns the raw derivation, flagged not consistent, so a wrong
// number can never be shown as if it came from the engine.
func Explain(ctx context.Context, c Completer, derivation string, engineFractions []string) (Explanation, error) {
	valid := map[string]bool{}
	for _, f := range engineFractions {
		if fr, err := rational.Parse(f); err == nil {
			valid[fr.String()] = true
		}
	}

	resp, err := c.Complete(ctx, Request{System: explainPrompt(), Prompt: derivation, MaxTokens: 1024})
	if err != nil {
		return Explanation{}, err
	}

	if !fractionsConsistent(resp.Text, valid) {
		return Explanation{Text: derivation, Consistent: false, Experimental: true}, nil
	}
	return Explanation{Text: resp.Text, Consistent: true, Experimental: true}, nil
}

// fractionsConsistent reports whether every fraction in the prose matches, by
// value, one of the engine's fractions.
func fractionsConsistent(prose string, valid map[string]bool) bool {
	for _, m := range fractionRe.FindAllString(prose, -1) {
		fr, err := rational.Parse(m)
		if err != nil {
			continue
		}
		if !valid[fr.String()] {
			return false
		}
	}
	return true
}

func explainPrompt() string {
	return "You rephrase a deterministic Islamic inheritance derivation into clear, " +
		"plain-language prose for a layperson. The derivation is the ground truth. Do not " +
		"change, add, recompute, or round any number or fraction, and do not introduce any " +
		"fraction that is not already present in the derivation. Explain only what the " +
		"derivation states, in a few short sentences."
}
