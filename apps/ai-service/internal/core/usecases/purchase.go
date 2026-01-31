package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/financial-ai-service/internal/core/ports"
)

// PurchaseUseCase implementa el caso de uso para decisiones de compra
type PurchaseUseCase struct {
	openaiClient ports.OpenAIClient
	cacheClient  ports.CacheClient
}

// NewPurchaseUseCase crea un nuevo caso de uso de compra
func NewPurchaseUseCase(openaiClient ports.OpenAIClient, cacheClient ports.CacheClient) ports.PurchaseDecisionPort {
	return &PurchaseUseCase{
		openaiClient: openaiClient,
		cacheClient:  cacheClient,
	}
}

// CanIBuy analiza si el usuario puede realizar una compra específica
func (p *PurchaseUseCase) CanIBuy(ctx context.Context, request ports.PurchaseAnalysisRequest) (*ports.PurchaseDecision, error) {
	log.Printf("🛒 Analyzing purchase decision for user: %s, item: %s ($%.2f)",
		request.UserID, request.ItemName, request.Amount)

	// Intentar obtener del cache (TTL corto para decisiones de compra)
	cacheKey := fmt.Sprintf("purchase_decision:%s:%s:%.2f", request.UserID, request.ItemName, request.Amount)
	if cached, err := p.cacheClient.Get(ctx, cacheKey); err == nil {
		log.Printf("📋 Cache hit for purchase decision: %s", cacheKey)
		var decision ports.PurchaseDecision
		if err := json.Unmarshal(cached, &decision); err == nil {
			return &decision, nil
		}
	}

	// Construir prompts para IA
	systemPrompt := `Eres un asesor financiero experto especializado en análisis de compras inteligentes. 
	Tu trabajo es analizar la situación financiera del usuario y dar recomendaciones precisas y personalizadas. 
	Debes considerar el tipo de pago, la necesidad real del artículo, y el impacto a largo plazo. 
	Sé específico con números y porcentajes. 
	Responde ÚNICAMENTE con un JSON válido en el formato solicitado.`

	userPrompt := p.buildPurchaseAnalysisPrompt(request)

	// Llamar a OpenAI
	response, err := p.openaiClient.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		log.Printf("❌ Error generating purchase decision: %v", err)
		return nil, fmt.Errorf("error analyzing purchase decision: %w", err)
	}

	// Parsear respuesta
	var decision ports.PurchaseDecision
	if err := json.Unmarshal([]byte(response), &decision); err != nil {
		log.Printf("❌ Error parsing purchase decision response: %v", err)
		return nil, fmt.Errorf("error parsing purchase decision response: %w", err)
	}

	// Agregar timestamp
	decision.GeneratedAt = time.Now()

	// Guardar en cache (30 minutos)
	if cacheData, err := json.Marshal(decision); err == nil {
		p.cacheClient.Set(ctx, cacheKey, cacheData, 30*time.Minute)
	}

	log.Printf("✅ Purchase decision completed for user: %s (Can buy: %v, Confidence: %.2f)",
		request.UserID, decision.CanBuy, decision.Confidence)
	return &decision, nil
}

// SuggestAlternatives sugiere alternativas de compra más económicas
func (p *PurchaseUseCase) SuggestAlternatives(ctx context.Context, request ports.PurchaseAnalysisRequest) ([]ports.Alternative, error) {
	log.Printf("🔍 Generating alternatives for user: %s, item: %s", request.UserID, request.ItemName)

	// Intentar obtener del cache
	cacheKey := fmt.Sprintf("alternatives:%s:%s", request.UserID, request.ItemName)
	if cached, err := p.cacheClient.Get(ctx, cacheKey); err == nil {
		log.Printf("📋 Cache hit for alternatives: %s", cacheKey)
		var alternatives []ports.Alternative
		if err := json.Unmarshal(cached, &alternatives); err == nil {
			return alternatives, nil
		}
	}

	// Construir prompts para IA
	systemPrompt := `Eres un asesor financiero experto que genera alternativas de compra más económicas y viables. 
	Debes proporcionar opciones realistas y específicas que ayuden al usuario a ahorrar dinero.
	Responde ÚNICAMENTE con un JSON válido que contenga un array de alternativas.`

	userPrompt := p.buildAlternativesPrompt(request)

	// Llamar a OpenAI
	response, err := p.openaiClient.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		log.Printf("❌ Error generating alternatives: %v", err)
		return nil, fmt.Errorf("error generating alternatives: %w", err)
	}

	// Parsear respuesta
	var alternatives []ports.Alternative
	if err := json.Unmarshal([]byte(response), &alternatives); err != nil {
		log.Printf("❌ Error parsing alternatives response: %v", err)
		return nil, fmt.Errorf("error parsing alternatives response: %w", err)
	}

	// Guardar en cache (2 horas)
	if cacheData, err := json.Marshal(alternatives); err == nil {
		p.cacheClient.Set(ctx, cacheKey, cacheData, 2*time.Hour)
	}

	log.Printf("✅ Generated %d alternatives for user: %s", len(alternatives), request.UserID)
	return alternatives, nil
}

// buildPurchaseAnalysisPrompt construye el prompt para análisis de compra
func (p *PurchaseUseCase) buildPurchaseAnalysisPrompt(request ports.PurchaseAnalysisRequest) string {
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
}`,
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
func (p *PurchaseUseCase) buildAlternativesPrompt(request ports.PurchaseAnalysisRequest) string {
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
]`,
		request.ItemName,
		request.Amount,
		request.UserFinancialProfile.MonthlyIncome,
	)
}
