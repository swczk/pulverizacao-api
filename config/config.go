package config

import (
	"log"
	"os"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	Port         string
}

func Load() *Config {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}

	return &Config{
		MongoURI:     mongoURI,
		DatabaseName: getEnv("DATABASE_NAME", "pulverizacao"),
		Port:         getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
