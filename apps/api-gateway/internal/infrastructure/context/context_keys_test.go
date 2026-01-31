package context_test

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/infrastructure/context"
	"github.com/stretchr/testify/assert"
)

func TestContentTypeValues(t *testing.T) {
	assert.Equal(t, "action", context.ActionKey.String())
	assert.Equal(t, "x-caller-id", context.XCallerIDKey.String())
}
