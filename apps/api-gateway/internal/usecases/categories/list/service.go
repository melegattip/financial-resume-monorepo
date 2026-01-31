package list

import (
	"strings"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

type ListCategories struct {
	CategoryRepository baseRepo.CategoryRepository
}

func NewListCategories(repo baseRepo.CategoryRepository) *ListCategories {
	return &ListCategories{
		CategoryRepository: repo,
	}
}

func (s *ListCategories) Execute(userID string) ([]*domain.Category, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.NewBadRequest("User ID cannot be empty")
	}

	categories, err := s.CategoryRepository.List(userID)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
