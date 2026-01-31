package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/ports"
)

// AIServiceProxy maneja la comunicación HTTP con el microservicio AI
type AIServiceProxy struct {
	aiServiceURL string
	httpClient   *http.Client
}

// NewAIServiceProxy crea una nueva instancia del proxy
func NewAIServiceProxy(serviceURL string) *AIServiceProxy {
	return &AIServiceProxy{
		aiServiceURL: serviceURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AnalyzeFinancialHealth implementa ports.AIAnalysisPort
func (p *AIServiceProxy) AnalyzeFinancialHealth(ctx context.Context, data ports.FinancialAnalysisData) (*ports.HealthAnalysis, error) {
	// Preparar request payload
	requestPayload := map[string]interface{}{
		"user_id":              data.UserID,
		"total_income":         data.TotalIncome,
		"total_expenses":       data.TotalExpenses,
		"savings_rate":         data.SavingsRate,
		"income_stability":     data.IncomeStability,
		"financial_score":      data.FinancialScore,
		"period":               data.Period,
		"expenses_by_category": data.ExpensesByCategory,
	}

	// Hacer petición HTTP
	response, err := p.makeRequest(ctx, "POST", "/api/v1/ai/health-analysis", requestPayload)
	if err != nil {
		return nil, fmt.Errorf("error llamando al AI Service: %w", err)
	}

	// Parsear respuesta
	var aiResponse struct {
		Data struct {
			Score       int               `json:"score"`
			Level       string            `json:"level"`
			Message     string            `json:"message"`
			Insights    []ports.AIInsight `json:"insights"`
			GeneratedAt time.Time         `json:"generated_at"`
		} `json:"data"`
		Success bool   `json:"success"`
		Source  string `json:"source"`
	}

	if err := json.Unmarshal(response, &aiResponse); err != nil {
		return nil, fmt.Errorf("error parseando respuesta del AI Service: %w", err)
	}

	if !aiResponse.Success {
		return nil, fmt.Errorf("AI Service retornó error")
	}

	return &ports.HealthAnalysis{
		Score:       aiResponse.Data.Score,
		Level:       aiResponse.Data.Level,
		Message:     aiResponse.Data.Message,
		Insights:    aiResponse.Data.Insights,
		GeneratedAt: aiResponse.Data.GeneratedAt,
	}, nil
}

// GenerateInsights implementa ports.AIAnalysisPort
func (p *AIServiceProxy) GenerateInsights(ctx context.Context, data ports.FinancialAnalysisData) ([]ports.AIInsight, error) {
	// Reutilizar el análisis de salud financiera
	healthAnalysis, err := p.AnalyzeFinancialHealth(ctx, data)
	if err != nil {
		return nil, err
	}

	return healthAnalysis.Insights, nil
}

// CanIBuy implementa ports.PurchaseDecisionPort
func (p *AIServiceProxy) CanIBuy(ctx context.Context, request ports.PurchaseAnalysisRequest) (*ports.PurchaseDecision, error) {
	// Preparar request payload
	requestPayload := map[string]interface{}{
		"user_id":       request.UserID,
		"item_name":     request.ItemName,
		"amount":        request.Amount,
		"description":   request.Description,
		"is_necessary":  request.IsNecessary,
		"payment_types": request.PaymentTypes,
		"user_financial_profile": map[string]interface{}{
			"current_balance":      request.UserFinancialProfile.CurrentBalance,
			"monthly_income":       request.UserFinancialProfile.MonthlyIncome,
			"monthly_expenses":     request.UserFinancialProfile.MonthlyExpenses,
			"savings_rate":         request.UserFinancialProfile.SavingsRate,
			"income_stability":     request.UserFinancialProfile.IncomeStability,
			"financial_discipline": request.UserFinancialProfile.FinancialDiscipline,
		},
	}

	// Hacer petición HTTP
	response, err := p.makeRequest(ctx, "POST", "/api/v1/ai/can-i-buy", requestPayload)
	if err != nil {
		return nil, fmt.Errorf("error llamando al AI Service: %w", err)
	}

	// Parsear respuesta
	var aiResponse struct {
		Data struct {
			CanBuy       bool     `json:"can_buy"`
			Confidence   float64  `json:"confidence"`
			Reasoning    string   `json:"reasoning"`
			Alternatives []string `json:"alternatives"`
			ImpactScore  int      `json:"impact_score"`
		} `json:"data"`
		Success bool   `json:"success"`
		Source  string `json:"source"`
	}

	if err := json.Unmarshal(response, &aiResponse); err != nil {
		return nil, fmt.Errorf("error parseando respuesta del AI Service: %w", err)
	}

	if !aiResponse.Success {
		return nil, fmt.Errorf("AI Service retornó error")
	}

	return &ports.PurchaseDecision{
		CanBuy:       aiResponse.Data.CanBuy,
		Confidence:   aiResponse.Data.Confidence,
		Reasoning:    aiResponse.Data.Reasoning,
		Alternatives: aiResponse.Data.Alternatives,
		ImpactScore:  aiResponse.Data.ImpactScore,
	}, nil
}

// GenerateImprovementPlan implementa ports.CreditAnalysisPort
func (p *AIServiceProxy) GenerateImprovementPlan(ctx context.Context, data ports.FinancialAnalysisData) (*ports.CreditPlan, error) {
	// Preparar request payload
	requestPayload := map[string]interface{}{
		"user_id":              data.UserID,
		"total_income":         data.TotalIncome,
		"total_expenses":       data.TotalExpenses,
		"savings_rate":         data.SavingsRate,
		"income_stability":     data.IncomeStability,
		"financial_score":      data.FinancialScore,
		"period":               data.Period,
		"expenses_by_category": data.ExpensesByCategory,
	}

	// Hacer petición HTTP
	response, err := p.makeRequest(ctx, "POST", "/api/v1/ai/credit-plan", requestPayload)
	if err != nil {
		return nil, fmt.Errorf("error llamando al AI Service: %w", err)
	}

	// Parsear respuesta
	var aiResponse struct {
		Data struct {
			CurrentScore   int                    `json:"current_score"`
			TargetScore    int                    `json:"target_score"`
			TimelineMonths int                    `json:"timeline_months"`
			Actions        []ports.CreditAction   `json:"actions"`
			KeyMetrics     map[string]interface{} `json:"key_metrics"`
		} `json:"data"`
		Success bool   `json:"success"`
		Source  string `json:"source"`
	}

	if err := json.Unmarshal(response, &aiResponse); err != nil {
		return nil, fmt.Errorf("error parseando respuesta del AI Service: %w", err)
	}

	if !aiResponse.Success {
		return nil, fmt.Errorf("AI Service retornó error")
	}

	return &ports.CreditPlan{
		CurrentScore:   aiResponse.Data.CurrentScore,
		TargetScore:    aiResponse.Data.TargetScore,
		TimelineMonths: aiResponse.Data.TimelineMonths,
		Actions:        aiResponse.Data.Actions,
		KeyMetrics:     aiResponse.Data.KeyMetrics,
	}, nil
}

// SuggestAlternatives implementa ports.PurchaseDecisionPort
func (p *AIServiceProxy) SuggestAlternatives(ctx context.Context, request ports.PurchaseAnalysisRequest) ([]ports.Alternative, error) {
	// Preparar request payload
	requestPayload := map[string]interface{}{
		"user_id":       request.UserID,
		"item_name":     request.ItemName,
		"amount":        request.Amount,
		"description":   request.Description,
		"is_necessary":  request.IsNecessary,
		"payment_types": request.PaymentTypes,
		"user_financial_profile": map[string]interface{}{
			"current_balance":      request.UserFinancialProfile.CurrentBalance,
			"monthly_income":       request.UserFinancialProfile.MonthlyIncome,
			"monthly_expenses":     request.UserFinancialProfile.MonthlyExpenses,
			"savings_rate":         request.UserFinancialProfile.SavingsRate,
			"income_stability":     request.UserFinancialProfile.IncomeStability,
			"financial_discipline": request.UserFinancialProfile.FinancialDiscipline,
		},
	}

	// Hacer petición HTTP
	response, err := p.makeRequest(ctx, "POST", "/api/v1/ai/alternatives", requestPayload)
	if err != nil {
		return nil, fmt.Errorf("error llamando al AI Service: %w", err)
	}

	// Parsear respuesta
	var aiResponse struct {
		Data struct {
			Alternatives []ports.Alternative `json:"alternatives"`
		} `json:"data"`
		Success bool   `json:"success"`
		Source  string `json:"source"`
	}

	if err := json.Unmarshal(response, &aiResponse); err != nil {
		return nil, fmt.Errorf("error parseando respuesta del AI Service: %w", err)
	}

	if !aiResponse.Success {
		return nil, fmt.Errorf("AI Service retornó error")
	}

	return aiResponse.Data.Alternatives, nil
}

// CalculateCreditScore implementa ports.CreditAnalysisPort
func (p *AIServiceProxy) CalculateCreditScore(ctx context.Context, data ports.FinancialAnalysisData) (int, error) {
	// Preparar request payload
	requestPayload := map[string]interface{}{
		"user_id":              data.UserID,
		"total_income":         data.TotalIncome,
		"total_expenses":       data.TotalExpenses,
		"savings_rate":         data.SavingsRate,
		"income_stability":     data.IncomeStability,
		"financial_score":      data.FinancialScore,
		"period":               data.Period,
		"expenses_by_category": data.ExpensesByCategory,
	}

	// Hacer petición HTTP
	response, err := p.makeRequest(ctx, "POST", "/api/v1/ai/credit-score", requestPayload)
	if err != nil {
		return 0, fmt.Errorf("error llamando al AI Service: %w", err)
	}

	// Parsear respuesta
	var aiResponse struct {
		Data struct {
			Score int `json:"score"`
		} `json:"data"`
		Success bool   `json:"success"`
		Source  string `json:"source"`
	}

	if err := json.Unmarshal(response, &aiResponse); err != nil {
		return 0, fmt.Errorf("error parseando respuesta del AI Service: %w", err)
	}

	if !aiResponse.Success {
		return 0, fmt.Errorf("AI Service retornó error")
	}

	return aiResponse.Data.Score, nil
}

// makeRequest es un método auxiliar para hacer peticiones HTTP
func (p *AIServiceProxy) makeRequest(ctx context.Context, method, endpoint string, payload interface{}) ([]byte, error) {
	// Serializar payload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %w", err)
	}

	// Crear request
	url := p.aiServiceURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	// Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Ejecutar request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando request: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI Service retornó status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// HealthCheck verifica que el AI Service esté disponible
func (p *AIServiceProxy) HealthCheck() error {
	resp, err := p.httpClient.Get(p.aiServiceURL + "/health")
	if err != nil {
		return fmt.Errorf("error conectando con AI Service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI Service no está saludable: status %d", resp.StatusCode)
	}

	return nil
}
