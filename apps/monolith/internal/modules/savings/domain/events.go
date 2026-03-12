package domain

import "time"

// SavingsGoalCreatedEvent is published when a new savings goal is created.
type SavingsGoalCreatedEvent struct {
	GoalID    string
	User      string
	Name      string
	Amount    float64
	Timestamp time.Time
}

func (e SavingsGoalCreatedEvent) EventType() string   { return "savings_goal.created" }
func (e SavingsGoalCreatedEvent) AggregateID() string { return e.GoalID }
func (e SavingsGoalCreatedEvent) UserID() string      { return e.User }
func (e SavingsGoalCreatedEvent) OccurredAt() string  { return e.Timestamp.Format(time.RFC3339) }

// SavingsGoalUpdatedEvent is published when a savings goal is updated.
type SavingsGoalUpdatedEvent struct {
	GoalID    string
	User      string
	Timestamp time.Time
}

func (e SavingsGoalUpdatedEvent) EventType() string   { return "savings_goal.updated" }
func (e SavingsGoalUpdatedEvent) AggregateID() string { return e.GoalID }
func (e SavingsGoalUpdatedEvent) UserID() string      { return e.User }
func (e SavingsGoalUpdatedEvent) OccurredAt() string  { return e.Timestamp.Format(time.RFC3339) }

// SavingsGoalAchievedEvent is published when a savings goal reaches 100% completion.
type SavingsGoalAchievedEvent struct {
	GoalID    string
	User      string
	GoalName  string
	Amount    float64
	Timestamp time.Time
}

func (e SavingsGoalAchievedEvent) EventType() string   { return "savings_goal.achieved" }
func (e SavingsGoalAchievedEvent) AggregateID() string { return e.GoalID }
func (e SavingsGoalAchievedEvent) UserID() string      { return e.User }
func (e SavingsGoalAchievedEvent) OccurredAt() string  { return e.Timestamp.Format(time.RFC3339) }

// SavingsGoalDeletedEvent is published when a savings goal is deleted.
type SavingsGoalDeletedEvent struct {
	GoalID    string
	User      string
	Timestamp time.Time
}

func (e SavingsGoalDeletedEvent) EventType() string   { return "savings_goal.deleted" }
func (e SavingsGoalDeletedEvent) AggregateID() string { return e.GoalID }
func (e SavingsGoalDeletedEvent) UserID() string      { return e.User }
func (e SavingsGoalDeletedEvent) OccurredAt() string  { return e.Timestamp.Format(time.RFC3339) }
