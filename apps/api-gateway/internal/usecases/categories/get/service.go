package get

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

type GetCategory struct {
	CategoryRepository baseRepo.CategoryRepository
}

func NewGetCategory(repo baseRepo.CategoryRepository) *GetCategory {
	return &GetCategory{CategoryRepository: repo}
}

func (s *GetCategory) Execute(userID string, name string) (*domain.Category, error) {
	return s.CategoryRepository.GetByName(userID, name)
}
