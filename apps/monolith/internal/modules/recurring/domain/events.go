package domain

import "time"

// RecurringTransactionCreatedEvent is published when a new recurring transaction is created
type RecurringTransactionCreatedEvent struct {
	RecurringID string
	User        string
	Type        string
	Amount      float64
	Frequency   string
	NextDate    time.Time
	Timestamp   time.Time
}

func (e RecurringTransactionCreatedEvent) EventType() string { return "recurring.created" }
func (e RecurringTransactionCreatedEvent) AggregateID() string { return e.RecurringID }
func (e RecurringTransactionCreatedEvent) UserID() string { return e.User }
func (e RecurringTransactionCreatedEvent) OccurredAt() string {
	return e.Timestamp.Format(time.RFC3339)
}

// RecurringTransactionUpdatedEvent is published when a recurring transaction is updated
type RecurringTransactionUpdatedEvent struct {
	RecurringID string
	User        string
	Timestamp   time.Time
}

func (e RecurringTransactionUpdatedEvent) EventType() string  { return "recurring.updated" }
func (e RecurringTransactionUpdatedEvent) AggregateID() string { return e.RecurringID }
func (e RecurringTransactionUpdatedEvent) UserID() string      { return e.User }
func (e RecurringTransactionUpdatedEvent) OccurredAt() string {
	return e.Timestamp.Format(time.RFC3339)
}

// RecurringTransactionDeletedEvent is published when a recurring transaction is deleted
type RecurringTransactionDeletedEvent struct {
	RecurringID string
	User        string
	Timestamp   time.Time
}

func (e RecurringTransactionDeletedEvent) EventType() string  { return "recurring.deleted" }
func (e RecurringTransactionDeletedEvent) AggregateID() string { return e.RecurringID }
func (e RecurringTransactionDeletedEvent) UserID() string      { return e.User }
func (e RecurringTransactionDeletedEvent) OccurredAt() string {
	return e.Timestamp.Format(time.RFC3339)
}

// RecurringTransactionExecutedEvent is published when a recurring transaction is executed
// (either manually or by the cron job)
type RecurringTransactionExecutedEvent struct {
	RecurringID   string
	User          string
	Type          string  // "income" or "expense"
	Amount        float64
	TransactionID string // ID of the newly created expense or income
	Timestamp     time.Time
}

func (e RecurringTransactionExecutedEvent) EventType() string  { return "recurring.executed" }
func (e RecurringTransactionExecutedEvent) AggregateID() string { return e.RecurringID }
func (e RecurringTransactionExecutedEvent) UserID() string      { return e.User }
func (e RecurringTransactionExecutedEvent) OccurredAt() string {
	return e.Timestamp.Format(time.RFC3339)
}

// RecurringTransactionPausedEvent is published when a recurring transaction is paused
type RecurringTransactionPausedEvent struct {
	RecurringID string
	User        string
	Timestamp   time.Time
}

func (e RecurringTransactionPausedEvent) EventType() string  { return "recurring.paused" }
func (e RecurringTransactionPausedEvent) AggregateID() string { return e.RecurringID }
func (e RecurringTransactionPausedEvent) UserID() string      { return e.User }
func (e RecurringTransactionPausedEvent) OccurredAt() string {
	return e.Timestamp.Format(time.RFC3339)
}

// RecurringTransactionResumedEvent is published when a recurring transaction is resumed
type RecurringTransactionResumedEvent struct {
	RecurringID string
	User        string
	Timestamp   time.Time
}

func (e RecurringTransactionResumedEvent) EventType() string  { return "recurring.resumed" }
func (e RecurringTransactionResumedEvent) AggregateID() string { return e.RecurringID }
func (e RecurringTransactionResumedEvent) UserID() string      { return e.User }
func (e RecurringTransactionResumedEvent) OccurredAt() string {
	return e.Timestamp.Format(time.RFC3339)
}
