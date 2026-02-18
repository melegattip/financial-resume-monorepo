package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/domain"
)

// PurchaseService handles AI-powered purchase decision analysis.
type PurchaseService struct {
	openai *OpenAIClient
}

// NewPurchaseService creates a new PurchaseService.
func NewPurchaseService(openai *OpenAIClient) *PurchaseService {
	return &PurchaseService{openai: openai}
}

// CanIBuy analyses whether the user can afford a specific purchase.
func (s *PurchaseService) CanIBuy(ctx context.Context, req domain.PurchaseAnalysisRequest) (*domain.PurchaseDecision, error) {
	systemPrompt := `Eres un asesor financiero experto especializado en análisis de compras inteligentes.
Tu trabajo es analizar la situación financiera del usuario y dar recomendaciones precisas y personalizadas.
Debes considerar el tipo de pago, la necesidad real del artículo, y el impacto a largo plazo.
Sé específico con números y porcentajes.
Responde ÚNICAMENTE con un JSON válido en el formato solicitado.`

	userPrompt := s.buildPurchaseAnalysisPrompt(req)

	raw, err := s.openai.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("error analyzing purchase decision: %w", err)
	}

	raw = cleanJSONResponse(raw)

	var decision domain.PurchaseDecision
	if err := json.Unmarshal([]byte(raw), &decision); err != nil {
		return &domain.PurchaseDecision{
			CanBuy:       false,
			Confidence:   0.5,
			Reasoning:    "Análisis no disponible en este momento. Considera tu balance actual antes de decidir.",
			Alternatives: []string{},
			ImpactScore:  50,
			GeneratedAt:  time.Now(),
		}, nil
	}

	decision.GeneratedAt = time.Now()
	return &decision, nil
}

// SuggestAlternatives generates cheaper or more viable alternatives to a purchase.
func (s *PurchaseService) SuggestAlternatives(ctx context.Context, req domain.PurchaseAnalysisRequest) ([]domain.Alternative, error) {
	systemPrompt := `Eres un asesor financiero experto que genera alternativas de compra más económicas y viables.
Debes proporcionar opciones realistas y específicas que ayuden al usuario a ahorrar dinero.
Responde ÚNICAMENTE con un JSON válido que contenga un array de alternativas.`

	userPrompt := s.buildAlternativesPrompt(req)

	raw, err := s.openai.GenerateAnalysis(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("error generating alternatives: %w", err)
	}

	raw = cleanJSONResponse(raw)

	var alternatives []domain.Alternative
	if err := json.Unmarshal([]byte(raw), &alternatives); err != nil {
		return []domain.Alternative{}, nil
	}
	return alternatives, nil
}

// buildPurchaseAnalysisPrompt builds the structured prompt for purchase analysis.
func (s *PurchaseService) buildPurchaseAnalysisPrompt(req domain.PurchaseAnalysisRequest) string {
	available := req.UserFinancialProfile.MonthlyIncome - req.UserFinancialProfile.MonthlyExpenses

	necessityText := "No es urgente"
	if req.IsNecessary {
		necessityText = "Es una necesidad urgente/esencial"
	}

	descriptionText := ""
	if req.Description != "" {
		descriptionText = fmt.Sprintf("\nDescripción: %s", req.Description)
	}

	paymentTypeText := "No especificado"
	if len(req.PaymentTypes) > 0 {
		var labels []string
		for _, pt := range req.PaymentTypes {
			switch pt {
			case "contado":
				labels = append(labels, "Pago de contado")
			case "cuotas":
				labels = append(labels, "Plan de pagos/cuotas")
			case "ahorro":
				labels = append(labels, "Necesita ahorrar para esto")
			default:
				labels = append(labels, pt)
			}
		}
		paymentTypeText = strings.Join(labels, " + ")
	}

	savingsGoalsInfo := ""
	if len(req.UserFinancialProfile.SavingsGoals) > 0 {
		savingsGoalsInfo = fmt.Sprintf("\nMetas de ahorro activas: %d", len(req.UserFinancialProfile.SavingsGoals))
		for _, goal := range req.UserFinancialProfile.SavingsGoals {
			savingsGoalsInfo += fmt.Sprintf("\n- %s: $%.0f de $%.0f (%.1f%%)",
				goal.Name, goal.CurrentAmount, goal.TargetAmount, goal.Progress*100)
		}
	}

	return fmt.Sprintf(`
Analiza si el usuario puede realizar esta compra:
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
- Disciplina financiera: %d/1000%s

Responde con este JSON:
{
  "can_buy": true|false,
  "confidence": 0.0-1.0,
  "reasoning": "Explicación detallada con números concretos",
  "alternatives": ["alternativa1", "alternativa2", "alternativa3"],
  "impact_score": 1-100
}`,
		req.ItemName,
		req.Amount,
		descriptionText,
		necessityText,
		paymentTypeText,
		req.UserFinancialProfile.MonthlyIncome,
		req.UserFinancialProfile.MonthlyExpenses,
		available,
		req.UserFinancialProfile.CurrentBalance,
		req.UserFinancialProfile.SavingsRate*100,
		req.UserFinancialProfile.FinancialDiscipline,
		savingsGoalsInfo,
	)
}

// buildAlternativesPrompt builds the structured prompt for generating purchase alternatives.
func (s *PurchaseService) buildAlternativesPrompt(req domain.PurchaseAnalysisRequest) string {
	return fmt.Sprintf(`
Genera 3 alternativas más económicas o viables para esta compra:
Artículo: %s
Precio original: $%.0f
Ingresos mensuales del usuario: $%.0f

Responde con un array JSON:
[
  {
    "name": "Nombre de la alternativa",
    "description": "Descripción breve y específica",
    "savings": 0.0,
    "feasibility": "alta|media|baja"
  }
]`,
		req.ItemName,
		req.Amount,
		req.UserFinancialProfile.MonthlyIncome,
	)
}
