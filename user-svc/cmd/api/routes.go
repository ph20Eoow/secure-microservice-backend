package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	r := chi.NewRouter()

	// CORS Settings
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/user", func(r chi.Router) {
		r.With(app.isAuthByBasicAuth).Get("/{id}", app.getUserProfile)
		r.Put("/", app.createUser)
		r.With(app.isAuthByBasicAuth).Post("/password", app.updatePassword)
	})
	return r
}
