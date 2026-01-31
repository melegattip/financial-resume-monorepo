package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// AIService maneja la integración con OpenAI para análisis financiero
type AIService struct {
	client  *openai.Client
	useMock bool
}

// NewAIService crea una nueva instancia del servicio de IA
func NewAIService() *AIService {
	useMock := os.Getenv("USE_AI_MOCK") == "true"
	log.Printf("🤖 AI Service - USE_AI_MOCK env var: '%s'", os.Getenv("USE_AI_MOCK"))

	var client *openai.Client
	if !useMock {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Println("⚠️ OPENAI_API_KEY not set, using mock responses")
			useMock = true
		} else {
			// Mostrar solo los primeros y últimos caracteres de la API key por seguridad
			maskedKey := apiKey[:10] + "..." + apiKey[len(apiKey)-10:]
			log.Printf("✅ OPENAI_API_KEY configured: %s", maskedKey)
			client = openai.NewClient(apiKey)
		}
	}

	if useMock {
		log.Println("🎭 AI Service initialized in MOCK mode")
	} else {
		log.Println("🧠 AI Service initialized in REAL AI mode (OpenAI GPT-4)")
	}

	return &AIService{
		client:  client,
		useMock: useMock,
	}
}

// FinancialAnalysisRequest representa la solicitud de análisis financiero
type FinancialAnalysisRequest struct {
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

// SavingsGoalInfo representa información de una meta de ahorro relevante
type SavingsGoalInfo struct {
	Name          string  `json:"name"`
	Category      string  `json:"category"`
	CurrentAmount float64 `json:"current_amount"`
	TargetAmount  float64 `json:"target_amount"`
	Progress      float64 `json:"progress"`
}

// UserFinancialProfile representa el perfil financiero completo del usuario
type UserFinancialProfile struct {
	// Datos básicos actuales
	CurrentBalance  float64 `json:"current_balance"`
	MonthlyIncome   float64 `json:"monthly_income"`
	MonthlyExpenses float64 `json:"monthly_expenses"`
	SavingsRate     float64 `json:"savings_rate"`

	// Comportamiento financiero (scores 0-1)
	SavingsConsistency    float64 `json:"savings_consistency"`    // Qué tan consistente es ahorrando
	BudgetCompliance      float64 `json:"budget_compliance"`      // Cumplimiento de presupuesto
	ExpensePredictability float64 `json:"expense_predictability"` // Predictibilidad de gastos
	IncomeStability       float64 `json:"income_stability"`       // Estabilidad de ingresos
	FinancialDiscipline   int     `json:"financial_discipline"`   // Score 0-1000

	// Patrones de gasto (top 5 categorías)
	TopExpenseCategories map[string]float64 `json:"top_expense_categories"`
	SeasonalMultiplier   float64            `json:"seasonal_multiplier"` // Factor estacional actual

	// Metas y objetivos
	SavingsGoals        []SavingsGoalInfo `json:"savings_goals"`
	GoalAchievementRate float64           `json:"goal_achievement_rate"` // % de metas logradas históricamente
	EmergencyFundMonths float64           `json:"emergency_fund_months"` // Meses de gastos cubiertos

	// Contexto de compras recientes (últimas 3 grandes)
	RecentLargePurchases []RecentPurchase `json:"recent_large_purchases"`

	// Alertas y límites activos
	BudgetAlerts   []BudgetAlert  `json:"budget_alerts"`
	SpendingLimits SpendingLimits `json:"spending_limits"`
}

// RecentPurchase representa una compra grande reciente
type RecentPurchase struct {
	ItemName string  `json:"item_name"`
	Amount   float64 `json:"amount"`
	Date     string  `json:"date"`
	Category string  `json:"category"`
	Outcome  string  `json:"outcome"` // "successful", "struggled", "regretted"
}

// BudgetAlert representa una alerta de presupuesto activa
type BudgetAlert struct {
	Category   string  `json:"category"`
	Limit      float64 `json:"limit"`
	Spent      float64 `json:"spent"`
	Remaining  float64 `json:"remaining"`
	AlertLevel string  `json:"alert_level"` // "warning", "danger", "exceeded"
}

// SpendingLimits representa límites de gasto configurados
type SpendingLimits struct {
	DailyDiscretionary float64 `json:"daily_discretionary"`
	MonthlyLuxury      float64 `json:"monthly_luxury"`
	EmergencyThreshold float64 `json:"emergency_threshold"`
}

// CanIBuyRequest representa una consulta de compra
type CanIBuyRequest struct {
	UserID               string               `json:"user_id"`
	ItemName             string               `json:"item_name"`
	Amount               float64              `json:"amount"`
	Description          string               `json:"description"`       // Descripción opcional
	PaymentTypes         []string             `json:"payment_types"`     // Array de tipos de pago
	IsNecessary          bool                 `json:"is_necessary"`      // Si es una necesidad urgente
	UserFinancialProfile UserFinancialProfile `json:"financial_profile"` // Perfil financiero completo
	SavingsGoal          float64              `json:"savings_goal"`      // Mantener por compatibilidad
}

// CanIBuyResponse representa la respuesta de decisión de compra
type CanIBuyResponse struct {
	CanBuy       bool     `json:"can_buy"`
	Confidence   float64  `json:"confidence"`
	Reasoning    string   `json:"reasoning"`
	Alternatives []string `json:"alternatives"`
	ImpactScore  int      `json:"impact_score"`
}

// CreditImprovementPlan representa un plan de mejora crediticia
type CreditImprovementPlan struct {
	CurrentScore   int                    `json:"current_score"`
	TargetScore    int                    `json:"target_score"`
	TimelineMonths int                    `json:"timeline_months"`
	Actions        []CreditAction         `json:"actions"`
	KeyMetrics     map[string]interface{} `json:"key_metrics"`
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

// GenerateFinancialInsights genera insights financieros usando IA
func (s *AIService) GenerateFinancialInsights(ctx context.Context, request FinancialAnalysisRequest) ([]AIInsight, error) {
	if s.useMock {
		return s.getMockInsights(request), nil
	}

	prompt := s.buildInsightsPrompt(request)

	// Crear contexto con timeout más corto
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Eres un asesor financiero experto. Analiza los datos financieros y genera insights útiles en español. Responde SOLO con un JSON válido que contenga un array de insights.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   800, // Reducido para respuestas más rápidas
	})

	if err != nil {
		log.Printf("Error calling OpenAI: %v", err)
		return s.getMockInsights(request), nil
	}

	var insights []AIInsight
	content := resp.Choices[0].Message.Content

	// Limpiar el contenido si viene con markdown
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &insights); err != nil {
		log.Printf("Error parsing AI response: %v", err)
		return s.getMockInsights(request), nil
	}

	return insights, nil
}

