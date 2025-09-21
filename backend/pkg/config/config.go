// Package config provides configuration management for the FinanceTracker application.
// It loads configuration from environment variables with sensible defaults.
package config

import (
	"os"
)

// Config holds all configuration values for the application.
type Config struct {
	// App
	AppEnv  string
	AppPort string
	GinMode string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Auth0
	Auth0Domain       string
	Auth0ClientID     string
	Auth0ClientSecret string
	Auth0Audience     string

	// CORS
	CORSAllowedOrigins string
	CORSAllowedMethods string
	CORSAllowedHeaders string

	// Logging
	LogLevel  string
	LogFormat string

	// Atlas
	AtlasDatabaseURL string
}

// Load creates and returns a new Config instance with values from environment variables.
func Load() *Config {
	return &Config{
		// App
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "financetracker_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Auth0
		Auth0Domain:       getEnv("AUTH0_DOMAIN", ""),
		Auth0ClientID:     getEnv("AUTH0_CLIENT_ID", ""),
		Auth0ClientSecret: getEnv("AUTH0_CLIENT_SECRET", ""),
		Auth0Audience:     getEnv("AUTH0_AUDIENCE", ""),

		// CORS
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		CORSAllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		CORSAllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "debug"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		// Atlas
		AtlasDatabaseURL: getEnv("ATLAS_DATABASE_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDSN returns PostgreSQL connection string
func (c *Config) GetDSN() string {
	return "host=" + c.DBHost + " user=" + c.DBUser + " password=" + c.DBPassword +
		" dbname=" + c.DBName + " port=" + c.DBPort + " sslmode=" + c.DBSSLMode
}
