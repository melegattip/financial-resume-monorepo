package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/financial-ai-service/internal/core/ports"
	"github.com/financial-ai-service/internal/core/usecases"
	"github.com/financial-ai-service/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// CreditUseCaseTestSuite suite de tests para el caso de uso de crédito
type CreditUseCaseTestSuite struct {
	suite.Suite
	mockOpenAI *mocks.MockOpenAIClient
	mockCache  *mocks.MockCacheClient
	useCase    ports.CreditAnalysisPort
	ctx        context.Context
}

// SetupTest configura cada test individual
func (suite *CreditUseCaseTestSuite) SetupTest() {
	suite.mockOpenAI = new(mocks.MockOpenAIClient)
	suite.mockCache = new(mocks.MockCacheClient)
	suite.useCase = usecases.NewCreditUseCase(suite.mockOpenAI, suite.mockCache)
	suite.ctx = context.Background()
}

// TearDownTest verifica que todas las expectativas se cumplieron
func (suite *CreditUseCaseTestSuite) TearDownTest() {
	suite.mockOpenAI.AssertExpectations(suite.T())
	suite.mockCache.AssertExpectations(suite.T())
}

// TestGenerateImprovementPlan_Success testa la generación exitosa del plan de mejora
func (suite *CreditUseCaseTestSuite) TestGenerateImprovementPlan_Success() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	expectedPlan := mocks.TestCreditPlan
	planJSON, _ := json.Marshal(expectedPlan)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(planJSON), nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.GenerateImprovementPlan(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 750, result.CurrentScore)
	assert.Equal(suite.T(), 800, result.TargetScore)
	assert.Equal(suite.T(), 12, result.TimelineMonths)
	assert.Len(suite.T(), result.Actions, 2)
	assert.NotZero(suite.T(), result.GeneratedAt)
}

// TestGenerateImprovementPlan_CacheHit testa el plan con cache hit
func (suite *CreditUseCaseTestSuite) TestGenerateImprovementPlan_CacheHit() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	cachedPlan := mocks.TestCreditPlan
	cachedData, _ := json.Marshal(cachedPlan)

	// Mock cache hit
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return(cachedData, nil)

	// Act
	result, err := suite.useCase.GenerateImprovementPlan(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), cachedPlan.CurrentScore, result.CurrentScore)
	assert.Equal(suite.T(), cachedPlan.TargetScore, result.TargetScore)

	// Verificar que no se llamó a OpenAI
	suite.mockOpenAI.AssertNotCalled(suite.T(), "GenerateAnalysis")
}

// TestGenerateImprovementPlan_OpenAIError testa el manejo de errores de OpenAI
func (suite *CreditUseCaseTestSuite) TestGenerateImprovementPlan_OpenAIError() {
	// Arrange
	data := mocks.TestFinancialAnalysisData

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI error
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", errors.New("OpenAI API error"))

	// Act
	result, err := suite.useCase.GenerateImprovementPlan(suite.ctx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error generating credit improvement plan")
}

// TestGenerateImprovementPlan_InvalidJSON testa el manejo de JSON inválido
func (suite *CreditUseCaseTestSuite) TestGenerateImprovementPlan_InvalidJSON() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	invalidJSON := "invalid json response"

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response with invalid JSON
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(invalidJSON, nil)

	// Act
	result, err := suite.useCase.GenerateImprovementPlan(suite.ctx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error parsing credit plan response")
}

// TestCalculateCreditScore_Success testa el cálculo exitoso del score
func (suite *CreditUseCaseTestSuite) TestCalculateCreditScore_Success() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	expectedScore := 750
	scoreResponse := struct {
		Score int `json:"score"`
	}{Score: expectedScore}
	scoreJSON, _ := json.Marshal(scoreResponse)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(scoreJSON), nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.CalculateCreditScore(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedScore, result)
}

// TestCalculateCreditScore_CacheHit testa el cálculo con cache hit
func (suite *CreditUseCaseTestSuite) TestCalculateCreditScore_CacheHit() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	expectedScore := 750
	scoreResponse := struct {
		Score int `json:"score"`
	}{Score: expectedScore}
	cachedData, _ := json.Marshal(scoreResponse)

	// Mock cache hit
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return(cachedData, nil)

	// Act
	result, err := suite.useCase.CalculateCreditScore(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedScore, result)

	// Verificar que no se llamó a OpenAI
	suite.mockOpenAI.AssertNotCalled(suite.T(), "GenerateAnalysis")
}

// TestCalculateCreditScore_OpenAIError testa el manejo de errores de OpenAI
func (suite *CreditUseCaseTestSuite) TestCalculateCreditScore_OpenAIError() {
	// Arrange
	data := mocks.TestFinancialAnalysisData

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI error
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", errors.New("OpenAI API error"))

	// Act
	result, err := suite.useCase.CalculateCreditScore(suite.ctx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, result)
	assert.Contains(suite.T(), err.Error(), "error calculating credit score")
}

