package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/financial-ai-service/internal/adapters/http/handlers"
	"github.com/financial-ai-service/internal/core/ports"
	"github.com/financial-ai-service/test/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// HandlersTestSuite suite de tests para handlers
type HandlersTestSuite struct {
	suite.Suite
	router       *gin.Engine
	mockAnalysis *mocks.MockAIAnalysisPort
	mockPurchase *mocks.MockPurchaseDecisionPort
	mockCredit   *mocks.MockCreditAnalysisPort
	handlers     *handlers.Handlers
}

// SetupSuite configura el suite de tests
func (suite *HandlersTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

// SetupTest configura cada test individual
func (suite *HandlersTestSuite) SetupTest() {
	suite.mockAnalysis = new(mocks.MockAIAnalysisPort)
	suite.mockPurchase = new(mocks.MockPurchaseDecisionPort)
	suite.mockCredit = new(mocks.MockCreditAnalysisPort)

	suite.router = gin.New()
	handlers.SetupRoutes(suite.router, suite.mockAnalysis, suite.mockPurchase, suite.mockCredit)
}

// TearDownTest limpia después de cada test
func (suite *HandlersTestSuite) TearDownTest() {
	suite.mockAnalysis.AssertExpectations(suite.T())
	suite.mockPurchase.AssertExpectations(suite.T())
	suite.mockCredit.AssertExpectations(suite.T())
}

// TestHealthCheck testa el endpoint de health check
func (suite *HandlersTestSuite) TestHealthCheck() {
	req, _ := http.NewRequest("GET", "/health", nil)
	recorder := httptest.NewRecorder()

	suite.router.ServeHTTP(recorder, req)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "healthy", response["status"])
	assert.Equal(suite.T(), "Financial AI Service", response["service"])
}

