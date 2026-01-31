package repository

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// CategoryRepository define la interfaz para el repositorio de categorías
type CategoryRepository interface {
	Create(category *domain.Category) error
	Get(userID string, id string) (*domain.Category, error)
	GetByName(userID string, name string) (*domain.Category, error)
	List(userID string) ([]*domain.Category, error)
	Update(category *domain.Category) error
	Delete(userID string, id string) error

	// Métodos para transacciones relacionadas con categorías
	UpdateCategoryReferences(categoryID string) error
}
