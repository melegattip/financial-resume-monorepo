package domain

import "errors"

var (
	// ErrInvalidTransactionType se devuelve cuando el tipo de transacción no es válido
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)
