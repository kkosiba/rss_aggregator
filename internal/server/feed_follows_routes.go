package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	decoder := json.NewDecoder(r.Body)
	var jsonBody struct {
		FeedId uuid.UUID
	}
	err := decoder.Decode(&jsonBody)
	if err != nil {
		baseMessage := "Failed to decode JSON body"
		respondWithError(w, http.StatusBadRequest, []string{fmt.Sprintf("%s: Error: %s", baseMessage, err)}, []string{baseMessage})
		return
	}
	apiKey, err := extractApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, []string{string(err.Error())}, []string{err.Error()})
		return
	}

	connection := rs.database.Connect()

	var userId uuid.UUID
	err = connection.QueryRow(
		context.Background(),
		"SELECT id FROM users WHERE api_key = $1",
		apiKey,
	).Scan(&userId)
	if err != nil {
		baseMessage := fmt.Sprintf("Failed to retrieve user ID")
		respondWithError(
			w,
			http.StatusBadRequest,
			[]string{fmt.Sprintf("%s: Error: %s", baseMessage, err)},
			[]string{baseMessage},
		)
		return
	}

	ID := uuid.New()
	CreatedAt := time.Now().UTC()
	UpdatedAt := time.Now().UTC()

	_, err = connection.Query(
		context.Background(),
		"INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id) values ($1, $2, $3, $4, $5)",
		ID, CreatedAt, UpdatedAt, jsonBody.FeedId, userId,
	)
	if err != nil {
		baseMessage := "Failed to create feed follow."
		respondWithError(
			w,
			http.StatusBadRequest,
			[]string{fmt.Sprintf("%s: Error: %s", baseMessage, err)},
			[]string{baseMessage},
		)
		return
	}
	respondWithJSON(
		w,
		http.StatusOK,
		database.FeedFollowModel{
			ID:        ID,
			CreatedAt: CreatedAt,
			UpdatedAt: UpdatedAt,
			FeedId:    jsonBody.FeedId,
			UserId:    userId,
		},
	)
	log.Printf("User '%s' follows feed '%s'. Feed follow ID: '%s'", userId, jsonBody.FeedId, ID)
}

// Defines a handler for DELETE /feed_follows/{feedId}
func (rs feedFollowsResource) Delete(w http.ResponseWriter, r *http.Request) {
	// todo: implement me!
}

// Defines a handler for GET /feed_follows
func (rs feedFollowsResource) Get(w http.ResponseWriter, r *http.Request) {
	// todo: implement me!
}
