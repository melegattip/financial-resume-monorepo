package repository

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"gorm.io/gorm"
)

// Expense implementa el repositorio específico para gastos
type Expense struct {
	*baseRepo.BaseTransactionRepository[*domain.Expense]
	db *gorm.DB
}

// NewExpenseRepository crea una nueva instancia del repositorio de gastos
func NewExpenseRepository(db *gorm.DB) *Expense {
	return &Expense{
		BaseTransactionRepository: baseRepo.NewBaseTransactionRepository[*domain.Expense](domain.NewExpenseFactory()),
		db:                        db,
	}
}

func (r *Expense) Create(expense *domain.Expense) error {
	return r.db.Create(expense).Error
}

func (r *Expense) Get(userID string, id string) (*domain.Expense, error) {
	var expense domain.Expense
	result := r.db.Where("user_id = ? AND id = ?", userID, id).First(&expense)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.NewResourceNotFound("Expense not found")
		}
		return nil, result.Error
	}
	return &expense, nil
}

func (r *Expense) List(userID string) ([]*domain.Expense, error) {
	var expenses []domain.Expense
	result := r.db.Where("user_id = ?", userID).Find(&expenses)
	if result.Error != nil {
		return nil, result.Error
	}

	expensePointers := make([]*domain.Expense, len(expenses))
	for i := range expenses {
		expensePointers[i] = &expenses[i]
	}
	return expensePointers, nil
}

func (r *Expense) Update(expense *domain.Expense) error {
	return r.db.Save(expense).Error
}

func (r *Expense) Delete(userID string, id string) error {
	result := r.db.Where("user_id = ? AND id = ?", userID, id).Delete(&domain.Expense{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NewResourceNotFound("Expense not found")
	}
	return nil
}

// ListUnpaid retorna todos los gastos no pagados de un usuario
func (r *Expense) ListUnpaid(userID string) ([]*domain.Expense, error) {
	var expenses []domain.Expense
	result := r.db.Where("user_id = ? AND paid = ?", userID, false).Find(&expenses)
	if result.Error != nil {
		return nil, result.Error
	}

	expensePointers := make([]*domain.Expense, len(expenses))
	for i := range expenses {
		expensePointers[i] = &expenses[i]
	}
	return expensePointers, nil
}
