package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func intPtr(i int) *int {
	return &i
}

func TestAnalyticsHandlers_GetExpensesSummary_Pagination(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedLimit  int
		expectedOffset int
		totalItems     int
		hasMore        bool
	}{
		{
			name:           "default pagination",
			queryParams:    "",
			expectedLimit:  50,  // valor por defecto
			expectedOffset: 0,   // valor por defecto
			totalItems:     120, // más que una página
			hasMore:        true,
		},
		{
			name:           "custom pagination",
			queryParams:    "limit=10&offset=20",
			expectedLimit:  10,
			expectedOffset: 20,
			totalItems:     45,
			hasMore:        true,
		},
		{
			name:           "last page",
			queryParams:    "limit=50&offset=100",
			expectedLimit:  50,
			expectedOffset: 100,
			totalItems:     120,
			hasMore:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Preparar mocks
			mockExpensesService := &MockExpensesAnalyticsService{}
			mockCategoriesService := &MockCategoriesAnalyticsService{}
			mockIncomesService := &MockIncomesAnalyticsService{}

			// Crear datos de respuesta con paginación
			expectedResponse := &usecases.ExpensesSummary{
				Expenses: make([]usecases.ExpenseItem, tt.expectedLimit), // Array del tamaño del límite
				Summary: usecases.ExpenseSummary{
					TotalAmount: 1000.0,
				},
				Pagination: usecases.PaginationInfo{
					Total:   tt.totalItems,
					Limit:   tt.expectedLimit,
					Offset:  tt.expectedOffset,
					HasMore: tt.hasMore,
				},
			}

			// Configurar mock para verificar parámetros de paginación
			mockExpensesService.On("GetExpensesSummary", mock.Anything, mock.MatchedBy(func(params usecases.ExpensesSummaryParams) bool {
				return params.Pagination.Limit == tt.expectedLimit && params.Pagination.Offset == tt.expectedOffset
			})).Return(expectedResponse, nil)

			// Crear handlers y router
			handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
			router := setupTestRouter(handlers)

			// Crear petición
			url := fmt.Sprintf("/api/v1/expenses/summary?%s", tt.queryParams)
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("X-Caller-ID", "user-123")
			w := httptest.NewRecorder()

			// Ejecutar
			router.ServeHTTP(w, req)

			// Verificar
			assert.Equal(t, http.StatusOK, w.Code)

			var response usecases.ExpensesSummary
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Verificar paginación
			assert.Equal(t, tt.totalItems, response.Pagination.Total)
			assert.Equal(t, tt.expectedLimit, response.Pagination.Limit)
			assert.Equal(t, tt.expectedOffset, response.Pagination.Offset)
			assert.Equal(t, tt.hasMore, response.Pagination.HasMore)

			mockExpensesService.AssertExpectations(t)
		})
	}
}

func TestAnalyticsHandlers_GetExpensesSummary_Filters(t *testing.T) {
	tests := []struct {
		name        string
		queryParams string
		expected    usecases.ExpensesSummaryParams
	}{
		{
			name:        "filter by period",
			queryParams: "year=2024&month=1",
			expected: usecases.ExpensesSummaryParams{
				Period: usecases.DatePeriod{
					Year:  intPtr(2024),
					Month: intPtr(1),
				},
			},
		},
		{
			name:        "sorting by amount desc",
			queryParams: "sort_by=amount&order=desc",
			expected: usecases.ExpensesSummaryParams{
				Sorting: usecases.SortingCriteria{
					Field: usecases.SortByAmount,
					Order: usecases.Descending,
				},
			},
		},
		{
			name:        "sorting by date asc",
			queryParams: "sort_by=date&order=asc",
			expected: usecases.ExpensesSummaryParams{
				Sorting: usecases.SortingCriteria{
					Field: usecases.SortByDate,
					Order: usecases.Ascending,
				},
			},
		},
		{
			name:        "sorting by category",
			queryParams: "sort_by=category&order=desc",
			expected: usecases.ExpensesSummaryParams{
				Sorting: usecases.SortingCriteria{
					Field: usecases.SortByCategory,
					Order: usecases.Descending,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Preparar mocks
			mockExpensesService := &MockExpensesAnalyticsService{}
			mockCategoriesService := &MockCategoriesAnalyticsService{}
			mockIncomesService := &MockIncomesAnalyticsService{}

			// Configurar mock para verificar parámetros de filtrado
			mockExpensesService.On("GetExpensesSummary", mock.Anything, mock.MatchedBy(func(params usecases.ExpensesSummaryParams) bool {
				// Verificar que los parámetros coincidan con los esperados
				if tt.expected.Period.Year != nil && (params.Period.Year == nil || *params.Period.Year != *tt.expected.Period.Year) {
					return false
				}
				if tt.expected.Period.Month != nil && (params.Period.Month == nil || *params.Period.Month != *tt.expected.Period.Month) {
					return false
				}
				if tt.expected.Sorting.Field != "" && params.Sorting.Field != tt.expected.Sorting.Field {
					return false
				}
				if tt.expected.Sorting.Order != "" && params.Sorting.Order != tt.expected.Sorting.Order {
					return false
				}
				return true
			})).Return(&usecases.ExpensesSummary{}, nil)

			// Crear handlers y router
			handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
			router := setupTestRouter(handlers)

			// Crear petición
			url := fmt.Sprintf("/api/v1/expenses/summary?%s", tt.queryParams)
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("X-Caller-ID", "user-123")
			w := httptest.NewRecorder()

			// Ejecutar
			router.ServeHTTP(w, req)

			// Verificar
			assert.Equal(t, http.StatusOK, w.Code)
			mockExpensesService.AssertExpectations(t)
		})
	}
}

func TestAnalyticsHandlers_GetExpensesSummary_InvalidFilters(t *testing.T) {
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
		{
			name:           "invalid pagination limit",
			queryParams:    "limit=1001",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Límite inválido (debe ser entre 1 y 100)",
		},
		{
			name:           "negative offset",
			queryParams:    "offset=-1",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Offset inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Preparar mocks
			mockExpensesService := &MockExpensesAnalyticsService{}
			mockCategoriesService := &MockCategoriesAnalyticsService{}
			mockIncomesService := &MockIncomesAnalyticsService{}

			// No necesitamos configurar el mock para casos de error de validación
			// ya que el error ocurre antes de llamar al servicio

			// Crear handlers y router
			handlers := NewAnalyticsHandlers(mockExpensesService, mockCategoriesService, mockIncomesService)
			router := setupTestRouter(handlers)

			// Crear petición
			url := fmt.Sprintf("/api/v1/expenses/summary?%s", tt.queryParams)
			req := httptest.NewRequest("GET", url, nil)
			req.Header.Set("X-Caller-ID", "user-123")
			w := httptest.NewRecorder()

			// Ejecutar
			router.ServeHTTP(w, req)

			// Verificar
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar el mensaje de error
			var response struct {
				Error string `json:"error"`
			}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedError, response.Error)

			// No debería haber llamadas al servicio en casos de error de validación
			mockExpensesService.AssertNotCalled(t, "GetExpensesSummary")
		})
	}
}
