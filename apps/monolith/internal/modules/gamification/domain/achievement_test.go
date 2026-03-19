package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultAchievements_Count(t *testing.T) {
	achs := DefaultAchievements("user-1")
	// 8 original + 6 new flywheel + 1 education achievement = 15 total.
	assert.Len(t, achs, 15)
}

func TestDefaultAchievements_NewFlywheelTypesPresent(t *testing.T) {
	achs := DefaultAchievements("user-1")

	types := make(map[string]Achievement, len(achs))
	for _, a := range achs {
		types[a.Type] = a
	}

	newTypes := []string{
		"savings_starter",
		"savings_champion",
		"planner_pro",
		"budget_beginner",
		"budget_disciplined",
		"ai_executor",
	}
	for _, typ := range newTypes {
		a, ok := types[typ]
		require.True(t, ok, "achievement type %q not found in DefaultAchievements", typ)
		assert.False(t, a.Completed, "new achievement %q should start incomplete", typ)
		assert.Equal(t, 0, a.Progress, "new achievement %q should start at progress 0", typ)
		assert.Greater(t, a.Target, 0, "achievement %q must have a positive target", typ)
		assert.Greater(t, a.Points, 0, "achievement %q must have positive points", typ)
		assert.Equal(t, "user-1", a.UserID)
		assert.NotEmpty(t, a.ID)
	}
}

func TestDefaultAchievements_FlywheelTargets(t *testing.T) {
	achs := DefaultAchievements("u")

	types := make(map[string]Achievement, len(achs))
	for _, a := range achs {
		types[a.Type] = a
	}

	assert.Equal(t, 1, types["savings_starter"].Target)
	assert.Equal(t, 1, types["savings_champion"].Target)
	assert.Equal(t, 3, types["planner_pro"].Target)
	assert.Equal(t, 1, types["budget_beginner"].Target)
	assert.Equal(t, 3, types["budget_disciplined"].Target)
	assert.Equal(t, 5, types["ai_executor"].Target)
}

func TestDefaultAchievements_UniqueIDs(t *testing.T) {
	achs := DefaultAchievements("user-2")

	ids := make(map[string]struct{}, len(achs))
	for _, a := range achs {
		_, dup := ids[a.ID]
		assert.False(t, dup, "duplicate achievement ID %q", a.ID)
		ids[a.ID] = struct{}{}
	}
}

func TestDefaultAchievements_DifferentUsersGetDistinctIDs(t *testing.T) {
	achs1 := DefaultAchievements("user-A")
	achs2 := DefaultAchievements("user-B")

	ids1 := make(map[string]struct{}, len(achs1))
	for _, a := range achs1 {
		ids1[a.ID] = struct{}{}
	}
	for _, a := range achs2 {
		_, collision := ids1[a.ID]
		assert.False(t, collision, "IDs should not collide across users")
	}
}

func TestAchievement_UpdateProgress_Completion(t *testing.T) {
	a := Achievement{Target: 3, Progress: 0, Completed: false}
	a.UpdateProgress(3)

	assert.True(t, a.Completed)
	assert.NotNil(t, a.UnlockedAt)
}

func TestAchievement_UpdateProgress_NotYetComplete(t *testing.T) {
	a := Achievement{Target: 5, Progress: 0, Completed: false}
	a.UpdateProgress(4)

	assert.False(t, a.Completed)
	assert.Nil(t, a.UnlockedAt)
}

func TestAchievement_UpdateProgress_AlreadyCompletedNotOverwritten(t *testing.T) {
	a := Achievement{Target: 1, Progress: 1, Completed: true}
	first := a.UnlockedAt
	a.UpdateProgress(2)

	// UnlockedAt must not be updated a second time.
	assert.Equal(t, first, a.UnlockedAt)
}

func TestAchievement_IsCompleted(t *testing.T) {
	a := Achievement{Target: 5, Progress: 5}
	assert.True(t, a.IsCompleted())

	a.Progress = 4
	assert.False(t, a.IsCompleted())
}

func TestDefaultAchievements_ContainsFinancialLearner(t *testing.T) {
	achs := DefaultAchievements("user_test")
	var found *Achievement
	for i := range achs {
		if achs[i].Type == "financial_learner" {
			found = &achs[i]
			break
		}
	}
	require.NotNil(t, found, "financial_learner achievement should be in DefaultAchievements")
	assert.Equal(t, 3, found.Target)
	assert.Equal(t, 50, found.Points)
	assert.False(t, found.Completed)
	assert.Equal(t, "user_test", found.UserID)
}
