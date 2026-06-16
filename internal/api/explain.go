package api

import (
	"encoding/json"
	"net/http"

	"github.com/YASSERRMD/Faraid/internal/core/solver"
	"github.com/YASSERRMD/Faraid/internal/llm"
)

type explanationDTO struct {
	Text         string `json:"text"`
	Consistent   bool   `json:"consistent"`
	Experimental bool   `json:"experimental"`
}

func (s *Server) handleExplain(w http.ResponseWriter, r *http.Request) {
	metricExplainTotal.Add(1)
	if s.completer == nil {
		writeError(w, http.StatusNotFound, "the trial explanation feature is disabled")
		return
	}
	var req solveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	c, m, err := toCase(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := solver.Solve(c, m)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	exp, err := llm.Explain(r.Context(), s.completer, result.Derivation.String(), engineFractions(result))
	if err != nil {
		writeError(w, http.StatusBadGateway, "the explanation provider failed")
		return
	}
	writeJSON(w, http.StatusOK, explanationDTO{
		Text:         exp.Text,
		Consistent:   exp.Consistent,
		Experimental: exp.Experimental,
	})
}

// engineFractions collects every fraction the engine produced, as the ground
// truth for the explanation consistency guard.
func engineFractions(r solver.Result) []string {
	out := make([]string, 0, len(r.Shares)+len(r.Derivation.Steps))
	for _, sh := range r.Shares {
		out = append(out, sh.Fraction.String())
	}
	if r.Residue.Sign() > 0 {
		out = append(out, r.Residue.String())
	}
	if r.Derivation != nil {
		for _, step := range r.Derivation.Steps {
			if !step.Fraction.IsZero() {
				out = append(out, step.Fraction.String())
			}
		}
	}
	return out
}
