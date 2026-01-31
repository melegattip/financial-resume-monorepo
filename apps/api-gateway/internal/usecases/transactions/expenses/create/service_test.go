package create

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ErrNotFound = errors.New("not found")

// Asegurarse de que MockPercentageObserver implementa la interfaz PercentageTransactionObserver
var _ transactions.PercentageTransactionObserver = (*MockPercentageObserver)(nil)

func TestCreateExpense(t *testing.T) {
	mockExpenseRepo := new(MockExpenseRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockPercentageObserver := new(MockPercentageObserver)

	service := NewService(mockExpenseRepo, mockCategoryRepo, mockPercentageObserver)

	ctx := context.Background()
	request := &expenses.CreateExpenseRequest{
		UserID:      "user123",
		Amount:      100.0,
		Description: "Test Expense",
		CategoryID:  "cat1",
		DueDate:     "2024-04-16",
	}

	// Configurar el mock de la categoría
	mockCategoryRepo.On("Get", request.UserID, request.CategoryID).Return(&domain.Category{}, nil)

	// Configurar el mock del PercentageObserver para GetTotalIncome
	mockPercentageObserver.On("GetTotalIncome", ctx, request.UserID).Return(1000.0, nil)

	// Configurar el mock del repositorio de gastos
	mockExpenseRepo.On("Create", mock.MatchedBy(func(exp *domain.Expense) bool {
		return exp.UserID == request.UserID &&
			exp.Amount == request.Amount &&
			exp.Description == request.Description &&
			exp.CategoryID != nil && *exp.CategoryID == request.CategoryID &&
			exp.DueDate.Equal(time.Date(2024, 4, 16, 0, 0, 0, 0, time.UTC)) &&
			exp.Percentage == 10.0 // 100.0 / 1000.0 = 10%
	})).Return(nil)

	// Configurar el mock del PercentageObserver para OnTransactionCreated
	mockPercentageObserver.On("OnTransactionCreated", ctx, mock.MatchedBy(func(exp *domain.Expense) bool {
		return exp.UserID == request.UserID &&
			exp.Amount == request.Amount &&
			exp.Description == request.Description &&
			exp.CategoryID != nil && *exp.CategoryID == request.CategoryID &&
			exp.DueDate.Equal(time.Date(2024, 4, 16, 0, 0, 0, 0, time.UTC)) &&
			exp.Percentage == 10.0
	})).Return(nil)

	response, err := service.CreateExpense(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.UserID, response.UserID)
	assert.Equal(t, request.Amount, response.Amount)
	assert.Equal(t, request.Description, response.Description)
	assert.Equal(t, request.CategoryID, response.CategoryID)
	assert.Equal(t, request.DueDate, response.DueDate)
	assert.Equal(t, 10.0, response.Percentage)

	mockExpenseRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
	mockPercentageObserver.AssertExpectations(t)
}

func TestCreateExpense_Validation(t *testing.T) {
	mockExpenseRepo := new(MockExpenseRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockPercentageObserver := new(MockPercentageObserver)

	service := NewService(mockExpenseRepo, mockCategoryRepo, mockPercentageObserver)

	ctx := context.Background()

	tests := []struct {
		name    string
		request *expenses.CreateExpenseRequest
		errMsg  string
	}{
		{
			name: "Empty description",
			request: &expenses.CreateExpenseRequest{
				UserID:      "user123",
				Amount:      100.0,
				Description: "",
				CategoryID:  "cat1",
			},
			errMsg: "El nombre del gasto es requerido",
		},
		{
			name: "Zero amount",
			request: &expenses.CreateExpenseRequest{
				UserID:      "user123",
				Amount:      0,
				Description: "Test Expense",
				CategoryID:  "cat1",
			},
			errMsg: "El monto del gasto debe ser mayor a 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.CreateExpense(ctx, tt.request)
			assert.Nil(t, response)
			assert.Error(t, err)
			assert.Equal(t, tt.errMsg, err.Error())
		})
	}
}

func TestCreateExpense_InvalidCategory(t *testing.T) {
	mockExpenseRepo := new(MockExpenseRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockPercentageObserver := new(MockPercentageObserver)

	service := NewService(mockExpenseRepo, mockCategoryRepo, mockPercentageObserver)

	ctx := context.Background()
	request := &expenses.CreateExpenseRequest{
		UserID:      "user123",
		Amount:      100.0,
		Description: "Test Expense",
		CategoryID:  "invalid_cat",
	}

	// Configurar el mock de la categoría para retornar error
	mockCategoryRepo.On("Get", request.UserID, request.CategoryID).Return(nil, ErrNotFound)

	response, err := service.CreateExpense(ctx, request)
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Equal(t, "La categoría especificada no existe", err.Error())

	mockCategoryRepo.AssertExpectations(t)
}
