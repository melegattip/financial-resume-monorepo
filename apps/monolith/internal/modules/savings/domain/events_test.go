package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// SavingsGoalAchievedEvent — interface contract
// ---------------------------------------------------------------------------

func TestSavingsGoalAchievedEvent_EventType(t *testing.T) {
	e := SavingsGoalAchievedEvent{GoalID: "g1", User: "u1"}
	assert.Equal(t, "savings_goal.achieved", e.EventType())
}

func TestSavingsGoalAchievedEvent_AggregateID(t *testing.T) {
	e := SavingsGoalAchievedEvent{GoalID: "goal-123", User: "user-1"}
	assert.Equal(t, "goal-123", e.AggregateID())
}

func TestSavingsGoalAchievedEvent_UserID(t *testing.T) {
	e := SavingsGoalAchievedEvent{GoalID: "g1", User: "user-abc"}
	assert.Equal(t, "user-abc", e.UserID())
}

func TestSavingsGoalAchievedEvent_OccurredAt_RFC3339(t *testing.T) {
	ts := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)
	e := SavingsGoalAchievedEvent{Timestamp: ts}
	assert.Equal(t, "2026-03-12T10:00:00Z", e.OccurredAt())
}

func TestSavingsGoalAchievedEvent_DistinctFromCreated(t *testing.T) {
	achieved := SavingsGoalAchievedEvent{}
	created := SavingsGoalCreatedEvent{}
	assert.NotEqual(t, achieved.EventType(), created.EventType())
}

// ---------------------------------------------------------------------------
// Existing events — regression checks
// ---------------------------------------------------------------------------

func TestSavingsGoalCreatedEvent_EventType(t *testing.T) {
	e := SavingsGoalCreatedEvent{}
	assert.Equal(t, "savings_goal.created", e.EventType())
}

func TestSavingsGoalUpdatedEvent_EventType(t *testing.T) {
	e := SavingsGoalUpdatedEvent{}
	assert.Equal(t, "savings_goal.updated", e.EventType())
}

func TestSavingsGoalDeletedEvent_EventType(t *testing.T) {
	e := SavingsGoalDeletedEvent{}
	assert.Equal(t, "savings_goal.deleted", e.EventType())
}

func TestAllSavingsEvents_UniqueEventTypes(t *testing.T) {
	types := []string{
		SavingsGoalCreatedEvent{}.EventType(),
		SavingsGoalUpdatedEvent{}.EventType(),
		SavingsGoalAchievedEvent{}.EventType(),
		SavingsGoalDeletedEvent{}.EventType(),
	}
	seen := make(map[string]struct{})
	for _, et := range types {
		_, dup := seen[et]
		assert.False(t, dup, "duplicate event type: %q", et)
		seen[et] = struct{}{}
	}
}
