package insights

import (
	"math"
	"time"
)

// CreditScoringMatrix define la matriz de scoring crediticio
type CreditScoringMatrix struct {
	PaymentHistory     CreditFactor `json:"payment_history"`     // 35%
	CreditUtilization  CreditFactor `json:"credit_utilization"`  // 30%
	IncomeStability    CreditFactor `json:"income_stability"`    // 15%
	SavingsCapacity    CreditFactor `json:"savings_capacity"`    // 10%
	FinancialDiversity CreditFactor `json:"financial_diversity"` // 10%
}

// CreditFactor representa un factor individual del score crediticio
type CreditFactor struct {
	Score       int      `json:"score"`        // 0-1000
	Weight      float64  `json:"weight"`       // Peso en el cálculo total
	Status      string   `json:"status"`       // Excelente, Bueno, Regular, Malo
	Trend       string   `json:"trend"`        // Mejorando, Estable, Empeorando
	Description string   `json:"description"`  // Descripción del factor
	ActionItems []string `json:"action_items"` // Acciones para mejorar
}

// CreditProfile representa el perfil crediticio completo del usuario
type CreditProfile struct {
	OverallScore     int                   `json:"overall_score"`     // Score general 0-1000
	CreditLevel      string                `json:"credit_level"`      // Excelente, Bueno, Regular, Malo
	ScoringMatrix    CreditScoringMatrix   `json:"scoring_matrix"`    // Matriz detallada
	ImprovementPlan  CreditImprovementPlan `json:"improvement_plan"`  // Plan de mejora
	PurchaseCapacity PurchaseCapacity      `json:"purchase_capacity"` // Capacidad de compra
	RiskAssessment   RiskAssessment        `json:"risk_assessment"`   // Evaluación de riesgo
	LastUpdated      time.Time             `json:"last_updated"`      // Última actualización
}

// CreditImprovementPlan define un plan para mejorar el score crediticio
type CreditImprovementPlan struct {
	CurrentScore    int                 `json:"current_score"`
	TargetScore     int                 `json:"target_score"`
	EstimatedMonths int                 `json:"estimated_months"`
	Priority        []ImprovementAction `json:"priority"`   // Acciones prioritarias
	QuickWins       []ImprovementAction `json:"quick_wins"` // Mejoras rápidas
	LongTerm        []ImprovementAction `json:"long_term"`  // Mejoras a largo plazo
}

// ImprovementAction representa una acción específica para mejorar el score
type ImprovementAction struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Impact        int    `json:"impact"`         // Puntos que puede sumar (0-100)
	Difficulty    string `json:"difficulty"`     // Fácil, Medio, Difícil
	EstimatedDays int    `json:"estimated_days"` // Días estimados para completar
	Category      string `json:"category"`       // payment, savings, income, etc.
	XPReward      int    `json:"xp_reward"`      // Puntos XP por completar
	Status        string `json:"status"`         // pending, in_progress, completed
}

// PurchaseCapacity evalúa la capacidad de realizar compras
type PurchaseCapacity struct {
	MonthlyDisposableIncome float64            `json:"monthly_disposable_income"`
	RecommendedSpendingCap  float64            `json:"recommended_spending_cap"`
	EmergencyFundMonths     float64            `json:"emergency_fund_months"`
	CategoryLimits          map[string]float64 `json:"category_limits"` // Límites por categoría
	SafetyBuffer            float64            `json:"safety_buffer"`   // Buffer de seguridad
}

// RiskAssessment evalúa los riesgos financieros
type RiskAssessment struct {
	OverallRisk       string       `json:"overall_risk"`       // Bajo, Medio, Alto
	RiskFactors       []RiskFactor `json:"risk_factors"`       // Factores de riesgo identificados
	ProtectiveFactors []string     `json:"protective_factors"` // Factores que reducen riesgo
	Recommendations   []string     `json:"recommendations"`    // Recomendaciones de mitigación
}

// RiskFactor representa un factor de riesgo específico
type RiskFactor struct {
	Type        string  `json:"type"`        // income_volatility, high_expenses, etc.
	Severity    string  `json:"severity"`    // Bajo, Medio, Alto
	Impact      float64 `json:"impact"`      // Impacto en el score (negativo)
	Description string  `json:"description"` // Descripción del riesgo
}

