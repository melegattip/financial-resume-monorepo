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

// BudgetThresholdCrossedEvent is published when a budget crosses a status threshold
// (on_track → warning, on_track → exceeded, or warning → exceeded).
type BudgetThresholdCrossedEvent struct {
	BudgetID    string
	User        string
	TenantID    string
	CategoryID  string
	SpentAmount float64
	BudgetLimit float64
	SpentPct    float64
	Period      string
	NewStatus   string
	Timestamp   time.Time
}

func (e BudgetThresholdCrossedEvent) EventType() string   { return "budget.threshold_crossed" }
func (e BudgetThresholdCrossedEvent) AggregateID() string { return e.BudgetID }
func (e BudgetThresholdCrossedEvent) UserID() string      { return e.User }
func (e BudgetThresholdCrossedEvent) OccurredAt() string  { return e.Timestamp.Format(time.RFC3339) }
