package server

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
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

type usersResource struct {
	database *database.Database
}

// Defines routes for /users namespace
func (rs usersResource) Routes() chi.Router {
	router := chi.NewRouter()

	router.Post("/", rs.Create)
	router.Get("/", rs.Get)

	return router
}

// Defines a handler for POST /users to create a new user
func (rs usersResource) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonBody struct{ Name string }
	err := decoder.Decode(&jsonBody)
	if err != nil {
		baseMessage := "Failed to decode JSON body"
		respondWithError(w, 400, []string{fmt.Sprintf("%s: Error: %s", baseMessage, err)}, []string{baseMessage})
		return
	}

	connection := rs.database.Connect()

	data := make([]byte, 10)
	if _, err := rand.Read(data); err == nil {
		apiKey := fmt.Sprintf("%x", sha256.Sum256(data))
		_, err := connection.Query(
			context.Background(),
			"INSERT INTO users (id, created_at, updated_at, name, api_key) values (gen_random_uuid(), $1, $2, $3, $4)",
			time.Now().UTC(), time.Now().UTC(), jsonBody.Name, apiKey,
		)
		if err != nil {
			baseMessage := fmt.Sprintf("Failed to create user '%s'", jsonBody.Name)
			respondWithError(w, 400, []string{fmt.Sprintf("%s: Error: %s", baseMessage, err)}, []string{baseMessage})
			return
		}
		msg := fmt.Sprintf("User '%s' created successfully", jsonBody.Name)
		respondWithJSON(w, 200, map[string]string{"details": msg})
		log.Print(msg)
	} else {
		baseMessage := "Failed to generate API key for user '%s'"
		respondWithError(w, 500, []string{fmt.Sprintf("%s. Error: %s", baseMessage, err)}, []string{baseMessage})
	}
}

// Defines a handler for GET /users to fetch user details.
// This endpoint requires Authorization header to be set.
func (rs usersResource) Get(w http.ResponseWriter, r *http.Request) {
	apiKey, err := extractApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, []string{string(err.Error())}, []string{err.Error()})
		return
	}
	connection := rs.database.Connect()

	var (
		ID        uuid.UUID
		CreatedAt time.Time
		UpdatedAt time.Time
		Name      string
		ApiKey    string
	)

	err = connection.QueryRow(
		context.Background(),
		"SELECT id, created_at, updated_at, name, api_key FROM users WHERE api_key = $1", apiKey,
	).Scan(&ID, &CreatedAt, &UpdatedAt, &Name, &ApiKey)

	if err != nil {
		baseMessage := "Failed to retrieve user"
		respondWithError(w, 400, []string{fmt.Sprintf("%s. Error: %s", baseMessage, err)}, []string{baseMessage})
		return
	}
	respondWithJSON(w, 200, &database.UserModel{ID: ID, CreatedAt: CreatedAt, UpdatedAt: UpdatedAt, Name: Name, ApiKey: ApiKey})
}
