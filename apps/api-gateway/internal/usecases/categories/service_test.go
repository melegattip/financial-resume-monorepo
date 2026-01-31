package categories

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type GetCategoryRequest struct {
	ID     string
	Name   string
	UserID string
}

// MockCategoryRepository es un mock del repositorio de categorías
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

func TestService(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		mockSetup      func(*MockCategoryRepository)
		expectedOutput interface{}
		expectedError  error
	}{
		{
			name: "successfully create category",
			input: CreateCategoryRequest{
				Name:   "Test Category",
				UserID: "user_id",
			},
			mockSetup: func(m *MockCategoryRepository) {
				m.On("GetByName", "user_id", "Test Category").Return(nil, errors.NewResourceNotFound("category not found"))
				m.On("Create", mock.AnythingOfType("*domain.Category")).Return(nil)
			},
			expectedOutput: &domain.Category{
				Name:   "Test Category",
				UserID: "user_id",
			},
			expectedError: nil,
		},
		{
			name: "error creating category that already exists",
			input: CreateCategoryRequest{
				Name:   "Existing Category",
				UserID: "user_id",
			},
			mockSetup: func(m *MockCategoryRepository) {
				existingCategory := &domain.Category{
					ID:     "cat_12345678",
					Name:   "Existing Category",
					UserID: "user_id",
				}
				m.On("GetByName", "user_id", "Existing Category").Return(existingCategory, nil)
			},
			expectedOutput: nil,
			expectedError:  errors.NewResourceAlreadyExists("Error creating category"),
		},
		{
			name:  "successfully list categories",
			input: "user_id",
			mockSetup: func(m *MockCategoryRepository) {
				categories := []*domain.Category{
					{
						ID:     "cat_12345678",
						Name:   "Category 1",
						UserID: "user_id",
					},
					{
						ID:     "cat_87654321",
						Name:   "Category 2",
						UserID: "user_id",
					},
				}
				m.On("List", "user_id").Return(categories, nil)
			},
			expectedOutput: []*domain.Category{
				{
					Name:   "Category 1",
					UserID: "user_id",
				},
				{
					Name:   "Category 2",
					UserID: "user_id",
				},
			},
			expectedError: nil,
		},
		{
			name:           "error listing categories without user_id",
			input:          "",
			mockSetup:      func(m *MockCategoryRepository) {},
			expectedOutput: nil,
			expectedError:  errors.NewBadRequest("User ID cannot be empty"),
		},
		{
			name: "successfully get category",
			input: GetCategoryRequest{
				Name:   "Test_Category",
				UserID: "user_id",
			},
			mockSetup: func(m *MockCategoryRepository) {
				category := &domain.Category{
					ID:     "cat_12345678",
					Name:   "Test_Category",
					UserID: "user_id",
				}
				m.On("GetByName", "user_id", "Test_Category").Return(category, nil)
			},
			expectedOutput: &domain.Category{
				Name:   "Test_Category",
				UserID: "user_id",
			},
			expectedError: nil,
		},
		{
			name: "error getting non-existent category",
			input: GetCategoryRequest{
				Name:   "non_existent",
				UserID: "user_id",
			},
			mockSetup: func(m *MockCategoryRepository) {
				m.On("GetByName", "user_id", "non_existent").Return(nil, errors.NewResourceNotFound("Category not found"))
			},
			expectedOutput: nil,
			expectedError:  errors.NewResourceNotFound("Category not found"),
		},
		{
			name: "error getting category without name",
			input: GetCategoryRequest{
				Name:   "",
				UserID: "user_id",
			},
			mockSetup: func(m *MockCategoryRepository) {
				m.On("GetByName", "user_id", "").Return(nil, errors.NewBadRequest("Name is required"))
			},
			expectedOutput: nil,
			expectedError:  errors.NewBadRequest("Name is required"),
		},
		{
			name: "successfully update category",
			input: UpdateCategoryRequest{
				Name:    "Original Category",
				UserID:  "user_id",
				NewName: "Updated Category",
			},
			mockSetup: func(m *MockCategoryRepository) {
				category := &domain.Category{
					ID:     "cat_12345678",
					Name:   "Original Category",
					UserID: "user_id",
				}
				m.On("Get", "user_id", "Original Category").Return(category, nil)
				m.On("GetByName", "user_id", "Updated Category").Return(nil, errors.NewResourceNotFound("Category not found"))
				m.On("Update", mock.MatchedBy(func(c *domain.Category) bool {
					return c.ID == "cat_12345678" && c.Name == "Updated Category" && c.UserID == "user_id"
				})).Return(nil)
				m.On("UpdateCategoryReferences", "cat_12345678").Return(nil)
			},
			expectedOutput: &update.UpdateCategoryResponse{
				ID:     "cat_12345678",
				Name:   "Updated Category",
				UserID: "user_id",
			},
			expectedError: nil,
		},
		{
			name: "successfully delete category",
			input: DeleteCategoryRequest{
				ID:     "cat_12345678",
				UserID: "user_id",
			},
			mockSetup: func(m *MockCategoryRepository) {
				m.On("Delete", "user_id", "cat_12345678").Return(nil)
			},
			expectedOutput: nil,
			expectedError:  nil,
		},
		{
			name: "error deleting non-existent category",
			input: DeleteCategoryRequest{
				ID:     "non_existent",
				UserID: "user_id",
			},
			mockSetup: func(m *MockCategoryRepository) {
				m.On("Delete", "user_id", "non_existent").Return(errors.NewResourceNotFound("Category not found"))
			},
			expectedOutput: nil,
			expectedError:  errors.NewResourceNotFound("Category not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(MockCategoryRepository)
			test.mockSetup(mockRepo)
			service := NewService(mockRepo)

			var result interface{}
			var err error

			switch input := test.input.(type) {
			case CreateCategoryRequest:
				result, err = service.Create(input)
			case string:
				result, err = service.List(input)
			case GetCategoryRequest:
				result, err = service.GetByName(input.UserID, input.Name)
			case UpdateCategoryRequest:
				result, err = service.Update(update.UpdateCategoryRequest{
					ID:      input.Name,
					UserID:  input.UserID,
					NewName: input.NewName,
				})
			case DeleteCategoryRequest:
				err = service.Delete(input)
			}

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			switch result := result.(type) {
			case *domain.Category:
				expectedCategory := test.expectedOutput.(*domain.Category)
				assert.Equal(t, expectedCategory.Name, result.Name)
				assert.Equal(t, expectedCategory.UserID, result.UserID)
				assert.NotEmpty(t, result.ID)
			case []*domain.Category:
				expectedCategories := test.expectedOutput.([]*domain.Category)
				assert.Equal(t, len(expectedCategories), len(result))
				for i, expectedCategory := range expectedCategories {
					assert.Equal(t, expectedCategory.Name, result[i].Name)
					assert.Equal(t, expectedCategory.UserID, result[i].UserID)
					assert.NotEmpty(t, result[i].ID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
