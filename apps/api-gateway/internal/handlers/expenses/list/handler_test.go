package list

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	expensesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService es un mock del servicio de listado de gastos
type MockService struct {
	mock.Mock
}

func (m *MockService) ListExpenses(ctx context.Context, userID string) (*expensesDomain.ListExpensesResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expensesDomain.ListExpensesResponse), args.Error(1)
}

func (m *MockService) ListUnpaidExpenses(ctx context.Context, userID string) (*expensesDomain.ListExpensesResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expensesDomain.ListExpensesResponse), args.Error(1)
}

func TestListExpensesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockResponse   *expensesDomain.ListExpensesResponse
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Listado exitoso de gastos",
			userID: "user123",
			mockResponse: &expensesDomain.ListExpensesResponse{
				Expenses: []expensesDomain.GetExpenseResponse{
					{
						UserID:        "user123",
						Amount:        100.50,
						AmountPaid:    0,
						PendingAmount: 100.50,
						Description:   "Descripción 1",
						DueDate:       "2024-04-30",
						CategoryID:    "Categoría 1",
						Paid:          false,
					},
					{
						UserID:        "user123",
						Amount:        200.75,
						AmountPaid:    200.75,
						PendingAmount: 0,
						Description:   "Descripción 2",
						DueDate:       "2024-05-15",
						CategoryID:    "Categoría 2",
						Paid:          true,
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"expenses":[{"user_id":"user123","amount":100.5,"amount_paid":0,"pending_amount":100.5,"description":"Descripción 1","due_date":"2024-04-30","category_id":"Categoría 1","paid":false,"id":"","created_at":"","updated_at":""},{"user_id":"user123","amount":200.75,"amount_paid":200.75,"pending_amount":0,"description":"Descripción 2","due_date":"2024-05-15","category_id":"Categoría 2","paid":true,"id":"","created_at":"","updated_at":""}]}`,
		},
		{
			name:           "Error - ID de usuario vacío",
			userID:         "",
			mockError:      errors.New("user_id is required"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Usuario no autenticado","message":""}`,
		},
		{
			name:           "Error - Usuario no encontrado",
			userID:         "user999",
			mockError:      errors.New("user_id is required"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Internal server error","message":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el mock
			mockService := new(MockService)
			handler := NewHandler(mockService, nil)

			// Configurar la solicitud
			req, err := http.NewRequest(http.MethodGet, "/expenses/"+tt.userID, nil)
			require.NoError(t, err)
			req.Header.Set("X-Caller-ID", tt.userID)

			// Configurar el contexto de Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{gin.Param{Key: "user_id", Value: tt.userID}}

			// Simular el comportamiento del middleware de autenticación
			if tt.userID != "" {
				c.Set("user_id", tt.userID)
			}

			// Configurar las expectativas del mock
			if tt.mockResponse != nil {
				mockService.On("ListExpenses", mock.Anything, tt.userID).
					Return(tt.mockResponse, nil)
			} else if tt.userID != "" {
				mockService.On("ListExpenses", mock.Anything, tt.userID).
					Return(nil, tt.mockError)
			}

			// Ejecutar el handler
			handler.ListExpenses(c)

			// Verificar el código de estado
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar el tipo de contenido
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			// Verificar el cuerpo de la respuesta
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// Verificar que se llamaron todas las expectativas del mock
			mockService.AssertExpectations(t)
		})
	}
}

func TestListUnpaidExpensesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockResponse   *expensesDomain.ListExpensesResponse
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Listado exitoso de gastos no pagados",
			userID: "user123",
			mockResponse: &expensesDomain.ListExpensesResponse{
				Expenses: []expensesDomain.GetExpenseResponse{
					{
						UserID:        "user123",
						Amount:        100.50,
						AmountPaid:    0,
						PendingAmount: 100.50,
						Description:   "Descripción 1",
						DueDate:       "2024-04-30",
						CategoryID:    "Categoría 1",
						Paid:          false,
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"expenses":[{"user_id":"user123","amount":100.5,"amount_paid":0,"pending_amount":100.5,"description":"Descripción 1","due_date":"2024-04-30","category_id":"Categoría 1","paid":false,"id":"","created_at":"","updated_at":""}]}`,
		},
		{
			name:           "Error - ID de usuario vacío",
			userID:         "",
			mockError:      errors.New("user_id is required"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Usuario no autenticado","message":""}`,
		},
		{
			name:           "Error - Usuario no encontrado",
			userID:         "user999",
			mockError:      errors.New("user_id is required"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Internal server error","message":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configurar el mock
			mockService := new(MockService)
			handler := NewHandler(mockService, nil)

			// Configurar la solicitud
			req, err := http.NewRequest(http.MethodGet, "/expenses/unpaid", nil)
			require.NoError(t, err)
			req.Header.Set("X-Caller-ID", tt.userID)

			// Configurar el contexto de Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Simular el comportamiento del middleware de autenticación
			if tt.userID != "" {
				c.Set("user_id", tt.userID)
			}

			// Configurar las expectativas del mock
			if tt.mockResponse != nil {
				mockService.On("ListUnpaidExpenses", mock.Anything, tt.userID).
					Return(tt.mockResponse, nil)
			} else if tt.userID != "" {
				mockService.On("ListUnpaidExpenses", mock.Anything, tt.userID).
					Return(nil, tt.mockError)
			}

			// Ejecutar el handler
			handler.ListUnpaidExpenses(c)

			// Verificar el código de estado
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar el tipo de contenido
			assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

			// Verificar el cuerpo de la respuesta
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			// Verificar que se llamaron todas las expectativas del mock
			mockService.AssertExpectations(t)
		})
	}
}
