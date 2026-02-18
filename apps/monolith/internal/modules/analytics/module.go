package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/analytics/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/analytics/repository"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the analytics module.
type Module struct {
	analyticsHandler *handlers.AnalyticsHandler
	logger           zerolog.Logger
	authMW           *middleware.AuthMiddleware
}

// New creates and wires all dependencies for the analytics module.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware) *Module {
	analyticsRepo := repository.NewAnalyticsRepository(db)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsRepo, logger)

	return &Module{
		analyticsHandler: analyticsHandler,
		logger:           logger,
		authMW:           authMW,
	}
}

// RegisterRoutes registers all HTTP routes for the analytics module.
// All routes are protected by the JWT auth middleware.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	// Dashboard — top-level route.
	dashboard := r.Group("/dashboard")
	dashboard.Use(m.authMW.RequireAuth())
	{
		dashboard.GET("", m.analyticsHandler.GetDashboard)
	}

	// Analytics sub-routes.
	analytics := r.Group("/analytics")
	analytics.Use(m.authMW.RequireAuth())
	{
		analytics.GET("/expenses", m.analyticsHandler.GetExpenseSummary)
		analytics.GET("/incomes", m.analyticsHandler.GetIncomeSummary)
		analytics.GET("/categories", m.analyticsHandler.GetCategoryAnalysis)
		analytics.GET("/monthly", m.analyticsHandler.GetMonthlyTrends)
	}

	m.logger.Info().Msg("analytics module routes registered")
}
