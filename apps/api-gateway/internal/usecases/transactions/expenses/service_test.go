package expenses

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

func (m *MockExpenseRepository) Get(userID string, expenseID string) (*domain.Expense, error) {
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

func (m *MockExpenseRepository) ListUnpaid(userID string) ([]*domain.Expense, error) {
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

func (m *MockExpenseRepository) Delete(userID string, expenseID string) error {
	args := m.Called(userID, expenseID)
	return args.Error(0)
}

func TestCreateExpense(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := NewExpenseService(mockRepo)

	request := &CreateExpenseRequest{
		UserID:      "user1",
		Amount:      100.0,
		Description: "Test expense",
		CategoryID:  "food",
		Paid:        false,
		DueDate:     time.Now().Format("2006-01-02"),
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Expense")).Return(nil)

	result, err := service.CreateExpense(context.Background(), request)
	assert.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, request.UserID, result.UserID)
	assert.Equal(t, request.Amount, result.Amount)
	assert.Equal(t, request.Description, result.Description)
	assert.Equal(t, request.CategoryID, result.CategoryID)
	assert.Equal(t, request.Paid, result.Paid)
	assert.Equal(t, request.DueDate, result.DueDate)
	mockRepo.AssertExpectations(t)
}

func TestGetExpense(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := NewExpenseService(mockRepo)

	expectedExpense := domain.NewExpenseBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(100.0).
		SetDescription("Test expense").
		SetCategoryID("food").
		SetPaid(false).
		SetDueDate(time.Now()).
		Build()

	mockRepo.On("Get", "user1", "1").Return(expectedExpense, nil)

	result, err := service.GetExpense(context.Background(), "user1", "1")
	assert.NoError(t, err)
	assert.Equal(t, expectedExpense.ID, result.ID)
	assert.Equal(t, expectedExpense.UserID, result.UserID)
	assert.Equal(t, expectedExpense.Amount, result.Amount)
	assert.Equal(t, expectedExpense.Description, result.Description)
	assert.Equal(t, expectedExpense.GetCategoryID(), result.CategoryID)
	assert.Equal(t, expectedExpense.Paid, result.Paid)
	mockRepo.AssertExpectations(t)
}

func TestListExpenses(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := NewExpenseService(mockRepo)

	expense1 := domain.NewExpenseBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(100.0).
		SetDescription("Expense 1").
		SetCategoryID("food").
		SetPaid(false).
		SetDueDate(time.Now()).
		Build()

	expense2 := domain.NewExpenseBuilder().
		SetID("2").
		SetUserID("user1").
		SetAmount(200.0).
		SetDescription("Expense 2").
		SetCategoryID("transport").
		SetPaid(true).
		SetDueDate(time.Now()).
		Build()

	expectedExpenses := []*domain.Expense{expense1, expense2}

	mockRepo.On("List", "user1").Return(expectedExpenses, nil)

	result, err := service.ListExpenses(context.Background(), "user1")
	assert.NoError(t, err)
	assert.Len(t, result.Expenses, 2)
	assert.Equal(t, expense1.ID, result.Expenses[0].ID)
	assert.Equal(t, expense2.ID, result.Expenses[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateExpense(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := NewExpenseService(mockRepo)

	existingExpense := domain.NewExpenseBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(100.0).
		SetDescription("Test expense").
		SetCategoryID("food").
		SetPaid(false).
		SetDueDate(time.Now()).
		Build()

	request := &UpdateExpenseRequest{
		Amount:      150.0,
		Description: "Updated expense",
		CategoryID:  "entertainment",
		Paid:        true,
		DueDate:     time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
	}

	updatedExpense := domain.NewExpenseBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(request.Amount).
		SetDescription(request.Description).
		SetCategoryID(request.CategoryID).
		SetPaid(request.Paid).
		SetDueDate(time.Now()).
		Build()

	mockRepo.On("Get", "user1", "1").Return(existingExpense, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Expense")).Return(nil)

	result, err := service.UpdateExpense(context.Background(), "user1", "1", request)
	assert.NoError(t, err)
	assert.Equal(t, updatedExpense.ID, result.ID)
	assert.Equal(t, request.Amount, result.Amount)
	assert.Equal(t, request.Description, result.Description)
	assert.Equal(t, request.CategoryID, result.CategoryID)
	assert.Equal(t, request.Paid, result.Paid)
	assert.Equal(t, request.DueDate, result.DueDate)
	mockRepo.AssertExpectations(t)
}

func TestDeleteExpense(t *testing.T) {
	mockRepo := new(MockExpenseRepository)
	service := NewExpenseService(mockRepo)

	expense := domain.NewExpenseBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(100.0).
		SetDescription("Test expense").
		SetCategoryID("food").
		SetPaid(false).
		SetDueDate(time.Now()).
		Build()

	mockRepo.On("Get", "user1", "1").Return(expense, nil)
	mockRepo.On("Delete", "user1", "1").Return(nil)

	err := service.DeleteExpense(context.Background(), "user1", "1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
