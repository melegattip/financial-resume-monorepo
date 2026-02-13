package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/http/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
)

// NewRouter creates a configured Gin engine with middleware and routes.
func NewRouter(logger zerolog.Logger, corsOrigins string, healthHandler *handlers.HealthHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestID())
	engine.Use(middleware.CORS(corsOrigins))
	engine.Use(middleware.RequestLogging(logger))

	engine.GET("/health", healthHandler.Handle)

	return engine
}
