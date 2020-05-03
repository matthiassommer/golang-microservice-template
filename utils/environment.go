package utils

import (
	"os"
	"strings"
)

// Keywords for Environment identification
const (
	ENV_PROD  = "production" // keyword for production environment
	ENV_STAGE = "staging"    // keyword for staging environment
	ENV_DEV   = "develop"    // keyword for develop environment
	ENV_LOCAL = "local"      // keyword for local environment
)

// Environment returns the current environment according to environment variable ENV
func Environment() string {
	value, _ := os.LookupEnv("ENV")

	switch strings.ToLower(value) {
	case ENV_PROD, "prod":
		return ENV_PROD
	case ENV_STAGE, "stage":
		return ENV_STAGE
	case ENV_DEV, "dev":
		return ENV_DEV
	default:
		return ENV_LOCAL
	}
}

func IsProduction() bool {
	return Environment() == ENV_PROD
}
