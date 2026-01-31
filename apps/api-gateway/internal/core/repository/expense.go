package repository

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// ExpenseRepository define las operaciones para el repositorio de gastos
type ExpenseRepository interface {
	Create(expense *domain.Expense) error
	Get(userID string, id string) (*domain.Expense, error)
	List(userID string) ([]*domain.Expense, error)
	Update(expense *domain.Expense) error
	Delete(userID string, id string) error
	ListUnpaid(userID string) ([]*domain.Expense, error)
}