// CanIBuyThis analiza si el usuario puede permitirse una compra
func (s *AIService) CanIBuyThis(ctx context.Context, request CanIBuyRequest) (*CanIBuyResponse, error) {
	if s.useMock {
		log.Printf("❌ CanIBuy feature disabled - AI not available")
		return nil, fmt.Errorf("análisis de compra no disponible: IA no configurada")
	}

	log.Printf("🧠 Using REAL AI analysis for item: %s ($%.0f)", request.ItemName, request.Amount)

	prompt := s.buildCanIBuyPrompt(request)

	// Crear contexto con timeout más corto
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Eres un asesor financiero experto especializado en análisis de compras inteligentes. Tu trabajo es analizar la situación financiera del usuario y dar recomendaciones precisas y personalizadas. Debes considerar el tipo de pago, la necesidad real del artículo, y el impacto a largo plazo. Sé específico con números y porcentajes. Responde ÚNICAMENTE con un JSON válido en el formato solicitado.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   500, // Reducido para respuestas más rápidas
	})

	if err != nil {
		log.Printf("Error calling OpenAI: %v", err)
		return nil, fmt.Errorf("error conectando con OpenAI: %w", err)
	}

	var response CanIBuyResponse
	content := resp.Choices[0].Message.Content
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &response); err != nil {
		log.Printf("Error parsing AI response: %v", err)
		return nil, fmt.Errorf("error procesando respuesta de IA: %w", err)
	}

	return &response, nil
}

