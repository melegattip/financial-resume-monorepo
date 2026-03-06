package budgets

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/repository"
	transactionsdomain "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/domain"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the budgets module.
type Module struct {
	budgetHandler *handlers.BudgetHandler
	repo          *repository.BudgetRepo
	logger        zerolog.Logger
	authMW        *middleware.AuthMiddleware
	permMW        *middleware.PermissionMiddleware
}

// New creates a new budgets Module, wiring the repository and handler.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware, permMW *middleware.PermissionMiddleware) *Module {
	budgetRepo := repository.NewBudgetRepository(db)
	budgetHandler := handlers.NewBudgetHandler(budgetRepo, logger, eventBus)

	return &Module{
		budgetHandler: budgetHandler,
		repo:          budgetRepo,
		logger:        logger,
		authMW:        authMW,
		permMW:        permMW,
	}
}

// RegisterRoutes registers all HTTP routes for the budgets module.
// Literal routes (/status, /dashboard) are registered before /:id so that
// Gin matches them before the wildcard parameter.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	budgets := r.Group("/budgets")
	budgets.Use(m.authMW.RequireAuth())
	{
		budgets.POST("", m.permMW.Require("manage_budgets"), m.budgetHandler.Create)
		budgets.GET("", m.permMW.Require("view_data"), m.budgetHandler.List)
		budgets.GET("/status", m.permMW.Require("view_data"), m.budgetHandler.GetStatus)
		budgets.GET("/dashboard", m.permMW.Require("view_data"), m.budgetHandler.GetDashboard)
		budgets.GET("/:id", m.permMW.Require("view_data"), m.budgetHandler.GetByID)
		budgets.PUT("/:id", m.permMW.Require("manage_budgets"), m.budgetHandler.Update)
		budgets.DELETE("/:id", m.permMW.Require("manage_budgets"), m.budgetHandler.Delete)
	}

	m.logger.Info().Msg("budgets module routes registered")
}

// RegisterSubscribers registers event subscribers for the budgets module.
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	// Recalculate budget spent amounts when an expense is created.
	eventBus.Subscribe("expense.created", func(ctx context.Context, event sharedports.Event) error {
		ev, ok := event.(transactionsdomain.ExpenseCreatedEvent)
		if !ok || ev.TenantID == "" || ev.CategoryID == "" {
			return nil
		}

		budgets, err := m.repo.ListActive(ctx, ev.TenantID)
		if err != nil {
			m.logger.Warn().Err(err).Str("tenant_id", ev.TenantID).Msg("budget subscriber: failed to list active budgets")
			return nil
		}

		now := time.Now().UTC()
		for _, b := range budgets {
			if b.CategoryID != ev.CategoryID {
				continue
			}
			// Only process budgets whose period includes today.
			if now.Before(b.PeriodStart) || now.After(b.PeriodEnd) {
				continue
			}

			oldStatus := b.Status

			newSpent, err := m.repo.GetExpensesForPeriod(ctx, ev.TenantID, ev.CategoryID, b.PeriodStart, b.PeriodEnd)
			if err != nil {
				m.logger.Warn().Err(err).Str("budget_id", b.ID).Msg("budget subscriber: failed to get expenses for period")
				continue
			}

			b.SpentAmount = newSpent
			b.UpdatedAt = now

			// Recalculate status.
			if b.Amount > 0 {
				pct := newSpent / b.Amount
				switch {
				case pct >= 1.0:
					b.Status = domain.BudgetStatusExceeded
				case pct >= b.AlertAt:
					b.Status = domain.BudgetStatusWarning
				default:
					b.Status = domain.BudgetStatusOnTrack
				}
			}

			if err := m.repo.Update(ctx, b); err != nil {
				m.logger.Warn().Err(err).Str("budget_id", b.ID).Msg("budget subscriber: failed to update budget spent_amount")
				continue
			}

			// Publish threshold event if status worsened.
			if b.Status != oldStatus && b.Status != domain.BudgetStatusOnTrack {
				spentPct := 0.0
				if b.Amount > 0 {
					spentPct = (b.SpentAmount / b.Amount) * 100
				}
				thresholdEvent := domain.BudgetThresholdCrossedEvent{
					BudgetID:    b.ID,
					User:        b.UserID,
					TenantID:    b.TenantID,
					CategoryID:  b.CategoryID,
					SpentAmount: b.SpentAmount,
					BudgetLimit: b.Amount,
					SpentPct:    spentPct,
					Period:      string(b.Period),
					NewStatus:   string(b.Status),
					Timestamp:   now,
				}
				if err := eventBus.Publish(ctx, thresholdEvent); err != nil {
					m.logger.Warn().Err(err).Str("budget_id", b.ID).Msg("budget subscriber: failed to publish threshold event")
				}
			}
		}
		return nil
	})

	m.logger.Info().Msg("budgets module subscribers registered")
}
