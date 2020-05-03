package utils

import (
	"os"
)

// DefaultOrEnv Checks if the given environment variable key exists and returns its value.
// If the key was not set, the default value is returned.
func DefaultOrEnv(defaultValue, envKey string) string {
	if envVar, exists := os.LookupEnv(envKey); exists {
		return envVar
	}
	Log.Infof("environment variable %s not found, using default value %s", envKey, defaultValue)
	return defaultValue
}
