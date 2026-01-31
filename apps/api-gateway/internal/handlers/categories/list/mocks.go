package list

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/update"
	"github.com/stretchr/testify/mock"
)

// MockCategoryService es un mock del servicio de categorías
type MockCategoryService struct {
	mock.Mock
}

// Create implementa el método Create de la interfaz CategoryService
func (m *MockCategoryService) Create(request categories.CreateCategoryRequest) (*domain.Category, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

// GetByName implementa el método GetByName de la interfaz CategoryService
func (m *MockCategoryService) GetByName(userID string, name string) (*domain.Category, error) {
	args := m.Called(userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

// GetByID implementa el método GetByID de la interfaz CategoryService
func (m *MockCategoryService) GetByID(userID string, id string) (*domain.Category, error) {
	args := m.Called(userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

// List implementa el método List de la interfaz CategoryService
func (m *MockCategoryService) List(userID string) ([]*domain.Category, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

// Update implementa el método Update de la interfaz CategoryService
func (m *MockCategoryService) Update(request update.UpdateCategoryRequest) (*update.UpdateCategoryResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*update.UpdateCategoryResponse), args.Error(1)
}

// Delete implementa el método Delete de la interfaz CategoryService
func (m *MockCategoryService) Delete(request categories.DeleteCategoryRequest) error {
	args := m.Called(request)
	return args.Error(0)
}
