package repository

import (
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// BaseTransactionRepository es una estructura genérica para repositorios de transacciones
type BaseTransactionRepository[T domain.Transaction] struct {
	factory domain.TransactionFactory
}

// NewBaseTransactionRepository crea una nueva instancia de BaseTransactionRepository
func NewBaseTransactionRepository[T domain.Transaction](factory domain.TransactionFactory) *BaseTransactionRepository[T] {
	return &BaseTransactionRepository[T]{
		factory: factory,
	}
}

// GetFactory retorna la fábrica de transacciones
func (r *BaseTransactionRepository[T]) GetFactory() domain.TransactionFactory {
	return r.factory
}
