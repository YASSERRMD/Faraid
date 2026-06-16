package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/YASSERRMD/Faraid/internal/store"
)

// Server holds the HTTP handlers and their dependencies. Later phases add the
// trial LLM tier.
type Server struct {
	store store.Store
}

// NewServerWithStore returns a Server backed by the given store.
func NewServerWithStore(st store.Store) *Server {
	return &Server{store: st}
}

// NewServer returns a Server backed by an in-memory store.
func NewServer() *Server {
	return NewServerWithStore(store.NewMemory())
}

// Router returns the HTTP handler mounting every endpoint under /api/v1.
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", s.handleHealth)
		r.Get("/madhahib", s.handleMadhahib)
		r.Post("/solve", s.handleSolve)
		r.Post("/compare", s.handleCompare)

		r.Get("/cases", s.handleListCases)
		r.Post("/cases", s.handleCreateCase)
		r.Get("/cases/{id}", s.handleGetCase)
		r.Delete("/cases/{id}", s.handleDeleteCase)
	})
	return r
}