// PurchaseDecision representa la decisión sobre una compra
type PurchaseDecision struct {
	CanAfford      bool                  `json:"can_afford"`
	Confidence     float64               `json:"confidence"`      // 0-1
	Recommendation string                `json:"recommendation"`  // approve, reject, conditional
	Reasoning      string                `json:"reasoning"`       // Explicación de la decisión
	Alternatives   []PurchaseAlternative `json:"alternatives"`    // Alternativas sugeridas
	ImpactAnalysis PurchaseImpact        `json:"impact_analysis"` // Análisis de impacto
	Conditions     []string              `json:"conditions"`      // Condiciones para aprobar
}

// PurchaseAlternative representa una alternativa a la compra
type PurchaseAlternative struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Savings     float64 `json:"savings"`
	Reasoning   string  `json:"reasoning"`
}

// PurchaseImpact analiza el impacto de una compra en las finanzas
type PurchaseImpact struct {
	ScoreImpact         int     `json:"score_impact"`          // Cambio en score crediticio
	SavingsRateImpact   float64 `json:"savings_rate_impact"`   // Cambio en tasa de ahorro
	EmergencyFundImpact float64 `json:"emergency_fund_impact"` // Impacto en fondo de emergencia
	MonthlyBudgetImpact float64 `json:"monthly_budget_impact"` // Impacto en presupuesto mensual
	RecoveryTimeMonths  int     `json:"recovery_time_months"`  // Meses para recuperar estabilidad
}

// calculateCreditProfile calcula el perfil crediticio completo
func (s *Service) calculateCreditProfile(data *AnalyzedFinancialData) *CreditProfile {
	// Calcular matriz de scoring
	matrix := s.calculateCreditScoringMatrix(data)

	// Calcular score general
	overallScore := s.calculateOverallCreditScore(matrix)

	// Determinar nivel crediticio
	creditLevel := s.getCreditLevel(overallScore)

	// Generar plan de mejora
	improvementPlan := s.generateImprovementPlan(matrix, overallScore)

	// Calcular capacidad de compra
	purchaseCapacity := s.calculatePurchaseCapacity(data)

	// Evaluar riesgos
	riskAssessment := s.assessFinancialRisks(data, matrix)

	return &CreditProfile{
		OverallScore:     overallScore,
		CreditLevel:      creditLevel,
		ScoringMatrix:    matrix,
		ImprovementPlan:  improvementPlan,
		PurchaseCapacity: purchaseCapacity,
		RiskAssessment:   riskAssessment,
		LastUpdated:      time.Now(),
	}
}

// calculateCreditScoringMatrix calcula la matriz de scoring detallada
func (s *Service) calculateCreditScoringMatrix(data *AnalyzedFinancialData) CreditScoringMatrix {
	return CreditScoringMatrix{
		PaymentHistory:     s.calculatePaymentHistoryFactor(data),
		CreditUtilization:  s.calculateCreditUtilizationFactor(data),
		IncomeStability:    s.calculateIncomeStabilityFactor(data),
		SavingsCapacity:    s.calculateSavingsCapacityFactor(data),
		FinancialDiversity: s.calculateFinancialDiversityFactor(data),
	}
}

// calculatePaymentHistoryFactor calcula el factor de historial de pagos
func (s *Service) calculatePaymentHistoryFactor(data *AnalyzedFinancialData) CreditFactor {
	// Analizar pagos puntuales, gastos pendientes, etc.
	score := 800 // Base score

	// Penalizar por gastos pendientes
	if data.TotalExpenses > 0 {
		pendingRatio := 0.1 // Mock - calcular ratio real de gastos pendientes
		if pendingRatio > 0.1 {
			score -= int(pendingRatio * 300)
		}
	}

	score = int(math.Max(0, math.Min(1000, float64(score))))

	return CreditFactor{
		Score:       score,
		Weight:      0.35,
		Status:      s.getFactorStatus(score),
		Trend:       "Estable",
		Description: "Historial de pagos y cumplimiento de obligaciones",
		ActionItems: s.getPaymentHistoryActions(score),
	}
}

