package transactions

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionComponent_CreateAndUpdate(t *testing.T) {
	// Configurar mocks
	mockIncomeService := new(MockIncomeService)
	mockExpenseService := new(MockExpenseService)
	factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

	// Caso 1: Crear un gasto y actualizarlo
	t.Run("Crear y actualizar gasto", func(t *testing.T) {
		dueDate := time.Now().AddDate(0, 1, 0)
		ctx := context.Background()

		// Mock para crear gasto
		mockExpenseService.On("CreateExpense", ctx, &expenses.CreateExpenseRequest{
			UserID:      "user1",
			Amount:      100.0,
			Description: "Test expense",
			CategoryID:  "food",
			DueDate:     dueDate.Format("2006-01-02T15:04:05Z07:00"),
			Paid:        false,
		}).Return(&expenses.CreateExpenseResponse{
			ID:          "1",
			UserID:      "user1",
			Amount:      100.0,
			Description: "Test expense",
			CategoryID:  "food",
			DueDate:     dueDate.Format("2006-01-02T15:04:05Z07:00"),
			Paid:        false,
		}, nil)

		// Mock para actualizar gasto
		mockExpenseService.On("UpdateExpense", ctx, "user1", "1", &expenses.UpdateExpenseRequest{
			Amount:      200.0,
			Description: "Updated expense",
			CategoryID:  "food",
			DueDate:     dueDate.Format("2006-01-02T15:04:05Z07:00"),
			Paid:        true,
		}).Return(&expenses.UpdateExpenseResponse{
			ID:          "1",
			UserID:      "user1",
			Amount:      200.0,
			Description: "Updated expense",
			CategoryID:  "food",
			DueDate:     dueDate.Format("2006-01-02T15:04:05Z07:00"),
			Paid:        true,
		}, nil)

		// Crear gasto
		expense, err := factory.CreateTransaction(ctx, "user1", 100.0, "Test expense", "food", &dueDate)
		assert.NoError(t, err)
		assert.NotNil(t, expense)
		assert.Equal(t, "1", expense.GetID())
		assert.Equal(t, "user1", expense.GetUserID())
		assert.Equal(t, 100.0, expense.GetAmount())
		assert.Equal(t, "Test expense", expense.GetDescription())
		assert.Equal(t, "food", expense.GetCategoryID())

		// Actualizar gasto
		updatedExpense, err := factory.UpdateTransaction(ctx, "user1", "1", 200.0, "Updated expense", "food", &dueDate, ExpenseType)
		assert.NoError(t, err)
		assert.NotNil(t, updatedExpense)
		assert.Equal(t, "1", updatedExpense.GetID())
		assert.Equal(t, "user1", updatedExpense.GetUserID())
		assert.Equal(t, 200.0, updatedExpense.GetAmount())
		assert.Equal(t, "Updated expense", updatedExpense.GetDescription())
		assert.Equal(t, "food", updatedExpense.GetCategoryID())
	})

	// Caso 2: Crear un ingreso y actualizarlo
	t.Run("Crear y actualizar ingreso", func(t *testing.T) {
		ctx := context.Background()

		// Mock para crear ingreso
		mockIncomeService.On("CreateIncome", ctx, &incomes.CreateIncomeRequest{
			UserID:      "user1",
			Amount:      500.0,
			Description: "Test income",
		}).Return(&incomes.CreateIncomeResponse{
			ID:          "1",
			UserID:      "user1",
			Amount:      500.0,
			Description: "Test income",
		}, nil)

		// Mock para actualizar ingreso
		mockIncomeService.On("UpdateIncome", ctx, "user1", "1", &incomes.UpdateIncomeRequest{
			Amount:      1000.0,
			Description: "Updated income",
		}).Return(&incomes.UpdateIncomeResponse{
			ID:          "1",
			UserID:      "user1",
			Amount:      1000.0,
			Description: "Updated income",
		}, nil)

		// Crear ingreso
		income, err := factory.CreateTransaction(ctx, "user1", 500.0, "Test income", "salary", nil)
		assert.NoError(t, err)
		assert.NotNil(t, income)
		assert.Equal(t, "1", income.GetID())
		assert.Equal(t, "user1", income.GetUserID())
		assert.Equal(t, 500.0, income.GetAmount())
		assert.Equal(t, "Test income", income.GetDescription())

		// Actualizar ingreso
		updatedIncome, err := factory.UpdateTransaction(ctx, "user1", "1", 1000.0, "Updated income", "salary", nil, IncomeType)
		assert.NoError(t, err)
		assert.NotNil(t, updatedIncome)
		assert.Equal(t, "1", updatedIncome.GetID())
		assert.Equal(t, "user1", updatedIncome.GetUserID())
		assert.Equal(t, 1000.0, updatedIncome.GetAmount())
		assert.Equal(t, "Updated income", updatedIncome.GetDescription())
	})
}

