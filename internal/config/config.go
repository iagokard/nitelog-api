package config

import (
	"log"
	"os"
)

type Config struct {
	MongoURI   string
	DBName     string
	ServerAddr string
	Timezone   string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		MongoURI:   mustGetEnv("MONGO_URI"),
		DBName:     mustGetEnv("DB_NAME"),
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		Timezone:   getEnv("NITELOG_TIMEZONE", "America/Sao_Paulo"),
		JWTSecret:  mustGetEnv("JWT_SECRET"),
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
