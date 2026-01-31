package list

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestListCategoriesHandler(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		mockSetup      func(*MockCategoryService)
		expectedStatus int
		expectedError  string
		expectedData   []*domain.Category
	}{
		{
			testName: "listar categorías exitosamente",
			userID:   "user_id",
			mockSetup: func(m *MockCategoryService) {
				m.On("List", "user_id").Return([]*domain.Category{
					{
						ID:     "cat_123",
						Name:   "Test Category 1",
						UserID: "user_id",
					},
					{
						ID:     "cat_456",
						Name:   "Test Category 2",
						UserID: "user_id",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			expectedData: []*domain.Category{
				{
					ID:     "cat_123",
					Name:   "Test Category 1",
					UserID: "user_id",
				},
				{
					ID:     "cat_456",
					Name:   "Test Category 2",
					UserID: "user_id",
				},
			},
		},
		{
			testName: "error al listar categorías",
			userID:   "user_id",
			mockSetup: func(m *MockCategoryService) {
				m.On("List", "user_id").Return(nil, errors.New("error al listar categorías"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Internal server error",
			expectedData:   nil,
		},
		{
			testName: "user_id no proporcionado",
			userID:   "",
			mockSetup: func(m *MockCategoryService) {
				// No se espera ninguna llamada al servicio
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Usuario no autenticado",
			expectedData:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(MockCategoryService)
			handler := NewHandler(mockService)
			tt.mockSetup(mockService)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Siempre inicializamos la petición
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.userID != "" {
				c.Request.Header.Set("X-Caller-ID", tt.userID)
				// Simular el comportamiento del middleware de autenticación
				c.Set("user_id", tt.userID)
			}

			handler.Handle(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, response["error"])
			} else {
				assert.Equal(t, "Categories retrieved successfully", response["message"])
				data, ok := response["data"].([]interface{})
				assert.True(t, ok)
				assert.Equal(t, len(tt.expectedData), len(data))
			}

			mockService.AssertExpectations(t)
		})
	}
}
