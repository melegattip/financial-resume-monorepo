package update

import (
	"context"
	"log"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
)

// Service maneja la lógica de negocio para actualizar gastos
type Service struct {
	repository         baseRepo.ExpenseRepository
	percentageObserver transactions.PercentageTransactionObserver
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.ExpenseRepository, percentageObserver transactions.PercentageTransactionObserver) *Service {
	return &Service{
		repository:         repository,
		percentageObserver: percentageObserver,
	}
}

// UpdateExpense actualiza un gasto existente
func (s *Service) UpdateExpense(ctx context.Context, userID, id string, request *expenses.UpdateExpenseRequest) (*expenses.UpdateExpenseResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}

	if id == "" {
		return nil, errors.NewBadRequest("El ID del gasto es requerido")
	}

	// Solo validar amount si se está actualizando (diferente de 0)
	if request.Amount != 0 && request.Amount <= 0 {
		return nil, errors.NewBadRequest("El monto del gasto debe ser mayor a 0")
	}

	// Obtenemos el gasto existente
	expense, err := s.repository.Get(userID, id)
	if err != nil {
		return nil, err
	}

	// Actualizamos los campos solo si se proporcionan
	if request.Amount != 0 {
		expense.Amount = request.Amount
	}

	// Manejar pagos parciales
	if request.PaymentAmount > 0 {
		expense.AddPayment(request.PaymentAmount)
	}

	if request.Description != "" {
		expense.Description = request.Description
	}
	if request.CategoryID != "" {
		if request.CategoryID == "" {
			expense.CategoryID = nil
		} else {
			expense.CategoryID = &request.CategoryID
		}
	}
	if request.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", request.DueDate)
		if err != nil {
			return nil, errors.NewBadRequest("Formato de fecha inválido")
		}
		expense.DueDate = dueDate
	}
	// Actualizar paid siempre (es un bool, no se puede saber si se envió o no fácilmente)
	// Pero solo si no se está procesando un pago (AddPayment ya maneja el estado paid)
	if request.PaymentAmount == 0 {
		expense.Paid = request.Paid
	}
	expense.UpdatedAt = time.Now()

	// Recalcular el porcentaje solo si el monto cambió
	if request.Amount != 0 {
		totalIncome, err := s.percentageObserver.GetTotalIncome(ctx, userID)
		if err != nil {
			return nil, err
		}
		expense.CalculatePercentage(totalIncome)
	}

	// Guardamos los cambios
	if err := s.repository.Update(expense); err != nil {
		return nil, err
	}

	// Notificar al observer para que actualice los porcentajes de los demás gastos solo si cambió el amount
	if request.Amount != 0 {
		if err := s.percentageObserver.OnTransactionUpdated(ctx, expense); err != nil {
			// Log del error pero no lo retornamos para no fallar la actualización del gasto
			log.Printf("Error al recalcular porcentajes: %v", err)
		}
	}

	return &expenses.UpdateExpenseResponse{
		ID:            expense.ID,
		UserID:        expense.UserID,
		Amount:        expense.Amount,
		AmountPaid:    expense.AmountPaid,
		PendingAmount: expense.GetPendingAmount(),
		Description:   expense.Description,
		CategoryID: func() string {
			if expense.CategoryID != nil {
				return *expense.CategoryID
			}
			return ""
		}(),
		Paid:       expense.Paid,
		DueDate:    expense.DueDate.Format("2006-01-02"),
		CreatedAt:  expense.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  expense.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Percentage: expense.Percentage,
	}, nil
}
