package api

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/YASSERRMD/Faraid/internal/llm"
	"github.com/YASSERRMD/Faraid/internal/store"
)

// Server holds the HTTP handlers and their dependencies.
type Server struct {
	store store.Store
	// completer powers the trial explanation tier. It is nil when the trial
	// tier is disabled, in which case /explain returns 404.
	completer llm.Completer
}

// NewServerWithStore returns a Server backed by the given store.
func NewServerWithStore(st store.Store) *Server {
	return &Server{store: st}
}

// NewServer returns a Server backed by an in-memory store.
func NewServer() *Server {
	return NewServerWithStore(store.NewMemory())
}

// WithLLM enables the trial explanation tier with the given completer. Passing
// nil leaves the tier disabled.
func (s *Server) WithLLM(c llm.Completer) *Server {
	s.completer = c
	return s
}

// Router returns the HTTP handler mounting every endpoint under /api/v1.
// Requests are subject to security headers, a 512 KiB body size limit, and a
// per-IP rate limit of 120 requests per minute.
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(secureHeaders)
	r.Use(requestSizeLimit(512 * 1024))
	r.Use(rateLimit(120))
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", s.handleHealth)
		r.Get("/readyz", s.handleReadyz)
		r.Get("/madhahib", s.handleMadhahib)
		r.Post("/solve", s.handleSolve)
		r.Post("/compare", s.handleCompare)

		r.Get("/cases", s.handleListCases)
		r.Post("/cases", s.handleCreateCase)
		r.Get("/cases/{id}", s.handleGetCase)
		r.Delete("/cases/{id}", s.handleDeleteCase)

		r.Post("/export", s.handleExport)
		r.Post("/parse", s.handleParse)
		r.Post("/explain", s.handleExplain)
	})
	// Expose runtime counters for monitoring. This endpoint is intentionally
	// not under /api/v1 so it can be restricted to an internal network in
	// production without affecting the API prefix.
	r.Get("/debug/vars", expvar.Handler().ServeHTTP)
	return r
}