// GenerateCreditImprovementPlan genera un plan personalizado de mejora crediticia
func (s *AIService) GenerateCreditImprovementPlan(ctx context.Context, request FinancialAnalysisRequest) (*CreditImprovementPlan, error) {
	if s.useMock {
		return s.getMockCreditPlan(request), nil
	}

	prompt := s.buildCreditPlanPrompt(request)

	// Crear contexto con timeout más corto
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Eres un experto en mejora de perfil crediticio. Crea un plan detallado y realista para mejorar la salud financiera en español. Responde SOLO con un JSON válido.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.5,
		MaxTokens:   1000, // Reducido para respuestas más rápidas
	})

	if err != nil {
		log.Printf("Error calling OpenAI: %v", err)
		return s.getMockCreditPlan(request), nil
	}

	var plan CreditImprovementPlan
	content := resp.Choices[0].Message.Content
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &plan); err != nil {
		log.Printf("Error parsing AI response: %v", err)
		return s.getMockCreditPlan(request), nil
	}

	return &plan, nil
}

// buildInsightsPrompt construye el prompt para generar insights (optimizado)
func (s *AIService) buildInsightsPrompt(request FinancialAnalysisRequest) string {
	// Solo enviar las top 3 categorías de gastos para reducir tokens
	topCategories := getTopExpenseCategories(request.ExpensesByCategory, 3)

	return fmt.Sprintf(`
Datos: Ingresos $%.0f, Gastos $%.0f, Ahorro %.0f%%, Score %d/1000
Top gastos: %s

Genera 3 insights JSON:
[{"title":"","description":"","impact":"high/medium/low","score":100-1000,"action_type":"save/optimize/alert/invest","category":""}]
`,
		request.TotalIncome,
		request.TotalExpenses,
		request.SavingsRate*100,
		request.FinancialScore,
		formatTopCategories(topCategories),
	)
}

