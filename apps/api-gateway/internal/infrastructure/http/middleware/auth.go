package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/melegattip/financial-resume-engine/internal/core/errors"
	httpUtil "github.com/melegattip/financial-resume-engine/internal/infrastructure/http"
)

// AuthMiddleware valida que el user_id exista y coincida con el X-Caller-ID
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		// Paths que deben ser ignorados por el middleware de autenticación
		ignoredPaths := []string{
			"/health",
			"/favicon.ico",
			"/robots.txt",
			"/manifest.json",
		}

		// Verificar si el path debe ser ignorado
		for _, ignoredPath := range ignoredPaths {
			if path == ignoredPath {
				c.Next()
				return
			}
		}

		// Endpoints que solo requieren X-Caller-ID (no user_id en body)
		headerOnlyPaths := []string{
			"/api/v1/categories",
			"/api/v1/reports",
			"/api/v1/recurring-transactions",
		}

		isHeaderOnlyPath := false
		for _, headerPath := range headerOnlyPaths {
			if strings.HasPrefix(path, headerPath) {
				isHeaderOnlyPath = true
				break
			}
		}

		// Solo validamos user_id en body para métodos POST, PUT y PATCH que NO sean header-only
		if (method == "POST" || method == "PUT" || method == "PATCH") && !isHeaderOnlyPath {
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err != nil {
				httpUtil.BadRequest(c, coreErrors.NewBadRequest("Invalid request body"))
				c.Abort()
				return
			}

			userID, exists := body["user_id"].(string)
			if !exists || userID == "" {
				httpUtil.BadRequest(c, coreErrors.NewBadRequest("user_id is required in body"))
				c.Abort()
				return
			}

			callerID := c.GetHeader("X-Caller-ID")
			if callerID == "" {
				httpUtil.Unauthorized(c, coreErrors.NewUnauthorizedRequest("X-Caller-ID is required"))
				c.Abort()
				return
			}

			if userID != callerID {
				httpUtil.SendError(c, http.StatusForbidden, "User ID does not match X-Caller-ID")
				c.Abort()
				return
			}

			// Restauramos el body para que pueda ser leído nuevamente
			jsonData, _ := json.Marshal(body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))
			c.Set("user_id", userID)
		} else {
			// Para GET, DELETE y endpoints header-only, solo validamos X-Caller-ID
			userID := c.GetHeader("X-Caller-ID")
			if userID == "" {
				httpUtil.Unauthorized(c, coreErrors.NewUnauthorizedRequest("X-Caller-ID is required"))
				c.Abort()
				return
			}
			c.Set("user_id", userID)
		}

		c.Next()
	}
}
