package adapters

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/financial-ai-service/internal/adapters/openai"
	"github.com/financial-ai-service/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// OpenAIClientTestSuite suite de tests para el cliente OpenAI
type OpenAIClientTestSuite struct {
	suite.Suite
	client     ports.OpenAIClient
	mockClient ports.OpenAIClient
	ctx        context.Context
}

// SetupTest configura cada test individual
func (suite *OpenAIClientTestSuite) SetupTest() {
	// Cliente con mock habilitado
	suite.mockClient = openai.NewClient("", true)

	// Cliente con API key (para tests específicos)
	suite.client = openai.NewClient("test-api-key", false)

	suite.ctx = context.Background()
}

// TestNewClient_WithMock testa la creación del cliente con mock
func (suite *OpenAIClientTestSuite) TestNewClient_WithMock() {
	// Act
	client := openai.NewClient("", true)

	// Assert
	assert.NotNil(suite.T(), client)

	// Verificar que funciona con mock
	result, err := client.GenerateCompletion(suite.ctx, "test prompt")
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
	assert.Contains(suite.T(), result, "Mock response")
}

// TestNewClient_WithAPIKey testa la creación del cliente con API key
func (suite *OpenAIClientTestSuite) TestNewClient_WithAPIKey() {
	// Act
	client := openai.NewClient("sk-test123456789", false)

	// Assert
	assert.NotNil(suite.T(), client)
}

// TestNewClient_WithEmptyAPIKey testa la creación del cliente con API key vacía
func (suite *OpenAIClientTestSuite) TestNewClient_WithEmptyAPIKey() {
	// Act
	client := openai.NewClient("", false)

	// Assert
	assert.NotNil(suite.T(), client)

	// Debería funcionar en modo mock automáticamente
	result, err := client.GenerateCompletion(suite.ctx, "test prompt")
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

// TestGenerateCompletion_MockMode testa la generación de completions en modo mock
func (suite *OpenAIClientTestSuite) TestGenerateCompletion_MockMode() {
	testCases := []struct {
		name   string
		prompt string
	}{
		{
			name:   "Prompt simple",
			prompt: "¿Cómo está mi salud financiera?",
		},
		{
			name:   "Prompt complejo",
			prompt: "Analiza mi situación financiera con ingresos de $5,000,000 y gastos de $3,500,000",
		},
		{
			name:   "Prompt vacío",
			prompt: "",
		},
		{
			name:   "Prompt muy largo",
			prompt: strings.Repeat("test prompt ", 100),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Act
			result, err := suite.mockClient.GenerateCompletion(suite.ctx, tc.prompt)

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotEmpty(suite.T(), result)
			assert.Contains(suite.T(), result, "Mock response")
			assert.Contains(suite.T(), result, "status")
			assert.Contains(suite.T(), result, "timestamp")
		})
	}
}

// TestGenerateAnalysis_MockMode testa la generación de análisis en modo mock
func (suite *OpenAIClientTestSuite) TestGenerateAnalysis_MockMode() {
	testCases := []struct {
		name         string
		systemPrompt string
		userPrompt   string
		expectedType string
	}{
		{
			name:         "Análisis financiero",
			systemPrompt: "Eres un asesor financiero experto",
			userPrompt:   "Analiza mi salud financiera",
			expectedType: "financial",
		},
		{
			name:         "Análisis de compra",
			systemPrompt: "Analiza decisiones de compra",
			userPrompt:   "¿Puedo comprar un MacBook?",
			expectedType: "purchase",
		},
		{
			name:         "Análisis de crédito",
			systemPrompt: "Analiza situación crediticia",
			userPrompt:   "¿Cómo mejorar mi score?",
			expectedType: "credit",
		},
		{
			name:         "Análisis genérico",
			systemPrompt: "Eres un asistente",
			userPrompt:   "Ayúdame con mis finanzas",
			expectedType: "financial",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Act
			result, err := suite.mockClient.GenerateAnalysis(suite.ctx, tc.systemPrompt, tc.userPrompt)

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotEmpty(suite.T(), result)

			// Verificar que es un JSON válido
			assert.True(suite.T(), strings.HasPrefix(result, "{"))
			assert.True(suite.T(), strings.HasSuffix(result, "}"))

			// Verificar contenido específico basado en el tipo
			switch tc.expectedType {
			case "financial":
				assert.Contains(suite.T(), result, "score")
				assert.Contains(suite.T(), result, "insights")
			case "purchase":
				assert.Contains(suite.T(), result, "can_buy")
				assert.Contains(suite.T(), result, "confidence")
			case "credit":
				assert.Contains(suite.T(), result, "score")
			}
		})
	}
}

