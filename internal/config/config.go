package config

import (
	"log"
	"os"
)

type Config struct {
	ProjectID  string
	ServerAddr string
	Timezone   string
	JWTSecret  string
	ENV        string
}

func Load() *Config {
	return &Config{
		ProjectID:  mustGetEnv("GOOGLE_PROJECT_ID"),
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		Timezone:   getEnv("NITELOG_TIMEZONE", "America/Sao_Paulo"),
		JWTSecret:  mustGetEnv("JWT_SECRET"),
		ENV:        getEnv("NITELOG_ENV", "PRODUCTION"),
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return value
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
