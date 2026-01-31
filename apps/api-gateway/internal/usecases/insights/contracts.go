package insights

import (
	"context"
	"time"
)

// InsightsUseCase define el caso de uso principal de insights financieros
type InsightsUseCase interface {
	GetFinancialHealth(ctx context.Context, params InsightsParams) (*FinancialHealthResponse, error)
	GetAIInsights(ctx context.Context, params InsightsParams) (*AIInsightsResponse, error)
	CanIBuy(ctx context.Context, params CanIBuyParams) (*CanIBuyResponse, error)
	GetCreditImprovementPlan(ctx context.Context, params InsightsParams) (*CreditImprovementPlanResponse, error)
}

// InsightsParams representa los parámetros de entrada
type InsightsParams struct {
	UserID string
	Period DatePeriod
}

// DatePeriod representa un período de fechas (reutilizando la lógica existente)
type DatePeriod struct {
	Year  *int
	Month *int
}

// FinancialHealthResponse es la respuesta principal del endpoint
type FinancialHealthResponse struct {
	HealthScore  int                   `json:"health_score"`  // Score de 0-1000
	Level        string                `json:"level"`         // Excelente, Bueno, Regular, Mejorable
	Message      string                `json:"message"`       // Mensaje motivacional personalizado
	Insights     []FinancialInsight    `json:"insights"`      // Lista de insights/recomendaciones
	AnalyzedData AnalyzedFinancialData `json:"analyzed_data"` // Datos analizados para transparencia
	GeneratedAt  time.Time             `json:"generated_at"`  // Timestamp de generación
}

// FinancialInsight representa una recomendación/insight individual
type FinancialInsight struct {
	ID          string                 `json:"id"`          // Identificador único
	Title       string                 `json:"title"`       // Título del insight
	Description string                 `json:"description"` // Descripción detallada
	Type        InsightType            `json:"type"`        // Tipo de insight
	Impact      InsightImpact          `json:"impact"`      // Impacto/prioridad
	Score       int                    `json:"score"`       // Score individual (0-1000)
	Icon        string                 `json:"icon"`        // Emoji para mostrar
	Action      *RecommendedAction     `json:"action"`      // Acción recomendada (opcional)
	Metadata    map[string]interface{} `json:"metadata"`    // Datos adicionales
}

// InsightType define los tipos de insights
type InsightType string

const (
	InsightTypeSpending    InsightType = "spending"    // Relacionado con gastos
	InsightTypeSaving      InsightType = "saving"      // Relacionado con ahorros
	InsightTypeIncome      InsightType = "income"      // Relacionado con ingresos
	InsightTypeBudget      InsightType = "budget"      // Relacionado con presupuesto
	InsightTypePattern     InsightType = "pattern"     // Patrones de comportamiento
	InsightTypeOpportunity InsightType = "opportunity" // Oportunidades de mejora
)

// InsightImpact define el impacto/prioridad del insight
type InsightImpact string

const (
	InsightImpactHigh   InsightImpact = "high"   // Alto impacto (rojo/crítico)
	InsightImpactMedium InsightImpact = "medium" // Impacto medio (amarillo)
	InsightImpactLow    InsightImpact = "low"    // Bajo impacto (azul/informativo)
	InsightImpactGood   InsightImpact = "good"   // Positivo (verde)
)

// RecommendedAction representa una acción recomendada
type RecommendedAction struct {
	Title       string `json:"title"`       // Título de la acción
	Description string `json:"description"` // Descripción de qué hacer
	Difficulty  string `json:"difficulty"`  // Fácil, Medio, Difícil
	XPReward    int    `json:"xp_reward"`   // Puntos XP por completar
}

// AnalyzedFinancialData contiene los datos analizados
type AnalyzedFinancialData struct {
	Period            PeriodInfo               `json:"period"`
	TotalIncome       float64                  `json:"total_income"`
	TotalExpenses     float64                  `json:"total_expenses"`
	Balance           float64                  `json:"balance"`
	SavingsRate       float64                  `json:"savings_rate"`       // % de ingresos ahorrados
	ExpenseCategories []CategoryAnalysis       `json:"expense_categories"` // Análisis por categoría
	SpendingPatterns  SpendingPatternAnalysis  `json:"spending_patterns"`  // Patrones de gasto
	IncomeStability   IncomeStabilityAnalysis  `json:"income_stability"`   // Análisis de estabilidad de ingresos
	BudgetCompliance  BudgetComplianceAnalysis `json:"budget_compliance"`  // Cumplimiento de presupuesto
}

// PeriodInfo contiene información del período analizado
type PeriodInfo struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	Label        string `json:"label"`
	DaysInPeriod int    `json:"days_in_period"`
}

// CategoryAnalysis representa el análisis de una categoría de gastos
type CategoryAnalysis struct {
	CategoryID       string  `json:"category_id"`
	CategoryName     string  `json:"category_name"`
	Amount           float64 `json:"amount"`
	Percentage       float64 `json:"percentage"` // % del total de gastos
	TransactionCount int     `json:"transaction_count"`
	AverageAmount    float64 `json:"average_amount"`
	IsRecurring      bool    `json:"is_recurring"` // Si tiene gastos recurrentes
}

