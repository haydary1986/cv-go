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

	// In production, require secrets to be explicitly set
	if ginMode == "release" {
		if jwtSecret == "" {
			log.Fatal("FATAL: JWT_SECRET environment variable is required in production")
		}
		if len(jwtSecret) < 32 {
			log.Fatal("FATAL: JWT_SECRET must be at least 32 characters")
		}
		if aesKey == "" {
			log.Fatal("FATAL: AES_KEY environment variable is required in production")
		}
		if len(aesKey) != 32 {
			log.Fatal("FATAL: AES_KEY must be exactly 32 characters for AES-256")
		}
	} else {
		// Development defaults with warning
		if jwtSecret == "" {
			jwtSecret = "dev-only-jwt-secret-do-not-use-in-production"
			fmt.Println("WARNING: Using default JWT_SECRET. Set JWT_SECRET env var for production.")
		}
		if aesKey == "" {
			aesKey = "devonlyaeskey0000devonlyaeskey00"
			fmt.Println("WARNING: Using default AES_KEY. Set AES_KEY env var for production.")
		}
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
