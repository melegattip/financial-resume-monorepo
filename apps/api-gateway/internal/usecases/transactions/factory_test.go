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

// MockIncomeService es un mock del servicio de ingresos
type MockIncomeService struct {
	mock.Mock
}

func (m *MockIncomeService) CreateIncome(ctx context.Context, request *incomes.CreateIncomeRequest) (*incomes.CreateIncomeResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.CreateIncomeResponse), args.Error(1)
}

func (m *MockIncomeService) GetIncome(ctx context.Context, userID string, incomeID string) (*incomes.GetIncomeResponse, error) {
	args := m.Called(ctx, userID, incomeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.GetIncomeResponse), args.Error(1)
}

func (m *MockIncomeService) ListIncomes(ctx context.Context, userID string) (*incomes.ListIncomesResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.ListIncomesResponse), args.Error(1)
}

func (m *MockIncomeService) UpdateIncome(ctx context.Context, userID string, incomeID string, request *incomes.UpdateIncomeRequest) (*incomes.UpdateIncomeResponse, error) {
	args := m.Called(ctx, userID, incomeID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.UpdateIncomeResponse), args.Error(1)
}

func (m *MockIncomeService) DeleteIncome(ctx context.Context, userID string, incomeID string) error {
	args := m.Called(ctx, userID, incomeID)
	return args.Error(0)
}

func (m *MockIncomeService) ListUnreceivedIncomes(ctx context.Context, userID string) (*incomes.ListIncomesResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.ListIncomesResponse), args.Error(1)
}

type MockExpenseService struct {
	mock.Mock
}

func (m *MockExpenseService) CreateExpense(ctx context.Context, request *expenses.CreateExpenseRequest) (*expenses.CreateExpenseResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expenses.CreateExpenseResponse), args.Error(1)
}

func (m *MockExpenseService) GetExpense(ctx context.Context, userID string, expenseID string) (*expenses.GetExpenseResponse, error) {
	args := m.Called(ctx, userID, expenseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expenses.GetExpenseResponse), args.Error(1)
}

func (m *MockExpenseService) ListExpenses(ctx context.Context, userID string) (*expenses.ListExpensesResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expenses.ListExpensesResponse), args.Error(1)
}

func (m *MockExpenseService) ListUnpaidExpenses(ctx context.Context, userID string) (*expenses.ListExpensesResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expenses.ListExpensesResponse), args.Error(1)
}

func (m *MockExpenseService) UpdateExpense(ctx context.Context, userID string, expenseID string, request *expenses.UpdateExpenseRequest) (*expenses.UpdateExpenseResponse, error) {
	args := m.Called(ctx, userID, expenseID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expenses.UpdateExpenseResponse), args.Error(1)
}

func (m *MockExpenseService) DeleteExpense(ctx context.Context, userID string, expenseID string) error {
	args := m.Called(ctx, userID, expenseID)
	return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
	mockIncomeService := new(MockIncomeService)
	mockExpenseService := new(MockExpenseService)
	factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

	// Test crear un gasto
	dueDate := time.Now().AddDate(0, 1, 0)
	expenseRequest := &expenses.CreateExpenseRequest{
		UserID:      "user1",
		Amount:      100.0,
		Description: "Test expense",
		CategoryID:  "food",
		Paid:        false,
		DueDate:     dueDate.Format("2006-01-02T15:04:05Z07:00"),
	}

	expenseResponse := &expenses.CreateExpenseResponse{
		ID:          "1",
		UserID:      "user1",
		Amount:      100.0,
		Description: "Test expense",
		CategoryID:  "food",
		Paid:        false,
		DueDate:     dueDate.Format("2006-01-02T15:04:05Z07:00"),
	}

	mockExpenseService.On("CreateExpense", context.Background(), expenseRequest).Return(expenseResponse, nil)

	result, err := factory.CreateTransaction(context.Background(), "user1", 100.0, "Test expense", "food", &dueDate)
	assert.NoError(t, err)
	assert.Equal(t, expenseResponse.ID, result.GetID())
	assert.Equal(t, expenseResponse.UserID, result.GetUserID())
	assert.Equal(t, expenseResponse.Amount, result.GetAmount())
	assert.Equal(t, expenseResponse.Description, result.GetDescription())
	assert.Equal(t, expenseResponse.CategoryID, result.GetCategoryID())
	mockExpenseService.AssertExpectations(t)

	// Test crear un ingreso
	incomeRequest := &incomes.CreateIncomeRequest{
		UserID:      "user1",
		Amount:      200.0,
		Description: "Test income",
	}

	incomeResponse := &incomes.CreateIncomeResponse{
		ID:          "2",
		UserID:      "user1",
		Amount:      200.0,
		Description: "Test income",
	}

	mockIncomeService.On("CreateIncome", context.Background(), incomeRequest).Return(incomeResponse, nil)

	result, err = factory.CreateTransaction(context.Background(), "user1", 200.0, "Test income", "salary", nil)
	assert.NoError(t, err)
	assert.Equal(t, incomeResponse.ID, result.GetID())
	assert.Equal(t, incomeResponse.UserID, result.GetUserID())
	assert.Equal(t, incomeResponse.Amount, result.GetAmount())
	assert.Equal(t, incomeResponse.Description, result.GetDescription())
	mockIncomeService.AssertExpectations(t)
}

