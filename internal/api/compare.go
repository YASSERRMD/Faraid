package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/YASSERRMD/Faraid/internal/core/solver"
)

type comparisonDTO struct {
	Results     map[string]solveResultDTO `json:"results"`
	Divergences []string                  `json:"divergences,omitempty"`
}

func (s *Server) handleCompare(w http.ResponseWriter, r *http.Request) {
	metricCompareTotal.Add(1)
	var req solveRequest // the madhhab field is ignored here
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	c, err := toCaseInput(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	results := make(map[string]solveResultDTO, 4)
	for _, m := range solver.Madhahib() {
		res, err := solver.Solve(c, m)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		results[m.Name] = toSolveResult(res)
	}

	writeJSON(w, http.StatusOK, comparisonDTO{
		Results:     results,
		Divergences: computeDivergences(results),
	})
}

// computeDivergences reports, per heir, the schools whose share fractions
// differ, in a stable, readable form.
func computeDivergences(results map[string]solveResultDTO) []string {
	schools := make([]string, 0, len(results))
	for _, m := range solver.Madhahib() {
		if _, ok := results[m.Name]; ok {
			schools = append(schools, m.Name)
		}
	}

	relSet := map[string]bool{}
	for _, res := range results {
		for _, sh := range res.Shares {
			relSet[sh.Relation] = true
		}
	}
	rels := make([]string, 0, len(relSet))
	for rel := range relSet {
		rels = append(rels, rel)
	}
	sort.Strings(rels)

	var out []string
	for _, rel := range rels {
		fracs := make(map[string]string, len(schools))
		distinct := map[string]bool{}
		for _, sch := range schools {
			f := "0"
			for _, sh := range results[sch].Shares {
				if sh.Relation == rel {
					f = sh.Fraction
					break
				}
			}
			fracs[sch] = f
			distinct[f] = true
		}
		if len(distinct) > 1 {
			parts := make([]string, 0, len(schools))
			for _, sch := range schools {
				parts = append(parts, fmt.Sprintf("%s %s", sch, fracs[sch]))
			}
			out = append(out, fmt.Sprintf("%s: %s", rel, strings.Join(parts, ", ")))
		}
	}
	return out
}