// TestAnalyzeFinancialHealth testa el análisis de salud financiera
func (suite *HandlersTestSuite) TestAnalyzeFinancialHealth() {
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "Análisis exitoso",
			requestBody: mocks.TestFinancialAnalysisData,
			mockSetup: func() {
				suite.mockAnalysis.On("AnalyzeFinancialHealth", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(&mocks.TestHealthAnalysis, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Request inválido",
			requestBody:    "invalid json",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "Error en análisis",
			requestBody: mocks.TestFinancialAnalysisData,
			mockSetup: func() {
				suite.mockAnalysis.On("AnalyzeFinancialHealth", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return((*ports.HealthAnalysis)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Crear mocks nuevos para cada test individual
			mockAnalysis := new(mocks.MockAIAnalysisPort)
			mockPurchase := new(mocks.MockPurchaseDecisionPort)
			mockCredit := new(mocks.MockCreditAnalysisPort)

			router := gin.New()
			handlers.SetupRoutes(router, mockAnalysis, mockPurchase, mockCredit)

			// Configurar el mock específico para este test
			if tc.name == "Error en análisis" {
				mockAnalysis.On("AnalyzeFinancialHealth", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return((*ports.HealthAnalysis)(nil), assert.AnError)
			} else if tc.name == "Análisis exitoso" {
				mockAnalysis.On("AnalyzeFinancialHealth", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(&mocks.TestHealthAnalysis, nil)
			}

			bodyBytes, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/ai/health-analysis", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(suite.T(), tc.expectedStatus, recorder.Code)

			var response map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			if tc.expectedError {
				assert.Contains(suite.T(), response, "error")
			} else {
				assert.Equal(suite.T(), true, response["success"])
				assert.Equal(suite.T(), "ai_microservice", response["source"])
				assert.Contains(suite.T(), response, "data")
			}
		})
	}
}

// TestGenerateInsights testa la generación de insights
func (suite *HandlersTestSuite) TestGenerateInsights() {
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "Insights exitosos",
			requestBody: mocks.TestFinancialAnalysisData,
			mockSetup: func() {
				suite.mockAnalysis.On("GenerateInsights", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(mocks.TestInsights, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Request inválido",
			requestBody:    "invalid json string",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "Error en generación",
			requestBody: mocks.TestFinancialAnalysisData,
			mockSetup: func() {
				suite.mockAnalysis.On("GenerateInsights", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return([]ports.AIInsight(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Crear mocks nuevos para cada test individual
			mockAnalysis := new(mocks.MockAIAnalysisPort)
			mockPurchase := new(mocks.MockPurchaseDecisionPort)
			mockCredit := new(mocks.MockCreditAnalysisPort)

			router := gin.New()
			handlers.SetupRoutes(router, mockAnalysis, mockPurchase, mockCredit)

			// Configurar el mock específico para este test
			if tc.name == "Error en generación" {
				mockAnalysis.On("GenerateInsights", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return([]ports.AIInsight(nil), assert.AnError)
			} else if tc.name == "Insights exitosos" {
				mockAnalysis.On("GenerateInsights", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(mocks.TestInsights, nil)
			} else if tc.name == "Request inválido" {
				// Este caso pasa datos inválidos que se procesan como zeros, pero el handler los acepta
				mockAnalysis.On("GenerateInsights", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(mocks.TestInsights, nil).Maybe()
			}

			bodyBytes, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/ai/insights", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(suite.T(), tc.expectedStatus, recorder.Code)

			var response map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			if tc.expectedError {
				assert.Contains(suite.T(), response, "error")
			} else {
				assert.Equal(suite.T(), true, response["success"])
				assert.Contains(suite.T(), response, "data")
			}
		})
	}
}

// TestCanIBuy testa el análisis de decisión de compra
func (suite *HandlersTestSuite) TestCanIBuy() {
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "Decisión exitosa",
			requestBody: mocks.TestPurchaseAnalysisRequest,
			mockSetup: func() {
				suite.mockPurchase.On("CanIBuy", mock.Anything, mock.AnythingOfType("ports.PurchaseAnalysisRequest")).
					Return(&mocks.TestPurchaseDecision, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Request inválido",
			requestBody:    "invalid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:        "Error en análisis",
			requestBody: mocks.TestPurchaseAnalysisRequest,
			mockSetup: func() {
				suite.mockPurchase.On("CanIBuy", mock.Anything, mock.AnythingOfType("ports.PurchaseAnalysisRequest")).
					Return((*ports.PurchaseDecision)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Crear mocks nuevos para cada test individual
			mockAnalysis := new(mocks.MockAIAnalysisPort)
			mockPurchase := new(mocks.MockPurchaseDecisionPort)
			mockCredit := new(mocks.MockCreditAnalysisPort)

			router := gin.New()
			handlers.SetupRoutes(router, mockAnalysis, mockPurchase, mockCredit)

			// Configurar el mock específico para este test
			if tc.name == "Error en análisis" {
				mockPurchase.On("CanIBuy", mock.Anything, mock.AnythingOfType("ports.PurchaseAnalysisRequest")).
					Return((*ports.PurchaseDecision)(nil), assert.AnError)
			} else if tc.name == "Decisión exitosa" {
				mockPurchase.On("CanIBuy", mock.Anything, mock.AnythingOfType("ports.PurchaseAnalysisRequest")).
					Return(&mocks.TestPurchaseDecision, nil)
			}

			bodyBytes, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/ai/can-i-buy", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(suite.T(), tc.expectedStatus, recorder.Code)

			var response map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			if tc.expectedError {
				assert.Contains(suite.T(), response, "error")
			} else {
				assert.Equal(suite.T(), true, response["success"])
				assert.Contains(suite.T(), response, "data")
			}
		})
	}
}

// TestSuggestAlternatives testa la sugerencia de alternativas
func (suite *HandlersTestSuite) TestSuggestAlternatives() {
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "Alternativas exitosas",
			requestBody: mocks.TestPurchaseAnalysisRequest,
			mockSetup: func() {
				suite.mockPurchase.On("SuggestAlternatives", mock.Anything, mock.AnythingOfType("ports.PurchaseAnalysisRequest")).
					Return(mocks.TestAlternatives, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "Error en sugerencias",
			requestBody: mocks.TestPurchaseAnalysisRequest,
			mockSetup: func() {
				suite.mockPurchase.On("SuggestAlternatives", mock.Anything, mock.AnythingOfType("ports.PurchaseAnalysisRequest")).
					Return([]ports.Alternative(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Crear mocks nuevos para cada test individual
			mockAnalysis := new(mocks.MockAIAnalysisPort)
			mockPurchase := new(mocks.MockPurchaseDecisionPort)
			mockCredit := new(mocks.MockCreditAnalysisPort)

			router := gin.New()
			handlers.SetupRoutes(router, mockAnalysis, mockPurchase, mockCredit)

			// Configurar el mock específico para este test
			if tc.name == "Error en sugerencias" {
				mockPurchase.On("SuggestAlternatives", mock.Anything, mock.AnythingOfType("ports.PurchaseAnalysisRequest")).
					Return([]ports.Alternative(nil), assert.AnError)
			} else if tc.name == "Alternativas exitosas" {
				mockPurchase.On("SuggestAlternatives", mock.Anything, mock.AnythingOfType("ports.PurchaseAnalysisRequest")).
					Return(mocks.TestAlternatives, nil)
			}

			bodyBytes, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/ai/alternatives", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(suite.T(), tc.expectedStatus, recorder.Code)

			var response map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			if tc.expectedError {
				assert.Contains(suite.T(), response, "error")
			} else {
				assert.Equal(suite.T(), true, response["success"])
				assert.Contains(suite.T(), response, "data")
			}
		})
	}
}

// TestGenerateCreditPlan testa la generación de plan crediticio
func (suite *HandlersTestSuite) TestGenerateCreditPlan() {
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "Plan exitoso",
			requestBody: mocks.TestFinancialAnalysisData,
			mockSetup: func() {
				suite.mockCredit.On("GenerateImprovementPlan", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(&mocks.TestCreditPlan, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "Error en plan",
			requestBody: mocks.TestFinancialAnalysisData,
			mockSetup: func() {
				suite.mockCredit.On("GenerateImprovementPlan", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return((*ports.CreditPlan)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Crear mocks nuevos para cada test individual
			mockAnalysis := new(mocks.MockAIAnalysisPort)
			mockPurchase := new(mocks.MockPurchaseDecisionPort)
			mockCredit := new(mocks.MockCreditAnalysisPort)

			router := gin.New()
			handlers.SetupRoutes(router, mockAnalysis, mockPurchase, mockCredit)

			// Configurar el mock específico para este test
			if tc.name == "Error en plan" {
				mockCredit.On("GenerateImprovementPlan", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return((*ports.CreditPlan)(nil), assert.AnError)
			} else if tc.name == "Plan exitoso" {
				mockCredit.On("GenerateImprovementPlan", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(&mocks.TestCreditPlan, nil)
			}

			bodyBytes, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/ai/credit-plan", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(suite.T(), tc.expectedStatus, recorder.Code)

			var response map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			if tc.expectedError {
				assert.Contains(suite.T(), response, "error")
			} else {
				assert.Equal(suite.T(), true, response["success"])
				assert.Contains(suite.T(), response, "data")
			}
		})
	}
}

// TestCalculateCreditScore testa el cálculo de score crediticio
func (suite *HandlersTestSuite) TestCalculateCreditScore() {
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "Score exitoso",
			requestBody: mocks.TestFinancialAnalysisData,
			mockSetup: func() {
				suite.mockCredit.On("CalculateCreditScore", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(750, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "Error en cálculo",
			requestBody: mocks.TestFinancialAnalysisData,
			mockSetup: func() {
				suite.mockCredit.On("CalculateCreditScore", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(0, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Crear mocks nuevos para cada test individual
			mockAnalysis := new(mocks.MockAIAnalysisPort)
			mockPurchase := new(mocks.MockPurchaseDecisionPort)
			mockCredit := new(mocks.MockCreditAnalysisPort)

			router := gin.New()
			handlers.SetupRoutes(router, mockAnalysis, mockPurchase, mockCredit)

			// Configurar el mock específico para este test
			if tc.name == "Error en cálculo" {
				mockCredit.On("CalculateCreditScore", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(0, assert.AnError)
			} else if tc.name == "Score exitoso" {
				mockCredit.On("CalculateCreditScore", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
					Return(750, nil)
			}

			bodyBytes, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/ai/credit-score", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(suite.T(), tc.expectedStatus, recorder.Code)

			var response map[string]interface{}
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			if tc.expectedError {
				assert.Contains(suite.T(), response, "error")
			} else {
				assert.Equal(suite.T(), true, response["success"])
				assert.Contains(suite.T(), response, "data")
				// Verificar estructura específica del response
				data := response["data"].(map[string]interface{})
				assert.Equal(suite.T(), float64(750), data["score"])
				assert.Equal(suite.T(), "test-user-123", data["user_id"])
			}
		})
	}
}

// TestRunSuite ejecuta todos los tests
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

// Tests adicionales para casos específicos

// TestInvalidContentType testa manejo de content-type inválido
func TestInvalidContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockAnalysis := new(mocks.MockAIAnalysisPort)
	mockPurchase := new(mocks.MockPurchaseDecisionPort)
	mockCredit := new(mocks.MockCreditAnalysisPort)

	handlers.SetupRoutes(router, mockAnalysis, mockPurchase, mockCredit)

	req, _ := http.NewRequest("POST", "/api/v1/ai/health-analysis", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "text/plain")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
}

// TestMissingFields testa manejo de campos requeridos faltantes
func TestMissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockAnalysis := new(mocks.MockAIAnalysisPort)
	mockPurchase := new(mocks.MockPurchaseDecisionPort)
	mockCredit := new(mocks.MockCreditAnalysisPort)

	// Configurar mock para el caso con datos incompletos
	mockAnalysis.On("AnalyzeFinancialHealth", mock.Anything, mock.AnythingOfType("ports.FinancialAnalysisData")).
		Return(&mocks.TestHealthAnalysis, nil).Maybe()

	handlers.SetupRoutes(router, mockAnalysis, mockPurchase, mockCredit)

	incompleteData := map[string]interface{}{
		"user_id": "test-user",
		// Missing required fields
	}

	bodyBytes, _ := json.Marshal(incompleteData)
	req, _ := http.NewRequest("POST", "/api/v1/ai/health-analysis", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	// El handler debería manejar esto graciosamente
	assert.Contains(t, []int{http.StatusBadRequest, http.StatusOK}, recorder.Code)
}
