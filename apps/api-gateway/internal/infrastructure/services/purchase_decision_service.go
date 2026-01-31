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

// PurchaseDecisionService maneja las decisiones de compra con IA
// Implementa el principio de Single Responsibility
type PurchaseDecisionService struct {
	client  *openai.Client
	useMock bool
}

// NewPurchaseDecisionService crea una nueva instancia del servicio de decisiones de compra
func NewPurchaseDecisionService() *PurchaseDecisionService {
	useMock := os.Getenv("USE_AI_MOCK") == "true"
	log.Printf("🛒 Purchase Decision Service - USE_AI_MOCK env var: '%s'", os.Getenv("USE_AI_MOCK"))

	var client *openai.Client
	if !useMock {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Println("⚠️ OPENAI_API_KEY not set, using mock responses")
			useMock = true
		} else {
			maskedKey := apiKey[:10] + "..." + apiKey[len(apiKey)-10:]
			log.Printf("✅ OPENAI_API_KEY configured for Purchase Decision Service: %s", maskedKey)
			client = openai.NewClient(apiKey)
		}
	}

	if useMock {
		log.Println("🎭 Purchase Decision Service initialized in MOCK mode")
	} else {
		log.Println("🧠 Purchase Decision Service initialized in REAL AI mode (OpenAI GPT-4)")
	}

	return &PurchaseDecisionService{
		client:  client,
		useMock: useMock,
	}
}

// CanIBuy analiza si el usuario puede realizar una compra específica
func (s *PurchaseDecisionService) CanIBuy(ctx context.Context, request ports.PurchaseAnalysisRequest) (*ports.PurchaseDecision, error) {
	if s.useMock {
		return s.getMockPurchaseDecision(request), nil
	}

	log.Printf("🧠 Using REAL AI analysis for purchase: %s ($%.0f)", request.ItemName, request.Amount)

	prompt := s.buildCanIBuyPrompt(request)

	// Crear contexto con timeout
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
		MaxTokens:   500,
	})

	if err != nil {
		log.Printf("Error calling OpenAI for purchase decision: %v", err)
		return nil, fmt.Errorf("error conectando con OpenAI: %w", err)
	}

	var response struct {
		CanBuy       bool     `json:"can_buy"`
		Confidence   float64  `json:"confidence"`
		Reasoning    string   `json:"reasoning"`
		Alternatives []string `json:"alternatives"`
		ImpactScore  int      `json:"impact_score"`
	}

	content := resp.Choices[0].Message.Content
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &response); err != nil {
		log.Printf("Error parsing AI purchase decision response: %v", err)
		return nil, fmt.Errorf("error procesando respuesta de IA: %w", err)
	}

	return &ports.PurchaseDecision{
		CanBuy:       response.CanBuy,
		Confidence:   response.Confidence,
		Reasoning:    response.Reasoning,
		Alternatives: response.Alternatives,
		ImpactScore:  response.ImpactScore,
		GeneratedAt:  time.Now(),
	}, nil
}

// SuggestAlternatives sugiere alternativas de compra más económicas
func (s *PurchaseDecisionService) SuggestAlternatives(ctx context.Context, request ports.PurchaseAnalysisRequest) ([]ports.Alternative, error) {
	if s.useMock {
		return s.getMockAlternatives(request), nil
	}

	log.Printf("🧠 Generating alternatives for: %s ($%.0f)", request.ItemName, request.Amount)

	prompt := s.buildAlternativesPrompt(request)

	// Crear contexto con timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Eres un asesor financiero experto. Genera alternativas de compra más económicas y viables. Responde ÚNICAMENTE con un JSON válido.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.4,
		MaxTokens:   400,
	})

	if err != nil {
		log.Printf("Error calling OpenAI for alternatives: %v", err)
		return nil, fmt.Errorf("error conectando con OpenAI: %w", err)
	}

	var alternatives []ports.Alternative
	content := resp.Choices[0].Message.Content
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if err := json.Unmarshal([]byte(content), &alternatives); err != nil {
		log.Printf("Error parsing AI alternatives response: %v", err)
		return nil, fmt.Errorf("error procesando respuesta de IA: %w", err)
	}

	return alternatives, nil
}

// === MÉTODOS PRIVADOS ===

