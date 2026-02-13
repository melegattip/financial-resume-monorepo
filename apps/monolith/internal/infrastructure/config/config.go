package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// JWTConfig holds JWT token configuration.
type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

// SecurityConfig holds account security policy configuration.
type SecurityConfig struct {
	PasswordMinLength int
	MaxLoginAttempts  int
	LockoutDuration   time.Duration
}

// AppConfig holds all application configuration loaded from environment variables.
type AppConfig struct {
	ServerPort       string
	Environment      string
	LogLevel         string
	CORSAllowOrigins string

	// Database
	DatabaseURL    string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBMaxIdleConns int
	DBMaxOpenConns int

	// Migration source databases (only needed for cmd/migrate)
	UsersDBURL        string
	GamificationDBURL string

	// Auth
	JWT      JWTConfig
	Security SecurityConfig
}

// Load reads configuration from environment variables.
// It returns an error if required variables are missing.
func Load() (*AppConfig, error) {
	cfg := &AppConfig{
		ServerPort:       getEnv("PORT", "8080"),
		Environment:      getEnv("APP_ENV", "development"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           getEnv("DB_NAME", "financial_resume"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		DBMaxIdleConns:   getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns:   getEnvInt("DB_MAX_OPEN_CONNS", 25),
		UsersDBURL:        os.Getenv("USERS_DB_URL"),
		GamificationDBURL: os.Getenv("GAMIFICATION_DB_URL"),
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "financial_resume_secret_key_2024"),
			AccessExpiry:  time.Duration(getEnvInt("JWT_ACCESS_EXPIRY_HOURS", 24)) * time.Hour,
			RefreshExpiry: time.Duration(getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7)) * 24 * time.Hour,
			Issuer:        getEnv("JWT_ISSUER", "financial-resume"),
		},
		Security: SecurityConfig{
			PasswordMinLength: getEnvInt("PASSWORD_MIN_LENGTH", 8),
			MaxLoginAttempts:  getEnvInt("MAX_LOGIN_ATTEMPTS", 5),
			LockoutDuration:   time.Duration(getEnvInt("LOCKOUT_DURATION_MINUTES", 15)) * time.Minute,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string.
// If DATABASE_URL is set, it takes precedence over individual DB_* variables.
func (c *AppConfig) DSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func (c *AppConfig) validate() error {
	if c.DatabaseURL != "" {
		return nil
	}

	if c.DBUser == "" {
		return fmt.Errorf("required environment variable DB_USER is not set (or set DATABASE_URL)")
	}
	if c.DBPassword == "" {
		return fmt.Errorf("required environment variable DB_PASSWORD is not set (or set DATABASE_URL)")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}
