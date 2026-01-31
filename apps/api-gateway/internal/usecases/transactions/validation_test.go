package transactions

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransaction_Validation(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		amount      float64
		description string
		categoryID  string
		dueDate     *time.Time
		expectedErr error
		setupMocks  func(*MockIncomeService, *MockExpenseService)
	}{
		{
			name:        "Monto negativo",
			userID:      "user1",
			amount:      -100.0,
			description: "Test expense",
			categoryID:  "food",
			dueDate:     nil,
			expectedErr: domain.ErrInvalidAmount,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				// No necesitamos configurar mocks ya que la validación debe fallar antes
			},
		},
		{
			name:        "Descripción vacía",
			userID:      "user1",
			amount:      100.0,
			description: "",
			categoryID:  "food",
			dueDate:     nil,
			expectedErr: domain.ErrEmptyDescription,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				// No necesitamos configurar mocks ya que la validación debe fallar antes
			},
		},
		{
			name:        "Categoría inválida",
			userID:      "user1",
			amount:      100.0,
			description: "Test expense",
			categoryID:  "",
			dueDate:     nil,
			expectedErr: domain.ErrInvalidCategory,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				// No necesitamos configurar mocks ya que la validación debe fallar antes
			},
		},
		{
			name:        "Fecha inválida",
			userID:      "user1",
			amount:      100.0,
			description: "Test expense",
			categoryID:  "food",
			dueDate:     &time.Time{}, // Zero time
			expectedErr: domain.ErrInvalidDate,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				// No necesitamos configurar mocks ya que la validación debe fallar antes
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIncomeService := new(MockIncomeService)
			mockExpenseService := new(MockExpenseService)
			factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

			// Configurar los mocks según el caso de prueba
			tt.setupMocks(mockIncomeService, mockExpenseService)

			_, err := factory.CreateTransaction(context.Background(), tt.userID, tt.amount, tt.description, tt.categoryID, tt.dueDate)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestCreateTransaction_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		amount      float64
		description string
		categoryID  string
		dueDate     *time.Time
		expectedErr error
		setupMocks  func(*MockIncomeService, *MockExpenseService)
	}{
		{
			name:        "Monto muy grande",
			amount:      1e18, // 1 quintillón
			description: "Test expense",
			categoryID:  "food",
			dueDate:     nil,
			expectedErr: domain.ErrAmountTooLarge,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				// No necesitamos configurar mocks ya que la validación debe fallar antes
			},
		},
		{
			name:        "Fecha límite",
			amount:      100.0,
			description: "Test expense",
			categoryID:  "food",
			dueDate:     func() *time.Time { t := time.Now().AddDate(100, 0, 0); return &t }(),
			expectedErr: domain.ErrInvalidDate,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				// No necesitamos configurar mocks ya que la validación debe fallar antes
			},
		},
		{
			name:        "Categoría especial",
			amount:      100.0,
			description: "Test expense",
			categoryID:  "special_category",
			dueDate:     nil,
			expectedErr: nil,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				// Configurar el mock para una transacción exitosa
				income.On("CreateIncome", mock.Anything, mock.Anything).Return(&incomes.CreateIncomeResponse{
					ID:          "1",
					UserID:      "user1",
					Amount:      100.0,
					Description: "Test expense",
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIncomeService := new(MockIncomeService)
			mockExpenseService := new(MockExpenseService)
			factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

			// Configurar los mocks según el caso de prueba
			tt.setupMocks(mockIncomeService, mockExpenseService)

			_, err := factory.CreateTransaction(context.Background(), "user1", tt.amount, tt.description, tt.categoryID, tt.dueDate)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateTransaction_PercentageCalculation(t *testing.T) {
	mockIncomeService := new(MockIncomeService)
	mockExpenseService := new(MockExpenseService)
	factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

	// Configurar el mock para simular un error al calcular porcentajes
	mockIncomeService.On("CreateIncome", mock.Anything, mock.Anything).Return(nil, errors.New("error calculating percentages"))

	_, err := factory.CreateTransaction(context.Background(), "user1", 100.0, "Test expense", "food", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error calculating percentages")
	mockIncomeService.AssertExpectations(t)
}

func TestCreateTransaction_ObserverNotification(t *testing.T) {
	mockIncomeService := new(MockIncomeService)
	mockExpenseService := new(MockExpenseService)
	factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

	// Configurar el mock para simular un error al notificar observadores
	mockIncomeService.On("CreateIncome", mock.Anything, mock.Anything).Return(nil, errors.New("error notifying observers"))

	_, err := factory.CreateTransaction(context.Background(), "user1", 100.0, "Test expense", "food", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error notifying observers")
	mockIncomeService.AssertExpectations(t)
}

func TestCreateTransaction_StateUpdate(t *testing.T) {
	mockIncomeService := new(MockIncomeService)
	mockExpenseService := new(MockExpenseService)
	factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

	// Configurar el mock para simular un error al actualizar estados
	mockIncomeService.On("CreateIncome", mock.Anything, mock.Anything).Return(nil, errors.New("error updating state"))

	_, err := factory.CreateTransaction(context.Background(), "user1", 100.0, "Test expense", "food", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error updating state")
	mockIncomeService.AssertExpectations(t)
}
