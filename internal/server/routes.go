package server

import (
	"github.com/go-chi/chi/v5"

	"github.com/ConflictHQ/boilerworks-go-micro/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-micro/internal/handler"
	"github.com/ConflictHQ/boilerworks-go-micro/internal/middleware"
)

func RegisterRoutes(r *chi.Mux, q *queries.Queries) {
	eventHandler := handler.NewEventHandler(q)
	apiKeyHandler := handler.NewApiKeyHandler(q)

	// Public routes
	r.Get("/health", handler.HealthCheck)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.ApiKeyAuth(q))

		// Events
		r.Route("/events", func(r chi.Router) {
			r.With(middleware.RequireScope("events.write")).Post("/", eventHandler.Create)
			r.With(middleware.RequireScope("events.read")).Get("/", eventHandler.List)
			r.With(middleware.RequireScope("events.read")).Get("/{id}", eventHandler.Get)
			r.With(middleware.RequireScope("events.write")).Delete("/{id}", eventHandler.Delete)
		})

		// API Keys
		r.Route("/api-keys", func(r chi.Router) {
			r.With(middleware.RequireScope("keys.manage")).Post("/", apiKeyHandler.Create)
			r.With(middleware.RequireScope("keys.manage")).Get("/", apiKeyHandler.List)
			r.With(middleware.RequireScope("keys.manage")).Delete("/{id}", apiKeyHandler.Revoke)
		})
	})
}
