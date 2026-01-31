package mocks

import (
	"time"

	"github.com/financial-ai-service/internal/core/ports"
)

// TestFinancialAnalysisData datos de prueba para análisis financiero
var TestFinancialAnalysisData = ports.FinancialAnalysisData{
	UserID:          "test-user-123",
	TotalIncome:     5000000,
	TotalExpenses:   3500000,
	SavingsRate:     0.3,
	IncomeStability: 0.8,
	FinancialScore:  750,
	Period:          "monthly",
	ExpensesByCategory: map[string]float64{
		"Alimentación": 1200000,
		"Transporte":   800000,
		"Vivienda":     1500000,
	},
}

// TestHealthAnalysis análisis de salud esperado
var TestHealthAnalysis = ports.HealthAnalysis{
	Score:   750,
	Level:   "Bueno",
	Message: "Tu salud financiera es sólida con algunas áreas de mejora",
	Insights: []ports.AIInsight{
		{
			Title:       "Excelente tasa de ahorro",
			Description: "Tu tasa de ahorro del 30% está muy por encima del promedio",
			Impact:      "Positivo",
			Score:       85,
			ActionType:  "maintain",
			Category:    "savings",
		},
		{
			Title:       "Oportunidad de optimización",
			Description: "Podrías reducir gastos en entretenimiento",
			Impact:      "Medio",
			Score:       65,
			ActionType:  "optimize",
			Category:    "expenses",
		},
	},
	GeneratedAt: time.Now(),
}

// TestInsights insights de prueba
var TestInsights = []ports.AIInsight{
	{
		Title:       "Optimización de gastos",
		Description: "Puedes reducir 15% en gastos de entretenimiento",
		Impact:      "Alto",
		Score:       80,
		ActionType:  "reduce",
		Category:    "entertainment",
	},
	{
		Title:       "Incremento de ahorros",
		Description: "Considera aumentar tu fondo de emergencia",
		Impact:      "Medio",
		Score:       70,
		ActionType:  "increase",
		Category:    "emergency_fund",
	},
}

// TestPurchaseAnalysisRequest solicitud de análisis de compra
var TestPurchaseAnalysisRequest = ports.PurchaseAnalysisRequest{
	UserID:       "test-user-123",
	ItemName:     "MacBook Pro",
	Amount:       8000000,
	Description:  "Laptop para trabajo",
	IsNecessary:  true,
	PaymentTypes: []string{"contado"},
	UserFinancialProfile: ports.UserFinancialProfile{
		CurrentBalance:      10000000,
		MonthlyIncome:       5000000,
		MonthlyExpenses:     3500000,
		SavingsRate:         0.3,
		IncomeStability:     0.8,
		FinancialDiscipline: 750,
		TopExpenseCategories: map[string]float64{
			"Alimentación": 1200000,
			"Transporte":   800000,
		},
	},
}

// TestPurchaseDecision decisión de compra de prueba
var TestPurchaseDecision = ports.PurchaseDecision{
	CanBuy:       true,
	Confidence:   0.85,
	Reasoning:    "Basado en tu situación financiera actual, puedes realizar esta compra sin comprometer tu estabilidad",
	Alternatives: []string{"Buscar ofertas", "Considerar modelo anterior"},
	ImpactScore:  25,
	GeneratedAt:  time.Now(),
}

// TestAlternatives alternativas de prueba
var TestAlternatives = []ports.Alternative{
	{
		Name:        "MacBook Air",
		Description: "Opción más económica con rendimiento similar",
		Savings:     3000000,
		Feasibility: "Alta",
	},
	{
		Name:        "Modelo anterior",
		Description: "MacBook Pro del año pasado con descuento",
		Savings:     2000000,
		Feasibility: "Media",
	},
}

// TestCreditPlan plan crediticio de prueba
var TestCreditPlan = ports.CreditPlan{
	CurrentScore:   750,
	TargetScore:    800,
	TimelineMonths: 12,
	Actions: []ports.CreditAction{
		{
			Title:       "Aumentar fondo de emergencia",
			Description: "Mantener 6 meses de gastos como reserva",
			Priority:    "alta",
			Timeline:    "3-6 meses",
			Impact:      30,
			Difficulty:  "media",
		},
		{
			Title:       "Diversificar inversiones",
			Description: "Invertir en diferentes instrumentos financieros",
			Priority:    "media",
			Timeline:    "6-12 meses",
			Impact:      20,
			Difficulty:  "media",
		},
	},
	KeyMetrics: map[string]interface{}{
		"savings_rate_improvement": 0.05,
		"debt_reduction_target":    0.15,
		"emergency_fund_months":    6,
	},
	GeneratedAt: time.Now(),
}

// TestOpenAIResponse respuesta mock de OpenAI
var TestOpenAIResponse = `{
	"score": 750,
	"level": "Bueno",
	"message": "Tu salud financiera es sólida con algunas áreas de mejora",
	"insights": [
		{
			"title": "Excelente tasa de ahorro",
			"description": "Tu tasa de ahorro del 30% está muy por encima del promedio",
			"impact": "Positivo",
			"score": 85,
			"action_type": "maintain",
			"category": "savings"
		}
	]
}`

// TestCacheData datos de cache de prueba
var TestCacheData = []byte(`{"score": 750, "level": "Bueno", "message": "Test cached data"}`)

// TestPurchaseDecisionCannotBuy decisión de compra negativa
var TestPurchaseDecisionCannotBuy = ports.PurchaseDecision{
	CanBuy:       false,
	Confidence:   0.2,
	Reasoning:    "Esta compra comprometería tu estabilidad financiera",
	Alternatives: []string{"Buscar ofertas", "Considerar modelo anterior"},
	ImpactScore:  80,
	GeneratedAt:  time.Now(),
}

// TestPurchaseDecisionLowConfidence decisión de compra con baja confianza
var TestPurchaseDecisionLowConfidence = ports.PurchaseDecision{
	CanBuy:       true,
	Confidence:   0.6,
	Reasoning:    "Puedes realizar esta compra pero con ciertas reservas",
	Alternatives: []string{"Buscar ofertas", "Considerar modelo anterior"},
	ImpactScore:  50,
	GeneratedAt:  time.Now(),
}
