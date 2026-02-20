package recurring

import (
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
}

// New creates and wires up the recurring transactions module.
// db is passed to the handler so that ManualExecute can create
// records in the existing expenses/incomes tables.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware) *Module {
	recurringRepo := repository.NewRecurringRepository(db)
	recurringHandler := handlers.NewRecurringHandler(recurringRepo, db, eventBus, logger)

	return &Module{
		recurringHandler: recurringHandler,
		logger:           logger,
		authMW:           authMW,
	}
}

// RegisterRoutes registers all HTTP routes for the recurring module.
// All routes are under /recurring-transactions and require a valid JWT access token.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	recurring := r.Group("/recurring-transactions")
	recurring.Use(m.authMW.RequireAuth())
	{
		recurring.POST("", m.recurringHandler.Create)
		recurring.GET("", m.recurringHandler.List)
		recurring.GET("/dashboard", m.recurringHandler.GetDashboard)
		recurring.GET("/due", m.recurringHandler.ListDue)
		recurring.GET("/:id", m.recurringHandler.GetByID)
		recurring.PUT("/:id", m.recurringHandler.Update)
		recurring.DELETE("/:id", m.recurringHandler.Delete)
		recurring.POST("/:id/pause", m.recurringHandler.Pause)
		recurring.POST("/:id/resume", m.recurringHandler.Resume)
		recurring.POST("/:id/execute", m.recurringHandler.ManualExecute)
	}

	m.logger.Info().Msg("recurring module routes registered")
}

// RegisterSubscribers registers event subscribers for the recurring module.
// Reserved for future use (e.g., reacting to UserDeletedEvent).
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	m.logger.Info().Msg("recurring module subscribers registered")
}
