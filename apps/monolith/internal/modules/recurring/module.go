package recurring

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring/repository"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the recurring transactions module
type Module struct {
	recurringHandler *handlers.RecurringHandler
	logger           zerolog.Logger
}

// New creates and wires up the recurring transactions module.
// db is passed to the handler so that ManualExecute can create
// records in the existing expenses/incomes tables.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus) *Module {
	recurringRepo := repository.NewRecurringRepository(db)
	recurringHandler := handlers.NewRecurringHandler(recurringRepo, db, eventBus, logger)

	return &Module{
		recurringHandler: recurringHandler,
		logger:           logger,
	}
}

// RegisterRoutes registers all HTTP routes for the recurring module.
// All routes require a JWT auth middleware to be applied to the router group beforehand.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/recurring", m.recurringHandler.Create)
	r.GET("/recurring", m.recurringHandler.List)
	r.GET("/recurring/due", m.recurringHandler.ListDue)
	r.GET("/recurring/:id", m.recurringHandler.GetByID)
	r.PUT("/recurring/:id", m.recurringHandler.Update)
	r.DELETE("/recurring/:id", m.recurringHandler.Delete)
	r.POST("/recurring/:id/pause", m.recurringHandler.Pause)
	r.POST("/recurring/:id/resume", m.recurringHandler.Resume)
	r.POST("/recurring/:id/execute", m.recurringHandler.ManualExecute)

	m.logger.Info().Msg("recurring module routes registered")
}

// RegisterSubscribers registers event subscribers for the recurring module.
// Reserved for future use (e.g., reacting to UserDeletedEvent).
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	m.logger.Info().Msg("recurring module subscribers registered")
}
