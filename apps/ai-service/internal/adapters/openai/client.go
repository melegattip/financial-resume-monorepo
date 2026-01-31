package openai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/financial-ai-service/internal/core/ports"
	"github.com/sashabaranov/go-openai"
)

// Client implementa el adaptador de OpenAI
type Client struct {
	client  *openai.Client
	useMock bool
}

// NewClient crea un nuevo cliente de OpenAI
func NewClient(apiKey string, useMock bool) ports.OpenAIClient {
	var client *openai.Client

	if !useMock && apiKey != "" {
		client = openai.NewClient(apiKey)
		var maskedKey string
		if len(apiKey) > 20 {
			maskedKey = apiKey[:10] + "..." + apiKey[len(apiKey)-10:]
		} else if len(apiKey) > 6 {
			maskedKey = apiKey[:3] + "..." + apiKey[len(apiKey)-3:]
		} else {
			maskedKey = "***"
		}
		log.Printf("✅ OpenAI client initialized with API key: %s", maskedKey)
	} else {
		log.Println("🎭 OpenAI client initialized in MOCK mode")
		useMock = true
	}

	return &Client{
		client:  client,
		useMock: useMock,
	}
}

// GenerateCompletion genera una respuesta simple de OpenAI
func (c *Client) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
	if c.useMock {
		return c.getMockCompletion(prompt), nil
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := c.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   500,
	})

	if err != nil {
		log.Printf("❌ Error calling OpenAI completion: %v", err)
		return "", fmt.Errorf("error generating completion: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateAnalysis genera un análisis especializado con prompts del sistema
func (c *Client) GenerateAnalysis(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	if c.useMock {
		return c.getMockAnalysis(systemPrompt, userPrompt), nil
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := c.client.CreateChatCompletion(ctxWithTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   800,
	})

	if err != nil {
		log.Printf("❌ Error calling OpenAI analysis: %v", err)
		return "", fmt.Errorf("error generating analysis: %w", err)
	}

	content := resp.Choices[0].Message.Content
	// Limpiar respuesta JSON si viene envuelta en markdown
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	return content, nil
}

// getMockCompletion genera una respuesta mock simple
func (c *Client) getMockCompletion(prompt string) string {
	preview := prompt
	if len(prompt) > 50 {
		preview = prompt[:50] + "..."
	}
	log.Printf("🎭 Generating mock completion for prompt: %s", preview)
	return `{
		"status": "success",
		"message": "Mock response generated successfully",
		"timestamp": "` + time.Now().Format(time.RFC3339) + `"
	}`
}

// getMockAnalysis genera un análisis mock basado en el tipo de prompt
func (c *Client) getMockAnalysis(systemPrompt, userPrompt string) string {
	preview := systemPrompt
	if len(systemPrompt) > 50 {
		preview = systemPrompt[:50] + "..."
	}
	log.Printf("🎭 Generating mock analysis for system: %s", preview)

	// Determinar tipo de análisis basado en el prompt del sistema
	if strings.Contains(systemPrompt, "asesor financiero") {
		return c.getMockFinancialAnalysis()
	} else if strings.Contains(systemPrompt, "compra") {
		return c.getMockPurchaseAnalysis()
	} else if strings.Contains(systemPrompt, "crédito") {
		return c.getMockCreditAnalysis()
	}

	return c.getMockFinancialAnalysis()
}

// getMockFinancialAnalysis genera un análisis financiero mock
func (c *Client) getMockFinancialAnalysis() string {
	return `{
		"score": 750,
		"level": "Bueno",
		"message": "Tu salud financiera es sólida con algunas áreas de mejora",
		"insights": [
			{
				"title": "Excelente tasa de ahorro",
				"description": "Tu tasa de ahorro del 25% está por encima del promedio",
				"impact": "Positivo",
				"score": 85,
				"action_type": "maintain",
				"category": "savings"
			},
			{
				"title": "Oportunidad de optimización",
				"description": "Podrías reducir gastos en entretenimiento",
				"impact": "Medio",
				"score": 65,
				"action_type": "optimize",
				"category": "expenses"
			}
		]
	}`
}

// getMockPurchaseAnalysis genera un análisis de compra mock
func (c *Client) getMockPurchaseAnalysis() string {
	return `{
		"can_buy": true,
		"confidence": 0.8,
		"reasoning": "Basado en tu situación financiera actual, puedes realizar esta compra sin comprometer tu estabilidad",
		"alternatives": ["Buscar ofertas", "Considerar modelo anterior", "Esperar promociones"],
		"impact_score": 25
	}`
}

// getMockCreditAnalysis genera un análisis crediticio mock
func (c *Client) getMockCreditAnalysis() string {
	return `{
		"current_score": 720,
		"target_score": 800,
		"timeline_months": 12,
		"actions": [
			{
				"title": "Aumentar tasa de ahorro",
				"description": "Incrementar ahorro mensual en 5%",
				"priority": "alta",
				"timeline": "3 meses",
				"impact": 30,
				"difficulty": "media"
			},
			{
				"title": "Diversificar ingresos",
				"description": "Explorar fuentes de ingresos adicionales",
				"priority": "media",
				"timeline": "6 meses",
				"impact": 25,
				"difficulty": "alta"
			}
		],
		"key_metrics": {
			"savings_rate_improvement": 0.05,
			"debt_reduction_target": 0.15,
			"emergency_fund_months": 6
		}
	}`
}
