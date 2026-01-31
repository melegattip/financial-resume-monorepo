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

// CreditAnalysisService maneja el análisis y mejora crediticia con IA
// Implementa el principio de Single Responsibility
type CreditAnalysisService struct {
	client  *openai.Client
	useMock bool
}

// NewCreditAnalysisService crea una nueva instancia del servicio de análisis crediticio
func NewCreditAnalysisService() *CreditAnalysisService {
	useMock := os.Getenv("USE_AI_MOCK") == "true"
	log.Printf("📊 Credit Analysis Service - USE_AI_MOCK env var: '%s'", os.Getenv("USE_AI_MOCK"))

	var client *openai.Client
	if !useMock {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Println("⚠️ OPENAI_API_KEY not set, using mock responses")
			useMock = true
		} else {
			maskedKey := apiKey[:10] + "..." + apiKey[len(apiKey)-10:]
			log.Printf("✅ OPENAI_API_KEY configured for Credit Analysis Service: %s", maskedKey)
			client = openai.NewClient(apiKey)
		}
	}

	if useMock {
		log.Println("🎭 Credit Analysis Service initialized in MOCK mode")
	} else {
		log.Println("🧠 Credit Analysis Service initialized in REAL AI mode (OpenAI GPT-4)")
	}

	return &CreditAnalysisService{
		client:  client,
		useMock: useMock,
	}
}

// GenerateImprovementPlan genera un plan personalizado de mejora crediticia
func (s *CreditAnalysisService) GenerateImprovementPlan(ctx context.Context, data ports.FinancialAnalysisData) (*ports.CreditPlan, error) {
	if s.useMock {
		return s.getMockCreditPlan(data), nil
	}

	log.Printf("🧠 Using REAL AI analysis for credit improvement plan: user %s", data.UserID)

	prompt := s.buildCreditPlanPrompt(data)

	// Crear contexto con timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Eres un asesor financiero especializado en mejora crediticia. Analiza la situación financiera del usuario y genera un plan detallado y realista para mejorar su score crediticio. Responde ÚNICAMENTE con un JSON válido en el formato solicitado.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.2,
		MaxTokens:   1000,
	})

	if err != nil {
		log.Printf("Error calling OpenAI for credit plan: %v", err)
		return nil, fmt.Errorf("error conectando con OpenAI: %w", err)
	}

	var response struct {
		CurrentScore   int                    `json:"current_score"`
		TargetScore    int                    `json:"target_score"`
		TimelineMonths int                    `json:"timeline_months"`
		Actions        []ports.CreditAction   `json:"actions"`
		KeyMetrics     map[string]interface{} `json:"key_metrics"`
	}

	content := resp.Choices[0].Message.Content
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &response); err != nil {
		log.Printf("Error parsing AI credit plan response: %v", err)
		return nil, fmt.Errorf("error procesando respuesta de IA: %w", err)
	}

	return &ports.CreditPlan{
		CurrentScore:   response.CurrentScore,
		TargetScore:    response.TargetScore,
		TimelineMonths: response.TimelineMonths,
		Actions:        response.Actions,
		KeyMetrics:     response.KeyMetrics,
		GeneratedAt:    time.Now(),
	}, nil
}

// CalculateCreditScore calcula un score crediticio basado en datos financieros
func (s *CreditAnalysisService) CalculateCreditScore(ctx context.Context, data ports.FinancialAnalysisData) (int, error) {
	if s.useMock {
		return s.calculateMockCreditScore(data), nil
	}

	log.Printf("🧠 Calculating credit score for user: %s", data.UserID)

	prompt := s.buildCreditScorePrompt(data)

	// Crear contexto con timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Eres un experto en análisis crediticio. Calcula un score crediticio realista basado en los datos financieros proporcionados. Responde ÚNICAMENTE con un número entero entre 1 y 1000.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.1,
		MaxTokens:   50,
	})

	if err != nil {
		log.Printf("Error calling OpenAI for credit score: %v", err)
		return 0, fmt.Errorf("error conectando con OpenAI: %w", err)
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	var score int
	if _, err := fmt.Sscanf(content, "%d", &score); err != nil {
		log.Printf("Error parsing credit score response: %v", err)
		return s.calculateMockCreditScore(data), nil // Fallback a mock
	}

	// Validar rango
	if score < 1 {
		score = 1
	} else if score > 1000 {
		score = 1000
	}

	return score, nil
}

// === MÉTODOS PRIVADOS ===

