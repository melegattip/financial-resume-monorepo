package repository

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// IncomeRepository define las operaciones para el repositorio de ingresos
type IncomeRepository interface {
	Create(income *domain.Income) error
	Get(userID string, id string) (*domain.Income, error)
	List(userID string) ([]*domain.Income, error)
	Update(income *domain.Income) error
	Delete(userID string, id string) error
}
