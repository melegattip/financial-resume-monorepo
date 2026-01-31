package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewUnauthorizedRequest(t *testing.T) {
	message := "Unauthorized request message"
	err := errors.NewUnauthorizedRequest(message)

	assert.Equal(t, message, err.Error())
}
