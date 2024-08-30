package api

import (
	"log"
	"net/http"
	"os"


	// Third-party libraries
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func StartServer() {
	httpServerPort := os.Getenv("HTTP_SERVER_PORT")

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET"},
	}))

	v1router := chi.NewRouter()
	v1router.HandleFunc("/healthcheck", healthCheck)
	router.Mount("/v1", v1router)

	log.Printf("Starting an HTTP server on port %v", httpServerPort)
	server := &http.Server{
		Handler: router,
		Addr:    ":" + httpServerPort,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
