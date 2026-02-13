package ports

import "context"

// Event is the base interface for all domain events.
type Event interface {
	// EventType returns the event type identifier (e.g., "expense.created").
	EventType() string
	// AggregateID returns the ID of the aggregate that produced this event.
	AggregateID() string
	// UserID returns the ID of the user who triggered this event.
	UserID() string
	// OccurredAt returns the timestamp when the event occurred.
	OccurredAt() string
}

// EventHandler processes a domain event.
type EventHandler func(ctx context.Context, event Event) error

// EventBus defines the publish/subscribe contract for domain events.
type EventBus interface {
	// Publish dispatches an event to all registered handlers for its type.
	// Handlers execute asynchronously. Publishing to a type with no subscribers is a no-op.
	Publish(ctx context.Context, event Event) error
	// Subscribe registers a handler for a specific event type.
	// Multiple handlers can be registered for the same event type.
	Subscribe(eventType string, handler EventHandler)
}
