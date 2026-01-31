package list

import (
	"context"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	expensesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
)

// Service maneja la lógica de negocio para listar gastos
type Service struct {
	repository baseRepo.ExpenseRepository
}

// ExpenseLister define la interfaz para listar gastos
type ExpenseLister interface {
	ListExpenses(ctx context.Context, userID string) (*expensesDomain.ListExpensesResponse, error)
	ListUnpaidExpenses(ctx context.Context, userID string) (*expensesDomain.ListExpensesResponse, error)
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.ExpenseRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// ListExpenses lista todos los gastos de un usuario
func (s *Service) ListExpenses(ctx context.Context, userID string) (*expensesDomain.ListExpensesResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}

	expenses, err := s.repository.List(userID)
	if err != nil {
		return nil, err
	}

	expensesResponse := make([]expensesDomain.GetExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expensesResponse[i] = expensesDomain.GetExpenseResponse{
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
		}
	}

	return &expensesDomain.ListExpensesResponse{
		Expenses: expensesResponse,
	}, nil
}

// ListUnpaidExpenses lista todos los gastos no pagados de un usuario
func (s *Service) ListUnpaidExpenses(ctx context.Context, userID string) (*expensesDomain.ListExpensesResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}

	expenses, err := s.repository.ListUnpaid(userID)
	if err != nil {
		return nil, err
	}

	expensesResponse := make([]expensesDomain.GetExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expensesResponse[i] = expensesDomain.GetExpenseResponse{
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
		}
	}

	return &expensesDomain.ListExpensesResponse{
		Expenses: expensesResponse,
	}, nil
}
