package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	return router
}

func (server *HTTPServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	// Could check something useful here, but it's good enough for now
	respondWithJSON(w, 200, struct{}{})
}

func (server *HTTPServer) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonBody struct{ Name string }
	err := decoder.Decode(&jsonBody)
	if err != nil {
		msg := "Failed to decode JSON body"
		respondWithJSON(w, 400, map[string]string{"error": msg})
		log.Print(msg)
		return
	}

	connection, err := server.database.Connect()
	if err != nil {
		// We could consider app panic here if db connection can't be established
		log.Print(err)
		return
	}
	result := connection.Create(
		&database.UserModel{CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: jsonBody.Name},
	)
	if result.Error != nil {
		msg := fmt.Sprintf("Failed to create user %s", jsonBody.Name)
		respondWithJSON(w, 400, map[string]string{"error": msg})
		log.Print(msg)
		return
	}
	msg := fmt.Sprintf("User %s created successfully", jsonBody.Name)
	respondWithJSON(w, 200, map[string]string{"details": msg})
	log.Print(msg)
}
