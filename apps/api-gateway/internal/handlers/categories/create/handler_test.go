package create

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategoryHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    map[string]interface{}
		mockSetup      func(*MockCategoryService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "crear categoría exitosamente",
			userID: "user_id",
			requestBody: map[string]interface{}{
				"name": "Test Category",
			},
			mockSetup: func(m *MockCategoryService) {
				m.On("Create", categories.CreateCategoryRequest{
					Name:   "Test Category",
					UserID: "user_id",
				}).Return(&domain.Category{
					ID:     "cat_123",
					Name:   "Test Category",
					UserID: "user_id",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name:   "error al crear categoría con nombre vacío",
			userID: "user_id",
			requestBody: map[string]interface{}{
				"name": "",
			},
			mockSetup:      func(m *MockCategoryService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:   "error al crear categoría que ya existe",
			userID: "user_id",
			requestBody: map[string]interface{}{
				"name": "Existing Category",
			},
			mockSetup: func(m *MockCategoryService) {
				m.On("Create", categories.CreateCategoryRequest{
					Name:   "Existing Category",
					UserID: "user_id",
				}).Return(nil, errors.NewResourceAlreadyExists("la categoría ya existe"))
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "la categoría ya existe",
		},
		{
			name:   "error sin user_id en contexto",
			userID: "",
			requestBody: map[string]interface{}{
				"name": "Test Category",
			},
			mockSetup:      func(m *MockCategoryService) {},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Usuario no autenticado",
		},
		{
			name:   "error con request body inválido",
			userID: "user_id",
			requestBody: map[string]interface{}{
				"invalid": "field",
			},
			mockSetup:      func(m *MockCategoryService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCategoryService)
			handler := NewHandler(mockService)
			tt.mockSetup(mockService)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

			// Simular user_id en contexto (middleware de auth lo setea en runtime)
			if tt.userID != "" {
				c.Set("user_id", tt.userID)
			}

			handler.Handle(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}

			mockService.AssertExpectations(t)
		})
	}
}
