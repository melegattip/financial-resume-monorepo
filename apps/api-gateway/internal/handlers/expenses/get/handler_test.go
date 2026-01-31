package get

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	expensesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService es un mock del servicio de obtención de gastos
type MockService struct {
	mock.Mock
}

func (m *MockService) GetExpense(ctx context.Context, userID, id string) (*expensesDomain.GetExpenseResponse, error) {
	args := m.Called(ctx, userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expensesDomain.GetExpenseResponse), args.Error(1)
}

func TestGetExpenseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		expenseID      string
		mockResponse   *expensesDomain.GetExpenseResponse
		mockError      error
		expectedStatus int
	}{
		{
			name:      "Obtención exitosa de gasto",
			userID:    "user123",
			expenseID: "expense123",
			mockResponse: &expensesDomain.GetExpenseResponse{
				UserID:      "user123",
				Amount:      100.50,
				Description: "Test Description",
				DueDate:     "2024-04-30",
				CategoryID:  "Test Category",
				Paid:        false,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - gasto no encontrado",
			userID:         "user123",
			expenseID:      "expense123",
			mockError:      errors.NewResourceNotFound("Gasto no encontrado"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Error - ID de usuario vacío",
			userID:         "",
			expenseID:      "expense123",
			mockError:      errors.NewBadRequest("El ID del usuario es requerido"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error - ID de gasto vacío",
			userID:         "user123",
			expenseID:      "",
			mockError:      errors.NewBadRequest("El ID del gasto es requerido"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el mock
			mockService := new(MockService)
			handler := NewHandler(mockService, nil)

			// Configurar la solicitud
			req, _ := http.NewRequest(http.MethodGet, "/expenses/"+tt.expenseID, nil)

			// Configurar el contexto de Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("user_id", tt.userID)
			c.Params = gin.Params{{Key: "id", Value: tt.expenseID}}

			// Configurar las expectativas del mock
			if tt.mockResponse != nil {
				mockService.On("GetExpense", c.Request.Context(), tt.userID, tt.expenseID).
					Return(tt.mockResponse, nil)
			} else {
				mockService.On("GetExpense", c.Request.Context(), tt.userID, tt.expenseID).
					Return(nil, tt.mockError)
			}

			// Ejecutar el handler
			handler.GetExpense(c)

			// Verificar el código de estado
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar que se llamaron todas las expectativas del mock
			mockService.AssertExpectations(t)
		})
	}
}
