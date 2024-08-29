package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var expectedEnvVars = []string{"HTTP_SERVER_PORT", "POSTGRES_HOST", "POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD"}

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
