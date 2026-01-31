package delete

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCategoryHandler(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		categoryID     string
		mockSetup      func(*MockCategoryService)
		expectedStatus int
		expectedError  string
	}{
		{
			testName:   "eliminar categoría exitosamente",
			userID:     "user_id",
			categoryID: "cat_123",
			mockSetup: func(m *MockCategoryService) {
				m.On("Delete", categories.DeleteCategoryRequest{
					UserID: "user_id",
					ID:     "cat_123",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			testName:   "error al eliminar categoría que no existe",
			userID:     "user_id",
			categoryID: "non_existent",
			mockSetup: func(m *MockCategoryService) {
				m.On("Delete", categories.DeleteCategoryRequest{
					UserID: "user_id",
					ID:     "non_existent",
				}).Return(errors.NewResourceNotFound("Category not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Category not found",
		},
		{
			testName:       "error al eliminar categoría con ID vacío",
			userID:         "user_id",
			categoryID:     "",
			mockSetup:      func(m *MockCategoryService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Category ID is required",
		},
		{
			testName:   "error al eliminar categoría sin user_id",
			userID:     "",
			categoryID: "cat_123",
			mockSetup: func(m *MockCategoryService) {
				m.On("Delete", categories.DeleteCategoryRequest{
					UserID: "",
					ID:     "cat_123",
				}).Return(errors.NewUnauthorizedRequest("user_id is required"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "user_id is required",
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

			c.Params = []gin.Param{
				{Key: "id", Value: tt.categoryID},
			}
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
