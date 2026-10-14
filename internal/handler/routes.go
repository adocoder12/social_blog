package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) SetupRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.ClientIPFromRemoteAddr)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(8 * time.Second))

	// JSON errors instead of chi's plain-text defaults
	mux.NotFound(app.notFound)
	mux.MethodNotAllowed(app.methodNotAllowed)

	mux.Route("/api/v1", func(api chi.Router) {
		api.Get("/health", app.HealthCheckHandler)

		api.Route("/posts", func(posts chi.Router) {
			posts.Get("/", app.GetAllPostsHandler)
			posts.Post("/", app.CreatePostHandler)
		})
	})

	return mux
}