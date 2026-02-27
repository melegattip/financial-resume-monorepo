package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/domain"
)

// IncomeModel is the GORM model for incomes table
type IncomeModel struct {
	ID           string     `gorm:"column:id;type:varchar(255);primaryKey"`
	UserID       string     `gorm:"column:user_id;type:varchar(255);not null;index"`
	TenantID     string     `gorm:"column:tenant_id;type:varchar(50);index"`
	Amount       float64    `gorm:"column:amount;not null"`
	Source       string     `gorm:"column:source"`
	Description  string     `gorm:"column:description"`
	ReceivedDate time.Time  `gorm:"column:received_date;not null;index"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index"`
}

func (IncomeModel) TableName() string {
	return "incomes"
}

// ToIncome converts GORM model to domain Income
func (m *IncomeModel) ToIncome() *domain.Income {
	return &domain.Income{
		ID:           m.ID,
		UserID:       m.UserID,
		TenantID:     m.TenantID,
		Amount:       m.Amount,
		Source:       m.Source,
		Description:  m.Description,
		ReceivedDate: m.ReceivedDate,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    m.DeletedAt,
	}
}

// FromIncome converts domain Income to GORM model
func FromIncome(i *domain.Income) *IncomeModel {
	return &IncomeModel{
		ID:           i.ID,
		UserID:       i.UserID,
		TenantID:     i.TenantID,
		Amount:       i.Amount,
		Source:       i.Source,
		Description:  i.Description,
		ReceivedDate: i.ReceivedDate,
		CreatedAt:    i.CreatedAt,
		UpdatedAt:    i.UpdatedAt,
		DeletedAt:    i.DeletedAt,
	}
}

// IncomeRepo implements ports.IncomeRepository
type IncomeRepo struct {
	db *gorm.DB
}

// NewIncomeRepository creates a new income repository
func NewIncomeRepository(db *gorm.DB) *IncomeRepo {
	return &IncomeRepo{db: db}
}

func (r *IncomeRepo) Create(ctx context.Context, income *domain.Income) error {
	if income.ID == "" {
		income.ID = uuid.New().String()
	}

	model := FromIncome(income)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	income.ID = model.ID
	return nil
}

func (r *IncomeRepo) FindByID(ctx context.Context, id string) (*domain.Income, error) {
	var model IncomeModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToIncome(), nil
}

func (r *IncomeRepo) FindByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*domain.Income, error) {
	var models []IncomeModel
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("received_date DESC, created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	incomes := make([]*domain.Income, len(models))
	for i, m := range models {
		incomes[i] = m.ToIncome()
	}

	return incomes, nil
}

func (r *IncomeRepo) Update(ctx context.Context, income *domain.Income) error {
	model := FromIncome(income)
	return r.db.WithContext(ctx).
		Model(&IncomeModel{}).
		Where("id = ? AND deleted_at IS NULL", income.ID).
		Updates(model).Error
}

func (r *IncomeRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&IncomeModel{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now().UTC()).Error
}

func (r *IncomeRepo) Count(ctx context.Context, tenantID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&IncomeModel{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Count(&count).Error
	return count, err
}
