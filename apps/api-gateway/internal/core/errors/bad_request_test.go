package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewBadRequestError(t *testing.T) {
	err := errors.NewBadRequest("error")
	assert.Equal(t, "error", err.Error())
}

func TestIsBadRequestError(t *testing.T) {
	err := errors.NewBadRequest("error")
	assert.True(t, errors.IsBadRequestError(err))
}
