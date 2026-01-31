package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/financial-ai-service/internal/core/ports"
	"github.com/financial-ai-service/internal/core/usecases"
	"github.com/financial-ai-service/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// AnalysisUseCaseTestSuite suite de tests para el caso de uso de análisis
type AnalysisUseCaseTestSuite struct {
	suite.Suite
	mockOpenAI *mocks.MockOpenAIClient
	mockCache  *mocks.MockCacheClient
	useCase    ports.AIAnalysisPort
	ctx        context.Context
}

// SetupTest configura cada test individual
func (suite *AnalysisUseCaseTestSuite) SetupTest() {
	suite.mockOpenAI = new(mocks.MockOpenAIClient)
	suite.mockCache = new(mocks.MockCacheClient)
	suite.useCase = usecases.NewAnalysisUseCase(suite.mockOpenAI, suite.mockCache)
	suite.ctx = context.Background()
}

// TearDownTest verifica que todas las expectativas se cumplieron
func (suite *AnalysisUseCaseTestSuite) TearDownTest() {
	suite.mockOpenAI.AssertExpectations(suite.T())
	suite.mockCache.AssertExpectations(suite.T())
}

// TestAnalyzeFinancialHealth_Success testa el análisis exitoso de salud financiera
func (suite *AnalysisUseCaseTestSuite) TestAnalyzeFinancialHealth_Success() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	expectedResponse := mocks.TestOpenAIResponse

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(expectedResponse, nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.AnalyzeFinancialHealth(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 750, result.Score)
	assert.Equal(suite.T(), "Bueno", result.Level)
	assert.NotEmpty(suite.T(), result.Insights)
	assert.NotZero(suite.T(), result.GeneratedAt)
}

// TestAnalyzeFinancialHealth_CacheHit testa el análisis con cache hit
func (suite *AnalysisUseCaseTestSuite) TestAnalyzeFinancialHealth_CacheHit() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	cachedAnalysis := mocks.TestHealthAnalysis
	cachedData, _ := json.Marshal(cachedAnalysis)

	// Mock cache hit
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return(cachedData, nil)

	// Act
	result, err := suite.useCase.AnalyzeFinancialHealth(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), cachedAnalysis.Score, result.Score)
	assert.Equal(suite.T(), cachedAnalysis.Level, result.Level)

	// Verificar que no se llamó a OpenAI
	suite.mockOpenAI.AssertNotCalled(suite.T(), "GenerateAnalysis")
}

// TestAnalyzeFinancialHealth_OpenAIError testa el manejo de errores de OpenAI
func (suite *AnalysisUseCaseTestSuite) TestAnalyzeFinancialHealth_OpenAIError() {
	// Arrange
	data := mocks.TestFinancialAnalysisData

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI error
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", errors.New("OpenAI API error"))

	// Act
	result, err := suite.useCase.AnalyzeFinancialHealth(suite.ctx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error analyzing financial health")
}

// TestAnalyzeFinancialHealth_InvalidJSON testa el manejo de JSON inválido
func (suite *AnalysisUseCaseTestSuite) TestAnalyzeFinancialHealth_InvalidJSON() {
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
	result, err := suite.useCase.AnalyzeFinancialHealth(suite.ctx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error parsing analysis response")
}

// TestGenerateInsights_Success testa la generación exitosa de insights
func (suite *AnalysisUseCaseTestSuite) TestGenerateInsights_Success() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	expectedInsights := mocks.TestInsights
	insightsJSON, _ := json.Marshal(expectedInsights)

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(string(insightsJSON), nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.GenerateInsights(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "Optimización de gastos", result[0].Title)
	assert.Equal(suite.T(), "Incremento de ahorros", result[1].Title)
}

// TestGenerateInsights_CacheHit testa la generación con cache hit
func (suite *AnalysisUseCaseTestSuite) TestGenerateInsights_CacheHit() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	cachedInsights := mocks.TestInsights
	cachedData, _ := json.Marshal(cachedInsights)

	// Mock cache hit
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return(cachedData, nil)

	// Act
	result, err := suite.useCase.GenerateInsights(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result, 2)

	// Verificar que no se llamó a OpenAI
	suite.mockOpenAI.AssertNotCalled(suite.T(), "GenerateAnalysis")
}

// TestGenerateInsights_OpenAIError testa el manejo de errores de OpenAI
func (suite *AnalysisUseCaseTestSuite) TestGenerateInsights_OpenAIError() {
	// Arrange
	data := mocks.TestFinancialAnalysisData

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI error
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", errors.New("OpenAI API error"))

	// Act
	result, err := suite.useCase.GenerateInsights(suite.ctx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error generating insights")
}

// TestGenerateInsights_InvalidJSON testa el manejo de JSON inválido
func (suite *AnalysisUseCaseTestSuite) TestGenerateInsights_InvalidJSON() {
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
	result, err := suite.useCase.GenerateInsights(suite.ctx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error parsing insights response")
}

// TestCacheFailure_DoesNotAffectOperation testa que los errores de cache no afectan la operación
func (suite *AnalysisUseCaseTestSuite) TestCacheFailure_DoesNotAffectOperation() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	expectedResponse := mocks.TestOpenAIResponse

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(expectedResponse, nil)

	// Mock cache set failure (no debería afectar la operación)
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(errors.New("cache set failed"))

	// Act
	result, err := suite.useCase.AnalyzeFinancialHealth(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err) // El error de cache no debería afectar el resultado
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 750, result.Score)
}

// TestContextCancellation testa el manejo de cancelación de contexto
func (suite *AnalysisUseCaseTestSuite) TestContextCancellation() {
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
	result, err := suite.useCase.AnalyzeFinancialHealth(cancelledCtx, data)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "error analyzing financial health")
}

