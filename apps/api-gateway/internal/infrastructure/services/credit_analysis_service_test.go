package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCreditAnalysisService(t *testing.T) {
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
			expectedUseMock: false,
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
			service := NewCreditAnalysisService()

			// Assert
			assert.NotNil(t, service)
			assert.Equal(t, tt.expectedUseMock, service.useMock)
		})
	}
}

func TestCreditAnalysisService_CalculateCreditScore_Mock(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewCreditAnalysisService()

	tests := []struct {
		name        string
		data        ports.FinancialAnalysisData
		expectedMin int
		expectedMax int
		description string
	}{
		{
			name: "Excellent financial profile should yield high credit score",
			data: ports.FinancialAnalysisData{
				UserID:        "test-user-excellent",
				TotalIncome:   10000000,
				TotalExpenses: 6000000,
				SavingsRate:   0.4,
				ExpensesByCategory: map[string]float64{
					"Alimentación": 2000000,
					"Transporte":   1500000,
					"Vivienda":     2500000,
				},
				IncomeStability: 0.9,
				FinancialScore:  850,
				Period:          "2024-01",
			},
			expectedMin: 700,
			expectedMax: 1000,
			description: "High income and savings should yield excellent credit score",
		},
		{
			name: "Poor financial profile should yield low credit score",
			data: ports.FinancialAnalysisData{
				TotalIncome:     2000000,
				TotalExpenses:   2100000, // Gastos > ingresos
				SavingsRate:     -0.05,   // Tasa de ahorro negativa
				IncomeStability: 0.3,     // Baja estabilidad
				FinancialScore:  400,
			},
			expectedMin: 300,
			expectedMax: 600, // Ajustado: baseScore=600 - 50 (savings<0) + 30 (stability) - 25 (expenses>income) - 10 (low score) ≈ 545
			description: "Negative savings and low stability should yield poor credit score",
		},
		{
			name: "Average financial profile should yield moderate credit score",
			data: ports.FinancialAnalysisData{
				TotalIncome:     3000000,
				TotalExpenses:   2400000,
				SavingsRate:     0.15, // 15% ahorro
				IncomeStability: 0.8,  // Buena estabilidad
				FinancialScore:  700,
			},
			expectedMin: 700,
			expectedMax: 1000,
			description: "Moderate profile should yield average credit score",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			ctx := context.Background()
			score, err := service.CalculateCreditScore(ctx, tt.data)

			// Assert
			require.NoError(t, err)
			assert.GreaterOrEqual(t, score, tt.expectedMin, tt.description)
			assert.LessOrEqual(t, score, tt.expectedMax, tt.description)
			assert.GreaterOrEqual(t, score, 0, "Credit score should never be negative")
			assert.LessOrEqual(t, score, 1000, "Credit score should not exceed 1000")
		})
	}
}

func TestCreditAnalysisService_GenerateImprovementPlan_Mock(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewCreditAnalysisService()

	testData := ports.FinancialAnalysisData{
		UserID:        "test-user-plan",
		TotalIncome:   4000000,
		TotalExpenses: 3200000,
		SavingsRate:   0.2,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1200000,
			"Transporte":   1000000,
			"Vivienda":     1000000,
		},
		IncomeStability: 0.7,
		FinancialScore:  550,
		Period:          "2024-01",
	}

	// Execute
	ctx := context.Background()
	plan, err := service.GenerateImprovementPlan(ctx, testData)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, plan)

	// Verificar estructura del plan
	assert.Greater(t, plan.CurrentScore, 0)
	assert.LessOrEqual(t, plan.CurrentScore, 1000)
	assert.Greater(t, plan.TargetScore, plan.CurrentScore)
	assert.LessOrEqual(t, plan.TargetScore, 1000)
	assert.Greater(t, plan.TimelineMonths, 0)
	assert.LessOrEqual(t, plan.TimelineMonths, 36) // Máximo 3 años
	assert.NotEmpty(t, plan.Actions)
	assert.Len(t, plan.Actions, 3, "Mock should return exactly 3 actions")
	assert.NotNil(t, plan.KeyMetrics)
	assert.True(t, plan.GeneratedAt.After(time.Now().Add(-1*time.Minute)))

	// Verificar estructura de acciones
	for _, action := range plan.Actions {
		assert.NotEmpty(t, action.Title)
		assert.NotEmpty(t, action.Description)
		assert.Contains(t, []string{"high", "medium", "low"}, action.Priority)
		assert.NotEmpty(t, action.Timeline)
		assert.Greater(t, action.Impact, 0)
		assert.LessOrEqual(t, action.Impact, 100)
		assert.Contains(t, []string{"easy", "medium", "hard"}, action.Difficulty)
	}

	// Verificar métricas clave
	expectedMetrics := []string{"target_savings_rate", "debt_to_income_target", "emergency_fund_months", "diversification_score"}
	for _, metric := range expectedMetrics {
		assert.Contains(t, plan.KeyMetrics, metric, fmt.Sprintf("Key metrics should contain %s", metric))
	}
}

