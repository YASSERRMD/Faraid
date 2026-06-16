package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestCompareDivergence(t *testing.T) {
	// The mushtaraka case diverges: the full brother shares under Maliki and
	// Shafi'i but not under Hanafi and Hanbali.
	body := `{"deceasedSex":"female","heirs":{"husband":1,"mother":1,"uterine brother":2,"full brother":1}}`
	rec := contractCheck(t, http.MethodPost, "/api/v1/compare", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body.String())
	}
	var cmp comparisonDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &cmp); err != nil {
		t.Fatal(err)
	}
	if len(cmp.Results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(cmp.Results))
	}
	if cmp.Results["Maliki"].Shares == nil || cmp.Results["Hanafi"].Madhhab != "Hanafi" {
		t.Errorf("unexpected results: %+v", cmp.Results)
	}
	joined := strings.Join(cmp.Divergences, "\n")
	if !strings.Contains(joined, "full brother") {
		t.Errorf("divergences should mention the full brother:\n%s", joined)
	}
}

func TestCompareAgreement(t *testing.T) {
	// All schools agree on this case, so no divergences are reported.
	body := `{"deceasedSex":"female","heirs":{"husband":1,"son":1,"daughter":1}}`
	rec := contractCheck(t, http.MethodPost, "/api/v1/compare", body)
	var cmp comparisonDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &cmp); err != nil {
		t.Fatal(err)
	}
	if len(cmp.Divergences) != 0 {
		t.Errorf("expected no divergences, got %v", cmp.Divergences)
	}
}

func TestCompareBadRequest(t *testing.T) {
	srv := NewServer()
	bodies := []string{
		`{bad`, // invalid JSON
		`{"deceasedSex":"other","heirs":{"son":1}}`,    // rejected while building the case
		`{"deceasedSex":"male","heirs":{"husband":1}}`, // builds, but fails validation when solved
	}
	for _, body := range bodies {
		if got := statusFor(srv, http.MethodPost, "/api/v1/compare", body); got != http.StatusBadRequest {
			t.Errorf("body %q = %d, want 400", body, got)
		}
	}
}
