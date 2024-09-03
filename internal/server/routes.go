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
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kkosiba/rss_aggregator/internal/database"
)

func (server *HTTPServer) RegisterRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/healthcheck", server.healthCheck)
	router.Post("/users", server.createUser)
	router.Get("/users", server.getUser)
	return router
}

func (server *HTTPServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	// Could check something useful here, but it's good enough for now
	respondWithJSON(w, 200, struct{}{})
}

func (server *HTTPServer) getUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		message := "Authorization header is not set"
		respondWithError(w, 401, []string{message}, []string{message})
		return
	}

	// Strip redundant prefix
	apiKey, _ := strings.CutPrefix(authHeader, "ApiKey ")

	dbpool := server.database.Connect()

	var user database.UserModel
	err := dbpool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE api_key = $1", apiKey,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.ApiKey)

	if err != nil {
		baseMessage := "Failed to retrieve user"
		respondWithError(w, 400, []string{fmt.Sprintf("%s. Error: %s", baseMessage, err)}, []string{baseMessage})
		return
	}
	respondWithJSON(w, 200, &user)
}

func (server *HTTPServer) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonBody struct{ Name string }
	err := decoder.Decode(&jsonBody)
	if err != nil {
		baseMessage := "Failed to decode JSON body"
		respondWithError(w, 400, []string{fmt.Sprintf("%s: Error: %s", baseMessage, err)}, []string{baseMessage})
		return
	}

	connection := server.database.Connect()
	defer connection.Close()

	data := make([]byte, 10)
	if _, err := rand.Read(data); err == nil {
		apiKey := fmt.Sprintf("%x", sha256.Sum256(data))
		_, err := connection.Query(
			context.Background(),
			"INSERT INTO users (created_at, updated_at, name, api_key) values ($1, $2, $3)",
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
		return
	} else {
		baseMessage := "Failed to generate API key for user '%s'"
		respondWithError(w, 500, []string{fmt.Sprintf("%s. Error: %s", baseMessage, err)}, []string{baseMessage})
		return
	}
}
