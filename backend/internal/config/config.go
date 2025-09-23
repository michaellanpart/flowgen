package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port         string
	Environment  string
	DatabaseURL  string
	DiagramsPath string
	JiraBaseURL  string
	JiraUsername string
	JiraAPIToken string
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "3001"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		DiagramsPath: getEnv("DIAGRAMS_PATH", "./diagrams"),
		JiraBaseURL:  getEnv("JIRA_BASE_URL", ""),
		JiraUsername: getEnv("JIRA_USERNAME", ""),
		JiraAPIToken: getEnv("JIRA_API_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
