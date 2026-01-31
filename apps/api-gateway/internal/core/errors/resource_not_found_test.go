package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewResourceNotFound(t *testing.T) {
	message := "error message"
	err := errors.NewResourceNotFound(message)

	assert.Equal(t, message, err.Error())
}

func TestIsNewResourceNotFound(t *testing.T) {
	err := errors.NewResourceNotFound("error")
	assert.True(t, errors.IsResourceNotFound(err))
}
