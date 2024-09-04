package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type healthcheckResource struct{}

// Defines routes for /healthcheck namespace
func (rs healthcheckResource) Routes() chi.Router {
	router := chi.NewRouter()
	router.Get("/", rs.Get)
	return router
}

// Defines a handler for GET /healthcheck
func (rs healthcheckResource) Get(w http.ResponseWriter, r *http.Request) {
	// Could check something useful here, but it's good enough for now
	respondWithJSON(w, 200, struct{}{})
}
