package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/domain"
)

// CategoryModel is the GORM model for categories table
type CategoryModel struct {
	ID        string     `gorm:"column:id;type:varchar(255);primaryKey"`
	UserID    string     `gorm:"column:user_id;type:varchar(255);not null;index"`
	Name      string     `gorm:"column:name;not null"`
	Color     string     `gorm:"column:color"`
	Icon      string     `gorm:"column:icon"`
	Priority  int        `gorm:"column:priority;default:0"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

func (CategoryModel) TableName() string {
	return "categories"
}

// ToCategory converts GORM model to domain Category
func (m *CategoryModel) ToCategory() *domain.Category {
	return &domain.Category{
		ID:        m.ID,
		UserID:    m.UserID,
		Name:      m.Name,
		Color:     m.Color,
		Icon:      m.Icon,
		Priority:  m.Priority,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

// FromCategory converts domain Category to GORM model
func FromCategory(c *domain.Category) *CategoryModel {
	return &CategoryModel{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Color:     c.Color,
		Icon:      c.Icon,
		Priority:  c.Priority,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
}

// CategoryRepo implements ports.CategoryRepository
type CategoryRepo struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	if category.ID == "" {
		category.ID = uuid.New().String()
	}

	model := FromCategory(category)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	category.ID = model.ID
	return nil
}

func (r *CategoryRepo) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	var model CategoryModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToCategory(), nil
}

func (r *CategoryRepo) FindByUserID(ctx context.Context, userID string) ([]*domain.Category, error) {
	var models []CategoryModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("CASE WHEN priority = 0 THEN 1 ELSE 0 END ASC, priority ASC, name ASC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, len(models))
	for i, m := range models {
		categories[i] = m.ToCategory()
	}

	return categories, nil
}

func (r *CategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	model := FromCategory(category)
	return r.db.WithContext(ctx).
		Model(&CategoryModel{}).
		Where("id = ? AND deleted_at IS NULL", category.ID).
		Updates(model).Error
}

func (r *CategoryRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&CategoryModel{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now().UTC()).Error
}
