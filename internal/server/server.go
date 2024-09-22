package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/kkosiba/rss_aggregator/internal/database"
)

type HTTPServer struct {
	port     string
	database *database.Database
}

func New() *http.Server {
	database := &database.Database{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_DB"),
	}

	httpServer := &HTTPServer{
		port:     os.Getenv("HTTP_SERVER_PORT"),
		database: database,
	}
	return &http.Server{
		Handler: httpServer.RegisterRoutes(),
		Addr:    fmt.Sprintf(":%s", httpServer.port),
	}
}

func (server *HTTPServer) RegisterRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.CleanPath)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	// API v1
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/healthcheck", healthcheckResource{}.Routes())
		r.Mount("/users", usersResource{database: server.database}.Routes())
		// Add custom API key auth middleware to /feeds endpoints
		// todo: GET /feeds should not require authentication
		r.With(ApiKeyAuth(server.database)).Mount("/feeds", feedsResource{database: server.database}.Routes())
	})

	return router
}
