package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// PermissionChecker defines the interface for checking RBAC permissions.
type PermissionChecker interface {
	HasPermission(ctx context.Context, tenantID, role, permission string) (bool, error)
}

// PermissionMiddleware provides dynamic RBAC middleware for protected routes.
type PermissionMiddleware struct {
	checker PermissionChecker
	logger  zerolog.Logger
}

// NewPermissionMiddleware creates a new permission middleware.
func NewPermissionMiddleware(checker PermissionChecker, logger zerolog.Logger) *PermissionMiddleware {
	return &PermissionMiddleware{checker: checker, logger: logger}
}

// Require returns a Gin middleware that enforces the given permission.
// Must be used after RequireAuth() which sets tenant_id and role in the context.
func (m *PermissionMiddleware) Require(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		role := c.GetString("role")

		if tenantID == "" || role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing tenant context"})
			c.Abort()
			return
		}

		hasPermission, err := m.checker.HasPermission(c.Request.Context(), tenantID, role, permission)
		if err != nil {
			m.logger.Error().Err(err).
				Str("tenant_id", tenantID).
				Str("role", role).
				Str("permission", permission).
				Msg("permission check failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "permission check failed"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
