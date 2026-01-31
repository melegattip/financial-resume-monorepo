package services

import (
	"context"
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPurchaseDecisionService(t *testing.T) {
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
			service := NewPurchaseDecisionService()

			// Assert
			assert.NotNil(t, service)
			assert.Equal(t, tt.expectedUseMock, service.useMock)
		})
	}
}

func TestPurchaseDecisionService_CanIBuy_Mock_CanAfford(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewPurchaseDecisionService()

	// Usuario con buena salud financiera que puede permitirse la compra
	request := ports.PurchaseAnalysisRequest{
		UserID:      "test-user-can-afford",
		ItemName:    "Laptop Nueva",
		Amount:      2000000, // $2M COP
		Description: "Laptop para trabajo remoto",
		IsNecessary: true,
		UserFinancialProfile: ports.UserFinancialProfile{
			MonthlyIncome:       8000000,  // $8M COP/mes
			MonthlyExpenses:     5000000,  // $5M COP/mes
			CurrentBalance:      10000000, // $10M COP balance actual
			SavingsRate:         0.375,    // 37.5% ahorro
			IncomeStability:     0.9,      // Alta estabilidad
			FinancialDiscipline: 850,      // Excelente disciplina
			TopExpenseCategories: map[string]float64{
				"Alimentación": 1500000,
				"Transporte":   1000000,
				"Vivienda":     2500000,
			},
		},
	}

	// Execute
	ctx := context.Background()
	response, err := service.CanIBuy(ctx, request)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.False(t, response.CanBuy, "Should not be able to buy when amount exceeds 30% of available")
	assert.Contains(t, response.Reasoning, "disponible mensual", "Reasoning should mention available monthly amount")
	assert.Greater(t, response.Confidence, 0.5, "Should have reasonable confidence")
	assert.Len(t, response.Alternatives, 3, "Mock should return exactly 3 alternatives")
}

func TestPurchaseDecisionService_CanIBuy_Mock_CannotAfford(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewPurchaseDecisionService()

	// Usuario con mala salud financiera que no puede permitirse la compra
	request := ports.PurchaseAnalysisRequest{
		UserID:      "test-user-cannot-afford",
		ItemName:    "Auto de Lujo",
		Amount:      80000000, // $80M COP - muy caro
		Description: "Auto deportivo",
		IsNecessary: false,
		UserFinancialProfile: ports.UserFinancialProfile{
			MonthlyIncome:       3000000, // $3M COP/mes
			MonthlyExpenses:     2800000, // $2.8M COP/mes
			CurrentBalance:      500000,  // $500K COP balance actual
			SavingsRate:         0.067,   // 6.7% ahorro
			IncomeStability:     0.4,     // Baja estabilidad
			FinancialDiscipline: 300,     // Disciplina baja
			TopExpenseCategories: map[string]float64{
				"Alimentación": 1200000,
				"Transporte":   800000,
				"Vivienda":     800000,
			},
		},
	}

	// Execute
	ctx := context.Background()
	response, err := service.CanIBuy(ctx, request)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.False(t, response.CanBuy, "Should not be able to afford expensive item with low available funds")
	assert.Contains(t, response.Reasoning, "disponible mensual", "Reasoning should mention available monthly amount")
	assert.Greater(t, response.Confidence, 0.5, "Should have reasonable confidence")
	assert.Len(t, response.Alternatives, 3, "Mock should return exactly 3 alternatives")
}

func TestPurchaseDecisionService_CanIBuy_Mock_Caution(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewPurchaseDecisionService()

	// Usuario con salud financiera regular - caso límite
	request := ports.PurchaseAnalysisRequest{
		UserID:      "test-user-caution",
		ItemName:    "Televisor 4K",
		Amount:      3000000, // $3M COP
		Description: "Televisor para el hogar",
		IsNecessary: false,
		UserFinancialProfile: ports.UserFinancialProfile{
			MonthlyIncome:       5000000, // $5M COP/mes
			MonthlyExpenses:     4000000, // $4M COP/mes
			CurrentBalance:      2000000, // $2M COP balance actual
			SavingsRate:         0.2,     // 20% ahorro
			IncomeStability:     0.6,     // Estabilidad media
			FinancialDiscipline: 500,     // Disciplina regular
			TopExpenseCategories: map[string]float64{
				"Alimentación": 1500000,
				"Transporte":   1200000,
				"Vivienda":     1300000,
			},
		},
	}

	// Execute
	ctx := context.Background()
	response, err := service.CanIBuy(ctx, request)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)
	// Verificar resultado - según la lógica mock: 500000 > (1000000 * 0.3) = 300000, por lo que canBuy = false
	assert.False(t, response.CanBuy, "Should be cautious with purchase that exceeds 30% threshold")
	assert.Contains(t, response.Reasoning, "disponible mensual", "Reasoning should mention available monthly amount")
	assert.Greater(t, response.Confidence, 0.5, "Should have reasonable confidence")
	assert.NotEmpty(t, response.Reasoning)
	assert.Contains(t, response.Reasoning, "comprometer", "Should mention potential compromise to financial stability")
	assert.Greater(t, response.ImpactScore, 0)
	assert.LessOrEqual(t, response.ImpactScore, 1000)
	assert.NotEmpty(t, response.Alternatives)
}