// buildCanIBuyPrompt construye el prompt para análisis de compra (optimizado)
func (s *AIService) buildCanIBuyPrompt(request CanIBuyRequest) string {
	available := request.UserFinancialProfile.MonthlyIncome - request.UserFinancialProfile.MonthlyExpenses

	// Construir información adicional
	necessityText := "No es urgente"
	if request.IsNecessary {
		necessityText = "Es una necesidad urgente/esencial"
	}

	descriptionText := ""
	if request.Description != "" {
		descriptionText = fmt.Sprintf("\nDescripción: %s", request.Description)
	}

	// Construir texto de tipos de pago múltiples
	paymentTypeText := "No especificado"
	if len(request.PaymentTypes) > 0 {
		var paymentLabels []string
		for _, paymentType := range request.PaymentTypes {
			switch paymentType {
			case "contado":
				paymentLabels = append(paymentLabels, "Pago de contado")
			case "cuotas":
				paymentLabels = append(paymentLabels, "Plan de pagos/cuotas")
			case "ahorro":
				paymentLabels = append(paymentLabels, "Necesita ahorrar para esto")
			}
		}
		paymentTypeText = strings.Join(paymentLabels, " + ")
	}

	// Construir información de metas de ahorro
	savingsGoalsInfo := ""
	if len(request.UserFinancialProfile.SavingsGoals) > 0 {
		savingsGoalsInfo = "\nMETAS DE AHORRO EXISTENTES:\n"
		for _, goal := range request.UserFinancialProfile.SavingsGoals {
			savingsGoalsInfo += fmt.Sprintf("- %s (%s): $%.0f de $%.0f (%.1f%% completado)\n",
				goal.Name, goal.Category, goal.CurrentAmount, goal.TargetAmount, goal.Progress*100)
		}
	} else {
		savingsGoalsInfo = "\nMETAS DE AHORRO: No hay metas relacionadas con esta compra"
	}

	// Construir información del perfil financiero de manera compacta
	profile := request.UserFinancialProfile

	// Información de comportamiento financiero
	behaviorInfo := fmt.Sprintf("\nCOMPORTAMIENTO FINANCIERO (Scores 0-100):\n- Disciplina: %d/1000 | Consistencia ahorro: %.0f%% | Cumplimiento presupuesto: %.0f%%\n- Estabilidad ingresos: %.0f%% | Predictibilidad gastos: %.0f%% | Factor estacional: %.1fx",
		profile.FinancialDiscipline, profile.SavingsConsistency*100, profile.BudgetCompliance*100,
		profile.IncomeStability*100, profile.ExpensePredictability*100, profile.SeasonalMultiplier)

	// Top categorías de gasto
	expenseInfo := "\nPATRONES DE GASTO (Top 5):\n"
	for category, amount := range profile.TopExpenseCategories {
		percentage := (amount / profile.MonthlyExpenses) * 100
		expenseInfo += fmt.Sprintf("- %s: $%.0f (%.1f%%) ", category, amount, percentage)
	}

	// Contexto de compras recientes
	purchaseHistory := ""
	if len(profile.RecentLargePurchases) > 0 {
		purchaseHistory = "\nCOMPRAS GRANDES RECIENTES:\n"
		for _, purchase := range profile.RecentLargePurchases {
			purchaseHistory += fmt.Sprintf("- %s: $%.0f (%s, %s) ", purchase.ItemName, purchase.Amount, purchase.Date, purchase.Outcome)
		}
	}

	// Alertas de presupuesto activas
	budgetAlerts := ""
	if len(profile.BudgetAlerts) > 0 {
		budgetAlerts = "\nALERTAS PRESUPUESTO ACTIVAS:\n"
		for _, alert := range profile.BudgetAlerts {
			budgetAlerts += fmt.Sprintf("- %s: $%.0f/$%.0f (%s) ", alert.Category, alert.Spent, alert.Limit, alert.AlertLevel)
		}
	}

	return fmt.Sprintf(`
ANÁLISIS FINANCIERO INTELIGENTE CON PERFIL COMPLETO:

SITUACIÓN FINANCIERA MENSUAL:
- Balance actual: $%.0f 
- Ingresos MENSUALES: $%.0f (anuales: $%.0f)
- Gastos MENSUALES: $%.0f (anuales: $%.0f)
- Disponible MENSUAL: $%.0f
- Tasa ahorro: %.1f%% | Fondo emergencia: %.1f meses | Logro metas: %.0f%%
%s%s%s%s%s

COMPRA A ANALIZAR:
- Artículo: %s | Monto: $%.0f%s | Tipos: %s | Urgente: %s

INSTRUCCIONES CRÍTICAS PARA LA IA:
1. IMPORTANTE: Los ingresos son MENSUALES, no anuales. Multiplica por 12 para cálculos anuales.
2. ANALIZAR DESCRIPCIÓN PARA DINERO DISPONIBLE: 
   - Buscar ventas de activos (auto, casa, inversiones)
   - Buscar intercambios o permutas
   - Buscar fondos específicos mencionados
   - SUMAR TODO EL DINERO DISPONIBLE antes de calcular faltante
3. MATEMÁTICA CORRECTA: Dinero total = Balance + Ventas + Fondos, Faltante = Precio - Dinero total
4. CONSIDERAR CAPACIDAD REAL: Con ingresos mensuales de $%.0f, la capacidad anual es $%.0f
5. USAR TODO EL CONTEXTO: Comportamiento, patrones, historial, alertas, metas
6. RECOMENDAR BASADO EN PERFIL: Si es disciplinado → más flexible, si es impulsivo → más estricto

ANÁLISIS MATEMÁTICO REQUERIDO:
1. Leer descripción y extraer TODOS los activos vendibles o fondos mencionados
2. Sumar: Balance actual + Ventas de activos + Fondos específicos = DINERO TOTAL DISPONIBLE
3. Calcular: Precio del artículo - DINERO TOTAL DISPONIBLE = FALTANTE REAL
4. Si FALTANTE > 0: calcular meses necesarios = FALTANTE ÷ ahorro mensual
5. Basar decisión en FALTANTE REAL, no en balance actual solamente

EJEMPLO CONCRETO para esta compra:
- Si menciona "auto que se vende en 14 millones" → agregar $14,000,000 al dinero disponible
- Si menciona "fondo de ahorro" → extraer cantidad específica y agregar
- Dinero total disponible = $%.0f + ventas + fondos
- Solo entonces calcular si puede comprar el artículo de $%.0f

CRITERIOS ADAPTATIVOS:
- Contado: Balance + capacidad reposición + historial compras similares
- Cuotas: Capacidad mensual + cumplimiento presupuesto + estabilidad ingresos  
- Ahorro: Consistencia + metas existentes + disciplina histórica
- IMPORTANTE: Ajustar criterios según perfil (disciplinado vs impulsivo)

RESPUESTA PERSONALIZADA (JSON):
{
  "can_buy": true/false,
  "confidence": 0.0-1.0,
  "reasoning": "Análisis PERSONALIZADO considerando comportamiento, patrones y contexto completo. RECORDAR: ingresos son MENSUALES.",
  "alternatives": ["Basadas en perfil y comportamiento"],
  "impact_score": 1-1000
}
`,
		profile.CurrentBalance,
		profile.MonthlyIncome, profile.MonthlyIncome*12,
		profile.MonthlyExpenses, profile.MonthlyExpenses*12,
		available,
		profile.SavingsRate*100, profile.EmergencyFundMonths, profile.GoalAchievementRate*100,
		behaviorInfo, expenseInfo, purchaseHistory, budgetAlerts, savingsGoalsInfo,
		request.ItemName,
		request.Amount,
		descriptionText,
		paymentTypeText,
		necessityText,
		profile.MonthlyIncome, profile.MonthlyIncome*12,
		profile.CurrentBalance, // Balance actual para ejemplo
		request.Amount,         // Precio del artículo para ejemplo
	)
}

