package di

import (
	"os"
	"strconv"
	"time"
)

// Config は全てのアプリケーション設定を保持する
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth0    Auth0Config
	JWT      JWTConfig
	Log      LogConfig
}

// ServerConfig はサーバー関連の設定
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig はデータベース関連の設定
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Auth0Config はAuth0認証関連の設定
type Auth0Config struct {
	Domain       string
	ClientID     string
	ClientSecret string
	CallbackURL  string
	Audience     string
}

// JWTConfig はJWT関連の設定
type JWTConfig struct {
	Secret      string
	ExpiryHours time.Duration
}

// LogConfig はログ関連の設定
type LogConfig struct {
	Level  string
	Format string // json or console
}

// loadConfig は環境変数から設定を読み込む
func (c *Container) loadConfig() error {
	c.Config = &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getIntEnv("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "financetracker"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth0: Auth0Config{
			Domain:       getEnv("AUTH0_DOMAIN", ""),
			ClientID:     getEnv("AUTH0_CLIENT_ID", ""),
			ClientSecret: getEnv("AUTH0_CLIENT_SECRET", ""),
			CallbackURL:  getEnv("AUTH0_CALLBACK_URL", ""),
			Audience:     getEnv("AUTH0_AUDIENCE", ""),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "default-secret"),
			ExpiryHours: getDurationEnv("JWT_EXPIRY_HOURS", 24*time.Hour),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
	return nil
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv は環境変数を整数として取得し、存在しない場合はデフォルト値を返す
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv は環境変数をtime.Durationとして取得し、存在しない場合はデフォルト値を返す
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
