package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// openAIURL is the OpenAI chat completions endpoint.
// Declared as a var (not const) so that tests can redirect requests to a local httptest.Server.
var openAIURL = "https://api.openai.com/v1/chat/completions"

// OpenAIClient is a lightweight OpenAI client using only stdlib net/http.
type OpenAIClient struct {
	apiKey     string
	useMock    bool
	httpClient *http.Client
}

// NewOpenAIClient creates a new OpenAI client.
// If apiKey is empty, the client will return mock responses instead of calling the API.
func NewOpenAIClient(apiKey string) *OpenAIClient {
	useMock := apiKey == ""
	return &OpenAIClient{
		apiKey:  apiKey,
		useMock: useMock,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GenerateAnalysis sends a chat completion request to OpenAI and returns the content string.
// Falls back to a realistic mock response when OPENAI_API_KEY is not set.
func (c *OpenAIClient) GenerateAnalysis(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	if c.useMock {
		return c.getMockResponse(userPrompt), nil
	}

	payload := chatRequest{
		Model: "gpt-4.1",
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
		MaxTokens:   2000,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal openai request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create openai request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read openai response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openai returned %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse openai response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in openai response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// getMockResponse returns a realistic hardcoded JSON response for testing without a real API key.
// The mock response adapts slightly based on keywords in the prompt.
func (c *OpenAIClient) getMockResponse(userPrompt string) string {
	// Return education cards mock.
	if containsKeyword(userPrompt, "tarjetas educativas") {
		return `{
  "cards": [
    {
      "topic": "emergencia",
      "title": "Tu fondo de emergencia: el primer paso",
      "summary": "Tener 3-6 meses de gastos guardados te protege ante imprevistos. Es la base de cualquier plan financiero sólido. Sin este colchón, cualquier gasto inesperado puede desestabilizar tu situación.",
      "key_concept": "3-6 meses de gastos = tranquilidad financiera",
      "cta": "Crear meta de emergencia",
      "deep_link": "/savings-goals",
      "difficulty": "básico"
    },
    {
      "topic": "presupuesto",
      "title": "El presupuesto 50/30/20 simplificado",
      "summary": "Dividí tus ingresos en 3 partes: 50% para necesidades, 30% para deseos y 20% para ahorro e inversión. Es el método más usado en Latinoamérica por su simplicidad. Empezá con una categoría a la vez.",
      "key_concept": "50% necesidades · 30% deseos · 20% ahorro",
      "cta": "Crear mi presupuesto",
      "deep_link": "/budgets",
      "difficulty": "básico"
    },
    {
      "topic": "ahorro",
      "title": "Págate primero: automatizá tu ahorro",
      "summary": "Transferí a tu cuenta de ahorro el día que cobrás, antes de gastar. Así el ahorro deja de ser 'lo que sobra' y se convierte en una prioridad. Incluso el 5% de tus ingresos hace la diferencia.",
      "key_concept": "Primero ahorrá, después gastá el resto",
      "cta": "Ver mis metas de ahorro",
      "deep_link": "/savings-goals",
      "difficulty": "básico"
    }
  ]
}`
	}

	// Return monthly coaching report mock.
	if containsKeyword(userPrompt, "reporte de coaching") {
		return `{
  "sentiment": "neutral",
  "summary": "Tu situación financiera este mes muestra un balance estable. Hay oportunidades de mejora en el control de gastos variables, pero también logros importantes que vale la pena reconocer.",
  "wins": [
    {
      "title": "Registraste tus transacciones",
      "description": "Mantener el registro de ingresos y gastos es el hábito más valioso. Eso ya te pone adelante de la mayoría."
    },
    {
      "title": "Ingresos estables",
      "description": "Tu flujo de ingresos se mantuvo consistente este mes, lo que te da una base sólida para planificar."
    }
  ],
  "improvements": [
    {
      "title": "Revisá tus gastos por categoría",
      "description": "Identificá las 2-3 categorías donde más gastás y evaluá si hay margen para reducir sin afectar tu calidad de vida."
    },
    {
      "title": "Configurá un presupuesto mensual",
      "description": "Tener límites por categoría te ayuda a tomar decisiones más conscientes antes de gastar."
    }
  ],
  "actions": [
    {
      "title": "Revisá tus gastos de este mes",
      "detail": "Entrá a la sección de gastos y filtrá por categoría para identificar dónde podés optimizar.",
      "deep_link": "/expenses"
    },
    {
      "title": "Creá un presupuesto",
      "detail": "Configurá límites mensuales para tus categorías principales. 15 minutos que cambian tus hábitos.",
      "deep_link": "/budgets"
    },
    {
      "title": "Revisá tus metas de ahorro",
      "detail": "Asegurate de que tus metas estén activas y hacé un depósito aunque sea pequeño.",
      "deep_link": "/savings-goals"
    }
  ],
  "behavior_note": "Seguís usando la app regularmente, lo que demuestra compromiso con tu salud financiera. El próximo nivel es pasar de registrar a planificar activamente."
}`
	}

	// Return credit-score-specific mock when asked for just a score.
	if containsKeyword(userPrompt, "score crediticio") {
		return `{"score": 720}`
	}

	// Return credit-plan mock.
	if containsKeyword(userPrompt, "plan de mejora crediticia") {
		return `{
  "current_score": 650,
  "target_score": 800,
  "timeline_months": 12,
  "actions": [
    {
      "title": "Reducir gastos discrecionales",
      "description": "Identificar y eliminar gastos innecesarios para aumentar el margen de ahorro mensual.",
      "priority": "alta",
      "timeline": "1-3 meses",
      "impact": 40,
      "difficulty": "media"
    },
    {
      "title": "Construir fondo de emergencia",
      "description": "Acumular al menos 3 meses de gastos en una cuenta de ahorro de fácil acceso.",
      "priority": "alta",
      "timeline": "3-6 meses",
      "impact": 35,
      "difficulty": "media"
    },
    {
      "title": "Diversificar fuentes de ingreso",
      "description": "Explorar ingresos adicionales para mejorar la estabilidad financiera.",
      "priority": "media",
      "timeline": "6-12 meses",
      "impact": 25,
      "difficulty": "difícil"
    }
  ],
  "key_metrics": {
    "savings_rate_improvement": 0.05,
    "debt_reduction_target": 0.15,
    "emergency_fund_months": 3
  }
}`
	}

	// Return purchase-decision mock.
	if containsKeyword(userPrompt, "can_buy") || containsKeyword(userPrompt, "compra") {
		return `{
  "can_buy": true,
  "confidence": 0.72,
  "reasoning": "Dado tu balance actual y tu tasa de ahorro, puedes realizar esta compra sin comprometer tu estabilidad financiera. Sin embargo, considera si es el mejor momento o si podrías esperar para ahorrar más.",
  "alternatives": [
    "Buscar el producto en oferta o reacondicionado",
    "Dividir el pago en cuotas sin interés si está disponible",
    "Esperar 30 días para evaluar si sigue siendo necesario"
  ],
  "impact_score": 35
}`
	}

	// Return insights mock.
	if containsKeyword(userPrompt, "insights") {
		return `[
  {
    "title": "Tasa de ahorro saludable",
    "description": "Tu tasa de ahorro actual está por encima del promedio. Mantén este ritmo para alcanzar tus metas financieras antes.",
    "impact": "Positivo",
    "score": 82,
    "action_type": "maintain",
    "category": "savings"
  },
  {
    "title": "Gastos en entretenimiento elevados",
    "description": "Los gastos en entretenimiento representan un porcentaje alto de tu presupuesto mensual. Reducirlos un 20% liberaría fondos para inversión.",
    "impact": "Negativo",
    "score": 55,
    "action_type": "improve",
    "category": "expenses"
  },
  {
    "title": "Estabilidad de ingresos",
    "description": "Tu fuente de ingresos muestra estabilidad consistente. Esto es un factor clave para acceder a mejores productos financieros.",
    "impact": "Positivo",
    "score": 78,
    "action_type": "maintain",
    "category": "income"
  }
]`
	}

	// Default: return health analysis mock.
	return `{
  "score": 750,
  "level": "Bueno",
  "message": "Tu salud financiera es sólida. Tienes buenos hábitos de ahorro y estabilidad en tus ingresos. Con algunos ajustes podrías alcanzar el nivel Excelente.",
  "insights": [
    {
      "title": "Buen nivel de ahorro",
      "description": "Tu tasa de ahorro está por encima del promedio del 20%. Continúa con este hábito y considera invertir el excedente.",
      "impact": "Positivo",
      "score": 85,
      "action_type": "maintain",
      "category": "savings"
    },
    {
      "title": "Optimiza tus gastos variables",
      "description": "Los gastos en categorías discrecionales como entretenimiento y restaurantes pueden reducirse para mejorar tu margen financiero.",
      "impact": "Neutro",
      "score": 60,
      "action_type": "optimize",
      "category": "expenses"
    },
    {
      "title": "Construye tu fondo de emergencia",
      "description": "Asegúrate de tener al menos 6 meses de gastos cubiertos en un fondo de emergencia de fácil acceso.",
      "impact": "Neutro",
      "score": 65,
      "action_type": "improve",
      "category": "savings"
    }
  ]
}`
}

// containsKeyword is a simple helper to detect keywords in the prompt for mock routing.
func containsKeyword(s, keyword string) bool {
	return len(s) > 0 && len(keyword) > 0 && func() bool {
		for i := 0; i <= len(s)-len(keyword); i++ {
			if s[i:i+len(keyword)] == keyword {
				return true
			}
		}
		return false
	}()
}