// TestCalculateCreditScore_InvalidJSON testa el manejo de JSON inválido
func (suite *CreditUseCaseTestSuite) TestCalculateCreditScore_InvalidJSON() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	invalidJSON := "invalid json response"

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response with invalid JSON
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(invalidJSON, nil)

	// Act
	result, err := suite.useCase.CalculateCreditScore(suite.ctx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, result)
	assert.Contains(suite.T(), err.Error(), "error parsing credit score response")
}

// TestCalculateCreditScore_InvalidScore testa el manejo de score inválido
func (suite *CreditUseCaseTestSuite) TestCalculateCreditScore_InvalidScore() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	invalidScore := 0 // Score fuera del rango válido (1-1000)
	scoreResponse := struct {
		Score int `json:"score"`
	}{Score: invalidScore}
	scoreJSON, _ := json.Marshal(scoreResponse)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response with invalid score
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(scoreJSON), nil)

	// Mock cache set (para el score calculado por defecto)
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.CalculateCreditScore(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), result, 300)
	assert.LessOrEqual(suite.T(), result, 850)
	// Debería usar el algoritmo de fallback
}

// TestDefaultScoreCalculation testa el algoritmo de cálculo por defecto
func (suite *CreditUseCaseTestSuite) TestDefaultScoreCalculation() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	extremeScore := 1000 // Score extremadamente alto para forzar el fallback
	scoreResponse := struct {
		Score int `json:"score"`
	}{Score: extremeScore}
	scoreJSON, _ := json.Marshal(scoreResponse)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response with extreme score
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(scoreJSON), nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.CalculateCreditScore(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), result, 1)
	assert.LessOrEqual(suite.T(), result, 1000)
	// Debería usar el algoritmo de fallback y estar en el rango válido
}

// TestScoreCalculationWithPoorFinances testa el cálculo con finanzas pobres
func (suite *CreditUseCaseTestSuite) TestScoreCalculationWithPoorFinances() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	data.SavingsRate = 0.0                      // Sin ahorros
	data.IncomeStability = 0.2                  // Ingresos muy inestables
	data.TotalExpenses = data.TotalIncome * 1.2 // Gastos mayores que ingresos

	invalidScore := 100 // Score inválido para forzar el fallback
	scoreResponse := struct {
		Score int `json:"score"`
	}{Score: invalidScore}
	scoreJSON, _ := json.Marshal(scoreResponse)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response with invalid score
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(scoreJSON), nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.CalculateCreditScore(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), result, 1)
	assert.LessOrEqual(suite.T(), result, 1000)
	// Con finanzas pobres, el score debería ser relativamente bajo
	assert.Less(suite.T(), result, 600) // Score debería ser bajo
}

// TestScoreCalculationWithExcellentFinances testa el cálculo con finanzas excelentes
func (suite *CreditUseCaseTestSuite) TestScoreCalculationWithExcellentFinances() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	data.SavingsRate = 0.5                      // Excelente tasa de ahorro
	data.IncomeStability = 1.0                  // Ingresos completamente estables
	data.TotalExpenses = data.TotalIncome * 0.5 // Gastos bien controlados

	invalidScore := 1100 // Score inválido para forzar el fallback
	scoreResponse := struct {
		Score int `json:"score"`
	}{Score: invalidScore}
	scoreJSON, _ := json.Marshal(scoreResponse)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response with invalid score
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(scoreJSON), nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.CalculateCreditScore(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), result, 1)
	assert.LessOrEqual(suite.T(), result, 1000)
	// Con finanzas excelentes, el score debería ser alto
	assert.Greater(suite.T(), result, 700) // Score debería ser alto
}

// TestContextCancellation testa el manejo de cancelación de contexto
func (suite *CreditUseCaseTestSuite) TestContextCancellation() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancelar inmediatamente

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI que debería recibir el contexto cancelado
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", context.Canceled)

	// Act
	result, err := suite.useCase.CalculateCreditScore(cancelledCtx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, result)
	assert.Contains(suite.T(), err.Error(), "error calculating credit score")
}

// TestRunSuite ejecuta todos los tests del suite
func TestCreditUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CreditUseCaseTestSuite))
}

// Tests adicionales para casos específicos

// TestCreditPlanWithDifferentScores testa planes con diferentes scores iniciales
func TestCreditPlanWithDifferentScores(t *testing.T) {
	scores := []int{1, 250, 500, 750, 1000}

	for _, score := range scores {
		t.Run(fmt.Sprintf("Score_%d", score), func(t *testing.T) {
			// Arrange
			mockOpenAI := new(mocks.MockOpenAIClient)
			mockCache := new(mocks.MockCacheClient)
			useCase := usecases.NewCreditUseCase(mockOpenAI, mockCache)

			data := mocks.TestFinancialAnalysisData
			data.FinancialScore = score

			expectedPlan := mocks.TestCreditPlan
			expectedPlan.CurrentScore = score
			planJSON, _ := json.Marshal(expectedPlan)

			// Mock cache miss
			mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
				Return([]byte(nil), errors.New("cache miss"))

			// Mock OpenAI response
			mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
				Return(string(planJSON), nil)

			// Mock cache set
			mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
				Return(nil)

			// Act
			result, err := useCase.GenerateImprovementPlan(context.Background(), data)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, score, result.CurrentScore)

			// Verificar expectativas
			mockOpenAI.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}
