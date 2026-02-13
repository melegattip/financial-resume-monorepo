package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := &AppError{Code: 404, Message: "not found"}
	assert.Equal(t, "not found", err.Error())
}

func TestNewNotFound(t *testing.T) {
	err := NewNotFound("resource missing")
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Equal(t, "resource missing", err.Message)
}

func TestNewValidation(t *testing.T) {
	err := NewValidation("invalid input")
	assert.Equal(t, http.StatusBadRequest, err.Code)
}

func TestNewUnauthorized(t *testing.T) {
	err := NewUnauthorized("no token")
	assert.Equal(t, http.StatusUnauthorized, err.Code)
}

func TestNewForbidden(t *testing.T) {
	err := NewForbidden("insufficient permissions")
	assert.Equal(t, http.StatusForbidden, err.Code)
}

func TestNewInternal(t *testing.T) {
	err := NewInternal("server error")
	assert.Equal(t, http.StatusInternalServerError, err.Code)
}

func TestNewConflict(t *testing.T) {
	err := NewConflict("already exists")
	assert.Equal(t, http.StatusConflict, err.Code)
}
