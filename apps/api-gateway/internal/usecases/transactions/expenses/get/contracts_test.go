package get

import (
	"context"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
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

// TestGetExpense_Contract verifica que el contrato de obtención de gastos se cumpla
func TestGetExpense_Contract(t *testing.T) {
	mockExpenseRepo := new(MockExpenseRepository)
	service := NewService(mockExpenseRepo)

	ctx := context.Background()
	userID := "user123"
	expenseID := "exp123"

	// Configurar el mock del repositorio para Get
	categoryID := "cat1"
	expectedExpense := &domain.Expense{
		ID:          expenseID,
		UserID:      userID,
		Amount:      150.0,
		Description: "Test Expense",
		CategoryID:  &categoryID,
		DueDate:     time.Date(2024, 4, 17, 0, 0, 0, 0, time.UTC),
		Paid:        false,
		Percentage:  15.0, // 150.0 / 1000.0 = 15%
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockExpenseRepo.On("Get", userID, expenseID).Return(expectedExpense, nil)

	// Ejecutar el servicio
	response, err := service.GetExpense(ctx, userID, expenseID)
	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Verificar el contrato de la respuesta
	assert.Equal(t, expenseID, response.ID)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, expectedExpense.Amount, response.Amount)
	assert.Equal(t, expectedExpense.Description, response.Description)
	assert.Equal(t, expectedExpense.GetCategoryID(), response.CategoryID)
	assert.Equal(t, expectedExpense.DueDate.Format("2006-01-02"), response.DueDate)
	assert.Equal(t, expectedExpense.Paid, response.Paid)
	assert.Equal(t, expectedExpense.Percentage, response.Percentage) // Verificar que el porcentaje se incluye en la respuesta
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)

	// Verificar que se llamaron todos los mocks
	mockExpenseRepo.AssertExpectations(t)
}
