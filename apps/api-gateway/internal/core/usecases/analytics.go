package usecases

import "context"

// ExpensesAnalyticsUseCase define el caso de uso para analytics de gastos
type ExpensesAnalyticsUseCase interface {
	GetExpensesSummary(ctx context.Context, params ExpensesSummaryParams) (*ExpensesSummary, error)
}

// CategoriesAnalyticsUseCase define el caso de uso para analytics de categorías
type CategoriesAnalyticsUseCase interface {
	GetCategoriesAnalytics(ctx context.Context, params CategoriesAnalyticsParams) (*CategoriesAnalytics, error)
}

// IncomesAnalyticsUseCase define el caso de uso para analytics de ingresos
type IncomesAnalyticsUseCase interface {
	GetIncomesSummary(ctx context.Context, params IncomesSummaryParams) (*IncomesSummary, error)
}

// ExpensesSummaryParams representa los parámetros para resumen de gastos
type ExpensesSummaryParams struct {
	UserID     string
	Period     DatePeriod
	Sorting    SortingCriteria
	Pagination PaginationParams
}

// CategoriesAnalyticsParams representa los parámetros para analytics de categorías
type CategoriesAnalyticsParams struct {
	UserID string
	Period DatePeriod
}

// IncomesSummaryParams representa los parámetros para resumen de ingresos
type IncomesSummaryParams struct {
	UserID  string
	Period  DatePeriod
	Sorting SortingCriteria
}

// SortingCriteria define criterios de ordenamiento
type SortingCriteria struct {
	Field SortField
	Order SortOrder
}

// SortField representa los campos por los que se puede ordenar
type SortField string

const (
	SortByDate     SortField = "date"
	SortByAmount   SortField = "amount"
	SortByCategory SortField = "category"
)

// SortOrder representa el orden de clasificación
type SortOrder string

const (
	Ascending  SortOrder = "asc"
	Descending SortOrder = "desc"
)

// PaginationParams define parámetros de paginación
type PaginationParams struct {
	Limit  int
	Offset int
}

// ExpensesSummary representa el resumen de gastos
type ExpensesSummary struct {
	Expenses   []ExpenseItem
	Summary    ExpenseSummary
	Pagination PaginationInfo
}

// CategoriesAnalytics representa el análisis de categorías
type CategoriesAnalytics struct {
	Categories []CategoryItem
	Summary    CategorySummary
}

// IncomesSummary representa el resumen de ingresos
type IncomesSummary struct {
	Incomes []IncomeItem
	Summary IncomeSummary
}

// ExpenseItem representa un item de gasto en analytics
type ExpenseItem struct {
	ID                 string
	Description        string
	Amount             float64
	AmountPaid         float64
	PendingAmount      float64
	PercentageOfIncome float64
	CategoryID         string
	CategoryName       string
	Paid               bool
	DueDate            string
	CreatedAt          string
	DaysUntilDue       *int
}

// CategoryItem representa un item de categoría en analytics
type CategoryItem struct {
	CategoryID            string
	CategoryName          string
	TotalAmount           float64
	PercentageOfExpenses  float64
	PercentageOfIncome    float64
	TransactionCount      int
	AveragePerTransaction float64
	ColorSeed             int
}

// IncomeItem representa un item de ingreso en analytics
type IncomeItem struct {
	ID           string
	Description  string
	Amount       float64
	CategoryID   string
	CategoryName string
	CreatedAt    string
}

// ExpenseSummary contiene resumen agregado de gastos
type ExpenseSummary struct {
	TotalAmount             float64
	PaidAmount              float64
	PendingAmount           float64
	AverageTransaction      float64
	PercentageOfTotalIncome float64
}

// CategorySummary contiene resumen de categorías
type CategorySummary struct {
	TotalCategories  int
	LargestCategory  string
	SmallestCategory string
	TotalAmount      float64
}

// IncomeSummary contiene resumen de ingresos
type IncomeSummary struct {
	TotalAmount        float64
	AverageTransaction float64
	TransactionCount   int
}

// PaginationInfo contiene información de paginación
type PaginationInfo struct {
	Total   int
	Limit   int
	Offset  int
	HasMore bool
}

// AnalyticsCalculator define la interfaz para cálculos analíticos
type AnalyticsCalculator interface {
	CalculatePercentages(amount, total float64) float64
	CalculateAverage(total float64, count int) float64
	GenerateColorSeed(identifier string) int
}

// CategoryService define la interfaz para servicios de categorías
type CategoryService interface {
	GetCategoryName(categoryID string) (string, error)
	GetCategoryNames(categoryIDs []string) (map[string]string, error)
}