func TestCreditAnalysisService_GenerateImprovementPlan_HighScore(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewCreditAnalysisService()

	// Usuario con score ya alto
	testData := ports.FinancialAnalysisData{
		UserID:        "test-user-high-score",
		TotalIncome:   8000000,
		TotalExpenses: 4800000,
		SavingsRate:   0.4,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1600000,
			"Transporte":   1200000,
			"Vivienda":     2000000,
		},
		IncomeStability: 0.9,
		FinancialScore:  800,
		Period:          "2024-01",
	}

	// Execute
	ctx := context.Background()
	plan, err := service.GenerateImprovementPlan(ctx, testData)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Greater(t, plan.CurrentScore, 700)      // Score alto
	assert.LessOrEqual(t, plan.TimelineMonths, 12) // Timeline más corto para mejoras menores
	assert.NotEmpty(t, plan.Actions)

	// Para scores altos, las acciones deberían ser de mantenimiento
	hasMaintenanceAction := false
	for _, action := range plan.Actions {
		if action.Priority == "low" || action.Difficulty == "easy" {
			hasMaintenanceAction = true
			break
		}
	}
	assert.True(t, hasMaintenanceAction, "High score plans should include maintenance actions")
}

func TestCreditAnalysisService_GenerateImprovementPlan_LowScore(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewCreditAnalysisService()

	// Usuario con score bajo que necesita mejoras urgentes
	testData := ports.FinancialAnalysisData{
		UserID:        "test-user-low-score",
		TotalIncome:   3000000,
		TotalExpenses: 3200000, // Gastos > ingresos
		SavingsRate:   -0.067,  // Ahorro negativo
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1200000,
			"Transporte":   1000000,
			"Vivienda":     1000000,
		},
		IncomeStability: 0.4,
		FinancialScore:  300,
		Period:          "2024-01",
	}

	// Execute
	ctx := context.Background()
	plan, err := service.GenerateImprovementPlan(ctx, testData)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Less(t, plan.CurrentScore, 600, "Low profile should yield low credit score")
	assert.Greater(t, plan.TimelineMonths, 12) // Timeline más largo para mejoras significativas
	assert.NotEmpty(t, plan.Actions)

	// Para scores bajos, debería haber acciones de alta prioridad
	hasHighPriorityAction := false
	for _, action := range plan.Actions {
		if action.Priority == "high" {
			hasHighPriorityAction = true
			break
		}
	}
	assert.True(t, hasHighPriorityAction, "Low score plans should include high priority actions")
}

