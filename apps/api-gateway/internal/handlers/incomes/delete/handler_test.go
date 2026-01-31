package delete

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteIncomeHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		incomeID       string
		mockSetup      func(*MockIncomeService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "eliminar ingreso exitosamente",
			userID:   "user123",
			incomeID: "inc123",
			mockSetup: func(m *MockIncomeService) {
				m.On("DeleteIncome", mock.Anything, "user123", "inc123").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "error - ingreso no encontrado",
			userID:   "user123",
			incomeID: "nonexistent",
			mockSetup: func(m *MockIncomeService) {
				m.On("DeleteIncome", mock.Anything, "user123", "nonexistent").Return(errors.NewResourceNotFound("Ingreso no encontrado"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Ingreso no encontrado",
		},
		{
			name:           "error - sin user_id",
			userID:         "",
			incomeID:       "inc123",
			mockSetup:      func(m *MockIncomeService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "User ID is required",
		},
		{
			name:           "error - sin income_id",
			userID:         "user123",
			incomeID:       "",
			mockSetup:      func(m *MockIncomeService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Income ID is required",
		},
		{
			name:     "error - usuario no autorizado",
			userID:   "user123",
			incomeID: "inc123",
			mockSetup: func(m *MockIncomeService) {
				m.On("DeleteIncome", mock.Anything, "user123", "inc123").Return(errors.NewUnauthorizedRequest("Usuario no autorizado"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Usuario no autorizado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar mock
			mockService := new(MockIncomeService)
			handler := NewHandler(mockService, nil)
			tt.mockSetup(mockService)

			// Configurar request
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodDelete, "/", nil)

			// Simular middleware de autenticación
			if tt.userID != "" {
				c.Set("user_id", tt.userID)
			}

			// Simular parámetros de ruta
			if tt.incomeID != "" {
				c.Params = []gin.Param{{Key: "id", Value: tt.incomeID}}
			}

			// Ejecutar handler
			handler.Handle(c)

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
