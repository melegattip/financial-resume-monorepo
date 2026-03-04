package recurring

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring/repository"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the recurring transactions module
type Module struct {
	recurringHandler *handlers.RecurringHandler
	logger           zerolog.Logger
	authMW           *middleware.AuthMiddleware
	permMW           *middleware.PermissionMiddleware
}

// New creates and wires up the recurring transactions module.
// db is passed to the handler so that ManualExecute can create
// records in the existing expenses/incomes tables.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware, permMW *middleware.PermissionMiddleware) *Module {
	recurringRepo := repository.NewRecurringRepository(db)
	recurringHandler := handlers.NewRecurringHandler(recurringRepo, db, eventBus, logger)

	return &Module{
		recurringHandler: recurringHandler,
		logger:           logger,
		authMW:           authMW,
		permMW:           permMW,
	}
}

// RegisterRoutes registers all HTTP routes for the recurring module.
// All routes are under /recurring-transactions and require a valid JWT access token.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	recurring := r.Group("/recurring-transactions")
	recurring.Use(m.authMW.RequireAuth())
	{
		recurring.POST("", m.permMW.Require("manage_recurring"), m.recurringHandler.Create)
		recurring.GET("", m.permMW.Require("view_data"), m.recurringHandler.List)
		recurring.GET("/dashboard", m.permMW.Require("view_data"), m.recurringHandler.GetDashboard)
		recurring.GET("/due", m.permMW.Require("view_data"), m.recurringHandler.ListDue)
		recurring.GET("/projection", m.permMW.Require("view_data"), m.recurringHandler.GetProjection)
		recurring.POST("/batch/process", m.permMW.Require("manage_recurring"), m.recurringHandler.ProcessPending)
		recurring.POST("/batch/notify", m.permMW.Require("manage_recurring"), m.recurringHandler.SendNotifications)
		recurring.GET("/:id", m.permMW.Require("view_data"), m.recurringHandler.GetByID)
		recurring.PUT("/:id", m.permMW.Require("manage_recurring"), m.recurringHandler.Update)
		recurring.DELETE("/:id", m.permMW.Require("manage_recurring"), m.recurringHandler.Delete)
		recurring.POST("/:id/pause", m.permMW.Require("manage_recurring"), m.recurringHandler.Pause)
		recurring.POST("/:id/resume", m.permMW.Require("manage_recurring"), m.recurringHandler.Resume)
		recurring.POST("/:id/execute", m.permMW.Require("manage_recurring"), m.recurringHandler.ManualExecute)
	}

	m.logger.Info().Msg("recurring module routes registered")
}

// RegisterSubscribers registers event subscribers for the recurring module.
// Reserved for future use (e.g., reacting to UserDeletedEvent).
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	m.logger.Info().Msg("recurring module subscribers registered")
}

// StartScheduler runs a background goroutine that processes all due recurring
// transactions once on startup and then every hour. Stops when ctx is cancelled.
func (m *Module) StartScheduler(ctx context.Context) {
	go func() {
		m.logger.Info().Msg("recurring scheduler started (interval: 1h)")
		m.recurringHandler.RunAllDue(ctx)

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.recurringHandler.RunAllDue(ctx)
			case <-ctx.Done():
				m.logger.Info().Msg("recurring scheduler stopped")
				return
			}
		}
	}()
}
