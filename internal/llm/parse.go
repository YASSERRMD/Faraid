package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YASSERRMD/Faraid/internal/core/estate"
	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

// CaseProposal is the LLM's structured interpretation of a natural-language
// case description. It is experimental and unconfirmed: the caller must show it
// to a human for confirmation before building an engine case and solving. The
// parser never solves and is never the source of a legal result.
type CaseProposal struct {
	DeceasedSex string         `json:"deceasedSex"`
	EstateTotal int64          `json:"estateTotal"`
	Heirs       map[string]int `json:"heirs"`
	Notes       string         `json:"notes"`
}

// ToCase converts a proposal into a validated engine case. It returns an error
// when the proposal does not describe a valid case.
func (p CaseProposal) ToCase() (estate.Case, error) {
	var sex heir.Sex
	switch p.DeceasedSex {
	case "male":
		sex = heir.Male
	case "female":
		sex = heir.Female
	default:
		return estate.Case{}, fmt.Errorf("llm: invalid deceasedSex %q", p.DeceasedSex)
	}

	h := heir.New()
	for name, count := range p.Heirs {
		r, ok := heir.ParseRelation(name)
		if !ok {
			return estate.Case{}, fmt.Errorf("llm: unknown heir %q", name)
		}
		h.Set(r, count)
	}

	c := estate.Case{DeceasedSex: sex, Estate: estate.Estate{Total: p.EstateTotal}, Heirs: h}
	if err := c.Validate(); err != nil {
		return estate.Case{}, err
	}
	return c, nil
}

// ParseCase asks the model to interpret natural-language text (Arabic or
// English) as a structured case, then validates the result. It returns an error
// when the model output is not valid JSON, fails schema validation, or
// describes an impossible case, so the caller can fall back to manual entry. It
// never solves: the returned proposal is unconfirmed and must be confirmed by a
// human before the deterministic engine is run.
func ParseCase(ctx context.Context, c Completer, text string) (CaseProposal, error) {
	resp, err := c.Complete(ctx, Request{System: caseParsePrompt(), Prompt: text, MaxTokens: 1024})
	if err != nil {
		return CaseProposal{}, err
	}

	var p CaseProposal
	if err := json.Unmarshal([]byte(extractJSONObject(resp.Text)), &p); err != nil {
		return CaseProposal{}, fmt.Errorf("llm: model output was not valid JSON: %w", err)
	}
	if _, err := p.ToCase(); err != nil {
		return CaseProposal{}, err
	}
	return p, nil
}

// extractJSONObject returns the substring from the first opening brace to the
// last closing brace, which tolerates markdown fences or surrounding prose.
func extractJSONObject(s string) string {
	i := strings.Index(s, "{")
	j := strings.LastIndex(s, "}")
	if i >= 0 && j >= i {
		return s[i : j+1]
	}
	return s
}

// caseParsePrompt builds the strict-JSON instruction, listing the recognized
// heir labels so the model uses only those.
func caseParsePrompt() string {
	labels := make([]string, 0, len(heir.AllRelations()))
	for _, r := range heir.AllRelations() {
		labels = append(labels, r.String())
	}
	return "You convert a natural-language description of an Islamic inheritance case " +
		"into strict JSON. Output ONLY a JSON object, with no prose and no markdown fences. " +
		"Schema: {\"deceasedSex\": \"male\" or \"female\", \"estateTotal\": integer in the " +
		"smallest currency unit (0 if unknown), \"heirs\": an object mapping heir label to a " +
		"positive integer count, \"notes\": a short string}. " +
		"Use only these exact heir labels: " + strings.Join(labels, ", ") + ". " +
		"Do not compute shares; only describe the heirs present."
}
