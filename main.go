package main

import (
	"github.com/kkosiba/rss_aggregator/api"
	"github.com/kkosiba/rss_aggregator/db"
	"github.com/kkosiba/rss_aggregator/utils"
)

func main() {
	// Check if expected env vars are set
	utils.ValidateEnv()

	connection := db.ConnectToDatabase()
	connection.AutoMigrate(&db.UserModel{})

	api.StartServer()
}
