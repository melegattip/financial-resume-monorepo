package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/domain"
)

// ExpenseModel is the GORM model for expenses table
type ExpenseModel struct {
	ID              string     `gorm:"column:id;type:varchar(255);primaryKey"`
	UserID          string     `gorm:"column:user_id;type:varchar(255);not null;index"`
	TenantID        string     `gorm:"column:tenant_id;type:varchar(50);index"`
	CategoryID      string     `gorm:"column:category_id;type:varchar(255);not null"`
	Amount          float64    `gorm:"column:amount;not null"`
	Description     string     `gorm:"column:description;not null"`
	TransactionDate time.Time  `gorm:"column:transaction_date;not null;index"`
	PaymentMethod   string     `gorm:"column:payment_method"`
	Notes           string     `gorm:"column:notes"`
	Paid            bool       `gorm:"column:paid;default:false"`
	AmountPaid      float64    `gorm:"column:amount_paid;default:0"`
	PendingAmount   float64    `gorm:"column:pending_amount;default:0"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;index"`
}

func (ExpenseModel) TableName() string {
	return "expenses"
}

// ToExpense converts GORM model to domain Expense
func (m *ExpenseModel) ToExpense() *domain.Expense {
	return &domain.Expense{
		ID:              m.ID,
		UserID:          m.UserID,
		TenantID:        m.TenantID,
		CategoryID:      m.CategoryID,
		Amount:          m.Amount,
		Description:     m.Description,
		TransactionDate: m.TransactionDate,
		PaymentMethod:   m.PaymentMethod,
		Notes:           m.Notes,
		Paid:            m.Paid,
		AmountPaid:      m.AmountPaid,
		PendingAmount:   m.PendingAmount,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DeletedAt:       m.DeletedAt,
	}
}

// FromExpense converts domain Expense to GORM model
func FromExpense(e *domain.Expense) *ExpenseModel {
	return &ExpenseModel{
		ID:              e.ID,
		UserID:          e.UserID,
		TenantID:        e.TenantID,
		CategoryID:      e.CategoryID,
		Amount:          e.Amount,
		Description:     e.Description,
		TransactionDate: e.TransactionDate,
		PaymentMethod:   e.PaymentMethod,
		Notes:           e.Notes,
		Paid:            e.Paid,
		AmountPaid:      e.AmountPaid,
		PendingAmount:   e.PendingAmount,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		DeletedAt:       e.DeletedAt,
	}
}

// ExpenseRepo implements ports.ExpenseRepository
type ExpenseRepo struct {
	db *gorm.DB
}

// NewExpenseRepository creates a new expense repository
func NewExpenseRepository(db *gorm.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Create(ctx context.Context, expense *domain.Expense) error {
	if expense.ID == "" {
		expense.ID = uuid.New().String()
	}

	model := FromExpense(expense)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	expense.ID = model.ID
	return nil
}

func (r *ExpenseRepo) FindByID(ctx context.Context, id string) (*domain.Expense, error) {
	var model ExpenseModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToExpense(), nil
}

func (r *ExpenseRepo) FindByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*domain.Expense, error) {
	var models []ExpenseModel
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("transaction_date DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	expenses := make([]*domain.Expense, len(models))
	for i, m := range models {
		expenses[i] = m.ToExpense()
	}

	return expenses, nil
}

func (r *ExpenseRepo) Update(ctx context.Context, expense *domain.Expense) error {
	// Use map to ensure zero-values (paid=false, amount_paid=0) are persisted.
	updates := map[string]interface{}{
		"category_id":      expense.CategoryID,
		"amount":           expense.Amount,
		"description":      expense.Description,
		"transaction_date": expense.TransactionDate,
		"payment_method":   expense.PaymentMethod,
		"notes":            expense.Notes,
		"paid":             expense.Paid,
		"amount_paid":      expense.AmountPaid,
		"pending_amount":   expense.PendingAmount,
		"updated_at":       time.Now().UTC(),
	}
	return r.db.WithContext(ctx).
		Model(&ExpenseModel{}).
		Where("id = ? AND deleted_at IS NULL", expense.ID).
		Updates(updates).Error
}

func (r *ExpenseRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&ExpenseModel{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now().UTC()).Error
}

func (r *ExpenseRepo) Count(ctx context.Context, tenantID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&ExpenseModel{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Count(&count).Error
	return count, err
}
