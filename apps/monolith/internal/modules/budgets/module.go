package budgets

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/repository"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the budgets module.
type Module struct {
	budgetHandler *handlers.BudgetHandler
	logger        zerolog.Logger
}

// New creates a new budgets Module, wiring the repository and handler.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus) *Module {
	budgetRepo := repository.NewBudgetRepository(db)
	budgetHandler := handlers.NewBudgetHandler(budgetRepo, logger)

	return &Module{
		budgetHandler: budgetHandler,
		logger:        logger,
	}
}

// RegisterRoutes registers all HTTP routes for the budgets module.
// The /budgets/status route must be registered before /budgets/:id so that
// Gin matches the literal segment before the wildcard parameter.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/budgets", m.budgetHandler.Create)
	r.GET("/budgets", m.budgetHandler.List)
	r.GET("/budgets/status", m.budgetHandler.GetStatus)
	r.GET("/budgets/:id", m.budgetHandler.GetByID)
	r.PUT("/budgets/:id", m.budgetHandler.Update)
	r.DELETE("/budgets/:id", m.budgetHandler.Delete)

	m.logger.Info().Msg("budgets module routes registered")
}

// RegisterSubscribers registers event subscribers for the budgets module.
// Currently the module does not subscribe to any events; this hook is provided
// for future use (e.g. recalculating spent amounts on ExpenseCreatedEvent).
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	m.logger.Info().Msg("budgets module subscribers registered")
}
