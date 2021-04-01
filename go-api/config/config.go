package config

import "os"

var (
	URI           string
	IS_DEBUG_MODE bool
	PORT          string
)

func init() {
	URI = getEnv("URI", "")
	IS_DEBUG_MODE = true
	PORT = getEnv("PORT", ":8080")
}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
