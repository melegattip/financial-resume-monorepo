package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config representa la configuración completa del servicio
type Config struct {
	Server ServerConfig `json:"server"`
	OpenAI OpenAIConfig `json:"openai"`
	Redis  RedisConfig  `json:"redis"`
	Cache  CacheConfig  `json:"cache"`
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// OpenAIConfig configuración de OpenAI
type OpenAIConfig struct {
	APIKey  string `json:"api_key"`
	UseMock bool   `json:"use_mock"`
}

// RedisConfig configuración de Redis
type RedisConfig struct {
	URL      string `json:"url"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// CacheConfig configuración del cache
type CacheConfig struct {
	DefaultTTLMinutes int `json:"default_ttl_minutes"`
	InsightsTTLHours  int `json:"insights_ttl_hours"`
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8082"),
			Host: getEnv("HOST", "localhost"),
		},
		OpenAI: OpenAIConfig{
			APIKey:  getEnv("OPENAI_API_KEY", ""),
			UseMock: getEnvBool("USE_AI_MOCK", true),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Cache: CacheConfig{
			DefaultTTLMinutes: getEnvInt("CACHE_DEFAULT_TTL_MINUTES", 30),
			InsightsTTLHours:  getEnvInt("CACHE_INSIGHTS_TTL_HOURS", 20),
		},
	}

	// Validar configuración requerida
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// validate valida la configuración
func (c *Config) validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if !c.OpenAI.UseMock && c.OpenAI.APIKey == "" {
		return fmt.Errorf("OpenAI API key is required when not using mock")
	}

	if c.Redis.URL == "" {
		return fmt.Errorf("Redis URL is required")
	}

	return nil
}

// getEnv obtiene una variable de entorno con un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool obtiene una variable de entorno booleana
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getEnvInt obtiene una variable de entorno entera
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
