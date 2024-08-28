package main

import (
	"log"
	"net/http"
	"os"

	// Third-party libraries
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load env vars from .env file
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Environment variable PORT is not set")
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET"},
	}))

	v1router := chi.NewRouter()
	v1router.HandleFunc("/healthcheck", healthCheck)
	router.Mount("/v1", v1router)

	log.Printf("Starting an HTTP server on port %v", portString)
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
