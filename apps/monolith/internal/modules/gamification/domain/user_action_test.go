package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXPForAction_KnownActions(t *testing.T) {
	cases := []struct {
		action string
		xp     int
	}{
		{ActionCreateExpense, 8},
		{ActionCreateIncome, 8},
		{ActionUpdateExpense, 5},
		{ActionUpdateIncome, 5},
		{ActionDeleteExpense, 3},
		{ActionDeleteIncome, 3},
		{ActionCreateCategory, 10},
		{ActionUpdateCategory, 5},
		{ActionAssignCategory, 3},
		{ActionDailyLogin, 5},
		{ActionWeeklyStreak, 25},
		{ActionMonthlyStreak, 100},
		{ActionCompleteProfile, 50},
		{ActionViewDashboard, 2},
		{ActionViewExpenses, 1},
		{ActionViewIncomes, 1},
		{ActionViewCategories, 1},
		{ActionViewAnalytics, 3},
		{ActionViewMonthlyReport, 5},
		{ActionViewCategoryBreakdown, 3},
		{ActionExportData, 10},
		// Pre-existing quality actions.
		{ActionCreateSavingsGoal, 15},
		{ActionDepositSavings, 8},
		{ActionAchieveSavingsGoal, 100},
		{ActionCreateBudget, 20},
		{ActionStayWithinBudget, 15},
		// New flywheel action types.
		{ActionCreateRecurringTransaction, 30},
		{ActionApplyAIRecommendation, 20},
		{ActionCompleteMonthlyReview, 15},
	}

	for _, tc := range cases {
		got := XPForAction(tc.action)
		assert.Equal(t, tc.xp, got, "unexpected XP for action %q", tc.action)
	}
}

func TestXPForAction_UnknownReturnsOne(t *testing.T) {
	assert.Equal(t, 1, XPForAction("unknown_action_type"))
	assert.Equal(t, 1, XPForAction(""))
}

func TestXPForAction_NewFlywheelActions_HigherThanViewActions(t *testing.T) {
	// Quality financial actions should reward more XP than passive view actions.
	assert.Greater(t, XPForAction(ActionCreateRecurringTransaction), XPForAction(ActionViewDashboard))
	assert.Greater(t, XPForAction(ActionApplyAIRecommendation), XPForAction(ActionViewAnalytics))
	assert.Greater(t, XPForAction(ActionCompleteMonthlyReview), XPForAction(ActionViewExpenses))
}

func TestXPForAction_AchieveSavingsGoal_HighestReward(t *testing.T) {
	// Completing a savings goal is the highest XP action along with monthly streak.
	assert.Equal(t, XPForAction(ActionAchieveSavingsGoal), XPForAction(ActionMonthlyStreak))
	assert.Greater(t, XPForAction(ActionAchieveSavingsGoal), XPForAction(ActionCreateRecurringTransaction))
}
