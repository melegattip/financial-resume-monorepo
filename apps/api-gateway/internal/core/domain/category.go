// Package domain defines the core business entities and their behavior
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrEmptyCategoryName is returned when attempting to create a category with an empty name
	ErrEmptyCategoryName = errors.New("category name cannot be empty")
)

// Category representa una categoría en el sistema
type Category struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	Name      string    `json:"name" gorm:"column:name;not null"`
	UserID    string    `json:"user_id" gorm:"column:user_id;not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (Category) TableName() string {
	return "categories"
}

// CategoryBuilder es un constructor para crear categorías
type CategoryBuilder struct {
	category *Category
}

// NewCategoryBuilder crea una nueva instancia de CategoryBuilder
func NewCategoryBuilder() *CategoryBuilder {
	return &CategoryBuilder{
		category: &Category{},
	}
}

// SetID establece el ID de la categoría
func (b *CategoryBuilder) SetID(id string) *CategoryBuilder {
	b.category.ID = id
	return b
}

// SetName establece el nombre de la categoría
func (b *CategoryBuilder) SetName(name string) *CategoryBuilder {
	b.category.Name = name
	return b
}

// SetUserID establece el ID del usuario
func (b *CategoryBuilder) SetUserID(userID string) *CategoryBuilder {
	b.category.UserID = userID
	return b
}

// Build construye y retorna la categoría
func (b *CategoryBuilder) Build() *Category {
	if b.category.ID == "" {
		b.category.ID = "cat_" + uuid.New().String()[:8]
	}
	return b.category
}

// Validate checks if the category is valid
// Returns ErrEmptyCategoryName if the name is empty
func (c *Category) Validate() error {
	if c.Name == "" {
		return ErrEmptyCategoryName
	}
	return nil
}
