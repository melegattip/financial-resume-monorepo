package delete

import (
	"strings"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

type DeleteCategory struct {
	CategoryRepository baseRepo.CategoryRepository
}

func NewDeleteCategory(repo baseRepo.CategoryRepository) *DeleteCategory {
	return &DeleteCategory{
		CategoryRepository: repo,
	}
}

func (s *DeleteCategory) Execute(userID string, id string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.NewBadRequest("User ID cannot be empty")
	}

	if strings.TrimSpace(id) == "" {
		return errors.NewBadRequest("Category name is required")
	}

	return s.CategoryRepository.Delete(userID, id)
}
