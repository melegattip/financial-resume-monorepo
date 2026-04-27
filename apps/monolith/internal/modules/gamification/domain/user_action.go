package domain

import "time"

// Action type constants define all trackable user actions in the gamification system.
const (
	ActionCreateExpense       = "create_expense"
	ActionCreateIncome        = "create_income"
	ActionUpdateExpense       = "update_expense"
	ActionUpdateIncome        = "update_income"
	ActionDeleteExpense       = "delete_expense"
	ActionDeleteIncome        = "delete_income"
	ActionCreateCategory      = "create_category"
	ActionUpdateCategory      = "update_category"
	ActionAssignCategory      = "assign_category"
	ActionDailyLogin          = "daily_login"
	ActionWeeklyStreak        = "weekly_streak"
	ActionMonthlyStreak       = "monthly_streak"
	ActionCompleteProfile     = "complete_profile"
	ActionViewDashboard       = "view_dashboard"
	ActionViewExpenses        = "view_expenses"
	ActionViewIncomes         = "view_incomes"
	ActionViewCategories      = "view_categories"
	ActionViewAnalytics       = "view_analytics"
	ActionViewMonthlyReport   = "view_monthly_report"
	ActionViewCategoryBreakdown = "view_category_breakdown"
	ActionExportData          = "export_data"
	ActionCreateSavingsGoal   = "create_savings_goal"
	ActionDepositSavings      = "deposit_savings"
	ActionAchieveSavingsGoal  = "achieve_savings_goal"
	ActionCreateBudget               = "create_budget"
	ActionStayWithinBudget           = "stay_within_budget"
	ActionCreateRecurringTransaction = "create_recurring_transaction"
	ActionApplyAIRecommendation      = "apply_ai_recommendation"
	ActionCompleteMonthlyReview      = "complete_monthly_review"
	ActionReadEducationCard          = "read_education_card"
)

// UserAction records a single action performed by a user.
type UserAction struct {
	ID          string
	UserID      string
	ActionType  string
	EntityType  string
	EntityID    string
	Description string
	XPEarned    int
	CreatedAt   time.Time
}

// XPForAction returns the point reward for a given action type.
// Values are calibrated for the 0–1000 scoring scale (max level at 1000 pts).
func XPForAction(actionType string) int {
	xpTable := map[string]int{
		ActionViewDashboard:              1,
		ActionViewExpenses:               1,
		ActionViewIncomes:                1,
		ActionViewCategories:             1,
		ActionViewAnalytics:              1,
		ActionCreateExpense:              2,
		ActionCreateIncome:               2,
		ActionUpdateExpense:              1,
		ActionUpdateIncome:               1,
		ActionDeleteExpense:              1,
		ActionDeleteIncome:               1,
		ActionCreateCategory:             2,
		ActionUpdateCategory:             1,
		ActionAssignCategory:             1,
		ActionDailyLogin:                 1,
		ActionWeeklyStreak:               5,
		ActionMonthlyStreak:              18,
		ActionCompleteProfile:            9,
		ActionViewMonthlyReport:          1,
		ActionViewCategoryBreakdown:      1,
		ActionExportData:                 2,
		ActionCreateSavingsGoal:          3,
		ActionDepositSavings:             1,
		ActionAchieveSavingsGoal:         18,
		ActionCreateBudget:               4,
		ActionStayWithinBudget:           3,
		ActionCreateRecurringTransaction: 5,
		ActionApplyAIRecommendation:      4,
		ActionCompleteMonthlyReview:      3,
		ActionReadEducationCard:          1,
	}
	if xp, ok := xpTable[actionType]; ok {
		return xp
	}
	return 1
}