// buildCreditPlanPrompt construye el prompt para plan de mejora crediticia
func (s *AIService) buildCreditPlanPrompt(request FinancialAnalysisRequest) string {
	return fmt.Sprintf(`
Crea un plan de mejora de perfil crediticio basado en estos datos:

Situación actual:
- Score financiero: %d/1000
- Ingresos: $%.2f
- Gastos: $%.2f
- Tasa de ahorro: %.1f%%
- Estabilidad: %.2f

Gastos por categoría:
%s

Crea un plan que incluya:
1. Score objetivo realista
2. Timeline en meses
3. Acciones específicas priorizadas
4. Métricas clave a mejorar

Formato JSON requerido:
{
  "current_score": %d,
  "target_score": 850-950,
  "timeline_months": 6-24,
  "actions": [
    {
      "title": "Acción específica",
      "description": "Descripción detallada",
      "priority": "high|medium|low",
      "timeline": "1-3 meses",
      "impact": 50-200,
      "difficulty": "easy|medium|hard"
    }
  ],
  "key_metrics": {
    "target_savings_rate": 0.25,
    "debt_to_income_target": 0.30,
    "emergency_fund_months": 6
  }
}
`,
		request.FinancialScore,
		request.TotalIncome,
		request.TotalExpenses,
		request.SavingsRate*100,
		request.IncomeStability,
		formatExpensesByCategory(request.ExpensesByCategory),
		request.FinancialScore,
	)
}

// Funciones auxiliares para mocks y formateo
func formatExpensesByCategory(expenses map[string]float64) string {
	var lines []string
	for category, amount := range expenses {
		lines = append(lines, fmt.Sprintf("- %s: $%.2f", category, amount))
	}
	return strings.Join(lines, "\n")
}

// getTopExpenseCategories obtiene las N categorías con mayor gasto
func getTopExpenseCategories(expenses map[string]float64, limit int) map[string]float64 {
	type categoryAmount struct {
		category string
		amount   float64
	}

	var categories []categoryAmount
	for cat, amount := range expenses {
		categories = append(categories, categoryAmount{cat, amount})
	}

	// Ordenar por cantidad descendente
	for i := 0; i < len(categories)-1; i++ {
		for j := i + 1; j < len(categories); j++ {
			if categories[i].amount < categories[j].amount {
				categories[i], categories[j] = categories[j], categories[i]
			}
		}
	}

	// Tomar solo las top N
	result := make(map[string]float64)
	for i := 0; i < len(categories) && i < limit; i++ {
		result[categories[i].category] = categories[i].amount
	}

	return result
}