// SpendingPatternAnalysis analiza patrones de gasto
type SpendingPatternAnalysis struct {
	AverageTransactionAmount float64           `json:"average_transaction_amount"`
	LargestExpense           float64           `json:"largest_expense"`
	SmallestExpense          float64           `json:"smallest_expense"`
	FrequentCategories       []string          `json:"frequent_categories"` // Categorías más usadas
	UnusualSpending          []UnusualSpending `json:"unusual_spending"`    // Gastos inusuales
	DailyAverageSpending     float64           `json:"daily_average_spending"`
}

// UnusualSpending representa un gasto inusual detectado
type UnusualSpending struct {
	Amount     float64   `json:"amount"`
	CategoryID string    `json:"category_id"`
	Date       time.Time `json:"date"`
	Reason     string    `json:"reason"` // Por qué se considera inusual
}

// IncomeStabilityAnalysis analiza la estabilidad de ingresos
type IncomeStabilityAnalysis struct {
	IsStable             bool    `json:"is_stable"`
	AverageMonthlyIncome float64 `json:"average_monthly_income"`
	IncomeVariation      float64 `json:"income_variation"`       // Coeficiente de variación
	RecurringIncomeRatio float64 `json:"recurring_income_ratio"` // % de ingresos recurrentes
}

// BudgetComplianceAnalysis analiza el cumplimiento del presupuesto
type BudgetComplianceAnalysis struct {
	HasBudget           bool     `json:"has_budget"`
	BudgetCompliance    float64  `json:"budget_compliance"`    // % de cumplimiento (si tiene presupuesto)
	OverspentCategories []string `json:"overspent_categories"` // Categorías que excedieron presupuesto
}

// === NUEVOS ENDPOINTS DE IA ===

// AIInsightsResponse representa insights generados por IA
type AIInsightsResponse struct {
	Insights    []AIInsight `json:"insights"`
	GeneratedAt time.Time   `json:"generated_at"`
	Source      string      `json:"source"` // "ai" o "mock"
}

// AIInsight representa un insight generado por IA
type AIInsight struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`      // "high", "medium", "low"
	Score       int    `json:"score"`       // 0-1000
	ActionType  string `json:"action_type"` // "save", "optimize", "alert", "invest"
	Category    string `json:"category"`
}

// SavingsGoalInfo representa información de una meta de ahorro relevante
type SavingsGoalInfo struct {
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	CurrentAmount float64 `json:"current_amount"`
	TargetAmount  float64 `json:"target_amount"`
	Progress      float64 `json:"progress"`
}

// CanIBuyParams representa los parámetros para análisis de compra (LEGACY - mantener compatibilidad)
type CanIBuyParams struct {
	UserID          string            `json:"user_id"`
	ItemName        string            `json:"item_name"`
	Amount          float64           `json:"amount"`
	Description     string            `json:"description"`   // Descripción opcional
	PaymentTypes    []string          `json:"payment_types"` // Array de tipos de pago ["contado", "cuotas", "ahorro"]
	IsNecessary     bool              `json:"is_necessary"`  // Si es una necesidad urgente
	CurrentBalance  float64           `json:"current_balance"`
	MonthlyIncome   float64           `json:"monthly_income"`
	MonthlyExpenses float64           `json:"monthly_expenses"`
	SavingsGoal     float64           `json:"savings_goal"`
	SavingsGoals    []SavingsGoalInfo `json:"savings_goals"` // Metas de ahorro existentes
}

// CanIBuyResponse representa la respuesta de decisión de compra
type CanIBuyResponse struct {
	CanBuy       bool      `json:"can_buy"`
	Confidence   float64   `json:"confidence"` // 0.0-1.0
	Reasoning    string    `json:"reasoning"`
	Alternatives []string  `json:"alternatives"`
	ImpactScore  int       `json:"impact_score"` // 0-1000
	GeneratedAt  time.Time `json:"generated_at"`
	Source       string    `json:"source"` // "ai" o "mock"
}

// CreditImprovementPlanResponse representa un plan de mejora crediticia
type CreditImprovementPlanResponse struct {
	CurrentScore   int                    `json:"current_score"`
	TargetScore    int                    `json:"target_score"`
	TimelineMonths int                    `json:"timeline_months"`
	Actions        []CreditAction         `json:"actions"`
	KeyMetrics     map[string]interface{} `json:"key_metrics"`
	GeneratedAt    time.Time              `json:"generated_at"`
	Source         string                 `json:"source"` // "ai" o "mock"
}

// CreditAction representa una acción específica para mejorar el crédito
type CreditAction struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`   // "high", "medium", "low"
	Timeline    string `json:"timeline"`   // "1-3 meses", etc.
	Impact      int    `json:"impact"`     // Puntos de mejora esperados
	Difficulty  string `json:"difficulty"` // "easy", "medium", "hard"
}
