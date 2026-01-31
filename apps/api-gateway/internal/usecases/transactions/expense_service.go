package transactions

import (
	"context"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

// ExpenseService define las operaciones relacionadas con gastos
type ExpenseService interface {
	// GetTotalIncome obtiene el total de ingresos para un usuario
	GetTotalIncome(ctx context.Context, userID string) (float64, error)

	// GetAllExpenses obtiene todos los gastos de un usuario
	GetAllExpenses(ctx context.Context, userID string) ([]domain.Transaction, error)

	// UpdateExpense actualiza un gasto existente
	UpdateExpense(ctx context.Context, expense domain.Transaction) error

	// UpdateExpensePercentages actualiza los porcentajes de gastos para un usuario
	UpdateExpensePercentages(ctx context.Context, userID string) error

	// GetTransaction obtiene una transacción por su ID
	GetTransaction(ctx context.Context, transactionID string) (domain.Transaction, error)
}

type expenseService struct {
	expenseRepo baseRepo.ExpenseRepository
	incomeRepo  baseRepo.IncomeRepository
}

func NewExpenseService(expenseRepo baseRepo.ExpenseRepository, incomeRepo baseRepo.IncomeRepository) ExpenseService {
	return &expenseService{
		expenseRepo: expenseRepo,
		incomeRepo:  incomeRepo,
	}
}

func (s *expenseService) GetTotalIncome(ctx context.Context, userID string) (float64, error) {
	incomes, err := s.incomeRepo.List(userID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, income := range incomes {
		total += income.Amount
	}

	return total, nil
}

func (s *expenseService) GetAllExpenses(ctx context.Context, userID string) ([]domain.Transaction, error) {
	expenses, err := s.expenseRepo.List(userID)
	if err != nil {
		return nil, err
	}

	transactions := make([]domain.Transaction, len(expenses))
	for i, expense := range expenses {
		transactions[i] = expense
	}

	return transactions, nil
}

func (s *expenseService) UpdateExpense(ctx context.Context, expense domain.Transaction) error {
	exp, ok := expense.(*domain.Expense)
	if !ok {
		return nil // No actualizamos si no es un gasto
	}
	return s.expenseRepo.Update(exp)
}

func (s *expenseService) UpdateExpensePercentages(ctx context.Context, userID string) error {
	// Obtener el total de ingresos
	totalIncome, err := s.GetTotalIncome(ctx, userID)
	if err != nil {
		return err
	}

	// Obtener todos los gastos
	expenses, err := s.expenseRepo.List(userID)
	if err != nil {
		return err
	}

	// Actualizar el porcentaje de cada gasto
	for _, expense := range expenses {
		expense.CalculatePercentage(totalIncome)
		if err := s.expenseRepo.Update(expense); err != nil {
			return err
		}
	}

	return nil
}

func (s *expenseService) GetTransaction(ctx context.Context, transactionID string) (domain.Transaction, error) {
	// Primero intentamos obtener el gasto
	expense, err := s.expenseRepo.Get("", transactionID)
	if err == nil {
		return expense, nil
	}

	// Si no es un gasto, intentamos obtener el ingreso
	income, err := s.incomeRepo.Get("", transactionID)
	if err == nil {
		return income, nil
	}

	// Si no encontramos ni gasto ni ingreso, retornamos error
	return nil, err
}
