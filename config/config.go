package config

import (
	"os"
)

type Config struct {
	ServerPort string
	Postgres   PostgresConfig
}

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
}

func New() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "3000"),
		Postgres: PostgresConfig{
			PostgresqlHost:     getEnv("POSTGRES_HOST", "localhost"),
			PostgresqlPort:     getEnv("POSTGRES_PORT", "5432"),
			PostgresqlUser:     getEnv("POSTGRES_USER", "postgres"),
			PostgresqlPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
			PostgresqlDbname:   getEnv("POSTGRES_DBNAME", "postgres"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
