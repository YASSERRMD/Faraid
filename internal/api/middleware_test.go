package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func TestSecureHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	secureHeaders(okHandler).ServeHTTP(w, r)

	want := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Content-Security-Policy": "default-src 'self'",
		"Referrer-Policy":        "same-origin",
		"Permissions-Policy":     "interest-cohort=()",
	}
	for header, val := range want {
		if got := w.Header().Get(header); got != val {
			t.Errorf("header %s: got %q, want %q", header, got, val)
		}
	}
}

func TestRequestSizeLimitAllows(t *testing.T) {
	body := strings.NewReader(`{"key":"value"}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", body)
	requestSizeLimit(512 * 1024)(okHandler).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequestSizeLimitRejectsOnRead(t *testing.T) {
	large := bytes.Repeat([]byte("x"), 600*1024)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(large))

	// The handler that actually reads the body is where the limit fires.
	handler := requestSizeLimit(512*1024)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 600*1024)
		n, err := r.Body.Read(buf)
		if err == nil && n == 600*1024 {
			// Should not reach here if the limit fired.
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusRequestEntityTooLarge)
	}))
	handler.ServeHTTP(w, r)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Code)
	}
}

func TestRateLimitAllowsUnderLimit(t *testing.T) {
	mw := rateLimit(5)
	for i := range 5 {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "10.0.0.1:1234"
		mw(okHandler).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimitBlocksOverLimit(t *testing.T) {
	mw := rateLimit(3)
	for range 3 {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "10.0.0.2:1234"
		mw(okHandler).ServeHTTP(w, r)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "10.0.0.2:1234"
	mw(okHandler).ServeHTTP(w, r)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}
}

func TestRateLimitIsolatesIPs(t *testing.T) {
	mw := rateLimit(2)
	// Exhaust IP A.
	for range 2 {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "10.0.0.3:1"
		mw(okHandler).ServeHTTP(w, r)
	}
	// IP B should still be allowed.
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "10.0.0.4:1"
	mw(okHandler).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for different IP, got %d", w.Code)
	}
}

func TestRateLimitReadsXForwardedFor(t *testing.T) {
	mw := rateLimit(1)
	make := func(xff string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "10.0.0.5:1"
		r.Header.Set("X-Forwarded-For", xff)
		mw(okHandler).ServeHTTP(w, r)
		return w
	}
	// First request from client A passes.
	if got := make("192.168.1.1, 10.0.0.5").Code; got != http.StatusOK {
		t.Fatalf("expected 200, got %d", got)
	}
	// Second request from same forwarded IP is blocked.
	if got := make("192.168.1.1, 10.0.0.5").Code; got != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", got)
	}
}

func TestRouterHasSecurityHeaders(t *testing.T) {
	srv := NewServer()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	srv.Router().ServeHTTP(w, r)
	if got := w.Header().Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("X-Frame-Options: got %q, want DENY", got)
	}
}
