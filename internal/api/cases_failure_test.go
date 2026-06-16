package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/store"
)

// failStore returns an error from every method.
type failStore struct{}

func (failStore) SaveCase(context.Context, string, json.RawMessage, json.RawMessage) (store.SavedCase, error) {
	return store.SavedCase{}, errors.New("boom")
}
func (failStore) GetCase(context.Context, string) (store.SavedCase, error) {
	return store.SavedCase{}, errors.New("boom")
}
func (failStore) ListCases(context.Context) ([]store.SavedCase, error) {
	return nil, errors.New("boom")
}
func (failStore) DeleteCase(context.Context, string) error { return errors.New("boom") }

// corruptStore succeeds but returns cases whose stored input is not valid JSON.
type corruptStore struct{}

func (corruptStore) SaveCase(context.Context, string, json.RawMessage, json.RawMessage) (store.SavedCase, error) {
	return store.SavedCase{ID: "x", Input: json.RawMessage("{bad")}, nil
}
func (corruptStore) GetCase(context.Context, string) (store.SavedCase, error) {
	return store.SavedCase{ID: "x", Input: json.RawMessage("{bad")}, nil
}
func (corruptStore) ListCases(context.Context) ([]store.SavedCase, error) {
	return []store.SavedCase{{ID: "x", Input: json.RawMessage("{bad")}}, nil
}
func (corruptStore) DeleteCase(context.Context, string) error { return nil }

func statusFor(srv *Server, method, path, body string) int {
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, newJSONRequest(method, path, body))
	return rec.Code
}

func TestCaseEndpointStoreErrors(t *testing.T) {
	srv := NewServerWithStore(failStore{})
	checks := []struct {
		method, path, body string
	}{
		{http.MethodPost, "/api/v1/cases", sampleCaseBody},
		{http.MethodGet, "/api/v1/cases", ""},
		{http.MethodGet, "/api/v1/cases/x", ""},
		{http.MethodDelete, "/api/v1/cases/x", ""},
	}
	for _, c := range checks {
		if got := statusFor(srv, c.method, c.path, c.body); got != http.StatusInternalServerError {
			t.Errorf("%s %s = %d, want 500", c.method, c.path, got)
		}
	}
}

func TestCaseEndpointEncodeErrors(t *testing.T) {
	srv := NewServerWithStore(corruptStore{})
	for _, c := range []struct{ method, path, body string }{
		{http.MethodPost, "/api/v1/cases", sampleCaseBody},
		{http.MethodGet, "/api/v1/cases", ""},
		{http.MethodGet, "/api/v1/cases/x", ""},
	} {
		if got := statusFor(srv, c.method, c.path, c.body); got != http.StatusInternalServerError {
			t.Errorf("%s %s = %d, want 500", c.method, c.path, got)
		}
	}
}
