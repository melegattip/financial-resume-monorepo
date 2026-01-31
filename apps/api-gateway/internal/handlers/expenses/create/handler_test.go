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
	expensesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService es un mock del servicio de creación de gastos
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateExpense(ctx context.Context, request *expensesDomain.CreateExpenseRequest) (*expensesDomain.CreateExpenseResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expensesDomain.CreateExpenseResponse), args.Error(1)
}

func TestCreateExpenseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		userID         string
		mockResponse   *expensesDomain.CreateExpenseResponse
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Creación exitosa de gasto",
			requestBody: expensesDomain.CreateExpenseRequest{
				Amount:      100.50,
				Description: "Test Description",
				DueDate:     "2024-04-30",
				CategoryID:  "Test Category",
			},
			userID: "user123",
			mockResponse: &expensesDomain.CreateExpenseResponse{
				UserID:      "user123",
				Amount:      100.50,
				Description: "Test Description",
				DueDate:     "2024-04-30",
				CategoryID:  "Test Category",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   nil,
		},
		{
			name: "Error de validación - nombre vacío",
			requestBody: expensesDomain.CreateExpenseRequest{
				Amount:      100.50,
				Description: "Test Description",
				DueDate:     "2024-04-30",
				CategoryID:  "Test Category",
			},
			userID:         "user123",
			mockError:      errors.NewBadRequest("Invalid request body"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Error de validación - monto inválido",
			requestBody: expensesDomain.CreateExpenseRequest{
				UserID:      "user123",
				Amount:      0,
				Description: "Test Description",
				DueDate:     "2024-04-30",
				CategoryID:  "Test Category",
			},
			userID:         "user123",
			mockError:      errors.NewBadRequest("Invalid request body"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el mock
			mockService := &MockService{}

			// Configurar las expectativas del mock solo para el caso exitoso
			if tt.mockResponse != nil {
				mockService.On("CreateExpense", mock.Anything, mock.MatchedBy(func(req *expensesDomain.CreateExpenseRequest) bool {
					return true
				})).Return(tt.mockResponse, nil)
			} else if tt.mockError != nil {
				mockService.On("CreateExpense", mock.Anything, mock.MatchedBy(func(req *expensesDomain.CreateExpenseRequest) bool {
					return true
				})).Return(nil, tt.mockError)
			}

			handler := NewHandler(mockService, nil)

			// Configurar la solicitud
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Configurar el contexto de Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("user_id", tt.userID)

			// Ejecutar el handler
			handler.CreateExpense(c)

			// Verificar el código de estado
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar el cuerpo de la respuesta
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			} else if tt.mockResponse != nil {
				var response struct {
					Data expensesDomain.CreateExpenseResponse `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse, &response.Data)
			}

			// Verificar que se llamaron todas las expectativas del mock
			mockService.AssertExpectations(t)
		})
	}
}
