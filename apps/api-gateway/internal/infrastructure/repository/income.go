package repository

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"gorm.io/gorm"
)

// Income implementa el repositorio específico para ingresos
type Income struct {
	*baseRepo.BaseTransactionRepository[*domain.Income]
	db *gorm.DB
}

// NewIncomeRepository crea una nueva instancia del repositorio de ingresos
func NewIncomeRepository(db *gorm.DB) *Income {
	return &Income{
		BaseTransactionRepository: baseRepo.NewBaseTransactionRepository[*domain.Income](domain.NewIncomeFactory()),
		db:                        db,
	}
}

func (r *Income) Create(income *domain.Income) error {
	return r.db.Create(income).Error
}

func (r *Income) Get(userID string, id string) (*domain.Income, error) {
	var income domain.Income
	result := r.db.Where("user_id = ? AND id = ?", userID, id).First(&income)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewResourceNotFound("Income not found")
		}
		return nil, result.Error
	}
	return &income, nil
}

func (r *Income) List(userID string) ([]*domain.Income, error) {
	var incomes []domain.Income
	result := r.db.Where("user_id = ?", userID).Find(&incomes)
	if result.Error != nil {
		return nil, result.Error
	}

	incomePointers := make([]*domain.Income, len(incomes))
	for i := range incomes {
		incomePointers[i] = &incomes[i]
	}
	return incomePointers, nil
}

func (r *Income) Update(income *domain.Income) error {
	return r.db.Save(income).Error
}

func (r *Income) Delete(userID string, id string) error {
	result := r.db.Where("user_id = ? AND id = ?", userID, id).Delete(&domain.Income{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewResourceNotFound("Income not found")
	}
	return nil
}

func (r *Income) ListUnreceived(userID string) ([]*domain.Income, error) {
	var incomes []domain.Income
	result := r.db.Where("user_id = ? AND received = ?", userID, false).Find(&incomes)
	if result.Error != nil {
		return nil, result.Error
	}

	incomePointers := make([]*domain.Income, len(incomes))
	for i := range incomes {
		incomePointers[i] = &incomes[i]
	}
	return incomePointers, nil
}
