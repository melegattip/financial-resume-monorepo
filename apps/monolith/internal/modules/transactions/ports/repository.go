package ports

import (
	"context"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/domain"
)

// ExpenseRepository defines operations for expense persistence
type ExpenseRepository interface {
	Create(ctx context.Context, expense *domain.Expense) error
	FindByID(ctx context.Context, id string) (*domain.Expense, error)
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Expense, error)
	Update(ctx context.Context, expense *domain.Expense) error
	Delete(ctx context.Context, id string) error // Soft delete
	Count(ctx context.Context, userID string) (int64, error)
}

// IncomeRepository defines operations for income persistence
type IncomeRepository interface {
	Create(ctx context.Context, income *domain.Income) error
	FindByID(ctx context.Context, id string) (*domain.Income, error)
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Income, error)
	Update(ctx context.Context, income *domain.Income) error
	Delete(ctx context.Context, id string) error // Soft delete
	Count(ctx context.Context, userID string) (int64, error)
}

// CategoryRepository defines operations for category persistence
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id string) error // Soft delete
}