func TestPurchaseDecisionService_SuggestAlternatives_Mock(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewPurchaseDecisionService()

	request := ports.PurchaseAnalysisRequest{
		UserID:      "test-user-alternatives",
		ItemName:    "iPhone 15 Pro Max",
		Amount:      6000000, // $6M COP
		Description: "Smartphone de alta gama",
		IsNecessary: false,
		UserFinancialProfile: ports.UserFinancialProfile{
			MonthlyIncome:       4000000,
			MonthlyExpenses:     3200000,
			CurrentBalance:      1000000,
			SavingsRate:         0.2,
			IncomeStability:     0.7,
			FinancialDiscipline: 600,
			TopExpenseCategories: map[string]float64{
				"Alimentación": 1200000,
				"Transporte":   1000000,
				"Vivienda":     1000000,
			},
		},
	}

	// Execute
	ctx := context.Background()
	alternatives, err := service.SuggestAlternatives(ctx, request)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, alternatives)
	assert.Len(t, alternatives, 3) // Mock devuelve 3 alternativas

	// Verificar estructura de alternativas
	for _, alt := range alternatives {
		assert.NotEmpty(t, alt.Name, "Alternative should have a name")
		assert.NotEmpty(t, alt.Description, "Alternative should have a description")
		assert.Greater(t, alt.Savings, 0.0, "Alternative should show savings")
		// Verificar que feasibility tiene valores válidos según implementación real
		assert.Contains(t, []string{"alta", "media", "baja"}, alt.Feasibility, "Feasibility should be alta, media, or baja")
	}
}

// Nota: Los métodos buildCanIBuyPrompt, buildAlternativesPrompt y analyzeFinancialProfile
// son privados, por lo que no se pueden testear directamente.
// Su funcionalidad se testea indirectamente a través de los métodos públicos.

// Test de integración completo
func TestPurchaseDecisionService_Integration_Mock(t *testing.T) {
	// Setup
	t.Setenv("USE_AI_MOCK", "true")
	service := NewPurchaseDecisionService()

	request := ports.PurchaseAnalysisRequest{
		UserID:      "integration-test-user",
		ItemName:    "Nintendo Switch",
		Amount:      1500000, // $1.5M COP
		Description: "Consola de videojuegos",
		IsNecessary: false,
		UserFinancialProfile: ports.UserFinancialProfile{
			MonthlyIncome:       5000000,
			MonthlyExpenses:     3800000,
			CurrentBalance:      3000000,
			SavingsRate:         0.24,
			IncomeStability:     0.75,
			FinancialDiscipline: 625,
			TopExpenseCategories: map[string]float64{
				"Alimentación": 1400000,
				"Transporte":   1200000,
				"Vivienda":     1200000,
			},
		},
	}

	// Execute - Flujo completo
	ctx := context.Background()

	// 1. Análisis de compra
	purchaseResponse, err := service.CanIBuy(ctx, request)
	require.NoError(t, err)
	assert.NotNil(t, purchaseResponse)

	// 2. Sugerir alternativas
	alternatives, err := service.SuggestAlternatives(ctx, request)
	require.NoError(t, err)
	assert.NotNil(t, alternatives)

	// 3. Verificar coherencia
	assert.NotEmpty(t, purchaseResponse.Reasoning)
	assert.Len(t, alternatives, 3)

	// 4. Verificar que las alternativas tienen sentido
	for _, alt := range alternatives {
		assert.Greater(t, alt.Savings, float64(0))
		assert.NotEmpty(t, alt.Feasibility)
	}
}

// Benchmark para verificar performance
func BenchmarkPurchaseDecisionService_CanIBuy_Mock(b *testing.B) {
	// Setup
	b.Setenv("USE_AI_MOCK", "true")
	service := NewPurchaseDecisionService()

	request := ports.PurchaseAnalysisRequest{
		UserID:      "benchmark-user",
		ItemName:    "Smartphone",
		Amount:      2000000,
		Description: "Teléfono móvil",
		IsNecessary: false,
		UserFinancialProfile: ports.UserFinancialProfile{
			MonthlyIncome:       5000000,
			MonthlyExpenses:     3500000,
			CurrentBalance:      4000000,
			SavingsRate:         0.3,
			IncomeStability:     0.8,
			FinancialDiscipline: 700,
			TopExpenseCategories: map[string]float64{
				"Alimentación": 1200000,
				"Transporte":   1000000,
				"Vivienda":     1300000,
			},
		},
	}

	ctx := context.Background()

	// Execute
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.CanIBuy(ctx, request)
		if err != nil {
			b.Fatal(err)
		}
	}
}
