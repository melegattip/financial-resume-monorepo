package analytics

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

func TestExpensesHandler_GetExpensesSummary(t *testing.T) {
	tests := []struct {
		name           string
		callerID       string
		queryParams    string
		mockSetup      func(*MockExpensesAnalytics)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "obtener resumen exitosamente",
			callerID:    "user123",
			queryParams: "year=2024&month=3&sort_by=amount&order=desc&limit=10&offset=0",
			mockSetup: func(m *MockExpensesAnalytics) {
				m.On("GetExpensesSummary", mock.Anything, usecases.ExpensesSummaryParams{
					UserID: "user123",
					Period: usecases.DatePeriod{
						Year:  intPtr(2024),
						Month: intPtr(3),
					},
					Sorting: usecases.SortingCriteria{
						Field: usecases.SortByAmount,
						Order: usecases.Descending,
					},
					Pagination: usecases.PaginationParams{
						Limit:  10,
						Offset: 0,
					},
				}).Return(&usecases.ExpensesSummary{
					Expenses: []usecases.ExpenseItem{
						{
							ID:          "exp1",
							Description: "Gasto 1",
							Amount:      1000.0,
							CategoryID:  "cat1",
							Paid:        true,
						},
						{
							ID:          "exp2",
							Description: "Gasto 2",
							Amount:      500.0,
							CategoryID:  "cat2",
							Paid:        false,
						},
					},
					Summary: usecases.ExpenseSummary{
						TotalAmount:   1500.0,
						PaidAmount:    1000.0,
						PendingAmount: 500.0,
					},
					Pagination: usecases.PaginationInfo{
						Total:   2,
						Limit:   10,
						Offset:  0,
						HasMore: false,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - sin X-Caller-ID",
			callerID:       "",
			queryParams:    "",
			mockSetup:      func(m *MockExpensesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "X-Caller-ID header es requerido",
		},
		{
			name:           "error - mes inválido",
			callerID:       "user123",
			queryParams:    "month=13",
			mockSetup:      func(m *MockExpensesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Mes debe estar entre 1 y 12",
		},
		{
			name:           "error - año inválido",
			callerID:       "user123",
			queryParams:    "year=invalid",
			mockSetup:      func(m *MockExpensesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Año inválido",
		},
		{
			name:           "error - orden inválido",
			callerID:       "user123",
			queryParams:    "order=invalid",
			mockSetup:      func(m *MockExpensesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Orden inválido (debe ser asc o desc)",
		},
		{
			name:           "error - campo de ordenamiento inválido",
			callerID:       "user123",
			queryParams:    "sort_by=invalid",
			mockSetup:      func(m *MockExpensesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Campo de ordenamiento inválido",
		},
		{
			name:           "error - límite inválido",
			callerID:       "user123",
			queryParams:    "limit=1001",
			mockSetup:      func(m *MockExpensesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Límite inválido (debe ser entre 1 y 1000)",
		},
		{
			name:           "error - offset negativo",
			callerID:       "user123",
			queryParams:    "offset=-1",
			mockSetup:      func(m *MockExpensesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Offset inválido (debe ser >= 0)",
		},
		{
			name:        "error interno",
			callerID:    "user123",
			queryParams: "",
			mockSetup: func(m *MockExpensesAnalytics) {
				m.On("GetExpensesSummary", mock.Anything, mock.Anything).Return(nil, errors.NewInternalServerError("Error interno del servidor"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			mockService := new(MockExpensesAnalytics)
			handler := NewExpensesHandler(mockService)
			tt.mockSetup(mockService)

			// Configurar request
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/expenses/summary?"+tt.queryParams, nil)

			// Simular header X-Caller-ID
			if tt.callerID != "" {
				c.Request.Header.Set("X-Caller-ID", tt.callerID)
			}

			// Ejecutar handler
			handler.GetExpensesSummary(c)

			// Verificar respuesta
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}

			// Verificar que se llamaron los mocks según lo esperado
			mockService.AssertExpectations(t)
		})
	}
}

// MockExpensesAnalytics es un mock del servicio de analytics de gastos
type MockExpensesAnalytics struct {
	mock.Mock
}

func (m *MockExpensesAnalytics) GetExpensesSummary(ctx context.Context, params usecases.ExpensesSummaryParams) (*usecases.ExpensesSummary, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.ExpensesSummary), args.Error(1)
}
