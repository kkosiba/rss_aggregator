package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var expectedEnvVars = []string{
	// HTTP Server
	"HTTP_SERVER_PORT",

	// Database
	"POSTGRES_HOST",
	"POSTGRES_PORT",
	"POSTGRES_DB",
	"POSTGRES_USER",
	"POSTGRES_PASSWORD",
}

func validateEnv() {
	var errors []string
	for _, envVar := range expectedEnvVars {
		_, exists := os.LookupEnv(envVar)
		if !exists {
			errors = append(errors, envVar)
		}
	}
	if errors != nil {
		msg := fmt.Sprintf("Missing env vars: %s", strings.Join(errors, ","))
		log.Fatal(msg)
	}
}
