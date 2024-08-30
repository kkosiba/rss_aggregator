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
	database := database.New()
	database.Migrate()
	
	httpServer := &HTTPServer{
		port: os.Getenv("HTTP_SERVER_PORT"),
		database: &database,
	}
	return &http.Server{
		Handler: httpServer.RegisterRoutes(),
		Addr:    fmt.Sprintf(":%s", httpServer.port),
	}
}
