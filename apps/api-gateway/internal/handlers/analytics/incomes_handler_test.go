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

func TestIncomesHandler_GetIncomesSummary(t *testing.T) {
	tests := []struct {
		name           string
		callerID       string
		queryParams    string
		mockSetup      func(*MockIncomesAnalytics)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "obtener resumen exitosamente",
			callerID:    "user123",
			queryParams: "year=2024&month=3&sort_by=amount&order=desc",
			mockSetup: func(m *MockIncomesAnalytics) {
				m.On("GetIncomesSummary", mock.Anything, usecases.IncomesSummaryParams{
					UserID: "user123",
					Period: usecases.DatePeriod{
						Year:  intPtr(2024),
						Month: intPtr(3),
					},
					Sorting: usecases.SortingCriteria{
						Field: usecases.SortByAmount,
						Order: usecases.Descending,
					},
				}).Return(&usecases.IncomesSummary{
					Incomes: []usecases.IncomeItem{
						{
							ID:           "inc1",
							Description:  "Salario",
							Amount:       5000.0,
							CategoryID:   "cat1",
							CategoryName: "Trabajo",
							CreatedAt:    "2024-03-01",
						},
						{
							ID:           "inc2",
							Description:  "Bonus",
							Amount:       1000.0,
							CategoryID:   "cat2",
							CategoryName: "Extras",
							CreatedAt:    "2024-03-15",
						},
					},
					Summary: usecases.IncomeSummary{
						TotalAmount:        6000.0,
						AverageTransaction: 3000.0,
						TransactionCount:   2,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - sin X-Caller-ID",
			callerID:       "",
			queryParams:    "",
			mockSetup:      func(m *MockIncomesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "X-Caller-ID header es requerido",
		},
		{
			name:           "error - mes inválido",
			callerID:       "user123",
			queryParams:    "month=13",
			mockSetup:      func(m *MockIncomesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Mes debe estar entre 1 y 12",
		},
		{
			name:           "error - año inválido",
			callerID:       "user123",
			queryParams:    "year=invalid",
			mockSetup:      func(m *MockIncomesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Año inválido",
		},
		{
			name:           "error - orden inválido",
			callerID:       "user123",
			queryParams:    "order=invalid",
			mockSetup:      func(m *MockIncomesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Orden inválido (debe ser asc o desc)",
		},
		{
			name:           "error - campo de ordenamiento inválido",
			callerID:       "user123",
			queryParams:    "sort_by=invalid",
			mockSetup:      func(m *MockIncomesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Campo de ordenamiento inválido",
		},
		{
			name:        "error interno",
			callerID:    "user123",
			queryParams: "",
			mockSetup: func(m *MockIncomesAnalytics) {
				m.On("GetIncomesSummary", mock.Anything, mock.Anything).Return(nil, errors.NewInternalServerError("Error interno del servidor"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			mockService := new(MockIncomesAnalytics)
			handler := NewIncomesHandler(mockService)
			tt.mockSetup(mockService)

			// Configurar request
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/incomes/summary?"+tt.queryParams, nil)

			// Simular header X-Caller-ID
			if tt.callerID != "" {
				c.Request.Header.Set("X-Caller-ID", tt.callerID)
			}

			// Ejecutar handler
			handler.GetIncomesSummary(c)

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

// MockIncomesAnalytics es un mock del servicio de analytics de ingresos
type MockIncomesAnalytics struct {
	mock.Mock
}

func (m *MockIncomesAnalytics) GetIncomesSummary(ctx context.Context, params usecases.IncomesSummaryParams) (*usecases.IncomesSummary, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.IncomesSummary), args.Error(1)
}
