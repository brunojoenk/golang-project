package config

import (
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
}

func New() *Config {
	return &Config{
		ServerPort:  getEnv("PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
