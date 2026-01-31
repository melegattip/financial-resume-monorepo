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

func TestGetIncomesSummary_Success(t *testing.T) {
	// Arrange
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	expectedResponse := &usecases.IncomesSummary{
		Incomes: []usecases.IncomeItem{
			{
				ID:           "inc-1",
				Description:  "Salary",
				Amount:       5000.0,
				CategoryID:   "cat-1",
				CategoryName: "Work",
				CreatedAt:    "2024-01-01T00:00:00Z",
			},
			{
				ID:           "inc-2",
				Description:  "Freelance",
				Amount:       1000.0,
				CategoryID:   "cat-2",
				CategoryName: "Extra",
				CreatedAt:    "2024-01-15T00:00:00Z",
			},
		},
		Summary: usecases.IncomeSummary{
			TotalAmount:      6000.0,
			TransactionCount: 2,
		},
	}

	mockIncomesService.On("GetIncomesSummary", mock.Anything, mock.MatchedBy(func(params usecases.IncomesSummaryParams) bool {
		return params.UserID == "user-123" &&
			params.Sorting.Field == usecases.SortByDate &&
			params.Sorting.Order == usecases.Descending
	})).Return(expectedResponse, nil)

	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Act
	req := httptest.NewRequest("GET", "/api/v1/incomes/summary?year=2024&month=1", nil)
	req.Header.Set("X-Caller-ID", "user-123")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response usecases.IncomesSummary
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Incomes, 2)
	assert.Equal(t, "Salary", response.Incomes[0].Description)
	assert.Equal(t, 5000.0, response.Incomes[0].Amount)
	assert.Equal(t, 6000.0, response.Summary.TotalAmount)
	assert.Equal(t, 2, response.Summary.TransactionCount)

	mockIncomesService.AssertExpectations(t)
}

func TestGetIncomesSummary_WithSorting(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		sortField   usecases.SortField
		sortOrder   usecases.SortOrder
	}{
		{
			name:        "sort by amount ascending",
			queryParams: "sort_by=amount&order=asc",
			sortField:   usecases.SortByAmount,
			sortOrder:   usecases.Ascending,
		},
		{
			name:        "sort by date descending",
			queryParams: "sort_by=date&order=desc",
			sortField:   usecases.SortByDate,
			sortOrder:   usecases.Descending,
		},
		{
			name:        "sort by category ascending",
			queryParams: "sort_by=category&order=asc",
			sortField:   usecases.SortByCategory,
			sortOrder:   usecases.Ascending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockExpensesService := &MockExpensesAnalyticsService{}
			mockCategoriesService := &MockCategoriesAnalyticsService{}
			mockIncomesService := &MockIncomesAnalyticsService{}

			mockIncomesService.On("GetIncomesSummary", mock.Anything, mock.MatchedBy(func(params usecases.IncomesSummaryParams) bool {
				return params.UserID == "user-123" &&
					params.Sorting.Field == tt.sortField &&
					params.Sorting.Order == tt.sortOrder
			})).Return(&usecases.IncomesSummary{}, nil)

			handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
			router := setupTestRouter(handlers)

			// Act
			url := "/api/v1/incomes/summary?" + tt.queryParams
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("X-Caller-ID", "user-123")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)
			mockIncomesService.AssertExpectations(t)
		})
	}
}

func TestGetIncomesSummary_NoAuth(t *testing.T) {
	// Arrange
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Act - Request sin header de autenticación
	req := httptest.NewRequest("GET", "/api/v1/incomes/summary", nil)
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

	mockIncomesService.AssertNotCalled(t, "GetIncomesSummary")
}

func TestGetIncomesSummary_InvalidParams(t *testing.T) {
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
		{
			name:           "invalid sort field",
			queryParams:    "sort_by=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Campo de ordenamiento inválido",
		},
		{
			name:           "invalid sort order",
			queryParams:    "order=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Orden inválido",
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
			url := "/api/v1/incomes/summary?" + tt.queryParams
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

			mockIncomesService.AssertNotCalled(t, "GetIncomesSummary")
		})
	}
}

func TestGetIncomesSummary_ServiceError(t *testing.T) {
	// Arrange
	mockExpensesService := &MockExpensesAnalyticsService{}
	mockCategoriesService := &MockCategoriesAnalyticsService{}
	mockIncomesService := &MockIncomesAnalyticsService{}

	mockIncomesService.On("GetIncomesSummary", mock.Anything, mock.Anything).
		Return(nil, errors.NewInternalServerError("Error al obtener resumen"))

	handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
	router := setupTestRouter(handlers)

	// Act
	req := httptest.NewRequest("GET", "/api/v1/incomes/summary", nil)
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

	mockIncomesService.AssertExpectations(t)
}
