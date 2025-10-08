package config

import (
	"os"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	GitHubClientID     string
	GitHubSecret       string
	GCPProjectID       string
	GCPBucketName      string
	GCPCredentialsPath string
	RedirectURL        string
}

func Load() *Config {
	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://user:password@localhost/devsync?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubSecret:       getEnv("GITHUB_CLIENT_SECRET", ""),
		GCPProjectID:       getEnv("GCP_PROJECT_ID", ""),
		GCPBucketName:      getEnv("GCP_BUCKET_NAME", ""),
		GCPCredentialsPath: getEnv("GCP_CREDENTIALS_PATH", ""),
		RedirectURL:        getEnv("REDIRECT_URL", "http://localhost:3000/auth/callback"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
