package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kkosiba/rss_aggregator/internal/database"
)

type HTTPServer struct {
	port     string
	database *database.Database
}

func New() *http.Server {
	database := &database.Database{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_DB"),
	}
	database.Migrate()

	httpServer := &HTTPServer{
		port:     os.Getenv("HTTP_SERVER_PORT"),
		database: database,
	}
	return &http.Server{
		Handler: httpServer.RegisterRoutes(),
		Addr:    fmt.Sprintf(":%s", httpServer.port),
	}
}
