package errors

import (
	"errors"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Common errors
var (
	ErrNotFound       = &AppError{Code: http.StatusNotFound, Message: "resource not found"}
	ErrUnauthorized   = &AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden      = &AppError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrBadRequest     = &AppError{Code: http.StatusBadRequest, Message: "bad request"}
	ErrInternalServer = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrConflict       = &AppError{Code: http.StatusConflict, Message: "resource already exists"}
)

// New creates a new AppError
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap wraps an error with AppError
func Wrap(err error, code int, message string) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// IsNotFound checks if error is a not found error
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == http.StatusNotFound
	}
	return false
}

// IsUnauthorized checks if error is an unauthorized error
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == http.StatusUnauthorized
	}
	return false
}

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err *AppError) ErrorResponse {
	return ErrorResponse{
		Error: err.Message,
		Code:  err.Code,
	}
}