// formatTopCategories formatea las top categorías de manera concisa
func formatTopCategories(expenses map[string]float64) string {
	var parts []string
	for category, amount := range expenses {
		parts = append(parts, fmt.Sprintf("%s $%.0f", category, amount))
	}
	return strings.Join(parts, ", ")
}

// getMockInsights devuelve insights de prueba
func (s *AIService) getMockInsights(request FinancialAnalysisRequest) []AIInsight {
	insights := []AIInsight{
		{
			Title:       "Excelente capacidad de ahorro",
			Description: fmt.Sprintf("Estás ahorrando %.1f%% de tus ingresos, superando el promedio nacional. Considera explorar opciones de inversión para hacer crecer tu dinero.", request.SavingsRate*100),
			Impact:      "high",
			Score:       920,
			ActionType:  "invest",
			Category:    "ahorro",
		},
	}

	// Agregar insight basado en categoría con mayor gasto
	var maxCategory string
	var maxAmount float64
	for category, amount := range request.ExpensesByCategory {
		if amount > maxAmount {
			maxAmount = amount
			maxCategory = category
		}
	}

	if maxCategory != "" {
		// ✅ Calcular porcentaje de forma segura
		percentage := 0.0
		if request.TotalExpenses > 0 {
			percentage = (maxAmount / request.TotalExpenses) * 100
		} else {
			percentage = 100.0 // Si no hay gastos totales, asignar 100%
		}

		insights = append(insights, AIInsight{
			Title:       fmt.Sprintf("Mayor gasto: %s", maxCategory),
			Description: fmt.Sprintf("El %s representa %.1f%% de tus gastos ($%.0f). Revisa si hay oportunidades de optimización en esta categoría.", maxCategory, percentage, maxAmount),
			Impact:      GetImpactLevel(percentage),
			Score:       CalculateCategoryScore(percentage),
			ActionType:  "optimize",
			Category:    maxCategory,
		})
	}

	// Insight sobre estabilidad de ingresos
	if request.IncomeStability < 0.8 {
		insights = append(insights, AIInsight{
			Title:       "Ingresos variables",
			Description: "Tus ingresos muestran variabilidad. Considera diversificar fuentes de ingresos o crear un fondo de emergencia más robusto.",
			Impact:      "medium",
			Score:       600,
			ActionType:  "save",
			Category:    "ingresos",
		})
	}

	return insights
}

// getMockCanIBuyResponse devuelve una respuesta inteligente para decisiones de compra
// BuildUserFinancialProfile construye un perfil financiero completo basado en datos disponibles
func (s *AIService) BuildUserFinancialProfile(userID string, currentBalance, monthlyIncome, monthlyExpenses float64, savingsGoals []SavingsGoalInfo) UserFinancialProfile {
	// Datos básicos
	savingsRate := 0.0
	if monthlyIncome > 0 {
		savingsRate = math.Max(0, (monthlyIncome-monthlyExpenses)/monthlyIncome)
	}

	// Calcular scores de comportamiento basados en datos disponibles
	savingsConsistency := s.calculateSavingsConsistency(userID, savingsRate)
	budgetCompliance := s.calculateBudgetCompliance(userID, monthlyExpenses)
	expensePredictability := s.calculateExpensePredictability(userID)
	incomeStability := s.calculateIncomeStability(userID, monthlyIncome)
	financialDiscipline := s.calculateFinancialDiscipline(savingsConsistency, budgetCompliance, len(savingsGoals))

	// Patrones de gasto (simulados inteligentemente)
	topExpenseCategories := s.generateTopExpenseCategories(monthlyExpenses)
	seasonalMultiplier := s.calculateSeasonalMultiplier()

	// Metas y objetivos
	goalAchievementRate := s.calculateGoalAchievementRate(userID, savingsGoals)
	emergencyFundMonths := currentBalance / math.Max(monthlyExpenses, 1)

	// Compras recientes y alertas (simuladas)
	recentPurchases := s.generateRecentPurchases(userID, monthlyIncome)
	budgetAlerts := s.generateBudgetAlerts(monthlyExpenses, topExpenseCategories)
	spendingLimits := s.generateSpendingLimits(monthlyIncome, monthlyExpenses)

	return UserFinancialProfile{
		CurrentBalance:        currentBalance,
		MonthlyIncome:         monthlyIncome,
		MonthlyExpenses:       monthlyExpenses,
		SavingsRate:           savingsRate,
		SavingsConsistency:    savingsConsistency,
		BudgetCompliance:      budgetCompliance,
		ExpensePredictability: expensePredictability,
		IncomeStability:       incomeStability,
		FinancialDiscipline:   financialDiscipline,
		TopExpenseCategories:  topExpenseCategories,
		SeasonalMultiplier:    seasonalMultiplier,
		SavingsGoals:          savingsGoals,
		GoalAchievementRate:   goalAchievementRate,
		EmergencyFundMonths:   emergencyFundMonths,
		RecentLargePurchases:  recentPurchases,
		BudgetAlerts:          budgetAlerts,
		SpendingLimits:        spendingLimits,
	}
}

