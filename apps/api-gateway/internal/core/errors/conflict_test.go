package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewConflict(t *testing.T) {
	message := "conflict error message"
	err := errors.NewConflict(message)

	assert.Equal(t, message, err.Error())
}
