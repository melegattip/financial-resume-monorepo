package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXPForAction_KnownActions(t *testing.T) {
	// Expected values are calibrated for the 0–1000 scoring scale.
	cases := []struct {
		action string
		xp     int
	}{
		{ActionCreateExpense, 2},
		{ActionCreateIncome, 2},
		{ActionUpdateExpense, 1},
		{ActionUpdateIncome, 1},
		{ActionDeleteExpense, 1},
		{ActionDeleteIncome, 1},
		{ActionCreateCategory, 2},
		{ActionUpdateCategory, 1},
		{ActionAssignCategory, 1},
		{ActionDailyLogin, 1},
		{ActionWeeklyStreak, 5},
		{ActionMonthlyStreak, 18},
		{ActionCompleteProfile, 9},
		{ActionViewDashboard, 1},
		{ActionViewExpenses, 1},
		{ActionViewIncomes, 1},
		{ActionViewCategories, 1},
		{ActionViewAnalytics, 1},
		{ActionViewMonthlyReport, 1},
		{ActionViewCategoryBreakdown, 1},
		{ActionExportData, 2},
		// Quality financial actions.
		{ActionCreateSavingsGoal, 3},
		{ActionDepositSavings, 1},
		{ActionAchieveSavingsGoal, 18},
		{ActionCreateBudget, 4},
		{ActionStayWithinBudget, 3},
		// Flywheel action types.
		{ActionCreateRecurringTransaction, 5},
		{ActionApplyAIRecommendation, 4},
		{ActionCompleteMonthlyReview, 3},
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

func TestXPForAction_ReadEducationCard(t *testing.T) {
	assert.Equal(t, 1, XPForAction(ActionReadEducationCard))
}

func TestActionReadEducationCard_Constant(t *testing.T) {
	assert.Equal(t, "read_education_card", ActionReadEducationCard)
}

func TestXPForAction_ReadEducationCard_LessThanCreateBudget(t *testing.T) {
	// Education card reading gives less XP than creating a budget (exploratory vs active).
	assert.Less(t, XPForAction(ActionReadEducationCard), XPForAction(ActionCreateBudget))
}
