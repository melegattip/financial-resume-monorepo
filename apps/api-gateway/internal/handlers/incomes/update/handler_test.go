package update

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateIncomeHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		incomeID       string
		requestBody    map[string]interface{}
		mockSetup      func(*MockIncomeService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "actualizar ingreso exitosamente",
			userID:   "user123",
			incomeID: "inc123",
			requestBody: map[string]interface{}{
				"amount":      1500.0,
				"description": "Salario actualizado",
				"category_id": "cat123",
			},
			mockSetup: func(m *MockIncomeService) {
				m.On("UpdateIncome", mock.Anything, "user123", "inc123", mock.MatchedBy(func(req *incomes.UpdateIncomeRequest) bool {
					return req.Amount == 1500.0 &&
						req.Description == "Salario actualizado" &&
						req.CategoryID == "cat123"
				})).Return(&incomes.UpdateIncomeResponse{
					ID:          "inc123",
					UserID:      "user123",
					Amount:      1500.0,
					Description: "Salario actualizado",
					CategoryID:  "cat123",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "error - ingreso no encontrado",
			userID:   "user123",
			incomeID: "nonexistent",
			requestBody: map[string]interface{}{
				"amount": 1500.0,
			},
			mockSetup: func(m *MockIncomeService) {
				m.On("UpdateIncome", mock.Anything, "user123", "nonexistent", mock.Anything).Return(nil, coreErrors.NewResourceNotFound("Ingreso no encontrado"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Ingreso no encontrado",
		},
		{
			name:     "error - monto negativo",
			userID:   "user123",
			incomeID: "inc123",
			requestBody: map[string]interface{}{
				"amount": -100.0,
			},
			mockSetup: func(m *MockIncomeService) {
				m.On("UpdateIncome", mock.Anything, "user123", "inc123", mock.MatchedBy(func(req *incomes.UpdateIncomeRequest) bool {
					return req.Amount == -100.0
				})).Return(nil, coreErrors.NewBadRequest("El monto debe ser positivo"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "El monto debe ser positivo",
		},
		{
			name:     "error - categoría no existe",
			userID:   "user123",
			incomeID: "inc123",
			requestBody: map[string]interface{}{
				"category_id": "nonexistent",
			},
			mockSetup: func(m *MockIncomeService) {
				m.On("UpdateIncome", mock.Anything, "user123", "inc123", mock.MatchedBy(func(req *incomes.UpdateIncomeRequest) bool {
					return req.CategoryID == "nonexistent"
				})).Return(nil, coreErrors.NewResourceNotFound("La categoría no existe"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "La categoría no existe",
		},
		{
			name:           "error - sin user_id",
			userID:         "",
			incomeID:       "inc123",
			requestBody:    map[string]interface{}{},
			mockSetup:      func(m *MockIncomeService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "User ID is required",
		},
		{
			name:           "error - sin income_id",
			userID:         "user123",
			incomeID:       "",
			requestBody:    map[string]interface{}{},
			mockSetup:      func(m *MockIncomeService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Income ID is required",
		},
		{
			name:     "error - request body inválido",
			userID:   "user123",
			incomeID: "inc123",
			requestBody: map[string]interface{}{
				"amount": "invalid",
			},
			mockSetup:      func(m *MockIncomeService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			mockService := new(MockIncomeService)
			handler := NewHandler(mockService, nil)
			tt.mockSetup(mockService)

			// Configurar request
			body, _ := json.Marshal(tt.requestBody)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Simular middleware de autenticación
			if tt.userID != "" {
				c.Set("user_id", tt.userID)
			}

			// Simular parámetros de ruta
			if tt.incomeID != "" {
				c.Params = []gin.Param{{Key: "id", Value: tt.incomeID}}
			}

			// Ejecutar handler
			handler.UpdateIncome(c)

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

// MockIncomeService es un mock del servicio de ingresos
type MockIncomeService struct {
	mock.Mock
}

func (m *MockIncomeService) CreateIncome(ctx context.Context, request *incomes.CreateIncomeRequest) (*incomes.CreateIncomeResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.CreateIncomeResponse), args.Error(1)
}

func (m *MockIncomeService) GetIncome(ctx context.Context, userID, id string) (*incomes.GetIncomeResponse, error) {
	args := m.Called(ctx, userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.GetIncomeResponse), args.Error(1)
}

func (m *MockIncomeService) ListIncomes(ctx context.Context, userID string) (*incomes.ListIncomesResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.ListIncomesResponse), args.Error(1)
}

func (m *MockIncomeService) UpdateIncome(ctx context.Context, userID string, incomeID string, request *incomes.UpdateIncomeRequest) (*incomes.UpdateIncomeResponse, error) {
	args := m.Called(ctx, userID, incomeID, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incomes.UpdateIncomeResponse), args.Error(1)
}

func (m *MockIncomeService) DeleteIncome(ctx context.Context, userID, id string) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}