func TestGetTransaction(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		transactionID   string
		transactionType TransactionType
		setupMocks      func(*MockIncomeService, *MockExpenseService)
		expectedErr     error
	}{
		{
			name:            "Tipo de transacción inválido",
			userID:          "user1",
			transactionID:   "1",
			transactionType: "invalid",
			setupMocks:      func(income *MockIncomeService, expense *MockExpenseService) {},
			expectedErr:     domain.ErrInvalidTransactionType,
		},
		{
			name:            "Error al obtener ingreso",
			userID:          "user1",
			transactionID:   "1",
			transactionType: IncomeType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				income.On("GetIncome", mock.Anything, "user1", "1").Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			name:            "Error al obtener gasto",
			userID:          "user1",
			transactionID:   "1",
			transactionType: ExpenseType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				expense.On("GetExpense", mock.Anything, "user1", "1").Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIncomeService := new(MockIncomeService)
			mockExpenseService := new(MockExpenseService)
			factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

			tt.setupMocks(mockIncomeService, mockExpenseService)

			_, err := factory.GetTransaction(context.Background(), tt.userID, tt.transactionID, tt.transactionType)
			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListTransactions(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		transactionType TransactionType
		setupMocks      func(*MockIncomeService, *MockExpenseService)
		expectedLen     int
		expectedErr     error
	}{
		{
			name:            "Lista vacía de ingresos",
			userID:          "user1",
			transactionType: IncomeType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				income.On("ListIncomes", mock.Anything, "user1").Return(&incomes.ListIncomesResponse{
					Incomes: []incomes.GetIncomeResponse{},
				}, nil)
			},
			expectedLen: 0,
			expectedErr: nil,
		},
		{
			name:            "Lista vacía de gastos",
			userID:          "user1",
			transactionType: ExpenseType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				expense.On("ListExpenses", mock.Anything, "user1").Return(&expenses.ListExpensesResponse{
					Expenses: []expenses.GetExpenseResponse{},
				}, nil)
			},
			expectedLen: 0,
			expectedErr: nil,
		},
		{
			name:            "Tipo de transacción inválido",
			userID:          "user1",
			transactionType: "invalid",
			setupMocks:      func(income *MockIncomeService, expense *MockExpenseService) {},
			expectedLen:     0,
			expectedErr:     domain.ErrInvalidTransactionType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIncomeService := new(MockIncomeService)
			mockExpenseService := new(MockExpenseService)
			factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

			tt.setupMocks(mockIncomeService, mockExpenseService)

			result, err := factory.ListTransactions(context.Background(), tt.userID, tt.transactionType)
			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}
		})
	}
}

