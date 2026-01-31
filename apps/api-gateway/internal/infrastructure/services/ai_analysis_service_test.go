package services

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAIAnalysisService(t *testing.T) {
	tests := []struct {
		name            string
		mockEnv         string
		expectedUseMock bool
	}{
		{
			name:            "Should use mock when USE_AI_MOCK is true",
			mockEnv:         "true",
			expectedUseMock: true,
		},
		{
			name:            "Should use real AI when USE_AI_MOCK is false and API key exists",
			mockEnv:         "false",
			expectedUseMock: false, // Nota: Podría cambiar a true si no hay API key
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			t.Setenv("USE_AI_MOCK", tt.mockEnv)
			if tt.mockEnv == "false" {
				t.Setenv("OPENAI_API_KEY", "test-api-key-1234567890abcdef")
			}

			// Execute
			service := NewAIAnalysisService()

			// Assert
			assert.NotNil(t, service)
			assert.Equal(t, tt.expectedUseMock, service.useMock)
		})
	}
}

func TestAIAnalysisService_AnalyzeFinancialHealth_Mock(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewAIAnalysisService()

	testData := ports.FinancialAnalysisData{
		UserID:        "test-user-123",
		TotalIncome:   5000000, // $5M COP
		TotalExpenses: 3500000, // $3.5M COP
		SavingsRate:   0.3,     // 30%
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1500000,
			"Transporte":   1000000,
			"Vivienda":     1000000,
		},
		IncomeStability: 0.8,
		FinancialScore:  750,
		Period:          "2024-01",
	}

	// Execute
	ctx := context.Background()
	result, err := service.AnalyzeFinancialHealth(ctx, testData)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, result.Score, 0)
	assert.LessOrEqual(t, result.Score, 1000)
	assert.NotEmpty(t, result.Level)
	assert.NotEmpty(t, result.Message)
	assert.NotEmpty(t, result.Insights)
	assert.True(t, result.GeneratedAt.After(time.Now().Add(-1*time.Minute)))

	// Verificar que el score se calcula correctamente
	expectedScore := calculateExpectedScore(testData)
	assert.Equal(t, expectedScore, result.Score)
}

func TestAIAnalysisService_GenerateInsights_Mock(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewAIAnalysisService()

	testData := ports.FinancialAnalysisData{
		UserID:        "test-user-456",
		TotalIncome:   3000000,
		TotalExpenses: 2800000,
		SavingsRate:   0.067, // ~6.7%
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1200000,
			"Transporte":   800000,
			"Servicios":    800000,
		},
		IncomeStability: 0.6,
		FinancialScore:  450,
		Period:          "2024-01",
	}

	// Execute
	ctx := context.Background()
	insights, err := service.GenerateInsights(ctx, testData)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, insights)
	assert.Len(t, insights, 3) // Mock devuelve exactamente 3 insights

	// Verificar estructura de insights
	for _, insight := range insights {
		assert.NotEmpty(t, insight.Title)
		assert.NotEmpty(t, insight.Description)
		assert.Contains(t, []string{"high", "medium", "low"}, insight.Impact)
		assert.Greater(t, insight.Score, 0)
		assert.LessOrEqual(t, insight.Score, 1000)
		assert.Contains(t, []string{"save", "optimize", "alert", "invest"}, insight.ActionType)
		assert.NotEmpty(t, insight.Category)
	}
}

