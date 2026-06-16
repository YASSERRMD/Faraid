package api

import (
	"encoding/json"
	"net/http"

	"github.com/YASSERRMD/Faraid/internal/llm"
)

type parseRequest struct {
	Text string `json:"text"`
}

type parseProposalDTO struct {
	DeceasedSex string         `json:"deceasedSex"`
	Heirs       map[string]int `json:"heirs"`
	Notes       string         `json:"notes,omitempty"`
}

func (s *Server) handleParse(w http.ResponseWriter, r *http.Request) {
	metricParseTotal.Add(1)
	if s.completer == nil {
		writeError(w, http.StatusNotFound, "the trial parse feature is disabled")
		return
	}
	var req parseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Text == "" {
		writeError(w, http.StatusBadRequest, "text is required")
		return
	}
	prop, err := llm.ParseCase(r.Context(), s.completer, req.Text)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, parseProposalDTO{
		DeceasedSex: prop.DeceasedSex,
		Heirs:       prop.Heirs,
		Notes:       prop.Notes,
	})
}
