package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring/domain"
)

// RecurringTransactionModel is the GORM model for the recurring_transactions table
type RecurringTransactionModel struct {
	ID             string     `gorm:"column:id;type:varchar(255);primaryKey"`
	UserID         string     `gorm:"column:user_id;type:varchar(255);not null;index"`
	Amount         float64    `gorm:"column:amount;not null"`
	Description    string     `gorm:"column:description;not null"`
	CategoryID     *string    `gorm:"column:category_id;type:varchar(255);index"`
	Type           string     `gorm:"column:type;type:varchar(20);not null"`
	Frequency      string     `gorm:"column:frequency;type:varchar(20);not null"`
	NextDate       time.Time  `gorm:"column:next_date;not null;index"`
	LastExecuted   *time.Time `gorm:"column:last_executed"`
	IsActive       bool       `gorm:"column:is_active;default:true"`
	AutoCreate     bool       `gorm:"column:auto_create;default:true"`
	NotifyBefore   int        `gorm:"column:notify_before;default:1"`
	EndDate        *time.Time `gorm:"column:end_date"`
	ExecutionCount int        `gorm:"column:execution_count;default:0"`
	MaxExecutions  *int       `gorm:"column:max_executions"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index"`
}

// TableName returns the database table name
func (RecurringTransactionModel) TableName() string {
	return "recurring_transactions"
}

// toDomain converts the GORM model to the domain entity
func (m *RecurringTransactionModel) toDomain() *domain.RecurringTransaction {
	return &domain.RecurringTransaction{
		ID:             m.ID,
		UserID:         m.UserID,
		Amount:         m.Amount,
		Description:    m.Description,
		CategoryID:     m.CategoryID,
		Type:           m.Type,
		Frequency:      m.Frequency,
		NextDate:       m.NextDate,
		LastExecuted:   m.LastExecuted,
		IsActive:       m.IsActive,
		AutoCreate:     m.AutoCreate,
		NotifyBefore:   m.NotifyBefore,
		EndDate:        m.EndDate,
		ExecutionCount: m.ExecutionCount,
		MaxExecutions:  m.MaxExecutions,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
}

// fromDomain converts the domain entity to the GORM model
func fromDomain(rt *domain.RecurringTransaction) *RecurringTransactionModel {
	return &RecurringTransactionModel{
		ID:             rt.ID,
		UserID:         rt.UserID,
		Amount:         rt.Amount,
		Description:    rt.Description,
		CategoryID:     rt.CategoryID,
		Type:           rt.Type,
		Frequency:      rt.Frequency,
		NextDate:       rt.NextDate,
		LastExecuted:   rt.LastExecuted,
		IsActive:       rt.IsActive,
		AutoCreate:     rt.AutoCreate,
		NotifyBefore:   rt.NotifyBefore,
		EndDate:        rt.EndDate,
		ExecutionCount: rt.ExecutionCount,
		MaxExecutions:  rt.MaxExecutions,
		CreatedAt:      rt.CreatedAt,
		UpdatedAt:      rt.UpdatedAt,
		DeletedAt:      rt.DeletedAt,
	}
}

// RecurringRepo implements ports.RecurringTransactionRepository
type RecurringRepo struct {
	db *gorm.DB
}

// NewRecurringRepository creates a new recurring transaction repository
func NewRecurringRepository(db *gorm.DB) *RecurringRepo {
	return &RecurringRepo{db: db}
}

// Create persists a new recurring transaction
func (r *RecurringRepo) Create(ctx context.Context, rt *domain.RecurringTransaction) error {
	if rt.ID == "" {
		rt.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	if rt.CreatedAt.IsZero() {
		rt.CreatedAt = now
	}
	if rt.UpdatedAt.IsZero() {
		rt.UpdatedAt = now
	}

	model := fromDomain(rt)
	return r.db.WithContext(ctx).Create(model).Error
}

// GetByID returns a recurring transaction by ID scoped to the user.
// Returns nil, nil when not found.
func (r *RecurringRepo) GetByID(ctx context.Context, userID, id string) (*domain.RecurringTransaction, error) {
	var model RecurringTransactionModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.toDomain(), nil
}

// List returns all non-deleted recurring transactions for a user
func (r *RecurringRepo) List(ctx context.Context, userID string) ([]*domain.RecurringTransaction, error) {
	var models []RecurringTransactionModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	result := make([]*domain.RecurringTransaction, len(models))
	for i := range models {
		result[i] = models[i].toDomain()
	}
	return result, nil
}

// ListActive returns only active non-deleted recurring transactions for a user
func (r *RecurringRepo) ListActive(ctx context.Context, userID string) ([]*domain.RecurringTransaction, error) {
	var models []RecurringTransactionModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = true AND deleted_at IS NULL", userID).
		Order("next_date ASC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	result := make([]*domain.RecurringTransaction, len(models))
	for i := range models {
		result[i] = models[i].toDomain()
	}
	return result, nil
}

// ListDue returns all active recurring transactions whose NextDate is on or before now.
// Intended for use by a cron job.
func (r *RecurringRepo) ListDue(ctx context.Context, now time.Time) ([]*domain.RecurringTransaction, error) {
	var models []RecurringTransactionModel
	err := r.db.WithContext(ctx).
		Where("is_active = true AND next_date <= ? AND deleted_at IS NULL", now).
		Order("next_date ASC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	result := make([]*domain.RecurringTransaction, len(models))
	for i := range models {
		result[i] = models[i].toDomain()
	}
	return result, nil
}

// Update persists changes to an existing recurring transaction
func (r *RecurringRepo) Update(ctx context.Context, rt *domain.RecurringTransaction) error {
	rt.UpdatedAt = time.Now().UTC()
	model := fromDomain(rt)
	return r.db.WithContext(ctx).
		Model(&RecurringTransactionModel{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", rt.ID, rt.UserID).
		Updates(model).Error
}

// Delete soft-deletes a recurring transaction scoped to the given user
func (r *RecurringRepo) Delete(ctx context.Context, userID, id string) error {
	return r.db.WithContext(ctx).
		Model(&RecurringTransactionModel{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("deleted_at", time.Now().UTC()).Error
}
