package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	"github.com/sashabaranov/go-openai"
)

// AIAnalysisService maneja el análisis financiero con IA
// Implementa el principio de Single Responsibility
type AIAnalysisService struct {
	client  *openai.Client
	useMock bool
}

// NewAIAnalysisService crea una nueva instancia del servicio de análisis IA
func NewAIAnalysisService() *AIAnalysisService {
	useMock := os.Getenv("USE_AI_MOCK") == "true"
	log.Printf("🧠 AI Analysis Service - USE_AI_MOCK env var: '%s'", os.Getenv("USE_AI_MOCK"))

	var client *openai.Client
	if !useMock {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Println("⚠️ OPENAI_API_KEY not set, using mock responses")
			useMock = true
		} else {
			maskedKey := apiKey[:10] + "..." + apiKey[len(apiKey)-10:]
			log.Printf("✅ OPENAI_API_KEY configured for Analysis Service: %s", maskedKey)
			client = openai.NewClient(apiKey)
		}
	}

	if useMock {
		log.Println("🎭 AI Analysis Service initialized in MOCK mode")
	} else {
		log.Println("🧠 AI Analysis Service initialized in REAL AI mode (OpenAI GPT-4)")
	}

	return &AIAnalysisService{
		client:  client,
		useMock: useMock,
	}
}

// AnalyzeFinancialHealth analiza la salud financiera del usuario
func (s *AIAnalysisService) AnalyzeFinancialHealth(ctx context.Context, data ports.FinancialAnalysisData) (*ports.HealthAnalysis, error) {
	if s.useMock {
		return s.getMockHealthAnalysis(data), nil
	}

	// Generar insights primero
	insights, err := s.GenerateInsights(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("error generando insights para análisis de salud: %w", err)
	}

	// Calcular score basado en los insights y datos
	score := s.calculateHealthScore(data)
	level, message := s.getHealthLevelAndMessage(score)

	return &ports.HealthAnalysis{
		Score:       score,
		Level:       level,
		Message:     message,
		Insights:    insights,
		GeneratedAt: time.Now(),
	}, nil
}

// GenerateInsights genera insights financieros usando IA
func (s *AIAnalysisService) GenerateInsights(ctx context.Context, data ports.FinancialAnalysisData) ([]ports.AIInsight, error) {
	if s.useMock {
		return s.getMockInsights(data), nil
	}

	log.Printf("🧠 Using REAL AI analysis for user: %s", data.UserID)

	prompt := s.buildInsightsPrompt(data)

	// Crear contexto con timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Eres un asesor financiero experto. Analiza los datos financieros y genera insights precisos y accionables. Responde ÚNICAMENTE con un JSON válido en el formato solicitado.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   800,
	})

	if err != nil {
		log.Printf("Error calling OpenAI for insights: %v", err)
		return nil, fmt.Errorf("error conectando con OpenAI: %w", err)
	}

	var insights []ports.AIInsight
	content := resp.Choices[0].Message.Content
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &insights); err != nil {
		log.Printf("Error parsing AI insights response: %v", err)
		return nil, fmt.Errorf("error procesando respuesta de IA: %w", err)
	}

	return insights, nil
}

// === MÉTODOS PRIVADOS ===

// buildInsightsPrompt construye el prompt para análisis de insights
func (s *AIAnalysisService) buildInsightsPrompt(data ports.FinancialAnalysisData) string {
	topCategories := GetTopExpenseCategories(data.ExpensesByCategory, 3)

	return fmt.Sprintf(`
Datos financieros para análisis:
- Ingresos: $%.0f
- Gastos: $%.0f  
- Tasa de ahorro: %.1f%%
- Score financiero: %d/1000
- Top categorías de gasto: %s
- Período: %s

Genera exactamente 3 insights financieros personalizados en formato JSON:
[
  {
    "title": "Título del insight",
    "description": "Descripción detallada y accionable",
    "impact": "high|medium|low",
    "score": 100-1000,
    "action_type": "save|optimize|alert|invest",
    "category": "categoria_relevante"
  }
]
`,
		data.TotalIncome,
		data.TotalExpenses,
		data.SavingsRate*100,
		data.FinancialScore,
		FormatTopCategories(topCategories),
		data.Period,
	)
}

