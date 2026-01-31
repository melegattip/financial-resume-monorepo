package create

import (
	"context"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExpenseRepository es un mock del repositorio de gastos
type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) Create(expense *domain.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) Get(userID, expenseID string) (*domain.Expense, error) {
	args := m.Called(userID, expenseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepository) List(userID string) ([]*domain.Expense, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepository) Update(expense *domain.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) Delete(userID, expenseID string) error {
	args := m.Called(userID, expenseID)
	return args.Error(0)
}

func (m *MockExpenseRepository) ListUnpaid(userID string) ([]*domain.Expense, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Expense), args.Error(1)
}

// MockCategoryRepository es un mock del repositorio de categorías
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Get(userID, categoryID string) (*domain.Category, error) {
	args := m.Called(userID, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByName(userID, name string) (*domain.Category, error) {
	args := m.Called(userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) List(userID string) ([]*domain.Category, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(userID, categoryID string) error {
	args := m.Called(userID, categoryID)
	return args.Error(0)
}

func (m *MockCategoryRepository) UpdateCategoryReferences(categoryID string) error {
	args := m.Called(categoryID)
	return args.Error(0)
}

// MockPercentageObserver es un mock del PercentageTransactionObserver
type MockPercentageObserver struct {
	mock.Mock
}

func (m *MockPercentageObserver) OnTransactionCreated(ctx context.Context, transaction domain.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockPercentageObserver) OnTransactionUpdated(ctx context.Context, transaction domain.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockPercentageObserver) OnTransactionDeleted(ctx context.Context, transactionID string) error {
	args := m.Called(ctx, transactionID)
	return args.Error(0)
}

func (m *MockPercentageObserver) GetTotalIncome(ctx context.Context, userID string) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

// TestCreateExpense_Contract verifica que el contrato de creación de gastos se cumpla
func TestCreateExpense_Contract(t *testing.T) {
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

	// Ejecutar el servicio
	response, err := service.CreateExpense(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Verificar el contrato de la respuesta
	assert.Equal(t, request.UserID, response.UserID)
	assert.Equal(t, request.Amount, response.Amount)
	assert.Equal(t, request.Description, response.Description)
	assert.Equal(t, request.CategoryID, response.CategoryID)
	assert.Equal(t, request.DueDate, response.DueDate)
	assert.Equal(t, 10.0, response.Percentage) // Verificar que el porcentaje se incluye en la respuesta
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)

	// Verificar que se llamaron todos los mocks
	mockExpenseRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
	mockPercentageObserver.AssertExpectations(t)
}
