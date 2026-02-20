package budgets

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/repository"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the budgets module.
type Module struct {
	budgetHandler *handlers.BudgetHandler
	logger        zerolog.Logger
	authMW        *middleware.AuthMiddleware
}

// New creates a new budgets Module, wiring the repository and handler.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware) *Module {
	budgetRepo := repository.NewBudgetRepository(db)
	budgetHandler := handlers.NewBudgetHandler(budgetRepo, logger)

	return &Module{
		budgetHandler: budgetHandler,
		logger:        logger,
		authMW:        authMW,
	}
}

// RegisterRoutes registers all HTTP routes for the budgets module.
// Literal routes (/status, /dashboard) are registered before /:id so that
// Gin matches them before the wildcard parameter.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	budgets := r.Group("/budgets")
	budgets.Use(m.authMW.RequireAuth())
	{
		budgets.POST("", m.budgetHandler.Create)
		budgets.GET("", m.budgetHandler.List)
		budgets.GET("/status", m.budgetHandler.GetStatus)
		budgets.GET("/dashboard", m.budgetHandler.GetDashboard)
		budgets.GET("/:id", m.budgetHandler.GetByID)
		budgets.PUT("/:id", m.budgetHandler.Update)
		budgets.DELETE("/:id", m.budgetHandler.Delete)
	}

	m.logger.Info().Msg("budgets module routes registered")
}

// RegisterSubscribers registers event subscribers for the budgets module.
// Currently the module does not subscribe to any events; this hook is provided
// for future use (e.g. recalculating spent amounts on ExpenseCreatedEvent).
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	m.logger.Info().Msg("budgets module subscribers registered")
}
