package repository

import (
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// TransactionRepository define las operaciones para el repositorio de transacciones
type TransactionRepository interface {
	Create(transaction *domain.Transaction) error
	Get(id string) (*domain.Transaction, error)
	List() ([]*domain.Transaction, error)
	Update(transaction *domain.Transaction) error
	Delete(id string) error
	ListByCategory(categoryID string) ([]*domain.Transaction, error)
	ListByDateRange(startDate, endDate time.Time) ([]*domain.Transaction, error)
	ListByType(transactionType string) ([]*domain.Transaction, error)
	GetMonthlyTotal(transactionType string, date time.Time) (float64, error)
	GetYearlyTotal(transactionType string, year int) (float64, error)
	GetMonthlyTotalByCategory(categoryID string, date time.Time) (float64, error)
	GetYearlyTotalByCategory(categoryID string, year int) (float64, error)
}
