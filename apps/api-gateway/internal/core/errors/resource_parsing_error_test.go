package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewResourceParsingError(t *testing.T) {
	message := "Resource parsing error message"
	err := errors.NewResourceParsingError(message)

	assert.Equal(t, message, err.Error())
}
