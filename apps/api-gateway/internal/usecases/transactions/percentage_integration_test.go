package transactions

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExpenseRepositoryIntegration es un mock más completo para tests de integración
type MockExpenseRepositoryIntegration struct {
	mock.Mock
	expenses map[string]*domain.Expense // Simula almacenamiento en memoria
}

func NewMockExpenseRepositoryIntegration() *MockExpenseRepositoryIntegration {
	return &MockExpenseRepositoryIntegration{
		expenses: make(map[string]*domain.Expense),
	}
}

func (m *MockExpenseRepositoryIntegration) Create(expense *domain.Expense) error {
	args := m.Called(expense)
	if args.Error(0) == nil {
		m.expenses[expense.ID] = expense
	}
	return args.Error(0)
}

func (m *MockExpenseRepositoryIntegration) Get(userID, expenseID string) (*domain.Expense, error) {
	args := m.Called(userID, expenseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepositoryIntegration) List(userID string) ([]*domain.Expense, error) {
	args := m.Called(userID)

	// Retornar gastos del almacenamiento simulado
	var userExpenses []*domain.Expense
	for _, expense := range m.expenses {
		if expense.UserID == userID {
			userExpenses = append(userExpenses, expense)
		}
	}
	return userExpenses, args.Error(0)
}

func (m *MockExpenseRepositoryIntegration) Update(expense *domain.Expense) error {
	args := m.Called(expense)
	if args.Error(0) == nil {
		m.expenses[expense.ID] = expense
	}
	return args.Error(0)
}

func (m *MockExpenseRepositoryIntegration) Delete(userID, expenseID string) error {
	args := m.Called(userID, expenseID)
	if args.Error(0) == nil {
		delete(m.expenses, expenseID)
	}
	return args.Error(0)
}

func (m *MockExpenseRepositoryIntegration) ListUnpaid(userID string) ([]*domain.Expense, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Expense), args.Error(1)
}

// MockIncomeRepositoryIntegration simula el repositorio de ingresos
type MockIncomeRepositoryIntegration struct {
	mock.Mock
	incomes map[string]*domain.Income
}

func NewMockIncomeRepositoryIntegration() *MockIncomeRepositoryIntegration {
	return &MockIncomeRepositoryIntegration{
		incomes: make(map[string]*domain.Income),
	}
}

func (m *MockIncomeRepositoryIntegration) Create(income *domain.Income) error {
	args := m.Called(income)
	if args.Error(0) == nil {
		m.incomes[income.ID] = income
	}
	return args.Error(0)
}

func (m *MockIncomeRepositoryIntegration) Get(userID, incomeID string) (*domain.Income, error) {
	args := m.Called(userID, incomeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Income), args.Error(1)
}

func (m *MockIncomeRepositoryIntegration) List(userID string) ([]*domain.Income, error) {
	args := m.Called(userID)

	// Retornar ingresos del almacenamiento simulado
	var userIncomes []*domain.Income
	for _, income := range m.incomes {
		if income.UserID == userID {
			userIncomes = append(userIncomes, income)
		}
	}
	return userIncomes, args.Error(0)
}

func (m *MockIncomeRepositoryIntegration) Update(income *domain.Income) error {
	args := m.Called(income)
	if args.Error(0) == nil {
		m.incomes[income.ID] = income
	}
	return args.Error(0)
}

func (m *MockIncomeRepositoryIntegration) Delete(userID, incomeID string) error {
	args := m.Called(userID, incomeID)
	if args.Error(0) == nil {
		delete(m.incomes, incomeID)
	}
	return args.Error(0)
}

func (m *MockIncomeRepositoryIntegration) ListUnreceived(userID string) ([]*domain.Income, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Income), args.Error(1)
}

