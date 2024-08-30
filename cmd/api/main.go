package main

import (
	"github.com/kkosiba/rss_aggregator/internal/database"
	"github.com/kkosiba/rss_aggregator/internal/server"
	"github.com/kkosiba/rss_aggregator/internal/utils"
)

func main() {
	// Check if expected env vars are set
	utils.ValidateEnv()

	connection := database.ConnectToDatabase()
	connection.AutoMigrate(&database.UserModel{})

	server.StartServer()
}
