package categories

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/update"
	"github.com/stretchr/testify/assert"
)

func TestCategoryHandler(t *testing.T) {
	tests := []struct {
		testName       string
		method         string
		path           string
		body           string
		userID         string
		mockSetup      func(*MockCategoryService)
		expectedStatus int
		expectedError  string
	}{
		{
			testName: "crear categoría exitosamente",
			method:   "POST",
			path:     "/categories",
			body:     `{"name": "Test Category", "user_id": "user_id"}`,
			userID:   "user_id",
			mockSetup: func(m *MockCategoryService) {
				m.On("Create", categories.CreateCategoryRequest{
					Name:   "Test Category",
					UserID: "user_id",
				}).Return(&domain.Category{
					ID:     "cat_12345678",
					Name:   "Test Category",
					UserID: "user_id",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			testName: "error al crear categoría que ya existe",
			method:   "POST",
			path:     "/categories",
			body:     `{"name": "Existing Category", "user_id": "user_id"}`,
			userID:   "user_id",
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
			testName: "listar categorías exitosamente",
			method:   "GET",
			path:     "/categories",
			userID:   "user_id",
			mockSetup: func(m *MockCategoryService) {
				m.On("List", "user_id").Return([]*domain.Category{
					{
						ID:     "cat_12345678",
						Name:   "Test Category 1",
						UserID: "user_id",
					},
					{
						ID:     "cat_87654321",
						Name:   "Test Category 2",
						UserID: "user_id",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			testName: "obtener categoría exitosamente",
			method:   "GET",
			path:     "/categories/Test_Category",
			userID:   "user_id",
			mockSetup: func(m *MockCategoryService) {
				m.On("GetByID", "user_id", "Test_Category").Return(&domain.Category{
					ID:     "cat_12345678",
					Name:   "Test_Category",
					UserID: "user_id",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			testName: "actualizar categoría exitosamente",
			method:   "PATCH",
			path:     "/categories/Test_Category",
			body:     `{"name": "Test_Category", "new_name": "Updated Category", "user_id": "user_id"}`,
			userID:   "user_id",
			mockSetup: func(m *MockCategoryService) {
				m.On("Update", update.UpdateCategoryRequest{
					UserID:  "user_id",
					ID:      "Test_Category",
					NewName: "Updated Category",
				}).Return(&update.UpdateCategoryResponse{
					UserID: "user_id",
					Name:   "Updated Category",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			testName: "eliminar categoría exitosamente",
			method:   "DELETE",
			path:     "/categories/Test_Category",
			userID:   "user_id",
			mockSetup: func(m *MockCategoryService) {
				m.On("Delete", categories.DeleteCategoryRequest{
					ID:     "Test_Category",
					UserID: "user_id",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockService := new(MockCategoryService)
			// Para tests, usamos nil como gamificationHelper ya que no lo necesitamos
			var gamificationHelper *services.GamificationHelper = nil
			handler := NewHandler(mockService, gamificationHelper)
			tt.mockSetup(mockService)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Configurar la petición
			c.Request = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.userID != "" {
				c.Request.Header.Set("X-Caller-ID", tt.userID)
				c.Set("user_id", tt.userID)
			}

			// Configurar parámetros de ruta si es necesario
			if strings.Contains(tt.path, "/categories/") {
				parts := strings.Split(tt.path, "/")
				if len(parts) > 2 {
					c.Params = []gin.Param{
						{Key: "id", Value: parts[2]},
					}
				}
			}

			// Ejecutar el handler correspondiente
			switch tt.method {
			case "POST":
				handler.CreateCategory(c)
			case "GET":
				if strings.Contains(tt.path, "/categories/") {
					handler.GetCategory(c)
				} else {
					handler.ListCategories(c)
				}
			case "PATCH":
				handler.UpdateCategory(c)
			case "DELETE":
				handler.DeleteCategory(c)
			}

			// Verificar status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar mensaje de error si se espera uno
			if tt.expectedError != "" {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}

			// Verificar que se llamaron todos los mocks esperados
			mockService.AssertExpectations(t)
		})
	}
}