// Funciones auxiliares para calcular métricas de comportamiento
func (s *AIService) calculateSavingsConsistency(userID string, savingsRate float64) float64 {
	// Simulación inteligente basada en tasa de ahorro
	if savingsRate >= 0.2 {
		return 0.9 + (savingsRate-0.2)*0.5 // Alta consistencia
	} else if savingsRate >= 0.1 {
		return 0.6 + (savingsRate-0.1)*3 // Consistencia media
	} else {
		return savingsRate * 6 // Baja consistencia
	}
}

func (s *AIService) calculateBudgetCompliance(userID string, monthlyExpenses float64) float64 {
	// Simulación basada en estabilidad de gastos
	if monthlyExpenses < 1000000 { // Gastos bajos = mejor control
		return 0.85
	} else if monthlyExpenses < 3000000 { // Gastos medios
		return 0.75
	} else { // Gastos altos = más difícil controlar
		return 0.65
	}
}

func (s *AIService) calculateExpensePredictability(userID string) float64 {
	// Simulación estándar con variación por usuario
	hash := 0
	for _, c := range userID {
		hash += int(c)
	}
	return 0.7 + float64(hash%20)/100 // 0.7-0.89
}

func (s *AIService) calculateIncomeStability(userID string, monthlyIncome float64) float64 {
	// Ingresos más altos tienden a ser más estables
	if monthlyIncome >= 5000000 { // Ingresos altos
		return 0.9
	} else if monthlyIncome >= 2000000 { // Ingresos medios
		return 0.8
	} else { // Ingresos bajos
		return 0.7
	}
}

func (s *AIService) calculateFinancialDiscipline(savingsConsistency, budgetCompliance float64, goalsCount int) int {
	// Score 0-1000 basado en múltiples factores
	base := (savingsConsistency + budgetCompliance) / 2 * 600 // 0-600
	goalBonus := math.Min(float64(goalsCount)*50, 200)        // 0-200 por metas
	randomFactor := 200                                       // Factor base
	return int(base + goalBonus + float64(randomFactor))
}

func (s *AIService) generateTopExpenseCategories(monthlyExpenses float64) map[string]float64 {
	// Distribución típica de gastos
	categories := map[string]float64{
		"Alimentación": monthlyExpenses * 0.25,
		"Transporte":   monthlyExpenses * 0.20,
		"Vivienda":     monthlyExpenses * 0.30,
		"Servicios":    monthlyExpenses * 0.15,
		"Ocio":         monthlyExpenses * 0.10,
	}
	return categories
}

func (s *AIService) calculateSeasonalMultiplier() float64 {
	// Factor estacional basado en el mes actual
	month := time.Now().Month()
	switch month {
	case 12, 1: // Diciembre, Enero - temporada alta
		return 1.3
	case 6, 7: // Junio, Julio - vacaciones
		return 1.2
	case 3, 9: // Marzo, Septiembre - inicio de períodos
		return 1.1
	default:
		return 1.0
	}
}