// calculateHealthScore calcula el score de salud financiera
func (s *AIAnalysisService) calculateHealthScore(data ports.FinancialAnalysisData) int {
	score := 0

	// Tasa de ahorro (40% del score)
	savingsScore := int(data.SavingsRate * 400)
	if savingsScore > 400 {
		savingsScore = 400
	}
	score += savingsScore

	// Estabilidad de ingresos (30% del score)
	incomeScore := int(data.IncomeStability * 300)
	score += incomeScore

	// Balance ingresos vs gastos (30% del score)
	if data.TotalIncome > 0 {
		balanceRatio := (data.TotalIncome - data.TotalExpenses) / data.TotalIncome
		if balanceRatio > 0 {
			balanceScore := int(balanceRatio * 300)
			if balanceScore > 300 {
				balanceScore = 300
			}
			score += balanceScore
		}
	}

	// Asegurar que el score esté en el rango 0-1000
	if score > 1000 {
		score = 1000
	}
	if score < 0 {
		score = 0
	}

	return score
}

// getHealthLevelAndMessage determina el nivel y mensaje basado en el score
func (s *AIAnalysisService) getHealthLevelAndMessage(score int) (string, string) {
	if score >= 800 {
		return "Excelente", "¡Tu salud financiera es excelente! Sigue así."
	} else if score >= 600 {
		return "Bueno", "Tu salud financiera es buena, con oportunidades de mejora."
	} else if score >= 400 {
		return "Regular", "Tu salud financiera es regular. Hay varias áreas de mejora."
	} else {
		return "Mejorable", "Tu salud financiera necesita atención. ¡Pero puedes mejorar!"
	}
}

// getMockHealthAnalysis genera un análisis de salud mock
func (s *AIAnalysisService) getMockHealthAnalysis(data ports.FinancialAnalysisData) *ports.HealthAnalysis {
	score := s.calculateHealthScore(data)
	level, message := s.getHealthLevelAndMessage(score)
	insights := s.getMockInsights(data)

	return &ports.HealthAnalysis{
		Score:       score,
		Level:       level,
		Message:     message,
		Insights:    insights,
		GeneratedAt: time.Now(),
	}
}

// getMockInsights genera insights mock para testing
func (s *AIAnalysisService) getMockInsights(data ports.FinancialAnalysisData) []ports.AIInsight {
	// ✅ Calcular porcentaje de gastos de forma segura
	expensePercentage := 0.0
	if data.TotalIncome > 0 {
		expensePercentage = (data.TotalExpenses / data.TotalIncome) * 100
	} else {
		// Caso especial: ingresos cero, asignar alto impacto a gastos
		expensePercentage = 100.0
	}

	// ✅ Calcular score de ahorro de forma segura (mínimo 50 puntos)
	savingsScore := int(data.SavingsRate * 500)
	if savingsScore < 50 {
		savingsScore = 50
	}

	insights := []ports.AIInsight{
		{
			Title:       "Optimiza tus gastos principales",
			Description: fmt.Sprintf("Tus gastos representan %.1f%% de tus ingresos. Considera reducir gastos en las categorías principales.", expensePercentage),
			Impact:      GetImpactLevel(expensePercentage),
			Score:       CalculateCategoryScore(expensePercentage),
			ActionType:  "optimize",
			Category:    "gastos",
		},
		{
			Title:       "Mejora tu tasa de ahorro",
			Description: fmt.Sprintf("Tu tasa de ahorro actual es %.1f%%. Intenta alcanzar al menos 20%% para una mejor salud financiera.", data.SavingsRate*100),
			Impact:      GetImpactLevel(data.SavingsRate * 100),
			Score:       savingsScore,
			ActionType:  "save",
			Category:    "ahorro",
		},
		{
			Title:       "Diversifica tus ingresos",
			Description: "Considera crear fuentes adicionales de ingresos para mejorar tu estabilidad financiera.",
			Impact:      "medium",
			Score:       400,
			ActionType:  "invest",
			Category:    "ingresos",
		},
	}

	return insights
}