// buildCanIBuyPrompt construye el prompt para análisis de compra
func (s *PurchaseDecisionService) buildCanIBuyPrompt(request ports.PurchaseAnalysisRequest) string {
	available := request.UserFinancialProfile.MonthlyIncome - request.UserFinancialProfile.MonthlyExpenses

	necessityText := "No es urgente"
	if request.IsNecessary {
		necessityText = "Es una necesidad urgente/esencial"
	}

	descriptionText := ""
	if request.Description != "" {
		descriptionText = fmt.Sprintf("\nDescripción: %s", request.Description)
	}

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

	savingsGoalsInfo := ""
	if len(request.UserFinancialProfile.SavingsGoals) > 0 {
		savingsGoalsInfo = fmt.Sprintf("\nMetas de ahorro activas: %d", len(request.UserFinancialProfile.SavingsGoals))
		for _, goal := range request.UserFinancialProfile.SavingsGoals {
			savingsGoalsInfo += fmt.Sprintf("\n- %s: $%.0f de $%.0f (%.1f%%)",
				goal.Name, goal.CurrentAmount, goal.TargetAmount, goal.Progress*100)
		}
	}

	return fmt.Sprintf(`
Analiza esta compra:
Artículo: %s
Precio: $%.0f%s
Necesidad: %s
Forma de pago: %s

Perfil financiero:
- Ingresos mensuales: $%.0f
- Gastos mensuales: $%.0f
- Disponible mensual: $%.0f
- Balance actual: $%.0f
- Tasa de ahorro: %.1f%%
- Score financiero: %d/1000%s

Responde en JSON:
{
  "can_buy": true/false,
  "confidence": 0.0-1.0,
  "reasoning": "Explicación detallada",
  "alternatives": ["alternativa1", "alternativa2"],
  "impact_score": 1-100
}
`,
		request.ItemName,
		request.Amount,
		descriptionText,
		necessityText,
		paymentTypeText,
		request.UserFinancialProfile.MonthlyIncome,
		request.UserFinancialProfile.MonthlyExpenses,
		available,
		request.UserFinancialProfile.CurrentBalance,
		request.UserFinancialProfile.SavingsRate*100,
		request.UserFinancialProfile.FinancialDiscipline,
		savingsGoalsInfo,
	)
}

// buildAlternativesPrompt construye el prompt para generar alternativas
func (s *PurchaseDecisionService) buildAlternativesPrompt(request ports.PurchaseAnalysisRequest) string {
	return fmt.Sprintf(`
Genera 3 alternativas más económicas para:
Artículo: %s
Precio original: $%.0f
Ingresos mensuales: $%.0f

Responde en JSON:
[
  {
    "name": "Nombre de la alternativa",
    "description": "Descripción breve",
    "savings": 0.0,
    "feasibility": "alta|media|baja"
  }
]
`,
		request.ItemName,
		request.Amount,
		request.UserFinancialProfile.MonthlyIncome,
	)
}

// getMockPurchaseDecision genera una decisión de compra mock
func (s *PurchaseDecisionService) getMockPurchaseDecision(request ports.PurchaseAnalysisRequest) *ports.PurchaseDecision {
	available := request.UserFinancialProfile.MonthlyIncome - request.UserFinancialProfile.MonthlyExpenses

	canBuy := request.Amount <= available*0.3 // Máximo 30% del disponible
	if request.IsNecessary {
		canBuy = request.Amount <= available*0.5 // 50% si es necesario
	}

	confidence := 0.7
	if canBuy {
		confidence = 0.8
	}

	reasoning := fmt.Sprintf("Basado en tu disponible mensual de $%.0f, ", available)
	if canBuy {
		reasoning += "puedes realizar esta compra sin comprometer tu estabilidad financiera."
	} else {
		reasoning += "esta compra podría comprometer tu estabilidad financiera. Considera ahorrar primero."
	}

	impactScore := int((request.Amount / available) * 100)
	if impactScore > 100 {
		impactScore = 100
	}

	return &ports.PurchaseDecision{
		CanBuy:       canBuy,
		Confidence:   confidence,
		Reasoning:    reasoning,
		Alternatives: []string{"Buscar ofertas", "Comprar usado", "Ahorrar y comprar después"},
		ImpactScore:  impactScore,
		GeneratedAt:  time.Now(),
	}
}

// getMockAlternatives genera alternativas mock
func (s *PurchaseDecisionService) getMockAlternatives(request ports.PurchaseAnalysisRequest) []ports.Alternative {
	return []ports.Alternative{
		{
			Name:        "Versión usada/reacondicionada",
			Description: "Buscar el mismo artículo en condición usada pero funcional",
			Savings:     request.Amount * 0.3,
			Feasibility: "alta",
		},
		{
			Name:        "Modelo anterior o similar",
			Description: "Considerar modelos anteriores con características similares",
			Savings:     request.Amount * 0.2,
			Feasibility: "alta",
		},
		{
			Name:        "Esperar ofertas o promociones",
			Description: "Aguardar a temporadas de descuentos como Black Friday",
			Savings:     request.Amount * 0.25,
			Feasibility: "media",
		},
	}
}
