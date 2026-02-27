package ports

import (
	"context"
	"time"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring/domain"
)

// RecurringTransactionRepository defines persistence operations for recurring transactions
type RecurringTransactionRepository interface {
	// Create persists a new recurring transaction
	Create(ctx context.Context, rt *domain.RecurringTransaction) error

	// GetByID returns a recurring transaction by its ID, scoped to the given tenant.
	// Returns nil, nil when not found.
	GetByID(ctx context.Context, tenantID, id string) (*domain.RecurringTransaction, error)

	// List returns all non-deleted recurring transactions for a tenant
	List(ctx context.Context, tenantID string) ([]*domain.RecurringTransaction, error)

	// ListActive returns only active (IsActive=true) non-deleted recurring transactions for a tenant
	ListActive(ctx context.Context, tenantID string) ([]*domain.RecurringTransaction, error)

	// ListDue returns all active recurring transactions whose NextDate is on or before now.
	// Intended for use by a cron job to find transactions that need to be executed.
	ListDue(ctx context.Context, now time.Time) ([]*domain.RecurringTransaction, error)

	// Update persists changes to an existing recurring transaction
	Update(ctx context.Context, rt *domain.RecurringTransaction) error

	// Delete soft-deletes a recurring transaction scoped to the given tenant
	Delete(ctx context.Context, tenantID, id string) error
}
