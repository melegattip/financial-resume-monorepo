package domain

import "time"

// FinancialAnalysisData represents the data used for financial analysis.
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

// AIInsight represents a single AI-generated financial insight.
type AIInsight struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Score       int    `json:"score"`
	ActionType  string `json:"action_type"`
	Category    string `json:"category"`
}

// HealthAnalysis represents the result of a financial health analysis.
type HealthAnalysis struct {
	Score       int         `json:"score"`
	Level       string      `json:"level"`
	Message     string      `json:"message"`
	Insights    []AIInsight `json:"insights"`
	GeneratedAt time.Time   `json:"generated_at"`
}

// SavingsGoalInfo represents information about a savings goal.
type SavingsGoalInfo struct {
	Name          string    `json:"name"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	Progress      float64   `json:"progress"`
	TargetDate    time.Time `json:"target_date"`
}

// UserFinancialProfile represents the user's current financial profile for purchase analysis.
type UserFinancialProfile struct {
	CurrentBalance       float64            `json:"current_balance"`
	MonthlyIncome        float64            `json:"monthly_income"`
	MonthlyExpenses      float64            `json:"monthly_expenses"`
	SavingsRate          float64            `json:"savings_rate"`
	IncomeStability      float64            `json:"income_stability"`
	FinancialDiscipline  int                `json:"financial_discipline"`
	TopExpenseCategories map[string]float64 `json:"top_expense_categories"`
	SavingsGoals         []SavingsGoalInfo  `json:"savings_goals"`
}

// PurchaseAnalysisRequest represents a request to analyze a potential purchase.
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

// PurchaseDecision represents the AI decision about a potential purchase.
type PurchaseDecision struct {
	CanBuy       bool      `json:"can_buy"`
	Confidence   float64   `json:"confidence"`
	Reasoning    string    `json:"reasoning"`
	Alternatives []string  `json:"alternatives"`
	ImpactScore  int       `json:"impact_score"`
	GeneratedAt  time.Time `json:"generated_at"`
}

// Alternative represents a cheaper or more viable alternative to a purchase.
type Alternative struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Savings     float64 `json:"savings"`
	Feasibility string  `json:"feasibility"`
}

// CreditPlan represents a personalized credit improvement plan.
type CreditPlan struct {
	CurrentScore   int                    `json:"current_score"`
	TargetScore    int                    `json:"target_score"`
	TimelineMonths int                    `json:"timeline_months"`
	Actions        []CreditAction         `json:"actions"`
	KeyMetrics     map[string]interface{} `json:"key_metrics"`
	GeneratedAt    time.Time              `json:"generated_at"`
}

// CreditAction represents a specific action to improve credit score.
type CreditAction struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Timeline    string `json:"timeline"`
	Impact      int    `json:"impact"`
	Difficulty  string `json:"difficulty"`
}

// CreditScoreResponse wraps the calculated credit score with metadata.
type CreditScoreResponse struct {
	Score       int       `json:"score"`
	UserID      string    `json:"user_id"`
	CalculatedAt time.Time `json:"calculated_at"`
}
