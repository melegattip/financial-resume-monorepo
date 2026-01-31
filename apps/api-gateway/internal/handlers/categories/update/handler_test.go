package update

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/usecases/categories/update"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCategoryHandler(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		categoryID     string
		request        update.UpdateCategoryRequest
		mockSetup      func(*MockCategoryService)
		expectedStatus int
		expectedError  string
	}{
		{
			testName:   "actualizar categoría exitosamente",
			userID:     "user_id",
			categoryID: "cat_123_original",
			request: update.UpdateCategoryRequest{
				UserID:  "user_id",
				ID:      "cat_123_original",
				NewName: "cat_123_updated",
			},
			mockSetup: func(m *MockCategoryService) {
				response := &update.UpdateCategoryResponse{
					UserID: "user_id",
					ID:     "cat_123_original",
					Name:   "cat_123_updated",
				}
				m.On("Update", update.UpdateCategoryRequest{
					ID:      "cat_123_original",
					NewName: "cat_123_updated",
					UserID:  "user_id",
				}).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
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

			body, _ := json.Marshal(tt.request)
			c.Request = httptest.NewRequest("PATCH", "/", bytes.NewBuffer(body))
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
