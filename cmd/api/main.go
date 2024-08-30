package main

import (
	"log"
	"strings"

	"github.com/kkosiba/rss_aggregator/internal/server"
	"github.com/kkosiba/rss_aggregator/internal/utils"
)

func main() {
	// Check if expected env vars are set
	utils.ValidateEnv()

	server := server.New()
	log.Printf("Starting an HTTP server on port %v", strings.Split(server.Addr, ":")[1])
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
