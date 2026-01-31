package config

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/config/environment"
	"gorm.io/gorm"
)

// Handler maneja las peticiones de configuración
type Handler struct {
	db *gorm.DB
}

// NewHandler crea una nueva instancia de Handler
func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// GetConfig retorna la configuración del sistema
func (h *Handler) GetConfig(c *gin.Context) {
	config := map[string]interface{}{
		"environment":       environment.GlobalConfig.Environment,
		"api_base_url":      environment.GlobalConfig.APIBaseURL,
		"gamification_url":  environment.GlobalConfig.GamificationURL,
		"ai_service_url":    environment.GlobalConfig.AIServiceURL,
		"users_service_url": environment.GlobalConfig.UsersServiceURL,
		"version":           "1.0.0",
		"timestamp":         time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, config)
}

// GetDiagnostics retorna información de diagnóstico del sistema
func (h *Handler) GetDiagnostics(c *gin.Context) {
	diagnostics := map[string]interface{}{
		"environment":           environment.GlobalConfig.Environment,
		"timestamp":             time.Now().Format(time.RFC3339),
		"database":              map[string]interface{}{},
		"services":              map[string]interface{}{},
		"environment_variables": map[string]interface{}{},
	}

	// Verificar base de datos
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			diagnostics["database"] = map[string]interface{}{
				"status": "error",
				"error":  err.Error(),
			}
		} else {
			// Test de conexión
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := sqlDB.PingContext(ctx); err != nil {
				diagnostics["database"] = map[string]interface{}{
					"status": "error",
					"error":  err.Error(),
				}
			} else {
				stats := sqlDB.Stats()
				diagnostics["database"] = map[string]interface{}{
					"status":               "healthy",
					"max_open_connections": stats.MaxOpenConnections,
					"open_connections":     stats.OpenConnections,
					"in_use":               stats.InUse,
					"idle":                 stats.Idle,
				}
			}
		}
	}

	// Verificar variables de entorno críticas
	criticalEnvVars := []string{
		"DB_HOST", "DB_USER", "DB_NAME", "DB_PORT", "DB_SSLMODE",
		"JWT_SECRET", "GAMIFICATION_SERVICE_URL", "AI_SERVICE_URL",
		"USERS_SERVICE_URL", "DATABASE_URL",
	}

	for _, envVar := range criticalEnvVars {
		value := os.Getenv(envVar)
		if value != "" {
			// Ocultar valores sensibles
			if envVar == "JWT_SECRET" || envVar == "DB_PASSWORD" {
				diagnostics["environment_variables"].(map[string]interface{})[envVar] = "***"
			} else {
				diagnostics["environment_variables"].(map[string]interface{})[envVar] = value
			}
		} else {
			diagnostics["environment_variables"].(map[string]interface{})[envVar] = "NOT_SET"
		}
	}

	// Verificar servicios externos
	services := map[string]string{
		"gamification":  environment.GlobalConfig.GamificationURL,
		"ai_service":    environment.GlobalConfig.AIServiceURL,
		"users_service": environment.GlobalConfig.UsersServiceURL,
	}

	for serviceName, serviceURL := range services {
		if serviceURL != "" {
			// Aquí podrías hacer un health check real a los servicios
			diagnostics["services"].(map[string]interface{})[serviceName] = map[string]interface{}{
				"url":    serviceURL,
				"status": "configured", // Placeholder - en producción harías un health check real
			}
		}
	}

	c.JSON(http.StatusOK, diagnostics)
}
