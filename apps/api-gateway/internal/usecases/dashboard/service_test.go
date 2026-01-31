package dashboard

import (
	"context"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) List(userID string) ([]*domain.Expense, error) {
	args := m.Called(userID)
	return args.Get(0).([]*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepository) Get(userID, id string) (*domain.Expense, error) {
	args := m.Called(userID, id)
	return args.Get(0).(*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepository) Create(expense *domain.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) Update(expense *domain.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) Delete(userID, id string) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

func (m *MockExpenseRepository) ListUnpaid(userID string) ([]*domain.Expense, error) {
	args := m.Called(userID)
	return args.Get(0).([]*domain.Expense), args.Error(1)
}

type MockIncomeRepository struct {
	mock.Mock
}

func (m *MockIncomeRepository) List(userID string) ([]*domain.Income, error) {
	args := m.Called(userID)
	return args.Get(0).([]*domain.Income), args.Error(1)
}

func (m *MockIncomeRepository) Get(userID, id string) (*domain.Income, error) {
	args := m.Called(userID, id)
	return args.Get(0).(*domain.Income), args.Error(1)
}

func (m *MockIncomeRepository) Create(income *domain.Income) error {
	args := m.Called(income)
	return args.Error(0)
}

func (m *MockIncomeRepository) Update(income *domain.Income) error {
	args := m.Called(income)
	return args.Error(0)
}

func (m *MockIncomeRepository) Delete(userID, id string) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

type MockPeriodCalculator struct {
	mock.Mock
}

func (m *MockPeriodCalculator) FilterTransactionsByPeriod(transactions []usecases.Transaction, period usecases.DatePeriod) []usecases.Transaction {
	args := m.Called(transactions, period)
	return args.Get(0).([]usecases.Transaction)
}

func (m *MockPeriodCalculator) CalculateMetrics(transactions []usecases.Transaction) usecases.FinancialMetrics {
	args := m.Called(transactions)
	return args.Get(0).(usecases.FinancialMetrics)
}

func (m *MockPeriodCalculator) FormatPeriodLabel(period usecases.DatePeriod) string {
	args := m.Called(period)
	return args.String(0)
}

type MockAnalyticsCalculator struct {
	mock.Mock
}

func (m *MockAnalyticsCalculator) CalculatePercentages(amount, total float64) float64 {
	args := m.Called(amount, total)
	return args.Get(0).(float64)
}

func (m *MockAnalyticsCalculator) CalculateAverage(total float64, count int) float64 {
	args := m.Called(total, count)
	return args.Get(0).(float64)
}

func (m *MockAnalyticsCalculator) GenerateColorSeed(identifier string) int {
	args := m.Called(identifier)
	return args.Int(0)
}

func TestService_GetDashboardOverview_Success(t *testing.T) {
	// Preparar mocks
	mockExpenseRepo := &MockExpenseRepository{}
	mockIncomeRepo := &MockIncomeRepository{}
	mockPeriodCalc := &MockPeriodCalculator{}
	mockAnalyticsCalc := &MockAnalyticsCalculator{}

	// Datos de prueba
	userID := "user-123"
	year := 2024
	month := 12
	categoryID := "cat-123"

	expenses := []*domain.Expense{
		{
			ID:         "exp-1",
			UserID:     userID,
			Amount:     100.0,
			CategoryID: &categoryID,
			Paid:       false,
			CreatedAt:  time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	incomes := []*domain.Income{
		{
			ID:         "inc-1",
			UserID:     userID,
			Amount:     200.0,
			CategoryID: &categoryID,
			CreatedAt:  time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	expectedMetrics := usecases.FinancialMetrics{
		TotalIncome:          200.0,
		TotalExpenses:        100.0,
		Balance:              100.0,
		PendingExpenses:      100.0,
		PendingExpensesCount: 1,
	}

	// Configurar mocks
	mockExpenseRepo.On("List", userID).Return(expenses, nil)
	mockIncomeRepo.On("List", userID).Return(incomes, nil)
	mockPeriodCalc.On("FilterTransactionsByPeriod", mock.Anything, mock.Anything).Return([]usecases.Transaction{})
	mockPeriodCalc.On("CalculateMetrics", mock.Anything).Return(expectedMetrics)
	mockPeriodCalc.On("FormatPeriodLabel", mock.Anything).Return("Diciembre 2024")
	mockAnalyticsCalc.On("CalculatePercentages", 100.0, 200.0).Return(50.0)

	// Crear servicio
	service := NewService(mockExpenseRepo, mockIncomeRepo, mockPeriodCalc, mockAnalyticsCalc)

	// Ejecutar
	params := usecases.DashboardParams{
		UserID: userID,
		Period: usecases.DatePeriod{
			Year:  &year,
			Month: &month,
		},
	}

	result, err := service.GetDashboardOverview(context.Background(), params)

	// Verificar
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Diciembre 2024", result.Period.Label)
	assert.Equal(t, expectedMetrics.TotalIncome, result.Metrics.TotalIncome)
	assert.Equal(t, expectedMetrics.TotalExpenses, result.Metrics.TotalExpenses)
	assert.Equal(t, expectedMetrics.Balance, result.Metrics.Balance)

	// Verificar que todos los mocks fueron llamados
	mockExpenseRepo.AssertExpectations(t)
	mockIncomeRepo.AssertExpectations(t)
	mockPeriodCalc.AssertExpectations(t)
	mockAnalyticsCalc.AssertExpectations(t)
}

func TestService_GetDashboardOverview_ValidationError(t *testing.T) {
	// Preparar mocks
	mockExpenseRepo := &MockExpenseRepository{}
	mockIncomeRepo := &MockIncomeRepository{}
	mockPeriodCalc := &MockPeriodCalculator{}
	mockAnalyticsCalc := &MockAnalyticsCalculator{}

	// Crear servicio
	service := NewService(mockExpenseRepo, mockIncomeRepo, mockPeriodCalc, mockAnalyticsCalc)

	// Ejecutar con parámetros inválidos
	params := usecases.DashboardParams{
		UserID: "", // UserID vacío debería fallar
		Period: usecases.DatePeriod{},
	}

	result, err := service.GetDashboardOverview(context.Background(), params)

	// Verificar
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "El ID del usuario es requerido")
}

func TestService_validateParams(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name    string
		params  usecases.DashboardParams
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid params",
			params: usecases.DashboardParams{
				UserID: "user-123",
				Period: usecases.DatePeriod{},
			},
			wantErr: false,
		},
		{
			name: "empty user ID",
			params: usecases.DashboardParams{
				UserID: "",
				Period: usecases.DatePeriod{},
			},
			wantErr: true,
			errMsg:  "El ID del usuario es requerido",
		},
		{
			name: "invalid year",
			params: usecases.DashboardParams{
				UserID: "user-123",
				Period: usecases.DatePeriod{
					Year: func() *int { y := 1800; return &y }(),
				},
			},
			wantErr: true,
			errMsg:  "Año inválido",
		},
		{
			name: "invalid month",
			params: usecases.DashboardParams{
				UserID: "user-123",
				Period: usecases.DatePeriod{
					Month: func() *int { m := 13; return &m }(),
				},
			},
			wantErr: true,
			errMsg:  "Mes inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateParams(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
