package errors_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewDatabaseFailure(t *testing.T) {
	err := errors.NewDatabaseFailure("error")
	assert.Equal(t, "error", err.Error())
}
