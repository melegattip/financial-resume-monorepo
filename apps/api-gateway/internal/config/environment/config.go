// Package environment provides configuration for different deployment environments
package environment

import (
	"log"
	"os"

	environment "github.com/melegattip/financial-resume-engine/internal/config/environment/constants"
)

// ServiceConfig holds all service URLs and configurations
type ServiceConfig struct {
	Environment     string
	APIBaseURL      string
	GamificationURL string
	AIServiceURL    string
	UsersServiceURL string
	DatabaseURL     string
	RedisURL        string
	LogLevel        string
	EnableCache     string
	MaxConnections  string
	GinMode         string
}

// Global configuration instance
var GlobalConfig ServiceConfig

// detectEnvironment automatically detects the deployment environment based on cloud provider variables
func detectEnvironment() string {
	// Check for Render.com environment variables
	if os.Getenv("RENDER") == "true" || os.Getenv("RENDER_SERVICE_ID") != "" {
		log.Println("🔍 [Environment] Detected: Render.com")
		return environment.Render
	}

	// Check for GCP Cloud Run environment variables
	if os.Getenv("K_SERVICE") != "" || os.Getenv("GOOGLE_CLOUD_PROJECT") != "" || os.Getenv("K_REVISION") != "" {
		log.Println("🔍 [Environment] Detected: Google Cloud Platform")
		return environment.GCP
	}

	// Check for explicit environment variable
	if goEnv := os.Getenv("GO_ENVIRONMENT"); goEnv != "" {
		log.Printf("🔍 [Environment] Explicit GO_ENVIRONMENT: %s", goEnv)
		if goEnv == environment.Production {
			return environment.GCP // Production defaults to GCP
		}
		return goEnv
	}

	// Default to development for local execution
	log.Println("🔍 [Environment] Detected: Development (localhost)")
	return environment.Development
}

// SetUp initializes the environment configuration based on automatic environment detection
func SetUp() {
	log.Println("🚀 [Environment] Global Environment Setup - Auto Detection")

	_ = os.Setenv("APPLICATION", environment.Application)

	// Detect environment automatically
	detectedEnv := detectEnvironment()
	os.Setenv("GO_ENVIRONMENT", detectedEnv)

	// Configure based on detected environment
	switch detectedEnv {
	case environment.Render:
		GlobalConfig = setupRenderEnvironment()
	case environment.Development:
		GlobalConfig = setupDevelopmentEnvironment()
	case environment.Beta:
		GlobalConfig = setupBetaEnvironment()
	default:
		log.Printf("⚠️ [Environment] Unknown environment '%s', falling back to development", detectedEnv)
		GlobalConfig = setupDevelopmentEnvironment()
	}

	// Apply configuration to environment variables
	applyConfigToEnv(GlobalConfig)

	log.Printf("✅ [Environment] Setup complete for: %s", GlobalConfig.Environment)
	log.Printf("🔗 [Environment] API Base URL: %s", GlobalConfig.APIBaseURL)
	log.Printf("🎮 [Environment] Gamification URL: %s", GlobalConfig.GamificationURL)
	log.Printf("🤖 [Environment] AI Service URL: %s", GlobalConfig.AIServiceURL)
	log.Printf("👤 [Environment] Users Service URL: %s", GlobalConfig.UsersServiceURL)
}

// setupDevelopmentEnvironment configures the environment variables for local development
func setupDevelopmentEnvironment() ServiceConfig {
	log.Println("🔧 [Environment] Configuring for Development")

	// Use environment variables if available, otherwise fall back to localhost
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	if usersServiceURL == "" {
		usersServiceURL = "http://localhost:8083/api/v1"
	}

	gamificationServiceURL := os.Getenv("GAMIFICATION_SERVICE_URL")
	if gamificationServiceURL == "" {
		gamificationServiceURL = "http://localhost:8081/api/v1"
	}

	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "http://localhost:8082/api/v1"
	}

	return ServiceConfig{
		Environment:     environment.Development,
		APIBaseURL:      "http://localhost:8080",
		GamificationURL: gamificationServiceURL,
		AIServiceURL:    aiServiceURL,
		UsersServiceURL: usersServiceURL,
		DatabaseURL:     "postgresql://postgres:postgres@localhost:5432/financial_resume",
		RedisURL:        "redis://localhost:6379",
		LogLevel:        "debug",
		EnableCache:     "false",
		MaxConnections:  "20",
		GinMode:         "debug",
	}
}

