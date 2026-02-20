package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
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
	authMW          *middleware.AuthMiddleware
}

// New creates a new transactions module
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware) *Module {
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
		authMW:          authMW,
	}
}

// RegisterRoutes registers all HTTP routes for the transactions module.
// All routes require a valid JWT access token.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	tx := r.Group("")
	tx.Use(m.authMW.RequireAuth())
	{
		// Expenses
		tx.POST("/expenses", m.expenseHandler.Create)
		tx.GET("/expenses", m.expenseHandler.List)
		tx.GET("/expenses/:id", m.expenseHandler.GetByID)
		tx.PUT("/expenses/:id", m.expenseHandler.Update)
		tx.DELETE("/expenses/:id", m.expenseHandler.Delete)

		// Incomes
		tx.POST("/incomes", m.incomeHandler.Create)
		tx.GET("/incomes", m.incomeHandler.List)
		tx.GET("/incomes/:id", m.incomeHandler.GetByID)
		tx.PUT("/incomes/:id", m.incomeHandler.Update)
		tx.DELETE("/incomes/:id", m.incomeHandler.Delete)

		// Categories
		tx.GET("/categories", m.categoryHandler.List)
	}

	m.logger.Info().Msg("transactions module routes registered")
}

// RegisterSubscribers registers event subscribers (for future use)
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	// In the future, this module could subscribe to events from other modules
	// For example: subscribe to UserDeletedEvent to cascade delete user's expenses
	m.logger.Info().Msg("transactions module subscribers registered")
}
