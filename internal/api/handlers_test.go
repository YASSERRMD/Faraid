package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

// contractCheck serves a request with the real handler and validates both the
// request and the response against the OpenAPI contract.
func contractCheck(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	doc := loadSpec(t)
	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		t.Fatalf("build spec router: %v", err)
	}

	handler := NewServer().Router()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newJSONRequest(method, path, body))

	valReq := newJSONRequest(method, path, body)
	route, pathParams, err := router.FindRoute(valReq)
	if err != nil {
		t.Fatalf("find route %s %s: %v", method, path, err)
	}
	reqInput := &openapi3filter.RequestValidationInput{Request: valReq, PathParams: pathParams, Route: route}
	if err := openapi3filter.ValidateRequest(context.Background(), reqInput); err != nil {
		t.Fatalf("request does not match the spec: %v", err)
	}
	respInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: reqInput,
		Status:                 rec.Code,
		Header:                 rec.Header(),
	}
	respInput.SetBodyBytes(rec.Body.Bytes())
	if err := openapi3filter.ValidateResponse(context.Background(), respInput); err != nil {
		t.Fatalf("response does not match the spec: %v\nbody: %s", err, rec.Body.String())
	}
	return rec
}

func newJSONRequest(method, path, body string) *http.Request {
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, strings.NewReader(body))
	}
	r.Header.Set("Content-Type", "application/json")
	return r
}

func TestSolveEndpoint(t *testing.T) {
	body := `{"deceasedSex":"female","heirs":{"husband":1,"son":1,"daughter":1},"madhhab":"Hanafi"}`
	rec := contractCheck(t, http.MethodPost, "/api/v1/solve", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", rec.Code, rec.Body.String())
	}
	var res solveResultDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	shares := map[string]string{}
	for _, s := range res.Shares {
		shares[s.Relation] = s.Fraction
	}
	if shares["husband"] != "1/4" || shares["son"] != "1/2" || shares["daughter"] != "1/4" {
		t.Errorf("unexpected shares: %v", shares)
	}
	if res.Base != 4 || len(res.Derivation) == 0 {
		t.Errorf("base = %d, derivation steps = %d", res.Base, len(res.Derivation))
	}
}

func TestSolveBadRequests(t *testing.T) {
	cases := map[string]string{
		"bad json":       `{not json`,
		"unknown heir":   `{"deceasedSex":"male","heirs":{"cousin twice removed":1},"madhhab":"Hanafi"}`,
		"unknown school": `{"deceasedSex":"male","heirs":{"son":1},"madhhab":"Jafari"}`,
		"bad sex":        `{"deceasedSex":"other","heirs":{"son":1},"madhhab":"Hanafi"}`,
		"inconsistent":   `{"deceasedSex":"male","heirs":{"husband":1},"madhhab":"Hanafi"}`,
	}
	handler := NewServer().Router()
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, newJSONRequest(http.MethodPost, "/api/v1/solve", body))
			if rec.Code != http.StatusBadRequest {
				t.Errorf("%s: status = %d, want 400", name, rec.Code)
			}
		})
	}
}

func TestSolveTreasuryAndExcluded(t *testing.T) {
	// Husband alone: a treasury residue is reported.
	rec := contractCheck(t, http.MethodPost, "/api/v1/solve",
		`{"deceasedSex":"female","heirs":{"husband":1},"madhhab":"Hanafi"}`)
	var soleSpouse solveResultDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &soleSpouse); err != nil {
		t.Fatal(err)
	}
	if soleSpouse.Residue != "1/2" {
		t.Errorf("residue = %q, want 1/2", soleSpouse.Residue)
	}

	// A son excludes the full brother, who is reported as excluded.
	rec = contractCheck(t, http.MethodPost, "/api/v1/solve",
		`{"deceasedSex":"female","heirs":{"husband":1,"son":1,"full brother":1},"madhhab":"Hanafi"}`)
	var blocked solveResultDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &blocked); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, ex := range blocked.Excluded {
		if ex == "full brother" {
			found = true
		}
	}
	if !found {
		t.Errorf("full brother should be excluded, got %v", blocked.Excluded)
	}
}

func TestHealthAndMadhahibContract(t *testing.T) {
	rec := contractCheck(t, http.MethodGet, "/api/v1/healthz", "")
	if rec.Code != http.StatusOK {
		t.Errorf("health status = %d", rec.Code)
	}
	rec = contractCheck(t, http.MethodGet, "/api/v1/madhahib", "")
	var schools []madhhabDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &schools); err != nil {
		t.Fatal(err)
	}
	if len(schools) != 4 {
		t.Errorf("expected 4 schools, got %d", len(schools))
	}
}