func (s *AIService) calculateGoalAchievementRate(userID string, savingsGoals []SavingsGoalInfo) float64 {
	if len(savingsGoals) == 0 {
		return 0.5 // Valor neutral sin metas
	}

	// Calcular progreso promedio de metas activas
	totalProgress := 0.0
	for _, goal := range savingsGoals {
		totalProgress += goal.Progress
	}
	return totalProgress / float64(len(savingsGoals))
}

func (s *AIService) generateRecentPurchases(userID string, monthlyIncome float64) []RecentPurchase {
	// Generar compras recientes basadas en nivel de ingresos
	purchases := []RecentPurchase{}

	if monthlyIncome >= 3000000 {
		purchases = append(purchases, RecentPurchase{
			ItemName: "Laptop nueva",
			Amount:   2500000,
			Date:     "2024-01-15",
			Category: "Tecnología",
			Outcome:  "successful",
		})
	}

	if monthlyIncome >= 1500000 {
		purchases = append(purchases, RecentPurchase{
			ItemName: "Vacaciones familiares",
			Amount:   1800000,
			Date:     "2023-12-20",
			Category: "Ocio",
			Outcome:  "successful",
		})
	}

	return purchases
}

func (s *AIService) generateBudgetAlerts(monthlyExpenses float64, categories map[string]float64) []BudgetAlert {
	alerts := []BudgetAlert{}

	// Generar alerta si los gastos de ocio son altos
	if ocio, exists := categories["Ocio"]; exists && ocio > monthlyExpenses*0.15 {
		alerts = append(alerts, BudgetAlert{
			Category:   "Ocio",
			Limit:      monthlyExpenses * 0.15,
			Spent:      ocio,
			Remaining:  math.Max(0, monthlyExpenses*0.15-ocio),
			AlertLevel: "warning",
		})
	}

	return alerts
}

func (s *AIService) generateSpendingLimits(monthlyIncome, monthlyExpenses float64) SpendingLimits {
	available := monthlyIncome - monthlyExpenses

	return SpendingLimits{
		DailyDiscretionary: available * 0.1 / 30, // 10% del disponible diario
		MonthlyLuxury:      available * 0.3,      // 30% del disponible para lujos
		EmergencyThreshold: monthlyExpenses * 3,  // 3 meses de gastos
	}
}

// getMockCreditPlan devuelve un plan de mejora crediticia de prueba
func (s *AIService) getMockCreditPlan(request FinancialAnalysisRequest) *CreditImprovementPlan {
	targetScore := request.FinancialScore + 100
	if targetScore > 1000 {
		targetScore = 1000
	}

	timeline := 12
	if request.SavingsRate > 0.2 {
		timeline = 8
	}

	actions := []CreditAction{
		{
			Title:       "Optimizar tasa de ahorro",
			Description: "Aumentar la tasa de ahorro al 25% para mejorar el perfil financiero",
			Priority:    "high",
			Timeline:    "3 meses",
			Impact:      80,
			Difficulty:  "medium",
		},
		{
			Title:       "Diversificar ingresos",
			Description: "Crear fuentes adicionales de ingresos para mejorar la estabilidad financiera",
			Priority:    "medium",
			Timeline:    "6 meses",
			Impact:      120,
			Difficulty:  "hard",
		},
	}

	return &CreditImprovementPlan{
		CurrentScore:   request.FinancialScore,
		TargetScore:    targetScore,
		TimelineMonths: timeline,
		Actions:        actions,
		KeyMetrics: map[string]interface{}{
			"target_savings_rate":   0.25,
			"debt_to_income_target": 0.30,
			"emergency_fund_months": 6,
			"diversification_score": 0.8,
		},
	}
}

// Funciones auxiliares
// ✅ NOTA: Funciones getImpactLevel y calculateCategoryScore movidas a ai_utils.go
// para evitar duplicación y usar las versiones exportadas GetImpactLevel y CalculateCategoryScore
