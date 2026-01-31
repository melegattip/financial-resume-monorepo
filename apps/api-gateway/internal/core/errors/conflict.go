// Package errors provides custom error types for handling different error scenarios
package errors

// Conflict represents an error that occurs when there is a resource conflict
type Conflict struct {
	message string
}

// NewConflict creates a new Conflict error with the specified message
func NewConflict(message string) Conflict {
	return Conflict{
		message: message,
	}
}

// Error returns the error message
func (e Conflict) Error() string {
	return e.message
}
