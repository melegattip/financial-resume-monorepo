package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewTooManyRequests(t *testing.T) {
	message := "Too many requests message"
	err := errors.NewTooManyRequests(message)

	assert.Equal(t, message, err.Error())
}

func TestIsTooManyRequestsError(t *testing.T) {
	err := errors.NewTooManyRequests("Too many requests message")
	assert.True(t, errors.IsTooManyRequestError(err))
}
