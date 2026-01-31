package update

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

type UpdateCategory struct {
	CategoryRepository baseRepo.CategoryRepository
}

func NewUpdateCategory(repo baseRepo.CategoryRepository) *UpdateCategory {
	return &UpdateCategory{
		CategoryRepository: repo,
	}
}

func (s *UpdateCategory) Execute(request UpdateCategoryRequest) (*UpdateCategoryResponse, error) {
	// Buscar la categoría existente por su ID
	existingCategory, err := s.CategoryRepository.Get(request.UserID, request.ID)
	if err != nil {
		return nil, err
	}

	// Verificar si ya existe una categoría con el nuevo nombre
	_, err = s.CategoryRepository.GetByName(request.UserID, request.NewName)
	if err == nil {
		// Si no hay error, significa que ya existe una categoría con ese nombre
		return nil, errors.NewResourceAlreadyExists("Category with new name already exists")
	} else if !errors.IsResourceNotFound(err) {
		// Si hay un error diferente a "no encontrado", propagarlo
		return nil, err
	}

	// Crear la categoría actualizada
	updatedCategory := domain.NewCategoryBuilder().
		SetID(existingCategory.ID).
		SetName(request.NewName).
		SetUserID(request.UserID).
		Build()

	// Actualizar la categoría
	err = s.CategoryRepository.Update(updatedCategory)
	if err != nil {
		return nil, err
	}

	// Actualizar las referencias en las transacciones que usan esta categoría
	err = s.CategoryRepository.UpdateCategoryReferences(updatedCategory.ID)
	if err != nil {
		return nil, err
	}

	// Crear la respuesta
	response := &UpdateCategoryResponse{
		ID:     updatedCategory.ID,
		UserID: updatedCategory.UserID,
		Name:   updatedCategory.Name,
	}

	return response, nil
}
