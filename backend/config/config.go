package config

import (
	"os"
)

type Config struct {
	DBPath        string
	JWTSecret     string
	Port          string
	AESKey        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	FrontendURL   string
}

func Load() *Config {
	return &Config{
		DBPath:             getEnv("DB_PATH", "cvbuilder.db"),
		JWTSecret:          getEnv("JWT_SECRET", "cv-builder-secret-key-change-in-production-2024"),
		Port:               getEnv("PORT", "8080"),
		AESKey:             getEnv("AES_KEY", "0123456789abcdef0123456789abcdef"), // 32 bytes for AES-256
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
