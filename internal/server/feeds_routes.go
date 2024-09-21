package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kkosiba/rss_aggregator/internal/database"
)

type feedsResource struct {
	database *database.Database
}

// Defines routes for /feeds namespace
func (rs feedsResource) Routes() chi.Router {
	router := chi.NewRouter()
	router.Get("/", rs.Get)
	router.Post("/", rs.Create)
	return router
}

// Defines a handler for GET /feeds
func (rs feedsResource) Get(w http.ResponseWriter, r *http.Request) {
	connection := rs.database.Connect()

	var feeds []database.FeedModel

	rows, err := connection.Query(
		context.Background(),
		"SELECT id, created_at, updated_at, name, url, user_id FROM feeds",
	)
	if err != nil {
		baseMessage := "Failed to retrieve feeds"
		respondWithError(
			w,
			http.StatusBadRequest,
			[]string{fmt.Sprintf("%s: Error: %s", baseMessage, err)},
			[]string{baseMessage},
		)
		return
	}
	defer rows.Close()

	var errors []string
	for rows.Next() {
		var feed database.FeedModel
		err = rows.Scan(
			&feed.ID,
			&feed.CreatedAt,
			&feed.UpdatedAt,
			&feed.Name,
			&feed.Url,
			&feed.UserId,
		)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		} else {
			feeds = append(feeds, feed)
		}
	}
	if len(errors) > 0 {
		baseMessage := "Failed to retrieve feeds"
		respondWithError(
			w,
			http.StatusBadRequest,
			[]string{fmt.Sprintf("%s: Error: %s", baseMessage, strings.Join(errors, "\n"))},
			[]string{baseMessage},
		)
		return
	}

	respondWithJSON(w, http.StatusOK, &feeds)
}

// Defines a handler for POST /feeds
func (rs feedsResource) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonBody struct {
		Name string
		Url  string
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
		"INSERT INTO feeds (id, created_at, updated_at, name, url, user_id) values ($1, $2, $3, $4, $5, $6)",
		ID, CreatedAt, UpdatedAt, jsonBody.Name, jsonBody.Url, userId,
	)
	if err != nil {
		baseMessage := fmt.Sprintf("Failed to create feed '%s' with URL '%s'", jsonBody.Name, jsonBody.Url)
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
		database.FeedModel{
			ID:        ID,
			CreatedAt: CreatedAt,
			UpdatedAt: UpdatedAt,
			Name:      jsonBody.Name,
			Url:       jsonBody.Url,
			UserId:    userId,
		},
	)
	log.Printf("Feed '%s' created successfully for user ID '%s'", jsonBody.Name, userId)
}
