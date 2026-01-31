package dashboard

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

// Mock del servicio de dashboard
type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetDashboardOverview(ctx context.Context, params usecases.DashboardParams) (*usecases.DashboardOverview, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.DashboardOverview), args.Error(1)
}

func setupTestRouter(handler *Handler) *gin.Engine {
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
		v1.GET("/dashboard", handler.GetDashboard)
	}

	return router
}

func TestHandler_GetDashboard_Success(t *testing.T) {
	// Preparar mock service
	mockService := &MockDashboardService{}

	// Datos de respuesta esperados
	expectedResponse := &usecases.DashboardOverview{
		Period: usecases.PeriodInfo{
			Year:  "2024",
			Label: "2024",
		},
		Metrics: usecases.FinancialMetrics{
			TotalIncome:   2000.0,
			TotalExpenses: 1000.0,
			Balance:       1000.0,
		},
		Counts: usecases.TransactionCounts{
			IncomeTransactions:  2,
			ExpenseTransactions: 5,
		},
	}

	// Configurar mock
	mockService.On("GetDashboardOverview", mock.Anything, mock.MatchedBy(func(params usecases.DashboardParams) bool {
		return params.UserID == "user-123"
	})).Return(expectedResponse, nil)

	// Crear handler
	handler := NewHandler(mockService, nil)
	router := setupTestRouter(handler)

	// Crear petición
	req := httptest.NewRequest("GET", "/api/v1/dashboard?year=2024", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.DashboardOverview
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2000.0, response.Metrics.TotalIncome)
	assert.Equal(t, 1000.0, response.Metrics.TotalExpenses)
	assert.Equal(t, 1000.0, response.Metrics.Balance)

	mockService.AssertExpectations(t)
}

func TestHandler_GetDashboard_SuccessWithMonthFilter(t *testing.T) {
	// Preparar mock service
	mockService := &MockDashboardService{}

	// Datos de respuesta esperados
	expectedResponse := &usecases.DashboardOverview{
		Period: usecases.PeriodInfo{
			Year:  "2024",
			Month: "12",
			Label: "December 2024",
		},
		Metrics: usecases.FinancialMetrics{
			TotalIncome:   500.0,
			TotalExpenses: 300.0,
			Balance:       200.0,
		},
	}

	// Configurar mock para verificar que se pasan year y month
	mockService.On("GetDashboardOverview", mock.Anything, mock.MatchedBy(func(params usecases.DashboardParams) bool {
		return params.UserID == "user-123" &&
			params.Period.Year != nil && *params.Period.Year == 2024 &&
			params.Period.Month != nil && *params.Period.Month == 12
	})).Return(expectedResponse, nil)

	// Crear handler
	handler := NewHandler(mockService, nil)
	router := setupTestRouter(handler)

	// Crear petición con year y month
	req := httptest.NewRequest("GET", "/api/v1/dashboard?year=2024&month=12", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.DashboardOverview
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "December 2024", response.Period.Label)

	mockService.AssertExpectations(t)
}

func TestHandler_GetDashboard_MissingUserID(t *testing.T) {
	// Preparar mock service
	mockService := &MockDashboardService{}

	// Crear handler
	handler := NewHandler(mockService, nil)
	router := setupTestRouter(handler)

	// Crear petición sin header X-Caller-ID
	req := httptest.NewRequest("GET", "/api/v1/dashboard", nil)
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Usuario no autenticado")

	// El servicio no debería ser llamado
	mockService.AssertNotCalled(t, "GetDashboardOverview")
}

func TestHandler_GetDashboard_InvalidYear(t *testing.T) {
	// Preparar mock service
	mockService := &MockDashboardService{}

	// Crear handler
	handler := NewHandler(mockService, nil)
	router := setupTestRouter(handler)

	// Crear petición con año inválido
	req := httptest.NewRequest("GET", "/api/v1/dashboard?year=invalid", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Año inválido")

	// El servicio no debería ser llamado
	mockService.AssertNotCalled(t, "GetDashboardOverview")
}

func TestHandler_GetDashboard_InvalidMonth(t *testing.T) {
	// Preparar mock service
	mockService := &MockDashboardService{}

	// Crear handler
	handler := NewHandler(mockService, nil)
	router := setupTestRouter(handler)

	// Crear petición con mes inválido
	req := httptest.NewRequest("GET", "/api/v1/dashboard?month=invalid", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Mes inválido")

	// El servicio no debería ser llamado
	mockService.AssertNotCalled(t, "GetDashboardOverview")
}

func TestHandler_GetDashboard_ServiceError(t *testing.T) {
	// Preparar mock service
	mockService := &MockDashboardService{}

	// Configurar mock para devolver error
	mockService.On("GetDashboardOverview", mock.Anything, mock.Anything).Return(nil, errors.NewInternalServerError("database error"))

	// Crear handler
	handler := NewHandler(mockService, nil)
	router := setupTestRouter(handler)

	// Crear petición
	req := httptest.NewRequest("GET", "/api/v1/dashboard", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()

	// Ejecutar
	router.ServeHTTP(w, req)

	// Verificar que se devuelve error del servidor
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

func TestNewHandler(t *testing.T) {
	// Preparar mock service
	mockService := &MockDashboardService{}

	// Crear handler
	handler := NewHandler(mockService, nil)

	// Verificar
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}
