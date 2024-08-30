package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kkosiba/rss_aggregator/db"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		msg := fmt.Sprintf("Method %s is not allowed", r.Method)
		respondWithJSON(w, 405, map[string]string{"error": msg})
		log.Print(msg)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var jsonBody struct{Name string}
	err := decoder.Decode(&jsonBody)
	if err != nil {
		msg := "Failed to decode JSON body"
		respondWithJSON(w, 400, map[string]string{"error": msg})
		log.Print(msg)
		return
	}

	connection := db.ConnectToDatabase()
	result := connection.Create(
		&db.UserModel{CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: jsonBody.Name},
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
