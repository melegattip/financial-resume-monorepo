package config

import (
	"os"
	"strconv"
)

// GetEnv returns the value of an environment variable or a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt returns the value of an environment variable as int or a default value
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvBool returns the value of an environment variable as bool or a default value
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDatabaseConfigFromEnv creates a DatabaseConfig from environment variables
func NewDatabaseConfigFromEnv(prefix string) DatabaseConfig {
	return DatabaseConfig{
		Host:     GetEnv(prefix+"DB_HOST", "localhost"),
		Port:     GetEnvInt(prefix+"DB_PORT", 5432),
		User:     GetEnv(prefix+"DB_USER", "postgres"),
		Password: GetEnv(prefix+"DB_PASSWORD", "postgres"),
		DBName:   GetEnv(prefix+"DB_NAME", "financial_resume"),
		SSLMode:  GetEnv(prefix+"DB_SSLMODE", "disable"),
	}
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedisConfigFromEnv creates a RedisConfig from environment variables
func NewRedisConfigFromEnv() RedisConfig {
	return RedisConfig{
		Host:     GetEnv("REDIS_HOST", "localhost"),
		Port:     GetEnvInt("REDIS_PORT", 6379),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       GetEnvInt("REDIS_DB", 0),
	}
}

// JWTConfigEnv holds JWT configuration from environment
type JWTConfigEnv struct {
	Secret            string
	AccessExpiryHours int
	RefreshExpiryDays int
}

// NewJWTConfigFromEnv creates JWT config from environment variables
func NewJWTConfigFromEnv() JWTConfigEnv {
	return JWTConfigEnv{
		Secret:            GetEnv("JWT_SECRET", "default_secret_change_me"),
		AccessExpiryHours: GetEnvInt("JWT_ACCESS_EXPIRY_HOURS", 24),
		RefreshExpiryDays: GetEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7),
	}
}