func TestTransactionComponent_ListAndDelete(t *testing.T) {
	// Configurar mocks
	mockIncomeService := new(MockIncomeService)
	mockExpenseService := new(MockExpenseService)
	factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

	// Caso 1: Listar y eliminar gastos
	t.Run("Listar y eliminar gastos", func(t *testing.T) {
		ctx := context.Background()

		// Mock para listar gastos
		mockExpenseService.On("ListExpenses", ctx, "user1").Return(&expenses.ListExpensesResponse{
			Expenses: []expenses.GetExpenseResponse{
				{
					ID:          "1",
					UserID:      "user1",
					Amount:      100.0,
					Description: "Expense 1",
					CategoryID:  "food",
				},
				{
					ID:          "2",
					UserID:      "user1",
					Amount:      200.0,
					Description: "Expense 2",
					CategoryID:  "transport",
				},
			},
		}, nil)

		// Mock para eliminar gasto
		mockExpenseService.On("DeleteExpense", ctx, "user1", "1").Return(nil)

		// Listar gastos
		expenses, err := factory.ListTransactions(ctx, "user1", ExpenseType)
		assert.NoError(t, err)
		assert.Len(t, expenses, 2)
		assert.Equal(t, "1", expenses[0].GetID())
		assert.Equal(t, "2", expenses[1].GetID())

		// Eliminar gasto
		err = factory.DeleteTransaction(ctx, "user1", "1", ExpenseType)
		assert.NoError(t, err)
	})

	// Caso 2: Listar y eliminar ingresos
	t.Run("Listar y eliminar ingresos", func(t *testing.T) {
		ctx := context.Background()

		// Mock para listar ingresos
		mockIncomeService.On("ListIncomes", ctx, "user1").Return(&incomes.ListIncomesResponse{
			Incomes: []incomes.GetIncomeResponse{
				{
					ID:          "1",
					UserID:      "user1",
					Amount:      1000.0,
					Description: "Income 1",
				},
				{
					ID:          "2",
					UserID:      "user1",
					Amount:      2000.0,
					Description: "Income 2",
				},
			},
		}, nil)

		// Mock para eliminar ingreso
		mockIncomeService.On("DeleteIncome", ctx, "user1", "1").Return(nil)

		// Listar ingresos
		incomes, err := factory.ListTransactions(ctx, "user1", IncomeType)
		assert.NoError(t, err)
		assert.Len(t, incomes, 2)
		assert.Equal(t, "1", incomes[0].GetID())
		assert.Equal(t, "2", incomes[1].GetID())

		// Eliminar ingreso
		err = factory.DeleteTransaction(ctx, "user1", "1", IncomeType)
		assert.NoError(t, err)
	})
}

func TestTransactionComponent_ErrorHandling(t *testing.T) {
	// Configurar mocks
	mockIncomeService := new(MockIncomeService)
	mockExpenseService := new(MockExpenseService)
	factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

	// Caso 1: Error al crear transacción con monto inválido
	t.Run("Error al crear transacción con monto inválido", func(t *testing.T) {
		ctx := context.Background()

		// Intentar crear transacción con monto negativo
		_, err := factory.CreateTransaction(ctx, "user1", -100.0, "Test", "food", nil)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidAmount, err)

		// Intentar crear transacción con monto muy grande
		_, err = factory.CreateTransaction(ctx, "user1", 1e13, "Test", "food", nil)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrAmountTooLarge, err)
	})

	// Caso 2: Error al actualizar transacción inexistente
	t.Run("Error al actualizar transacción inexistente", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("not found")

		// Mock para actualizar ingreso inexistente
		mockIncomeService.On("UpdateIncome", ctx, "user1", "999", mock.Anything).Return(nil, expectedErr)

		// Intentar actualizar ingreso inexistente
		_, err := factory.UpdateTransaction(ctx, "user1", "999", 100.0, "Test", "", nil, IncomeType)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		// Mock para actualizar gasto inexistente
		mockExpenseService.On("UpdateExpense", ctx, "user1", "999", mock.Anything).Return(nil, expectedErr)

		// Intentar actualizar gasto inexistente
		_, err = factory.UpdateTransaction(ctx, "user1", "999", 100.0, "Test", "food", nil, ExpenseType)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	// Caso 3: Error al eliminar transacción inexistente
	t.Run("Error al eliminar transacción inexistente", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("not found")

		// Mock para eliminar ingreso inexistente
		mockIncomeService.On("DeleteIncome", ctx, "user1", "999").Return(expectedErr)

		// Intentar eliminar ingreso inexistente
		err := factory.DeleteTransaction(ctx, "user1", "999", IncomeType)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		// Mock para eliminar gasto inexistente
		mockExpenseService.On("DeleteExpense", ctx, "user1", "999").Return(expectedErr)

		// Intentar eliminar gasto inexistente
		err = factory.DeleteTransaction(ctx, "user1", "999", ExpenseType)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
