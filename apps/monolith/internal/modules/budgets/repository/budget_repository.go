package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/domain"
)

// BudgetModel is the GORM model for the budgets table.
type BudgetModel struct {
	ID          string     `gorm:"column:id;type:varchar(255);primaryKey"`
	UserID      string     `gorm:"column:user_id;type:varchar(255);not null;index"`
	TenantID    string     `gorm:"column:tenant_id;type:varchar(50);index"`
	CategoryID  string     `gorm:"column:category_id;type:varchar(255);not null;index"`
	Amount      float64    `gorm:"column:amount;not null"`
	SpentAmount float64    `gorm:"column:spent_amount;not null;default:0"`
	Period      string     `gorm:"column:period;type:varchar(20);not null"`
	PeriodStart time.Time  `gorm:"column:period_start;not null"`
	PeriodEnd   time.Time  `gorm:"column:period_end;not null"`
	AlertAt     float64    `gorm:"column:alert_at;not null;default:0.8"`
	Status      string     `gorm:"column:status;type:varchar(20);not null;default:'on_track'"`
	IsActive    bool       `gorm:"column:is_active;not null;default:true"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index"`
}

// TableName returns the database table name for budgets.
func (BudgetModel) TableName() string {
	return "budgets"
}

// ToBudget converts a GORM model to a domain Budget.
func (m *BudgetModel) ToBudget() *domain.Budget {
	return &domain.Budget{
		ID:          m.ID,
		UserID:      m.UserID,
		TenantID:    m.TenantID,
		CategoryID:  m.CategoryID,
		Amount:      m.Amount,
		SpentAmount: m.SpentAmount,
		Period:      domain.BudgetPeriod(m.Period),
		PeriodStart: m.PeriodStart,
		PeriodEnd:   m.PeriodEnd,
		AlertAt:     m.AlertAt,
		Status:      domain.BudgetStatus(m.Status),
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromBudget converts a domain Budget to a GORM model.
func FromBudget(b *domain.Budget) *BudgetModel {
	return &BudgetModel{
		ID:          b.ID,
		UserID:      b.UserID,
		TenantID:    b.TenantID,
		CategoryID:  b.CategoryID,
		Amount:      b.Amount,
		SpentAmount: b.SpentAmount,
		Period:      string(b.Period),
		PeriodStart: b.PeriodStart,
		PeriodEnd:   b.PeriodEnd,
		AlertAt:     b.AlertAt,
		Status:      string(b.Status),
		IsActive:    b.IsActive,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

// BudgetRepo implements ports.BudgetRepository using GORM.
type BudgetRepo struct {
	db *gorm.DB
}

// NewBudgetRepository creates a new BudgetRepo.
func NewBudgetRepository(db *gorm.DB) *BudgetRepo {
	return &BudgetRepo{db: db}
}

// Create persists a new budget.
func (r *BudgetRepo) Create(ctx context.Context, budget *domain.Budget) error {
	model := FromBudget(budget)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a budget by ID, scoped to the given tenant.
func (r *BudgetRepo) GetByID(ctx context.Context, tenantID, budgetID string) (*domain.Budget, error) {
	var model BudgetModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", budgetID, tenantID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToBudget(), nil
}

// GetByCategory retrieves the budget for a given tenant, category, and period.
func (r *BudgetRepo) GetByCategory(ctx context.Context, tenantID, categoryID string, period domain.BudgetPeriod) (*domain.Budget, error) {
	var model BudgetModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND category_id = ? AND period = ? AND deleted_at IS NULL", tenantID, categoryID, string(period)).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToBudget(), nil
}

// List returns all non-deleted budgets for a tenant, ordered by creation date descending.
func (r *BudgetRepo) List(ctx context.Context, tenantID string) ([]*domain.Budget, error) {
	var models []BudgetModel
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	budgets := make([]*domain.Budget, len(models))
	for i, m := range models {
		budgets[i] = m.ToBudget()
	}
	return budgets, nil
}

// ListActive returns all active, non-deleted budgets for a tenant.
func (r *BudgetRepo) ListActive(ctx context.Context, tenantID string) ([]*domain.Budget, error) {
	var models []BudgetModel
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = true AND deleted_at IS NULL", tenantID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	budgets := make([]*domain.Budget, len(models))
	for i, m := range models {
		budgets[i] = m.ToBudget()
	}
	return budgets, nil
}

// Update saves changes to an existing budget.
func (r *BudgetRepo) Update(ctx context.Context, budget *domain.Budget) error {
	model := FromBudget(budget)
	return r.db.WithContext(ctx).
		Model(&BudgetModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", budget.ID, budget.TenantID).
		Updates(map[string]interface{}{
			"amount":       model.Amount,
			"spent_amount": model.SpentAmount,
			"alert_at":     model.AlertAt,
			"status":       model.Status,
			"is_active":    model.IsActive,
			"period_start": model.PeriodStart,
			"period_end":   model.PeriodEnd,
			"updated_at":   model.UpdatedAt,
		}).Error
}

// Delete soft-deletes the budget identified by budgetID, scoped to the given tenant.
func (r *BudgetRepo) Delete(ctx context.Context, tenantID, budgetID string) error {
	return r.db.WithContext(ctx).
		Model(&BudgetModel{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", budgetID, tenantID).
		Update("deleted_at", time.Now().UTC()).Error
}

// GetExpensesForPeriod returns the total expense amount for a tenant's category
// within the given date range. It queries the expenses table directly.
func (r *BudgetRepo) GetExpensesForPeriod(ctx context.Context, tenantID, categoryID string, startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Table("expenses").
		Select("COALESCE(SUM(amount), 0)").
		Where("tenant_id = ? AND category_id = ? AND transaction_date BETWEEN ? AND ? AND deleted_at IS NULL",
			tenantID, categoryID, startDate, endDate).
		Scan(&total).Error

	if err != nil {
		return 0, err
	}
	return total, nil
}
