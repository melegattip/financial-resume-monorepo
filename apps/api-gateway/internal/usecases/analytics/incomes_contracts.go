package analytics

// IncomesSummaryRequest representa los parámetros para el resumen de ingresos
type IncomesSummaryRequest struct {
	UserID string `json:"user_id"`
	Year   string `json:"year,omitempty"`
	Month  string `json:"month,omitempty"`
	SortBy string `json:"sort_by,omitempty"` // date, amount, category
	Order  string `json:"order,omitempty"`   // asc, desc
}

// IncomesSummaryResponse representa la respuesta del resumen de ingresos
type IncomesSummaryResponse struct {
	Incomes []IncomeAnalyticsItem `json:"incomes"`
	Summary IncomeSummaryInfo     `json:"summary"`
}

// IncomeAnalyticsItem representa un ingreso con datos analíticos
type IncomeAnalyticsItem struct {
	ID           string  `json:"id"`
	Description  string  `json:"description"`
	Amount       float64 `json:"amount"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	CreatedAt    string  `json:"created_at"`
}

// IncomeSummaryInfo contiene métricas agregadas de ingresos
type IncomeSummaryInfo struct {
	TotalAmount        float64 `json:"total_amount"`
	AverageTransaction float64 `json:"average_transaction"`
	TransactionCount   int     `json:"transaction_count"`
}