// Test de integración básico
func TestCreditAnalysisService_Integration_Mock(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewCreditAnalysisService()

	testData := ports.FinancialAnalysisData{
		UserID:        "integration-test-user",
		TotalIncome:   6000000,
		TotalExpenses: 4500000,
		SavingsRate:   0.25,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1500000,
			"Transporte":   1200000,
			"Vivienda":     1800000,
		},
		IncomeStability: 0.8,
		FinancialScore:  650,
		Period:          "2024-01",
	}

	// Execute - Flujo completo
	ctx := context.Background()

	// 1. Calcular score crediticio
	creditScore, err := service.CalculateCreditScore(ctx, testData)
	require.NoError(t, err)
	assert.Greater(t, creditScore, 0)
	assert.LessOrEqual(t, creditScore, 1000)

	// 2. Generar plan de mejora
	improvementPlan, err := service.GenerateImprovementPlan(ctx, testData)
	require.NoError(t, err)
	assert.NotNil(t, improvementPlan)

	// 3. Verificar coherencia entre score calculado y plan
	assert.Equal(t, creditScore, improvementPlan.CurrentScore)
	assert.Greater(t, improvementPlan.TargetScore, improvementPlan.CurrentScore)

	// 4. Verificar que el plan tiene sentido
	assert.NotEmpty(t, improvementPlan.Actions)
	assert.Greater(t, improvementPlan.TimelineMonths, 0)
	assert.True(t, improvementPlan.GeneratedAt.After(time.Now().Add(-5*time.Second)))
}

// Test de casos edge
func TestCreditAnalysisService_EdgeCases_Mock(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewCreditAnalysisService()

	tests := []struct {
		name string
		data ports.FinancialAnalysisData
	}{
		{
			name: "Zero income case",
			data: ports.FinancialAnalysisData{
				UserID:        "test-zero-income",
				TotalIncome:   0,
				TotalExpenses: 1000000,
				SavingsRate:   -1.0,
				ExpensesByCategory: map[string]float64{
					"Alimentación": 1000000,
				},
				IncomeStability: 0.0,
				FinancialScore:  0,
				Period:          "2024-01",
			},
		},
		{
			name: "Very high income case",
			data: ports.FinancialAnalysisData{
				UserID:        "test-high-income",
				TotalIncome:   100000000, // $100M COP
				TotalExpenses: 20000000,  // $20M COP
				SavingsRate:   0.8,       // 80% ahorro
				ExpensesByCategory: map[string]float64{
					"Alimentación": 5000000,
					"Transporte":   5000000,
					"Vivienda":     10000000,
				},
				IncomeStability: 1.0,
				FinancialScore:  850,
				Period:          "2024-01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Execute - No debería fallar incluso con casos extremos
			score, err := service.CalculateCreditScore(ctx, tt.data)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, score, 0)
			assert.LessOrEqual(t, score, 1000)

			plan, err := service.GenerateImprovementPlan(ctx, tt.data)
			require.NoError(t, err)
			assert.NotNil(t, plan)
			assert.NotEmpty(t, plan.Actions)
		})
	}
}

// Benchmark para verificar performance
func BenchmarkCreditAnalysisService_CalculateCreditScore_Mock(b *testing.B) {
	// Setup
	b.Setenv("USE_AI_MOCK", "true")
	service := NewCreditAnalysisService()

	testData := ports.FinancialAnalysisData{
		UserID:        "benchmark-user",
		TotalIncome:   5000000,
		TotalExpenses: 3500000,
		SavingsRate:   0.3,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1200000,
			"Transporte":   1000000,
			"Vivienda":     1300000,
		},
		IncomeStability: 0.8,
		FinancialScore:  700,
		Period:          "2024-01",
	}

	ctx := context.Background()

	// Execute
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.CalculateCreditScore(ctx, testData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCreditAnalysisService_GenerateImprovementPlan_Mock(b *testing.B) {
	// Setup
	b.Setenv("USE_AI_MOCK", "true")
	service := NewCreditAnalysisService()

	testData := ports.FinancialAnalysisData{
		UserID:        "benchmark-user",
		TotalIncome:   5000000,
		TotalExpenses: 3500000,
		SavingsRate:   0.3,
		ExpensesByCategory: map[string]float64{
			"Alimentación": 1200000,
			"Transporte":   1000000,
			"Vivienda":     1300000,
		},
		IncomeStability: 0.8,
		FinancialScore:  700,
		Period:          "2024-01",
	}

	ctx := context.Background()

	// Execute
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateImprovementPlan(ctx, testData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
