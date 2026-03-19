package transactions

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/repository"
	sharedevents "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/events"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the transactions module
type Module struct {
	expenseHandler  *handlers.ExpenseHandler
	incomeHandler   *handlers.IncomeHandler
	categoryHandler *handlers.CategoryHandler
	categoryRepo    *repository.CategoryRepo
	logger          zerolog.Logger
	authMW          *middleware.AuthMiddleware
	permMW          *middleware.PermissionMiddleware
}

// New creates a new transactions module
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware, permMW *middleware.PermissionMiddleware) *Module {
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
		categoryRepo:    categoryRepo,
		logger:          logger,
		authMW:          authMW,
		permMW:          permMW,
	}
}

// RegisterRoutes registers all HTTP routes for the transactions module.
// All routes require a valid JWT access token and appropriate tenant permissions.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	tx := r.Group("")
	tx.Use(m.authMW.RequireAuth())
	{
		// Expenses
		tx.POST("/expenses", m.permMW.Require("create_transaction"), m.expenseHandler.Create)
		tx.GET("/expenses", m.permMW.Require("view_data"), m.expenseHandler.List)
		tx.GET("/expenses/:id", m.permMW.Require("view_data"), m.expenseHandler.GetByID)
		tx.PUT("/expenses/:id", m.permMW.Require("edit_any_transaction"), m.expenseHandler.Update)
		tx.DELETE("/expenses/:id", m.permMW.Require("delete_any_transaction"), m.expenseHandler.Delete)

		// Incomes
		tx.POST("/incomes", m.permMW.Require("create_transaction"), m.incomeHandler.Create)
		tx.GET("/incomes", m.permMW.Require("view_data"), m.incomeHandler.List)
		tx.GET("/incomes/:id", m.permMW.Require("view_data"), m.incomeHandler.GetByID)
		tx.PUT("/incomes/:id", m.permMW.Require("edit_any_transaction"), m.incomeHandler.Update)
		tx.DELETE("/incomes/:id", m.permMW.Require("delete_any_transaction"), m.incomeHandler.Delete)

		// Categories
		tx.GET("/categories", m.permMW.Require("view_data"), m.categoryHandler.List)
		tx.POST("/categories", m.permMW.Require("create_transaction"), m.categoryHandler.Create)
		tx.PATCH("/categories/:id", m.permMW.Require("edit_any_transaction"), m.categoryHandler.Update)
		tx.DELETE("/categories/:id", m.permMW.Require("delete_any_transaction"), m.categoryHandler.Delete)
	}

	m.logger.Info().Msg("transactions module routes registered")
}

// RegisterSubscribers registers event subscribers for the transactions module.
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	// Seed 15 default categories when a new user registers.
	eventBus.Subscribe("user.registered", func(ctx context.Context, event sharedports.Event) error {
		domainEvt, ok := event.(sharedevents.DomainEvent)
		if !ok {
			return nil
		}
		data, ok := domainEvt.Data.(map[string]string)
		if !ok {
			return nil
		}
		tenantID := data["tenant_id"]
		if tenantID == "" {
			return nil
		}
		if err := m.categoryRepo.SeedDefaultCategories(ctx, event.UserID(), tenantID); err != nil {
			m.logger.Warn().Err(err).Str("user_id", event.UserID()).Msg("failed to seed default categories for new user")
		}
		return nil
	})

	m.logger.Info().Msg("transactions module subscribers registered")
}
