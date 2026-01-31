package create

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateIncomeHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    map[string]interface{}
		mockSetup      func(*MockIncomeService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "crear ingreso exitosamente",
			userID: "user123",
			requestBody: map[string]interface{}{
				"amount":      1000.0,
				"description": "Salario",
				"category_id": "cat123",
			},
			mockSetup: func(m *MockIncomeService) {
				m.On("CreateIncome", mock.Anything, mock.MatchedBy(func(req *incomes.CreateIncomeRequest) bool {
					return req.UserID == "user123" &&
						req.Amount == 1000.0 &&
						req.Description == "Salario" &&
						req.CategoryID == "cat123"
				})).Return(&incomes.CreateIncomeResponse{
					ID:          "inc123",
					UserID:      "user123",
					Amount:      1000.0,
					Description: "Salario",
					CategoryID:  "cat123",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "error - monto negativo",
			userID: "user123",
			requestBody: map[string]interface{}{
				"amount":      -100.0,
				"description": "Salario",
				"category_id": "cat123",
			},
			mockSetup: func(m *MockIncomeService) {
				m.On("CreateIncome", mock.Anything, mock.MatchedBy(func(req *incomes.CreateIncomeRequest) bool {
					return req.Amount == -100.0
				})).Return(nil, errors.NewBadRequest("El monto debe ser positivo"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "El monto debe ser positivo",
		},
		{
			name:   "error - sin descripción",
			userID: "user123",
			requestBody: map[string]interface{}{
				"amount":      1000.0,
				"category_id": "cat123",
			},
			mockSetup: func(m *MockIncomeService) {
				m.On("CreateIncome", mock.Anything, mock.MatchedBy(func(req *incomes.CreateIncomeRequest) bool {
					return req.Description == ""
				})).Return(nil, errors.NewBadRequest("La descripción es requerida"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "La descripción es requerida",
		},
		{
			name:   "error - categoría no existe",
			userID: "user123",
			requestBody: map[string]interface{}{
				"amount":      1000.0,
				"description": "Salario",
				"category_id": "nonexistent",
			},
			mockSetup: func(m *MockIncomeService) {
				m.On("CreateIncome", mock.Anything, mock.MatchedBy(func(req *incomes.CreateIncomeRequest) bool {
					return req.CategoryID == "nonexistent"
				})).Return(nil, errors.NewResourceNotFound("La categoría no existe"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "La categoría no existe",
		},
		{
			name:   "error - request body inválido",
			userID: "user123",
			requestBody: map[string]interface{}{
				"amount": "invalid",
			},
			mockSetup:      func(m *MockIncomeService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:           "error - sin user_id",
			userID:         "",
			requestBody:    map[string]interface{}{},
			mockSetup:      func(m *MockIncomeService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "User ID is required",
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
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Simular middleware de autenticación
			if tt.userID != "" {
				c.Set("user_id", tt.userID)
			}

			// Ejecutar handler
			handler.CreateIncome(c)

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

// Verificar que MockIncomeService implementa la interfaz
var _ create.ServiceInterface = (*MockIncomeService)(nil)
