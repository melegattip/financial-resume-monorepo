package categories

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/update"
	"github.com/stretchr/testify/mock"
)

// MockCategoryService implementa un mock del servicio de categorías para pruebas
type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) Create(request categories.CreateCategoryRequest) (*domain.Category, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryService) GetByName(userID string, name string) (*domain.Category, error) {
	args := m.Called(userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryService) GetByID(userID string, id string) (*domain.Category, error) {
	args := m.Called(userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryService) List(userID string) ([]*domain.Category, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryService) Update(request update.UpdateCategoryRequest) (*update.UpdateCategoryResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*update.UpdateCategoryResponse), args.Error(1)
}

func (m *MockCategoryService) Delete(request categories.DeleteCategoryRequest) error {
	args := m.Called(request)
	return args.Error(0)
}