// calculateCreditUtilizationFactor calcula el factor de utilización de crédito
func (s *Service) calculateCreditUtilizationFactor(data *AnalyzedFinancialData) CreditFactor {
	score := 750 // Base score

	// Analizar ratio de gastos vs ingresos
	if data.TotalIncome > 0 {
		utilizationRatio := data.TotalExpenses / data.TotalIncome
		if utilizationRatio < 0.3 {
			score = 950 // Excelente utilización
		} else if utilizationRatio < 0.5 {
			score = 800 // Buena utilización
		} else if utilizationRatio < 0.7 {
			score = 600 // Regular utilización
		} else {
			score = 300 // Alta utilización
		}
	}

	return CreditFactor{
		Score:       score,
		Weight:      0.30,
		Status:      s.getFactorStatus(score),
		Trend:       "Mejorando",
		Description: "Ratio de utilización de ingresos vs gastos",
		ActionItems: s.getCreditUtilizationActions(score),
	}
}

// calculateIncomeStabilityFactor calcula el factor de estabilidad de ingresos
func (s *Service) calculateIncomeStabilityFactor(data *AnalyzedFinancialData) CreditFactor {
	score := 600 // Base score

	if data.IncomeStability.IsStable {
		score = 850 + int(data.IncomeStability.RecurringIncomeRatio*150)
	} else {
		score = 400 + int(data.IncomeStability.RecurringIncomeRatio*400)
	}

	score = int(math.Max(0, math.Min(1000, float64(score))))

	return CreditFactor{
		Score:       score,
		Weight:      0.15,
		Status:      s.getFactorStatus(score),
		Trend:       s.getIncomeStabilityTrend(data.IncomeStability),
		Description: "Estabilidad y consistencia de ingresos",
		ActionItems: s.getIncomeStabilityActions(score),
	}
}

// calculateSavingsCapacityFactor calcula el factor de capacidad de ahorro
func (s *Service) calculateSavingsCapacityFactor(data *AnalyzedFinancialData) CreditFactor {
	score := int(s.calculateSavingsScore(data.SavingsRate))

	return CreditFactor{
		Score:       score,
		Weight:      0.10,
		Status:      s.getFactorStatus(score),
		Trend:       s.getSavingsTrend(data.SavingsRate),
		Description: "Capacidad de generar y mantener ahorros",
		ActionItems: s.getSavingsCapacityActions(score),
	}
}

// calculateFinancialDiversityFactor calcula el factor de diversidad financiera
func (s *Service) calculateFinancialDiversityFactor(data *AnalyzedFinancialData) CreditFactor {
	score := int(s.calculateDiversificationScore(data.ExpenseCategories))

	return CreditFactor{
		Score:       score,
		Weight:      0.10,
		Status:      s.getFactorStatus(score),
		Trend:       "Estable",
		Description: "Diversificación de gastos y gestión financiera",
		ActionItems: s.getFinancialDiversityActions(score),
	}
}

// calculateOverallCreditScore calcula el score crediticio general
func (s *Service) calculateOverallCreditScore(matrix CreditScoringMatrix) int {
	weightedScore := float64(matrix.PaymentHistory.Score)*matrix.PaymentHistory.Weight +
		float64(matrix.CreditUtilization.Score)*matrix.CreditUtilization.Weight +
		float64(matrix.IncomeStability.Score)*matrix.IncomeStability.Weight +
		float64(matrix.SavingsCapacity.Score)*matrix.SavingsCapacity.Weight +
		float64(matrix.FinancialDiversity.Score)*matrix.FinancialDiversity.Weight

	return int(math.Max(0, math.Min(1000, weightedScore)))
}

// Funciones auxiliares para determinar estados y tendencias
func (s *Service) getFactorStatus(score int) string {
	if score >= 800 {
		return "Excelente"
	} else if score >= 600 {
		return "Bueno"
	} else if score >= 400 {
		return "Regular"
	} else {
		return "Malo"
	}
}

func (s *Service) getCreditLevel(score int) string {
	if score >= 800 {
		return "Excelente"
	} else if score >= 700 {
		return "Muy Bueno"
	} else if score >= 600 {
		return "Bueno"
	} else if score >= 500 {
		return "Regular"
	} else {
		return "Malo"
	}
}

