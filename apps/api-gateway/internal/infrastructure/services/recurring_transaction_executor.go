package services

import (
	"context"
	"fmt"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

// RecurringTransactionExecutorService implements the executor interface
type RecurringTransactionExecutorService struct {
	expenseRepo        baseRepo.ExpenseRepository
	incomeRepo         baseRepo.IncomeRepository
	gamificationHelper *GamificationHelper
}

// NewRecurringTransactionExecutorService creates a new executor service
func NewRecurringTransactionExecutorService(
	expenseRepo baseRepo.ExpenseRepository,
	incomeRepo baseRepo.IncomeRepository,
	gamificationHelper *GamificationHelper,
) ports.RecurringTransactionExecutorService {
	return &RecurringTransactionExecutorService{
		expenseRepo:        expenseRepo,
		incomeRepo:         incomeRepo,
		gamificationHelper: gamificationHelper,
	}
}

// ExecuteTransaction executes a recurring transaction (same as CreateTransactionFromRecurring)
func (s *RecurringTransactionExecutorService) ExecuteTransaction(ctx context.Context, transaction *domain.RecurringTransaction) error {
	return s.CreateTransactionFromRecurring(ctx, transaction)
}

// CreateTransactionFromRecurring creates a real transaction from a recurring transaction
func (s *RecurringTransactionExecutorService) CreateTransactionFromRecurring(ctx context.Context, recurring *domain.RecurringTransaction) error {
	switch recurring.Type {
	case "expense":
		return s.createExpenseFromRecurring(ctx, recurring)
	case "income":
		return s.createIncomeFromRecurring(ctx, recurring)
	default:
		return fmt.Errorf("tipo de transacción no soportado: %s", recurring.Type)
	}
}

// createExpenseFromRecurring creates an expense from a recurring transaction
func (s *RecurringTransactionExecutorService) createExpenseFromRecurring(ctx context.Context, recurring *domain.RecurringTransaction) error {
	// Create expense using domain builder
	expenseBuilder := domain.NewExpenseBuilder().
		SetID(domain.NewExpenseID()).
		SetUserID(recurring.UserID).
		SetAmount(recurring.Amount).
		SetDescription(recurring.Description).
		SetPaid(false).        // New expenses are unpaid by default
		SetDueDate(time.Now()) // Due today

		// Set category if provided
	if recurring.CategoryID != nil {
		expenseBuilder = expenseBuilder.SetCategoryID(*recurring.CategoryID)
	}

	expense := expenseBuilder.Build()

	// Save the expense
	if err := s.expenseRepo.Create(expense); err != nil {
		return fmt.Errorf("error creando gasto desde recurrente: %w", err)
	}

	// Gamificación: registrar creación de gasto, asignación de categoría (si aplica) y ejecución de recurrente
	if s.gamificationHelper != nil {
		s.gamificationHelper.RecordExpenseAction(
			recurring.UserID,
			ActionCreateExpense,
			expense.ID,
			"Gasto creado automáticamente desde transacción recurrente",
		)
		if recurring.CategoryID != nil && *recurring.CategoryID != "" {
			s.gamificationHelper.RecordActionAsync(
				recurring.UserID,
				ActionAssignCategory,
				EntityExpense,
				expense.ID,
				"Categoría asignada automáticamente desde transacción recurrente",
			)
		}
		s.gamificationHelper.RecordActionAsync(
			recurring.UserID,
			ActionExecuteRecurring,
			EntityRecurring,
			recurring.ID,
			"Ejecución de transacción recurrente (gasto)",
		)
	}

	return nil
}

// createIncomeFromRecurring creates an income from a recurring transaction
func (s *RecurringTransactionExecutorService) createIncomeFromRecurring(ctx context.Context, recurring *domain.RecurringTransaction) error {
	// Create income using domain builder
	incomeBuilder := domain.NewIncomeBuilder().
		SetID(domain.NewIncomeID()).
		SetUserID(recurring.UserID).
		SetAmount(recurring.Amount).
		SetDescription(recurring.Description)

		// Set category if provided
	if recurring.CategoryID != nil {
		incomeBuilder = incomeBuilder.SetCategoryID(*recurring.CategoryID)
	}

	income := incomeBuilder.Build()

	// Save the income
	if err := s.incomeRepo.Create(income); err != nil {
		return fmt.Errorf("error creando ingreso desde recurrente: %w", err)
	}

	// Gamificación: registrar creación de ingreso, asignación de categoría (si aplica) y ejecución de recurrente
	if s.gamificationHelper != nil {
		s.gamificationHelper.RecordIncomeAction(
			recurring.UserID,
			ActionCreateIncome,
			income.ID,
			"Ingreso creado automáticamente desde transacción recurrente",
		)
		if recurring.CategoryID != nil && *recurring.CategoryID != "" {
			s.gamificationHelper.RecordActionAsync(
				recurring.UserID,
				ActionAssignCategory,
				EntityIncome,
				income.ID,
				"Categoría asignada automáticamente desde transacción recurrente",
			)
		}
		s.gamificationHelper.RecordActionAsync(
			recurring.UserID,
			ActionExecuteRecurring,
			EntityRecurring,
			recurring.ID,
			"Ejecución de transacción recurrente (ingreso)",
		)
	}

	return nil
}
