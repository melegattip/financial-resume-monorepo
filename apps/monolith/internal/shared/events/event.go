package events

import "time"

// DomainEvent is the base struct implementing the ports.Event interface.
// All module-specific events should embed this struct.
type DomainEvent struct {
	Type      string    `json:"type"`
	Aggregate string    `json:"aggregate_id"`
	User      string    `json:"user_id"`
	Timestamp time.Time `json:"occurred_at"`
	Data      any       `json:"payload,omitempty"`
}

// EventType returns the event type identifier.
func (e DomainEvent) EventType() string { return e.Type }

// AggregateID returns the ID of the aggregate that produced this event.
func (e DomainEvent) AggregateID() string { return e.Aggregate }

// UserID returns the ID of the user who triggered this event.
func (e DomainEvent) UserID() string { return e.User }

// OccurredAt returns the timestamp as RFC3339 string.
func (e DomainEvent) OccurredAt() string { return e.Timestamp.Format(time.RFC3339) }

// NewDomainEvent creates a new DomainEvent with the current timestamp.
func NewDomainEvent(eventType, aggregateID, userID string, payload any) DomainEvent {
	return DomainEvent{
		Type:      eventType,
		Aggregate: aggregateID,
		User:      userID,
		Timestamp: time.Now().UTC(),
		Data:      payload,
	}
}
