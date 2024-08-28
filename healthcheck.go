package main

import "net/http"

func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}
