package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCategoriesAnalytics_Success(t *testing.T) {
	// Arrange
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	expectedResponse := &usecases.CategoriesAnalytics{
		Categories: []usecases.CategoryItem{
			{
				CategoryID:            "cat-1",
				CategoryName:          "Food",
				TotalAmount:           1000.0,
				PercentageOfExpenses:  30.0,
				PercentageOfIncome:    20.0,
				TransactionCount:      10,
				AveragePerTransaction: 100.0,
				ColorSeed:             1,
			},
			{
				CategoryID:            "cat-2",
				CategoryName:          "Transport",
				TotalAmount:           500.0,
				PercentageOfExpenses:  15.0,
				PercentageOfIncome:    10.0,
				TransactionCount:      5,
				AveragePerTransaction: 100.0,
				ColorSeed:             2,
			},
		},
		Summary: usecases.CategorySummary{
			TotalCategories:  2,
			LargestCategory:  "Food",
			SmallestCategory: "Transport",
			TotalAmount:      1500.0,
		},
	}

	mockCategoriesService.On("GetCategoriesAnalytics", mock.Anything, mock.MatchedBy(func(params usecases.CategoriesAnalyticsParams) bool {
		return params.UserID == "user-123"
	})).Return(expectedResponse, nil)

	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Act
	req := httptest.NewRequest("GET", "/api/v1/categories/analytics?year=2024&month=1", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.CategoriesAnalytics
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Categories, 2)
	assert.Equal(t, "Food", response.Categories[0].CategoryName)
	assert.Equal(t, 1000.0, response.Categories[0].TotalAmount)
	assert.Equal(t, 30.0, response.Categories[0].PercentageOfExpenses)
	assert.Equal(t, 2, response.Summary.TotalCategories)
	assert.Equal(t, 1500.0, response.Summary.TotalAmount)

	mockCategoriesService.AssertExpectations(t)
}

func TestGetCategoriesAnalytics_NoAuth(t *testing.T) {
	// Arrange
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Act - Request sin header de autenticación
	req := httptest.NewRequest("GET", "/api/v1/categories/analytics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response struct {
		Error string `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Usuario no autenticado", response.Error)

	mockCategoriesService.AssertNotCalled(t, "GetCategoriesAnalytics")
}

func TestGetCategoriesAnalytics_InvalidParams(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid year",
			queryParams:    "year=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Año inválido",
		},
		{
			name:           "invalid month",
			queryParams:    "month=13",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Mes debe estar entre 1 y 12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockExpensesService := &MockExpensesAnalyticsService{}
			mockCategoriesService := &MockCategoriesAnalyticsService{}
			mockIncomesService := &MockIncomesAnalyticsService{}

			handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
			router := setupTestRouter(handlers)

			// Act
			url := "/api/v1/categories/analytics?" + tt.queryParams
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("X-Caller-ID", "user-123")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response struct {
				Error string `json:"error"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedError, response.Error)

			mockCategoriesService.AssertNotCalled(t, "GetCategoriesAnalytics")
		})
	}
}

func TestGetCategoriesAnalytics_ServiceError(t *testing.T) {
	// Arrange
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	mockCategoriesService.On("GetCategoriesAnalytics", mock.Anything, mock.Anything).
		Return(nil, errors.NewInternalServerError("Error al obtener analytics"))

	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Act
	req := httptest.NewRequest("GET", "/api/v1/categories/analytics", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response struct {
		Error string `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Internal server error", response.Error)

	mockCategoriesService.AssertExpectations(t)
}
