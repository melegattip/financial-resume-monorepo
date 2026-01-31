package analytics

// CategoriesAnalyticsRequest representa los parámetros para analytics de categorías
type CategoriesAnalyticsRequest struct {
	UserID string `json:"user_id"`
	Year   string `json:"year,omitempty"`
	Month  string `json:"month,omitempty"`
}

// CategoriesAnalyticsResponse representa la respuesta de analytics de categorías
type CategoriesAnalyticsResponse struct {
	Categories []CategoryAnalyticsItem `json:"categories"`
	Summary    CategorySummaryInfo     `json:"summary"`
}

// CategoryAnalyticsItem representa una categoría con datos analíticos
type CategoryAnalyticsItem struct {
	CategoryID            string  `json:"category_id"`
	CategoryName          string  `json:"category_name"`
	TotalAmount           float64 `json:"total_amount"`
	PercentageOfExpenses  float64 `json:"percentage_of_expenses"`
	PercentageOfIncome    float64 `json:"percentage_of_income"`
	TransactionCount      int     `json:"transaction_count"`
	AveragePerTransaction float64 `json:"average_per_transaction"`
	ColorSeed             int     `json:"color_seed"`
}

// CategorySummaryInfo contiene resumen de categorías
type CategorySummaryInfo struct {
	TotalCategories  int     `json:"total_categories"`
	LargestCategory  string  `json:"largest_category"`
	SmallestCategory string  `json:"smallest_category"`
	TotalAmount      float64 `json:"total_amount"`
}
