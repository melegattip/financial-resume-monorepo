package analytics

// ExpensesSummaryRequest representa los parámetros para el resumen de gastos
type ExpensesSummaryRequest struct {
	UserID string `json:"user_id"`
	Year   string `json:"year,omitempty"`
	Month  string `json:"month,omitempty"`
	SortBy string `json:"sort_by,omitempty"` // date, amount, category
	Order  string `json:"order,omitempty"`   // asc, desc
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

// ExpensesSummaryResponse representa la respuesta del resumen de gastos
type ExpensesSummaryResponse struct {
	Expenses   []ExpenseAnalyticsItem `json:"expenses"`
	Summary    ExpenseSummaryInfo     `json:"summary"`
	Pagination PaginationInfo         `json:"pagination"`
}

// ExpenseAnalyticsItem representa un gasto con datos analíticos
type ExpenseAnalyticsItem struct {
	ID                 string  `json:"id"`
	Description        string  `json:"description"`
	Amount             float64 `json:"amount"`
	AmountPaid         float64 `json:"amount_paid"`
	PendingAmount      float64 `json:"pending_amount"`
	PercentageOfIncome float64 `json:"percentage_of_income"`
	CategoryID         string  `json:"category_id"`
	CategoryName       string  `json:"category_name"`
	Paid               bool    `json:"paid"`
	DueDate            string  `json:"due_date"`
	CreatedAt          string  `json:"created_at"`
	DaysUntilDue       *int    `json:"days_until_due"`
}

// ExpenseSummaryInfo contiene métricas agregadas de gastos
type ExpenseSummaryInfo struct {
	TotalAmount             float64 `json:"total_amount"`
	PaidAmount              float64 `json:"paid_amount"`
	PendingAmount           float64 `json:"pending_amount"`
	AverageTransaction      float64 `json:"average_transaction"`
	PercentageOfTotalIncome float64 `json:"percentage_of_total_income"`
}

// PaginationInfo contiene información de paginación
type PaginationInfo struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}
