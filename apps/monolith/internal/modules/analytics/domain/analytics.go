package domain

import "time"

// ExpenseSummary is the aggregated expense summary for a period.
type ExpenseSummary struct {
	TotalAmount   float64           `json:"total_amount"`
	Count         int               `json:"count"`
	AverageAmount float64           `json:"average_amount"`
	ByCategory    []CategorySummary `json:"by_category"`
	ByMonth       []MonthlySummary  `json:"by_month"`
	Period        string            `json:"period"` // "this_month", "last_month", "this_year", etc.
}

// CategorySummary holds aggregated data per category.
type CategorySummary struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
}

// MonthlySummary holds aggregated data per calendar month.
type MonthlySummary struct {
	Year   int     `json:"year"`
	Month  int     `json:"month"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}

// IncomeSummary is the aggregated income summary for a period.
type IncomeSummary struct {
	TotalAmount   float64          `json:"total_amount"`
	Count         int              `json:"count"`
	AverageAmount float64          `json:"average_amount"`
	ByMonth       []MonthlySummary `json:"by_month"`
	Period        string           `json:"period"`
}

// DashboardSummary is the top-level view shown on the user's dashboard.
type DashboardSummary struct {
	CurrentMonthExpenses float64          `json:"current_month_expenses"`
	CurrentMonthIncomes  float64          `json:"current_month_incomes"`
	CurrentMonthBalance  float64          `json:"current_month_balance"`
	TotalExpenses        float64          `json:"total_expenses"`
	TotalIncomes         float64          `json:"total_incomes"`
	SavingsRate          float64          `json:"savings_rate"` // (incomes - expenses) / incomes
	TopCategories        []CategorySummary `json:"top_categories"`
	RecentExpenses       []RecentItem     `json:"recent_expenses"`
	RecentIncomes        []RecentItem     `json:"recent_incomes"`
	UpdatedAt            time.Time        `json:"updated_at"`
}

// RecentItem is a lightweight representation of a recent transaction.
type RecentItem struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category,omitempty"`
	Date        time.Time `json:"date"`
}
