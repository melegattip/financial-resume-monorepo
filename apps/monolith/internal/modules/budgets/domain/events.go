package domain

import "time"

// BudgetCreatedEvent is published when a new budget is created.
type BudgetCreatedEvent struct {
	BudgetID  string
	User      string
	Amount    float64
	Period    string
	Timestamp time.Time
}

func (e BudgetCreatedEvent) EventType() string   { return "budget.created" }
func (e BudgetCreatedEvent) AggregateID() string { return e.BudgetID }
func (e BudgetCreatedEvent) UserID() string      { return e.User }
func (e BudgetCreatedEvent) OccurredAt() string  { return e.Timestamp.Format(time.RFC3339) }

// BudgetUpdatedEvent is published when a budget is updated.
type BudgetUpdatedEvent struct {
	BudgetID  string
	User      string
	Timestamp time.Time
}

func (e BudgetUpdatedEvent) EventType() string   { return "budget.updated" }
func (e BudgetUpdatedEvent) AggregateID() string { return e.BudgetID }
func (e BudgetUpdatedEvent) UserID() string      { return e.User }
func (e BudgetUpdatedEvent) OccurredAt() string  { return e.Timestamp.Format(time.RFC3339) }

// BudgetDeletedEvent is published when a budget is deleted.
type BudgetDeletedEvent struct {
	BudgetID  string
	User      string
	Timestamp time.Time
}

func (e BudgetDeletedEvent) EventType() string   { return "budget.deleted" }
func (e BudgetDeletedEvent) AggregateID() string { return e.BudgetID }
func (e BudgetDeletedEvent) UserID() string      { return e.User }
func (e BudgetDeletedEvent) OccurredAt() string  { return e.Timestamp.Format(time.RFC3339) }