// Funciones para generar acciones de mejora por factor
func (s *Service) getPaymentHistoryActions(score int) []string {
	if score < 600 {
		return []string{
			"Configurar pagos automáticos para evitar retrasos",
			"Priorizar el pago de deudas pendientes",
			"Establecer recordatorios de fechas de vencimiento",
		}
	}
	return []string{"Mantener el buen historial de pagos"}
}

func (s *Service) getCreditUtilizationActions(score int) []string {
	if score < 700 {
		return []string{
			"Reducir gastos en categorías no esenciales",
			"Aumentar ingresos con actividades adicionales",
			"Establecer límites de gasto por categoría",
		}
	}
	return []string{"Mantener el buen control de gastos"}
}

func (s *Service) getIncomeStabilityActions(score int) []string {
	if score < 600 {
		return []string{
			"Diversificar fuentes de ingreso",
			"Buscar ingresos más estables y recurrentes",
			"Crear un fondo de emergencia para volatilidad",
		}
	}
	return []string{"Mantener la estabilidad de ingresos"}
}

func (s *Service) getSavingsCapacityActions(score int) []string {
	if score < 600 {
		return []string{
			"Establecer un objetivo de ahorro mensual",
			"Automatizar transferencias a cuenta de ahorros",
			"Revisar y reducir gastos innecesarios",
		}
	}
	return []string{"Mantener o aumentar la capacidad de ahorro"}
}

func (s *Service) getFinancialDiversityActions(score int) []string {
	if score < 600 {
		return []string{
			"Diversificar gastos entre diferentes categorías",
			"Evitar concentrar gastos en una sola área",
			"Planificar presupuesto balanceado",
		}
	}
	return []string{"Mantener la diversificación financiera"}
}

// Funciones para determinar tendencias
func (s *Service) getIncomeStabilityTrend(stability IncomeStabilityAnalysis) string {
	if stability.IncomeVariation < 0.1 {
		return "Estable"
	} else if stability.IncomeVariation < 0.3 {
		return "Mejorando"
	} else {
		return "Empeorando"
	}
}

func (s *Service) getSavingsTrend(savingsRate float64) string {
	if savingsRate > 0.2 {
		return "Excelente"
	} else if savingsRate > 0.1 {
		return "Mejorando"
	} else if savingsRate > 0.05 {
		return "Estable"
	} else {
		return "Preocupante"
	}
}

// generateImprovementPlan genera un plan de mejora crediticia
func (s *Service) generateImprovementPlan(matrix CreditScoringMatrix, currentScore int) CreditImprovementPlan {
	targetScore := currentScore + 100
	if targetScore > 1000 {
		targetScore = 1000
	}

	var priority []ImprovementAction
	var quickWins []ImprovementAction
	var longTerm []ImprovementAction

	// Analizar cada factor y generar acciones
	if matrix.PaymentHistory.Score < 700 {
		priority = append(priority, ImprovementAction{
			ID:            "payment_history_1",
			Title:         "Mejorar historial de pagos",
			Description:   "Configurar pagos automáticos y eliminar pagos pendientes",
			Impact:        80,
			Difficulty:    "Medio",
			EstimatedDays: 30,
			Category:      "payment",
			XPReward:      100,
			Status:        "pending",
		})
	}

	if matrix.SavingsCapacity.Score < 600 {
		quickWins = append(quickWins, ImprovementAction{
			ID:            "savings_1",
			Title:         "Aumentar tasa de ahorro",
			Description:   "Reducir gastos no esenciales para ahorrar al menos 10% de ingresos",
			Impact:        60,
			Difficulty:    "Fácil",
			EstimatedDays: 15,
			Category:      "savings",
			XPReward:      75,
			Status:        "pending",
		})
	}

	if matrix.IncomeStability.Score < 700 {
		longTerm = append(longTerm, ImprovementAction{
			ID:            "income_1",
			Title:         "Diversificar fuentes de ingresos",
			Description:   "Crear fuentes adicionales de ingresos para mejorar la estabilidad",
			Impact:        120,
			Difficulty:    "Difícil",
			EstimatedDays: 90,
			Category:      "income",
			XPReward:      150,
			Status:        "pending",
		})
	}

	estimatedMonths := 6
	if currentScore < 500 {
		estimatedMonths = 12
	} else if currentScore > 700 {
		estimatedMonths = 3
	}

	return CreditImprovementPlan{
		CurrentScore:    currentScore,
		TargetScore:     targetScore,
		EstimatedMonths: estimatedMonths,
		Priority:        priority,
		QuickWins:       quickWins,
		LongTerm:        longTerm,
	}
}

