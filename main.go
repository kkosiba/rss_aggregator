package main

func main() {
	// Check if expected env vars are set
	validateEnv()

	connectToDatabase()
	startServer()
}
