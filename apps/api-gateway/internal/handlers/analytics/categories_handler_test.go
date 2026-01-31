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

func TestCategoriesHandler_GetCategoriesAnalytics(t *testing.T) {
	tests := []struct {
		name           string
		callerID       string
		queryParams    string
		mockSetup      func(*MockCategoriesAnalytics)
		expectedStatus int
		expectedError  string
	}{
		{
			name:        "obtener analytics exitosamente",
			callerID:    "user123",
			queryParams: "year=2024&month=3",
			mockSetup: func(m *MockCategoriesAnalytics) {
				m.On("GetCategoriesAnalytics", mock.Anything, usecases.CategoriesAnalyticsParams{
					UserID: "user123",
					Period: usecases.DatePeriod{
						Year:  intPtr(2024),
						Month: intPtr(3),
					},
				}).Return(&usecases.CategoriesAnalytics{
					Categories: []usecases.CategoryItem{
						{
							CategoryID:            "cat1",
							CategoryName:          "Comida",
							TotalAmount:           1000.0,
							PercentageOfExpenses:  66.67,
							PercentageOfIncome:    50.0,
							TransactionCount:      2,
							AveragePerTransaction: 500.0,
							ColorSeed:             1,
						},
						{
							CategoryID:            "cat2",
							CategoryName:          "Transporte",
							TotalAmount:           500.0,
							PercentageOfExpenses:  33.33,
							PercentageOfIncome:    25.0,
							TransactionCount:      1,
							AveragePerTransaction: 500.0,
							ColorSeed:             2,
						},
					},
					Summary: usecases.CategorySummary{
						TotalCategories:  2,
						LargestCategory:  "Comida",
						SmallestCategory: "Transporte",
						TotalAmount:      1500.0,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error - sin X-Caller-ID",
			callerID:       "",
			queryParams:    "",
			mockSetup:      func(m *MockCategoriesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "X-Caller-ID header es requerido",
		},
		{
			name:           "error - mes inválido",
			callerID:       "user123",
			queryParams:    "month=13",
			mockSetup:      func(m *MockCategoriesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Mes debe estar entre 1 y 12",
		},
		{
			name:           "error - año inválido",
			callerID:       "user123",
			queryParams:    "year=invalid",
			mockSetup:      func(m *MockCategoriesAnalytics) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Año inválido",
		},
		{
			name:        "error interno",
			callerID:    "user123",
			queryParams: "",
			mockSetup: func(m *MockCategoriesAnalytics) {
				m.On("GetCategoriesAnalytics", mock.Anything, mock.Anything).Return(nil, errors.NewInternalServerError("Error interno del servidor"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			mockService := new(MockCategoriesAnalytics)
			handler := NewCategoriesHandler(mockService)
			tt.mockSetup(mockService)

			// Configurar request
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/categories/analytics?"+tt.queryParams, nil)

			// Simular header X-Caller-ID
			if tt.callerID != "" {
				c.Request.Header.Set("X-Caller-ID", tt.callerID)
			}

			// Ejecutar handler
			handler.GetCategoriesAnalytics(c)

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

// MockCategoriesAnalytics es un mock del servicio de analytics de categorías
type MockCategoriesAnalytics struct {
	mock.Mock
}

func (m *MockCategoriesAnalytics) GetCategoriesAnalytics(ctx context.Context, params usecases.CategoriesAnalyticsParams) (*usecases.CategoriesAnalytics, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecases.CategoriesAnalytics), args.Error(1)
}
