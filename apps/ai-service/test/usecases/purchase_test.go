package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/financial-ai-service/internal/core/ports"
	"github.com/financial-ai-service/internal/core/usecases"
	"github.com/financial-ai-service/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// PurchaseUseCaseTestSuite suite de tests para el caso de uso de compra
type PurchaseUseCaseTestSuite struct {
	suite.Suite
	mockOpenAI *mocks.MockOpenAIClient
	mockCache  *mocks.MockCacheClient
	useCase    ports.PurchaseDecisionPort
	ctx        context.Context
}

// SetupTest configura cada test individual
func (suite *PurchaseUseCaseTestSuite) SetupTest() {
	suite.mockOpenAI = new(mocks.MockOpenAIClient)
	suite.mockCache = new(mocks.MockCacheClient)
	suite.useCase = usecases.NewPurchaseUseCase(suite.mockOpenAI, suite.mockCache)
	suite.ctx = context.Background()
}

// TearDownTest verifica que todas las expectativas se cumplieron
func (suite *PurchaseUseCaseTestSuite) TearDownTest() {
	suite.mockOpenAI.AssertExpectations(suite.T())
	suite.mockCache.AssertExpectations(suite.T())
}

// TestCanIBuy_Success testa el análisis exitoso de decisión de compra
func (suite *PurchaseUseCaseTestSuite) TestCanIBuy_Success() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest
	expectedDecision := mocks.TestPurchaseDecision
	decisionJSON, _ := json.Marshal(expectedDecision)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(decisionJSON), nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.CanIBuy(suite.ctx, request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), true, result.CanBuy)
	assert.Equal(suite.T(), 0.85, result.Confidence)
	assert.Equal(suite.T(), 25, result.ImpactScore)
	assert.NotEmpty(suite.T(), result.Reasoning)
	assert.NotZero(suite.T(), result.GeneratedAt)
}

// TestCanIBuy_CacheHit testa el análisis con cache hit
func (suite *PurchaseUseCaseTestSuite) TestCanIBuy_CacheHit() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest
	cachedDecision := mocks.TestPurchaseDecision
	cachedData, _ := json.Marshal(cachedDecision)

	// Mock cache hit
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return(cachedData, nil)

	// Act
	result, err := suite.useCase.CanIBuy(suite.ctx, request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), cachedDecision.CanBuy, result.CanBuy)
	assert.Equal(suite.T(), cachedDecision.Confidence, result.Confidence)

	// Verificar que no se llamó a OpenAI
	suite.mockOpenAI.AssertNotCalled(suite.T(), "GenerateAnalysis")
}

// TestCanIBuy_OpenAIError testa el manejo de errores de OpenAI
func (suite *PurchaseUseCaseTestSuite) TestCanIBuy_OpenAIError() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI error
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", errors.New("OpenAI API error"))

	// Act
	result, err := suite.useCase.CanIBuy(suite.ctx, request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error analyzing purchase decision")
}

// TestCanIBuy_InvalidJSON testa el manejo de JSON inválido
func (suite *PurchaseUseCaseTestSuite) TestCanIBuy_InvalidJSON() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest
	invalidJSON := "invalid json response"

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response with invalid JSON
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(invalidJSON, nil)

	// Act
	result, err := suite.useCase.CanIBuy(suite.ctx, request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error parsing purchase decision response")
}

// TestSuggestAlternatives_Success testa la sugerencia exitosa de alternativas
func (suite *PurchaseUseCaseTestSuite) TestSuggestAlternatives_Success() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest
	expectedAlternatives := mocks.TestAlternatives
	alternativesJSON, _ := json.Marshal(expectedAlternatives)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(alternativesJSON), nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.SuggestAlternatives(suite.ctx, request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "MacBook Air", result[0].Name)
	assert.Equal(suite.T(), "Modelo anterior", result[1].Name)
	assert.Equal(suite.T(), float64(3000000), result[0].Savings)
	assert.Equal(suite.T(), float64(2000000), result[1].Savings)
}

// TestSuggestAlternatives_CacheHit testa las sugerencias con cache hit
func (suite *PurchaseUseCaseTestSuite) TestSuggestAlternatives_CacheHit() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest
	cachedAlternatives := mocks.TestAlternatives
	cachedData, _ := json.Marshal(cachedAlternatives)

	// Mock cache hit
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return(cachedData, nil)

	// Act
	result, err := suite.useCase.SuggestAlternatives(suite.ctx, request)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)

	// Verificar que no se llamó a OpenAI
	suite.mockOpenAI.AssertNotCalled(suite.T(), "GenerateAnalysis")
}

