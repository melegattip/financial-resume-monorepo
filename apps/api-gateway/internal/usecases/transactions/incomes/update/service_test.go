package update

import (
	"context"
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	incomesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
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

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Get(userID string, id string) (*domain.Category, error) {
	args := m.Called(userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) List(userID string) ([]*domain.Category, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(userID string, id string) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByName(userID string, name string) (*domain.Category, error) {
	args := m.Called(userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) UpdateCategoryReferences(categoryID string) error {
	args := m.Called(categoryID)
	return args.Error(0)
}

func TestUpdateIncome(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewService(mockRepo, mockCategoryRepo)

	existingIncome := domain.NewIncomeBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(1000.50).
		SetDescription("Salary").
		SetCategoryID("work").
		Build()

	request := &incomesDomain.UpdateIncomeRequest{
		Amount:      1500.00,
		Description: "Updated Salary",
		CategoryID:  "work",
		Source:      "New Company",
	}

	mockRepo.On("Get", "user1", "1").Return(existingIncome, nil)
	mockCategoryRepo.On("Get", "user1", request.CategoryID).Return(&domain.Category{}, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	response, err := service.UpdateIncome(context.Background(), "user1", "1", request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, existingIncome.ID, response.ID)
	assert.Equal(t, request.Amount, response.Amount)
	assert.Equal(t, request.Description, response.Description)
	assert.Equal(t, request.CategoryID, response.CategoryID)
	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestUpdateIncome_ValidationErrors(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewService(mockRepo, mockCategoryRepo)

	tests := []struct {
		name        string
		userID      string
		id          string
		request     *incomesDomain.UpdateIncomeRequest
		expectedErr string
	}{
		{
			name:   "Empty user ID",
			userID: "",
			id:     "1",
			request: &incomesDomain.UpdateIncomeRequest{
				Amount:      1500.00,
				Description: "Updated Salary",
			},
			expectedErr: "El ID del usuario es requerido",
		},
		{
			name:   "Empty income ID",
			userID: "user1",
			id:     "",
			request: &incomesDomain.UpdateIncomeRequest{
				Amount:      1500.00,
				Description: "Updated Salary",
			},
			expectedErr: "El ID del ingreso es requerido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.UpdateIncome(context.Background(), tt.userID, tt.id, tt.request)
			assert.Nil(t, response)
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}

func TestUpdateIncome_NotFound(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewService(mockRepo, mockCategoryRepo)

	request := &incomesDomain.UpdateIncomeRequest{
		Amount:      1500.00,
		Description: "Updated Salary",
	}

	mockRepo.On("Get", "user1", "1").Return(nil, errors.NewResourceNotFound("Income not found"))

	response, err := service.UpdateIncome(context.Background(), "user1", "1", request)
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Equal(t, "Income not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateIncome_InvalidCategory(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewService(mockRepo, mockCategoryRepo)

	existingIncome := domain.NewIncomeBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(1000.50).
		SetDescription("Salary").
		SetCategoryID("work").
		Build()

	request := &incomesDomain.UpdateIncomeRequest{
		CategoryID: "invalid_category",
	}

	mockRepo.On("Get", "user1", "1").Return(existingIncome, nil)
	mockCategoryRepo.On("Get", "user1", request.CategoryID).Return(nil, errors.NewResourceNotFound("Category not found"))

	response, err := service.UpdateIncome(context.Background(), "user1", "1", request)
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Equal(t, "La categoría especificada no existe", err.Error())
	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}
