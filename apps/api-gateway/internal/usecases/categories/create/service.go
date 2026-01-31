package create

import (
	"github.com/google/uuid"
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/logs"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

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
