package get

import (
	"context"
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
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

func TestGetIncome(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewService(mockRepo)

	income := domain.NewIncomeBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(1000.50).
		SetDescription("Salary").
		SetCategoryID("work").
		Build()

	mockRepo.On("Get", "user1", "1").Return(income, nil)

	response, err := service.GetIncome(context.Background(), "user1", "1")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, income.ID, response.ID)
	assert.Equal(t, income.UserID, response.UserID)
	assert.Equal(t, income.Amount, response.Amount)
	assert.Equal(t, income.Description, response.Description)
	assert.Equal(t, income.GetCategoryID(), response.CategoryID)
	mockRepo.AssertExpectations(t)
}

func TestGetIncome_ValidationErrors(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewService(mockRepo)

	tests := []struct {
		name        string
		userID      string
		id          string
		expectedErr string
	}{
		{
			name:        "Empty user ID",
			userID:      "",
			id:          "1",
			expectedErr: "El ID del usuario es requerido",
		},
		{
			name:        "Empty income ID",
			userID:      "user1",
			id:          "",
			expectedErr: "El ID del ingreso es requerido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.GetIncome(context.Background(), tt.userID, tt.id)
			assert.Nil(t, response)
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}

func TestGetIncome_NotFound(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewService(mockRepo)

	mockRepo.On("Get", "user1", "1").Return(nil, errors.NewResourceNotFound("Income not found"))

	response, err := service.GetIncome(context.Background(), "user1", "1")
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Equal(t, "Income not found", err.Error())
	mockRepo.AssertExpectations(t)
}
