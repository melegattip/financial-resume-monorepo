package create

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

func TestCreateIncome(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewService(mockRepo, mockCategoryRepo)

	request := &incomesDomain.CreateIncomeRequest{
		UserID:      "user1",
		Amount:      1000.50,
		Description: "Salary",
		CategoryID:  "work",
		Source:      "Company XYZ",
	}

	mockCategoryRepo.On("Get", request.UserID, request.CategoryID).Return(&domain.Category{}, nil)
	mockRepo.On("Create", mock.Anything).Return(nil)

	response, err := service.CreateIncome(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.UserID, response.UserID)
	assert.Equal(t, request.Amount, response.Amount)
	assert.Equal(t, request.Description, response.Description)
	assert.Equal(t, request.CategoryID, response.CategoryID)
	mockRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestCreateIncome_ValidationErrors(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewService(mockRepo, mockCategoryRepo)

	tests := []struct {
		name        string
		request     *incomesDomain.CreateIncomeRequest
		expectedErr string
	}{
		{
			name: "Empty description",
			request: &incomesDomain.CreateIncomeRequest{
				UserID:      "user1",
				Amount:      1000.50,
				Description: "",
				CategoryID:  "work",
				Source:      "Company XYZ",
			},
			expectedErr: "La descripción del ingreso es requerida",
		},
		{
			name: "Zero amount",
			request: &incomesDomain.CreateIncomeRequest{
				UserID:      "user1",
				Amount:      0,
				Description: "Salary",
				CategoryID:  "work",
				Source:      "Company XYZ",
			},
			expectedErr: "El monto del ingreso debe ser mayor a 0",
		},
		{
			name: "Empty source",
			request: &incomesDomain.CreateIncomeRequest{
				UserID:      "user1",
				Amount:      1000.50,
				Description: "Salary",
				CategoryID:  "work",
				Source:      "",
			},
			expectedErr: "La fuente del ingreso es requerida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := service.CreateIncome(context.Background(), tt.request)
			assert.Nil(t, response)
			assert.Error(t, err)
			assert.Equal(t, tt.expectedErr, err.Error())
		})
	}
}

func TestCreateIncome_InvalidCategory(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewService(mockRepo, mockCategoryRepo)

	request := &incomesDomain.CreateIncomeRequest{
		UserID:      "user1",
		Amount:      1000.50,
		Description: "Salary",
		CategoryID:  "invalid_category",
		Source:      "Company XYZ",
	}

	mockCategoryRepo.On("Get", request.UserID, request.CategoryID).Return(nil, errors.NewResourceNotFound("Category not found"))

	response, err := service.CreateIncome(context.Background(), request)
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Equal(t, "La categoría especificada no existe", err.Error())
	mockCategoryRepo.AssertExpectations(t)
}
