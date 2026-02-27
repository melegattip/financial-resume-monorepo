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

	// GetByID retrieves a budget by its ID, scoped to the given tenant.
	GetByID(ctx context.Context, tenantID, budgetID string) (*domain.Budget, error)

	// GetByCategory retrieves the budget for a given tenant, category, and period.
	GetByCategory(ctx context.Context, tenantID, categoryID string, period domain.BudgetPeriod) (*domain.Budget, error)

	// List returns all (non-deleted) budgets for a tenant.
	List(ctx context.Context, tenantID string) ([]*domain.Budget, error)

	// ListActive returns all active (is_active=true, non-deleted) budgets for a tenant.
	ListActive(ctx context.Context, tenantID string) ([]*domain.Budget, error)

	// Update saves changes to an existing budget.
	Update(ctx context.Context, budget *domain.Budget) error

	// Delete soft-deletes the budget identified by budgetID, scoped to the given tenant.
	Delete(ctx context.Context, tenantID, budgetID string) error

	// GetExpensesForPeriod returns the total expense amount for a tenant's category
	// between startDate and endDate (inclusive).
	GetExpensesForPeriod(ctx context.Context, tenantID, categoryID string, startDate, endDate time.Time) (float64, error)
}
