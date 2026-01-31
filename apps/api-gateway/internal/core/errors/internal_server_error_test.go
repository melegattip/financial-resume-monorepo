package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestInternalServerError(t *testing.T) {
	message := "Internal server error message"
	err := errors.NewInternalServerError(message)

	assert.Equal(t, message, err.Error())
}
