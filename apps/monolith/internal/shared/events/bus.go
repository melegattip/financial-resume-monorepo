package events

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// InMemoryEventBus implements ports.EventBus with in-memory publish/subscribe.
// Handlers execute asynchronously in separate goroutines.
// Thread-safe for concurrent publish and subscribe operations.
type InMemoryEventBus struct {
	handlers map[string][]ports.EventHandler
	mu       sync.RWMutex
	logger   zerolog.Logger
}

// NewInMemoryEventBus creates a new in-memory event bus.
func NewInMemoryEventBus(logger zerolog.Logger) *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]ports.EventHandler),
		logger:   logger,
	}
}

// Publish dispatches an event to all registered handlers for its type.
// Handlers run in separate goroutines. If no subscribers exist, this is a no-op.
func (bus *InMemoryEventBus) Publish(ctx context.Context, event ports.Event) error {
	bus.mu.RLock()
	handlers := make([]ports.EventHandler, len(bus.handlers[event.EventType()]))
	copy(handlers, bus.handlers[event.EventType()])
	bus.mu.RUnlock()

	for _, handler := range handlers {
		go bus.safeExecute(ctx, event, handler)
	}

	return nil
}

// Subscribe registers a handler for a specific event type.
func (bus *InMemoryEventBus) Subscribe(eventType string, handler ports.EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
}

func (bus *InMemoryEventBus) safeExecute(ctx context.Context, event ports.Event, handler ports.EventHandler) {
	defer func() {
		if r := recover(); r != nil {
			bus.logger.Error().
				Str("event_type", event.EventType()).
				Str("panic", fmt.Sprintf("%v", r)).
				Msg("event handler panicked")
		}
	}()

	if err := handler(ctx, event); err != nil {
		bus.logger.Error().
			Err(err).
			Str("event_type", event.EventType()).
			Msg("event handler failed")
	}
}
