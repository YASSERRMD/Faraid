package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server holds the HTTP handlers and their dependencies. Later phases add the
// store and the trial LLM tier.
type Server struct{}

// NewServer returns a Server.
func NewServer() *Server {
	return &Server{}
}

// Router returns the HTTP handler mounting every endpoint under /api/v1.
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", s.handleHealth)
		r.Get("/madhahib", s.handleMadhahib)
		r.Post("/solve", s.handleSolve)
	})
	return r
}
