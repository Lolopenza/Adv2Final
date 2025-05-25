package config

import (
	"os"
	"strings"
)

// GetEnv retrieves an environment variable or returns a default value if not found
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Handle case where values are in one line with spaces
	if strings.Contains(key, "SMTP_") && strings.Contains(value, " ") {
		// Parse values that might be combined
		parts := strings.Fields(value)
		if len(parts) > 0 {
			return parts[0]
		}
	}

	return value
}
