package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	"gorm.io/gorm"
)

// BudgetRepository implements ports.BudgetRepository interface (Repository pattern)
type BudgetRepository struct {
	db *gorm.DB
}

// NewBudgetRepository creates a new BudgetRepository instance
func NewBudgetRepository(db *gorm.DB) ports.BudgetRepository {
	return &BudgetRepository{db: db}
}

// Create creates a new budget in the database
func (r *BudgetRepository) Create(ctx context.Context, budget *domain.Budget) error {
	if err := r.db.WithContext(ctx).Create(budget).Error; err != nil {
		return r.handleDBError(err, "error creando presupuesto")
	}
	return nil
}

// GetByID retrieves a budget by ID
func (r *BudgetRepository) GetByID(ctx context.Context, userID, budgetID string) (*domain.Budget, error) {
	var budget domain.Budget

	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", budgetID, userID).
		First(&budget).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewResourceNotFound("presupuesto no encontrado")
		}
		return nil, r.handleDBError(err, "error obteniendo presupuesto")
	}

	return &budget, nil
}

// GetByCategory retrieves a budget by category and period
func (r *BudgetRepository) GetByCategory(ctx context.Context, userID, categoryID string, period domain.BudgetPeriod) (*domain.Budget, error) {
	var budget domain.Budget

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND category_id = ? AND period = ? AND is_active = ?",
			userID, categoryID, period, true).
		First(&budget).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not finding a budget is not an error in this case
		}
		return nil, r.handleDBError(err, "error obteniendo presupuesto por categoría")
	}

	return &budget, nil
}

// List retrieves all budgets for a user
func (r *BudgetRepository) List(ctx context.Context, userID string) ([]*domain.Budget, error) {
	var budgets []domain.Budget

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&budgets).Error

	if err != nil {
		return nil, r.handleDBError(err, "error listando presupuestos")
	}

	// Convert to pointers
	result := make([]*domain.Budget, len(budgets))
	for i := range budgets {
		result[i] = &budgets[i]
	}

	return result, nil
}

// ListActive retrieves all active budgets for a user
func (r *BudgetRepository) ListActive(ctx context.Context, userID string) ([]*domain.Budget, error) {
	var budgets []domain.Budget

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&budgets).Error

	if err != nil {
		return nil, r.handleDBError(err, "error listando presupuestos activos")
	}

	// Convert to pointers
	result := make([]*domain.Budget, len(budgets))
	for i := range budgets {
		result[i] = &budgets[i]
	}

	return result, nil
}

// Update updates an existing budget
func (r *BudgetRepository) Update(ctx context.Context, budget *domain.Budget) error {
	budget.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", budget.ID, budget.UserID).
		Updates(budget)

	if result.Error != nil {
		return r.handleDBError(result.Error, "error actualizando presupuesto")
	}

	if result.RowsAffected == 0 {
		return errors.NewResourceNotFound("presupuesto no encontrado")
	}

	return nil
}

// Delete removes a budget
func (r *BudgetRepository) Delete(ctx context.Context, userID, budgetID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", budgetID, userID).
		Delete(&domain.Budget{})

	if result.Error != nil {
		return r.handleDBError(result.Error, "error eliminando presupuesto")
	}

	if result.RowsAffected == 0 {
		return errors.NewResourceNotFound("presupuesto no encontrado")
	}

	return nil
}

// GetExpensesForPeriod gets total expenses for a category in a specific period
func (r *BudgetRepository) GetExpensesForPeriod(ctx context.Context, userID, categoryID string, startDate, endDate time.Time) (float64, error) {
	var totalExpenses float64

	query := r.db.WithContext(ctx).
		Model(&domain.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate)

	// If category is specified, filter by it
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	err := query.Scan(&totalExpenses).Error

	if err != nil {
		return 0, r.handleDBError(err, "error calculando gastos del período")
	}

	return totalExpenses, nil
}

// GetExpensesByCategory gets expenses grouped by category for a period (for analytics)
func (r *BudgetRepository) GetExpensesByCategory(ctx context.Context, userID string, startDate, endDate time.Time) (map[string]float64, error) {
	type CategoryExpense struct {
		CategoryID string  `gorm:"column:category_id"`
		Amount     float64 `gorm:"column:total_amount"`
	}

	var results []CategoryExpense

	err := r.db.WithContext(ctx).
		Model(&domain.Expense{}).
		Select("category_id, SUM(amount) as total_amount").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Group("category_id").
		Scan(&results).Error

	if err != nil {
		return nil, r.handleDBError(err, "error obteniendo gastos por categoría")
	}

	expenseMap := make(map[string]float64)
	for _, result := range results {
		if result.CategoryID != "" { // Only include expenses with categories
			expenseMap[result.CategoryID] = result.Amount
		}
	}

	return expenseMap, nil
}

