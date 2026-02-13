package events

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

func newTestBus() *InMemoryEventBus {
	logger := zerolog.Nop()
	return NewInMemoryEventBus(logger)
}

func newTestEvent(eventType string) ports.Event {
	return NewDomainEvent(eventType, "agg-1", "user-1", nil)
}

func TestPublish_SingleSubscriber(t *testing.T) {
	bus := newTestBus()
	var received atomic.Bool

	bus.Subscribe("test.event", func(ctx context.Context, event ports.Event) error {
		received.Store(true)
		assert.Equal(t, "test.event", event.EventType())
		assert.Equal(t, "agg-1", event.AggregateID())
		assert.Equal(t, "user-1", event.UserID())
		return nil
	})

	err := bus.Publish(context.Background(), newTestEvent("test.event"))
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	assert.True(t, received.Load(), "subscriber should have received the event")
}

func TestPublish_MultipleSubscribers(t *testing.T) {
	bus := newTestBus()
	var count atomic.Int32

	for i := 0; i < 3; i++ {
		bus.Subscribe("test.event", func(ctx context.Context, event ports.Event) error {
			count.Add(1)
			return nil
		})
	}

	err := bus.Publish(context.Background(), newTestEvent("test.event"))
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, int32(3), count.Load(), "all 3 subscribers should receive the event")
}

func TestPublish_HandlerErrorDoesNotAffectOthers(t *testing.T) {
	bus := newTestBus()
	var received atomic.Bool

	bus.Subscribe("test.event", func(ctx context.Context, event ports.Event) error {
		return errors.New("handler failure")
	})

	bus.Subscribe("test.event", func(ctx context.Context, event ports.Event) error {
		received.Store(true)
		return nil
	})

	err := bus.Publish(context.Background(), newTestEvent("test.event"))
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	assert.True(t, received.Load(), "second handler should still execute despite first handler error")
}

func TestPublish_HandlerPanicDoesNotAffectOthers(t *testing.T) {
	bus := newTestBus()
	var received atomic.Bool

	bus.Subscribe("test.event", func(ctx context.Context, event ports.Event) error {
		panic("handler panic")
	})

	bus.Subscribe("test.event", func(ctx context.Context, event ports.Event) error {
		received.Store(true)
		return nil
	})

	err := bus.Publish(context.Background(), newTestEvent("test.event"))
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	assert.True(t, received.Load(), "second handler should still execute despite first handler panic")
}

func TestPublish_NoSubscribers(t *testing.T) {
	bus := newTestBus()

	err := bus.Publish(context.Background(), newTestEvent("unhandled.event"))
	assert.NoError(t, err, "publishing to no subscribers should not error")
}

func TestPublish_ConcurrentSafe(t *testing.T) {
	bus := newTestBus()
	var count atomic.Int32

	bus.Subscribe("test.event", func(ctx context.Context, event ports.Event) error {
		count.Add(1)
		return nil
	})

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = bus.Publish(context.Background(), newTestEvent("test.event"))
		}()
	}
	wg.Wait()

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, int32(100), count.Load(), "all 100 concurrent publishes should be handled")
}

func TestDomainEvent_OccurredAt(t *testing.T) {
	event := NewDomainEvent("test.event", "agg-1", "user-1", nil)
	ts := event.OccurredAt()
	assert.NotEmpty(t, ts, "OccurredAt should return a non-empty RFC3339 timestamp")
}

func TestSubscribe_ConcurrentSafe(t *testing.T) {
	bus := newTestBus()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bus.Subscribe("test.event", func(ctx context.Context, event ports.Event) error {
				return nil
			})
		}()
	}
	wg.Wait()

	bus.mu.RLock()
	count := len(bus.handlers["test.event"])
	bus.mu.RUnlock()
	assert.Equal(t, 100, count, "all 100 concurrent subscribes should register")
}