func TestAIAnalysisService_CalculateHealthScore(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewAIAnalysisService()

	tests := []struct {
		name        string
		data        ports.FinancialAnalysisData
		expectedMin int
		expectedMax int
		description string
	}{
		{
			name: "Excellent financial health",
			data: ports.FinancialAnalysisData{
				TotalIncome:     5000000,
				TotalExpenses:   2000000, // 40% de gastos
				SavingsRate:     0.6,     // 60% ahorro
				IncomeStability: 0.9,     // 90% estabilidad
			},
			expectedMin: 650,
			expectedMax: 750,
			description: "High savings rate and stability should yield good score",
		},
		{
			name: "Poor financial health",
			data: ports.FinancialAnalysisData{
				TotalIncome:     2000000,
				TotalExpenses:   2200000, // Gastos > ingresos
				SavingsRate:     -0.1,    // Ahorro negativo
				IncomeStability: 0.3,     // Baja estabilidad
			},
			expectedMin: 0,
			expectedMax: 300,
			description: "Negative savings and low stability should yield poor score",
		},
		{
			name: "Average financial health",
			data: ports.FinancialAnalysisData{
				TotalIncome:     5000000, // $5M
				TotalExpenses:   4000000, // $4M
				SavingsRate:     0.15,    // 15% tasa de ahorro
				IncomeStability: 0.7,     // 70% estabilidad
				FinancialScore:  600,     // Score promedio
			},
			expectedMin: 300, // Ajustado según algoritmo real: 0.15*400 + 0.7*300 + (1M/5M)*300 + bonus ≈ 350
			expectedMax: 500, // Rango más realista
			description: "Moderate savings and stability should yield average score",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			score := service.calculateHealthScore(tt.data)

			// Assert
			assert.GreaterOrEqual(t, score, tt.expectedMin, tt.description)
			assert.LessOrEqual(t, score, tt.expectedMax, tt.description)
			assert.GreaterOrEqual(t, score, 0, "Score should never be negative")
			assert.LessOrEqual(t, score, 1000, "Score should never exceed 1000")
		})
	}
}

func TestAIAnalysisService_GetHealthLevelAndMessage(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewAIAnalysisService()

	tests := []struct {
		score         int
		expectedLevel string
		description   string
	}{
		{score: 850, expectedLevel: "Excelente", description: "Score 850 should be Excelente"},
		{score: 750, expectedLevel: "Bueno", description: "Score 750 should be Bueno"},
		{score: 500, expectedLevel: "Regular", description: "Score 500 should be Regular"},
		{score: 300, expectedLevel: "Mejorable", description: "Score 300 should be Mejorable"},
		{score: 800, expectedLevel: "Excelente", description: "Boundary case: 800 should be Excelente"},
		{score: 600, expectedLevel: "Bueno", description: "Boundary case: 600 should be Bueno"},
		{score: 400, expectedLevel: "Regular", description: "Boundary case: 400 should be Regular"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Execute
			level, message := service.getHealthLevelAndMessage(tt.score)

			// Assert
			assert.Equal(t, tt.expectedLevel, level)
			assert.NotEmpty(t, message)
			assert.Contains(t, message, "financiera")
		})
	}
}

func TestAIAnalysisService_BuildInsightsPrompt(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewAIAnalysisService()

	testData := ports.FinancialAnalysisData{
		UserID:        "test-user",
		TotalIncome:   5000000,
		TotalExpenses: 3500000,
		SavingsRate:   0.3,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1500000,
			"Transporte":   1000000,
			"Vivienda":     1000000,
		},
		FinancialScore: 750,
		Period:         "2024-01",
	}

	// Execute
	prompt := service.buildInsightsPrompt(testData)

	// Assert
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "5000000")      // Ingresos
	assert.Contains(t, prompt, "3500000")      // Gastos
	assert.Contains(t, prompt, "30.0")         // Tasa de ahorro en %
	assert.Contains(t, prompt, "750")          // Score financiero
	assert.Contains(t, prompt, "2024-01")      // Período
	assert.Contains(t, prompt, "JSON")         // Formato esperado
	assert.Contains(t, prompt, "Alimentación") // Top categoría
}

// Función auxiliar para calcular el score esperado (replica la lógica del servicio)
func calculateExpectedScore(data ports.FinancialAnalysisData) int {
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

	// Asegurar rango válido
	if score > 1000 {
		score = 1000
	}
	if score < 0 {
		score = 0
	}

	return score
}