// calculatePurchaseCapacity calcula la capacidad de compra del usuario
func (s *Service) calculatePurchaseCapacity(data *AnalyzedFinancialData) PurchaseCapacity {
	monthlyIncome := data.TotalIncome
	monthlyExpenses := data.TotalExpenses
	disposableIncome := monthlyIncome - monthlyExpenses

	// Calcular límite recomendado (50% del ingreso disponible)
	recommendedCap := disposableIncome * 0.5
	if recommendedCap < 0 {
		recommendedCap = 0
	}

	// Calcular meses de fondo de emergencia
	emergencyFundMonths := 0.0
	if monthlyExpenses > 0 {
		emergencyFundMonths = math.Max(0, disposableIncome) / monthlyExpenses
	}

	// Calcular límites por categoría
	categoryLimits := make(map[string]float64)
	for _, category := range data.ExpenseCategories {
		// Límite basado en el promedio histórico + 20%
		categoryLimits[category.CategoryName] = category.Amount * 1.2
	}

	// Buffer de seguridad (20% del ingreso disponible)
	safetyBuffer := disposableIncome * 0.2

	return PurchaseCapacity{
		MonthlyDisposableIncome: disposableIncome,
		RecommendedSpendingCap:  recommendedCap,
		EmergencyFundMonths:     emergencyFundMonths,
		CategoryLimits:          categoryLimits,
		SafetyBuffer:            safetyBuffer,
	}
}

// assessFinancialRisks evalúa los riesgos financieros
func (s *Service) assessFinancialRisks(data *AnalyzedFinancialData, matrix CreditScoringMatrix) RiskAssessment {
	var riskFactors []RiskFactor
	var protectiveFactors []string
	var recommendations []string

	// Evaluar volatilidad de ingresos
	if !data.IncomeStability.IsStable {
		riskFactors = append(riskFactors, RiskFactor{
			Type:        "income_volatility",
			Severity:    "Medio",
			Impact:      -50,
			Description: "Ingresos variables pueden afectar la capacidad de pago",
		})
		recommendations = append(recommendations, "Crear un fondo de emergencia más robusto")
	}

	// Evaluar alta utilización de ingresos
	if data.TotalIncome > 0 && (data.TotalExpenses/data.TotalIncome) > 0.7 {
		riskFactors = append(riskFactors, RiskFactor{
			Type:        "high_expense_ratio",
			Severity:    "Alto",
			Impact:      -80,
			Description: "Gastos muy altos en relación a los ingresos",
		})
		recommendations = append(recommendations, "Reducir gastos no esenciales")
	}

	// Evaluar baja tasa de ahorro
	if data.SavingsRate < 0.05 {
		riskFactors = append(riskFactors, RiskFactor{
			Type:        "low_savings_rate",
			Severity:    "Medio",
			Impact:      -40,
			Description: "Baja capacidad de ahorro limita la flexibilidad financiera",
		})
		recommendations = append(recommendations, "Establecer un plan de ahorro automático")
	}

	// Identificar factores protectivos
	if data.SavingsRate > 0.2 {
		protectiveFactors = append(protectiveFactors, "Excelente capacidad de ahorro")
	}
	if data.IncomeStability.IsStable {
		protectiveFactors = append(protectiveFactors, "Ingresos estables y predecibles")
	}

	// Determinar riesgo general
	overallRisk := "Bajo"
	if len(riskFactors) > 2 {
		overallRisk = "Alto"
	} else if len(riskFactors) > 0 {
		overallRisk = "Medio"
	}

	return RiskAssessment{
		OverallRisk:       overallRisk,
		RiskFactors:       riskFactors,
		ProtectiveFactors: protectiveFactors,
		Recommendations:   recommendations,
	}
}
