package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YASSERRMD/Faraid/internal/llm"
)

type fakeCompleter struct {
	text string
	err  error
}

func (f fakeCompleter) Complete(context.Context, llm.Request) (llm.Response, error) {
	return llm.Response{Text: f.text}, f.err
}

const explainBody = `{"deceasedSex":"female","heirs":{"husband":1,"son":1,"daughter":1},"madhhab":"Hanafi"}`

func TestExplainDisabledReturns404(t *testing.T) {
	rec := httptest.NewRecorder()
	NewServer().Router().ServeHTTP(rec, newJSONRequest(http.MethodPost, "/api/v1/explain", explainBody))
	if rec.Code != http.StatusNotFound {
		t.Errorf("disabled explain = %d, want 404", rec.Code)
	}
}

func TestExplainConsistent(t *testing.T) {
	prose := "The husband takes 1/4 of the estate, the son 1/2, and the daughter 1/4."
	handler := NewServer().WithLLM(fakeCompleter{text: prose}).Router()
	rec := contractCheckH(t, handler, http.MethodPost, "/api/v1/explain", explainBody)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body.String())
	}
	var exp explanationDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &exp); err != nil {
		t.Fatal(err)
	}
	if !exp.Consistent || !exp.Experimental || exp.Text != prose {
		t.Errorf("expected consistent prose, got %+v", exp)
	}
}

func TestExplainDriftReturnsRawDerivation(t *testing.T) {
	// The prose invents 5/9, which the engine never produced.
	prose := "The husband takes 5/9 of everything."
	srv := NewServer().WithLLM(fakeCompleter{text: prose})
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, newJSONRequest(http.MethodPost, "/api/v1/explain", explainBody))
	var exp explanationDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &exp); err != nil {
		t.Fatal(err)
	}
	if exp.Consistent {
		t.Error("drifting prose should be flagged not consistent")
	}
	if exp.Text == prose {
		t.Error("on drift, the raw derivation should be returned, not the prose")
	}
}

func TestExplainProviderError(t *testing.T) {
	srv := NewServer().WithLLM(fakeCompleter{err: context.DeadlineExceeded})
	if got := statusFor(srv, http.MethodPost, "/api/v1/explain", explainBody); got != http.StatusBadGateway {
		t.Errorf("provider error = %d, want 502", got)
	}
}

func TestExplainBadRequest(t *testing.T) {
	srv := NewServer().WithLLM(fakeCompleter{text: "ok"})
	if got := statusFor(srv, http.MethodPost, "/api/v1/explain", `{bad`); got != http.StatusBadRequest {
		t.Errorf("bad json = %d, want 400", got)
	}
	if got := statusFor(srv, http.MethodPost, "/api/v1/explain", `{"deceasedSex":"male","heirs":{"husband":1},"madhhab":"Hanafi"}`); got != http.StatusBadRequest {
		t.Errorf("invalid case = %d, want 400", got)
	}
}
