package delete

import (
	"context"
	"log"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions"
)

// Service maneja la lógica de negocio para eliminar gastos
type Service struct {
	repository         baseRepo.ExpenseRepository
	percentageObserver transactions.PercentageTransactionObserver
}

// ExpenseDeleter define la interfaz para eliminar gastos
type ExpenseDeleter interface {
	DeleteExpense(ctx context.Context, userID, id string) error
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.ExpenseRepository, percentageObserver transactions.PercentageTransactionObserver) *Service {
	return &Service{
		repository:         repository,
		percentageObserver: percentageObserver,
	}
}

// DeleteExpense elimina un gasto existente
func (s *Service) DeleteExpense(ctx context.Context, userID, id string) error {
	if userID == "" {
		return errors.NewBadRequest("El ID del usuario es requerido")
	}

	if id == "" {
		return errors.NewBadRequest("El ID del gasto es requerido")
	}

	// Eliminamos el gasto
	err := s.repository.Delete(userID, id)
	if err != nil {
		return err
	}

	// Notificar al observer para que actualice los porcentajes de los demás gastos
	if err := s.percentageObserver.OnTransactionDeleted(ctx, id); err != nil {
		// Log del error pero no lo retornamos para no fallar la eliminación del gasto
		log.Printf("Error al recalcular porcentajes: %v", err)
	}

	return nil
}
