package update

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository es un mock del repositorio de categorías
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Get(userID string, name string) (*domain.Category, error) {
	args := m.Called(userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetByName(userID string, name string) (*domain.Category, error) {
	args := m.Called(userID, name)
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

func (m *MockCategoryRepository) UpdateCategoryReferences(categoryID string) error {
	args := m.Called(categoryID)
	return args.Error(0)
}

func TestUpdateCategory(t *testing.T) {
	tests := []struct {
		name           string
		input          UpdateCategoryRequest
		mockSetup      func(*MockCategoryRepository)
		expectedOutput *UpdateCategoryResponse
		expectedError  error
	}{
		{
			name: "actualizar categoría exitosamente",
			input: UpdateCategoryRequest{
				ID:      "cat_12345678",
				UserID:  "user_id",
				NewName: "Updated Category",
			},
			mockSetup: func(m *MockCategoryRepository) {
				existingCategory := &domain.Category{
					ID:     "cat_12345678",
					Name:   "Test Category",
					UserID: "user_id",
				}
				m.On("Get", "user_id", "cat_12345678").Return(existingCategory, nil)
				m.On("GetByName", "user_id", "Updated Category").Return(nil, errors.NewResourceNotFound("Category not found"))
				m.On("Update", mock.MatchedBy(func(c *domain.Category) bool {
					return c.ID == "cat_12345678" && c.Name == "Updated Category" && c.UserID == "user_id"
				})).Return(nil)
				m.On("UpdateCategoryReferences", "cat_12345678").Return(nil)
			},
			expectedOutput: &UpdateCategoryResponse{
				ID:     "cat_12345678",
				Name:   "Updated Category",
				UserID: "user_id",
			},
			expectedError: nil,
		},
		{
			name: "error al actualizar categoría que no existe",
			input: UpdateCategoryRequest{
				ID:      "cat_12345678",
				UserID:  "user_id",
				NewName: "Updated Category",
			},
			mockSetup: func(m *MockCategoryRepository) {
				m.On("Get", "user_id", "cat_12345678").Return(nil, errors.NewResourceNotFound("Category not found"))
			},
			expectedOutput: nil,
			expectedError:  errors.NewResourceNotFound("Category not found"),
		},
		{
			name: "error al actualizar categoría con nombre vacío",
			input: UpdateCategoryRequest{
				ID:      "cat_12345678",
				UserID:  "user_id",
				NewName: "",
			},
			mockSetup: func(m *MockCategoryRepository) {
				existingCategory := &domain.Category{
					ID:     "cat_12345678",
					Name:   "Test Category",
					UserID: "user_id",
				}
				m.On("Get", "user_id", "cat_12345678").Return(existingCategory, nil)
				m.On("GetByName", "user_id", "").Return(nil, errors.NewBadRequest("Category name cannot be empty"))
			},
			expectedOutput: nil,
			expectedError:  errors.NewBadRequest("Category name cannot be empty"),
		},
		{
			name: "error al actualizar categoría con userID vacío",
			input: UpdateCategoryRequest{
				UserID:  "",
				ID:      "cat_12345678",
				NewName: "New Name",
			},
			mockSetup: func(m *MockCategoryRepository) {
				m.On("Get", "", "cat_12345678").Return(nil, errors.NewBadRequest("User ID cannot be empty"))
			},
			expectedOutput: nil,
			expectedError:  errors.NewBadRequest("User ID cannot be empty"),
		},
		{
			name: "error al actualizar categoría con nuevo nombre vacío",
			input: UpdateCategoryRequest{
				UserID:  "user_id",
				ID:      "cat_12345678",
				NewName: "",
			},
			mockSetup: func(m *MockCategoryRepository) {
				existingCategory := &domain.Category{
					ID:     "cat_12345678",
					Name:   "Test Category",
					UserID: "user_id",
				}
				m.On("Get", "user_id", "cat_12345678").Return(existingCategory, nil)
				m.On("GetByName", "user_id", "").Return(nil, errors.NewBadRequest("New name cannot be empty"))
			},
			expectedOutput: nil,
			expectedError:  errors.NewBadRequest("New name cannot be empty"),
		},
		{
			name: "error al actualizar categoría con error interno",
			input: UpdateCategoryRequest{
				UserID:  "user_id",
				ID:      "cat_12345678",
				NewName: "Updated Name",
			},
			mockSetup: func(m *MockCategoryRepository) {
				existingCategory := &domain.Category{
					ID:     "cat_12345678",
					Name:   "Test Category",
					UserID: "user_id",
				}
				m.On("Get", "user_id", "cat_12345678").Return(existingCategory, nil)
				m.On("GetByName", "user_id", "Updated Name").Return(nil, errors.NewResourceNotFound("Category not found"))
				m.On("Update", mock.MatchedBy(func(c *domain.Category) bool {
					return c.ID == "cat_12345678" && c.Name == "Updated Name" && c.UserID == "user_id"
				})).Return(errors.NewInternalServerError("Internal server error"))
			},
			expectedOutput: nil,
			expectedError:  errors.NewInternalServerError("Internal server error"),
		},
		{
			name: "error al actualizar categoría que no existe",
			input: UpdateCategoryRequest{
				ID:      "cat_12345678",
				UserID:  "user_id",
				NewName: "Updated Category",
			},
			mockSetup: func(m *MockCategoryRepository) {
				m.On("Get", "user_id", "cat_12345678").Return(nil, errors.NewResourceNotFound("Category not found"))
			},
			expectedOutput: nil,
			expectedError:  errors.NewResourceNotFound("Category not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCategoryRepository)
			service := NewUpdateCategory(mockRepo)
			tt.mockSetup(mockRepo)

			response, err := service.Execute(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedError, err)
				assert.Nil(t, response)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, tt.input.UserID, response.UserID)
			assert.Equal(t, tt.input.NewName, response.Name)
			mockRepo.AssertExpectations(t)
		})
	}
}
