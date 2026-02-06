package get

import (
	"context"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	expensesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
)

// Service maneja la lógica de negocio para obtener gastos
type Service struct {
	repository baseRepo.ExpenseRepository
}

// ExpenseGetter define la interfaz para obtener gastos
type ExpenseGetter interface {
	GetExpense(ctx context.Context, userID, id string) (*expensesDomain.GetExpenseResponse, error)
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.ExpenseRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// GetExpense obtiene un gasto por su ID
func (s *Service) GetExpense(ctx context.Context, userID, id string) (*expensesDomain.GetExpenseResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}

	if id == "" {
		return nil, errors.NewBadRequest("El ID del gasto es requerido")
	}

	expense, err := s.repository.Get(userID, id)
	if err != nil {
		return nil, err
	}

	return &expensesDomain.GetExpenseResponse{
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
		Paid:            expense.Paid,
		DueDate:         expense.DueDate.Format("2006-01-02"),
		TransactionDate: expense.TransactionDate.Format("2006-01-02"),
		CreatedAt:       expense.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       expense.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Percentage:      expense.Percentage,
	}, nil
}