func TestPercentageUpdateIntegration_WhenIncomeChanges(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := "user123"

	// Crear repositorios mock
	expenseRepo := NewMockExpenseRepositoryIntegration()
	incomeRepo := NewMockIncomeRepositoryIntegration()

	// Crear servicio real
	expenseService := NewExpenseService(expenseRepo, incomeRepo)
	percentageObserver := NewPercentageObserver(expenseService)

	// Configurar datos iniciales
	// 1. Crear un ingreso inicial de $1000
	initialIncome := &domain.Income{
		ID:        "income1",
		UserID:    userID,
		Amount:    1000.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	incomeRepo.incomes[initialIncome.ID] = initialIncome

	// 2. Crear dos gastos: $100 y $200
	expense1 := &domain.Expense{
		ID:          "expense1",
		UserID:      userID,
		Amount:      100.0,
		Description: "Groceries",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Percentage:  10.0, // 100/1000 = 10%
	}
	expense2 := &domain.Expense{
		ID:          "expense2",
		UserID:      userID,
		Amount:      200.0,
		Description: "Gas",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Percentage:  20.0, // 200/1000 = 20%
	}
	expenseRepo.expenses[expense1.ID] = expense1
	expenseRepo.expenses[expense2.ID] = expense2

	// Configurar mocks para permitir las operaciones
	incomeRepo.On("List", userID).Return(nil)
	expenseRepo.On("List", userID).Return(nil)
	expenseRepo.On("Update", mock.AnythingOfType("*domain.Expense")).Return(nil)

	// Act: Agregar un nuevo ingreso de $500 (total ahora será $1500)
	newIncome := &domain.Income{
		ID:        "income2",
		UserID:    userID,
		Amount:    500.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	incomeRepo.incomes[newIncome.ID] = newIncome

	// Simular que se notifica al observer sobre el nuevo ingreso
	err := percentageObserver.OnTransactionCreated(ctx, newIncome)

	// Assert
	assert.NoError(t, err)

	// Verificar que los porcentajes se actualizaron correctamente
	// Con total de ingresos de $1500:
	// - expense1 ($100) debería ser 6.67% (100/1500)
	// - expense2 ($200) debería ser 13.33% (200/1500)

	updatedExpense1 := expenseRepo.expenses["expense1"]
	updatedExpense2 := expenseRepo.expenses["expense2"]

	assert.InDelta(t, 6.67, updatedExpense1.Percentage, 0.01, "Expense1 percentage should be updated to ~6.67%")
	assert.InDelta(t, 13.33, updatedExpense2.Percentage, 0.01, "Expense2 percentage should be updated to ~13.33%")

	// Verificar que se llamaron los métodos esperados
	expenseRepo.AssertExpectations(t)
	incomeRepo.AssertExpectations(t)
}

func TestPercentageUpdateIntegration_WhenIncomeDeleted(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := "user123"

	// Crear repositorios mock
	expenseRepo := NewMockExpenseRepositoryIntegration()
	incomeRepo := NewMockIncomeRepositoryIntegration()

	// Crear servicio real
	expenseService := NewExpenseService(expenseRepo, incomeRepo)
	percentageObserver := NewPercentageObserver(expenseService)

	// Configurar datos iniciales
	// 1. Crear dos ingresos: $1000 y $500 (total $1500)
	income1 := &domain.Income{
		ID:        "income1",
		UserID:    userID,
		Amount:    1000.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	income2 := &domain.Income{
		ID:        "income2",
		UserID:    userID,
		Amount:    500.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	incomeRepo.incomes[income1.ID] = income1
	incomeRepo.incomes[income2.ID] = income2

	// 2. Crear un gasto de $150
	expense1 := &domain.Expense{
		ID:          "expense1",
		UserID:      userID,
		Amount:      150.0,
		Description: "Shopping",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Percentage:  10.0, // 150/1500 = 10%
	}
	expenseRepo.expenses[expense1.ID] = expense1

	// Configurar mocks
	incomeRepo.On("List", userID).Return(nil)
	expenseRepo.On("List", userID).Return(nil)
	expenseRepo.On("Update", mock.AnythingOfType("*domain.Expense")).Return(nil)
	expenseRepo.On("Get", "", "income2").Return(nil, errors.New("not found")) // Gasto no encontrado
	incomeRepo.On("Get", "", "income2").Return(income2, nil)                  // Ingreso encontrado

	// Act: Eliminar income2 ($500), quedando solo $1000 de ingresos
	delete(incomeRepo.incomes, "income2")

	// Simular que se notifica al observer sobre la eliminación
	err := percentageObserver.OnTransactionDeleted(ctx, "income2")

	// Assert
	assert.NoError(t, err)

	// Verificar que el porcentaje se actualizó correctamente
	// Con total de ingresos de $1000:
	// - expense1 ($150) debería ser 15% (150/1000)

	updatedExpense1 := expenseRepo.expenses["expense1"]
	assert.Equal(t, 15.0, updatedExpense1.Percentage, "Expense1 percentage should be updated to 15%")

	// Verificar que se llamaron los métodos esperados
	expenseRepo.AssertExpectations(t)
	incomeRepo.AssertExpectations(t)
}

func TestPercentageUpdateIntegration_WhenNoIncomes(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := "user123"

	// Crear repositorios mock
	expenseRepo := NewMockExpenseRepositoryIntegration()
	incomeRepo := NewMockIncomeRepositoryIntegration()

	// Crear servicio real
	expenseService := NewExpenseService(expenseRepo, incomeRepo)
	percentageObserver := NewPercentageObserver(expenseService)

	// Configurar datos iniciales
	// Solo un gasto, sin ingresos
	expense1 := &domain.Expense{
		ID:          "expense1",
		UserID:      userID,
		Amount:      100.0,
		Description: "Emergency expense",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Percentage:  50.0, // Valor anterior cualquiera
	}
	expenseRepo.expenses[expense1.ID] = expense1

	// Configurar mocks
	incomeRepo.On("List", userID).Return(nil)
	expenseRepo.On("List", userID).Return(nil)
	expenseRepo.On("Update", mock.AnythingOfType("*domain.Expense")).Return(nil)

	// Act: Actualizar porcentajes sin ingresos
	err := percentageObserver.UpdatePercentages(ctx, userID)

	// Assert
	assert.NoError(t, err)

	// Verificar que el porcentaje se puso en 0% (sin ingresos)
	updatedExpense1 := expenseRepo.expenses["expense1"]
	assert.Equal(t, 0.0, updatedExpense1.Percentage, "Expense percentage should be 0% when no incomes exist")

	// Verificar que se llamaron los métodos esperados
	expenseRepo.AssertExpectations(t)
	incomeRepo.AssertExpectations(t)
}