func TestUpdateTransaction(t *testing.T) {
	dueDate := time.Now().AddDate(0, 1, 0)
	tests := []struct {
		name            string
		userID          string
		transactionID   string
		amount          float64
		description     string
		categoryID      string
		dueDate         *time.Time
		transactionType TransactionType
		setupMocks      func(*MockIncomeService, *MockExpenseService)
		expectedErr     error
	}{
		{
			name:            "Actualizar ingreso exitosamente",
			userID:          "user1",
			transactionID:   "1",
			amount:          100.0,
			description:     "Updated income",
			categoryID:      "",
			dueDate:         nil,
			transactionType: IncomeType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				income.On("UpdateIncome", mock.Anything, "user1", "1", &incomes.UpdateIncomeRequest{
					Amount:      100.0,
					Description: "Updated income",
				}).Return(&incomes.UpdateIncomeResponse{
					ID:          "1",
					UserID:      "user1",
					Amount:      100.0,
					Description: "Updated income",
				}, nil)
			},
			expectedErr: nil,
		},
		{
			name:            "Actualizar gasto exitosamente",
			userID:          "user1",
			transactionID:   "1",
			amount:          200.0,
			description:     "Updated expense",
			categoryID:      "food",
			dueDate:         &dueDate,
			transactionType: ExpenseType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				expense.On("UpdateExpense", mock.Anything, "user1", "1", &expenses.UpdateExpenseRequest{
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
			},
			expectedErr: nil,
		},
		{
			name:            "Error al actualizar ingreso",
			userID:          "user1",
			transactionID:   "1",
			amount:          100.0,
			description:     "Updated income",
			categoryID:      "",
			dueDate:         nil,
			transactionType: IncomeType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				income.On("UpdateIncome", mock.Anything, "user1", "1", &incomes.UpdateIncomeRequest{
					Amount:      100.0,
					Description: "Updated income",
				}).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			name:            "Error al actualizar gasto",
			userID:          "user1",
			transactionID:   "1",
			amount:          200.0,
			description:     "Updated expense",
			categoryID:      "food",
			dueDate:         &dueDate,
			transactionType: ExpenseType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				expense.On("UpdateExpense", mock.Anything, "user1", "1", &expenses.UpdateExpenseRequest{
					Amount:      200.0,
					Description: "Updated expense",
					CategoryID:  "food",
					DueDate:     dueDate.Format("2006-01-02T15:04:05Z07:00"),
					Paid:        true,
				}).Return(nil, errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			name:            "Tipo de transacción inválido",
			userID:          "user1",
			transactionID:   "1",
			amount:          100.0,
			description:     "Invalid",
			categoryID:      "",
			dueDate:         nil,
			transactionType: "invalid",
			setupMocks:      func(income *MockIncomeService, expense *MockExpenseService) {},
			expectedErr:     domain.ErrInvalidTransactionType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIncomeService := new(MockIncomeService)
			mockExpenseService := new(MockExpenseService)
			factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

			tt.setupMocks(mockIncomeService, mockExpenseService)

			_, err := factory.UpdateTransaction(context.Background(), tt.userID, tt.transactionID, tt.amount, tt.description, tt.categoryID, tt.dueDate, tt.transactionType)
			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteTransaction(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		transactionID   string
		transactionType TransactionType
		setupMocks      func(*MockIncomeService, *MockExpenseService)
		expectedErr     error
	}{
		{
			name:            "Eliminar ingreso exitosamente",
			userID:          "user1",
			transactionID:   "1",
			transactionType: IncomeType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				income.On("DeleteIncome", mock.Anything, "user1", "1").Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:            "Eliminar gasto exitosamente",
			userID:          "user1",
			transactionID:   "1",
			transactionType: ExpenseType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				expense.On("DeleteExpense", mock.Anything, "user1", "1").Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:            "Error al eliminar ingreso",
			userID:          "user1",
			transactionID:   "1",
			transactionType: IncomeType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				income.On("DeleteIncome", mock.Anything, "user1", "1").Return(errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			name:            "Error al eliminar gasto",
			userID:          "user1",
			transactionID:   "1",
			transactionType: ExpenseType,
			setupMocks: func(income *MockIncomeService, expense *MockExpenseService) {
				expense.On("DeleteExpense", mock.Anything, "user1", "1").Return(errors.New("error"))
			},
			expectedErr: errors.New("error"),
		},
		{
			name:            "Tipo de transacción inválido",
			userID:          "user1",
			transactionID:   "1",
			transactionType: "invalid",
			setupMocks:      func(income *MockIncomeService, expense *MockExpenseService) {},
			expectedErr:     domain.ErrInvalidTransactionType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockIncomeService := new(MockIncomeService)
			mockExpenseService := new(MockExpenseService)
			factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

			tt.setupMocks(mockIncomeService, mockExpenseService)

			err := factory.DeleteTransaction(context.Background(), tt.userID, tt.transactionID, tt.transactionType)
			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateTransactionError(t *testing.T) {
	mockIncomeService := new(MockIncomeService)
	mockExpenseService := new(MockExpenseService)
	factory := NewTransactionFactory(mockIncomeService, mockExpenseService)

	// Test error al crear un gasto
	dueDate := time.Now().AddDate(0, 1, 0)
	expenseRequest := &expenses.CreateExpenseRequest{
		UserID:      "user1",
		Amount:      100.0,
		Description: "Test expense",
		CategoryID:  "food",
		Paid:        false,
		DueDate:     dueDate.Format("2006-01-02T15:04:05Z07:00"),
	}

	expectedError := errors.New("error al crear gasto")
	mockExpenseService.On("CreateExpense", context.Background(), expenseRequest).Return(nil, expectedError)

	result, err := factory.CreateTransaction(context.Background(), "user1", 100.0, "Test expense", "food", &dueDate)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	mockExpenseService.AssertExpectations(t)

	// Test error al crear un ingreso
	incomeRequest := &incomes.CreateIncomeRequest{
		UserID:      "user1",
		Amount:      200.0,
		Description: "Test income",
	}

	expectedError = errors.New("error al crear ingreso")
	mockIncomeService.On("CreateIncome", context.Background(), incomeRequest).Return(nil, expectedError)

	result, err = factory.CreateTransaction(context.Background(), "user1", 200.0, "Test income", "salary", nil)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	mockIncomeService.AssertExpectations(t)
}

func TestHelperFunctions(t *testing.T) {
	// Test toTransaction con tipo inválido
	result := toTransaction("invalid")
	assert.Nil(t, result)

	// Test toTransactionSlice con tipo inválido
	results := toTransactionSlice("invalid")
	assert.Empty(t, results)
}
