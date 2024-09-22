package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kkosiba/rss_aggregator/internal/database"
)

type feedFollowsResource struct {
	database *database.Database
}

// Defines routes for /feed_follows namespace
func (rs feedFollowsResource) Routes() chi.Router {
	router := chi.NewRouter()
	router.Get("/", rs.Create)
	router.Post("/", rs.Create)
	router.Delete("/{feedId}", rs.Delete)
	return router
}

// Defines a handler for POST /feed_follows
func (rs feedFollowsResource) Create(w http.ResponseWriter, r *http.Request) {
	// todo: implement me!
}

// Defines a handler for DELETE /feed_follows/{feedId}
func (rs feedFollowsResource) Delete(w http.ResponseWriter, r *http.Request) {
	// todo: implement me!
}

// Defines a handler for GET /feed_follows
func (rs feedFollowsResource) Get(w http.ResponseWriter, r *http.Request) {
	// todo: implement me!
}