// Test de integración básico
func TestAIAnalysisService_Integration_Mock(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewAIAnalysisService()

	testData := ports.FinancialAnalysisData{
		UserID:        "integration-test-user",
		TotalIncome:   4000000,
		TotalExpenses: 3000000,
		SavingsRate:   0.25,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1200000,
			"Transporte":   900000,
			"Vivienda":     900000,
		},
		IncomeStability: 0.8,
		FinancialScore:  650,
		Period:          "2024-01",
	}

	// Execute - Flujo completo
	ctx := context.Background()

	// 1. Generar insights
	insights, err := service.GenerateInsights(ctx, testData)
	require.NoError(t, err)
	assert.Len(t, insights, 3)

	// 2. Analizar salud financiera
	healthAnalysis, err := service.AnalyzeFinancialHealth(ctx, testData)
	require.NoError(t, err)
	assert.NotNil(t, healthAnalysis)

	// 3. Verificar coherencia entre insights y análisis de salud
	assert.Equal(t, len(insights), len(healthAnalysis.Insights))
	assert.Greater(t, healthAnalysis.Score, 0)
	assert.LessOrEqual(t, healthAnalysis.Score, 1000)

	// 4. Verificar que el timestamp es reciente
	assert.True(t, healthAnalysis.GeneratedAt.After(time.Now().Add(-5*time.Second)))
}

// Benchmark para verificar performance
func BenchmarkAIAnalysisService_GenerateInsights_Mock(b *testing.B) {
	// Setup
	b.Setenv("USE_AI_MOCK", "true")
	service := NewAIAnalysisService()

	testData := ports.FinancialAnalysisData{
		UserID:        "benchmark-user",
		TotalIncome:   5000000,
		TotalExpenses: 3500000,
		SavingsRate:   0.3,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1500000,
			"Transporte":   1000000,
			"Vivienda":     1000000,
		},
		IncomeStability: 0.8,
		FinancialScore:  750,
		Period:          "2024-01",
	}

	ctx := context.Background()

	// Execute
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateInsights(ctx, testData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestAIAnalysisService_AnalyzeFinancialHealth_RealAI tests the real AI flow (without actually calling OpenAI)
func TestAIAnalysisService_AnalyzeFinancialHealth_RealAI_ErrorHandling(t *testing.T) {
	// Test error handling when API key is missing
	os.Setenv("USE_AI_MOCK", "false")
	os.Unsetenv("OPENAI_API_KEY")
	defer func() {
		os.Setenv("USE_AI_MOCK", "true")
		os.Setenv("OPENAI_API_KEY", "test-api-key-1234567890abcdef")
	}()

	service := NewAIAnalysisService()

	data := ports.FinancialAnalysisData{
		UserID:        "test-user",
		TotalIncome:   5000000,
		TotalExpenses: 4000000,
		SavingsRate:   0.2,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1500000,
			"Transporte":   1200000,
			"Vivienda":     1300000,
		},
		IncomeStability: 0.8,
		FinancialScore:  750,
		Period:          "2024-01",
	}

	ctx := context.Background()
	result, err := service.AnalyzeFinancialHealth(ctx, data)

	// Should still work with mock fallback or handle gracefully
	if err != nil {
		assert.Contains(t, err.Error(), "API key", "Should mention API key issue")
	} else {
		assert.NotNil(t, result, "Should return valid result even with missing API key")
	}
}

// TestAIAnalysisService_GenerateInsights_ErrorCases tests error handling
func TestAIAnalysisService_GenerateInsights_ErrorCases(t *testing.T) {
	tests := []struct {
		name string
		data ports.FinancialAnalysisData
	}{
		{
			name: "Zero income case",
			data: ports.FinancialAnalysisData{
				UserID:        "test-user-zero",
				TotalIncome:   0,
				TotalExpenses: 1000000,
				SavingsRate:   -1.0,
				ExpensesByCategory: map[string]float64{
					"Alimentación": 500000,
					"Transporte":   300000,
					"Vivienda":     200000,
				},
				IncomeStability: 0.0,
				FinancialScore:  100,
				Period:          "2024-01",
			},
		},
		{
			name: "Very high income case",
			data: ports.FinancialAnalysisData{
				UserID:        "test-user-rich",
				TotalIncome:   100000000,
				TotalExpenses: 5000000,
				SavingsRate:   0.95,
				ExpensesByCategory: map[string]float64{
					"Alimentación": 1000000,
					"Transporte":   2000000,
					"Vivienda":     2000000,
				},
				IncomeStability: 1.0,
				FinancialScore:  950,
				Period:          "2024-01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("USE_AI_MOCK", "true")
			service := NewAIAnalysisService()

			ctx := context.Background()
			insights, err := service.GenerateInsights(ctx, tt.data)

			require.NoError(t, err, "Should handle edge cases gracefully")
			assert.NotEmpty(t, insights, "Should return insights even for edge cases")
			assert.Len(t, insights, 3, "Should always return 3 insights")

			// Verify insights structure
			for _, insight := range insights {
				assert.NotEmpty(t, insight.Title, "Insight should have title")
				assert.NotEmpty(t, insight.Description, "Insight should have description")
				assert.NotEmpty(t, insight.Impact, "Insight should have impact level")
				assert.Greater(t, insight.Score, 0, "Insight should have positive score")
				assert.NotEmpty(t, insight.ActionType, "Insight should have action type")
				assert.NotEmpty(t, insight.Category, "Insight should have category")
			}
		})
	}
}

