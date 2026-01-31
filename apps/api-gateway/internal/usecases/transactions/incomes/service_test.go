package incomes

import (
	"context"
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockIncomeRepository es un mock del repositorio de ingresos
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

func TestCreateIncome(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewIncomeService(mockRepo)

	request := &CreateIncomeRequest{
		UserID:      "user1",
		Amount:      1000.50,
		Description: "Salary",
		CategoryID:  "work",
	}

	mockRepo.On("Create", mock.Anything).Return(nil)

	response, err := service.CreateIncome(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.UserID, response.UserID)
	assert.Equal(t, request.Amount, response.Amount)
	assert.Equal(t, request.Description, response.Description)
	assert.Equal(t, request.CategoryID, response.CategoryID)
	mockRepo.AssertExpectations(t)
}

func TestGetIncome(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewIncomeService(mockRepo)

	expectedIncome := domain.NewIncomeBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(1000.50).
		SetDescription("Salary").
		SetCategoryID("work").
		Build()

	mockRepo.On("Get", "user1", "1").Return(expectedIncome, nil)

	response, err := service.GetIncome(context.Background(), "user1", "1")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedIncome.ID, response.ID)
	assert.Equal(t, expectedIncome.UserID, response.UserID)
	assert.Equal(t, expectedIncome.Amount, response.Amount)
	assert.Equal(t, expectedIncome.Description, response.Description)
	assert.Equal(t, expectedIncome.GetCategoryID(), response.CategoryID)
	mockRepo.AssertExpectations(t)
}

func TestListIncomes(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewIncomeService(mockRepo)

	income1 := domain.NewIncomeBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(1000.0).
		SetDescription("Salary").
		SetCategoryID("work").
		Build()

	income2 := domain.NewIncomeBuilder().
		SetID("2").
		SetUserID("user1").
		SetAmount(2000.0).
		SetDescription("Bonus").
		SetCategoryID("work").
		Build()

	expectedIncomes := []*domain.Income{income1, income2}

	mockRepo.On("List", "user1").Return(expectedIncomes, nil)

	result, err := service.ListIncomes(context.Background(), "user1")
	assert.NoError(t, err)
	assert.Len(t, result.Incomes, 2)
	assert.Equal(t, income1.ID, result.Incomes[0].ID)
	assert.Equal(t, income2.ID, result.Incomes[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateIncome(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewIncomeService(mockRepo)

	existingIncome := domain.NewIncomeBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(1000.0).
		SetDescription("Salary").
		SetCategoryID("work").
		Build()

	request := &UpdateIncomeRequest{
		Amount:      1500.0,
		Description: "Updated Salary",
		CategoryID:  "work",
	}

	mockRepo.On("Get", "user1", "1").Return(existingIncome, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	response, err := service.UpdateIncome(context.Background(), "user1", "1", request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, request.Amount, response.Amount)
	assert.Equal(t, request.Description, response.Description)
	assert.Equal(t, request.CategoryID, response.CategoryID)
	mockRepo.AssertExpectations(t)
}

func TestDeleteIncome(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewIncomeService(mockRepo)

	existingIncome := domain.NewIncomeBuilder().
		SetID("1").
		SetUserID("user1").
		SetAmount(1000.0).
		SetDescription("Salary").
		SetCategoryID("work").
		Build()

	mockRepo.On("Get", "user1", "1").Return(existingIncome, nil)
	mockRepo.On("Delete", "user1", "1").Return(nil)

	err := service.DeleteIncome(context.Background(), "user1", "1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateIncome_Validation(t *testing.T) {
	mockRepo := new(MockIncomeRepository)
	service := NewIncomeService(mockRepo)

	tests := []struct {
		name        string
		request     *CreateIncomeRequest
		expectedErr string
	}{
		{
			name: "Empty description",
			request: &CreateIncomeRequest{
				UserID:      "user1",
				Amount:      1000.50,
				Description: "",
				CategoryID:  "work",
			},
			expectedErr: "La descripción del ingreso es requerida",
		},
		{
			name: "Zero amount",
			request: &CreateIncomeRequest{
				UserID:      "user1",
				Amount:      0,
				Description: "Salary",
				CategoryID:  "work",
			},
			expectedErr: "El monto del ingreso debe ser mayor a 0",
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
