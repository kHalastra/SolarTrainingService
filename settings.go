package main

import (
	"os"
)

// Config struct to hold the database credentials
type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
}

// LoadConfig loads database configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "solar"),
		DBHost:     getEnv("DB_HOST", "localhost"),
	}

	return config
}

// Helper function to read environment variables with a fallback default value
func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
