package domain

import "errors"

var (
	// ErrInvalidAmount se devuelve cuando el monto de la transacción no es válido
	ErrInvalidAmount = errors.New("invalid amount")

	// ErrEmptyDescription se devuelve cuando la descripción de la transacción está vacía
	ErrEmptyDescription = errors.New("empty description")

	// ErrInvalidCategory se devuelve cuando la categoría de la transacción no es válida
	ErrInvalidCategory = errors.New("invalid category")

	// ErrInvalidDate se devuelve cuando la fecha de la transacción no es válida
	ErrInvalidDate = errors.New("invalid date")

	// ErrAmountTooLarge se devuelve cuando el monto de la transacción es demasiado grande
	ErrAmountTooLarge = errors.New("amount too large")
)