// TestSuggestAlternatives_OpenAIError testa el manejo de errores de OpenAI
func (suite *PurchaseUseCaseTestSuite) TestSuggestAlternatives_OpenAIError() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI error
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", errors.New("OpenAI API error"))

	// Act
	result, err := suite.useCase.SuggestAlternatives(suite.ctx, request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error generating alternatives")
}

// TestSuggestAlternatives_InvalidJSON testa el manejo de JSON inválido
func (suite *PurchaseUseCaseTestSuite) TestSuggestAlternatives_InvalidJSON() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest
	invalidJSON := "invalid json response"

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response with invalid JSON
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(invalidJSON, nil)

	// Act
	result, err := suite.useCase.SuggestAlternatives(suite.ctx, request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error parsing alternatives response")
}

// TestPurchaseWithDifferentAmounts testa análisis con diferentes montos
func (suite *PurchaseUseCaseTestSuite) TestPurchaseWithDifferentAmounts() {
	amounts := []float64{1000000, 5000000, 10000000, 50000000}

	for _, amount := range amounts {
		suite.Run(fmt.Sprintf("Amount_%.0f", amount), func() {
			// Arrange
			request := mocks.TestPurchaseAnalysisRequest
			request.Amount = amount

			var expectedDecision ports.PurchaseDecision
			// Ajustar la decisión basada en el monto
			if amount <= 10000000 {
				expectedDecision = mocks.TestPurchaseDecision
			} else {
				expectedDecision = mocks.TestPurchaseDecisionCannotBuy
			}
			decisionJSON, _ := json.Marshal(expectedDecision)

			// Mock cache miss
			suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
				Return([]byte(nil), errors.New("cache miss")).Once()

			// Mock OpenAI response
			suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
				Return(string(decisionJSON), nil).Once()

			// Mock cache set
			suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
				Return(nil).Once()

			// Act
			result, err := suite.useCase.CanIBuy(suite.ctx, request)

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), result)
			assert.Equal(suite.T(), expectedDecision.CanBuy, result.CanBuy)
		})
	}
}

// TestPurchaseWithDifferentPaymentTypes testa análisis con diferentes tipos de pago
func (suite *PurchaseUseCaseTestSuite) TestPurchaseWithDifferentPaymentTypes() {
	paymentTypes := [][]string{
		{"contado"},
		{"credito"},
		{"contado", "credito"},
		{"debito"},
	}

	for _, payments := range paymentTypes {
		suite.Run(fmt.Sprintf("Payment_%v", payments), func() {
			// Arrange
			request := mocks.TestPurchaseAnalysisRequest
			request.PaymentTypes = payments

			expectedDecision := mocks.TestPurchaseDecision
			decisionJSON, _ := json.Marshal(expectedDecision)

			// Mock cache miss
			suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
				Return([]byte(nil), errors.New("cache miss"))

			// Mock OpenAI response
			suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
				Return(string(decisionJSON), nil)

			// Mock cache set
			suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
				Return(nil)

			// Act
			result, err := suite.useCase.CanIBuy(suite.ctx, request)

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), result)
		})
	}
}

// TestPurchaseNecessaryVsNonNecessary testa análisis entre compras necesarias y no necesarias
func (suite *PurchaseUseCaseTestSuite) TestPurchaseNecessaryVsNonNecessary() {
	testCases := []struct {
		name        string
		isNecessary bool
		expectedBuy bool
	}{
		{
			name:        "Compra necesaria",
			isNecessary: true,
			expectedBuy: true,
		},
		{
			name:        "Compra no necesaria",
			isNecessary: false,
			expectedBuy: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			request := mocks.TestPurchaseAnalysisRequest
			request.IsNecessary = tc.isNecessary

			var expectedDecision ports.PurchaseDecision
			if tc.expectedBuy {
				expectedDecision = mocks.TestPurchaseDecision
			} else {
				expectedDecision = mocks.TestPurchaseDecisionCannotBuy
			}
			decisionJSON, _ := json.Marshal(expectedDecision)

			// Mock cache miss
			suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
				Return([]byte(nil), errors.New("cache miss")).Once()

			// Mock OpenAI response
			suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
				Return(string(decisionJSON), nil).Once()

			// Mock cache set
			suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
				Return(nil).Once()

			// Act
			result, err := suite.useCase.CanIBuy(suite.ctx, request)

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), result)
			assert.Equal(suite.T(), tc.expectedBuy, result.CanBuy)
		})
	}
}