// TestAnalysisWithDifferentPeriods testa análisis con diferentes períodos
func (suite *AnalysisUseCaseTestSuite) TestAnalysisWithDifferentPeriods() {
	periods := []string{"monthly", "quarterly", "yearly"}

	for _, period := range periods {
		suite.Run("Period_"+period, func() {
			// Arrange
			data := mocks.TestFinancialAnalysisData
			data.Period = period
			expectedResponse := mocks.TestOpenAIResponse

			// Mock cache miss
			suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
				Return([]byte(nil), errors.New("cache miss"))

			// Mock OpenAI response
			suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
				Return(expectedResponse, nil)

			// Mock cache set
			suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
				Return(nil)

			// Act
			result, err := suite.useCase.AnalyzeFinancialHealth(suite.ctx, data)

			// Assert
			assert.NoError(suite.T(), err)
			assert.NotNil(suite.T(), result)
			assert.Equal(suite.T(), 750, result.Score)
		})
	}
}

// TestAnalysisWithEmptyCategories testa análisis con categorías vacías
func (suite *AnalysisUseCaseTestSuite) TestAnalysisWithEmptyCategories() {
	// Arrange
	data := mocks.TestFinancialAnalysisData
	data.ExpensesByCategory = map[string]float64{} // Categorías vacías
	expectedResponse := mocks.TestOpenAIResponse

	// Mock cache miss
	suite.mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	suite.mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(expectedResponse, nil)

	// Mock cache set
	suite.mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := suite.useCase.AnalyzeFinancialHealth(suite.ctx, data)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	// Debería funcionar incluso sin categorías
}

// TestRunSuite ejecuta todos los tests del suite
func TestAnalysisUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AnalysisUseCaseTestSuite))
}

// Tests adicionales para casos específicos

// TestAnalysisWithZeroIncome testa análisis con ingresos cero
func TestAnalysisWithZeroIncome(t *testing.T) {
	// Arrange
	mockOpenAI := new(mocks.MockOpenAIClient)
	mockCache := new(mocks.MockCacheClient)
	useCase := usecases.NewAnalysisUseCase(mockOpenAI, mockCache)

	data := mocks.TestFinancialAnalysisData
	data.TotalIncome = 0 // Ingresos cero

	expectedResponse := mocks.TestOpenAIResponse

	// Mock cache miss
	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(expectedResponse, nil)

	// Mock cache set
	mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := useCase.AnalyzeFinancialHealth(context.Background(), data)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verificar expectativas
	mockOpenAI.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

// TestAnalysisWithNegativeSavingsRate testa análisis con tasa de ahorro negativa
func TestAnalysisWithNegativeSavingsRate(t *testing.T) {
	// Arrange
	mockOpenAI := new(mocks.MockOpenAIClient)
	mockCache := new(mocks.MockCacheClient)
	useCase := usecases.NewAnalysisUseCase(mockOpenAI, mockCache)

	data := mocks.TestFinancialAnalysisData
	data.SavingsRate = -0.1 // Tasa de ahorro negativa

	expectedResponse := mocks.TestOpenAIResponse

	// Mock cache miss
	mockCache.On("Get", mock.Anything, mock.AnythingOfType("string")).
		Return([]byte(nil), errors.New("cache miss"))

	// Mock OpenAI response
	mockOpenAI.On("GenerateAnalysis", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(expectedResponse, nil)

	// Mock cache set
	mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("time.Duration")).
		Return(nil)

	// Act
	result, err := useCase.AnalyzeFinancialHealth(context.Background(), data)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verificar expectativas
	mockOpenAI.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