// TestContextCancellation testa el manejo de cancelación de contexto
func (suite *OpenAIClientTestSuite) TestContextCancellation() {
	// Arrange
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	// Act
	result, err := suite.mockClient.GenerateCompletion(cancelledCtx, "test prompt")

	// Assert
	// En modo mock, no debería fallar por cancelación de contexto
	// ya que no hace llamadas reales
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

// TestContextWithTimeout testa el manejo de timeout
func (suite *OpenAIClientTestSuite) TestContextWithTimeout() {
	// Arrange
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Act
	result, err := suite.mockClient.GenerateCompletion(ctxWithTimeout, "test prompt")

	// Assert
	// En modo mock, debería funcionar independientemente del timeout
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

// TestAnalysisJSONCleanup testa la limpieza de JSON en markdown
func (suite *OpenAIClientTestSuite) TestAnalysisJSONCleanup() {
	// Este test simula lo que haría el cliente real al limpiar markdown
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "JSON con markdown",
			input:    "```json\n{\"test\": \"value\"}\n```",
			expected: "{\"test\": \"value\"}",
		},
		{
			name:     "JSON sin markdown",
			input:    "{\"test\": \"value\"}",
			expected: "{\"test\": \"value\"}",
		},
		{
			name:     "JSON con espacios",
			input:    "  ```json  \n  {\"test\": \"value\"}  \n  ```  ",
			expected: "{\"test\": \"value\"}",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Simular la limpieza que hace el cliente (más completa)
			result := strings.TrimSpace(tc.input)

			// Remover bloques de código markdown
			if strings.Contains(result, "```json") {
				// Encontrar el inicio del bloque
				start := strings.Index(result, "```json")
				if start != -1 {
					result = result[start+7:] // Remover "```json"
				}

				// Encontrar el final del bloque
				end := strings.LastIndex(result, "```")
				if end != -1 {
					result = result[:end] // Remover "```"
				}
			}
			result = strings.TrimSpace(result)

			// Assert
			assert.Equal(suite.T(), tc.expected, result)
		})
	}
}

// TestMockResponseGeneration testa diferentes tipos de respuestas mock
func (suite *OpenAIClientTestSuite) TestMockResponseGeneration() {
	testCases := []struct {
		name            string
		systemPrompt    string
		expectedContent []string
	}{
		{
			name:         "Análisis financiero",
			systemPrompt: "asesor financiero",
			expectedContent: []string{
				"score",
				"level",
				"insights",
				"Bueno",
			},
		},
		{
			name:         "Análisis de compra",
			systemPrompt: "compra",
			expectedContent: []string{
				"can_buy",
				"confidence",
				"reasoning",
			},
		},
		{
			name:         "Análisis de crédito",
			systemPrompt: "crédito",
			expectedContent: []string{
				"current_score",
				"target_score",
				"actions",
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Act
			result, err := suite.mockClient.GenerateAnalysis(suite.ctx, tc.systemPrompt, "test user prompt")

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotEmpty(suite.T(), result)

			for _, expectedContent := range tc.expectedContent {
				assert.Contains(suite.T(), result, expectedContent)
			}
		})
	}
}

// TestConcurrentRequests testa solicitudes concurrentes
func (suite *OpenAIClientTestSuite) TestConcurrentRequests() {
	// Arrange
	numRequests := 10
	resultsChan := make(chan string, numRequests)
	errorsChan := make(chan error, numRequests)

	// Act
	for i := 0; i < numRequests; i++ {
		go func(index int) {
			result, err := suite.mockClient.GenerateCompletion(suite.ctx, fmt.Sprintf("prompt %d", index))
			if err != nil {
				errorsChan <- err
			} else {
				resultsChan <- result
			}
		}(i)
	}

	// Assert
	for i := 0; i < numRequests; i++ {
		select {
		case result := <-resultsChan:
			assert.NotEmpty(suite.T(), result)
		case err := <-errorsChan:
			assert.NoError(suite.T(), err) // No deberían haber errores
		case <-time.After(5 * time.Second):
			suite.T().Fatal("Timeout waiting for concurrent requests")
		}
	}
}

// TestLargePrompts testa prompts muy largos
func (suite *OpenAIClientTestSuite) TestLargePrompts() {
	// Arrange
	largePrompt := strings.Repeat("Esta es una consulta financiera muy larga. ", 1000)

	// Act
	result, err := suite.mockClient.GenerateCompletion(suite.ctx, largePrompt)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

// TestSpecialCharacters testa prompts con caracteres especiales
func (suite *OpenAIClientTestSuite) TestSpecialCharacters() {
	// Arrange
	specialPrompt := "Análisis financiero con símbolos: $5,000,000 💰 €1,000 ¥10,000 £500 🏦 📊 📈"

	// Act
	result, err := suite.mockClient.GenerateCompletion(suite.ctx, specialPrompt)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
}

// TestRunSuite ejecuta todos los tests del suite
func TestOpenAIClientTestSuite(t *testing.T) {
	suite.Run(t, new(OpenAIClientTestSuite))
}

// Tests adicionales para casos específicos

// TestClientInitializationWithDifferentKeys testa inicialización con diferentes tipos de API keys
func TestClientInitializationWithDifferentKeys(t *testing.T) {
	testCases := []struct {
		name       string
		apiKey     string
		useMock    bool
		shouldWork bool
	}{
		{
			name:       "API key válida con mock disabled",
			apiKey:     "sk-1234567890abcdef",
			useMock:    false,
			shouldWork: false, // API key falsa, esperamos error
		},
		{
			name:       "API key vacía con mock enabled",
			apiKey:     "",
			useMock:    true,
			shouldWork: true,
		},
		{
			name:       "API key corta",
			apiKey:     "sk-123",
			useMock:    false,
			shouldWork: false, // API key falsa, esperamos error
		},
		{
			name:       "Sin API key, mock disabled",
			apiKey:     "",
			useMock:    false,
			shouldWork: true, // Debería funcionar automáticamente en modo mock
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			client := openai.NewClient(tc.apiKey, tc.useMock)

			// Assert
			assert.NotNil(t, client)

			// Intentar generar completion para verificar funcionamiento
			result, err := client.GenerateCompletion(context.Background(), "test")

			if tc.shouldWork {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			} else {
				// Esperamos error para API keys falsas
				assert.Error(t, err)
				assert.Empty(t, result)
			}
		})
	}
}
