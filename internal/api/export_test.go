package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

const exportBody = `{"deceasedSex":"female","heirs":{"husband":1,"son":1,"daughter":1},"madhhab":"Hanafi"}`

func exportRec(t *testing.T, format, body string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	NewServer().Router().ServeHTTP(rec, newJSONRequest(http.MethodPost, "/api/v1/export?format="+format, body))
	return rec
}

func TestExportJSON(t *testing.T) {
	rec := exportRec(t, "json", exportBody)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if cd := rec.Header().Get("Content-Disposition"); cd == "" {
		t.Error("expected a Content-Disposition header")
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"madhhab":"Hanafi"`)) {
		t.Errorf("json export missing fields: %s", rec.Body.String())
	}
}

func TestExportPDFStructure(t *testing.T) {
	rec := exportRec(t, "pdf", exportBody)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Errorf("content type = %q, want application/pdf", ct)
	}
	body := rec.Body.Bytes()
	if !bytes.HasPrefix(body, []byte("%PDF-")) {
		t.Error("output is not a PDF")
	}
	if !bytes.Contains(body, []byte("%%EOF")) {
		t.Error("PDF is not terminated")
	}
	// Compression is off, so the labels are visible in the bytes.
	for _, want := range []string{"Faraid Inheritance Result", "Hanafi", "husband", "Derivation"} {
		if !bytes.Contains(body, []byte(want)) {
			t.Errorf("PDF missing %q", want)
		}
	}
}

func TestExportPDFSpecialCase(t *testing.T) {
	// A gharrawayn case puts a special-case line in the PDF.
	rec := exportRec(t, "pdf", `{"deceasedSex":"female","heirs":{"husband":1,"father":1,"mother":1},"madhhab":"Hanafi"}`)
	if !bytes.Contains(rec.Body.Bytes(), []byte("Special case: gharrawayn")) {
		t.Error("PDF should note the gharrawayn special case")
	}
}

func TestExportErrors(t *testing.T) {
	srv := NewServer()
	if got := statusFor(srv, http.MethodPost, "/api/v1/export?format=xml", exportBody); got != http.StatusBadRequest {
		t.Errorf("unknown format = %d, want 400", got)
	}
	if got := statusFor(srv, http.MethodPost, "/api/v1/export?format=json", `{bad`); got != http.StatusBadRequest {
		t.Errorf("bad json = %d, want 400", got)
	}
	if got := statusFor(srv, http.MethodPost, "/api/v1/export?format=pdf", `{"deceasedSex":"male","heirs":{"husband":1},"madhhab":"Hanafi"}`); got != http.StatusBadRequest {
		t.Errorf("invalid case = %d, want 400", got)
	}
}
