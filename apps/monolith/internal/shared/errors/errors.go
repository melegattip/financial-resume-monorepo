package errors

import "net/http"

// AppError represents an application-level error with an HTTP status code.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewNotFound creates a 404 Not Found error.
func NewNotFound(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

// NewValidation creates a 400 Bad Request error.
func NewValidation(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

// NewUnauthorized creates a 401 Unauthorized error.
func NewUnauthorized(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}

// NewForbidden creates a 403 Forbidden error.
func NewForbidden(message string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: message}
}

// NewInternal creates a 500 Internal Server Error.
func NewInternal(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}

// NewConflict creates a 409 Conflict error.
func NewConflict(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}
