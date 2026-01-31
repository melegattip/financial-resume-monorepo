package services

import (
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/repository"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (ts *TransactionService) CreateTransaction(transaction *domain.Transaction) error {
	return ts.repo.Create(transaction)
}

func (ts *TransactionService) GetTransaction(id string) (*domain.Transaction, error) {
	transaction, err := ts.repo.Get(id)
	if err != nil {
		return nil, errors.NewResourceNotFound("transaction not found")
	}
	return transaction, nil
}

func (ts *TransactionService) ListTransactions() ([]*domain.Transaction, error) {
	return ts.repo.List()
}

func (ts *TransactionService) UpdateTransaction(transaction *domain.Transaction) error {
	return ts.repo.Update(transaction)
}

func (ts *TransactionService) DeleteTransaction(id string) error {
	return ts.repo.Delete(id)
}

func (ts *TransactionService) ListTransactionsByCategory(categoryID string) ([]*domain.Transaction, error) {
	return ts.repo.ListByCategory(categoryID)
}

func (ts *TransactionService) ListTransactionsByDateRange(startDate, endDate time.Time) ([]*domain.Transaction, error) {
	return ts.repo.ListByDateRange(startDate, endDate)
}

func (ts *TransactionService) ListTransactionsByType(transactionType string) ([]*domain.Transaction, error) {
	return ts.repo.ListByType(transactionType)
}

func (ts *TransactionService) GetMonthlyTotal(transactionType string, date time.Time) (float64, error) {
	return ts.repo.GetMonthlyTotal(transactionType, date)
}

func (ts *TransactionService) GetYearlyTotal(transactionType string, year int) (float64, error) {
	return ts.repo.GetYearlyTotal(transactionType, year)
}

func (ts *TransactionService) GetMonthlyTotalByCategory(categoryID string, date time.Time) (float64, error) {
	return ts.repo.GetMonthlyTotalByCategory(categoryID, date)
}

func (ts *TransactionService) GetYearlyTotalByCategory(categoryID string, year int) (float64, error) {
	return ts.repo.GetYearlyTotalByCategory(categoryID, year)
}
