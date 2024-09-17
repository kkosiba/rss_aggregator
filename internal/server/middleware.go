package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/kkosiba/rss_aggregator/internal/database"
)

// ApiKeyAuth implements a simple authentication middleware handler using API key
func ApiKeyAuth(db *database.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			apiKey, err := extractApiKey(r.Header)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, []string{string(err.Error())}, []string{err.Error()})
				return
			}
			var userId string
			connection := db.Connect()
			err = connection.QueryRow(
				context.Background(),
				"SELECT id FROM users WHERE api_key = $1",
				apiKey,
			).Scan(&userId)
			connection.Close()
			if err != nil {
				baseMessage := "Authorization failure"
				respondWithError(
					w,
					http.StatusUnauthorized,
					[]string{fmt.Sprintf("%s. Error: %s", baseMessage, err)},
					[]string{baseMessage},
				)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
