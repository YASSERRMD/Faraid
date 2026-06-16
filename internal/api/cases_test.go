package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const sampleCaseBody = `{"name":"family A","input":{"deceasedSex":"female","heirs":{"husband":1,"son":1,"daughter":1},"madhhab":"Hanafi"}}`

func TestCreateListGetDeleteCase(t *testing.T) {
	// Create, validated against the contract.
	rec := contractCheck(t, http.MethodPost, "/api/v1/cases", sampleCaseBody)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body %s", rec.Code, rec.Body.String())
	}
	var created savedCaseDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Name != "family A" || created.Input.Madhhab != "Hanafi" {
		t.Fatalf("unexpected created case: %+v", created)
	}

	// A single shared server keeps state across the remaining calls.
	srv := NewServer()
	mk := func(method, path, body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		srv.Router().ServeHTTP(rec, newJSONRequest(method, path, body))
		return rec
	}

	post := mk(http.MethodPost, "/api/v1/cases", sampleCaseBody)
	var saved savedCaseDTO
	_ = json.Unmarshal(post.Body.Bytes(), &saved)

	if got := mk(http.MethodGet, "/api/v1/cases/"+saved.ID, ""); got.Code != http.StatusOK {
		t.Errorf("get status = %d", got.Code)
	}

	list := mk(http.MethodGet, "/api/v1/cases", "")
	var cases []savedCaseDTO
	if err := json.Unmarshal(list.Body.Bytes(), &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) != 1 {
		t.Errorf("list = %d, want 1", len(cases))
	}

	if del := mk(http.MethodDelete, "/api/v1/cases/"+saved.ID, ""); del.Code != http.StatusNoContent {
		t.Errorf("delete status = %d", del.Code)
	}
	if got := mk(http.MethodGet, "/api/v1/cases/"+saved.ID, ""); got.Code != http.StatusNotFound {
		t.Errorf("get after delete status = %d, want 404", got.Code)
	}
}

func TestCaseEndpointErrors(t *testing.T) {
	srv := NewServer()
	mk := func(method, path, body string) int {
		rec := httptest.NewRecorder()
		srv.Router().ServeHTTP(rec, newJSONRequest(method, path, body))
		return rec.Code
	}

	if c := mk(http.MethodPost, "/api/v1/cases", `{not json`); c != http.StatusBadRequest {
		t.Errorf("bad json = %d", c)
	}
	if c := mk(http.MethodPost, "/api/v1/cases", `{"input":{"deceasedSex":"male","heirs":{"son":1},"madhhab":"Hanafi"}}`); c != http.StatusBadRequest {
		t.Errorf("missing name = %d, want 400", c)
	}
	if c := mk(http.MethodPost, "/api/v1/cases", `{"name":"x","input":{"deceasedSex":"male","heirs":{"husband":1},"madhhab":"Hanafi"}}`); c != http.StatusBadRequest {
		t.Errorf("invalid case = %d, want 400", c)
	}
	if c := mk(http.MethodGet, "/api/v1/cases/missing", ""); c != http.StatusNotFound {
		t.Errorf("missing get = %d, want 404", c)
	}
	if c := mk(http.MethodDelete, "/api/v1/cases/missing", ""); c != http.StatusNotFound {
		t.Errorf("missing delete = %d, want 404", c)
	}
}

func TestListCasesContract(t *testing.T) {
	rec := contractCheck(t, http.MethodGet, "/api/v1/cases", "")
	if rec.Code != http.StatusOK {
		t.Errorf("list status = %d", rec.Code)
	}
}
