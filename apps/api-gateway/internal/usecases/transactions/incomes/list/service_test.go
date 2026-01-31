package list

import (
	"context"
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIncomeRepository struct {
	mock.Mock
}

func (m *MockIncomeRepository) Create(income *domain.Income) error {
	args := m.Called(income)
	return args.Error(0)
}

func (m *MockIncomeRepository) Get(userID string, id string) (*domain.Income, error) {
	args := m.Called(userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Income), args.Error(1)
}

func (m *MockIncomeRepository) List(userID string) ([]*domain.Income, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Income), args.Error(1)
}

func (m *MockIncomeRepository) Update(income *domain.Income) error {
	args := m.Called(income)
	return args.Error(0)
}

func (m *MockIncomeRepository) Delete(userID string, id string) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

func TestListIncomes(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewService(mockRepo)

	incomes := []*domain.Income{
		domain.NewIncomeBuilder().
			SetID("1").
			SetUserID("user1").
			SetAmount(1000.50).
			SetDescription("Salary").
			SetCategoryID("work").
			Build(),
		domain.NewIncomeBuilder().
			SetID("2").
			SetUserID("user1").
			SetAmount(2000.00).
			SetDescription("Bonus").
			SetCategoryID("work").
			Build(),
	}

	mockRepo.On("List", "user1").Return(incomes, nil)

	response, err := service.ListIncomes(context.Background(), "user1")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Incomes, 2)
	assert.Equal(t, incomes[0].ID, response.Incomes[0].ID)
	assert.Equal(t, incomes[0].UserID, response.Incomes[0].UserID)
	assert.Equal(t, incomes[0].Amount, response.Incomes[0].Amount)
	assert.Equal(t, incomes[0].Description, response.Incomes[0].Description)
	assert.Equal(t, incomes[0].GetCategoryID(), response.Incomes[0].CategoryID)
	assert.Equal(t, incomes[1].ID, response.Incomes[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestListIncomes_EmptyUserID(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewService(mockRepo)

	response, err := service.ListIncomes(context.Background(), "")
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Equal(t, "El ID del usuario es requerido", err.Error())
}

func TestListIncomes_EmptyList(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewService(mockRepo)

	mockRepo.On("List", "user1").Return([]*domain.Income{}, nil)

	response, err := service.ListIncomes(context.Background(), "user1")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Empty(t, response.Incomes)
	mockRepo.AssertExpectations(t)
}
