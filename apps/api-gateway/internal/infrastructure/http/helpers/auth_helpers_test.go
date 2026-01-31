package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userIDValue    interface{}
		userIDExists   bool
		expectedResult string
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "Valid string user ID",
			userIDValue:    "12345",
			userIDExists:   true,
			expectedResult: "12345",
			expectError:    false,
		},
		{
			name:           "Valid uint user ID",
			userIDValue:    uint(67890),
			userIDExists:   true,
			expectedResult: "67890",
			expectError:    false,
		},
		{
			name:         "Missing user ID",
			userIDExists: false,
			expectError:  true,
			errorMessage: "Usuario no autenticado",
		},
		{
			name:         "Invalid user ID type",
			userIDValue:  123.45, // float64
			userIDExists: true,
			expectError:  true,
			errorMessage: "Formato de user_id inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			if tt.userIDExists {
				c.Set("user_id", tt.userIDValue)
			}

			// Execute
			result, err := GetUserIDFromContext(c)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetUserIDFromContextAsUint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userIDValue    interface{}
		userIDExists   bool
		expectedResult uint
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "Valid uint user ID",
			userIDValue:    uint(12345),
			userIDExists:   true,
			expectedResult: 12345,
			expectError:    false,
		},
		{
			name:           "Valid string user ID",
			userIDValue:    "67890",
			userIDExists:   true,
			expectedResult: 67890,
			expectError:    false,
		},
		{
			name:         "Missing user ID",
			userIDExists: false,
			expectError:  true,
			errorMessage: "User ID not found in context",
		},
		{
			name:         "Invalid string format",
			userIDValue:  "invalid",
			userIDExists: true,
			expectError:  true,
			errorMessage: "Invalid user ID format",
		},
		{
			name:         "Invalid user ID type",
			userIDValue:  123.45, // float64
			userIDExists: true,
			expectError:  true,
			errorMessage: "Invalid user ID type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			if tt.userIDExists {
				c.Set("user_id", tt.userIDValue)
			}

			// Execute
			result, err := GetUserIDFromContextAsUint(c)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				assert.Equal(t, uint(0), result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
