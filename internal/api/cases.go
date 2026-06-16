package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/YASSERRMD/Faraid/internal/core/solver"
	"github.com/YASSERRMD/Faraid/internal/store"
)

type saveCaseRequest struct {
	Name  string       `json:"name"`
	Input solveRequest `json:"input"`
}

type savedCaseDTO struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Input     solveRequest `json:"input"`
	CreatedAt string       `json:"createdAt,omitempty"`
}

// savedCaseToDTO converts a stored case to its response form, parsing the
// stored input JSON back into a request.
func savedCaseToDTO(sc store.SavedCase) (savedCaseDTO, error) {
	var input solveRequest
	if err := json.Unmarshal(sc.Input, &input); err != nil {
		return savedCaseDTO{}, err
	}
	return savedCaseDTO{
		ID:        sc.ID,
		Name:      sc.Name,
		Input:     input,
		CreatedAt: sc.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

func (s *Server) handleCreateCase(w http.ResponseWriter, r *http.Request) {
	var req saveCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	c, m, err := toCase(req.Input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := solver.Solve(c, m)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	inputJSON, _ := json.Marshal(req.Input)
	resultJSON, _ := json.Marshal(toSolveResult(result))
	sc, err := s.store.SaveCase(r.Context(), req.Name, inputJSON, resultJSON)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save case")
		return
	}
	dto, err := savedCaseToDTO(sc)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not encode case")
		return
	}
	writeJSON(w, http.StatusCreated, dto)
}

func (s *Server) handleListCases(w http.ResponseWriter, r *http.Request) {
	cases, err := s.store.ListCases(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list cases")
		return
	}
	out := make([]savedCaseDTO, 0, len(cases))
	for _, sc := range cases {
		dto, err := savedCaseToDTO(sc)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not encode case")
			return
		}
		out = append(out, dto)
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleGetCase(w http.ResponseWriter, r *http.Request) {
	sc, err := s.store.GetCase(r.Context(), chi.URLParam(r, "id"))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "case not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch case")
		return
	}
	dto, err := savedCaseToDTO(sc)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not encode case")
		return
	}
	writeJSON(w, http.StatusOK, dto)
}

func (s *Server) handleDeleteCase(w http.ResponseWriter, r *http.Request) {
	err := s.store.DeleteCase(r.Context(), chi.URLParam(r, "id"))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "case not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete case")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
