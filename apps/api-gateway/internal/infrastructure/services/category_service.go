package services

import (
	"fmt"

	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// CategoryServiceImpl implementa CategoryService
type CategoryServiceImpl struct {
	categoryRepo baseRepo.CategoryRepository
}

// NewCategoryService crea una nueva instancia del servicio de categorías
func NewCategoryService(categoryRepo baseRepo.CategoryRepository) usecases.CategoryService {
	return &CategoryServiceImpl{
		categoryRepo: categoryRepo,
	}
}

// GetCategoryName obtiene el nombre de una categoría por su ID
func (s *CategoryServiceImpl) GetCategoryName(categoryID string) (string, error) {
	if categoryID == "" || categoryID == "sin-categoria" {
		return "Sin categoría", nil
	}
	// Por ahora retornamos el ID como fallback
	// Este método legacy no tiene user_id, así que no podemos consultar la BD
	return categoryID, nil
}

// GetCategoryNames obtiene los nombres de múltiples categorías
// NOTA: Este método tiene una limitación - no puede consultar la BD sin user_id
// Se usa principalmente desde contextos donde ya tenemos los nombres
func (s *CategoryServiceImpl) GetCategoryNames(categoryIDs []string) (map[string]string, error) {
	result := make(map[string]string)

	for _, categoryID := range categoryIDs {
		if categoryID == "" || categoryID == "sin-categoria" {
			result[categoryID] = "Sin categoría"
		} else {
			// Fallback: usar el ID como nombre
			result[categoryID] = categoryID
		}
	}

	return result, nil
}

// GetCategoryNamesWithUserID obtiene los nombres de categorías con user_id
// Este es el método que realmente consulta la base de datos
func (s *CategoryServiceImpl) GetCategoryNamesWithUserID(userID string, categoryIDs []string) (map[string]string, error) {
	result := make(map[string]string)

	if userID == "" {
		return nil, fmt.Errorf("user_id es requerido")
	}

	// Obtener todas las categorías del usuario
	categories, err := s.categoryRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo categorías del usuario: %w", err)
	}

	// Crear mapa de ID -> Nombre
	categoryMap := make(map[string]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}

	// Mapear los IDs solicitados
	for _, categoryID := range categoryIDs {
		if categoryID == "" || categoryID == "sin-categoria" {
			result[categoryID] = "Sin categoría"
		} else if name, exists := categoryMap[categoryID]; exists {
			result[categoryID] = name
		} else {
			result[categoryID] = "Categoría desconocida"
		}
	}

	return result, nil
}