// TestPurchaseWithDifferentFinancialProfiles testa análisis con diferentes perfiles financieros
func (suite *PurchaseUseCaseTestSuite) TestPurchaseWithDifferentFinancialProfiles() {
	testCases := []struct {
		name               string
		currentBalance     float64
		monthlyIncome      float64
		expectedCanBuy     bool
		expectedConfidence float64
	}{
		{
			name:               "Perfil excelente",
			currentBalance:     20000000,
			monthlyIncome:      10000000,
			expectedCanBuy:     true,
			expectedConfidence: 0.95,
		},
		{
			name:               "Perfil pobre",
			currentBalance:     1000000,
			monthlyIncome:      2000000,
			expectedCanBuy:     false,
			expectedConfidence: 0.2,
		},
		{
			name:               "Perfil promedio",
			currentBalance:     5000000,
			monthlyIncome:      4000000,
			expectedCanBuy:     true,
			expectedConfidence: 0.6,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Arrange
			request := mocks.TestPurchaseAnalysisRequest
			request.UserFinancialProfile.CurrentBalance = tc.currentBalance
			request.UserFinancialProfile.MonthlyIncome = tc.monthlyIncome

			// Crear decisión específica para cada perfil
			expectedDecision := ports.PurchaseDecision{
				CanBuy:       tc.expectedCanBuy,
				Confidence:   tc.expectedConfidence,
				Reasoning:    "Análisis basado en perfil financiero",
				Alternatives: []string{"Buscar ofertas", "Considerar modelo anterior"},
				ImpactScore:  25,
				GeneratedAt:  time.Now(),
			}
			decisionJSON, _ := json.Marshal(expectedDecision)

			// Mock cache miss
			suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
				Return([]byte(nil), errors.New("cache miss")).Once()

			// Mock OpenAI response
			suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
				Return(string(decisionJSON), nil).Once()

			// Mock cache set
			suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
				Return(nil).Once()

			// Act
			result, err := suite.useCase.CanIBuy(suite.ctx, request)

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), result)
			assert.Equal(suite.T(), tc.expectedCanBuy, result.CanBuy)
			assert.Equal(suite.T(), tc.expectedConfidence, result.Confidence)
		})
	}
}

// TestContextCancellation testa el manejo de cancelación de contexto
func (suite *PurchaseUseCaseTestSuite) TestContextCancellation() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI que debería recibir el contexto cancelado
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", context.Canceled)

	// Act
	result, err := suite.useCase.CanIBuy(cancelledCtx, request)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error analyzing purchase decision")
}

// TestCacheFailure_DoesNotAffectOperation testa que los errores de cache no afectan la operación
func (suite *PurchaseUseCaseTestSuite) TestCacheFailure_DoesNotAffectOperation() {
	// Arrange
	request := mocks.TestPurchaseAnalysisRequest
	expectedDecision := mocks.TestPurchaseDecision
	decisionJSON, _ := json.Marshal(expectedDecision)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(decisionJSON), nil)

	// Mock cache set failure (no debería afectar la operación)
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(errors.New("cache set failed"))

	// Act
	result, err := suite.useCase.CanIBuy(suite.ctx, request)

	// Assert
	assert.NoError(suite.T(), err) // El error de cache no debería afectar el resultado
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedDecision.CanBuy, result.CanBuy)
}

// TestRunSuite ejecuta todos los tests del suite
func TestPurchaseUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(PurchaseUseCaseTestSuite))
}

// Tests adicionales para casos específicos

// TestPurchaseWithEmptyDescription testa análisis con descripción vacía
func TestPurchaseWithEmptyDescription(t *testing.T) {
	// Arrange
	mockOpenAI := new(mocks.MockOpenAIClient)
	mockCache := new(mocks.MockCacheClient)
	useCase := usecases.NewPurchaseUseCase(mockOpenAI, mockCache)

	request := mocks.TestPurchaseAnalysisRequest
	request.Description = "" // Descripción vacía

	expectedDecision := mocks.TestPurchaseDecision
	decisionJSON, _ := json.Marshal(expectedDecision)

	// Mock cache miss
	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(decisionJSON), nil)

	// Mock cache set
	mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := useCase.CanIBuy(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Debería funcionar incluso sin descripción

	// Verificar expectativas
	mockOpenAI.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestPurchaseWithZeroAmount testa análisis con monto cero
func TestPurchaseWithZeroAmount(t *testing.T) {
	// Arrange
	mockOpenAI := new(mocks.MockOpenAIClient)
	mockCache := new(mocks.MockCacheClient)
	useCase := usecases.NewPurchaseUseCase(mockOpenAI, mockCache)

	request := mocks.TestPurchaseAnalysisRequest
	request.Amount = 0 // Monto cero

	expectedDecision := mocks.TestPurchaseDecision
	expectedDecision.CanBuy = true // Monto cero siempre debería ser "comprável"
	decisionJSON, _ := json.Marshal(expectedDecision)

	// Mock cache miss
	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(decisionJSON), nil)

	// Mock cache set
	mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := useCase.CanIBuy(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, true, result.CanBuy)

	// Verificar expectativas
	mockOpenAI.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
