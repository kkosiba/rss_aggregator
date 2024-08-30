package main

func main() {
	// Check if expected env vars are set
	validateEnv()

	connection := connectToDatabase()
	connection.AutoMigrate(&User{})
	
	startServer()
}
