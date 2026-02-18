package domain

import "time"

// ExpenseCreatedEvent is published when a new expense is created
type ExpenseCreatedEvent struct {
	ExpenseID       string
	User            string
	CategoryID      string
	Amount          float64
	Description     string
	TransactionDate time.Time
	Timestamp       time.Time
}

func (e ExpenseCreatedEvent) EventType() string {
	return "expense.created"
}

func (e ExpenseCreatedEvent) AggregateID() string {
	return e.ExpenseID
}

func (e ExpenseCreatedEvent) UserID() string {
	return e.User
}

func (e ExpenseCreatedEvent) OccurredAt() string {
	return e.Timestamp.Format(time.RFC3339)
}

// ExpenseUpdatedEvent is published when an expense is updated
type ExpenseUpdatedEvent struct {
	ExpenseID string
	User      string
	Timestamp time.Time
}

func (e ExpenseUpdatedEvent) EventType() string {
	return "expense.updated"
}

func (e ExpenseUpdatedEvent) AggregateID() string {
	return e.ExpenseID
}

func (e ExpenseUpdatedEvent) UserID() string {
	return e.User
}

func (e ExpenseUpdatedEvent) OccurredAt() string {
	return e.Timestamp.Format(time.RFC3339)
}

// ExpenseDeletedEvent is published when an expense is deleted
type ExpenseDeletedEvent struct {
	ExpenseID string
	User      string
	Amount    float64
	Timestamp time.Time
}

func (e ExpenseDeletedEvent) EventType() string {
	return "expense.deleted"
}

func (e ExpenseDeletedEvent) AggregateID() string {
	return e.ExpenseID
}

func (e ExpenseDeletedEvent) UserID() string {
	return e.User
}

func (e ExpenseDeletedEvent) OccurredAt() string {
	return e.Timestamp.Format(time.RFC3339)
}

// IncomeCreatedEvent is published when a new income is created
type IncomeCreatedEvent struct {
	IncomeID     string
	User         string
	Amount       float64
	Source       string
	Description  string
	ReceivedDate time.Time
	Timestamp    time.Time
}

func (i IncomeCreatedEvent) EventType() string {
	return "income.created"
}

func (i IncomeCreatedEvent) AggregateID() string {
	return i.IncomeID
}

func (i IncomeCreatedEvent) UserID() string {
	return i.User
}

func (i IncomeCreatedEvent) OccurredAt() string {
	return i.Timestamp.Format(time.RFC3339)
}

// IncomeUpdatedEvent is published when an income is updated
type IncomeUpdatedEvent struct {
	IncomeID  string
	User      string
	Timestamp time.Time
}

func (i IncomeUpdatedEvent) EventType() string {
	return "income.updated"
}

func (i IncomeUpdatedEvent) AggregateID() string {
	return i.IncomeID
}

func (i IncomeUpdatedEvent) UserID() string {
	return i.User
}

func (i IncomeUpdatedEvent) OccurredAt() string {
	return i.Timestamp.Format(time.RFC3339)
}

// IncomeDeletedEvent is published when an income is deleted
type IncomeDeletedEvent struct {
	IncomeID  string
	User      string
	Amount    float64
	Timestamp time.Time
}

func (i IncomeDeletedEvent) EventType() string {
	return "income.deleted"
}

func (i IncomeDeletedEvent) AggregateID() string {
	return i.IncomeID
}

func (i IncomeDeletedEvent) UserID() string {
	return i.User
}

func (i IncomeDeletedEvent) OccurredAt() string {
	return i.Timestamp.Format(time.RFC3339)
}
