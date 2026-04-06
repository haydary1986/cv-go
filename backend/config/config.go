package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	DBPath             string
	JWTSecret          string
	Port               string
	AESKey             string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	FrontendURL        string
}

func Load() *Config {
	ginMode := getEnv("GIN_MODE", "debug")

	jwtSecret := getEnv("JWT_SECRET", "")
	aesKey := getEnv("AES_KEY", "")

	// Validate secrets
	if jwtSecret == "" {
		if ginMode == "release" {
			log.Println("WARNING: JWT_SECRET not set in production! Generating temporary secret.")
		}
		jwtSecret = "auto-generated-change-me-" + fmt.Sprintf("%d", os.Getpid())
		log.Println("WARNING: Using auto-generated JWT_SECRET. Set JWT_SECRET env var for production.")
	}
	if aesKey == "" || len(aesKey) != 32 {
		if ginMode == "release" {
			log.Println("WARNING: AES_KEY not set or invalid length in production!")
		}
		aesKey = "autogen0key00000autogen0key00000"
		log.Println("WARNING: Using fallback AES_KEY. Set AES_KEY (32 chars) env var for production.")
	}

	return &Config{
		DBPath:             getEnv("DB_PATH", "cvbuilder.db"),
		JWTSecret:          jwtSecret,
		Port:               getEnv("PORT", "8080"),
		AESKey:             aesKey,
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