// GetBudgetUsageStats gets usage statistics for budget analytics
func (r *BudgetRepository) GetBudgetUsageStats(ctx context.Context, userID string) (*BudgetUsageStats, error) {
	type BudgetStat struct {
		TotalBudgets   int64   `gorm:"column:total_budgets"`
		ActiveBudgets  int64   `gorm:"column:active_budgets"`
		TotalAllocated float64 `gorm:"column:total_allocated"`
		TotalSpent     float64 `gorm:"column:total_spent"`
		OnTrackCount   int64   `gorm:"column:on_track_count"`
		WarningCount   int64   `gorm:"column:warning_count"`
		ExceededCount  int64   `gorm:"column:exceeded_count"`
	}

	var stats BudgetStat

	// Get basic counts and sums
	err := r.db.WithContext(ctx).
		Model(&domain.Budget{}).
		Select(`
			COUNT(*) as total_budgets,
			COUNT(CASE WHEN is_active = true THEN 1 END) as active_budgets,
			COALESCE(SUM(amount), 0) as total_allocated,
			COALESCE(SUM(spent_amount), 0) as total_spent,
			COUNT(CASE WHEN status = 'on_track' THEN 1 END) as on_track_count,
			COUNT(CASE WHEN status = 'warning' THEN 1 END) as warning_count,
			COUNT(CASE WHEN status = 'exceeded' THEN 1 END) as exceeded_count
		`).
		Where("user_id = ?", userID).
		Scan(&stats).Error

	if err != nil {
		return nil, r.handleDBError(err, "error obteniendo estadísticas de presupuesto")
	}

	result := &BudgetUsageStats{
		TotalBudgets:   int(stats.TotalBudgets),
		ActiveBudgets:  int(stats.ActiveBudgets),
		TotalAllocated: stats.TotalAllocated,
		TotalSpent:     stats.TotalSpent,
		OnTrackCount:   int(stats.OnTrackCount),
		WarningCount:   int(stats.WarningCount),
		ExceededCount:  int(stats.ExceededCount),
	}

	// Calculate average usage
	if stats.TotalBudgets > 0 && stats.TotalAllocated > 0 {
		result.AverageUsage = stats.TotalSpent / stats.TotalAllocated
	}

	return result, nil
}

// GetBudgetsNearingLimit gets budgets that are close to their limits
func (r *BudgetRepository) GetBudgetsNearingLimit(ctx context.Context, userID string, threshold float64) ([]*domain.Budget, error) {
	var budgets []domain.Budget

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ? AND (spent_amount / amount) >= ?",
			userID, true, threshold).
		Order("(spent_amount / amount) DESC").
		Find(&budgets).Error

	if err != nil {
		return nil, r.handleDBError(err, "error obteniendo presupuestos cerca del límite")
	}

	// Convert to pointers
	result := make([]*domain.Budget, len(budgets))
	for i := range budgets {
		result[i] = &budgets[i]
	}

	return result, nil
}

// RefreshOutdatedBudgets resets budgets that have moved to a new period
func (r *BudgetRepository) RefreshOutdatedBudgets(ctx context.Context) error {
	now := time.Now()

	// Find budgets where current time is past period_end
	var outdatedBudgets []domain.Budget
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND period_end < ?", true, now).
		Find(&outdatedBudgets).Error

	if err != nil {
		return r.handleDBError(err, "error obteniendo presupuestos vencidos")
	}

	// Reset each outdated budget for new period
	for _, budget := range outdatedBudgets {
		budget.ResetForNewPeriod()

		if err := r.Update(ctx, &budget); err != nil {
			// Log error but continue with other budgets
			continue
		}
	}

	return nil
}

// Helper methods

// handleDBError converts database errors to domain errors
func (r *BudgetRepository) handleDBError(err error, message string) error {
	if err == gorm.ErrRecordNotFound {
		return errors.NewResourceNotFound(message)
	}

	// Check for constraint violations
	if isDuplicateError(err) {
		return errors.NewConflict("presupuesto ya existe")
	}

	return fmt.Errorf("%s: %w", message, err)
}

// isDuplicateError checks if the error is a duplicate key constraint violation
func isDuplicateError(err error) bool {
	// This would depend on your database driver
	// For PostgreSQL, you'd check for specific error codes
	errStr := err.Error()
	return contains(errStr, "duplicate") || contains(errStr, "unique")
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			indexString(s, substr) >= 0))
}

// indexString returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func indexString(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// BudgetUsageStats represents budget usage statistics
type BudgetUsageStats struct {
	TotalBudgets   int     `json:"total_budgets"`
	ActiveBudgets  int     `json:"active_budgets"`
	TotalAllocated float64 `json:"total_allocated"`
	TotalSpent     float64 `json:"total_spent"`
	OnTrackCount   int     `json:"on_track_count"`
	WarningCount   int     `json:"warning_count"`
	ExceededCount  int     `json:"exceeded_count"`
	AverageUsage   float64 `json:"average_usage"`
}
