package update

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	expenses "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	expensesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService es un mock del servicio de actualización de gastos
type MockService struct {
	mock.Mock
}

func (m *MockService) UpdateExpense(ctx context.Context, userID, id string, request *expensesDomain.UpdateExpenseRequest) (*expensesDomain.UpdateExpenseResponse, error) {
	args := m.Called(ctx, userID, id, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expensesDomain.UpdateExpenseResponse), args.Error(1)
}

func TestUpdateExpenseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		userID         string
		expenseID      string
		mockResponse   *expensesDomain.UpdateExpenseResponse
		mockError      error
		expectedStatus int
	}{
		{
			name:      "Actualización exitosa de gasto",
			userID:    "user123",
			expenseID: "expense123",
			requestBody: expensesDomain.UpdateExpenseRequest{
				Amount:      150.75,
				Description: "Updated Description",
				DueDate:     "2024-05-15",
				CategoryID:  "Updated Category",
				Paid:        false,
			},
			mockResponse: &expensesDomain.UpdateExpenseResponse{
				ID:          "expense123",
				UserID:      "user123",
				Amount:      150.75,
				Description: "Updated Description",
				DueDate:     "2024-05-15",
				CategoryID:  "Updated Category",
				Paid:        false,
				CreatedAt:   "2024-04-01T00:00:00Z",
				UpdatedAt:   "2024-04-01T00:00:00Z",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Error - gasto no encontrado",
			userID:    "user123",
			expenseID: "expense123",
			requestBody: expensesDomain.UpdateExpenseRequest{
				Amount:      150.75,
				Description: "Updated Description",
				DueDate:     "2024-05-15",
				CategoryID:  "Updated Category",
				Paid:        false,
			},
			mockError:      errors.NewResourceNotFound("Gasto no encontrado"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "Error - ID de usuario vacío",
			userID:    "",
			expenseID: "expense123",
			requestBody: expenses.UpdateExpenseRequest{
				Amount:      150.75,
				Description: "Updated Description",
				DueDate:     "2024-05-15",
				CategoryID:  "Updated Category",
				Paid:        false,
			},
			mockError:      errors.NewBadRequest("El ID del usuario es requerido"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el mock
			mockService := new(MockService)
			handler := NewHandler(mockService, nil)

			// Configurar la solicitud
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/expenses/"+tt.expenseID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Configurar el contexto de Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("user_id", tt.userID)
			c.Params = gin.Params{{Key: "id", Value: tt.expenseID}}

			// Configurar las expectativas del mock
			if tt.mockResponse != nil {
				mockService.On("UpdateExpense", c.Request.Context(), tt.userID, tt.expenseID, mock.AnythingOfType("*expenses.UpdateExpenseRequest")).
					Return(tt.mockResponse, nil)
			} else {
				mockService.On("UpdateExpense", c.Request.Context(), tt.userID, tt.expenseID, mock.AnythingOfType("*expenses.UpdateExpenseRequest")).
					Return(nil, tt.mockError)
			}

			// Ejecutar el handler
			handler.UpdateExpense(c)

			// Verificar el código de estado
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar que se llamaron todas las expectativas del mock
			mockService.AssertExpectations(t)
		})
	}
}
