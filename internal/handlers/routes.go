package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	//todo: add logging middleware
	r.Get("/v1/healthcheck", app.healthcheckHandler)
	r.Post("/v1/movies", app.createMovieHandler)
	r.Get("/v1/movies/{id}", app.showMovieHandler)
	r.Put("/v1/movies/{id}", app.updateMovieHandler)
	r.Delete("/v1/movies/{id}", app.deleteMovieHandler)
	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)
	return r
}
