package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port string
	DatabaseURL string
	JWTSecret string
}

func Load() (Config, error) {
	requiredVars := []string{
		"DATABASE_URL",
		"JWT_SECRET",
	}
	
	var missing []string
	
	for _, key := range requiredVars {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	
	if len(missing) > 0 {
		return Config{}, fmt.Errorf(
			"missing required environment variables: %s",
			strings.Join(missing, ", "),
		)
	}
	
	return Config{
		Port: getEnv("PORT", "4849"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	
	return fallback
}