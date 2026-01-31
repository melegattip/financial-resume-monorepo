package context_test

import (
	"context"
	"testing"

	infrastructureCtx "github.com/melegattip/financial-resume-engine/internal/infrastructure/context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/assert"
)

func TestSetAction(t *testing.T) {
	ctx := context.Background()

	expectedAction := infrastructureCtx.SetAction(ctx, "TestAction")

	assert.NotNil(t, expectedAction)
	assert.Equal(t, "TestAction", expectedAction.Value(infrastructureCtx.ActionKey))
}

func TestNewGoroutineContextWithNewRelic(t *testing.T) {
	simulatedTxn := &newrelic.Transaction{
		Private: "TestTransaction",
	}

	newRelicCtx := newrelic.NewContext(context.Background(), simulatedTxn)

	expectedCtx := infrastructureCtx.NewGoroutineContext(newRelicCtx)

	assert.NotNil(t, expectedCtx)
	assert.NotEqualValues(t, expectedCtx, newRelicCtx)
}

func TestNewGoroutineContextWithoutNewRelic(t *testing.T) {
	expectedCtx := infrastructureCtx.NewGoroutineContext(context.Background())

	assert.NotNil(t, expectedCtx)
	assert.Equal(t, expectedCtx, context.Background())
}
