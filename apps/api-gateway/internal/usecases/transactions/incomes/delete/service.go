package delete

import (
	"context"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

// Service maneja la lógica de negocio para eliminar ingresos
type Service struct {
	repository baseRepo.IncomeRepository
}

// IncomeDeleter define la interfaz para eliminar ingresos
type IncomeDeleter interface {
	DeleteIncome(ctx context.Context, userID string, id string) error
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.IncomeRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// DeleteIncome elimina un ingreso existente
func (s *Service) DeleteIncome(ctx context.Context, userID string, id string) error {
	if userID == "" {
		return errors.NewBadRequest("El ID del usuario es requerido")
	}

	if id == "" {
		return errors.NewBadRequest("El ID del ingreso es requerido")
	}

	// Verificamos que el ingreso existe
	_, err := s.repository.Get(userID, id)
	if err != nil {
		return err
	}

	return s.repository.Delete(userID, id)
}
