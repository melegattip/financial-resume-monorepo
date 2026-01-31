package transactions

import (
	"context"
	"errors"
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPercentageExpenseService es un mock del servicio de gastos para testing del PercentageObserver
type MockPercentageExpenseService struct {
	mock.Mock
}

func (m *MockPercentageExpenseService) GetTotalIncome(ctx context.Context, userID string) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPercentageExpenseService) GetAllExpenses(ctx context.Context, userID string) ([]domain.Transaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockPercentageExpenseService) UpdateExpense(ctx context.Context, expense domain.Transaction) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockPercentageExpenseService) UpdateExpensePercentages(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPercentageExpenseService) GetTransaction(ctx context.Context, transactionID string) (domain.Transaction, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(domain.Transaction), args.Error(1)
}

func TestPercentageObserver_OnTransactionCreated(t *testing.T) {
	mockService := new(MockPercentageExpenseService)
	observer := NewPercentageObserver(mockService)

	ctx := context.Background()
	userID := "user123"

	// Test con un ingreso (sin categoría)
	income := &domain.Income{
		UserID: userID,
		Amount: 1000.0,
	}

	// Configurar los mocks necesarios para UpdatePercentages
	mockService.On("UpdateExpensePercentages", ctx, userID).Return(nil)

	err := observer.OnTransactionCreated(ctx, income)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)

	// Test con un gasto (con categoría)
	categoryID := "cat1"
	expense := &domain.Expense{
		UserID:     userID,
		Amount:     100.0,
		CategoryID: &categoryID,
	}

	// Configurar los mocks nuevamente para el segundo caso
	mockService.On("UpdateExpensePercentages", ctx, userID).Return(nil)

	err = observer.OnTransactionCreated(ctx, expense)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestPercentageObserver_OnTransactionUpdated(t *testing.T) {
	mockService := new(MockPercentageExpenseService)
	observer := NewPercentageObserver(mockService)

	ctx := context.Background()
	userID := "user123"

	// Test con un ingreso actualizado
	income := &domain.Income{
		UserID: userID,
		Amount: 1500.0,
	}

	// Configurar los mocks necesarios para UpdatePercentages
	mockService.On("UpdateExpensePercentages", ctx, userID).Return(nil)

	err := observer.OnTransactionUpdated(ctx, income)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)

	// Test con un gasto actualizado
	categoryID := "cat1"
	expense := &domain.Expense{
		UserID:     userID,
		Amount:     200.0,
		CategoryID: &categoryID,
	}

	// Configurar los mocks nuevamente para el segundo caso
	mockService.On("UpdateExpensePercentages", ctx, userID).Return(nil)

	err = observer.OnTransactionUpdated(ctx, expense)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestPercentageObserver_OnTransactionDeleted(t *testing.T) {
	mockService := new(MockPercentageExpenseService)
	observer := NewPercentageObserver(mockService)

	ctx := context.Background()
	transactionID := "tx123"
	userID := "user123"

	// Configurar el mock para GetTransaction
	transaction := &domain.Expense{
		UserID: userID,
		Amount: 100.0,
	}
	mockService.On("GetTransaction", ctx, transactionID).Return(transaction, nil)

	// Configurar los mocks necesarios para UpdatePercentages
	mockService.On("UpdateExpensePercentages", ctx, userID).Return(nil)

	err := observer.OnTransactionDeleted(ctx, transactionID)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestPercentageObserver_UpdatePercentages_Error(t *testing.T) {
	mockService := new(MockPercentageExpenseService)
	observer := NewPercentageObserver(mockService)

	ctx := context.Background()
	userID := "user123"

	// Test error al obtener total de ingresos
	expectedError := errors.New("error al obtener total de ingresos")
	mockService.On("UpdateExpensePercentages", ctx, userID).Return(expectedError)

	err := observer.UpdatePercentages(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockService.AssertExpectations(t)

	// Reiniciar el mock para el siguiente test
	mockService = new(MockPercentageExpenseService)
	observer = NewPercentageObserver(mockService)

	// Test error al obtener gastos
	mockService.On("UpdateExpensePercentages", ctx, userID).Return(expectedError)

	err = observer.UpdatePercentages(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockService.AssertExpectations(t)

	// Reiniciar el mock para el siguiente test
	mockService = new(MockPercentageExpenseService)
	observer = NewPercentageObserver(mockService)

	// Test error al actualizar un gasto
	mockService.On("UpdateExpensePercentages", ctx, userID).Return(expectedError)

	err = observer.UpdatePercentages(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockService.AssertExpectations(t)
}
