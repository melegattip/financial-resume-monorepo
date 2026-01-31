package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewResourceAlreadyExists(t *testing.T) {
	message := "error message"
	err := errors.NewResourceAlreadyExists(message)
	assert.Equal(t, message, err.Error())
}

func TestIsResourceAlreadyExists(t *testing.T) {
	err := errors.NewResourceAlreadyExists("error message")
	assert.True(t, errors.IsResourceAlreadyExists(err))
}
