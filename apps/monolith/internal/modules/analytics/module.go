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
	permMW           *middleware.PermissionMiddleware
}

// New creates and wires all dependencies for the analytics module.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware, permMW *middleware.PermissionMiddleware) *Module {
	analyticsRepo := repository.NewAnalyticsRepository(db)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsRepo, logger)

	return &Module{
		analyticsHandler: analyticsHandler,
		logger:           logger,
		authMW:           authMW,
		permMW:           permMW,
	}
}

// RegisterRoutes registers all HTTP routes for the analytics module.
// All routes are protected by the JWT auth middleware.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	// Dashboard — top-level route.
	dashboard := r.Group("/dashboard")
	dashboard.Use(m.authMW.RequireAuth())
	{
		dashboard.GET("", m.permMW.Require("view_data"), m.analyticsHandler.GetDashboard)
	}

	// Analytics sub-routes.
	analytics := r.Group("/analytics")
	analytics.Use(m.authMW.RequireAuth())
	{
		analytics.GET("/expenses", m.permMW.Require("view_data"), m.analyticsHandler.GetExpenseSummary)
		analytics.GET("/incomes", m.permMW.Require("view_data"), m.analyticsHandler.GetIncomeSummary)
		analytics.GET("/categories", m.permMW.Require("view_data"), m.analyticsHandler.GetCategoryAnalysis)
		analytics.GET("/monthly", m.permMW.Require("view_data"), m.analyticsHandler.GetMonthlyTrends)
	}

	// Insights sub-routes.
	insights := r.Group("/insights")
	insights.Use(m.authMW.RequireAuth())
	{
		insights.GET("/financial-health", m.permMW.Require("view_data"), m.analyticsHandler.GetFinancialHealth)
	}

	// Reports route.
	reports := r.Group("/reports")
	reports.Use(m.authMW.RequireAuth())
	{
		reports.GET("", m.permMW.Require("view_data"), m.analyticsHandler.GetReport)
	}

	m.logger.Info().Msg("analytics module routes registered")
}
