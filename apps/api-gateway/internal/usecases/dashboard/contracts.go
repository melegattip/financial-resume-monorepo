package dashboard

// DashboardRequest representa los parámetros de consulta para el dashboard
type DashboardRequest struct {
	UserID string `json:"user_id"`
	Year   string `json:"year,omitempty"`
	Month  string `json:"month,omitempty"`
}

// DashboardResponse representa la respuesta completa del dashboard
type DashboardResponse struct {
	Period  PeriodInfo  `json:"period"`
	Metrics MetricsInfo `json:"metrics"`
	Counts  CountsInfo  `json:"counts"`
	Trends  TrendsInfo  `json:"trends"`
}

// PeriodInfo contiene información sobre el período consultado
type PeriodInfo struct {
	Year    string `json:"year"`
	Month   string `json:"month"`
	Label   string `json:"label"`
	HasData bool   `json:"has_data"`
}

// MetricsInfo contiene las métricas financieras principales
type MetricsInfo struct {
	TotalIncome          float64 `json:"total_income"`
	TotalExpenses        float64 `json:"total_expenses"`
	Balance              float64 `json:"balance"`
	PendingExpenses      float64 `json:"pending_expenses"`
	PendingExpensesCount int     `json:"pending_expenses_count"`
}

// CountsInfo contiene contadores de transacciones
type CountsInfo struct {
	IncomeTransactions  int `json:"income_transactions"`
	ExpenseTransactions int `json:"expense_transactions"`
	CategoriesUsed      int `json:"categories_used"`
}

// TrendsInfo contiene información de tendencias
type TrendsInfo struct {
	ExpensePercentageOfIncome float64 `json:"expense_percentage_of_income"`
	MonthOverMonthChange      float64 `json:"month_over_month_change"`
}
