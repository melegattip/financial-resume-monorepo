package ports

import (
	"context"
	"time"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/domain"
)

// BudgetRepository defines the persistence operations for budgets.
type BudgetRepository interface {
	// Create persists a new budget.
	Create(ctx context.Context, budget *domain.Budget) error

	// GetByID retrieves a budget by its ID, scoped to the given user.
	GetByID(ctx context.Context, userID, budgetID string) (*domain.Budget, error)

	// GetByCategory retrieves the budget for a given user, category, and period.
	GetByCategory(ctx context.Context, userID, categoryID string, period domain.BudgetPeriod) (*domain.Budget, error)

	// List returns all (non-deleted) budgets for a user.
	List(ctx context.Context, userID string) ([]*domain.Budget, error)

	// ListActive returns all active (is_active=true, non-deleted) budgets for a user.
	ListActive(ctx context.Context, userID string) ([]*domain.Budget, error)

	// Update saves changes to an existing budget.
	Update(ctx context.Context, budget *domain.Budget) error

	// Delete soft-deletes the budget identified by budgetID, scoped to the given user.
	Delete(ctx context.Context, userID, budgetID string) error

	// GetExpensesForPeriod returns the total expense amount for a user's category
	// between startDate and endDate (inclusive).
	GetExpensesForPeriod(ctx context.Context, userID, categoryID string, startDate, endDate time.Time) (float64, error)
}