// buildCreditPlanPrompt construye el prompt para generar plan de mejora crediticia
func (s *CreditAnalysisService) buildCreditPlanPrompt(data ports.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Analiza este perfil financiero y genera un plan de mejora crediticia:

Datos actuales:
- Ingresos mensuales: $%.0f
- Gastos mensuales: $%.0f
- Tasa de ahorro: %.1f%%
- Score financiero actual: %d/1000
- Estabilidad de ingresos: %.1f%%
- Período de análisis: %s

Genera un plan en JSON:
{
  "current_score": 1-1000,
  "target_score": 1-1000,
  "timeline_months": 6-24,
  "actions": [
    {
      "title": "Título de la acción",
      "description": "Descripción detallada",
      "priority": "high|medium|low",
      "timeline": "1-3 meses|3-6 meses|6-12 meses",
      "impact": 10-100,
      "difficulty": "easy|medium|hard"
    }
  ],
  "key_metrics": {
    "target_savings_rate": 0.0-1.0,
    "debt_to_income_target": 0.0-1.0,
    "emergency_fund_months": 1-12,
    "diversification_score": 0.0-1.0
  }
}
`,
		data.TotalIncome,
		data.TotalExpenses,
		data.SavingsRate*100,
		data.FinancialScore,
		data.IncomeStability*100,
		data.Period,
	)
}

// buildCreditScorePrompt construye el prompt para calcular score crediticio
func (s *CreditAnalysisService) buildCreditScorePrompt(data ports.FinancialAnalysisData) string {
	return fmt.Sprintf(`
Calcula un score crediticio realista (300-850) basado en:

- Ingresos mensuales: $%.0f
- Gastos mensuales: $%.0f
- Tasa de ahorro: %.1f%%
- Estabilidad de ingresos: %.1f%%
- Ratio deuda/ingresos estimado: %.1f%%

Responde solo con el número del score (ejemplo: 720)
`,
		data.TotalIncome,
		data.TotalExpenses,
		data.SavingsRate*100,
		data.IncomeStability*100,
		(data.TotalExpenses/data.TotalIncome)*100,
	)
}

// calculateMockCreditScore calcula un score crediticio mock basado en datos
func (s *CreditAnalysisService) calculateMockCreditScore(data ports.FinancialAnalysisData) int {
	baseScore := 600

	// Ajustar por tasa de ahorro (30% del score)
	if data.SavingsRate >= 0.2 {
		baseScore += 100
	} else if data.SavingsRate >= 0.1 {
		baseScore += 50
	} else if data.SavingsRate <= 0 {
		baseScore -= 50
	}

	// Ajustar por estabilidad de ingresos (25% del score)
	stabilityBonus := int(data.IncomeStability * 100)
	baseScore += stabilityBonus

	// Ajustar por ratio gastos/ingresos (25% del score)
	if data.TotalIncome > 0 {
		expenseRatio := data.TotalExpenses / data.TotalIncome
		if expenseRatio <= 0.5 {
			baseScore += 75
		} else if expenseRatio <= 0.7 {
			baseScore += 50
		} else if expenseRatio <= 0.9 {
			baseScore += 25
		} else {
			baseScore -= 25
		}
	}

	// Ajustar por score financiero general (20% del score)
	financialBonus := (data.FinancialScore - 500) / 10
	baseScore += financialBonus

	// Asegurar rango válido
	if baseScore < 1 {
		baseScore = 1
	} else if baseScore > 1000 {
		baseScore = 1000
	}

	return baseScore
}

// getMockCreditPlan genera un plan de mejora crediticia mock
func (s *CreditAnalysisService) getMockCreditPlan(data ports.FinancialAnalysisData) *ports.CreditPlan {
	currentScore := s.calculateMockCreditScore(data)
	targetScore := currentScore + 100
	if targetScore > 1000 {
		targetScore = 1000
	}

	timeline := 12
	if data.SavingsRate > 0.2 {
		timeline = 8
	} else if data.SavingsRate < 0.1 {
		timeline = 18
	}

	actions := []ports.CreditAction{
		{
			Title:       "Optimizar tasa de ahorro",
			Description: fmt.Sprintf("Aumentar la tasa de ahorro del %.1f%% actual al 25%% para mejorar el perfil financiero", data.SavingsRate*100),
			Priority:    "high",
			Timeline:    "3 meses",
			Impact:      80,
			Difficulty:  "medium",
		},
		{
			Title:       "Estabilizar flujo de ingresos",
			Description: "Crear fuentes adicionales de ingresos para mejorar la estabilidad financiera",
			Priority:    "medium",
			Timeline:    "6 meses",
			Impact:      60,
			Difficulty:  "hard",
		},
		{
			Title:       "Reducir ratio gastos/ingresos",
			Description: fmt.Sprintf("Reducir gastos del %.1f%% actual de los ingresos a máximo 70%%", (data.TotalExpenses/data.TotalIncome)*100),
			Priority:    "high",
			Timeline:    "3 meses",
			Impact:      70,
			Difficulty:  "medium",
		},
	}

	// Ajustar acciones según situación específica
	if data.SavingsRate >= 0.2 {
		actions[0].Priority = "low"
		actions[0].Impact = 30
	}

	return &ports.CreditPlan{
		CurrentScore:   currentScore,
		TargetScore:    targetScore,
		TimelineMonths: timeline,
		Actions:        actions,
		KeyMetrics: map[string]interface{}{
			"target_savings_rate":   0.25,
			"debt_to_income_target": 0.30,
			"emergency_fund_months": 6,
			"diversification_score": 0.8,
		},
		GeneratedAt: time.Now(),
	}
}
