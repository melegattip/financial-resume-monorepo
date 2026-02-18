package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/repository"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the transactions module
type Module struct {
	expenseHandler  *handlers.ExpenseHandler
	incomeHandler   *handlers.IncomeHandler
	categoryHandler *handlers.CategoryHandler
	logger          zerolog.Logger
}

// New creates a new transactions module
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus) *Module {
	// Initialize repositories
	expenseRepo := repository.NewExpenseRepository(db)
	incomeRepo := repository.NewIncomeRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize handlers
	expenseHandler := handlers.NewExpenseHandler(expenseRepo, eventBus, logger)
	incomeHandler := handlers.NewIncomeHandler(incomeRepo, eventBus, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo, logger)

	return &Module{
		expenseHandler:  expenseHandler,
		incomeHandler:   incomeHandler,
		categoryHandler: categoryHandler,
		logger:          logger,
	}
}

// RegisterRoutes registers all HTTP routes for the transactions module
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	// Expenses
	r.POST("/expenses", m.expenseHandler.Create)
	r.GET("/expenses", m.expenseHandler.List)
	r.GET("/expenses/:id", m.expenseHandler.GetByID)
	r.PUT("/expenses/:id", m.expenseHandler.Update)
	r.DELETE("/expenses/:id", m.expenseHandler.Delete)

	// Incomes
	r.POST("/incomes", m.incomeHandler.Create)
	r.GET("/incomes", m.incomeHandler.List)
	r.GET("/incomes/:id", m.incomeHandler.GetByID)
	r.PUT("/incomes/:id", m.incomeHandler.Update)
	r.DELETE("/incomes/:id", m.incomeHandler.Delete)

	// Categories
	r.GET("/categories", m.categoryHandler.List)

	m.logger.Info().Msg("transactions module routes registered")
}

// RegisterSubscribers registers event subscribers (for future use)
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	// In the future, this module could subscribe to events from other modules
	// For example: subscribe to UserDeletedEvent to cascade delete user's expenses
	m.logger.Info().Msg("transactions module subscribers registered")
}
