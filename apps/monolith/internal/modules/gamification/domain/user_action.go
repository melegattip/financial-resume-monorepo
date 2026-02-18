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
	ActionCreateBudget        = "create_budget"
	ActionStayWithinBudget    = "stay_within_budget"
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

// XPForAction returns the XP reward for a given action type.
func XPForAction(actionType string) int {
	xpTable := map[string]int{
		ActionViewDashboard:       2,
		ActionViewExpenses:        1,
		ActionViewIncomes:         1,
		ActionViewCategories:      1,
		ActionViewAnalytics:       3,
		ActionCreateExpense:       8,
		ActionCreateIncome:        8,
		ActionUpdateExpense:       5,
		ActionUpdateIncome:        5,
		ActionDeleteExpense:       3,
		ActionDeleteIncome:        3,
		ActionCreateCategory:      10,
		ActionUpdateCategory:      5,
		ActionAssignCategory:      3,
		ActionDailyLogin:          5,
		ActionWeeklyStreak:        25,
		ActionMonthlyStreak:       100,
		ActionCompleteProfile:     50,
		ActionViewMonthlyReport:   5,
		ActionViewCategoryBreakdown: 3,
		ActionExportData:          10,
		ActionCreateSavingsGoal:   15,
		ActionDepositSavings:      8,
		ActionAchieveSavingsGoal:  100,
		ActionCreateBudget:        20,
		ActionStayWithinBudget:    15,
	}
	if xp, ok := xpTable[actionType]; ok {
		return xp
	}
	return 1
}
