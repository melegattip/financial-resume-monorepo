package update

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

// TestUpdateExpense_Contract verifica que el contrato de actualización de gastos se cumpla
func TestUpdateExpense_Contract(t *testing.T) {
	categoryID := "cat1"
	mockExpenseRepo := new(MockExpenseRepository)
	mockPercentageObserver := new(MockPercentageObserver)

	service := NewService(mockExpenseRepo, mockPercentageObserver)

	ctx := context.Background()
	userID := "user123"
	expenseID := "exp123"
	dueDate := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	request := &expenses.UpdateExpenseRequest{
		Amount:      150.0,
		Description: "Updated Expense",
		CategoryID:  categoryID,
		DueDate:     dueDate.Format("2006-01-02"),
		Paid:        true,
	}

	// Configurar el mock del repositorio para Get
	existingExpense := &domain.Expense{
		ID:          expenseID,
		UserID:      userID,
		Amount:      100.0,
		Description: "Original Expense",
		CategoryID:  &categoryID,
		DueDate:     time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		Paid:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockExpenseRepo.On("Get", userID, expenseID).Return(existingExpense, nil)

	// Configurar el mock del PercentageObserver para GetTotalIncome
	mockPercentageObserver.On("GetTotalIncome", ctx, userID).Return(1000.0, nil)

	// Configurar el mock del repositorio para Update
	mockExpenseRepo.On("Update", mock.MatchedBy(func(exp *domain.Expense) bool {
		return exp.ID == expenseID &&
			exp.UserID == userID &&
			exp.Amount == request.Amount &&
			exp.Description == request.Description &&
			exp.CategoryID != nil && *exp.CategoryID == request.CategoryID &&
			exp.DueDate.Equal(dueDate) &&
			exp.Paid == request.Paid &&
			exp.Percentage == 15.0 // 150.0 / 1000.0 = 15%
	})).Return(nil)

	// Configurar el mock del PercentageObserver para OnTransactionUpdated
	mockPercentageObserver.On("OnTransactionUpdated", ctx, mock.MatchedBy(func(exp *domain.Expense) bool {
		return exp.ID == expenseID &&
			exp.UserID == userID &&
			exp.Amount == request.Amount &&
			exp.Description == request.Description &&
			exp.CategoryID != nil && *exp.CategoryID == request.CategoryID &&
			exp.DueDate.Equal(dueDate) &&
			exp.Paid == request.Paid &&
			exp.Percentage == 15.0
	})).Return(nil)

	// Ejecutar el servicio
	response, err := service.UpdateExpense(ctx, userID, expenseID, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Verificar el contrato de la respuesta
	assert.Equal(t, expenseID, response.ID)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, request.Amount, response.Amount)
	assert.Equal(t, request.Description, response.Description)
	assert.Equal(t, request.CategoryID, response.CategoryID)
	assert.Equal(t, request.DueDate, response.DueDate)
	assert.Equal(t, request.Paid, response.Paid)
	assert.Equal(t, 15.0, response.Percentage) // Verificar que el porcentaje se incluye en la respuesta
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)

	// Verificar que se llamaron todos los mocks
	mockExpenseRepo.AssertExpectations(t)
	mockPercentageObserver.AssertExpectations(t)
}
