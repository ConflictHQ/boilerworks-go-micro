package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/ConflictHQ/boilerworks-go-micro/internal/database/queries"
	mw "github.com/ConflictHQ/boilerworks-go-micro/internal/middleware"
)

func New(q *queries.Queries) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	// Rate limiting: 60 requests per minute per IP
	limiter := mw.NewRateLimiter(60, time.Minute)
	r.Use(limiter.Handler)

	// Register routes
	RegisterRoutes(r, q)

	return r
}
