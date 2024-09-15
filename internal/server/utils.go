package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, errorCode int, errorMessagesToLog []string, errorMessagesToRender []string) {
	data, _ := json.Marshal(map[string][]string{"errors": errorMessagesToRender})
	w.WriteHeader(errorCode)
	w.Write(data)

	for _, message := range errorMessagesToLog {
		log.Printf("ERROR %s\n", message)
	}
}

// Extracts API key value from request headers
func extractApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is not set")
	}
	// Strip redundant prefix
	apiKey, _ := strings.CutPrefix(authHeader, "ApiKey ")
	return apiKey, nil
}
