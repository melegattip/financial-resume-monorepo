package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification/domain"
)

// ---------------------------------------------------------------------------
// countByType
// ---------------------------------------------------------------------------

func TestCountByType_Empty(t *testing.T) {
	counts := countByType(nil)
	assert.Empty(t, counts)
}

func TestCountByType_SingleAction(t *testing.T) {
	actions := []domain.UserAction{
		{ActionType: domain.ActionCreateBudget},
	}
	counts := countByType(actions)
	assert.Equal(t, 1, counts[domain.ActionCreateBudget])
}

func TestCountByType_MultipleActionsGrouped(t *testing.T) {
	actions := []domain.UserAction{
		{ActionType: domain.ActionCreateBudget},
		{ActionType: domain.ActionCreateBudget},
		{ActionType: domain.ActionCreateBudget},
		{ActionType: domain.ActionApplyAIRecommendation},
		{ActionType: domain.ActionApplyAIRecommendation},
		{ActionType: domain.ActionCreateRecurringTransaction},
	}
	counts := countByType(actions)

	assert.Equal(t, 3, counts[domain.ActionCreateBudget])
	assert.Equal(t, 2, counts[domain.ActionApplyAIRecommendation])
	assert.Equal(t, 1, counts[domain.ActionCreateRecurringTransaction])
	assert.Equal(t, 0, counts[domain.ActionCompleteMonthlyReview]) // absent key defaults to 0
}

func TestCountByType_AllFlywheelActions(t *testing.T) {
	actions := []domain.UserAction{
		{ActionType: domain.ActionCreateRecurringTransaction},
		{ActionType: domain.ActionApplyAIRecommendation},
		{ActionType: domain.ActionCompleteMonthlyReview},
		{ActionType: domain.ActionAchieveSavingsGoal},
		{ActionType: domain.ActionDepositSavings},
		{ActionType: domain.ActionStayWithinBudget},
	}
	counts := countByType(actions)

	assert.Equal(t, 1, counts[domain.ActionCreateRecurringTransaction])
	assert.Equal(t, 1, counts[domain.ActionApplyAIRecommendation])
	assert.Equal(t, 1, counts[domain.ActionCompleteMonthlyReview])
	assert.Equal(t, 1, counts[domain.ActionAchieveSavingsGoal])
	assert.Equal(t, 1, counts[domain.ActionDepositSavings])
	assert.Equal(t, 1, counts[domain.ActionStayWithinBudget])
}

func TestCountByType_TotalMatchesInputLength(t *testing.T) {
	actions := make([]domain.UserAction, 100)
	for i := range actions {
		actions[i] = domain.UserAction{ActionType: domain.ActionCreateExpense}
	}
	counts := countByType(actions)
	assert.Equal(t, 100, counts[domain.ActionCreateExpense])
}

// ---------------------------------------------------------------------------
// updateStreak (white-box — private method)
// ---------------------------------------------------------------------------

func newTestService() *GamificationService {
	return &GamificationService{}
}

func TestUpdateStreak_SameDay(t *testing.T) {
	svc := newTestService()
	now := time.Now().UTC()
	g := &domain.UserGamification{
		CurrentStreak: 5,
		LastActivity:  now,
	}
	svc.updateStreak(g)
	assert.Equal(t, 5, g.CurrentStreak, "same-day login must not change streak")
}

func TestUpdateStreak_ConsecutiveDay(t *testing.T) {
	svc := newTestService()
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	g := &domain.UserGamification{
		CurrentStreak: 5,
		LastActivity:  yesterday,
	}
	svc.updateStreak(g)
	assert.Equal(t, 6, g.CurrentStreak)
}

func TestUpdateStreak_GracePeriod(t *testing.T) {
	// Missed exactly one day — streak should still extend (grace period).
	svc := newTestService()
	twoDaysAgo := time.Now().UTC().AddDate(0, 0, -2)
	g := &domain.UserGamification{
		CurrentStreak: 5,
		LastActivity:  twoDaysAgo,
	}
	svc.updateStreak(g)
	assert.Equal(t, 6, g.CurrentStreak)
}

func TestUpdateStreak_Reset(t *testing.T) {
	svc := newTestService()
	threeDaysAgo := time.Now().UTC().AddDate(0, 0, -3)
	g := &domain.UserGamification{
		CurrentStreak: 20,
		LastActivity:  threeDaysAgo,
	}
	svc.updateStreak(g)
	assert.Equal(t, 1, g.CurrentStreak)
}
