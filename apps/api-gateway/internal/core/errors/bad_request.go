// Package errors provides custom error types for handling different error scenarios
package errors

import "reflect"

// BadRequest represents an error that occurs when the client sends an invalid request
type BadRequest struct {
	message string
}

// NewBadRequest creates a new BadRequest error with the specified message
func NewBadRequest(message string) *BadRequest {
	return &BadRequest{
		message: message,
	}
}

// Error returns the error message
func (e *BadRequest) Error() string {
	return e.message
}

// IsBadRequestError checks if the given error is of type BadRequest
func IsBadRequestError(err error) bool {
	return reflect.TypeOf(err).String() == "*errors.BadRequest"
}
