package main

import (
    "fmt"
    "log"
    "os"

    // Third-party libraries
    "github.com/joho/godotenv"
)

func main() {
    // Load env vars from .env file
    godotenv.Load(".env")

    portString := os.Getenv("PORT")
    if portString == "" {
        log.Fatal("Environment variable PORT is not set")
    }

    fmt.Println("PORT =", portString)
}
