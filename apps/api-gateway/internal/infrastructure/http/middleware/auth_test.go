package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		userID         string
		callerID       string
		body           string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "error cuando el body es inválido para POST",
			method:         http.MethodPost,
			path:           "/test",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:           "error cuando falta user_id en el body para POST",
			method:         http.MethodPost,
			path:           "/test",
			userID:         "",
			callerID:       "user_id",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "user_id is required in body",
		},
		{
			name:           "error cuando falta X-Caller-ID en POST",
			method:         http.MethodPost,
			path:           "/test",
			userID:         "user_id",
			callerID:       "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "X-Caller-ID is required",
		},
		{
			name:           "error cuando no coinciden los IDs en POST",
			method:         http.MethodPost,
			path:           "/test",
			userID:         "user_id",
			callerID:       "different_id",
			expectedStatus: http.StatusForbidden,
			expectedError:  "User ID does not match X-Caller-ID",
		},
		{
			name:           "autenticación exitosa en POST",
			method:         http.MethodPost,
			path:           "/test",
			userID:         "user_id",
			callerID:       "user_id",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "error cuando falta user_id en el body para PATCH",
			method:         http.MethodPatch,
			path:           "/test",
			userID:         "",
			callerID:       "user_id",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "user_id is required in body",
		},
		{
			name:           "autenticación exitosa en PATCH",
			method:         http.MethodPatch,
			path:           "/test",
			userID:         "user_id",
			callerID:       "user_id",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "error cuando falta X-Caller-ID en GET",
			method:         http.MethodGet,
			path:           "/test",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "X-Caller-ID is required",
		},
		{
			name:           "autenticación exitosa en GET",
			method:         http.MethodGet,
			path:           "/test",
			callerID:       "user_id",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			// Agregar middleware de autenticación
			router.Use(AuthMiddleware())

			// Agregar rutas de prueba
			router.Handle(tt.method, "/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Crear request
			var req *http.Request
			if tt.method == http.MethodPost || tt.method == http.MethodPut || tt.method == http.MethodPatch {
				var body []byte
				if tt.body != "" {
					body = []byte(tt.body)
				} else {
					bodyMap := map[string]interface{}{}
					if tt.userID != "" {
						bodyMap["user_id"] = tt.userID
					}
					body, _ = json.Marshal(bodyMap)
				}
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			if tt.callerID != "" {
				req.Header.Set("X-Caller-ID", tt.callerID)
			}

			// Crear response recorder
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verificar status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verificar mensaje de error si se espera uno
			if tt.expectedError != "" {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}
		})
	}
}
