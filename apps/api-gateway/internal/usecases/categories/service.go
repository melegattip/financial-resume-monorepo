package categories

import (
	"strings"

	"github.com/google/uuid"
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/logs"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/create"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/delete"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/get"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/list"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/update"
)

// CategoryService define las operaciones disponibles para el servicio de categorías
type CategoryService interface {
	Create(request CreateCategoryRequest) (*domain.Category, error)
	GetByName(userID string, name string) (*domain.Category, error)
	GetByID(userID string, id string) (*domain.Category, error)
	List(userID string) ([]*domain.Category, error)
	Update(request update.UpdateCategoryRequest) (*update.UpdateCategoryResponse, error)
	Delete(request DeleteCategoryRequest) error
}

// Service coordina las operaciones de categorías
type Service struct {
	createService      *create.CreateCategory
	getService         *get.GetCategory
	listService        *list.ListCategories
	updateService      *update.UpdateCategory
	deleteService      *delete.DeleteCategory
	CategoryRepository baseRepo.CategoryRepository
}

// NewService crea una nueva instancia de Service
func NewService(repo baseRepo.CategoryRepository) CategoryService {
	return &Service{
		createService:      create.NewCreateCategory(repo),
		getService:         get.NewGetCategory(repo),
		listService:        list.NewListCategories(repo),
		updateService:      update.NewUpdateCategory(repo),
		deleteService:      delete.NewDeleteCategory(repo),
		CategoryRepository: repo,
	}
}

// Create crea una nueva categoría
func (s *Service) Create(request CreateCategoryRequest) (*domain.Category, error) {
	category := domain.NewCategoryBuilder().
		SetName(request.Name).
		SetUserID(request.UserID).
		Build()

	category, err := s.createService.Execute(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetByName obtiene una categoría por su nombre
func (s *Service) GetByName(userID string, name string) (*domain.Category, error) {
	return s.getService.Execute(userID, name)
}

// GetByID obtiene una categoría por su ID
func (s *Service) GetByID(userID string, id string) (*domain.Category, error) {
	return s.CategoryRepository.Get(userID, id)
}

// List obtiene todas las categorías de un usuario
func (s *Service) List(userID string) ([]*domain.Category, error) {
	return s.listService.Execute(userID)
}

// Update actualiza una categoría existente
func (s *Service) Update(request update.UpdateCategoryRequest) (*update.UpdateCategoryResponse, error) {
	return s.updateService.Execute(request)
}

// Delete elimina una categoría
func (s *Service) Delete(request DeleteCategoryRequest) error {
	return s.deleteService.Execute(request.UserID, request.ID)
}

type CreateCategory struct {
	CategoryRepository baseRepo.CategoryRepository
}

func NewCreateCategory(repo baseRepo.CategoryRepository) *CreateCategory {
	return &CreateCategory{CategoryRepository: repo}
}

func (s *CreateCategory) Execute(category *domain.Category) (*domain.Category, error) {
	// Validate category is not nil
	if category == nil {
		return nil, errors.NewBadRequest("Category cannot be nil")
	}

	// Validate empty name
	if strings.TrimSpace(category.Name) == "" {
		return nil, errors.NewBadRequest("Category name cannot be empty")
	}

	// Validate if category already exists
	existingCategory, err := s.CategoryRepository.GetByName(category.UserID, category.Name)
	if err != nil && !errors.IsResourceNotFound(err) {
		return nil, err
	}
	if existingCategory != nil {
		return nil, errors.NewResourceAlreadyExists(logs.ErrorCreatingCategory.GetMessage())
	}

	// Generate unique ID for category
	categoryID := "cat_" + uuid.New().String()[:8]

	// Set creation date
	category.ID = categoryID

	// Create category in repository
	err = s.CategoryRepository.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

type ListCategories struct {
	CategoryRepository baseRepo.CategoryRepository
}

func NewListCategories(repo baseRepo.CategoryRepository) *ListCategories {
	return &ListCategories{CategoryRepository: repo}
}

func (s *ListCategories) Execute(userID string) ([]*domain.Category, error) {
	return s.CategoryRepository.List(userID)
}

type GetCategory struct {
	CategoryRepository baseRepo.CategoryRepository
}

func NewGetCategory(repo baseRepo.CategoryRepository) *GetCategory {
	return &GetCategory{CategoryRepository: repo}
}

func (s *GetCategory) Execute(userID string, name string) (*domain.Category, error) {
	// Validate empty name
	if strings.TrimSpace(name) == "" {
		return nil, errors.NewBadRequest("Category name cannot be empty")
	}

	return s.CategoryRepository.GetByName(userID, name)
}

type UpdateCategory struct {
	CategoryRepository baseRepo.CategoryRepository
}

func NewUpdateCategory(repo baseRepo.CategoryRepository) *UpdateCategory {
	return &UpdateCategory{CategoryRepository: repo}
}

func (s *UpdateCategory) Execute(category *domain.Category) error {
	// Validate category is not nil
	if category == nil {
		return errors.NewBadRequest("Category cannot be nil")
	}

	// Validate empty name
	if strings.TrimSpace(category.Name) == "" {
		return errors.NewBadRequest("Category name cannot be empty")
	}

	// Validate invalid ID
	if !strings.HasPrefix(category.ID, "cat_") {
		return errors.NewBadRequest("Invalid category ID")
	}

	return s.CategoryRepository.Update(category)
}

type DeleteCategory struct {
	CategoryRepository baseRepo.CategoryRepository
}

func NewDeleteCategory(repo baseRepo.CategoryRepository) *DeleteCategory {
	return &DeleteCategory{CategoryRepository: repo}
}

func (s *DeleteCategory) Execute(userID string, id string) error {
	// Validate invalid ID
	if !strings.HasPrefix(id, "cat_") {
		return errors.NewBadRequest("Invalid category ID")
	}

	return s.CategoryRepository.Delete(userID, id)
}
