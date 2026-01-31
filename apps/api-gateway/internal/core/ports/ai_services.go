package ports

import (
	"context"
	"time"
)

// ========== INTERFACES ESPECIALIZADAS PARA IA ==========
// Siguiendo el principio de Interface Segregation (ISP)

// AIAnalysisPort maneja el análisis financiero con IA
type AIAnalysisPort interface {
	AnalyzeFinancialHealth(ctx context.Context, data FinancialAnalysisData) (*HealthAnalysis, error)
	GenerateInsights(ctx context.Context, data FinancialAnalysisData) ([]AIInsight, error)
}

// PurchaseDecisionPort maneja las decisiones de compra inteligentes
type PurchaseDecisionPort interface {
	CanIBuy(ctx context.Context, request PurchaseAnalysisRequest) (*PurchaseDecision, error)
	SuggestAlternatives(ctx context.Context, request PurchaseAnalysisRequest) ([]Alternative, error)
}

// CreditAnalysisPort maneja el análisis y mejora crediticia
type CreditAnalysisPort interface {
	GenerateImprovementPlan(ctx context.Context, data FinancialAnalysisData) (*CreditPlan, error)
	CalculateCreditScore(ctx context.Context, data FinancialAnalysisData) (int, error)
}

// ========== DOMAIN ENTITIES PARA IA ==========

// FinancialAnalysisData representa los datos para análisis financiero
type FinancialAnalysisData struct {
	UserID             string             `json:"user_id"`
	TotalIncome        float64            `json:"total_income"`
	TotalExpenses      float64            `json:"total_expenses"`
	SavingsRate        float64            `json:"savings_rate"`
	ExpensesByCategory map[string]float64 `json:"expenses_by_category"`
	IncomeStability    float64            `json:"income_stability"`
	FinancialScore     int                `json:"financial_score"`
	Period             string             `json:"period"`
}

// AIInsight representa un insight generado por IA
type AIInsight struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Score       int    `json:"score"`
	ActionType  string `json:"action_type"`
	Category    string `json:"category"`
}

// HealthAnalysis representa el análisis de salud financiera
type HealthAnalysis struct {
	Score       int         `json:"score"`
	Level       string      `json:"level"`
	Message     string      `json:"message"`
	Insights    []AIInsight `json:"insights"`
	GeneratedAt time.Time   `json:"generated_at"`
}

// PurchaseAnalysisRequest representa una solicitud de análisis de compra
type PurchaseAnalysisRequest struct {
	UserID               string               `json:"user_id"`
	ItemName             string               `json:"item_name"`
	Amount               float64              `json:"amount"`
	Description          string               `json:"description,omitempty"`
	PaymentTypes         []string             `json:"payment_types,omitempty"`
	IsNecessary          bool                 `json:"is_necessary"`
	UserFinancialProfile UserFinancialProfile `json:"user_financial_profile"`
	SavingsGoal          *SavingsGoalInfo     `json:"savings_goal,omitempty"`
}

// PurchaseDecision representa la decisión sobre una compra
type PurchaseDecision struct {
	CanBuy       bool      `json:"can_buy"`
	Confidence   float64   `json:"confidence"`
	Reasoning    string    `json:"reasoning"`
	Alternatives []string  `json:"alternatives"`
	ImpactScore  int       `json:"impact_score"`
	GeneratedAt  time.Time `json:"generated_at"`
}

// Alternative representa una alternativa de compra
type Alternative struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Savings     float64 `json:"savings"`
	Feasibility string  `json:"feasibility"`
}

// CreditPlan representa un plan de mejora crediticia
type CreditPlan struct {
	CurrentScore   int                    `json:"current_score"`
	TargetScore    int                    `json:"target_score"`
	TimelineMonths int                    `json:"timeline_months"`
	Actions        []CreditAction         `json:"actions"`
	KeyMetrics     map[string]interface{} `json:"key_metrics"`
	GeneratedAt    time.Time              `json:"generated_at"`
}

// CreditAction representa una acción específica para mejorar el crédito
type CreditAction struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Timeline    string `json:"timeline"`
	Impact      int    `json:"impact"`
	Difficulty  string `json:"difficulty"`
}

// UserFinancialProfile representa el perfil financiero completo del usuario
type UserFinancialProfile struct {
	CurrentBalance        float64            `json:"current_balance"`
	MonthlyIncome         float64            `json:"monthly_income"`
	MonthlyExpenses       float64            `json:"monthly_expenses"`
	SavingsRate           float64            `json:"savings_rate"`
	SavingsConsistency    float64            `json:"savings_consistency"`
	BudgetCompliance      float64            `json:"budget_compliance"`
	ExpensePredictability float64            `json:"expense_predictability"`
	IncomeStability       float64            `json:"income_stability"`
	FinancialDiscipline   int                `json:"financial_discipline"`
	TopExpenseCategories  map[string]float64 `json:"top_expense_categories"`
	SeasonalMultiplier    float64            `json:"seasonal_multiplier"`
	SavingsGoals          []SavingsGoalInfo  `json:"savings_goals"`
	GoalAchievementRate   float64            `json:"goal_achievement_rate"`
	EmergencyFundMonths   float64            `json:"emergency_fund_months"`
	RecentLargePurchases  []RecentPurchase   `json:"recent_large_purchases"`
	BudgetAlerts          []BudgetAlert      `json:"budget_alerts"`
	SpendingLimits        map[string]float64 `json:"spending_limits"`
}

// SavingsGoalInfo representa información de una meta de ahorro
type SavingsGoalInfo struct {
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	CurrentAmount float64 `json:"current_amount"`
	TargetAmount  float64 `json:"target_amount"`
	Progress      float64 `json:"progress"`
}

// RecentPurchase representa una compra reciente
type RecentPurchase struct {
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

// BudgetAlert representa una alerta de presupuesto
type BudgetAlert struct {
	Category   string  `json:"category"`
	Spent      float64 `json:"spent"`
	Budget     float64 `json:"budget"`
	Percentage float64 `json:"percentage"`
	Status     string  `json:"status"`
}
