package ports

import (
	"context"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/savings/domain"
)

// SavingsGoalRepository defines the persistence operations for savings goals
type SavingsGoalRepository interface {
	Create(ctx context.Context, goal *domain.SavingsGoal) error
	GetByID(ctx context.Context, tenantID, goalID string) (*domain.SavingsGoal, error)
	List(ctx context.Context, tenantID string) ([]*domain.SavingsGoal, error)
	ListByStatus(ctx context.Context, tenantID string, status domain.SavingsGoalStatus) ([]*domain.SavingsGoal, error)
	Update(ctx context.Context, goal *domain.SavingsGoal) error
	Delete(ctx context.Context, tenantID, goalID string) error // Soft delete
	CreateTransaction(ctx context.Context, tx *domain.SavingsTransaction) error
	ListTransactions(ctx context.Context, goalID string) ([]*domain.SavingsTransaction, error)
}
