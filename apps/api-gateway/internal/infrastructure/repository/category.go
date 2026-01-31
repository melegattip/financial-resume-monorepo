package repository

import (
	"log"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"gorm.io/gorm"
)

// Category implementa el repositorio específico para categorías
type Category struct {
	db *gorm.DB
}

// NewCategoryRepository crea una nueva instancia del repositorio de categorías
func NewCategoryRepository(db *gorm.DB) baseRepo.CategoryRepository {
	return &Category{
		db: db,
	}
}

func (r *Category) Create(category *domain.Category) error {
	log.Printf("🔍 [CategoryRepo] Iniciando creación de categoría: %+v", category)

	err := r.db.Create(category).Error
	if err != nil {
		log.Printf("❌ [CategoryRepo] Error creando categoría: %v", err)
		return err
	}

	log.Printf("✅ [CategoryRepo] Categoría creada exitosamente: %s", category.ID)
	return nil
}

func (r *Category) Get(userID string, id string) (*domain.Category, error) {
	var category domain.Category
	result := r.db.Where("user_id = ? AND id = ?", userID, id).First(&category)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewResourceNotFound("Category not found")
		}
		return nil, result.Error
	}
	return &category, nil
}

func (r *Category) GetByName(userID string, name string) (*domain.Category, error) {
	var category domain.Category
	result := r.db.Where("user_id = ? AND name = ?", userID, name).First(&category)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewResourceNotFound("Category not found")
		}
		return nil, result.Error
	}
	return &category, nil
}

func (r *Category) List(userID string) ([]*domain.Category, error) {
	var categories []domain.Category
	result := r.db.Where("user_id = ?", userID).Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	categoryPointers := make([]*domain.Category, len(categories))
	for i := range categories {
		categoryPointers[i] = &categories[i]
	}
	return categoryPointers, nil
}

func (r *Category) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *Category) Delete(userID string, id string) error {
	result := r.db.Where("user_id = ? AND id = ?", userID, id).Delete(&domain.Category{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewResourceNotFound("Category not found")
	}
	return nil
}

// UpdateCategoryReferences actualiza todas las transacciones que tienen la categoría con el ID especificado
func (r *Category) UpdateCategoryReferences(categoryID string) error {
	// Iniciar una transacción
	tx := r.db.Begin()

	// En caso de error, realizar rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// No necesitamos actualizar ninguna transacción, ya que las relaciones
	// ahora están basadas en el ID de la categoría, que no cambia al modificar el nombre
	// Las transacciones ya referencian correctamente a la categoría a través del category_id

	return tx.Commit().Error
}
