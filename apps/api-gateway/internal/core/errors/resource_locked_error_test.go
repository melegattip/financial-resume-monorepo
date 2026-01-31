package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestResourceLockedError(t *testing.T) {
	err := errors.NewResourceLockedError("Resource locked error message")

	assert.True(t, errors.IsResourceLocked(err))
	assert.Equal(t, "Resource locked error message", err.Error())
}
