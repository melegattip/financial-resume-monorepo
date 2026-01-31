package repository

import "errors"

var (
	// ErrNotFound es retornado cuando un recurso no se encuentra
	ErrNotFound = errors.New("recurso no encontrado")
)