// setupRenderEnvironment configures the environment variables for Render.com deployment
func setupRenderEnvironment() ServiceConfig {
	log.Println("🎯 [Environment] Configuring for Render.com")

	// Get service URLs from environment variables with fallbacks
	gamificationServiceURL := os.Getenv("GAMIFICATION_SERVICE_URL")
	if gamificationServiceURL == "" {
		gamificationServiceURL = "https://financial-gamification-service.onrender.com/api/v1"
	}

	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "https://financial-ai-api.niloft.com/api/v1"
	}

	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	if usersServiceURL == "" {
		usersServiceURL = "https://users-service-mp5p.onrender.com/api/v1"
	}

	return ServiceConfig{
		Environment:     environment.Render,
		APIBaseURL:      "https://financial-resume-engine.onrender.com",
		GamificationURL: gamificationServiceURL,
		AIServiceURL:    aiServiceURL,
		UsersServiceURL: usersServiceURL,
		DatabaseURL:     os.Getenv("DATABASE_URL"), // Render provides this automatically
		RedisURL:        os.Getenv("REDIS_URL"),    // Render provides this automatically
		LogLevel:        "info",
		EnableCache:     "true",
		MaxConnections:  "50",
		GinMode:         "release",
	}
}

// setupBetaEnvironment configures the environment variables for beta deployment
func setupBetaEnvironment() ServiceConfig {
	log.Println("🧪 [Environment] Configuring for Beta")

	// Get service URLs from environment variables with beta fallbacks
	gamificationServiceURL := os.Getenv("GAMIFICATION_SERVICE_URL")
	if gamificationServiceURL == "" {
		gamificationServiceURL = "https://beta-financial-gamification-service.onrender.com"
	}

	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "https://beta-financial-ai-service.onrender.com"
	}

	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	if usersServiceURL == "" {
		usersServiceURL = "https://beta-users-service.onrender.com"
	}

	return ServiceConfig{
		Environment:     environment.Beta,
		APIBaseURL:      "https://beta-financial-resume-engine.onrender.com",
		GamificationURL: gamificationServiceURL,
		AIServiceURL:    aiServiceURL,
		UsersServiceURL: usersServiceURL,
		DatabaseURL:     "postgresql://user:password@beta-db.niloft.com:5432/financial_db",
		RedisURL:        "redis://beta-redis.niloft.com:6379",
		LogLevel:        "debug",
		EnableCache:     "true",
		MaxConnections:  "50",
		GinMode:         "debug",
	}
}

// applyConfigToEnv applies the ServiceConfig to environment variables
func applyConfigToEnv(config ServiceConfig) {
	// URLs
	os.Setenv("API_URL", config.APIBaseURL)
	os.Setenv("GAMIFICATION_SERVICE_URL", config.GamificationURL)
	os.Setenv("AI_SERVICE_URL", config.AIServiceURL)
	os.Setenv("USERS_SERVICE_URL", config.UsersServiceURL)

	// Database and Redis
	os.Setenv("DATABASE_URL", config.DatabaseURL)
	os.Setenv("REDIS_URL", config.RedisURL)

	// Application settings
	os.Setenv("LOG_LEVEL", config.LogLevel)
	os.Setenv("ENABLE_CACHE", config.EnableCache)
	os.Setenv("MAX_CONNECTIONS", config.MaxConnections)
	os.Setenv("GIN_MODE", config.GinMode)
}

// GetConfig returns the current global configuration
func GetConfig() ServiceConfig {
	return GlobalConfig
}

// GetAPIBaseURL returns the API base URL for the current environment
func GetAPIBaseURL() string {
	return GlobalConfig.APIBaseURL
}

// GetGamificationURL returns the gamification service URL for the current environment
func GetGamificationURL() string {
	return GlobalConfig.GamificationURL
}

// GetAIServiceURL returns the AI service URL for the current environment
func GetAIServiceURL() string {
	return GlobalConfig.AIServiceURL
}

// GetUsersServiceURL returns the users service URL for the current environment
func GetUsersServiceURL() string {
	return GlobalConfig.UsersServiceURL
}

// IsDevelopment returns true if the current environment is development
func IsDevelopment() bool {
	return GlobalConfig.Environment == environment.Development
}

// IsProduction returns true if the current environment is production (Render or GCP)
func IsProduction() bool {
	return GlobalConfig.Environment == environment.Render || GlobalConfig.Environment == environment.GCP
}

// getBaseUrls returns the API and internal URLs based on the environment
// Returns two strings: the API URL and the internal URL
// DEPRECATED: Use specific environment setup functions instead
func getBaseUrls(goEnvironment string) (string, string) {
	if environment.Production == goEnvironment {
		return "http://internal.niloft.com", "http://internal.niloft.com"
	}

	return "https://internal-api.niloft.com", "https://internal-api.niloft.com"
}
