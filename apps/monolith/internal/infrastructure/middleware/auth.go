package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/ports"
)

// AuthMiddleware provides JWT-based authentication middleware for protected routes.
type AuthMiddleware struct {
	jwtService ports.JWTService
}

// NewAuthMiddleware creates a new auth middleware with the given JWT service.
func NewAuthMiddleware(jwtService ports.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// RequireAuth returns a middleware that enforces a valid JWT access token.
// On success it sets "user_id" (string), "user_email" (string), and "token" (string) in the Gin context.
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims, err := m.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("token", token)

		c.Next()
	}
}

// OptionalAuth returns a middleware that extracts user info from a valid token if present,
// but does not block the request if the token is missing or invalid.
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token != "" {
			claims, err := m.jwtService.ValidateAccessToken(token)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("token", token)
			}
		}
		c.Next()
	}
}

func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}

	return parts[1]
}
