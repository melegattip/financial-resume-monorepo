package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks para los servicios
type MockExpensesAnalyticsService struct {
	mock.Mock
}

func (m *MockExpensesAnalyticsService) GetExpensesSummary(ctx context.Context, params usecases.ExpensesSummaryParams) (*usecases.ExpensesSummary, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.ExpensesSummary), args.Error(1)
}

type MockCategoriesAnalyticsService struct {
	mock.Mock
}

func (m *MockCategoriesAnalyticsService) GetCategoriesAnalytics(ctx context.Context, params usecases.CategoriesAnalyticsParams) (*usecases.CategoriesAnalytics, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.CategoriesAnalytics), args.Error(1)
}

type MockIncomesAnalyticsService struct {
	mock.Mock
}

func (m *MockIncomesAnalyticsService) GetIncomesSummary(ctx context.Context, params usecases.IncomesSummaryParams) (*usecases.IncomesSummary, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.IncomesSummary), args.Error(1)
}

func setupTestRouter(handlers *AnalyticsHandlers) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Middleware simulado para tests
	router.Use(func(c *gin.Context) {
		userID := c.GetHeader("X-Caller-ID")
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("/expenses/summary", handlers.GetExpensesSummary)
		v1.GET("/categories/analytics", handlers.GetCategoriesAnalytics)
		v1.GET("/incomes/summary", handlers.GetIncomesSummary)
	}

	return router
}

func TestAnalyticsHandlers_GetExpensesSummary_Success(t *testing.T) {
	// Preparar mocks
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	// Datos de respuesta esperados
	expectedResponse := &usecases.ExpensesSummary{
		Expenses: []usecases.ExpenseItem{
			{
				ID:           "exp-1",
				Description:  "Test Expense",
				Amount:       100.0,
				CategoryName: "Food",
			},
		},
		Summary: usecases.ExpenseSummary{
			TotalAmount: 100.0,
		},
		Pagination: usecases.PaginationInfo{
			Total:  1,
			Limit:  50,
			Offset: 0,
		},
	}

	// Configurar mock
	mockExpensesService.On("GetExpensesSummary", mock.Anything, mock.MatchedBy(func(params usecases.ExpensesSummaryParams) bool {
		return params.UserID == "user-123"
	})).Return(expectedResponse, nil)

	// Crear handlers
	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Crear petición
	req := httptest.NewRequest("GET", "/api/v1/expenses/summary?year=2024&month=12", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.ExpensesSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "exp-1", response.Expenses[0].ID)
	assert.Equal(t, 100.0, response.Summary.TotalAmount)

	mockExpensesService.AssertExpectations(t)
}

func TestAnalyticsHandlers_GetExpensesSummary_MissingUserID(t *testing.T) {
	// Preparar mocks
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	// Crear handlers
	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Crear petición sin header X-Caller-ID
	req := httptest.NewRequest("GET", "/api/v1/expenses/summary", nil)
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Usuario no autenticado")
}

func TestAnalyticsHandlers_GetExpensesSummary_InvalidParameters(t *testing.T) {
	tests := []struct {
		name        string
		queryParam  string
		expectedMsg string
	}{
		{
			name:        "invalid year",
			queryParam:  "year=invalid",
			expectedMsg: "Año inválido",
		},
		{
			name:        "invalid month",
			queryParam:  "month=invalid",
			expectedMsg: "Mes inválido",
		},
		{
			name:        "invalid sort_by",
			queryParam:  "sort_by=invalid",
			expectedMsg: "Campo de ordenamiento inválido",
		},
		{
			name:        "invalid order",
			queryParam:  "order=invalid",
			expectedMsg: "Orden inválido",
		},
		{
			name:        "invalid limit",
			queryParam:  "limit=101",
			expectedMsg: "Límite inválido",
		},
		{
			name:        "invalid offset",
			queryParam:  "offset=-1",
			expectedMsg: "Offset inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Preparar mocks
			mockExpensesService := &MockExpensesAnalyticsService{}
			mockCategoriesService := &MockCategoriesAnalyticsService{}
			mockIncomesService := &MockIncomesAnalyticsService{}

			// Crear handlers
			handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
			router := setupTestRouter(handlers)

			// Crear petición
			req := httptest.NewRequest("GET", "/api/v1/expenses/summary?"+tt.queryParam, nil)
			req.Header.Set("X-Caller-ID", "user-123")
			w := httptest.NewRecorder()

			// Ejecutar
			router.ServeHTTP(w, req)

			// Verificar
			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedMsg)
		})
	}
}

func TestAnalyticsHandlers_GetExpensesSummary_ServiceError(t *testing.T) {
	// Preparar mocks
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	// Configurar mock para devolver error
	mockExpensesService.On("GetExpensesSummary", mock.Anything, mock.Anything).Return(nil, errors.NewInternalServerError("service error"))

	// Crear handlers
	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Crear petición
	req := httptest.NewRequest("GET", "/api/v1/expenses/summary", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockExpensesService.AssertExpectations(t)
}

func TestAnalyticsHandlers_GetCategoriesAnalytics_Success(t *testing.T) {
	// Preparar mocks
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	// Datos de respuesta esperados
	expectedResponse := &usecases.CategoriesAnalytics{
		Categories: []usecases.CategoryItem{
			{
				CategoryID:   "cat-1",
				CategoryName: "Food",
				TotalAmount:  100.0,
			},
		},
		Summary: usecases.CategorySummary{
			TotalCategories: 1,
			TotalAmount:     100.0,
		},
	}

	// Configurar mock
	mockCategoriesService.On("GetCategoriesAnalytics", mock.Anything, mock.MatchedBy(func(params usecases.CategoriesAnalyticsParams) bool {
		return params.UserID == "user-123"
	})).Return(expectedResponse, nil)

	// Crear handlers
	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Crear petición
	req := httptest.NewRequest("GET", "/api/v1/categories/analytics?year=2024", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.CategoriesAnalytics
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "cat-1", response.Categories[0].CategoryID)
	assert.Equal(t, 100.0, response.Summary.TotalAmount)

	mockCategoriesService.AssertExpectations(t)
}

func TestAnalyticsHandlers_GetIncomesSummary_Success(t *testing.T) {
	// Preparar mocks
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	// Datos de respuesta esperados
	expectedResponse := &usecases.IncomesSummary{
		Incomes: []usecases.IncomeItem{
			{
				ID:           "inc-1",
				Description:  "Salary",
				Amount:       2000.0,
				CategoryName: "Work",
			},
		},
		Summary: usecases.IncomeSummary{
			TotalAmount:      2000.0,
			TransactionCount: 1,
		},
	}

	// Configurar mock
	mockIncomesService.On("GetIncomesSummary", mock.Anything, mock.MatchedBy(func(params usecases.IncomesSummaryParams) bool {
		return params.UserID == "user-123"
	})).Return(expectedResponse, nil)

	// Crear handlers
	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Crear petición
	req := httptest.NewRequest("GET", "/api/v1/incomes/summary?sort_by=amount&order=desc", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.IncomesSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "inc-1", response.Incomes[0].ID)
	assert.Equal(t, 2000.0, response.Summary.TotalAmount)

	mockIncomesService.AssertExpectations(t)
}

func TestNewAnalyticsHandlers(t *testing.T) {
	// Preparar mocks
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	// Crear handlers
	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)

	// Verificar
	assert.NotNil(t, handlers)
	assert.Equal(t, mockExpensesService, handlers.expensesService)
	assert.Equal(t, mockCategoriesService, handlers.categoriesService)
	assert.Equal(t, mockIncomesService, handlers.incomesService)
}
