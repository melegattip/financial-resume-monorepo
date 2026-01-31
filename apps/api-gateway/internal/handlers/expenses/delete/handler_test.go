package delete

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService es un mock del servicio de eliminación de gastos
type MockService struct {
	mock.Mock
}

func (m *MockService) DeleteExpense(ctx context.Context, userID, id string) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}

func TestDeleteExpenseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		expenseID      string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Eliminación exitosa de gasto",
			userID:         "user123",
			expenseID:      "expense123",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
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
			req, _ := http.NewRequest(http.MethodDelete, "/expenses/"+tt.expenseID, nil)

			// Configurar el contexto de Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("user_id", tt.userID)
			c.Params = gin.Params{{Key: "id", Value: tt.expenseID}}

			// Configurar las expectativas del mock
			mockService.On("DeleteExpense", c.Request.Context(), tt.userID, tt.expenseID).
				Return(tt.mockError)

			// Ejecutar el handler
			handler.DeleteExpense(c)

			// Verificar el código de estado
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar que se llamaron todas las expectativas del mock
			mockService.AssertExpectations(t)
		})
	}
}
