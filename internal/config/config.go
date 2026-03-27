package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	APIKeySeed  string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/boilerworks?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		APIKeySeed:  getEnv("API_KEY_SEED", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
