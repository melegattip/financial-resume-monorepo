package categories

import (
	"github.com/google/uuid"
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// CategoryBuilder es un constructor para crear categorías
type CategoryBuilder struct {
	category *domain.Category
}

// NewCategoryBuilder crea una nueva instancia de CategoryBuilder
func NewCategoryBuilder() *CategoryBuilder {
	return &CategoryBuilder{
		category: &domain.Category{},
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
func (b *CategoryBuilder) Build() *domain.Category {
	if b.category.ID == "" {
		b.category.ID = "cat_" + uuid.New().String()[:8]
	}
	return b.category
}