// TestAIAnalysisService_CalculateHealthScore_EdgeCases tests edge cases for score calculation
func TestAIAnalysisService_CalculateHealthScore_EdgeCases(t *testing.T) {
	os.Setenv("USE_AI_MOCK", "true")
	service := NewAIAnalysisService()

	tests := []struct {
		name        string
		data        ports.FinancialAnalysisData
		expectedMin int
		expectedMax int
		description string
	}{
		{
			name: "Maximum possible score",
			data: ports.FinancialAnalysisData{
				TotalIncome:     10000000,
				TotalExpenses:   1000000, // 10% expenses
				SavingsRate:     0.9,     // 90% savings (will be capped at 400)
				IncomeStability: 1.0,     // Perfect stability
				FinancialScore:  1000,
			},
			expectedMin: 900,
			expectedMax: 1000,
			description: "Perfect financial health should yield maximum score",
		},
		{
			name: "Minimum possible score",
			data: ports.FinancialAnalysisData{
				TotalIncome:     1000000,
				TotalExpenses:   2000000, // Expenses > income
				SavingsRate:     -1.0,    // Negative savings
				IncomeStability: 0.0,     // No stability
				FinancialScore:  0,
			},
			expectedMin: 0,
			expectedMax: 100,
			description: "Worst financial health should yield minimum score",
		},
		{
			name: "Zero income edge case",
			data: ports.FinancialAnalysisData{
				TotalIncome:     0,
				TotalExpenses:   1000000,
				SavingsRate:     0.0,
				IncomeStability: 0.0,
				FinancialScore:  0,
			},
			expectedMin: 0,
			expectedMax: 100,
			description: "Zero income should be handled gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := service.calculateHealthScore(tt.data)
			assert.GreaterOrEqual(t, score, tt.expectedMin, tt.description)
			assert.LessOrEqual(t, score, tt.expectedMax, tt.description)
			assert.GreaterOrEqual(t, score, 0, "Score should never be negative")
			assert.LessOrEqual(t, score, 1000, "Score should never exceed 1000")
		})
	}
}

// TestAIAnalysisService_BuildInsightsPrompt_Coverage tests prompt building
func TestAIAnalysisService_BuildInsightsPrompt_Coverage(t *testing.T) {
	os.Setenv("USE_AI_MOCK", "true")
	service := NewAIAnalysisService()

	data := ports.FinancialAnalysisData{
		UserID:        "test-user",
		TotalIncome:   5000000,
		TotalExpenses: 4000000,
		SavingsRate:   0.2,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1500000,
			"Transporte":   1200000,
			"Vivienda":     1300000,
		},
		IncomeStability: 0.8,
		FinancialScore:  750,
		Period:          "2024-01",
	}

	prompt := service.buildInsightsPrompt(data)

	assert.NotEmpty(t, prompt, "Prompt should not be empty")
	assert.Contains(t, prompt, "5000000", "Should contain income amount")
	assert.Contains(t, prompt, "4000000", "Should contain expenses amount")
	assert.Contains(t, prompt, "20.0", "Should contain savings rate percentage")
	assert.Contains(t, prompt, "750", "Should contain financial score")
	assert.Contains(t, prompt, "2024-01", "Should contain period")
	assert.Contains(t, prompt, "Alimentación", "Should contain top expense categories")
	assert.Contains(t, prompt, "JSON", "Should request JSON format")
}
